package database

import (
	"animar/v1/internal/pkg/domain"
	"animar/v1/internal/pkg/tools/tools"
	"strconv"

	"github.com/maru44/perr"
)

type BlogRepository struct {
	SqlHandler
}

func (repo *BlogRepository) ListAll() (blogs domain.TBlogJoinAnimes, err error) {
	rows, err := repo.Query(
		"SELECT * FROM blogs WHERE is_public = true",
	)
	if err != nil {
		return blogs, perr.Wrap(err, perr.InternalServerErrorWithUrgency)
	}
	defer rows.Close()

	for rows.Next() {
		var b domain.TBlogJoinAnime
		err := rows.Scan(
			&b.ID, &b.Slug, &b.Title, &b.Abstract, &b.Content,
			&b.UserId, &b.IsPublic, &b.CreatedAt, &b.UpdatedAt,
		)
		b.Animes, _ = repo.FilterByBlog(b.GetId())
		if err != nil {
			return blogs, perr.Wrap(err, perr.NotFound)
		}
		blogs = append(blogs, b)
	}
	return
}

func (repo *BlogRepository) FilterByUser(accessUserId string, blogUserId string) (blogs domain.TBlogJoinAnimes, err error) {
	var rows Rows
	if accessUserId == blogUserId {
		rows, err = repo.Query(
			"SELECT * FROM blogs WHERE user_id = ?", blogUserId,
		)
	} else {
		rows, err = repo.Query(
			"SELECT * FROM blogs WHERE user_id = ? AND is_public = true", blogUserId,
		)
	}
	if err != nil {
		return blogs, perr.Wrap(err, perr.InternalServerErrorWithUrgency)
	}
	defer rows.Close()

	for rows.Next() {
		var b domain.TBlogJoinAnime
		err = rows.Scan(
			&b.ID, &b.Slug, &b.Title, &b.Abstract, &b.Content,
			&b.UserId, &b.IsPublic, &b.CreatedAt, &b.UpdatedAt,
		)
		b.Animes, _ = repo.FilterByBlog(b.GetId())
		if err != nil {
			return blogs, perr.Wrap(err, perr.NotFound)
		}
		blogs = append(blogs, b)
	}
	return
}

func (repo *BlogRepository) GetUserId(id int) (userId string, err error) {
	rows, err := repo.Query(
		"SELECT user_id FROM blogs WHERE id = ?", id,
	)
	if err != nil {
		return userId, perr.Wrap(err, perr.InternalServerErrorWithUrgency)
	}
	defer rows.Close()

	rows.Next()
	err = rows.Scan(&userId)
	if err != nil {
		return userId, perr.Wrap(err, perr.NotFound)
	}
	return
}

func (repo *BlogRepository) FindById(id int) (b domain.TBlogJoinAnime, err error) {
	rows, err := repo.Query(
		"SELECT blogs.* FROM blogs WHERE id = ?", id,
	)
	if err != nil {
		return b, perr.Wrap(err, perr.InternalServerErrorWithUrgency)
	}
	defer rows.Close()

	rows.Next()
	err = rows.Scan(
		&b.ID, &b.Slug, &b.Title, &b.Abstract, &b.Content,
		&b.UserId, &b.IsPublic, &b.CreatedAt, &b.UpdatedAt,
	)
	if err != nil {
		return b, perr.Wrap(err, perr.NotFound)
	}
	return
}

func (repo *BlogRepository) FindBySlug(slug string) (b domain.TBlogJoinAnime, err error) {
	rows, err := repo.Query(
		"SELECT blogs.* FROM blogs WHERE slug = ?", slug,
	)
	if err != nil {
		return b, perr.Wrap(err, perr.InternalServerErrorWithUrgency)
	}
	defer rows.Close()

	rows.Next()
	err = rows.Scan(
		&b.ID, &b.Slug, &b.Title, &b.Abstract, &b.Content,
		&b.UserId, &b.IsPublic, &b.CreatedAt, &b.UpdatedAt,
	)
	if err != nil {
		return b, perr.Wrap(err, perr.NotFound)
	}
	return
}

func (repo *BlogRepository) Insert(b domain.TBlogInsert, userId string) (lastInserted int, err error) {
	slug := tools.GenRandSlug(8)
	exe, err := repo.Execute(
		"INSERT INTO blogs(slug, title, abstract, content, user_id, is_public) "+
			"VALUES(?, ?, ?, ?, ?, ?)", slug, b.Title, tools.NewNullString(b.Abstract), b.Content, userId, b.IsPublic,
	)
	if err != nil {
		return lastInserted, perr.Wrap(err, perr.InternalServerErrorWithUrgency)
	}

	rawId, err := exe.LastInsertId()
	if err != nil {
		return lastInserted, perr.Wrap(err, perr.BadRequest)
	}
	lastInserted = int(rawId)
	return
}

func (repo *BlogRepository) Update(b domain.TBlogInsert, id int) (rowsAffected int, err error) {
	exe, err := repo.Execute(
		"UPDATE blogs SET title = ?, abstract = ?, content = ?, is_public = ? WHERE id = ?",
		b.Title, tools.NewNullString(b.Abstract), b.Content, b.IsPublic, id,
	)
	if err != nil {
		return rowsAffected, perr.Wrap(err, perr.InternalServerErrorWithUrgency)
	}

	rawId, err := exe.RowsAffected()
	if err != nil {
		return rowsAffected, perr.Wrap(err, perr.BadRequest)
	}
	rowsAffected = int(rawId)
	return
}

func (repo *BlogRepository) Delete(id int) (rowsAffected int, err error) {
	exe, err := repo.Execute(
		"DELETE FROM blogs WHERE id = ?", id,
	)
	if err != nil {
		return rowsAffected, perr.Wrap(err, perr.InternalServerErrorWithUrgency)
	}

	rawId, err := exe.RowsAffected()
	if err != nil {
		return rowsAffected, perr.Wrap(err, perr.BadRequest)
	}
	rowsAffected = int(rawId)
	return
}

// relation

func (repo *BlogRepository) FilterByBlog(blogId int) (animes []domain.TJoinedAnime, err error) {
	rows, err := repo.Query(
		"SELECT relation_blog_animes.anime_id, animes.slug, animes.title FROM relation_blog_animes " +
			"LEFT JOIN animes ON animes.id = relation_blog_animes.anime_id " +
			"WHERE blog_id = " + strconv.Itoa(blogId),
	)
	if err != nil {
		return animes, perr.Wrap(err, perr.InternalServerErrorWithUrgency)
	}
	defer rows.Close()

	for rows.Next() {
		var a domain.TJoinedAnime
		err := rows.Scan(
			&a.AnimeId, &a.Slug, &a.Title,
		)
		if err != nil {
			return animes, perr.Wrap(err, perr.NotFound)
		}
		animes = append(animes, a)
	}
	return
}

// func (repo *BlogRepository) RelationBlogByAnime(animeId int) (blogs domain.TBlogs, err error) {
// 	rows, err := repo.Query(
// 		"SELECT relation_blog_animes.blog_id, blogs.slug, blogs.title, " +
// 			"blogs.abstract, blogs.created_at FROM relation_blog_animes " +
// 			"LEFT JOIN blogs ON blogs.id = relation_blog_animes.blog_id " +
// 			"WHERE anime_id = " + strconv.Itoa(animeId),
// 	)
// 	if err != nil {
// 		lg := domain.NewErrorLog(err.Error(), "")
// lg.Logging()
// 		return
// 	}
// 	defer rows.Close()
// 	for rows.Next() {
// 		var b domain.TBlog
// 		err = rows.Scan(
// 			b,
// 		)
// 		if err != nil {
// 			lg := domain.NewErrorLog(err.Error(), "")
// lg.Logging()
// 			return
// 		}
// 		blogs = append(blogs, b)
// 	}
// 	return
// }

func (repo *BlogRepository) InsertRelation(animeId int, blogId int) (is_success bool, err error) {
	exe, err := repo.Execute(
		"INSERT INTO relation_blog_animes(anime_id, blog_id) VALUES(?, ?)",
		animeId, blogId,
	)
	if err != nil {
		return is_success, perr.Wrap(err, perr.InternalServerErrorWithUrgency)
	}

	_, err = exe.RowsAffected()
	if err != nil {
		return is_success, perr.Wrap(err, perr.BadRequest)
	}
	is_success = true
	return
}

func (repo *BlogRepository) DeleteRelation(animeId int, blogId int) (err error) {
	exe, err := repo.Execute(
		"DELETE FROM relation_blog_animes WHERE anime_id = ? AND blog_id = ?",
		animeId, blogId,
	)
	if err != nil {
		return perr.Wrap(err, perr.InternalServerErrorWithUrgency)
	}

	_, err = exe.RowsAffected()
	return perr.Wrap(err, perr.BadRequest)
}
