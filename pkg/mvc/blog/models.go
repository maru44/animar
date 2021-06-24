package blog

import (
	"animar/v1/pkg/tools/connector"
	"animar/v1/pkg/tools/tools"
	"database/sql"
	"strconv"

	"firebase.google.com/go/v4/auth"
)

type TBlog struct {
	ID        int     `json:"id"`
	Slug      string  `json:"slug"`
	Title     string  `json:"title"`
	Abstract  *string `json:"abstract"`
	Content   string  `json:"content"`
	IsPublic  bool    `json:"is_public"`
	UserId    string  `json:"user_id"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt string  `json:"updated_at"`
}

type TJoinedAnime struct {
	AnimeId int    `json:"anime_id"`
	Slug    string `json:"slug"`
	Title   string `json:"title"`
}

type TJoinedBlog struct {
	BlogId    int    `json:"blog_id"`
	Title     string `json:"title"`
	Slug      string `json:"slug"`
	Abstract  string `json:"abstract"`
	CreatedAt string `json:"created_at"`
}

type TBlogJoinAnimes struct {
	ID        int            `json:"id"`
	Slug      string         `json:"slug"`
	Title     string         `json:"title"`
	Abstract  *string        `json:"abstract"`
	Content   string         `json:"content"`
	UserId    string         `json:"user_id"`
	IsPublic  bool           `json:"is_public"`
	CreatedAt string         `json:"created_at"`
	UpdatedAt string         `json:"updated_at"`
	Animes    []TJoinedAnime `json:"animes"`
}

type TBlogJoinAnimesUser struct {
	ID        int            `json:"id"`
	Slug      string         `json:"slug"`
	Title     string         `json:"title"`
	Abstract  *string        `json:"abstract,omitempty"`
	Content   string         `json:"content,omitempty"`
	UserId    string         `json:"user_id,omitempty"`
	IsPublic  bool           `json:"is_public"`
	CreatedAt string         `json:"created_at"`
	UpdatedAt string         `json:"updated_at,omitempty"`
	Animes    []TJoinedAnime `json:"animes,omitempty"`
	User      *auth.UserInfo `json:"user,omitempty"`
}

func ListBlog() *sql.Rows {
	db := connector.AccessDB()
	defer db.Close()
	rows, err := db.Query("SELECT * FROM blogs WHERE is_public = true")
	if err != nil {
		tools.ErrorLog(err)
	}
	return rows
}

func ListUsersBlog(uid string, userId string) *sql.Rows {
	db := connector.AccessDB()
	defer db.Close()
	var rows *sql.Rows
	var err error
	if uid == userId {
		rows, err = db.Query("SELECT * FROM blogs WHERE user_id = ?", uid)
	} else {
		rows, err = db.Query("SELECT * FROM blogs WHERE user_id = ? AND is_public = true", uid)
	}
	if err != nil {
		tools.ErrorLog(err)
	}
	return rows
}

func BlogUserId(id int) string {
	db := connector.AccessDB()
	defer db.Close()
	var userId string
	err := db.QueryRow("SELECT user_id FROM blogs WHERE id = ?", id).Scan(&userId)
	if err != nil {
		tools.ErrorLog(err)
	}
	return userId
}

func DetailBlog(id int) TBlogJoinAnimes {
	db := connector.AccessDB()
	defer db.Close()

	var b TBlogJoinAnimes
	err := db.QueryRow(
		"SELECT blogs.* FROM blogs WHERE id = ?", id,
	).Scan(
		&b.ID, &b.Slug, &b.Title, &b.Abstract, &b.Content,
		&b.UserId, &b.IsPublic, &b.CreatedAt, &b.UpdatedAt,
	)

	switch {
	case err == sql.ErrNoRows:
		b.ID = 0
	case err != nil:
		tools.ErrorLog(err)
	}
	return b
}

func DetailBlogWithUser(id int) TBlogJoinAnimesUser {
	db := connector.AccessDB()
	defer db.Close()

	var b TBlogJoinAnimesUser
	err := db.QueryRow(
		"SELECT blogs.* FROM blogs WHERE id = ?", id,
	).Scan(
		&b.ID, &b.Slug, &b.Title, &b.Abstract, &b.Content,
		&b.UserId, &b.IsPublic, &b.CreatedAt, &b.UpdatedAt,
	)

	switch {
	case err == sql.ErrNoRows:
		b.ID = 0
	case err != nil:
		tools.ErrorLog(err)
	}
	return b
}

func DetailBlogBySlug(slug string) TBlogJoinAnimes {
	db := connector.AccessDB()
	defer db.Close()

	var b TBlogJoinAnimes
	err := db.QueryRow("SELECT * FROM blogs WHERE slug = ?", slug).Scan(
		&b.ID, &b.Slug, &b.Title, &b.Abstract,
		&b.Content, &b.UserId, &b.IsPublic,
		&b.CreatedAt, &b.UpdatedAt,
	)

	switch {
	case err == sql.ErrNoRows:
		b.ID = 0
	case err != nil:
		tools.ErrorLog(err)
	}
	return b
}

func DetailBlogWithUserBySlug(slug string) TBlogJoinAnimesUser {
	db := connector.AccessDB()
	defer db.Close()

	var b TBlogJoinAnimesUser
	err := db.QueryRow("SELECT * FROM blogs WHERE slug = ?", slug).Scan(
		&b.ID, &b.Slug, &b.Title, &b.Abstract, &b.Content,
		&b.UserId, &b.IsPublic, &b.CreatedAt, &b.UpdatedAt,
	)

	switch {
	case err == sql.ErrNoRows:
		b.ID = 0
	case err != nil:
		tools.ErrorLog(err)
	}
	return b
}

func InsertBlog(title string, abstract string, content string, userId string, isPublic bool) int {
	db := connector.AccessDB()
	defer db.Close()

	stmt, err := db.Prepare(
		"INSERT INTO blogs(slug, title, abstract, content, user_id, is_public) " +
			"VALUES(?, ?, ?, ?, ?, ?)",
	)
	defer stmt.Close()

	slug := tools.GenRandSlug(8)
	exe, err := stmt.Exec(
		slug, title, abstract,
		content, userId, isPublic,
	)

	insertedId, err := exe.LastInsertId()
	if err != nil {
		tools.ErrorLog(err)
	}
	return int(insertedId)
}

// validation by userId @domain or view
func UpdateBlog(id int, title string, abstract string, content string, isPublic bool) int {
	db := connector.AccessDB()
	defer db.Close()

	stmt, err := db.Prepare("UPDATE blogs SET title = ?, abstract = ?, content = ?, is_public = ? WHERE id = ?")
	exe, err := stmt.Exec(
		title, abstract, content, id, isPublic,
	)
	defer stmt.Close()
	if err != nil {
		tools.ErrorLog(err)
	}
	updatedId, _ := exe.RowsAffected()
	return int(updatedId)
}

func DeleteBlog(id int) int {
	db := connector.AccessDB()
	defer db.Close()

	stmt, err := db.Prepare("DELETE FROM blogs WHERE id = ?")
	exe, err := stmt.Exec(id)
	if err != nil {
		tools.ErrorLog(err)
	}
	rowsAffect, _ := exe.RowsAffected()
	return int(rowsAffect)
}

// animes of blog
func RelationAnimeByBlog(blogId int) *sql.Rows {
	db := connector.AccessDB()
	defer db.Close()

	rows, err := db.Query(
		"SELECT relation_blog_animes.anime_id, animes.slug, animes.title FROM relation_blog_animes " +
			"LEFT JOIN animes ON animes.id = relation_blog_animes.anime_id " +
			"WHERE blog_id = " + strconv.Itoa(blogId),
	)
	if err != nil {
		tools.ErrorLog(err)
	}
	return rows
}

// blogs by anime
func RelationBlogByAnime(animeId int) *sql.Rows {
	db := connector.AccessDB()
	defer db.Close()

	rows, err := db.Query(
		"SELECT relation_blog_animes.blog_id, blogs.slug, blogs.title, " +
			"blogs.abstract, blogs.created_at FROM relation_blog_animes " +
			"LEFT JOIN blogs ON blogs.id = relation_blog_animes.blog_id " +
			"WHERE anime_id = " + strconv.Itoa(animeId),
	)
	if err != nil {
		tools.ErrorLog(err)
	}
	return rows
}

func InsertRelationAnimeBlog(animeId int, blogId int) bool {
	db := connector.AccessDB()
	defer db.Close()
	stmt, err := db.Prepare(
		"INSERT INTO relation_blog_animes(anime_id, blog_id) " +
			"VALUES(?, ?)",
	)
	defer stmt.Close()

	_, err = stmt.Exec(
		animeId, blogId,
	)
	if err != nil {
		return false
	}
	return true
}

func DeleteRelationAnimeBlog(animeId int, blogId int) error {
	db := connector.AccessDB()
	defer db.Close()

	stmt, err := db.Prepare("DELETE FROM relation_blog_animes WHERE anime_id = ? AND blog_id = ?")
	if err != nil {
		tools.ErrorLog(err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(animeId, blogId)
	if err != nil {
		return err
	}
	return nil
}
