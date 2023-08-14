package memory

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/evgen1067/hw12_13_14_15_calendar/internal/common"
	"github.com/stretchr/testify/require"
)

func TestStorage(t *testing.T) {
	var err error
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	store := NewStorage()
	events := make([]common.Event, 0)
	for i := 0; i < 10; i++ {
		e := common.Event{
			Title:       fmt.Sprintf("Title %v", i),
			Description: fmt.Sprintf("Description %v", i),
			DateStart:   time.Now().AddDate(0, 0, i*1),
			DateEnd:     time.Now().AddDate(0, 0, i*2),
		}
		events = append(events, e)
	}

	t.Run("Test create Event", func(t *testing.T) {
		for _, e := range events {
			_, err = store.Create(ctx, e)
			require.NoError(t, err)
		}
	})

	t.Run("Test Event creation error (date is busy)", func(t *testing.T) {
		_, err = store.Create(ctx, events[2])
		require.Error(t, err)
		require.ErrorIs(t, common.ErrDateBusy, err)

		_, err = store.Update(ctx, 1, events[1])
		require.NoError(t, err)
	})

	t.Run("Test Event update error (date is busy)", func(t *testing.T) {
		_, err = store.Update(ctx, 3, events[2])
		require.Error(t, err)
		require.ErrorIs(t, common.ErrDateBusy, err)
	})

	t.Run("Test of successful deletion", func(t *testing.T) {
		_, err = store.Delete(ctx, 0)
		require.NoError(t, err)
	})

	t.Run("Test of not successful deletion", func(t *testing.T) {
		_, err = store.Delete(ctx, 0)
		require.Error(t, err)
	})

	t.Run("Test of getting an event list", func(t *testing.T) {
		var periodEvents []common.Event

		periodEvents, err = store.DayList(ctx, time.Now().AddDate(0, 0, 1))
		require.NoError(t, err)
		require.Equal(t, 2, len(periodEvents))

		periodEvents, err = store.DayList(ctx, time.Now().AddDate(0, 0, -30))
		require.NoError(t, err)
		require.Equal(t, 0, len(periodEvents))

		periodEvents, err = store.WeekList(ctx, time.Now().AddDate(0, 0, 1))
		require.NoError(t, err)
		require.Equal(t, 8, len(periodEvents))

		periodEvents, err = store.MonthList(ctx, time.Now().AddDate(0, 0, 1))
		require.NoError(t, err)
		require.Equal(t, 9, len(periodEvents))
	})
}
