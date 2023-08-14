package services

import (
	"context"
	"time"

	"github.com/evgen1067/hw12_13_14_15_calendar/internal/common"
	"github.com/evgen1067/hw12_13_14_15_calendar/internal/logger"
	"github.com/evgen1067/hw12_13_14_15_calendar/internal/storage"
	"go.uber.org/zap"
)

type Events interface {
	Create(event common.Event) (common.EventID, error)
	Update(id common.EventID, event common.Event) (common.EventID, error)
	Delete(id common.EventID) (common.EventID, error)
	DayList(startDate time.Time) ([]common.Event, error)
	WeekList(startDate time.Time) ([]common.Event, error)
	MonthList(startDate time.Time) ([]common.Event, error)
}

type Logger interface {
	Error(msg string, fields ...zap.Field)
	Warn(msg string, fields ...zap.Field)
	Debug(msg string, fields ...zap.Field)
	Info(msg string, fields ...zap.Field)
}

type Services struct {
	Events
	Logger
}

func NewServices(ctx context.Context, db storage.EventsStorage, l *logger.Logger) *Services {
	events := NewEventsService(ctx, db)
	logg := NewLogger(l)
	return &Services{
		Events: events,
		Logger: logg,
	}
}
