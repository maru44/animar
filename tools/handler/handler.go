package handler

import (
	"animar/v1/tools/middleware"
	"net/http"
)

func Handle(handlers ...func(w http.ResponseWriter, r *http.Request) error) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		middleware.CorsMiddleware(w, r)
		middleware.AllowOptionsMiddleware(w, r)
		for _, handler := range handlers {
			if err := handler(w, r); err != nil {
				return
			}
		}
	}
}
