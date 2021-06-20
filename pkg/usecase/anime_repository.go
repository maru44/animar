package usecase

import "animar/v1/pkg/domain"

type AnimeRepository interface {
	ListAll() (domain.TAnimes, error)
}
