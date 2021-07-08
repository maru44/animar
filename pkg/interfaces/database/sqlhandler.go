package database

type SqlHandler interface {
	Query(string, ...interface{}) (Rows, error)
	Execute(string, ...interface{}) (Result, error)
	ErrNoRows() error
	Begin() (Tx, error)
}

// transaction
type Tx interface {
	Commit() error
	Execute(string, ...interface{}) (Result, error)
	Rollback() error
	Query(string, ...interface{}) (Rows, error)
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
