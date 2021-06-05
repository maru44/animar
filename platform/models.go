package platform

import (
	"animar/v1/tools"
	"database/sql"
	"fmt"
)

type TPlatform struct {
	ID        int     `json:"ID"`
	EngName   string  `json:"EngName"`
	PlatName  *string `json:"PlatName"`
	BaseUrl   *string `json:"BaseUrl"`
	Image     *string `json:"Image"`
	IsValid   bool    `json:"IsValid"`
	CreatedAt string  `json:"CreatedAt"`
	UpdatedAt string  `json:"UpdatedAt"`
}

type TPlatformInput struct {
	EngName  string  `json:"EngName"`
	PlatName *string `json:"PlatName"`
	BaseUrl  *string `json:"BaseUrl"`
	Image    *string `json:"Image"`
	IsValid  *bool   `json:"IsValid"`
}

func ListPlatform() *sql.Rows {
	db := tools.AccessDB()
	defer db.Close()
	rows, err := db.Query("Select * from platform")
	if err != nil {
		panic(err.Error())
	}
	return rows
}

func InsertPlatform(engName string, platName string, baseUrl string, image string, isValid bool) int {
	db := tools.AccessDB()
	defer db.Close()

	stmtInsert, err := db.Prepare(
		"INSERT INTO platform(eng_name, plat_name, base_url, image, is_valid) VALUES(?, ?, ?, ?, ?)",
	)
	defer stmtInsert.Close()

	exe, err := stmtInsert.Exec(
		engName, tools.NewNullString(platName),
		tools.NewNullString(baseUrl), tools.NewNullString(image),
		isValid,
	)
	insertedId, err := exe.LastInsertId()
	if err != nil {
		fmt.Print(err)
	}
	return int(insertedId)
}

func DetailPlatfrom(id int) TPlatform {
	db := tools.AccessDB()
	defer db.Close()

	var plat TPlatform
	err := db.QueryRow(
		"SELECT * FROM platform WHERE id = ?", id,
	).Scan(
		&plat.ID, &plat.EngName, &plat.PlatName, &plat.BaseUrl,
		&plat.Image, &plat.IsValid, &plat.CreatedAt, &plat.UpdatedAt,
	)

	switch {
	case err == sql.ErrNoRows:
		plat.ID = 0
	case err != nil:
		panic(err.Error())
	}
	return plat
}

// validation by userId @domain or view
func UpdatePlatform(engName string, platName string, baseUrl string, image string, isValid bool, id int) int {
	db := tools.AccessDB()
	defer db.Close()

	exe, err := db.Exec(
		"UPDATE platform SET eng_name = ?, plat_name = ?, base_url = ?, image = ?, is_valid = ? WHERE id = ?",
		engName, tools.NewNullString(platName),
		tools.NewNullString(baseUrl), tools.NewNullString(image),
		isValid, id,
	)
	if err != nil {
		fmt.Print(err)
	}
	updatedId, _ := exe.RowsAffected()
	return int(updatedId)
}

func DeletePlatform(id int) int {
	db := tools.AccessDB()
	defer db.Close()

	exe, err := db.Exec("DELETE FROM platform WHERE id = ?", id)
	if err != nil {
		panic(err.Error())
	}
	rowsAffect, _ := exe.RowsAffected()
	return int(rowsAffect)
}
