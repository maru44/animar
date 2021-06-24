package database

type SqlHandler interface {
	Query(string, ...interface{}) (Rows, error)
	Execute(string, ...interface{}) (Result, error)
	ErrNoRows() error
}

type Rows interface {
	Scan(...interface{}) error
	Next() bool
	Close() error
}

type Result interface {
	LastInsertId() (int64, error)
	RowsAffected() (int64, error)
}
