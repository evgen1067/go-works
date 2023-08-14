package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/evgen1067/anti-bruteforce/internal/bucket"
	"github.com/evgen1067/anti-bruteforce/internal/common"
	"github.com/evgen1067/anti-bruteforce/internal/config"
	"github.com/evgen1067/anti-bruteforce/internal/logger"
	"github.com/evgen1067/anti-bruteforce/internal/repository/psql"
	"github.com/evgen1067/anti-bruteforce/internal/rest"
	"github.com/evgen1067/anti-bruteforce/internal/service"
	"github.com/stretchr/testify/require"
)

var (
	ctx    context.Context
	cancel context.CancelFunc
)

func TestApp(t *testing.T) {
	ctx, cancel = context.WithCancel(context.Background())
	defer cancel()
	cfg, _ := config.Parse("../configs/local.json")
	zLog, _ := logger.NewLogger(cfg)
	db := psql.NewRepo(cfg)
	leakyBucket := bucket.NewLeakyBucket(cfg)
	go func() {
		leakyBucket.Repeat(ctx)
	}()
	services := service.NewServices(ctx, db, leakyBucket, zLog)
	rest.NewServer(services, cfg)
	_ = db.Connect(ctx)
	defer db.Close()

	t.Run("Authorization test with a leaky bucket", func(t *testing.T) {
		// проверяем, что после 10 запросов с логином заблокируемся
		for i := 0; i < 10; i++ {
			code, body := tryAuthorize(authReq)
			require.Equal(t, http.StatusOK, code)
			require.Equal(t, authSuccess, body)
		}
		// должна быть ошибка, тк превышен лимит запросов с логином
		code, body := tryAuthorize(authReq)
		require.Equal(t, http.StatusTooManyRequests, code)
		require.Equal(t, authFail, body)

		// проверяем, что мы можем сбросить ведро и снова отправлять запросы
		code = tryReset()
		require.Equal(t, http.StatusOK, code)

		// после сброса ведра все должно быть успешно
		code, body = tryAuthorize(authReq)
		require.Equal(t, http.StatusOK, code)
		require.Equal(t, authSuccess, body)
	})

	t.Run("Add to blacklist", func(t *testing.T) {
		code := tryAddToList(address, "blacklist")
		require.Equal(t, http.StatusCreated, code)

		code = tryAddToList(address, "blacklist")
		require.Equal(t, http.StatusConflict, code)
	})

	t.Run("Add to whitelist", func(t *testing.T) {
		code := tryAddToList(address, "whitelist")
		require.Equal(t, http.StatusCreated, code)

		code = tryAddToList(address, "whitelist")
		require.Equal(t, http.StatusConflict, code)
	})

	t.Run("Add or delete to non-existent list", func(t *testing.T) {
		code := tryAddToList(address, "test")
		require.Equal(t, http.StatusNotFound, code)

		code = tryDeleteFromList(address, "test")
		require.Equal(t, http.StatusNotFound, code)
	})

	t.Run("Remove from blacklist", func(t *testing.T) {
		code := tryDeleteFromList(address, "blacklist")
		require.Equal(t, http.StatusAccepted, code)

		code = tryDeleteFromList(address, "blacklist")
		require.Equal(t, http.StatusNotFound, code)
	})

	t.Run("Remove from whitelist", func(t *testing.T) {
		code := tryDeleteFromList(address, "whitelist")
		require.Equal(t, http.StatusAccepted, code)

		code = tryDeleteFromList(address, "whitelist")
		require.Equal(t, http.StatusNotFound, code)
	})

	t.Run("authorization test if the ip is in the black list", func(t *testing.T) {
		code := tryReset()
		require.Equal(t, http.StatusOK, code)

		// добавляем ip в blacklist
		code = tryAddToList(address, "blacklist")
		require.Equal(t, http.StatusCreated, code)

		// авторизация заблокирована
		code, body := tryAuthorize(authReq)
		require.Equal(t, http.StatusTooManyRequests, code)
		require.Equal(t, authFail, body)

		code = tryDeleteFromList(address, "blacklist")
		require.Equal(t, http.StatusAccepted, code)
	})

	t.Run("authorization test if the ip is in the white list", func(t *testing.T) {
		code := tryReset()
		require.Equal(t, http.StatusOK, code)

		// добавляем ip в whitelist
		code = tryAddToList(address, "whitelist")
		require.Equal(t, http.StatusCreated, code)

		// авторизация разрешена
		for i := 0; i < 20; i++ {
			code, body := tryAuthorize(authReq)
			require.Equal(t, http.StatusOK, code)
			require.Equal(t, authSuccess, body)
		}

		code = tryDeleteFromList(address, "whitelist")
		require.Equal(t, http.StatusAccepted, code)
	})

	err := os.RemoveAll("logs")
	require.NoError(t, err)
}

func tryAuthorize(r common.APIAuthRequest) (int, *common.APIAuthResponse) {
	body := new(bytes.Buffer)
	_ = json.NewEncoder(body).Encode(&r)

	code, resp := makeHTTPRequest(http.MethodPost, "http://localhost/auth", body)

	response := new(common.APIAuthResponse)
	all, _ := io.ReadAll(resp)
	_ = response.UnmarshalJSON(all)

	return code, response
}

func tryReset() int {
	code, _ := makeHTTPRequest(http.MethodPost, "http://localhost/reset/bucket", new(bytes.Buffer))
	return code
}

func tryAddToList(address, list string) int {
	addr := common.APIListRequest{Address: address}
	body := new(bytes.Buffer)
	_ = json.NewEncoder(body).Encode(&addr)

	code, _ := makeHTTPRequest(http.MethodPost, "http://localhost/list/"+list, body)
	return code
}

func tryDeleteFromList(address, list string) int {
	addr := common.APIListRequest{Address: address}
	body := new(bytes.Buffer)
	_ = json.NewEncoder(body).Encode(&addr)

	code, _ := makeHTTPRequest(http.MethodDelete, "http://localhost/list/"+list, body)
	return code
}

func makeHTTPRequest(method, url string, body *bytes.Buffer) (int, *bytes.Buffer) {
	recorder := httptest.NewRecorder()
	router := rest.Router()

	req, _ := http.NewRequestWithContext(ctx, method, url, body)
	router.ServeHTTP(recorder, req)

	return recorder.Code, recorder.Body
}
