package infrastractures

import (
	"animar/v1/configs"
	"animar/v1/interfaces/database"
	"database/sql"
	"fmt"
)

type SqlHandler struct {
	Conn *sql.DB
}

type SqlResult struct {
	Result sql.Result
}

func NewSqlHandler() *SqlHandler {
	conn, err := sql.Open("mysql", fmt.Sprintf("%s:%s%s/%s", configs.MysqlUser, configs.MysqlPassword, configs.MysqlHost, configs.MysqlDataBase))
	defer conn.Close()
	if err != nil {
		panic(err.Error())
	}
	sqlHandler := new(SqlHandler)
	sqlHandler.Conn = conn
	return sqlHandler
}

// func (handler *SqlHandler) Execute(statement string, args ...interface{}) (database.Result, error) {
// 	res := SqlResult{}
// 	stmt, err := handler.Conn.Prepare(statement)
// 	defer stmt.Close()
// 	if err != nil {
// 		return nil, err
// 	}

// 	exe, err := stmt.Exec(args...)
// 	if err != nil {
// 		return exe, err
// 	}
// 	res.Result = exe
// 	return exe, nil
// }

func (handler *SqlHandler) Query(statement string, args ...interface{}) (database.Rows, error) {
	rows, err := handler.Conn.Query(statement, args...)
	if err != nil {
		return new(SqlRows), err
	}
	row := new(SqlRows)
	row.Rows = rows
	return rows, nil
}

func (handler *SqlHandler) QueryRow(statement string, args ...interface{}) database.Row {
	result := handler.Conn.QueryRow(statement, args...)
	row := new(SqlRow)
	row.Row = result
	return result
}

func (r SqlResult) LastInsertId() (int64, error) {
	return r.Result.LastInsertId()
}

func (r SqlResult) RowsAffected() (int64, error) {
	return r.Result.RowsAffected()
}

type SqlRows struct {
	Rows *sql.Rows
}

func (r SqlRows) Scan(args ...interface{}) error {
	return r.Rows.Scan(args...)
}

func (r SqlRows) Next() bool {
	return r.Rows.Next()
}

func (r SqlRows) Close() error {
	return r.Rows.Close()
}

type SqlRow struct {
	Row *sql.Row
}

func (r SqlRow) Scan(args ...interface{}) error {
	return r.Row.Scan(args...)
}
