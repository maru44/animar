package database

type SqlHandler interface {
	Query(string, ...interface{}) (Rows, error)
	// Prepare(string, ...interface{}) (Statement, error)
}

type Rows interface {
	Scan(...interface{}) error
	Next() bool
	Close() error
}
