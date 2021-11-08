package infrastructure

import (
	"animar/v1/internal/pkg/domain"
	"net/http"
)

func Log(h http.Handler, lg *domain.LogA) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		lg.Logging(r)
		h.ServeHTTP(w, r)
	})
}
