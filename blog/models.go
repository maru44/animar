package blog

import (
	"animar/v1/tools"
	"database/sql"
	"fmt"
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
	db := tools.AccessDB()
	defer db.Close()
	rows, err := db.Query("Select * from blogs")
	if err != nil {
		panic(err.Error())
	}
	return rows
}

func ListUsersBlog(uid string) *sql.Rows {
	db := tools.AccessDB()
	defer db.Close()
	rows, err := db.Query("SELECT * from blogs where user_id = ?", uid)
	if err != nil {
		panic(err.Error())
	}
	return rows
}

// 使ってない Many to many のjoinに失敗
func ListBlogJoinAnime() *sql.Rows {
	db := tools.AccessDB()
	defer db.Close()
	// fail
	rows, err := db.Query(
		"SELECT blogs.*, animes FROM blogs " +
			"LEFT JOIN relation_blog_animes ON blog.id = relation_blog_animes.blog_id " +
			"LEFT JOIN animes ON anime.id = relation_blog_animes.anime_id",
	)
	if err != nil {
		panic(err.Error())
	}
	return rows
}

func BlogUserId(id int) string {
	db := tools.AccessDB()
	defer db.Close()
	var userId string
	err := db.QueryRow("SELECT user_id FROM blogs WHERE id = ?", id).Scan(&userId)
	if err != nil {
		fmt.Print(err)
	}
	return userId
}

func DetailBlog(id int) TBlogJoinAnimes {
	db := tools.AccessDB()
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
		panic(err.Error())
	}
	return b
}

func DetailBlogWithUser(id int) TBlogJoinAnimesUser {
	db := tools.AccessDB()
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
		panic(err.Error())
	}
	return b
}

func DetailBlogBySlug(slug string) TBlogJoinAnimes {
	db := tools.AccessDB()
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
		panic(err.Error())
	}
	return b
}

func DetailBlogWithUserBySlug(slug string) TBlogJoinAnimesUser {
	db := tools.AccessDB()
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
		panic(err.Error())
	}
	return b
}

func InsertBlog(title string, abstract string, content string, userId string) int {
	db := tools.AccessDB()
	defer db.Close()

	stmt, err := db.Prepare(
		"INSERT INTO blogs(slug, title, abstract, content, user_id) " +
			"VALUES(?, ?, ?, ?, ?)",
	)
	defer stmt.Close()

	slug := tools.GenRandSlug(8)
	exe, err := stmt.Exec(
		slug, title, abstract,
		content, userId,
	)

	insertedId, err := exe.LastInsertId()
	if err != nil {
		fmt.Print(err)
	}
	return int(insertedId)
}

// validation by userId @domain or view
func UpdateBlog(id int, title string, abstract string, content string) int {
	db := tools.AccessDB()
	defer db.Close()

	stmt, err := db.Prepare("UPDATE blogs SET title = ?, abstract = ?, content = ? WHERE id = ?")
	exe, err := stmt.Exec(
		title, abstract, content, id,
	)
	defer stmt.Close()
	if err != nil {
		fmt.Print(err)
	}
	updatedId, _ := exe.RowsAffected()
	return int(updatedId)
}

func DeleteBlog(id int) int {
	db := tools.AccessDB()
	defer db.Close()

	stmt, err := db.Prepare("DELETE FROM blogs WHERE id = ?")
	exe, err := stmt.Exec(id)
	if err != nil {
		panic(err.Error())
	}
	rowsAffect, _ := exe.RowsAffected()
	return int(rowsAffect)
}

// animes of blog
func RelationAnimeByBlog(blogId int) *sql.Rows {
	db := tools.AccessDB()
	defer db.Close()

	rows, err := db.Query(
		"SELECT relation_blog_animes.anime_id, animes.slug, animes.title FROM relation_blog_animes " +
			"LEFT JOIN animes ON animes.id = relation_blog_animes.anime_id " +
			"WHERE blog_id = " + strconv.Itoa(blogId),
	)
	if err != nil {
		panic(err.Error())
	}
	return rows
}

// blogs by anime
func RelationBlogByAnime(animeId int) *sql.Rows {
	db := tools.AccessDB()
	defer db.Close()

	rows, err := db.Query(
		"SELECT relation_blog_animes.blog_id, blogs.slug, blogs.title, " +
			"blogs.abstract, blogs.created_at FROM relation_blog_animes " +
			"LEFT JOIN blogs ON blogs.id = relation_blog_animes.blog_id " +
			"WHERE anime_id = " + strconv.Itoa(animeId),
	)
	if err != nil {
		panic(err.Error())
	}
	return rows
}

func InsertRelationAnimeBlog(animeId int, blogId int) bool {
	db := tools.AccessDB()
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
	db := tools.AccessDB()
	defer db.Close()

	stmt, err := db.Prepare("DELETE FROM relation_blog_animes WHERE anime_id = ? AND blog_id = ?")
	if err != nil {
		panic(err.Error())
	}
	defer stmt.Close()

	_, err = stmt.Exec(animeId, blogId)
	if err != nil {
		return err
	}
	return nil
}
