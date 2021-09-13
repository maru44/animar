package controllers

import (
	"animar/v1/internal/pkg/domain"
	"animar/v1/internal/pkg/interfaces/database"
	"animar/v1/internal/pkg/usecase"
	"encoding/json"
	"net/http"
	"strconv"
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
	articles, err := artc.interactor.FetchArticles()
	if err != nil {
		response(w, err, nil)
	} else {
		response(w, err, map[string]interface{}{"data": articles})
	}
	return
}

func (artc *ArticleController) ArticleDetailView(w http.ResponseWriter, r *http.Request) {
	slug := r.URL.Query().Get("slug")
	article, err := artc.interactor.GetArticleBySlug(slug)
	if err != nil {
		response(w, err, nil)
	} else {
		response(w, err, map[string]interface{}{"data": article})
	}
	return
}

func (artc *ArticleController) InsertArticleView(w http.ResponseWriter, r *http.Request) {
	var in domain.ArticleInput
	json.NewDecoder(r.Body).Decode(&in)
	userId := r.Context().Value(USER_ID).(string)
	inserted, err := artc.interactor.InsertArticle(in, userId)
	if err != nil {
		response(w, err, nil)
	} else {
		response(w, err, map[string]interface{}{"data": inserted})
	}
	return
}

func (artc *ArticleController) UpdateArticleView(w http.ResponseWriter, r *http.Request) {
	var in domain.ArticleInput
	json.NewDecoder(r.Body).Decode(&in)
	rawId := r.URL.Query().Get("id")
	id, err := strconv.Atoi(rawId)
	if err != nil {
		response(w, err, nil)
		return
	}
	affected, err := artc.interactor.UpdateArticle(in, id)
	if err != nil {
		response(w, err, nil)
	} else {
		response(w, err, map[string]interface{}{"data": affected})
	}
	return
}

func (artr *ArticleController) DeleteArticleView(w http.ResponseWriter, r *http.Request) {
	rawId := r.URL.Query().Get("id")
	id, err := strconv.Atoi(rawId)
	if err != nil {
		response(w, err, nil)
		return
	}
	affected, err := artr.interactor.DeleteArticle(id)
	if err != nil {
		response(w, err, nil)
	} else {
		response(w, err, map[string]interface{}{"data": affected})
	}
	return
}

// 	InsertChara(charaInput domain.ArticleCharacterInput, animeId int, userId string) (int, error)
// 	UpdateChara(charaInput domain.ArticleCharacterInput, id int, userId string) (int, error)
// 	DeleteChara(id int) (int, error)

/*  chara  */

func (artr *ArticleController) FilterCharaByArticleView(w http.ResponseWriter, r *http.Request) {
	rawId := r.URL.Query().Get("id")
	id, err := strconv.Atoi(rawId)
	if err != nil {
		response(w, err, nil)
		return
	}
	charas, err := artr.interactor.FetchArticleCharas(id)
	if err != nil {
		response(w, err, nil)
	} else {
		response(w, err, map[string]interface{}{"data": charas})
	}
	return
}

func (artr *ArticleController) FilterCharaByUserView(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(USER_ID).(string)
	charas, err := artr.interactor.FetchArticleCharasByUser(userId)
	if err != nil {
		response(w, err, nil)
	} else {
		response(w, err, map[string]interface{}{"data": charas})
	}
	return
}

func (artr *ArticleController) InsertArticleCharaView(w http.ResponseWriter, r *http.Request) {
	var in domain.ArticleCharacterInput
	json.NewDecoder(r.Body).Decode(&in)

	userId := r.Context().Value(USER_ID).(string)

	inserted, err := artr.interactor.InsertArticleChara(in, userId)
	if err != nil {
		response(w, err, nil)
	} else {
		response(w, err, map[string]interface{}{"data": inserted})
	}
	return
}

func (artr *ArticleController) UpdateArticleCharaView(w http.ResponseWriter, r *http.Request) {
	var in domain.ArticleCharacterInput
	json.NewDecoder(r.Body).Decode(&in)

	// userId := r.Context().Value(USER_ID).(string)
	rawId := r.URL.Query().Get("id")
	id, err := strconv.Atoi(rawId)
	if err != nil {
		response(w, err, nil)
		return
	}

	affected, err := artr.interactor.UpdateArticleChara(in, id)
	if err != nil {
		response(w, err, nil)
	} else {
		response(w, err, map[string]interface{}{"data": affected})
	}
	return
}

func (artr *ArticleController) DeleteArticleCharaView(w http.ResponseWriter, r *http.Request) {
	rawId := r.URL.Query().Get("id")
	id, err := strconv.Atoi(rawId)
	if err != nil {
		response(w, err, nil)
		return
	}

	affected, err := artr.interactor.DeleteArticleChara(id)
	if err != nil {
		response(w, err, nil)
	} else {
		response(w, err, map[string]interface{}{"data": affected})
	}
	return
}
