package domain

import (
	// "errors"
	"github.com/pkg/errors"
)

type ErrorLevel string

const (
	/* internal */

	ErrorLevelAlertInternal ErrorLevel = "ALERT"
	ErrorLevelInternal      ErrorLevel = "INTERNAL"

	// external

	ErrorLevelExternal ErrorLevel = "EXTERNAL"
)

// for response
var (
	SuccessCreated = errors.New("Created")
	/*   for trace   */

	ErrorBadRequest       = errors.New("Bad Request")
	ErrorDataNotFound     = errors.New("Not Found")
	ErrorTokenInValid     = errors.New("Invalid Token")
	ErrorTokenIsExpired   = errors.New("Token is expired")
	ErrorUnauthorized     = errors.New("Unauthorized")
	ErrorForbidden        = errors.New("Forbidden")
	ErrorMethodNotAllowed = errors.New("Method not allowed")
	ErrorCsrfInValid      = errors.New("Invalid csrf token")
	ErrorUnknownType      = errors.New("Unknown type")

	// internal not emergency

	ErrorInternalServer = errors.New("Internal Server Error")

	// internal emergency

	ErrorMySQLConncetion    = errors.New("Internal Server Error")
	ErrorFirebaseConnection = errors.New("Internal Server Error")
	ErrorS3Connection       = errors.New("Internal Server Error")
	ErrorHttpConnection     = errors.New("Internal Server Error")
)

func GetErrorLevel(e error) ErrorLevel {
	switch e {
	case ErrorMySQLConncetion, ErrorFirebaseConnection, ErrorS3Connection, ErrorHttpConnection:
		return ErrorLevelAlertInternal
	case ErrorInternalServer:
		return ErrorLevelInternal
	default:
		return ErrorLevelExternal
	}
}
