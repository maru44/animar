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

func slackErrorLogging(ctx context.Context, in error, additional ...string) (err error) {
	cli := http.Client{}
	var access string
	if a, ok := ctx.Value(accessContextKey).(accessContext); ok {
		access = fmt.Sprintf("URL: %s, METHOD: %s", a.URL, a.Method)
	}

	var message string
	if myErr, ok := in.(domain.MyError); ok {
		errLevel := myErr.Level()
		message = url.QueryEscape(fmt.Sprintf(fmt.Sprintf("%s\n%s\n\nStackTrace: %s\n%+v", errLevel, access, in.Error(), myErr.Traces())))
	} else {
		message = url.QueryEscape(fmt.Sprintf(fmt.Sprintf("%s\n%s\n\nStackTrace:\n%+v", domain.ErrExternal, access, in)))
	}
	// message = url.QueryEscape(fmt.Sprintf(fmt.Sprintf("%s\n%s\n\nStackTrace:\n%+v", errLevel, access, in)))
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
