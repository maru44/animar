package database_test

import (
	"animar/v1/internal/pkg/domain"
	infrastructure_test "animar/v1/internal/pkg/infrastructure/test"
	"animar/v1/internal/pkg/interfaces/database"
	"regexp"

	"github.com/DATA-DOG/go-sqlmock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

// success
var _ = Describe("Insert Anime OK", func() {
	lastInsertId := 21
	var lastInserted int

	kana := "どらえもん"
	engName := "doraemon"
	description := "あんなこといいな♪できたらいいな♪\nあんな夢..."
	state := "onair"
	countEpisodes := 8000
	copyright := "小学館?"

	BeforeEach(func() {
		db, mock, _ := sqlmock.New()
		defer db.Close()

		repo := &database.AdminAnimeRepository{
			SqlHandler: infrastructure_test.NewDummyHandler(db),
		}

		a := domain.AnimeInsert{
			Title:         "ドラえもん",
			Slug:          "oooooo",
			Abbreviation:  nil,
			Kana:          &kana,
			EngName:       &engName,
			Description:   &description,
			ThumbUrl:      nil,
			State:         &state,
			SeriesId:      nil,
			CountEpisodes: &countEpisodes,
			Copyright:     &copyright,
		}

		query := "INSERT INTO animes(title, slug, abbreviation, kana, eng_name, description, thumb_url, state, series_id, count_episodes, copyright) " +
			"VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"

		stmt := mock.ExpectPrepare(regexp.QuoteMeta(query))
		stmt.ExpectExec().WithArgs(
			a.Title, a.Slug, a.Abbreviation, a.Kana, a.EngName, a.Description,
			a.ThumbUrl, a.State, a.SeriesId, a.CountEpisodes, a.Copyright,
		).
			WillReturnResult(sqlmock.NewResult(int64(lastInsertId), 1))

		lastInserted, _ = repo.Insert(a)
	})

	Describe("アニメ追加", func() {
		Context("正常系", func() {
			It("lastInsertIdが0意外", func() {
				Expect(lastInserted).To(Equal(lastInsertId))
				// Ω(lastInserted).Should(lastInsertId)
			})
		})
	})
})
