package usecase_test

import (
	"animar/v1/internal/pkg/domain"
	"animar/v1/internal/pkg/interfaces/database"
	"animar/v1/internal/pkg/tools/tools"
	"animar/v1/internal/pkg/usecase"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var (
	rawContent = "良かった"
	rawRating  = 10
	rawUserId  = "aaaaaaaa"
	slug       = "slug"
	title      = "タイトル"
)

type fakeReviewRepository struct {
	database.ReviewRepository
}

func (repo *fakeReviewRepository) FindById(id int) (domain.ReviewWithAnimeSlug, error) {
	mockReview := domain.ReviewWithAnimeSlug{
		TReview: domain.TReview{
			ID:        id,
			Content:   tools.NewNullString(rawContent),
			Rating:    tools.NewNullInt(rawRating),
			AnimeId:   10,
			UserId:    tools.NewNullString(rawUserId),
			CreatedAt: time.Now().Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt: time.Now().Format("2006-01-02T15:04:05Z07:00"),
		},
		AnimeSlug:  slug,
		AnimeTitle: title,
	}
	return mockReview, nil
}

func TestFindById(t *testing.T) {
	mockReviewRepo := new(fakeReviewRepository)

	interactor := usecase.NewReviewInteractor(mockReviewRepo)
	r, err := interactor.GetReviewById(rawRating)

	t.Run("success", func(t *testing.T) {
		assert.NoError(t, err)
		assert.NotNil(t, r)

		assert.Equal(t, rawContent, *r.Content)
	})
}
