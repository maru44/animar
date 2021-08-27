package database

import (
	"animar/v1/internal/pkg/domain"
	"animar/v1/internal/pkg/tools/tools"
)

type ArticleRepository struct {
	SqlHandler
}

/*   interactor   */
// type ArticleInteractor interface {
// 	FetchArticles() ([]Article, error)
// 	GetArticleById(id int) (Article, error)
// 	GetArticleBySlug(slug string) (Article, error)
// 	InsertArticle(articleInput ArticleInput) (int, error)
// 	UpdateArticle(articleInput ArticleInput, id int) (int, error)
// 	DeleteArticle(id int) (int, error)
// 	FetchArticleByAnime(animeId int) ([]Article, error)

// 	// character

// 	FetchArticleCharas(articleId int) ([]ArticleCharacter, error)
// 	FetchArticleCharasByUser(userId string) ([]ArticleCharacter, error)
// 	InsertArticleChara(charaInput ArticleCharacterInput) (int, error)
// 	UpdateArticleChara(charaInput ArticleCharacterInput, id int) (int, error)
// 	DeleteArticleChara(id int) (int, error)

// 	// interview

// 	FetchInterviewQuotes(articleId int) ([]InterviewQuote, error)
// 	InsertInterviewQuote(interviewInput InterviewQuoteInput) (int, error)
// 	UpdateInterviewQuote(interviewInput InterviewQuoteInput, id int) (int, error)
// 	DeleteInterviewQuote(id int) (int, error)
// }

func (artr *ArticleRepository) Fetch() (articles []domain.Article, err error) {
	rows, err := artr.Query(
		"SELECT id, slug, article_type, abstract, content, image, author, is_public, user_id, created_at, updated_at, ",
		"FROM articles ",
		"WHERE is_public = true ",
		"ORDER BY created_at DESC",
	)
	if err != nil {
		domain.ErrorWarn(err)
		return
	}
	for rows.Next() {
		var a domain.Article
		err = rows.Scan(
			&a.ID, &a.Slug, &a.ArticleType, &a.Abstract, &a.Content,
			&a.Image, &a.Author, &a.IsPublic, &a.CreatedAt, &a.UpdatedAt,
		)
		a.Characters, err = artr.FilterCharaById(a.ID)
		if err != nil {
			domain.ErrorWarn(err)
			return
		}
		articles = append(articles, a)
	}
	return
}

func (artr *ArticleRepository) GetById(id int) (a domain.Article, err error) {
	rows, err := artr.Query(
		"SELECT id, slug, article_type, abstract, content, image, author, is_public, user_id, created_at, updated_at, ",
		"FROM articles ",
		"WHERE id = ?",
		id,
	)
	if err != nil {
		domain.ErrorWarn(err)
		return
	}
	rows.Next()
	err = rows.Scan(
		&a.ID, &a.Slug, &a.ArticleType, &a.Abstract, &a.Content, &a.Image,
		&a.Author, &a.CreatedAt, &a.UpdatedAt,
	)
	if err != nil {
		domain.ErrorWarn(err)
		return
	} else {
		if cs, err := artr.FilterCharaById(id); err != nil {
			domain.ErrorWarn(err)
		} else {
			a.Characters = cs
		}
	}
	return
}

func (artr *ArticleRepository) GetBySlug(slug string) (a domain.Article, err error) {
	rows, err := artr.Query(
		"SELECT id, slug, article_type, abstract, content, image, author, is_public, user_id, created_at, updated_at, ",
		"FROM articles ",
		"WHERE slug = ?",
		slug,
	)
	if err != nil {
		domain.ErrorWarn(err)
		return
	}
	rows.Next()
	err = rows.Scan(
		&a.ID, &a.Slug, &a.ArticleType, &a.Abstract, &a.Content, &a.Image,
		&a.Author, &a.CreatedAt, &a.UpdatedAt,
	)
	if err != nil {
		domain.ErrorWarn(err)
		return
	} else {
		// if cs, err := artr.FilterCharaById(s); err != nil {
		// 	domain.ErrorWarn(err)
		// } else {
		// 	a.Characters = cs
		// }
	}
	return
}

func (artr *ArticleRepository) Insert(a domain.ArticleInput) (inserted int, err error) {
	exe, err := artr.Execute(
		"INSERT INTO (slug, article_type, abstract, content, image, author, is_public, user_id) ",
		"VALUES (?, ?, ?, ?, ?, ?, ?, ?) ",
		tools.GenRandSlug(12), a.ArticleType, a.Abstract, a.Content, a.Image, a.IsPublic, a.UserId,
	)
	if err != nil {
		domain.ErrorWarn(err)
		return
	}
	rawInserted, err := exe.LastInsertId()
	if err != nil {
		domain.ErrorWarn(err)
		return
	}
	return int(rawInserted), err
}

func (artr *ArticleRepository) Update(a domain.ArticleInput, articleId int) (affected int, err error) {
	exe, err := artr.Execute(
		"UPDATE articles SET article_type, abstract = ?, content = ?, image = ?, author = ?, is_public = ? "+
			"WHERE id = ?",
		a.ArticleType, a.Abstract, a.Content, a.Image, a.Author, a.IsPublic, articleId,
	)
	if err != nil {
		domain.ErrorWarn(err)
		return
	}
	rawAffected, err := exe.RowsAffected()
	if err != nil {
		domain.ErrorWarn(err)
		return
	}
	return int(rawAffected), err
}

// func (artr *ArticleRepository) Delete(id int) (affected int, err error) {
// }

// func (artr *ArticleRepository) FilterByAnime(animeId int) (as []domain.Article, err error) {
// }

func (artr *ArticleRepository) FilterCharaById(articleId int) (cs []domain.ArticleCharacter, err error) {
	rows, err := artr.Query(
		"SELECT id, chara_name, image, created_at, updated_at ",
		"FROM anime_character ",
		"WHERE article_id = ?", articleId,
	)
	if err != nil {
		domain.ErrorWarn(err)
		return
	}
	for rows.Next() {
		var c domain.ArticleCharacter
		rows.Scan(&c.ID, &c.Name, &c.Image, &c.CreatedAt, &c.UpdatedAt)
		cs = append(cs, c)
	}
	return
}

// func (artr *ArticleRepository) FilterCharaByUserId(userId string) (cs []domain.ArticleCharacter, err error) {
// }

// func (artr *ArticleRepository) InsertChara(ci domain.ArticleCharacterInput) (inserted int, err error) {
// }

// func (artr *ArticleRepository) UpdateChara(ci domain.ArticleCharacterInput, id int) (affected int, err error) {}

// func (artr *ArticleRepository) DeleteChara(id int) (affected int, err error) {}

// func (artr *ArticleRepository) FetchInterview(articleId int) (ints []domain.InterviewQuote, err error) {}

// func (artr *ArticleRepository) InsertInterview(ii domain.InterviewQuoteInput) (int, error)

// func (artr *ArticleRepository) UpdateInterview(ii domain.InterviewQuoteInput, id int) (int, error) {}

// func (artr *ArticleRepository) DeleteInterview(id int) (int, error) {}
