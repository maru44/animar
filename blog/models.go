package blog

import (
	"animar/v1/tools"
	"database/sql"
	"fmt"
	"strconv"

	"firebase.google.com/go/v4/auth"
)

type TBlog struct {
	ID        int
	Slug      string
	Title     string
	Abstract  *string
	Content   string
	UserId    string
	CreatedAt string
	UpdatedAt string
}

type TJoinedAnime struct {
	AnimeId int
	Slug    string
	Title   string
}

type TJoinedBlog struct {
	BlogId    int
	Title     string
	Slug      string
	Abstract  string
	CreatedAt string
}

type TBlogJoinAnimes struct {
	ID        int
	Slug      string
	Title     string
	Abstract  *string
	Content   string
	UserId    string
	CreatedAt string
	UpdatedAt string
	Animes    []TJoinedAnime
}

type TBlogJoinAnimesUser struct {
	ID        int
	Slug      string
	Title     string
	Abstract  *string
	Content   string
	UserId    string
	CreatedAt string
	UpdatedAt string
	Animes    []TJoinedAnime
	User      *auth.UserInfo
}

func ListBlog() *sql.Rows {
	db := tools.AccessDB()
	defer db.Close()
	rows, err := db.Query("Select * from tbl_blog")
	if err != nil {
		panic(err.Error())
	}
	return rows
}

func ListUsersBlog(uid string) *sql.Rows {
	db := tools.AccessDB()
	defer db.Close()
	rows, err := db.Query("SELECT * from tbl_blog where user_id = ?", uid)
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
		"SELECT tbl_blog.*, anime.id FROM tbl_blog " +
			"LEFT JOIN relation_blog_animes ON tbl_blog.id = relation_blog_animes.blog_id " +
			"LEFT JOIN anime ON anime.id = relation_blog_animes.anime_id",
	)
	if err != nil {
		panic(err.Error())
	}
	return rows
}

func DetailBlog(id int) TBlog {
	db := tools.AccessDB()
	defer db.Close()

	var blog TBlog
	err := db.QueryRow(
		"SELECT tbl_blog.*, anime.id, anime.slug, anime.title FROM tbl_blog "+
			"WHERE id = ?", id,
	).Scan(
		&blog.ID, &blog.Slug, &blog.Title, &blog.Abstract, &blog.Content,
		&blog.CreatedAt, &blog.UserId, &blog.UpdatedAt,
	)

	switch {
	case err == sql.ErrNoRows:
		blog.ID = 0
	case err != nil:
		panic(err.Error())
	}
	return blog
}

func DetailBlogBySlug(slug string) TBlogJoinAnimes {
	db := tools.AccessDB()
	defer db.Close()

	var blog TBlogJoinAnimes
	err := db.QueryRow("SELECT * FROM tbl_blog WHERE slug = ?", slug).Scan(
		&blog.ID, &blog.Slug, &blog.Title, &blog.Abstract,
		&blog.Content, &blog.UserId, &blog.CreatedAt,
		&blog.UpdatedAt,
	)

	switch {
	case err == sql.ErrNoRows:
		blog.ID = 0
	case err != nil:
		panic(err.Error())
	}
	return blog
}

func DetailBlogWithUserBySlug(slug string) TBlogJoinAnimesUser {
	db := tools.AccessDB()
	defer db.Close()

	var blog TBlogJoinAnimesUser
	err := db.QueryRow("SELECT * FROM tbl_blog WHERE slug = ?", slug).Scan(
		&blog.ID, &blog.Slug, &blog.Title, &blog.Abstract,
		&blog.Content, &blog.UserId, &blog.CreatedAt,
		&blog.UpdatedAt,
	)

	switch {
	case err == sql.ErrNoRows:
		blog.ID = 0
	case err != nil:
		panic(err.Error())
	}
	return blog
}

func InsertBlog(title string, abstract string, content string, userId string) int {
	db := tools.AccessDB()
	defer db.Close()

	stmtInsert, err := db.Prepare(
		"INSERT INTO tbl_blog(slug, title, abstract, content, user_id) " +
			"VALUES(?, ?, ?, ?, ?)",
	)
	defer stmtInsert.Close()

	slug := tools.GenRandSlug(8)
	exe, err := stmtInsert.Exec(
		slug, title, tools.NewNullString(abstract).String,
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

	exe, err := db.Exec(
		"UPDATE tbl_blog SET title = ?, abstract = ?, content = ? WHERE id = ?",
		title, tools.NewNullString(abstract).String, content, id,
	)
	if err != nil {
		fmt.Print(err)
	}
	updatedId, _ := exe.RowsAffected()
	return int(updatedId)
}

func DeleteBlog(id int) int {
	db := tools.AccessDB()
	defer db.Close()

	exe, err := db.Exec("DELETE FROM tbl_blog WHERE id = ?")
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
		"SELECT relation_blog_animes.anime_id, anime.slug, anime.title FROM relation_blog_animes " +
			"LEFT JOIN anime ON anime.id = relation_blog_animes.anime_id " +
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
		"SELECT relation_blog_animes.blog_id, blog.slug, blog.title, " +
			"blog.abstract, blog.created_at FROM relation_blog_animes " +
			"LEFT JOIN blog ON blog.id = relation_blog_animes.blog_id " +
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
	stmtInsert, err := db.Prepare(
		"INSERT INTO relation_blog_animes(anime_id, blog_id) " +
			"VALUES(?, ?)",
	)
	defer stmtInsert.Close()

	_, err = stmtInsert.Exec(
		animeId, blogId,
	)
	if err != nil {
		return false
	}
	return true
}
