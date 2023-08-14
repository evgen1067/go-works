package storage

import (
	"context"
	"time"

	"github.com/evgen1067/hw12_13_14_15_calendar/internal/common"
	"github.com/evgen1067/hw12_13_14_15_calendar/internal/config"
)

type Period string

type DBStorage interface {
	Connect(ctx context.Context, cfg *config.Config) error
	Close() error
	EventsStorage
	SchedulerStorage
}

type EventsStorage interface {
	Create(ctx context.Context, event common.Event) (common.EventID, error)
	Update(ctx context.Context, id common.EventID, event common.Event) (common.EventID, error)
	Delete(ctx context.Context, id common.EventID) (common.EventID, error)
	DayList(ctx context.Context, startDate time.Time) ([]common.Event, error)
	WeekList(ctx context.Context, startDate time.Time) ([]common.Event, error)
	MonthList(ctx context.Context, startDate time.Time) ([]common.Event, error)
}

type SchedulerStorage interface {
	SchedulerList(ctx context.Context) ([]common.Notice, error)
	ClearOldEvents(ctx context.Context) error
}
