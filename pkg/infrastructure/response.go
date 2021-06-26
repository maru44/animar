package infrastructure

import (
	"animar/v1/pkg/domain"
	"animar/v1/pkg/interfaces/apis"
	"encoding/json"
	"net/http"
)

type ApiResponse struct{}

func NewApiResponse() apis.ApiResponse {
	ret := new(ApiResponse)
	return ret
}

func (api *ApiResponse) Response(w http.ResponseWriter, err error, body map[string]interface{}) error {
	status := getStatusCode(err)
	w.WriteHeader(status)
	if status == http.StatusOK {
		data, _ := json.Marshal(body)
		w.Write(data)
	}
	return err
}

func getStatusCode(err error) int {
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

// func GetResponse(status int, message string) (err error) {
// 	if message != "" {
// 		err = errors.New(message)
// 	} else {
// 		switch status {
// 		case http.StatusNotFound:
// 			err = errors.New("Not Found")
// 		case http.StatusForbidden:
// 			err = errors.New("Forbidden")
// 		case http.StatusUnauthorized:
// 			err = errors.New("Unauthorized")
// 		case http.StatusBadRequest:
// 			err = errors.New("Bad Request")
// 		case http.StatusOK:
// 			err = nil
// 		case http.StatusAccepted:
// 			err = nil
// 		default:
// 			err = errors.New("Internal Server Error")
// 		}
// 	}
// 	return err
// }
