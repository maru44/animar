package usecase

import "animar/v1/domain"

type AnimeRepository interface {
	ListAll() (domain.TAnimes, error)
	ListMinimumAll() (domain.TAnimeMinimums, error)
	ListOnAirAll() (domain.TAnimes, error)
	ListMinimumSearch(string) (domain.TAnimeMinimums, error)
	ListSearch(string) (domain.TAnimes, error)
	ListBySeason(string, string) (domain.TAnimes, error)
	// detail
	FindById(int) (domain.TAnime, error)
	FindBySlug(string) (domain.TAnimeWithSeries, error)
}
