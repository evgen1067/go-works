package grpc

import (
	"github.com/evgen1067/hw12_13_14_15_calendar/api"
	"github.com/evgen1067/hw12_13_14_15_calendar/internal/common"
	"github.com/golang/protobuf/ptypes/timestamp"
	"time"
)

func TransformEventToPb(e common.Event) *api.Event {
	return &api.Event{
		Id:          uint64(e.ID),
		Title:       e.Title,
		Description: e.Description,
		DateStart:   &timestamp.Timestamp{Seconds: e.DateStart.Unix(), Nanos: int32(e.DateStart.Nanosecond())},
		DateEnd:     &timestamp.Timestamp{Seconds: e.DateEnd.Unix(), Nanos: int32(e.DateEnd.Nanosecond())},
		NotifyIn:    uint64(e.NotifyIn),
		OwnerId:     uint64(e.OwnerID),
	}
}

func TransformPbToEvent(e *api.Event) common.Event {
	return common.Event{
		ID:          common.EventID(e.Id),
		Title:       e.Title,
		Description: e.Description,
		DateStart:   time.Unix(e.DateStart.Seconds, int64(e.DateStart.Nanos)),
		DateEnd:     time.Unix(e.DateEnd.Seconds, int64(e.DateEnd.Nanos)),
		NotifyIn:    int64(e.NotifyIn),
		OwnerID:     int64(e.OwnerId),
	}
}
