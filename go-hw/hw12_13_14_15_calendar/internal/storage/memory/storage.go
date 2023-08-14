package memory

import (
	"context"
	"sync"
	"time"

	"github.com/evgen1067/hw12_13_14_15_calendar/internal/common"
	"github.com/evgen1067/hw12_13_14_15_calendar/internal/storage"
)

type Storage struct {
	mu        sync.RWMutex
	Events    map[common.EventID]common.Event
	Increment common.EventID
	length    int
}

func NewStorage() storage.EventsStorage {
	return &Storage{
		Events: make(map[common.EventID]common.Event),
	}
}

func (store *Storage) CheckDate(ctx context.Context, event common.Event) error {
	for _, e := range store.Events {
		if e.DateStart.Format("2/Jan/2006:15:04") == event.DateStart.Format("2/Jan/2006:15:04") &&
			e.ID != event.ID && e.OwnerID == event.OwnerID {
			return common.ErrDateBusy
		}
	}
	return nil
}

func (store *Storage) Create(ctx context.Context, event common.Event) (common.EventID, error) {
	store.mu.Lock()
	defer store.mu.Unlock()
	event.ID = store.Increment
	err := store.CheckDate(ctx, event)
	if err != nil {
		return event.ID, err
	}
	store.Events[event.ID] = event
	store.Increment++
	store.length++
	return event.ID, nil
}

func (store *Storage) Update(ctx context.Context, id common.EventID, event common.Event) (common.EventID, error) {
	store.mu.Lock()
	defer store.mu.Unlock()

	_, ok := store.Events[id]
	if !ok {
		return id, common.ErrNotFound
	}

	event.ID = id

	err := store.CheckDate(ctx, event)
	if err != nil {
		return event.ID, err
	}
	store.Events[id] = event
	return event.ID, nil
}

func (store *Storage) Delete(ctx context.Context, id common.EventID) (common.EventID, error) {
	store.mu.Lock()
	defer store.mu.Unlock()

	_, ok := store.Events[id]
	if !ok {
		return id, common.ErrNotFound
	}
	delete(store.Events, id)
	store.length--
	return id, nil
}

func (store *Storage) PeriodList(
	ctx context.Context,
	startPeriod time.Time,
	period storage.Period,
) ([]common.Event, error) {
	var endPeriod time.Time
	switch period {
	case "Day":
		endPeriod = startPeriod.AddDate(0, 0, 1)
	case "Week":
		endPeriod = startPeriod.AddDate(0, 0, 7)
	case "Month":
		endPeriod = startPeriod.AddDate(0, 1, 0)
	}
	var events []common.Event
	for _, e := range store.Events {
		if e.DateEnd.After(startPeriod) && endPeriod.After(e.DateStart) {
			events = append(events, e)
		}
	}
	return events, nil
}

func (store *Storage) DayList(ctx context.Context, startDate time.Time) ([]common.Event, error) {
	store.mu.RLock()
	defer store.mu.RUnlock()
	period := storage.Period("Day")
	return store.PeriodList(ctx, startDate, period)
}

func (store *Storage) WeekList(ctx context.Context, startDate time.Time) ([]common.Event, error) {
	store.mu.RLock()
	defer store.mu.RUnlock()
	period := storage.Period("Week")
	return store.PeriodList(ctx, startDate, period)
}

func (store *Storage) MonthList(ctx context.Context, startDate time.Time) ([]common.Event, error) {
	store.mu.RLock()
	defer store.mu.RUnlock()
	period := storage.Period("Month")
	return store.PeriodList(ctx, startDate, period)
}
