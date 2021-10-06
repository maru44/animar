package controllers

import (
	"encoding/json"
	"net/http"
	"reflect"

	"github.com/maru44/perr"
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
		if perror, ok := perr.IsPerror(err); ok {
			mess = map[string]interface{}{"message": perror.Output().Error()}
		} else {
			mess = map[string]interface{}{"message": err.Error()}
		}

		data, _ := json.Marshal(mess)
		w.Write(data)
	}
	return err
}

func getStatusCode(err error, w http.ResponseWriter) int {
	if err == nil || reflect.ValueOf(err).IsNil() {
		return http.StatusOK
	}

	if perror, ok := perr.IsPerror(err); ok {
		switch perror.Output() {
		case perr.InternalServerError, perr.InternalServerErrorWithUrgency:
			return http.StatusInternalServerError
		case perr.NotFound:
			return http.StatusNotFound
		case perr.Forbidden:
			return http.StatusForbidden
		case perr.Unauthorized, perr.Expired, perr.InvalidToken:
			return http.StatusUnauthorized
		case perr.BadRequest:
			return http.StatusBadRequest
		case perr.Created:
			return http.StatusCreated
		case perr.UnsupportedMediaType:
			return http.StatusUnsupportedMediaType
		case perr.MethodNotAllowed:
			return http.StatusMethodNotAllowed
		case perr.UnsupportedMediaType:
			return http.StatusUnsupportedMediaType
		default:
			return http.StatusInternalServerError
		}
	}

	switch err {
	case perr.InternalServerError, perr.InternalServerErrorWithUrgency:
		return http.StatusInternalServerError
	case perr.NotFound:
		return http.StatusNotFound
	case perr.Forbidden:
		return http.StatusForbidden
	case perr.Unauthorized, perr.Expired, perr.InvalidToken:
		return http.StatusUnauthorized
	case perr.BadRequest:
		return http.StatusBadRequest
	case perr.Created:
		return http.StatusCreated
	case perr.UnsupportedMediaType:
		return http.StatusUnsupportedMediaType
	case perr.MethodNotAllowed:
		return http.StatusMethodNotAllowed
	case perr.UnsupportedMediaType:
		return http.StatusUnsupportedMediaType
	default:
		return http.StatusInternalServerError
	}
}
