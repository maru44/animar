package repository

import (
	"database/sql"
)

type BlogRepository interface {
	Insert(DB *sql.DB)
}
