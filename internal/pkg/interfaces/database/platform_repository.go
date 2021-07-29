package database

import (
	"animar/v1/internal/pkg/domain"
)

type PlatformRepository struct {
	SqlHandler
}

func (repo *PlatformRepository) FilterByAnime(animeId int) (platforms domain.TRelationPlatforms, err error) {
	rows, err := repo.Query(
		"Select relation_anime_platform.*, platforms.plat_name FROM relation_anime_platform "+
			"LEFT JOIN platforms ON relation_anime_platform.platform_id = platforms.id "+
			"WHERE anime_id = ?", animeId,
	)
	defer rows.Close()

	if err != nil {
		lg := domain.NewErrorLog()
		lg.Logging(err, "")
		return
	}
	for rows.Next() {
		var p domain.TRelationPlatform
		err = rows.Scan(
			&p.PlatformId, &p.AnimeId, &p.LinkUrl,
			&p.CreatedAt, &p.UpdatedAt, &p.PlatName,
		)
		if err != nil {
			lg := domain.NewErrorLog()
			lg.Logging(err, "")
			return
		}
		platforms = append(platforms, p)
	}
	return
}
