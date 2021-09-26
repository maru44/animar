package controllers

import (
	"animar/v1/internal/pkg/domain"
	"encoding/json"
	"net/http"
)

func response(w http.ResponseWriter, r *http.Request, err error, body map[string]interface{}) error {
	status := getStatusCode(err, w)
	if status == http.StatusOK {
		data, err := json.Marshal(body)

		// if failed marshal
		if err != nil {
			go slackErrorLogging(r.Context(), err)
			w.WriteHeader(http.StatusInternalServerError)
			mess, _ := json.Marshal(map[string]interface{}{"message": err.Error()})
			w.Write(mess)
			return err
		}
		w.WriteHeader(status)
		w.Write(data)
	} else {
		w.WriteHeader(status)
		go slackErrorLogging(r.Context(), err)

		var mess map[string]interface{}
		if myErr, ok := err.(domain.MyError); ok {
			mess = map[string]interface{}{"message": myErr.ErrorForOutput().Error()}
		} else {
			mess = map[string]interface{}{"message": err.Error()}
		}

		data, _ := json.Marshal(mess)
		w.Write(data)
	}
	return err
}

func getStatusCode(err error, w http.ResponseWriter) int {
	if err == nil {
		return http.StatusOK
	}

	if myErr, ok := err.(domain.MyError); ok {
		err = myErr.ErrorForOutput()
	}

	switch err {
	case domain.ErrInternalServerError:
		return http.StatusInternalServerError
	case domain.ErrNotFound:
		return http.StatusNotFound
	case domain.ErrForbidden:
		return http.StatusForbidden
	case domain.ErrUnauthorized, domain.ErrTokenIsExpired:
		return http.StatusUnauthorized
	case domain.ErrBadRequest:
		return http.StatusBadRequest
	case domain.StatusCreated:
		return http.StatusCreated
	case domain.ErrUnknownType:
		return http.StatusUnsupportedMediaType
	case domain.ErrMethodNotAllowed:
		return http.StatusMethodNotAllowed
	case domain.ErrCsrfNotValid:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}
