package batch

import (
	"animar/v1/configs"
	"animar/v1/internal/pkg/domain"
	"animar/v1/internal/pkg/interfaces/database"
	"animar/v1/internal/pkg/usecase"
	"fmt"
	"net/http"
	"net/url"

	"github.com/maru44/perr"
)

type PlatformBatch struct {
	interactor domain.PlatformBatchInteractor
}

func NewPlatformBatch(sql database.SqlHandler) *PlatformBatch {
	return &PlatformBatch{
		interactor: usecase.NewPlatformBatchInteractor(
			&database.PlatformBatchRepository{
				SqlHandler: sql,
			},
		),
	}
}

func (pb *PlatformBatch) SendBatch() error {
	// targets, err := pb.interactor.FilterNotificationTarget() // @TODO enable
	// if err != nil {
	// 	return perr.Wrap(err, perr.BadRequest)
	// }

	broadCasts, err := pb.interactor.TargetNotificationBroadcast()
	if err != nil {
		return perr.Wrap(err, perr.BadRequest)
	}

	unEscMess := pb.interactor.MakeSlackMessage(broadCasts)
	message := url.QueryEscape(unEscMess)

	cli := http.Client{}
	// @TODO multiplize
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("https://slack.com/api/chat.postMessage?channel=%s&text=%s&pretty=1", configs.SlackChannelId, message), nil)
	if err != nil {
		return perr.Wrap(err, perr.BadRequest)
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", configs.SlackBotToken))
	_, err = cli.Do(req)
	if err != nil {
		return perr.Wrap(err, perr.BadRequest)
	}
	return nil
}
