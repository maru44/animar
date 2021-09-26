package controllers

import (
	"animar/v1/internal/pkg/domain"
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"
)

func response(w http.ResponseWriter, r *http.Request, err error, body map[string]interface{}) error {
	if myErr, ok := err.(domain.MyError); ok {
		err = myErr.ErrorForOutput()
	}

	status := getStatusCode(err, w)
	if status == http.StatusOK {
		data, err := json.Marshal(body)

		// if failed marshal
		if err != nil {
			err = domain.Errors{Inner: errors.Wrap(err, "Failed to json.Marshal"), Flag: domain.InternalServerError}
			go slackErrorLogging(r.Context(), err)
			w.WriteHeader(http.StatusInternalServerError)
			mess, _ := json.Marshal(map[string]interface{}{"message": domain.ErrInternalServerError.Error()})
			w.Write(mess)
			return err
		}
		w.Write(data)
	} else {
		go slackErrorLogging(r.Context(), err)
	}
	w.WriteHeader(status)
	return err
}

func getStatusCode(err error, w http.ResponseWriter) int {
	if err == nil {
		return http.StatusOK
	}

	mess := map[string]interface{}{"message": err.Error()}
	data, _ := json.Marshal(mess)
	w.Write(data)

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
