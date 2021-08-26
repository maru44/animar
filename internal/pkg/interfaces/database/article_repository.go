package database

import (
	"animar/v1/internal/pkg/domain"
	"animar/v1/internal/pkg/tools/tools"
)

type ArticleRepository struct {
	SqlHandler
}

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

func (artr *ArticleRepository) GetArticleById(id int) (a domain.Article, err error) {
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

func (artr *ArticleRepository) FilterCharaById(articleId int) (charas []domain.ArticleCharacter, err error) {
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
		charas = append(charas, c)
	}
	return
}
