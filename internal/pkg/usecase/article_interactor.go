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
	Insert(articleInput domain.ArticleInput) (int, error)
	Update(articleInput domain.ArticleInput, id int) (int, error)
	Delete(id int) (int, error)
	FilterCharaById(articleId int) ([]domain.ArticleCharacter, error)
	FilterCharaByUserId(userId string) ([]domain.ArticleCharacter, error)
	InsertChara(charaInput domain.ArticleCharacterInput) (int, error)
	UpdateChara(charaInput domain.ArticleCharacterInput, id int) (int, error)
	DeleteChara(id int) (int, error)
	FetchInterview(articleId int) ([]domain.InterviewQuote, error)
	InsertInterview(interviewInput domain.InterviewQuoteInput) (int, error)
	UpdateInterview(interviewInput domain.InterviewQuoteInput, id int) (int, error)
	DeleteInterview(id int) (int, error)
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

func (arti *ArticleInteractor) InsertArticle(articleInput domain.ArticleInput) (int, error) {
	return arti.artr.Insert(articleInput)
}

func (arti *ArticleInteractor) UpdateArticle(articleInput domain.ArticleInput, id int) (int, error) {
	return arti.artr.Update(articleInput, id)
}

func (arti *ArticleInteractor) DeleteArticle(id int) (int, error) {
	return arti.artr.Delete(id)
}

func (arti *ArticleInteractor) FetchArticleCharas(articleId int) ([]domain.ArticleCharacter, error) {
	return arti.artr.FilterCharaById(articleId)
}

func (arti *ArticleInteractor) FetchArticleCharasByUser(userId string) ([]domain.ArticleCharacter, error) {
	return arti.artr.FilterCharaByUserId(userId)
}

func (arti *ArticleInteractor) InsertArticleChara(charaInput domain.ArticleCharacterInput) (int, error) {
	return arti.artr.InsertChara(charaInput)
}

func (arti *ArticleInteractor) UpdateArticleChara(charaInput domain.ArticleCharacterInput, id int) (int, error) {
	return arti.artr.UpdateChara(charaInput, id)
}

func (arti *ArticleInteractor) DeleteArticleChara(id int) (int, error) {
	return arti.artr.DeleteChara(id)
}

func (arti *ArticleInteractor) FetchInterviewQuotes(articleId int) ([]domain.InterviewQuote, error) {
	return arti.artr.FetchInterview(articleId)
}

func (arti *ArticleInteractor) InsertInterviewQuote(interviewInput domain.InterviewQuoteInput) (int, error) {
	return arti.artr.InsertInterview(interviewInput)
}

func (arti *ArticleInteractor) UpdateInterviewQuote(interviewInput domain.InterviewQuoteInput, id int) (int, error) {
	return arti.artr.UpdateInterview(interviewInput, id)
}

func (arti *ArticleInteractor) DeleteInterviewQuote(id int) (int, error) {
	return arti.artr.DeleteInterview(id)
}
