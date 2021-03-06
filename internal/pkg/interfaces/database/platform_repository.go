package database

import (
	"animar/v1/internal/pkg/domain"
	"animar/v1/internal/pkg/interfaces/database/queryset"

	"github.com/maru44/perr"
)

type PlatformRepository struct {
	SqlHandler
}

func (repo *PlatformRepository) FilterByAnime(animeId int) (platforms domain.TRelationPlatforms, err error) {
	rows, err := repo.Query(queryset.PlatformFilterByAnimeQuery, animeId)
	if err != nil {
		return platforms, perr.Wrap(err, perr.InternalServerErrorWithUrgency)
	}
	defer rows.Close()

	for rows.Next() {
		var p domain.TRelationPlatform
		err = rows.Scan(
			&p.PlatformId, &p.AnimeId, &p.LinkUrl,
			&p.DeliveryInterval, &p.FirstBroadcast,
			&p.CreatedAt, &p.UpdatedAt, &p.PlatName,
		)
		if err != nil {
			return platforms, perr.Wrap(err, perr.NotFound)
		}
		platforms = append(platforms, p)
	}
	return
}

func (repo *PlatformRepository) RegisterTarget(in domain.NotifiedTargetInput) (int, error) {
	exe, err := repo.Execute(queryset.RegisterNotifiedTargetQuery, in.SlackID, in.UserID)
	if err != nil {
		return 0, perr.Wrap(err, perr.InternalServerErrorWithUrgency)
	}

	rawInserted, err := exe.LastInsertId()
	if err != nil {
		return 0, perr.Wrap(err, perr.BadRequest)
	}
	return int(rawInserted), nil
}

func (repo *PlatformRepository) UpdateTarget(in domain.NotifiedTargetInput) (int, error) {
	exe, err := repo.Execute(queryset.UpdateNotifiedTargetQuery, in.SlackID, in.UserID)
	if err != nil {
		return 0, perr.Wrap(err, perr.InternalServerErrorWithUrgency)
	}

	rawAffected, err := exe.RowsAffected()
	if err != nil {
		return 0, perr.Wrap(err, perr.BadRequest)
	}
	return int(rawAffected), nil
}

func (repo *PlatformRepository) DeleteTarget(userId string) (int, error) {
	exe, err := repo.Execute(queryset.DeleteNotifiedTargetQuery, userId)
	if err != nil {
		return 0, perr.Wrap(err, perr.InternalServerErrorWithUrgency)
	}

	rawAffected, err := exe.RowsAffected()
	if err != nil {
		return 0, perr.Wrap(err, perr.BadRequest)
	}
	return int(rawAffected), nil
}

func (repo *PlatformRepository) GetUsersChannel(userId string) (*string, error) {
	rows, err := repo.Query(queryset.GetUsersSlackChannelQuery, userId)
	if err != nil {
		return nil, perr.Wrap(err, perr.InternalServerErrorWithUrgency)
	}

	rows.Next()
	var channelId *string
	if err := rows.Scan(&channelId); err != nil {
		return nil, perr.Wrap(err, perr.NotFound)
	}

	return channelId, nil
}
