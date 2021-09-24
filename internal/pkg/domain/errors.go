package domain

import "errors"

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
	Error struct {
		Inner error // stores the error returned by external dependencies
		Flag  uint32
		text  string
	}
)

func (e Error) Errors() string {
	if e.Inner != nil {
		return e.Inner.Error()
	} else if e.text != "" {
		return e.text
	} else {
		return ErrInternalServerError.Error()
	}
}

func (e Error) ErrorForOutput() error {
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

func NewError(text string, flag uint32) *Error {
	return &Error{
		Flag: flag,
		text: text,
	}
}
