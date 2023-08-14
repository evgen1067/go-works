package psql

import (
	"context"
	"testing"
	"time"

	"github.com/evgen1067/anti-bruteforce/internal/common"
	"github.com/evgen1067/anti-bruteforce/internal/config"
	"github.com/stretchr/testify/require"
)

func TestDatabase(t *testing.T) {
	t.Run("Test of database operations", func(t *testing.T) {
		cfg, err := config.Parse("../../../configs/local.json")
		require.NoError(t, err)
		repo := NewRepo(cfg)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		err = repo.Connect(ctx)
		require.NoError(t, err)
		defer repo.Close()
		addresses := []string{"192.1.2.0/12", "183.3.1.0/1", "118.1.5.0/22", "106.1.1.0/15"}
		for _, v := range addresses {
			err := repo.AddToBlacklist(ctx, v)
			require.NoError(t, err)

			err = repo.AddToWhitelist(ctx, v)
			require.NoError(t, err)
		}

		for _, v := range addresses {
			err := repo.AddToBlacklist(ctx, v)
			require.Error(t, err)
			require.ErrorIs(t, err, common.ErrIPExists)

			err = repo.AddToWhitelist(ctx, v)
			require.Error(t, err)
			require.ErrorIs(t, err, common.ErrIPExists)
		}

		for _, v := range addresses {
			flag, err := repo.ExistsInBlacklist(ctx, v)
			require.NoError(t, err)
			require.Equal(t, true, flag)

			flag, err = repo.ExistsInWhitelist(ctx, v)
			require.NoError(t, err)
			require.Equal(t, true, flag)
		}

		for _, v := range addresses {
			err := repo.DeleteFromBlacklist(ctx, v)
			require.NoError(t, err)

			err = repo.DeleteFromWhitelist(ctx, v)
			require.NoError(t, err)
		}

		for _, v := range addresses {
			flag, err := repo.ExistsInBlacklist(ctx, v)
			require.NoError(t, err)
			require.Equal(t, false, flag)

			flag, err = repo.ExistsInWhitelist(ctx, v)
			require.NoError(t, err)
			require.Equal(t, false, flag)
		}

		for _, v := range addresses {
			err := repo.DeleteFromBlacklist(ctx, v)
			require.Error(t, err)
			require.ErrorIs(t, err, common.ErrIPNotExists)

			err = repo.DeleteFromWhitelist(ctx, v)
			require.Error(t, err)
			require.ErrorIs(t, err, common.ErrIPNotExists)
		}
	})
}
