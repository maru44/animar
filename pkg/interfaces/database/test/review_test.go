package database_test

import (
	"animar/v1/pkg/domain"
	"animar/v1/pkg/interfaces/database"
)

type fakeReiviewReposioty struct {
	database.ReviewRepository
	FakeFilterByAnime func(int, string) (domain.TReviews, error)
}

func (f *fakeReiviewReposioty) FilterByAnime(animeId int, userId string) (domain.TReviews, error) {
	return f.FakeFilterByAnime(animeId, userId)
}
