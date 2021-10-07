package usecase_test

import (
	"animar/v1/internal/pkg/domain"
	"animar/v1/internal/pkg/interfaces/database"
	"animar/v1/internal/pkg/tools/tools"
	"animar/v1/internal/pkg/usecase"
	"testing"
	"time"

	"github.com/maru44/perr"
	"github.com/stretchr/testify/assert"
)

type reviewRepository struct {
	database.ReviewRepository
}

var (
	reviewFindByIdTable = []struct {
		testName string
		id       int
		review   domain.ReviewWithAnimeSlug
		err      error
	}{
		{
			"success only rating",
			22,
			domain.ReviewWithAnimeSlug{
				TReview: domain.TReview{
					ID:        22,
					Content:   tools.NewNullString(""),
					Rating:    tools.NewNullInt(8),
					AnimeId:   10,
					UserId:    tools.NewNullString("user_no_1"),
					CreatedAt: time.Now().Format("2006-01-02T15:04:05Z07:00"),
					UpdatedAt: time.Now().Format("2006-01-02T15:04:05Z07:00"),
				},
				AnimeSlug:  "anime_022",
				AnimeTitle: "rupin the third",
			},
			nil,
		},
		{
			"success only content",
			33,
			domain.ReviewWithAnimeSlug{
				TReview: domain.TReview{
					ID:        33,
					Content:   tools.NewNullString("good anime"),
					Rating:    tools.NewNullInt(0),
					AnimeId:   15,
					UserId:    tools.NewNullString("user_no_5"),
					CreatedAt: time.Now().Format("2006-01-02T15:04:05Z07:00"),
					UpdatedAt: time.Now().Format("2006-01-02T15:04:05Z07:00"),
				},
				AnimeSlug:  "anime_033",
				AnimeTitle: "cow boy bebup",
			},
			nil,
		},
	}
)

/* mock functions */

func (repo *reviewRepository) FindById(id int) (domain.ReviewWithAnimeSlug, error) {
	for _, r := range reviewFindByIdTable {
		if r.id == id {
			return r.review, nil
		}
	}
	return domain.ReviewWithAnimeSlug{}, perr.New("", perr.NotFound)
}

/* tests */

func TestFindById(t *testing.T) {
	mockReviewRepo := new(reviewRepository)
	interactor := usecase.NewReviewInteractor(mockReviewRepo)

	for _, tt := range reviewFindByIdTable {
		t.Run(tt.testName, func(t *testing.T) {
			r, err := interactor.GetReviewById(tt.id)

			assert.Equal(t, tt.err, err)
			assert.Equal(t, tt.review, r)
		})
	}
}
