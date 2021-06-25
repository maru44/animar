package infrastructure

import (
	"animar/v1/pkg/interfaces/apis"
	"encoding/json"
	"errors"
	"net/http"
)

type ApiResponse struct{}

func NewApiResponse() apis.ApiResponse {
	ret := new(ApiResponse)
	return ret
}

func (api *ApiResponse) Response(w http.ResponseWriter, status int, body map[string]interface{}) error {
	w.WriteHeader(status)
	if status == http.StatusOK {
		data, err := json.Marshal(body)
		if err != nil {
			return GetResponse(http.StatusInternalServerError, "")
		}
		w.Write(data)
	}
	return GetResponse(status, "")
}

func GetResponse(status int, message string) (err error) {
	if message != "" {
		err = errors.New(message)
	} else {
		switch status {
		case http.StatusNotFound:
			err = errors.New("Not Found")
		case http.StatusForbidden:
			err = errors.New("Forbidden")
		case http.StatusUnauthorized:
			err = errors.New("Unauthorized")
		case http.StatusBadRequest:
			err = errors.New("Bad Request")
		case http.StatusOK:
			err = nil
		case http.StatusAccepted:
			err = nil
		default:
			err = errors.New("Internal Server Error")
		}
	}
	return err
}
