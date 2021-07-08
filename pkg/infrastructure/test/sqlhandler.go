package infrastructure_test

import (
	"animar/v1/pkg/infrastructure"
	"animar/v1/pkg/interfaces/database"
	"database/sql"
)

func NewDummyHandler(db *sql.DB) database.SqlHandler {
	sqlHandler := new(infrastructure.SqlHandler)
	sqlHandler.Conn = db
	return sqlHandler
}
