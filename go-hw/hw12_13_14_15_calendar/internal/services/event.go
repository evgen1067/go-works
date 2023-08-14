package services

import (
	"context"
	"time"

	"github.com/evgen1067/hw12_13_14_15_calendar/internal/common"
	"github.com/evgen1067/hw12_13_14_15_calendar/internal/storage"
)

type EventsService struct {
	ctx context.Context
	db  storage.EventsStorage
}

func NewEventsService(ctx context.Context, db storage.EventsStorage) *EventsService {
	return &EventsService{
		ctx: ctx,
		db:  db,
	}
}

func (e *EventsService) Create(event common.Event) (common.EventID, error) {
	return e.db.Create(e.ctx, event)
}

func (e *EventsService) Update(id common.EventID, event common.Event) (common.EventID, error) {
	return e.db.Update(e.ctx, id, event)
}

func (e *EventsService) Delete(id common.EventID) (common.EventID, error) {
	return e.db.Delete(e.ctx, id)
}

func (e *EventsService) DayList(startDate time.Time) ([]common.Event, error) {
	return e.db.DayList(e.ctx, startDate)
}

func (e *EventsService) WeekList(startDate time.Time) ([]common.Event, error) {
	return e.db.WeekList(e.ctx, startDate)
}

func (e *EventsService) MonthList(startDate time.Time) ([]common.Event, error) {
	return e.db.MonthList(e.ctx, startDate)
}
