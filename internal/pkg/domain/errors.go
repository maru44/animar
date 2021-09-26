package domain

import (
	// "errors"
	"bytes"
	"fmt"
	"runtime"
	"strings"

	"github.com/pkg/errors"
)

// for response
var (
	ErrInternalServerError = errors.New("Internal Server Error")
	ErrNotFound            = errors.New("Not Found")
	ErrForbidden           = errors.New("Forbidden")
	ErrUnauthorized        = errors.New("Unauthorized")
	ErrTokenIsExpired      = errors.New("Token is expired")
	ErrBadRequest          = errors.New("Bad Request")
	ErrUnknownType         = errors.New("Unknow Type")
	ErrMethodNotAllowed    = errors.New("Method Not Allowed")
	ErrCsrfNotValid        = errors.New("Csrf Not Valid")
	StatusCreated          = errors.New("Created")
)

const (
	/*  external  */

	ExternalServerError uint32 = 1 << iota
	DataNotFoundError
	TokenIsInvalidError
	TokenIsExpiredError
	UnauthorizedError
	ForbiddenError
	MethodNotAllowedError
	CsrfNotValidError
	UnknownTypeError

	/*  Internal not emergency  */
	InternalServerError

	/*  Internal Emergency  */

	MySqlConnectionError
	FirebaseConnectionError
	S3ConnectionError
	HttpConnectionError
)

type (
	MyError interface {
		ErrorForOutput() error
		GetFlag() uint32
		Traces() string
	}

	// @TODO: add caller(for stack trace)
	Errors struct {
		Inner  error // stores the error returned by external dependencies
		Flag   uint32
		text   string
		traces []StackTraceFrame
	}

	StackTraceFrame struct {
		File           string
		Line           int
		Name           string
		ProgramCounter uintptr // origin data
	}
)

func (e Errors) Error() string {
	if e.Inner != nil {
		return e.Inner.Error()
	} else if e.text != "" {
		return e.text
	} else {
		return ErrInternalServerError.Error()
	}
}

func (e Errors) ErrorForOutput() error {
	switch e.Flag {
	case ExternalServerError, CsrfNotValidError:
		return ErrBadRequest
	case UnauthorizedError, TokenIsInvalidError:
		return ErrUnauthorized
	case TokenIsExpiredError:
		return ErrTokenIsExpired
	case ForbiddenError:
		return ErrForbidden
	case DataNotFoundError:
		return ErrNotFound
	case MethodNotAllowedError:
		return ErrMethodNotAllowed
	case UnknownTypeError:
		return ErrUnknownType
	case MySqlConnectionError, FirebaseConnectionError, S3ConnectionError, HttpConnectionError:
		return ErrInternalServerError
	default:
		return ErrInternalServerError
	}
}

func (e Errors) GetFlag() uint32 {
	return e.Flag
}

func (e Errors) Traces() string {
	var buf bytes.Buffer
	for _, fr := range e.traces {
		fmt.Fprintf(&buf, "%s: %d===>%v\n", fr.File, fr.Line, fr.Name)
	}
	return buf.String()
}

func NewError(text string, flag uint32) *Errors {
	return &Errors{
		Flag:   flag,
		text:   text,
		traces: NewTrace(callers()),
	}
}

func NewWrapError(e error, flag uint32) *Errors {
	return &Errors{
		Flag:   flag,
		Inner:  e,
		traces: NewTrace(callers()),
	}
}

// https://github.com/pkg/errors/blob/816c9085562cd7ee03e7f8188a1cfd942858cded/stack.go#L133
func callers() []uintptr {
	const depth = 32
	var pcs [depth]uintptr
	n := runtime.Callers(3, pcs[:])
	return pcs[0 : n-2]
}

func NewTrace(pcs []uintptr) []StackTraceFrame {
	traces := []StackTraceFrame{}

	for _, pc := range pcs {
		trace := StackTraceFrame{ProgramCounter: pc}
		fn := runtime.FuncForPC(pc)
		if fn == nil {
			return traces
		}
		trace.Name = trimPkgName(fn)
		trace.File, trace.Line = fn.FileLine(pc - 1)
		traces = append(traces, trace)
	}
	return traces
}

func trimPkgName(fn *runtime.Func) string {
	name := fn.Name()
	if ld := strings.LastIndex(name, "."); ld >= 0 {
		name = name[ld+1:]
	}
	return name
}
