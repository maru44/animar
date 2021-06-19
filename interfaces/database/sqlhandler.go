package database

type SqlHandler interface {
	Query(string, ...interface{}) (Rows, error)
	//QueryRow(string, ...interface{}) (Row, error)
	//Prepare(string, ...interface{}) (Statement, error)
}

type Statement interface {
	Exec(...interface{}) (Result, error)
}

type Result interface {
	LastInsertId() (int64, error)
	RowsAffected() (int64, error)
	// DeletedRow() (int64, error)
}

type Row interface {
	Scan(...interface{}) error
}

type Rows interface {
	Row
	Next() bool
	Close() error
}
