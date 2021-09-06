package controllers

import (
	"animar/v1/internal/pkg/domain"
	"animar/v1/internal/pkg/interfaces/database"
	"animar/v1/internal/pkg/usecase"
	"net/http"
)

// type ArticleRepository interface {
// 	Fetch() ([]domain.Article, error)
// 	GetById(id int) (domain.Article, error)
// 	GetBySlug(slug string) (domain.Article, error)
// 	Insert(articleInput domain.ArticleInput, userId string) (int, error)
// 	Update(articleInput domain.ArticleInput, id int, userId string) (int, error)
// 	Delete(id int) (int, error)
// 	FilterByAnime(animeId int) ([]domain.Article, error)
// 	FilterCharaById(articleId int) ([]domain.ArticleCharacter, error)
// 	FilterCharaByUserId(userId string) ([]domain.ArticleCharacter, error)
// 	InsertChara(charaInput domain.ArticleCharacterInput, animeId int, userId string) (int, error)
// 	UpdateChara(charaInput domain.ArticleCharacterInput, id int, userId string) (int, error)
// 	DeleteChara(id int) (int, error)
// 	FetchInterview(articleId int) ([]domain.InterviewQuote, error)
// 	InsertInterview(interviewInput domain.InterviewQuoteInput, userId string) (int, error)
// 	UpdateInterview(interviewInput domain.InterviewQuoteInput, id int) (int, error)
// 	DeleteInterview(id int) (int, error)
// 	InsertRelationArticleCharacter(in domain.RelationArticleCharacterInput) (int, error)
// 	DeleteRelationArticleCharacter(in domain.RelationArticleCharacterInput) (int, error)
// 	InsertRelationArticleAnime(in domain.RelationArticleAnimeInput) (int, error)
// 	DeleteRelationArticleAnime(in domain.RelationArticleAnimeInput) (int, error)
// }

type ArticleController struct {
	interactor domain.ArticleInteractor
}

func NewArticleController(sqlHandler database.SqlHandler) *ArticleController {
	return &ArticleController{
		interactor: usecase.NewArticleInteractor(
			&database.ArticleRepository{
				SqlHandler: sqlHandler,
			},
		),
	}
}

func (artc *ArticleController) ArticleListView(w http.ResponseWriter, r *http.Request) {

}
