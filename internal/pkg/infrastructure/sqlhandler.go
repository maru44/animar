package infrastructure

import (
	"animar/v1/internal/configs"
	"animar/v1/internal/pkg/interfaces/database"
	"animar/v1/internal/pkg/tools/tools"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

type SqlHandler struct {
	Conn *sql.DB
}

type SqlRows struct {
	Rows *sql.Rows
}

type SqlResult struct {
	Result sql.Result
}

type SqlTransaction struct {
	Tx *sql.Tx
}

// init

func NewSqlHandler() database.SqlHandler {
	conn, err := sql.Open("mysql", fmt.Sprintf("%s:%s%s/%s", configs.MysqlUser, configs.MysqlPassword, configs.MysqlHost, configs.MysqlDataBase))
	if err != nil {
		panic(err.Error())
	}
	sqlHandler := new(SqlHandler)
	sqlHandler.Conn = conn
	return sqlHandler
}

/********************
    sqlHandler methods
**************************/

func (handler *SqlHandler) ErrNoRows() error {
	return handler.ErrNoRows()
}

func (handler *SqlHandler) Query(statement string, args ...interface{}) (database.Rows, error) {
	rows, err := handler.Conn.Query(statement, args...)
	if err != nil {
		return new(SqlRows), err
	}
	row := new(SqlRows)
	row.Rows = rows
	return rows, nil
}

func (handler *SqlHandler) Execute(statement string, args ...interface{}) (database.Result, error) {
	res := SqlResult{}
	stmt, err := handler.Conn.Prepare(statement)
	defer stmt.Close()
	if err != nil {
		tools.ErrorLog(err)
		return res, err
	}
	exe, err := stmt.Exec(args...)
	if err != nil {
		tools.ErrorLog(err)
	}
	res.Result = exe
	return res, nil
}

// begin transaction
func (handler *SqlHandler) Begin() (database.Tx, error) {
	res := SqlTransaction{}
	transaction, err := handler.Conn.Begin()
	if err != nil {
		log.Print(err)
	}
	res.Tx = transaction
	return res, err
}

/***************************
         transaction
***************************/

func (t SqlTransaction) Commit() error {
	return t.Tx.Commit()
}

func (t SqlTransaction) Rollback() error {
	return t.Tx.Rollback()
}

func (t SqlTransaction) Execute(statement string, args ...interface{}) (database.Result, error) {
	res := SqlResult{}
	stmt, err := t.Tx.Prepare(statement)
	defer stmt.Close()
	if err != nil {
		log.Print(err)
		return res, err
	}
	exe, err := stmt.Exec(args...)
	if err != nil {
		log.Print(err)
	}
	res.Result = exe
	return res, nil
}

func (t SqlTransaction) Query(statement string, args ...interface{}) (database.Rows, error) {
	rows, err := t.Tx.Query(statement, args...)
	if err != nil {
		return new(SqlRows), err
	}
	row := new(SqlRows)
	row.Rows = rows
	return rows, nil
}

/********************
    Rows methods
**************************/

func (r SqlRows) Scan(dest ...interface{}) error {
	return r.Rows.Scan(dest...)
}

func (r SqlRows) Next() bool {
	return r.Rows.Next()
}

func (r SqlRows) Close() error {
	return r.Rows.Close()
}

/********************
    Results methods
**************************/

func (r SqlResult) LastInsertId() (int64, error) {
	return r.Result.LastInsertId()
}

func (r SqlResult) RowsAffected() (int64, error) {
	return r.Result.RowsAffected()
}
