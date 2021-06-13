package tools

import "net/http"

func Handle(handlers ...func(w http.ResponseWriter, r *http.Request) error) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		corsMiddleware(w, r)
		allowOptionsMiddleware(w, r)
		for _, handler := range handlers {
			if err := handler(w, r); err != nil {
				return
			}
		}
	}
}
