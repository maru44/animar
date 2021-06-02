package series

import (
	"animar/v1/tools"
	"database/sql"
	"fmt"
)

type TSeries struct {
	ID         int
	EngName    string
	SeriesName *string
	CreatedAt  string
	UpdatedAt  *string
}

type TSeriesInput struct {
	EngName    string `json:"EngName"`
	SeriesName string `json:"SeriesName,omitempty"`
}

func ListSeries() *sql.Rows {
	db := tools.AccessDB()
	defer db.Close()
	rows, err := db.Query("Select * from series")
	if err != nil {
		panic(err.Error())
	}
	return rows
}

func InsertSeries(engName string, seriesName string) int {
	db := tools.AccessDB()
	defer db.Close()

	stmtInsert, err := db.Prepare(
		"INSERT INTO anime(eng_name, series_name) VALUES(?, ?)",
	)
	defer stmtInsert.Close()

	exe, err := stmtInsert.Exec(
		engName, seriesName,
	)
	insertedId, err := exe.LastInsertId()
	if err != nil {
		fmt.Print(err)
	}

	return int(insertedId)
}
