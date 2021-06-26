package domain

import "errors"

var (
	ErrInternalServerError = errors.New("Internal Server Error")
	ErrNotFound            = errors.New("Not Found")
	ErrForbidden           = errors.New("Forbidden")
	ErrUnauthorized        = errors.New("Unauthorized")
	ErrBadRequest          = errors.New("Bad Request")
	ErrUnknownType         = errors.New("Unknow Type")
	StatusCreated          = errors.New("Created")
)
