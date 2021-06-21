package usecase

import "animar/v1/pkg/domain"

type AnimeRepository interface {
	ListAll() (domain.TAnimes, error)
	ListOnAirAll() (domain.TAnimes, error)
	ListMinimumSearch(string) (domain.TAnimeMinimums, error)
	ListSearch(string) (domain.TAnimes, error)
	ListBySeason(string, string) (domain.TAnimes, error)
	// detail
	FindById(int) (domain.TAnime, error)
	FindBySlug(string) (domain.TAnimeWithSeries, error)

	// review
	ReviewFilterByAnime(int, string) (domain.TReviews, error)
}
