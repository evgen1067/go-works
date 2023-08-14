package grpc

import (
	"context"
	"github.com/evgen1067/hw12_13_14_15_calendar/api"
	"github.com/evgen1067/hw12_13_14_15_calendar/internal/common"
	"time"
)

func (s *Server) Create(ctx context.Context, req *api.CreateRequest) (*api.CreateResponse, error) {
	e := TransformPbToEvent(req.Event)

	eID, err := s.services.Create(e)
	if err != nil {
		return nil, err
	}

	return &api.CreateResponse{
		Id: uint64(eID),
	}, nil
}

func (s *Server) Update(ctx context.Context, req *api.UpdateRequest) (*api.UpdateResponse, error) {
	id := common.EventID(req.Id)
	e := TransformPbToEvent(req.Event)

	eID, err := s.services.Update(id, e)
	if err != nil {
		return nil, err
	}

	return &api.UpdateResponse{
		Id: uint64(eID),
	}, nil
}

func (s *Server) Delete(ctx context.Context, req *api.DeleteRequest) (*api.DeleteResponse, error) {
	id := common.EventID(req.Id)

	eID, err := s.services.Delete(id)
	if err != nil {
		return nil, err
	}

	return &api.DeleteResponse{
		Id: uint64(eID),
	}, nil
}

func (s *Server) DayList(ctx context.Context, req *api.ListRequest) (*api.ListResponse, error) {
	return PeriodList(ctx, req, s.services.DayList)
}

func (s *Server) WeekList(ctx context.Context, req *api.ListRequest) (*api.ListResponse, error) {
	return PeriodList(ctx, req, s.services.WeekList)
}

func (s *Server) MonthList(ctx context.Context, req *api.ListRequest) (*api.ListResponse, error) {
	return PeriodList(ctx, req, s.services.MonthList)
}

func PeriodList(ctx context.Context,
	req *api.ListRequest,
	fn func(startDate time.Time) ([]common.Event, error),
) (*api.ListResponse, error) {
	startDate := time.Unix(req.Date.Seconds, int64(req.Date.Nanos))
	events, err := fn(startDate)
	if err != nil {
		return nil, err
	}
	periodList := make([]*api.Event, 0)
	for _, val := range events {
		periodList = append(periodList, TransformEventToPb(val))
	}
	return &api.ListResponse{Event: periodList}, nil
}
