package httphandle

import (
	"context"
	"io"
	"mime/multipart"
	"net/http"
)

type Httphandler interface {
	Handle(string, Handler)
	HandlerFunc(string, func(ResponseWriter, Request))
	ListenAndServe(string, Handler) error
	SetCookie(ResponseWriter, Cookie)
}

type Handler interface{}

type ResponseWriter interface {
	Write([]byte) (int, error)
	WriteHeader(int)
}

type Client interface {
	Do(*http.Request) (Response, error)
	// Get(string) (Response, error)
}

type Response interface {
	// Write(w io.Writer) error
}

type Header interface {
	Set(string, string) []string
	Write(io.Writer) error
}

type Request interface {
	// Cookies() []Cookie
	Cookie(string) (Cookie, error)
	FormValue(string) string
	FormFile(string) (multipart.File, *multipart.FileHeader, error)
	Write(io.Writer) error
	Context() context.Context
}

type Cookie interface {
	String() string
}
