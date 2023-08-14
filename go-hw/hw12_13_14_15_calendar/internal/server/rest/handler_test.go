package rest

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/evgen1067/hw12_13_14_15_calendar/internal/common"
	"github.com/evgen1067/hw12_13_14_15_calendar/internal/config"
	"github.com/evgen1067/hw12_13_14_15_calendar/internal/logger"
	"github.com/evgen1067/hw12_13_14_15_calendar/internal/services"
	"github.com/evgen1067/hw12_13_14_15_calendar/internal/storage/memory"
	"github.com/stretchr/testify/require"
)

var (
	ctx    context.Context
	cancel context.CancelFunc
)

func TestCreateUpdateDelete(t *testing.T) {
	ctx, cancel = context.WithCancel(context.Background())
	defer cancel()

	cfg, _ := config.Parse("../../../configs/local.json")
	zLog, _ := logger.NewLogger(cfg)
	store := memory.NewStorage()

	services := services.NewServices(ctx, store, zLog)

	NewServer(services, cfg)

	t.Run("Test Create, Update, Delete Event", func(t *testing.T) {
		e := common.Event{
			Title:       "Title",
			Description: "Desc",
			DateStart:   time.Now(),
			DateEnd:     time.Now().AddDate(2, 0, 1),
			NotifyIn:    1,
			OwnerID:     1,
		}
		body, err := eventToBytes(e)
		require.NoError(t, err)

		// не должно быть ошибок, событие новое
		code, _ := makeHTTPRequest(http.MethodPost, "/events/new", body)
		require.Equal(t, http.StatusCreated, code)

		body, err = eventToBytes(e)
		require.NoError(t, err)

		// должна быть ошибка о том, что дата занята
		code, _ = makeHTTPRequest(http.MethodPost, "/events/new", body)
		require.Equal(t, http.StatusConflict, code)

		e.DateStart = e.DateStart.AddDate(1, 0, 0)

		body, err = eventToBytes(e)
		require.NoError(t, err)

		code, body = makeHTTPRequest(http.MethodPost, "/events/new", body)
		require.Equal(t, http.StatusCreated, code)

		resp, err := bytesToIDResponse(body)
		require.NoError(t, err)

		body, err = eventToBytes(e)
		require.NoError(t, err)

		code, _ = makeHTTPRequest(http.MethodPut, "/events/-1", body)
		require.Equal(t, http.StatusNotFound, code)

		body, err = eventToBytes(e)
		require.NoError(t, err)

		code, _ = makeHTTPRequest(http.MethodPut, fmt.Sprintf("/events/%d", resp.EventID), body)
		require.Equal(t, http.StatusOK, code)

		// удаляем то, что создали (2 события с id - 0, 1)
		code, _ = makeHTTPRequest(http.MethodDelete, "/events/0", new(bytes.Buffer))
		require.Equal(t, http.StatusAccepted, code)

		code, _ = makeHTTPRequest(http.MethodDelete, "/events/1", new(bytes.Buffer))
		require.Equal(t, http.StatusAccepted, code)

		code, _ = makeHTTPRequest(http.MethodDelete, "/events/1", new(bytes.Buffer))
		require.Equal(t, http.StatusInternalServerError, code)
	})
	err := os.RemoveAll("logs")
	require.NoError(t, err)
}

func TestLists(t *testing.T) {
	ctx, cancel = context.WithCancel(context.Background())
	defer cancel()

	cfg, _ := config.Parse("../../../configs/local.json")
	zLog, _ := logger.NewLogger(cfg)
	store := memory.NewStorage()

	services := services.NewServices(ctx, store, zLog)

	NewServer(services, cfg)

	t.Run("Test invalid period", func(t *testing.T) {
		code, _ := makeHTTPRequest(http.MethodGet, fmt.Sprintf("/events/list/%v", "fail"), new(bytes.Buffer))
		require.Equal(t, http.StatusBadRequest, code)
	})
	t.Run("Test no start date specified", func(t *testing.T) {
		code, _ := makeHTTPRequest(http.MethodGet, fmt.Sprintf("/events/list/%v", "day"), new(bytes.Buffer))
		require.Equal(t, http.StatusBadRequest, code)
	})
	t.Run("Test date parsing error", func(t *testing.T) {
		code, _ := makeHTTPRequest(http.MethodGet,
			fmt.Sprintf("/events/list/%v?%v", "day", "date=2001-12-18"), new(bytes.Buffer))
		require.Equal(t, http.StatusBadRequest, code)
	})
	t.Run("Test success get list", func(t *testing.T) {
		e := common.Event{
			Title:       "Title",
			Description: "Desc",
			DateStart:   time.Now(),
			DateEnd:     time.Now().AddDate(2, 0, 1),
			NotifyIn:    1,
			OwnerID:     1,
		}
		body, err := eventToBytes(e)
		require.NoError(t, err)

		// не должно быть ошибок, событие новое
		code, _ := makeHTTPRequest(http.MethodPost, "/events/new", body)
		require.Equal(t, http.StatusCreated, code)

		code, body = makeHTTPRequest(http.MethodGet,
			fmt.Sprintf("/events/list/%v?date=%v", "week",
				time.Now().Add(-5*time.Second).Format("2006-01-02T15:04:05")), new(bytes.Buffer))
		require.Equal(t, http.StatusOK, code)
		resp, err := bytesToListResponse(body)
		require.NoError(t, err)
		require.Equal(t, 1, len(resp.Events))
	})
	err := os.RemoveAll("logs")
	require.NoError(t, err)
}

func eventToBytes(event common.Event) (*bytes.Buffer, error) {
	b := new(bytes.Buffer)
	err := json.NewEncoder(b).Encode(&event)
	return b, err
}

func bytesToIDResponse(body *bytes.Buffer) (common.ResponseEventID, error) {
	var resp common.ResponseEventID
	err := json.NewDecoder(body).Decode(&resp)
	return resp, err
}

func bytesToListResponse(body *bytes.Buffer) (common.ResponseEventList, error) {
	var resp common.ResponseEventList
	err := json.NewDecoder(body).Decode(&resp)
	return resp, err
}

func makeHTTPRequest(method, url string, body *bytes.Buffer) (int, *bytes.Buffer) {
	recorder := httptest.NewRecorder()
	router := Router()

	req, _ := http.NewRequestWithContext(ctx, method, url, body)
	router.ServeHTTP(recorder, req)

	return recorder.Code, recorder.Body
}
