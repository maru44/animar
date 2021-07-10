package httphandle

import (
	"context"
	"io"
	"mime/multipart"
)

type Httphandler interface {
	Handle(string, Handler)
	HandlerFunc(string, func(ResponseWriter, Request))
	ListenAndServe(string, Handler) error
	SetCookie(ResponseWriter, Cookie)
	CallClient() Client
}

//

type Handler interface{}

type ResponseWriter interface {
	Write([]byte) (int, error)
	WriteHeader(int)
}

type Client interface {
	Do(Request) (Response, error)
	Get(string) (Response, error)
}

//

type Response interface {
}

type Request interface {
	Cookie(string) (Cookie, error)
	FormValue(string) string
	FormFile(string) (multipart.File, *multipart.FileHeader, error)
	Write(io.Writer) error
	Context() context.Context
}

//

type Cookie interface {
	String() string
}
