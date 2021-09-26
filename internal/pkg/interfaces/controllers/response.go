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

		mess := map[string]interface{}{"message": err.Error()}
		data, _ := json.Marshal(mess)
		w.Write(data)
	}
	return err
}

func getStatusCode(err error, w http.ResponseWriter) int {
	if err == nil {
		return http.StatusOK
	}

	switch err {
	case domain.ErrorInternalServer, domain.ErrorMySQLConncetion, domain.ErrorFirebaseConnection, domain.ErrorS3Connection, domain.ErrorHttpConnection:
		return http.StatusInternalServerError
	case domain.ErrorDataNotFound:
		return http.StatusNotFound
	case domain.ErrorForbidden:
		return http.StatusForbidden
	case domain.ErrorUnauthorized, domain.ErrorTokenInValid, domain.ErrorTokenIsExpired:
		return http.StatusUnauthorized
	case domain.ErrorBadRequest:
		return http.StatusBadRequest
	case domain.SuccessCreated:
		return http.StatusCreated
	case domain.ErrorUnknownType:
		return http.StatusUnsupportedMediaType
	case domain.ErrorMethodNotAllowed:
		return http.StatusMethodNotAllowed
	case domain.ErrorCsrfInValid:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}
