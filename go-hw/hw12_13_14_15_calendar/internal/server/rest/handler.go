package rest

import (
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/evgen1067/hw12_13_14_15_calendar/internal/common"
	"github.com/evgen1067/hw12_13_14_15_calendar/internal/storage"
	"github.com/gorilla/mux"
)

func CustomNotFoundHandler(w http.ResponseWriter, r *http.Request) {
	ex := common.Exception{
		Code:    http.StatusNotFound,
		Message: "The page you are looking for does not exist.",
	}
	WriteException(w, ex)
}

func HelloWorld(w http.ResponseWriter, r *http.Request) {
	ex := common.Exception{
		Code:    http.StatusOK,
		Message: "Hello, World!",
	}
	WriteException(w, ex)
}

func CreateEvent(w http.ResponseWriter, r *http.Request) {
	var event common.Event
	bytes, err := io.ReadAll(r.Body)
	if err != nil {
		WriteException(w, common.Exception{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}
	err = event.UnmarshalJSON(bytes)
	if err != nil {
		WriteException(w, common.Exception{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}
	var eventID common.EventID
	eventID, err = s.Create(event)
	switch {
	case errors.Is(err, common.ErrDateBusy):
		WriteException(w, common.Exception{
			Code:    http.StatusConflict,
			Message: err.Error(),
		})
		return
	case err != nil:
		WriteException(w, common.Exception{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	default:
		WriteEventIDResponse(w, common.ResponseEventID{
			Code:    http.StatusCreated,
			EventID: eventID,
		})
		return
	}
}

func UpdateEvent(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	eventID := common.EventID(id)
	if err != nil {
		WriteException(w, common.Exception{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}
	var event common.Event
	bytes, err := io.ReadAll(r.Body)
	if err != nil {
		WriteException(w, common.Exception{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}
	err = event.UnmarshalJSON(bytes)
	if err != nil {
		WriteException(w, common.Exception{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}
	eventID, err = s.Update(eventID, event)
	switch {
	case errors.Is(err, common.ErrNotFound):
		WriteException(w, common.Exception{
			Code:    http.StatusNotFound,
			Message: err.Error(),
		})
		return
	case errors.Is(err, common.ErrDateBusy):
		WriteException(w, common.Exception{
			Code:    http.StatusConflict,
			Message: err.Error(),
		})
		return
	case err != nil:
		WriteException(w, common.Exception{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	default:
		WriteEventIDResponse(w, common.ResponseEventID{
			Code:    http.StatusOK,
			EventID: eventID,
		})
	}
}

func DeleteEvent(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		WriteException(w, common.Exception{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}
	eventID := common.EventID(id)
	eventID, err = s.Delete(eventID)
	if err != nil {
		WriteException(w, common.Exception{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}
	WriteEventIDResponse(w, common.ResponseEventID{
		Code:    http.StatusAccepted,
		EventID: eventID,
	})
}

func EventList(w http.ResponseWriter, r *http.Request) {
	_period := strings.ToLower(mux.Vars(r)["period"])
	var period storage.Period
	switch _period {
	case "day":
		period = storage.Period("Day")
	case "week":
		period = storage.Period("Week")
	case "month":
		period = storage.Period("Month")
	default:
		WriteException(w, common.Exception{
			Code:    http.StatusBadRequest,
			Message: "The period specified in the request is not supported by the service.",
		})
		return
	}
	dateParam := r.URL.Query().Get("date")
	if dateParam == "" {
		WriteException(w, common.Exception{
			Code:    http.StatusBadRequest,
			Message: "No start date specified.",
		})
		return
	}
	startDate, err := time.Parse("2006-01-02T15:04:05", dateParam)
	if err != nil {
		WriteException(w, common.Exception{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}
	var events []common.Event
	switch period {
	case "Day":
		events, err = s.DayList(startDate)
	case "Week":
		events, err = s.WeekList(startDate)
	case "Month":
		events, err = s.MonthList(startDate)
	}
	if err != nil {
		WriteException(w, common.Exception{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}
	WriteEventListResponse(w, common.ResponseEventList{
		Code:   http.StatusOK,
		Events: events,
	})
}

func WriteEventIDResponse(w http.ResponseWriter, r common.ResponseEventID) {
	w.WriteHeader(r.Code)
	jsonResponse, err := r.MarshalJSON()
	if err != nil {
		return
	}
	_, err = w.Write(jsonResponse)
	if err != nil {
		return
	}
}

func WriteEventListResponse(w http.ResponseWriter, r common.ResponseEventList) {
	w.WriteHeader(r.Code)
	jsonResponse, err := r.MarshalJSON()
	if err != nil {
		return
	}
	_, err = w.Write(jsonResponse)
	if err != nil {
		return
	}
}

func WriteException(w http.ResponseWriter, ex common.Exception) {
	w.WriteHeader(ex.Code)
	jsonResponse, err := ex.MarshalJSON()
	if err != nil {
		return
	}
	_, err = w.Write(jsonResponse)
	if err != nil {
		return
	}
}
