package usecase

import "animar/v1/internal/pkg/domain"

type ArticleInteractor struct {
	artr  ArticleRepository
	artcr ArticleCharacterReposiotry
	iqr   InterviewQuoteRepository
}

func NewArticleInteractor(artr ArticleRepository, artcr ArticleCharacterReposiotry, iqr InterviewQuoteRepository) domain.ArticleInteractor {
	return &ArticleInteractor{
		artr:  artr,
		artcr: artcr,
		iqr:   iqr,
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
}

type ArticleCharacterReposiotry interface {
	Insert(charaInput domain.ArticleCharacterInput) (int, error)
	Update(charaInput domain.ArticleCharacterInput, id int) (int, error)
	Delete(id int) (int, error)
}

type InterviewQuoteRepository interface {
	Fetch(articleId int) ([]domain.InterviewQuote, error)
	Insert(interviewInput domain.InterviewQuoteInput) (int, error)
	Update(interviewInput domain.InterviewQuoteInput, id int) (int, error)
	Delete(id int) (int, error)
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

func (arti *ArticleInteractor) InsertArticleChara(charaInput domain.ArticleCharacterInput) (int, error) {
	return arti.artcr.Insert(charaInput)
}

func (arti *ArticleInteractor) UpdateArticleChara(charaInput domain.ArticleCharacterInput, id int) (int, error) {
	return arti.artcr.Update(charaInput, id)
}

func (arti *ArticleInteractor) DeleteArticleChara(id int) (int, error) {
	return arti.artcr.Delete(id)
}

func (arti *ArticleInteractor) FetchInterviewQuotes(articleId int) ([]domain.InterviewQuote, error) {
	return arti.iqr.Fetch(articleId)
}

func (arti *ArticleInteractor) InsertInterviewQuote(interviewInput domain.InterviewQuoteInput) (int, error) {
	return arti.iqr.Insert(interviewInput)
}

func (arti *ArticleInteractor) UpdateInterviewQuote(interviewInput domain.InterviewQuoteInput, id int) (int, error) {
	return arti.iqr.Update(interviewInput, id)
}

func (arti *ArticleInteractor) DeleteInterviewQuote(id int) (int, error) {
	return arti.iqr.Delete(id)
}
