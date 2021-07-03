package infrastructure

import (
	"animar/v1/pkg/domain"
	"encoding/json"
	"net/http"
)

type Httphandler struct{}

func (h *Httphandler) response(w http.ResponseWriter, err error, body map[string]interface{}) error {
	status := h.getStatusCode(err)
	w.WriteHeader(status)
	if status == http.StatusOK {
		data, _ := json.Marshal(body)
		w.Write(data)
	}
	return err
}

func (h *Httphandler) getStatusCode(err error) int {
	if err == nil {
		return http.StatusOK
	}

	switch err {
	case domain.ErrInternalServerError:
		return http.StatusInternalServerError
	case domain.ErrNotFound:
		return http.StatusNotFound
	case domain.ErrForbidden:
		return http.StatusForbidden
	case domain.ErrUnauthorized:
		return http.StatusUnauthorized
	case domain.ErrBadRequest:
		return http.StatusBadRequest
	case domain.StatusCreated:
		return http.StatusCreated
	case domain.ErrUnknownType:
		return http.StatusUnsupportedMediaType
	default:
		return http.StatusInternalServerError
	}
}
