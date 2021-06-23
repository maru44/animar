package infrastructure

import (
	"animar/v1/configs"
	"animar/v1/pkg/interfaces/database"
	"animar/v1/pkg/tools/tools"
	"database/sql"
	"fmt"
)

type SqlHandler struct {
	Conn *sql.DB
}

type SqlRows struct {
	Rows *sql.Rows
}

func NewSqlHandler() database.SqlHandler {
	conn, err := sql.Open("mysql", fmt.Sprintf("%s:%s%s/%s", configs.MysqlUser, configs.MysqlPassword, configs.MysqlHost, configs.MysqlDataBase))
	if err != nil {
		panic(err.Error())
	}
	sqlHandler := new(SqlHandler)
	sqlHandler.Conn = conn
	return sqlHandler
}

func (handler *SqlHandler) ErrNoRows() error {
	return sql.ErrNoRows
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

func (r SqlRows) Scan(dest ...interface{}) error {
	return r.Rows.Scan(dest...)
}

func (r SqlRows) Next() bool {
	return r.Rows.Next()
}

func (r SqlRows) Close() error {
	return r.Rows.Close()
}

type SqlResult struct {
	Result sql.Result
}

func (handler *SqlHandler) Execute(statement string, args ...interface{}) (database.Result, error) {
	res := SqlResult{}
	stmt, err := handler.Conn.Prepare(statement)
	defer stmt.Close()
	if err != nil {
		return res, err
	}
	exe, err := stmt.Exec(args...)
	if err != nil {
		tools.ErrorLog(err)
	}
	res.Result = exe
	return res, nil
}

func (r SqlResult) LastInsertId() (int64, error) {
	return r.Result.LastInsertId()
}

func (r SqlResult) RowsAffected() (int64, error) {
	return r.Result.RowsAffected()
}
