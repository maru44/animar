package controllers

import (
	"animar/v1/internal/pkg/domain"
	"animar/v1/internal/pkg/interfaces/database"
	"animar/v1/internal/pkg/usecase"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/maru44/perr"
)

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

// @TODO:Add filter by userid(query params)
func (artc *ArticleController) ArticleListView(w http.ResponseWriter, r *http.Request) {
	articles, err := artc.interactor.FetchArticles()
	response(w, r, perr.Wrap(err, perr.NotFound), map[string]interface{}{"data": articles})
	return
}

func (artc *ArticleController) ArticleDetailView(w http.ResponseWriter, r *http.Request) {
	slug := r.URL.Query().Get("slug")
	article, err := artc.interactor.GetArticleBySlug(slug)
	if err != nil {
		response(w, r, perr.Wrap(err, perr.NotFound), nil)
		return
	}

	charas, err := artc.interactor.FetchArticleCharas(article.ID)
	if err != nil {
		response(w, r, perr.Wrap(err, perr.NotFound), nil)
		return
	}

	if article.ArticleType == domain.ArticleTypeInterview {
		interviews, err := artc.interactor.FetchInterviewQuotes(article.ID)
		response(w, r, perr.Wrap(err, perr.NotFound), map[string]interface{}{"article": article, "charas": charas, "interviews": interviews})
	} else {
		response(w, r, nil, map[string]interface{}{"article": article, "charas": charas})
	}

	return
}

func (artc *ArticleController) InsertArticleView(w http.ResponseWriter, r *http.Request) {
	var in domain.ArticleInput
	json.NewDecoder(r.Body).Decode(&in)
	userId := r.Context().Value(USER_ID).(string)
	inserted, err := artc.interactor.InsertArticle(in, userId)

	response(w, r, perr.Wrap(err, perr.BadRequest), map[string]interface{}{"data": inserted})
	return
}

func (artc *ArticleController) UpdateArticleView(w http.ResponseWriter, r *http.Request) {
	var in domain.ArticleInput
	json.NewDecoder(r.Body).Decode(&in)
	rawId := r.URL.Query().Get("id")
	id, err := strconv.Atoi(rawId)
	if err != nil {
		response(w, r, perr.Wrap(err, perr.BadRequest), nil)
		return
	}
	affected, err := artc.interactor.UpdateArticle(in, id)

	response(w, r, perr.Wrap(err, perr.BadRequest), map[string]interface{}{"data": affected})
	return
}

func (artr *ArticleController) DeleteArticleView(w http.ResponseWriter, r *http.Request) {
	rawId := r.URL.Query().Get("id")
	id, err := strconv.Atoi(rawId)
	if err != nil {
		response(w, r, perr.Wrap(err, perr.BadRequest), nil)
		return
	}

	affected, err := artr.interactor.DeleteArticle(id)
	response(w, r, perr.Wrap(err, perr.BadRequest), map[string]interface{}{"data": affected})
	return
}

/*  chara  */
func (artc *ArticleController) ArticleCharacterView(w http.ResponseWriter, r *http.Request) {
	rawId := r.URL.Query().Get("id")
	id, err := strconv.Atoi(rawId)
	var charas []domain.ArticleCharacter
	if err != nil {
		userId := r.Context().Value(USER_ID).(string)
		charas, err = artc.interactor.FetchArticleCharasByUser(userId)
		response(w, r, perr.Wrap(err, perr.BadRequest), map[string]interface{}{"data": charas})
	} else {
		charas, err = artc.interactor.FetchArticleCharas(id)
		response(w, r, perr.Wrap(err, perr.BadRequest), map[string]interface{}{"data": charas})
	}
	return
}

func (artc *ArticleController) InsertArticleCharaView(w http.ResponseWriter, r *http.Request) {
	var in domain.ArticleCharacterInput
	json.NewDecoder(r.Body).Decode(&in)

	userId := r.Context().Value(USER_ID).(string)

	inserted, err := artc.interactor.InsertArticleChara(in, userId)
	response(w, r, perr.Wrap(err, perr.BadRequest), map[string]interface{}{"data": inserted})
	return
}

func (artr *ArticleController) UpdateArticleCharaView(w http.ResponseWriter, r *http.Request) {
	var in domain.ArticleCharacterInput
	json.NewDecoder(r.Body).Decode(&in)

	rawId := r.URL.Query().Get("id")
	id, err := strconv.Atoi(rawId)
	if err != nil {
		response(w, r, perr.Wrap(err, perr.BadRequest), nil)
		return
	}

	affected, err := artr.interactor.UpdateArticleChara(in, id)
	response(w, r, perr.Wrap(err, perr.BadRequest), map[string]interface{}{"data": affected})
	return
}

func (artr *ArticleController) DeleteArticleCharaView(w http.ResponseWriter, r *http.Request) {
	rawId := r.URL.Query().Get("id")
	id, err := strconv.Atoi(rawId)
	if err != nil {
		response(w, r, perr.Wrap(err, perr.BadRequest), nil)
		return
	}

	affected, err := artr.interactor.DeleteArticleChara(id)
	response(w, r, perr.Wrap(err, perr.BadRequest), map[string]interface{}{"data": affected})
	return
}

/*  interview  */

func (artc *ArticleController) FetchInterviewView(w http.ResponseWriter, r *http.Request) {
	rawId := r.URL.Query().Get("id")
	id, err := strconv.Atoi(rawId)
	if err != nil {
		response(w, r, perr.Wrap(err, perr.BadRequest), nil)
		return
	}

	ints, err := artc.interactor.FetchInterviewQuotes(id)
	response(w, r, perr.Wrap(err, perr.BadRequest), map[string]interface{}{"data": ints})
	return
}

func (artc *ArticleController) InsertInterviewView(w http.ResponseWriter, r *http.Request) {
	var res domain.InterviewQuoteInput
	if err := json.NewDecoder(r.Body).Decode(&res); err != nil {
		response(w, r, perr.Wrap(err, perr.BadRequest), nil)
		return
	}

	userId := r.Context().Value(USER_ID).(string)
	inserted, err := artc.interactor.InsertInterviewQuote(res, userId)
	response(w, r, perr.Wrap(err, perr.BadRequest), map[string]interface{}{"data": inserted})
	return
}

func (artc *ArticleController) UpdateInterviewView(w http.ResponseWriter, r *http.Request) {
	rawId := r.URL.Query().Get("id")
	id, err := strconv.Atoi(rawId)
	if err != nil {
		response(w, r, perr.Wrap(err, perr.BadRequest), nil)
		return
	}
	var res domain.InterviewQuoteInput
	if err := json.NewDecoder(r.Body).Decode(&res); err != nil {
		response(w, r, perr.Wrap(err, perr.BadRequest), nil)
		return
	}

	// userId := r.Context().Value(USER_ID).(string)
	affected, err := artc.interactor.UpdateInterviewQuote(res, id)
	response(w, r, perr.Wrap(err, perr.BadRequest), map[string]interface{}{"data": affected})
	return
}

func (artc *ArticleController) DeleteInterviewView(w http.ResponseWriter, r *http.Request) {
	rawId := r.URL.Query().Get("id")
	id, err := strconv.Atoi(rawId)
	if err != nil {
		response(w, r, perr.Wrap(err, perr.BadRequest), nil)
		return
	}

	affected, err := artc.interactor.DeleteInterviewQuote(id)
	response(w, r, perr.Wrap(err, perr.BadRequest), map[string]interface{}{"data": affected})
	return
}

func (artc *ArticleController) InsertRelationArticleCharacterView(w http.ResponseWriter, r *http.Request) {
	var in domain.RelationArticleCharacterInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		response(w, r, perr.Wrap(err, perr.BadRequest), nil)
		return
	}

	inserted, err := artc.interactor.InsertRelationArticleCharacter(in)
	response(w, r, perr.Wrap(err, perr.BadRequest), map[string]interface{}{"data": inserted})
	return
}

func (artc *ArticleController) DeleteRelationArticleCharacterView(w http.ResponseWriter, r *http.Request) {
	var in domain.RelationArticleCharacterInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		response(w, r, perr.Wrap(err, perr.BadRequest), nil)
		return
	}

	affected, err := artc.interactor.DeleteRelationArticleCharacter(in)
	response(w, r, perr.Wrap(err, perr.BadRequest), map[string]interface{}{"data": affected})
	return
}

// func (artc *ArticleController) InsertRelationArticleAnimeView(w http.ResponseWriter, r *http.Request) {
// 	var in domain.RelationArticleAnimeInput
// 	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
// 		response(w, r, err, nil)
// 		return
// 	}

// 	inserted, err := artc.interactor.InsertRelationArticleAnime(in)
// 	if err != nil {
// 		response(w, r, err, nil)
// 	} else {
// 		response(w, r, err, map[string]interface{}{"data": inserted})
// 	}
// 	return
// }

// 	InsertRelationArticleCharacter(in domain.RelationArticleCharacterInput) (int, error)
// 	DeleteRelationArticleCharacter(in domain.RelationArticleCharacterInput) (int, error)
// 	InsertRelationArticleAnime(in domain.RelationArticleAnimeInput) (int, error)
// 	DeleteRelationArticleAnime(in domain.RelationArticleAnimeInput) (int, error)
// }
