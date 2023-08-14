package bucket

import (
	"strconv"
	"testing"

	"github.com/evgen1067/anti-bruteforce/internal/common"
	"github.com/evgen1067/anti-bruteforce/internal/config"
	"github.com/stretchr/testify/require"
)

func TestBucket(t *testing.T) {
	cfg, err := config.Parse("../../configs/local.json")
	require.NoError(t, err)

	lb := NewLeakyBucket(cfg)

	req := common.APIAuthRequest{
		Login:    "test_login",
		Password: "test_password",
		IP:       "127.0.0.1/25",
	}

	t.Run("check add login", func(t *testing.T) {
		for i := 0; i < 10; i++ {
			flag := lb.Add(req)
			require.Equal(t, true, flag)
		}
		flag := lb.Add(req)
		require.Equal(t, false, flag)
	})

	t.Run("check reset bucket", func(t *testing.T) {
		flag := lb.Add(req)
		require.Equal(t, false, flag)

		lb.ResetBucket()

		flag = lb.Add(req)
		require.Equal(t, true, flag)
	})

	lb.ResetBucket()

	t.Run("check add pass", func(t *testing.T) {
		for i := 0; i < 100; i++ {
			req.Login = strconv.Itoa(i)
			flag := lb.Add(req)
			require.Equal(t, true, flag)
		}
		flag := lb.Add(req)
		require.Equal(t, false, flag)
	})

	lb.ResetBucket()

	t.Run("check add ip", func(t *testing.T) {
		for i := 0; i < 1000; i++ {
			req.Login = strconv.Itoa(i)
			req.Password = strconv.Itoa(i)
			flag := lb.Add(req)
			require.Equal(t, true, flag)
		}
		flag := lb.Add(req)
		require.Equal(t, false, flag)
	})

	lb.ResetBucket()
}
