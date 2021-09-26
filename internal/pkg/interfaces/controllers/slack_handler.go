package controllers

import (
	"animar/v1/configs"
	"animar/v1/internal/pkg/domain"
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
)

type ErrLevel string

const (
	// internal

	ErrAlert    ErrLevel = "ALERT"
	ErrInternal ErrLevel = "INTERNAL ERROR"

	// external

	ErrExternal ErrLevel = "EXTERNAL ERROR"
)

func slackErrorLogging(ctx context.Context, in error, additional ...string) (err error) {
	cli := http.Client{}
	var access string
	if a, ok := ctx.Value(accessContextKey).(accessContext); ok {
		access = fmt.Sprintf("URL: %s, METHOD: %s", a.URL, a.Method)
	}

	errLevel := ErrExternal
	if myErr, ok := err.(domain.MyError); ok {
		flag := myErr.GetFlag()
		if flag > 512 {
			errLevel = ErrAlert
		} else if flag == domain.InternalServerError {
			errLevel = ErrInternal
		}
	}
	message := url.QueryEscape(fmt.Sprintf(fmt.Sprintf("%s\n%s\n\nStackTrace:\n%+v", errLevel, access, err)))
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("https://slack.com/api/chat.postMessage?channel=%s&text=%s&pretty=1", configs.SlackChannelId, message), nil)
	if err != nil {
		log.Println(err)
		return
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", configs.SlackBotToken))
	_, err = cli.Do(req)
	if err != nil {
		log.Println(err)
	}
	return err
}
