package psql

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/evgen1067/hw12_13_14_15_calendar/internal/common"
	"github.com/evgen1067/hw12_13_14_15_calendar/internal/config"
	"github.com/stretchr/testify/require"
)

func TestDatabase(t *testing.T) {
	cfg, err := config.Parse("../../../configs/local.json")
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	store := NewStorage()
	store.Connect(ctx, cfg)
	defer store.Close()

	e := common.Event{
		Title:       fmt.Sprintf("Title %v", 1),
		Description: fmt.Sprintf("Description %v", 1),
		DateStart:   time.Now().AddDate(0, 0, 1),
		DateEnd:     time.Now().AddDate(0, 0, 2),
	}

	t.Run("Test CRUD operation Event", func(t *testing.T) {
		// проверили, что событие создается
		firstID, err := store.Create(ctx, e)
		require.NoError(t, err)

		// проверили, что с той же датой не создается
		_, err = store.Create(ctx, e)
		require.Error(t, err)

		// поменяли дату события - теперь должен добавляться
		e.DateStart = e.DateStart.AddDate(0, 0, 1)

		secondID, err := store.Create(ctx, e)
		require.NoError(t, err)

		// првоерка об ошибке, что такого id нет при апдейте
		_, err = store.Update(ctx, -1, e)
		require.Error(t, err)
		require.ErrorIs(t, common.ErrNotFound, err)

		// должна быть ошибка, тк дата есть у 2 события
		_, err = store.Update(ctx, firstID, e)
		require.Error(t, err)

		// ошибки быть не должно
		e.Description = "TEST"
		_, err = store.Update(ctx, secondID, e)
		require.NoError(t, err)

		eventsList, err := store.WeekList(ctx, e.DateStart.AddDate(0, 0, -3))
		require.NoError(t, err)
		require.Equal(t, 2, len(eventsList))

		// првоерка об ошибке, что такого id нет удалении
		_, err = store.Delete(ctx, -1)
		require.Error(t, err)

		// ошибок быть не должно (мы создали эти события)
		_, err = store.Delete(ctx, firstID)
		require.NoError(t, err)
		_, err = store.Delete(ctx, secondID)
		require.NoError(t, err)
	})

	t.Run("Test clearing old events", func(t *testing.T) {
		e := common.Event{
			Title:       fmt.Sprintf("Title %v", 1),
			Description: fmt.Sprintf("Description %v", 1),
			DateStart:   time.Now().AddDate(-2, 0, 0),
			DateEnd:     time.Now().AddDate(-1, -6, 0),
		}

		// проверили, что событие создается
		_, err := store.Create(ctx, e)
		require.NoError(t, err)

		// проверили, что появилось в БД
		eventsList, err := store.WeekList(ctx, e.DateStart.AddDate(0, 0, -3))
		require.NoError(t, err)
		require.Equal(t, 1, len(eventsList))

		// запускаем чистку
		err = store.ClearOldEvents(ctx)
		require.NoError(t, err)

		// проверили, что появилось в БД
		eventsList, err = store.WeekList(ctx, e.DateStart.AddDate(0, 0, -3))
		require.NoError(t, err)
		require.Equal(t, 0, len(eventsList))
	})

	t.Run("Test Scheduler List", func(t *testing.T) {
		e := common.Event{
			Title:       fmt.Sprintf("Title %v", 1),
			Description: fmt.Sprintf("Description %v", 1),
			DateStart:   time.Now().AddDate(0, 0, 1),
			DateEnd:     time.Now().AddDate(0, 0, 1),
			NotifyIn:    48, // предупредить за 48 часов
		}

		// проверили, что событие создается
		firstID, err := store.Create(ctx, e)
		require.NoError(t, err)

		// получаем список уведомлений
		notifyList, err := store.SchedulerList(ctx)
		require.NoError(t, err)
		require.Equal(t, 1, len(notifyList))

		// ошибок быть не должно (мы создали эти события)
		_, err = store.Delete(ctx, firstID)
		require.NoError(t, err)
	})
}
