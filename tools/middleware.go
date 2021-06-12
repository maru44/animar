package tools

// type Meshods interface {
// 	MethodMiddleware()
// }

// func (methods []string) MethodMiddleware(next http.Handler, methods []string) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		for _, m := range methods {
// 			if m == r.Method {
// 				next.ServeHTTP(w, r)
// 			}
// 		}
// 		return
// 	})
// }

// func MethodMiddleware(w http.ResponseWriter, r *http.Request) error {
// 	for _, m := range methods {
// 		if m == r.Method {
// 			next.ServeHTTP(w, r)
// 		}
// 	}
// }
