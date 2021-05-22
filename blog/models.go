package blog

import (
	"animar/v1/helper"
	"database/sql"
	"fmt"
)

type TBlog struct {
	ID        int
	Slug      string
	Title     string
	Abstract  string
	Content   string
	UserId    string
	CreatedAt string
	UpdatedAt string
}

func ListBlog() *sql.Rows {
	db := helper.AccessDB()
	defer db.Close()
	rows, err := db.Query("Select * from tbl_blog")
	if err != nil {
		panic(err.Error())
	}
	return rows
}

func DetailBlog(id int) TBlog {
	db := helper.AccessDB()
	defer db.Close()

	var blog TBlog
	err := db.QueryRow("SELECT * FROM tbl_blog WHERE id = ?", id).Scan(
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

func DetailBlogBySlug(slug string) TBlog {
	db := helper.AccessDB()
	defer db.Close()

	var blog TBlog
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
	db := helper.AccessDB()
	defer db.Close()

	stmtInsert, err := db.Prepare(
		"INSERT INTO tbl_blog(slug, title, abstract, content, user_id) " +
			"VALUES(?, ?, ?, ?)",
	)
	defer stmtInsert.Close()

	slug := helper.GenRandSlug(8)
	exe, err := stmtInsert.Exec(
		slug, title, helper.NewNullString(abstract).String,
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
	db := helper.AccessDB()
	defer db.Close()

	exe, err := db.Exec(
		"UPDATE tbl_blog SET title = ?, abstract = ?, content = ? WHERE id = ?",
		title, helper.NewNullString(abstract).String, content, id,
	)
	if err != nil {
		fmt.Print(err)
	}
	updatedId, _ := exe.RowsAffected()
	return int(updatedId)
}

func DeleteBlog(id int) int {
	db := helper.AccessDB()
	defer db.Close()

	exe, err := db.Exec("DELETE FROM tbl_blog WHERE id = ?")
	if err != nil {
		panic(err.Error())
	}
	rowsAffect, _ := exe.RowsAffected()
	return int(rowsAffect)
}
