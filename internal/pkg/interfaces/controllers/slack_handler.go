package controllers

import (
	"animar/v1/configs"
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/maru44/perr"
)

func slackErrorLogging(ctx context.Context, in error, additional ...string) (err error) {
	cli := http.Client{}
	var access string
	if a, ok := ctx.Value(accessContextKey).(accessContext); ok {
		access = fmt.Sprintf("URL: %s, METHOD: %s", a.URL, a.Method)
	}

	var message string
	if perror, ok := in.(perr.Perror); ok {
		errLevel := perror.Level()
		message = url.QueryEscape(fmt.Sprintf(fmt.Sprintf("%s\n%s\n\nStackTrace: %s\n%+v", errLevel, access, in.Error(), perror.Traces())))
	} else {
		message = url.QueryEscape(fmt.Sprintf(fmt.Sprintf("%s\n%s\n\nStackTrace:\n%+v", perr.InternalServerError, access, in)))
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
