package controllers

import (
	"context"
	"fmt"
	"net/http"
)

func setUserToContext(r *http.Request, userId string) *http.Request {
	ctx := context.WithValue(
		r.Context(),
		USER_ID,
		userId,
	)
	r = r.WithContext(ctx)
	fmt.Println(ctx.Value(USER_ID).(string))
	return r
}
