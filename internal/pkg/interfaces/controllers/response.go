package controllers

import (
	"animar/v1/internal/pkg/domain"
	"encoding/json"
	"net/http"
)

func response(w http.ResponseWriter, err error, body map[string]interface{}) error {
	status := getStatusCode(err, w)
	w.WriteHeader(status)
	if status == http.StatusOK {
		data, _ := json.Marshal(body)
		w.Write(data)
	}
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
