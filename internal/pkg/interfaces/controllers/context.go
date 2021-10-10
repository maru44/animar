package controllers

import (
	"context"
	"net/http"
)

type (
	ContextKey string

	accessContext struct {
		Method string
		URL    string
	}
)

const (
	USER_ID          ContextKey = "userId"
	accessContextKey ContextKey = "access"
)

func setUserToContext(r *http.Request, userId string) *http.Request {
	ctx := context.WithValue(
		r.Context(),
		USER_ID,
		userId,
	)
	return r.WithContext(ctx)
}

func setAccessData(r *http.Request) *http.Request {
	ctx := context.WithValue(
		r.Context(),
		accessContextKey,
		accessContext{
			Method: r.Method,
			URL:    r.URL.Path,
		},
	)
	return r.WithContext(ctx)
}
