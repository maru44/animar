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

func slackErrorLogging(ctx context.Context, in error, additional ...string) (err error) {
	cli := http.Client{}
	var access string
	if a, ok := ctx.Value(accessContextKey).(accessContext); ok {
		access = fmt.Sprintf("URL: %s, METHOD: %s", a.URL, a.Method)
	}

	message := url.QueryEscape(fmt.Sprintf(fmt.Sprintf("%s\n%s\n\nStackTrace:\n%+v", domain.GetErrorLevel(err), access, in)))
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
