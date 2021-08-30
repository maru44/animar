package usecase

import "animar/v1/internal/pkg/domain"

type ArticleInteractor struct {
	artr ArticleRepository
}

func NewArticleInteractor(artr ArticleRepository) domain.ArticleInteractor {
	return &ArticleInteractor{
		artr: artr,
	}
}

/************************
        repository
************************/

type ArticleRepository interface {
	Fetch() ([]domain.Article, error)
	GetById(id int) (domain.Article, error)
	GetBySlug(slug string) (domain.Article, error)
	Insert(articleInput domain.ArticleInput, userId string) (int, error)
	Update(articleInput domain.ArticleInput, id int, userId string) (int, error)
	Delete(id int) (int, error)
	FilterByAnime(animeId int) ([]domain.Article, error)
	FilterCharaById(articleId int) ([]domain.ArticleCharacter, error)
	FilterCharaByUserId(userId string) ([]domain.ArticleCharacter, error)
	InsertChara(charaInput domain.ArticleCharacterInput, animeId int, userId string) (int, error)
	UpdateChara(charaInput domain.ArticleCharacterInput, id int, userId string) (int, error)
	DeleteChara(id int) (int, error)
	FetchInterview(articleId int) ([]domain.InterviewQuote, error)
	InsertInterview(interviewInput domain.InterviewQuoteInput, userId string) (int, error)
	UpdateInterview(interviewInput domain.InterviewQuoteInput, id int) (int, error)
	DeleteInterview(id int) (int, error)
	InsertRelationArticleCharacter(in domain.RelationArticleCharacterInput) (int, error)
	InsertRelationArticleAnime(in domain.RelationArticleAnimeInput) (int, error)
}

/**********************
   interactor methods
***********************/

func (arti *ArticleInteractor) FetchArticles() ([]domain.Article, error) {
	return arti.artr.Fetch()
}

func (arti *ArticleInteractor) GetArticleById(id int) (domain.Article, error) {
	return arti.artr.GetById(id)
}

func (arti *ArticleInteractor) GetArticleBySlug(slug string) (domain.Article, error) {
	return arti.artr.GetBySlug(slug)
}

func (arti *ArticleInteractor) InsertArticle(articleInput domain.ArticleInput, userId string) (int, error) {
	return arti.artr.Insert(articleInput, userId)
}

func (arti *ArticleInteractor) UpdateArticle(articleInput domain.ArticleInput, id int, userId string) (int, error) {
	return arti.artr.Update(articleInput, id, userId)
}

func (arti *ArticleInteractor) DeleteArticle(id int) (int, error) {
	return arti.artr.Delete(id)
}

func (arti *ArticleInteractor) FetchArticleByAnime(animeId int) ([]domain.Article, error) {
	return arti.artr.FilterByAnime(animeId)
}

func (arti *ArticleInteractor) FetchArticleCharas(articleId int) ([]domain.ArticleCharacter, error) {
	return arti.artr.FilterCharaById(articleId)
}

func (arti *ArticleInteractor) FetchArticleCharasByUser(userId string) ([]domain.ArticleCharacter, error) {
	return arti.artr.FilterCharaByUserId(userId)
}

func (arti *ArticleInteractor) InsertArticleChara(charaInput domain.ArticleCharacterInput, animeId int, userId string) (int, error) {
	return arti.artr.InsertChara(charaInput, animeId, userId)
}

func (arti *ArticleInteractor) UpdateArticleChara(charaInput domain.ArticleCharacterInput, id int, userId string) (int, error) {
	return arti.artr.UpdateChara(charaInput, id, userId)
}

func (arti *ArticleInteractor) DeleteArticleChara(id int) (int, error) {
	return arti.artr.DeleteChara(id)
}

func (arti *ArticleInteractor) FetchInterviewQuotes(articleId int) ([]domain.InterviewQuote, error) {
	return arti.artr.FetchInterview(articleId)
}

func (arti *ArticleInteractor) InsertInterviewQuote(interviewInput domain.InterviewQuoteInput, userId string) (int, error) {
	return arti.artr.InsertInterview(interviewInput, userId)
}

func (arti *ArticleInteractor) UpdateInterviewQuote(interviewInput domain.InterviewQuoteInput, id int) (int, error) {
	return arti.artr.UpdateInterview(interviewInput, id)
}

func (arti *ArticleInteractor) DeleteInterviewQuote(id int) (int, error) {
	return arti.artr.DeleteInterview(id)
}

func (arti *ArticleInteractor) InsertRelationArticleCharacter(in domain.RelationArticleCharacterInput) (int, error) {
	return arti.artr.InsertRelationArticleCharacter(in)
}

func (arti *ArticleInteractor) InsertRelationArticleAnime(in domain.RelationArticleAnimeInput) (int, error) {
	return arti.artr.InsertRelationArticleAnime(in)
}
