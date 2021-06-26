package apis

import "net/http"

type ApiResponse interface {
	Response(http.ResponseWriter, error, map[string]interface{}) error
}
