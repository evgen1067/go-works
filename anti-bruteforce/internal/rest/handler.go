package rest

import (
	"errors"
	"io"
	"net/http"

	"github.com/evgen1067/anti-bruteforce/internal/common"
	"github.com/gorilla/mux"
)

func CustomNotFoundHandler(w http.ResponseWriter, r *http.Request) {
	ex := common.APIException{
		Code:    http.StatusNotFound,
		Message: "The page you are looking for does not exist.",
	}
	WriteException(w, &ex)
}

func ResetBucket(w http.ResponseWriter, r *http.Request) {
	s.ResetBucket()
	ex := common.APIException{
		Code:    http.StatusOK,
		Message: "Leaky buckets dropped.",
	}
	WriteException(w, &ex)
}

func Auth(w http.ResponseWriter, r *http.Request) {
	var req common.APIAuthRequest
	bytes, err := io.ReadAll(r.Body)
	if err != nil {
		WriteException(w, &common.APIException{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}
	err = req.UnmarshalJSON(bytes)
	if err != nil {
		WriteException(w, &common.APIException{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}
	status := s.Authorize(req)
	response := &common.APIAuthResponse{
		Ok: status,
	}
	jsonResponse, err := response.MarshalJSON()
	if err != nil {
		return
	}
	if status {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusTooManyRequests)
	}
	_, err = w.Write(jsonResponse)
	if err != nil {
		return
	}
}

func Add(w http.ResponseWriter, r *http.Request) {
	WorkWithList(w, r, http.MethodPost)
}

func Delete(w http.ResponseWriter, r *http.Request) {
	WorkWithList(w, r, http.MethodDelete)
}

func WorkWithList(w http.ResponseWriter, r *http.Request, method string) {
	listName := common.TableName(mux.Vars(r)["value"])
	var req common.APIListRequest
	bytes, err := io.ReadAll(r.Body)
	if err != nil {
		WriteException(w, &common.APIException{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}
	err = req.UnmarshalJSON(bytes)
	if err != nil {
		WriteException(w, &common.APIException{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}
	switch listName {
	case common.Blacklist:
		if method == http.MethodPost {
			err = s.AddToBlacklist(req.Address)
		} else {
			err = s.DeleteFromBlacklist(req.Address)
		}
	case common.Whitelist:
		if method == http.MethodPost {
			err = s.AddToWhitelist(req.Address)
		} else {
			err = s.DeleteFromWhitelist(req.Address)
		}
	default:
		WriteException(w, &common.APIException{
			Code:    http.StatusNotFound,
			Message: "The page you are looking for does not exist.",
		})
		return
	}
	switch {
	case errors.Is(err, common.ErrIPExists):
		WriteException(w, &common.APIException{
			Code:    http.StatusConflict,
			Message: err.Error(),
		})
		return
	case errors.Is(err, common.ErrIPNotExists):
		WriteException(w, &common.APIException{
			Code:    http.StatusNotFound,
			Message: err.Error(),
		})
		return
	case err != nil:
		WriteException(w, &common.APIException{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}

	if method == http.MethodPost {
		w.WriteHeader(http.StatusCreated)
	} else {
		w.WriteHeader(http.StatusAccepted)
	}
}

func WriteException(w http.ResponseWriter, ex *common.APIException) {
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
