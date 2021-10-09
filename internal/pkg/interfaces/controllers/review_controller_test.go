package controllers

import (
	"animar/v1/internal/pkg/domain"
	"animar/v1/internal/pkg/usecase"
	"bytes"
	"encoding/json"
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

func (ri *revI) UpsertReviewRating(r domain.TReviewInput, userId string) (int, error) {
	return 4, nil
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

func Test_UpsertReviewRatingView(t *testing.T) {
	con := ReviewController{
		interactor: &revI{},
	}

	table := []struct {
		testName   string
		input      domain.TReviewInput
		userId     string
		wantStatus int
	}{
		{
			"success",
			domain.TReviewInput{
				AnimeId: 200,
				Content: "",
				Rating:  8,
			},
			"userId",
			200,
		},
	}

	for _, tt := range table {
		t.Run(tt.testName, func(t *testing.T) {
			j, err := json.Marshal(tt.input)
			if err != nil {
				t.Fatal(err)
			}

			req := httptest.NewRequest(http.MethodGet, "https://abc/def/", bytes.NewBuffer(j))
			req = setUserToContext(req, tt.userId)

			got := httptest.NewRecorder()
			con.UpsertReviewRatingView(got, req)
			assert.Equal(t, tt.wantStatus, got.Result().StatusCode)
		})
	}
}
