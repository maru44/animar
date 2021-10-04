package database

import (
	"animar/v1/internal/pkg/domain"
	"animar/v1/internal/pkg/tools/tools"

	"github.com/maru44/perr"
)

type ArticleRepository struct {
	SqlHandler
}

func (artr *ArticleRepository) Fetch() (articles []domain.Article, err error) {
	rows, err := artr.Query(
		"SELECT id, slug, article_type, abstract, content, image, author, is_public, user_id, created_at, updated_at, " +
			"FROM articles " +
			"WHERE is_public = true " +
			"ORDER BY created_at DESC",
	)
	if err != nil {
		return articles, perr.Wrap(err, perr.InternalServerErrorWithUrgency)
	}
	defer rows.Close()

	for rows.Next() {
		var a domain.Article
		err = rows.Scan(
			&a.ID, &a.Slug, &a.ArticleType, &a.Abstract, &a.Content,
			&a.Image, &a.Author, &a.IsPublic, &a.CreatedAt, &a.UpdatedAt,
		)
		if err != nil {
			return articles, perr.Wrap(err, perr.NotFound)
		}
		articles = append(articles, a)
	}
	return
}

func (artr *ArticleRepository) GetById(id int) (a domain.Article, err error) {
	rows, err := artr.Query(
		"SELECT id, slug, article_type, abstract, content, image, author, is_public, user_id, created_at, updated_at, "+
			"FROM articles "+
			"WHERE id = ?",
		id,
	)
	if err != nil {
		return a, perr.Wrap(err, perr.InsufficientStorageWithUrgency)
	}
	defer rows.Close()

	rows.Next()
	err = rows.Scan(
		&a.ID, &a.Slug, &a.ArticleType, &a.Abstract, &a.Content, &a.Image,
		&a.Author, &a.CreatedAt, &a.UpdatedAt,
	)
	if err != nil {
		return a, perr.Wrap(err, perr.NotFound)
	}
	return
}

func (artr *ArticleRepository) GetBySlug(slug string) (a domain.Article, err error) {
	rows, err := artr.Query(
		"SELECT id, slug, article_type, abstract, content, image, author, is_public, user_id, created_at, updated_at, "+
			"FROM articles "+
			"WHERE slug = ?",
		slug,
	)
	if err != nil {
		return a, perr.Wrap(err, perr.InternalServerErrorWithUrgency)
	}
	defer rows.Close()

	rows.Next()
	err = rows.Scan(
		&a.ID, &a.Slug, &a.ArticleType, &a.Abstract, &a.Content, &a.Image,
		&a.Author, &a.CreatedAt, &a.UpdatedAt,
	)
	if err != nil {
		return a, perr.Wrap(err, perr.NotFound)
	}
	return
}

func (artr *ArticleRepository) Insert(a domain.ArticleInput, userId string) (inserted int, err error) {
	exe, err := artr.Execute(
		"INSERT INTO articles(slug, article_type, abstract, content, image, author, is_public, user_id) "+
			"VALUES (?, ?, ?, ?, ?, ?, ?, ?) ",
		tools.GenRandSlug(12), a.ArticleType, a.Abstract, a.Content, a.Image, a.IsPublic, userId,
	)
	if err != nil {
		return inserted, perr.Wrap(err, perr.BadRequest)
	}

	rawInserted, err := exe.LastInsertId()
	if err != nil {
		return inserted, perr.Wrap(err, perr.BadRequest)
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
		return affected, perr.Wrap(err, perr.BadRequest)
	}

	rawAffected, err := exe.RowsAffected()
	if err != nil {
		return affected, perr.Wrap(err, perr.BadRequest)
	}
	return int(rawAffected), nil
}

func (artr *ArticleRepository) Delete(id int) (affected int, err error) {
	exe, err := artr.Execute(
		"DELETE FROM articles "+
			"WHERE id = ?",
		id,
	)
	if err != nil {
		return affected, perr.Wrap(err, perr.BadRequest)
	}

	rawAffected, err := exe.RowsAffected()
	if err != nil {
		return affected, perr.Wrap(err, perr.BadRequest)
	}
	return int(rawAffected), nil
}

/*   chara   */

func (artr *ArticleRepository) FilterByAnime(animeId int) (articles []domain.Article, err error) {
	rows, err := artr.Query(
		"SELECT id, slug, article_type, abstract, content, image, author, is_public, user_id, created_at, updated_at, "+
			"FROM articles "+
			"WHERE is_public = true AND id = ? "+
			"ORDER BY created_at DESC",
		animeId,
	)
	if err != nil {
		return articles, perr.Wrap(err, perr.InsufficientStorageWithUrgency)
	}
	defer rows.Close()

	for rows.Next() {
		var a domain.Article
		err = rows.Scan(
			&a.ID, &a.Slug, &a.ArticleType, &a.Abstract, &a.Content,
			&a.Image, &a.Author, &a.IsPublic, &a.CreatedAt, &a.UpdatedAt,
		)
		if err != nil {
			return articles, perr.Wrap(err, perr.NotFound)
		}
		articles = append(articles, a)
	}
	return
}

func (artr *ArticleRepository) FilterCharaById(articleId int) (cs []domain.ArticleCharacter, err error) {
	rows, err := artr.Query(
		"SELECT id, chara_name, image, created_at, updated_at "+
			"FROM anime_character "+
			"WHERE article_id = ?", articleId,
	)
	if err != nil {
		return cs, perr.Wrap(err, perr.InsufficientStorageWithUrgency)
	}
	defer rows.Close()

	for rows.Next() {
		var c domain.ArticleCharacter
		err = rows.Scan(&c.ID, &c.Name, &c.Image, &c.CreatedAt, &c.UpdatedAt)
		if err != nil {
			return cs, perr.Wrap(err, perr.NotFound)
		}
		cs = append(cs, c)
	}
	return
}

func (artr *ArticleRepository) FilterCharaByUserId(userId string) (cs []domain.ArticleCharacter, err error) {
	rows, err := artr.Query(
		"SELECT id, chara_name, image, created_at, updated_at "+
			"FROM article_chara "+
			"WHERE user_id = ?",
		userId,
	)
	if err != nil {
		return cs, perr.Wrap(err, perr.InsufficientStorageWithUrgency)
	}
	defer rows.Close()

	for rows.Next() {
		var c domain.ArticleCharacter
		err = rows.Scan(&c.ID, &c.Name, &c.Image, &c.CreatedAt, &c.UpdatedAt)
		if err != nil {
			return cs, perr.Wrap(err, perr.NotFound)
		}
		cs = append(cs, c)
	}
	return cs, nil
}

func (artr *ArticleRepository) InsertChara(ci domain.ArticleCharacterInput, userId string) (inserted int, err error) {
	exe, err := artr.Execute(
		"INSERT INTO article_chara(chara_name, image, user_id) "+
			"VALUES (?, ?, ?) ",
		ci.Name, ci.Image, userId,
	)
	if err != nil {
		return inserted, perr.Wrap(err, perr.BadRequest)
	}

	rawInserted, err := exe.LastInsertId()
	if err != nil {
		return inserted, perr.Wrap(err, perr.BadRequest)
	}
	return int(rawInserted), nil
}

func (artr *ArticleRepository) UpdateChara(ci domain.ArticleCharacterInput, id int) (affected int, err error) {
	exe, err := artr.Execute(
		"UPDATE article_chara SET chara_name = ?, image = ? "+
			"WHERE id = ?",
		id,
	)
	if err != nil {
		return affected, perr.Wrap(err, perr.BadRequest)
	}

	rawAffected, err := exe.RowsAffected()
	if err != nil {
		return affected, perr.Wrap(err, perr.BadRequest)
	}
	return int(rawAffected), nil
}

func (artr *ArticleRepository) DeleteChara(id int) (affected int, err error) {
	exe, err := artr.Execute(
		"DELETE FROM article_chara "+
			"WHERE id = ?",
		id,
	)
	if err != nil {
		return affected, perr.Wrap(err, perr.BadRequest)
	}

	rawAffected, err := exe.RowsAffected()
	if err != nil {
		return affected, perr.Wrap(err, perr.BadRequest)
	}
	return int(rawAffected), nil
}

/*   interview quote   */

func (artr *ArticleRepository) FetchInterview(articleId int) (ints []domain.InterviewQuote, err error) {
	rows, err := artr.Query(
		"SELECT id, chara_id, content, user_id, sequence, created_at, updated_at "+
			"FROM interview_quote "+
			"WHERE article_id = ?",
		articleId,
	)
	if err != nil {
		return ints, perr.Wrap(err, perr.InsufficientStorageWithUrgency)
	}
	defer rows.Close()

	for rows.Next() {
		var i domain.InterviewQuote
		err := rows.Scan(
			&i.ID, &i.CharaId, &i.Content, &i.UserId, &i.Sequence,
			&i.CreatedAt, &i.UpdatedAt,
		)
		if err != nil {
			return ints, perr.Wrap(err, perr.NotFound)
		}
		ints = append(ints, i)
	}
	return
}

func (artr *ArticleRepository) InsertInterview(ii domain.InterviewQuoteInput, userId string) (inserted int, err error) {
	exe, err := artr.Execute(
		"INSERT INTO interview_quote(article_id, chara_id, content, sequence, userId) "+
			"VALUES (?, ?, ?, ?, ?)",
		ii.ArticleId, ii.CharaId, ii.Content, ii.Sequence, userId,
	)
	if err != nil {
		return inserted, perr.Wrap(err, perr.BadRequest)
	}

	rawInserted, err := exe.LastInsertId()
	if err != nil {
		return inserted, perr.Wrap(err, perr.BadRequest)
	}
	return int(rawInserted), err
}

func (artr *ArticleRepository) UpdateInterview(ii domain.InterviewQuoteInput, id int) (affected int, err error) {
	exe, err := artr.Execute(
		"UPDATE interview_quote "+
			"SET chara_id = ?, content = ?, sequence = ? "+
			"WHERE id = ?",
		ii.CharaId, ii.Content, ii.Sequence, id,
	)
	if err != nil {
		return affected, perr.Wrap(err, perr.BadRequest)
	}

	rawAffected, err := exe.RowsAffected()
	if err != nil {
		return affected, perr.Wrap(err, perr.BadRequest)
	}
	return int(rawAffected), err
}

func (artr *ArticleRepository) DeleteInterview(id int) (affected int, err error) {
	exe, err := artr.Execute(
		"DELETE FROM interview_quote "+
			"WHERE id = ?",
		id,
	)
	if err != nil {
		return affected, perr.Wrap(err, perr.BadRequest)
	}

	rawAffected, err := exe.RowsAffected()
	if err != nil {
		return affected, perr.Wrap(err, perr.BadRequest)
	}
	return int(rawAffected), err
}

/*   relation chara   */

func (artr *ArticleRepository) InsertRelationArticleCharacter(in domain.RelationArticleCharacterInput) (inserted int, err error) {
	exe, err := artr.Execute(
		"INSERT INTO relation_article_chara(article_id, chara_id) "+
			"VALUES (?, ?)",
		in.ArticleId, in.CharaId,
	)
	if err != nil {
		return inserted, perr.Wrap(err, perr.BadRequest)
	}

	rawInserted, err := exe.LastInsertId()
	if err != nil {
		return inserted, perr.Wrap(err, perr.BadRequest)
	}
	return int(rawInserted), nil
}

/*   relation anime   */

func (artr *ArticleRepository) InsertRelationArticleAnime(in domain.RelationArticleAnimeInput) (affected int, err error) {
	exe, err := artr.Execute(
		"INSERT INTO relation_article_anime(anime_id, article_id) "+
			"VALUES (?, ?)",
		in.AnimeId, in.ArticleId,
	)
	if err != nil {
		return affected, perr.Wrap(err, perr.BadRequest)
	}

	rawAffected, err := exe.RowsAffected()
	if err != nil {
		return affected, perr.Wrap(err, perr.BadRequest)
	}
	return int(rawAffected), err
}

func (artr *ArticleRepository) DeleteRelationArticleCharacter(in domain.RelationArticleCharacterInput) (affected int, err error) {
	exe, err := artr.Execute(
		"DELETE FROM relation_article_chara "+
			"WHERE chara_id = ? AND article_id = ?",
		in.CharaId, in.ArticleId,
	)
	if err != nil {
		return affected, perr.Wrap(err, perr.BadRequest)
	}

	rawAffected, err := exe.RowsAffected()
	if err != nil {
		return affected, perr.Wrap(err, perr.BadRequest)
	}
	return int(rawAffected), err
}

func (artr *ArticleRepository) DeleteRelationArticleAnime(in domain.RelationArticleAnimeInput) (affected int, err error) {
	exe, err := artr.Execute(
		"DELETE FROM relation_article_anime "+
			"WHERE chara_id = ? AND article_id = ?",
		in.AnimeId, in.ArticleId,
	)
	if err != nil {
		return affected, perr.Wrap(err, perr.BadRequest)
	}

	rawAffected, err := exe.RowsAffected()
	if err != nil {
		return affected, perr.Wrap(err, perr.BadRequest)
	}
	return int(rawAffected), err
}
