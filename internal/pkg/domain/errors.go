package domain

import "errors"

var (
	ErrInternalServerError = errors.New("Internal Server Error")
	ErrNotFound            = errors.New("Not Found")
	ErrForbidden           = errors.New("Forbidden")
	ErrUnauthorized        = errors.New("Unauthorized")
	ErrBadRequest          = errors.New("Bad Request")
	ErrUnknownType         = errors.New("Unknow Type")
	ErrMethodNotAllowed    = errors.New("Method Not Allowed")
	ErrCsrfNotValid        = errors.New("Csrf Not Valid")
	StatusCreated          = errors.New("Created")
)
