package apis

import "net/http"

type ApiResponse interface {
	Response(http.ResponseWriter, int, map[string]interface{}) error
}
