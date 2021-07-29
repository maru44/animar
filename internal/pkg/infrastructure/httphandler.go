package infrastructure

import (
	"animar/v1/internal/pkg/domain"
	"animar/v1/internal/pkg/interfaces/httphandle"
	"bytes"
	"context"
	"io"
	"mime/multipart"
	"net/http"
)

type HttpHandler struct{}

type HttpClient struct {
	Client *http.Client
}

type HttpResponse struct {
	Response *http.Response
}

type HttpRequest struct {
	Request *http.Request
}

type HttpResponseWriter struct {
	ResponseWriter http.ResponseWriter
}

/*******************
	instance
*******************/

func NewHttpClient() httphandle.Client {
	cl := &http.Client{}
	client := new(HttpClient)
	client.Client = cl
	return client
}

func NewHttpRequest(method string, url string, b *bytes.Buffer) httphandle.Request {
	req, err := http.NewRequest(method, url, b)
	if err != nil {
		lg := domain.NewErrorLog(err.Error(), "")
		lg.Logging()
	}
	request := new(HttpRequest)
	request.Request = req
	return request
}

func NewHttpResponseWriter() httphandle.ResponseWriter {
	rw := new(http.ResponseWriter)
	writer := new(HttpResponseWriter)
	writer.ResponseWriter = *rw
	return writer
}

// func NewHttpHandler() httphandle.Httphandler {
// 	httpHandler := new(HttpHandler)
// 	return httpHandler
// }

/*******************
	client
*******************/

func (cl *HttpClient) Do(req *http.Request) (httphandle.Response, error) {
	resp, err := cl.Client.Do(req)
	response := new(HttpResponse)
	response.Response = resp
	return response, err
}

/*******************
	request
*******************/

func (req *HttpRequest) Cookie(key string) (httphandle.Cookie, error) {
	return req.Request.Cookie(key)
}

func (req *HttpRequest) Context() context.Context {
	return req.Request.Context()
}

func (req *HttpRequest) FormValue(key string) string {
	return req.Request.FormValue(key)
}

func (req *HttpRequest) FormFile(key string) (multipart.File, *multipart.FileHeader, error) {
	return req.Request.FormFile(key)
}

func (req *HttpRequest) Write(iow io.Writer) error {
	return req.Request.Write(iow)
}

/*******************
	response writer
*******************/

func (rw HttpResponseWriter) Write(b []byte) (int, error) {
	return rw.Write(b)
}

func (rw HttpResponseWriter) WriteHeader(status int) {
	rw.WriteHeader(status)
}

/*******************
   http logger
*******************/

func Log(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rAddress := r.RemoteAddr
		method := r.Method
		path := r.URL.Path
		lg := domain.NewAccessLog(rAddress, method, path)
		lg.Logging()
		// log.Printf("Access: %s %s %s", rAddress, method, path)
		h.ServeHTTP(w, r)
	})
}
