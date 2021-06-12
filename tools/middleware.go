package tools

import "net/http"

// type Meshods interface {
// 	MethodMiddleware()
// }

func MethodMiddleware(next http.Handler, methods []string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for _, m := range methods {
			if m == r.Method {
				next.ServeHTTP(w, r)
			}
		}
		http.Error(w, http.StatusText(405), 405)
		return
	})
}

// func MethodMiddleware2(w http.ResponseWriter, r *http.Request, methods []string) error {
// 	for _, m := range methods {
// 		if m == r.Method {
// 			return nil
// 		}
// 	}
// 	return errors.New("METHODS")
// }

// func MethodMiddleware(w http.ResponseWriter, r *http.Request) error {
// 	for _, m := range methods {
// 		if m == r.Method {
// 			next.ServeHTTP(w, r)
// 		}
// 	}
// }
