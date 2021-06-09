package series

import (
	"animar/v1/tools"
	"database/sql"
	"fmt"
)

type TSeries struct {
	ID         int     `json:"id"`
	EngName    string  `json:"eng_name"`
	SeriesName *string `json:"series_name,omitempty"`
	CreatedAt  string  `json:"created_at"`
	UpdatedAt  *string `json:"updated_at,omitempty"`
}

type TSeriesInput struct {
	EngName    string `json:"eng_name"`
	SeriesName string `json:"series_name,omitempty"`
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

	stmt, err := db.Prepare(
		"INSERT INTO series(eng_name, series_name) VALUES(?, ?)",
	)
	defer stmt.Close()

	exe, err := stmt.Exec(
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

	stmt, err := db.Prepare("UPDATE series SET eng_name = ?, series_name = ? WHERE id = ?")
	defer stmt.Close()
	exe, err := stmt.Exec(
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

	stmt, err := db.Prepare("DELETE FROM series WHERE id = ?")
	defer db.Close()
	exe, err := stmt.Exec(id)
	if err != nil {
		panic(err.Error())
	}
	rowsAffect, _ := exe.RowsAffected()
	return int(rowsAffect)
}
