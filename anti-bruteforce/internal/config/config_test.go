package config

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLogger(t *testing.T) {
	t.Run("Test Config with valid path", func(t *testing.T) {
		cfg, err := Parse("../../configs/local.json")
		require.NoError(t, err)

		require.Equal(t, "0.0.0.0", cfg.HTTP.Host)
		require.Equal(t, "8000", cfg.HTTP.Port)

		require.Equal(t, "0.0.0.0", cfg.DB.Host)
		require.Equal(t, "9000", cfg.DB.Port)
		require.Equal(t, "go_user", cfg.DB.User)
		require.Equal(t, "go_password", cfg.DB.Password)
		require.Equal(t, "bruteforce_db", cfg.DB.Database)
		require.Equal(t, "disable", cfg.DB.SSLMode)

		require.Equal(t, 10, cfg.Limitations.Login)
		require.Equal(t, 100, cfg.Limitations.Password)
		require.Equal(t, 1000, cfg.Limitations.IP)
	})

	t.Run("Test Config with invalid path", func(t *testing.T) {
		_, err := Parse("fail.json")
		require.Error(t, err)
	})
}
