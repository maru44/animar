package controllers

import (
	"animar/v1/internal/pkg/domain"
	"animar/v1/internal/pkg/usecase"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

type (
	revI struct {
		usecase.ReviewInteractor
	}
)

/*****************************
	dummy interactor func
******************************/

func (ri *revI) GetAnimeReviews(animeId int, userId string) (domain.TReviews, error) {
	return domain.TReviews{}, nil
}

/*****************************
	test
******************************/

func Test_GetAnimeREviewsView(t *testing.T) {
	con := ReviewController{
		interactor: &revI{},
	}
	table := []struct {
		name          string
		paramsUser    string
		paramsAnimeId string
		wantStatus    int
	}{
		{
			"success",
			"useraklklwfle",
			"21",
			200,
		},
		{
			"success without user",
			"",
			"34",
			200,
		},
		{
			"no anime id",
			"jaljfwejfe",
			"",
			400,
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodGet, "https://abc/def/", nil)
			qs := r.URL.Query()
			qs.Add("anime", tt.paramsAnimeId)
			qs.Add("user", tt.paramsUser)
			r.URL.RawQuery = qs.Encode()

			got := httptest.NewRecorder()
			// act
			con.GetAnimeReviewsView(got, r)
			assert.Equal(t, tt.wantStatus, got.Result().StatusCode)
		})
	}
}
