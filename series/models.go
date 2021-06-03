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
		"INSERT INTO series(eng_name, series_name) VALUES(?, ?)",
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

func DetailSeries(id int) TSeries {
	db := tools.AccessDB()
	defer db.Close()

	var ser TSeries
	err := db.QueryRow(
		"SELECT * FROM series WHERE id = ?", id,
	).Scan(
		&ser.ID, &ser.EngName, &ser.SeriesName, &ser.CreatedAt, &ser.UpdatedAt,
	)

	switch {
	case err == sql.ErrNoRows:
		ser.ID = 0
	case err != nil:
		panic(err.Error())
	}
	return ser
}

// validation by userId @domain or view
func UpdateSeries(engName string, seriesName string, id int) int {
	db := tools.AccessDB()
	defer db.Close()

	exe, err := db.Exec(
		"UPDATE series SET eng_name = ?, series_name = ? WHERE id = ?",
		engName, seriesName, id,
	)
	if err != nil {
		fmt.Print(err)
	}
	updatedId, _ := exe.RowsAffected()
	return int(updatedId)
}

func DeleteSeries(id int) int {
	db := tools.AccessDB()
	defer db.Close()

	exe, err := db.Exec("DELETE FROM series WHERE id = ?", id)
	if err != nil {
		panic(err.Error())
	}
	rowsAffect, _ := exe.RowsAffected()
	return int(rowsAffect)
}
