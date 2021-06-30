package database_test

import (
	"animar/v1/pkg/domain"
	"animar/v1/pkg/infrastructure"
	"animar/v1/pkg/interfaces/database"
	"animar/v1/pkg/usecase"
	"database/sql"
	"fmt"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func newSqlHandler(sql *sql.DB) database.SqlHandler {
	sqlHandler := new(infrastructure.SqlHandler)
	sqlHandler.Conn = sql
	return sqlHandler
}

func TestInsert(t *testing.T) {
	b := &domain.TBlogInsert{
		Title:    "test sample1",
		Abstract: "",
		Content:  "test \nkoreha test desu",
		IsPublic: true,
		AnimeIds: []int{},
	}
	userId := "aaaaaaaa"
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("DB FAIL: %s", err)
	}
	query := "INSERT INTO blogs(slug, title, abstract, content, user_id, is_public) VALUES(?, ?, ?, ?, ?, ?)"
	stmt := mock.ExpectPrepare(regexp.QuoteMeta(query))
	stmt.ExpectExec().WithArgs("bbbbbbbb", b.Title, b.Abstract, b.Content, userId, b.IsPublic)

	sqlHandler := newSqlHandler(db)
	us := usecase.NewBlogInteractor(
		&database.BlogRepository{
			SqlHandler: sqlHandler,
		},
	)
	id, err := us.InsertBlog(*b, "aaaaaaaa")
	fmt.Print(id)
	t.Errorf("INSERT FAIL: %s", err)
}
