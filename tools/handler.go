package tools

import "net/http"

func Handle(handlers ...func(w http.ResponseWriter, r *http.Request) error) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		for _, handler := range handlers {
			if err := handler(w, r); err != nil {
				return
			}
		}
	}
}
