package database

type SqlHandler interface {
	Query(string, ...interface{}) (Rows, error)
	// Prepare(string, ...interface{}) (Statement, error)
	Execute(string, ...interface{}) (Result, error)
}

// type Statement interface {
// 	Exec(...interface{}) (Result, error)
// 	Close() error
// }

type Rows interface {
	Scan(...interface{}) error
	Next() bool
	Close() error
}

type Result interface {
	LastInsertId() (int64, error)
	RowsAffected() (int64, error)
}
