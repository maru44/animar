package platform

import (
	"animar/v1/tools"
	"database/sql"
	"fmt"
)

type TPlatform struct {
	ID        int     `json:"id"`
	EngName   string  `json:"eng_name"`
	PlatName  *string `json:"plat_name,omitempty"`
	BaseUrl   *string `json:"base_url,omitempty"`
	Image     *string `json:"image,omitempty"`
	IsValid   bool    `json:"is_valid"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt string  `json:"updated_at,omitempty"`
}

// @NOTUSED
type TPlatformInput struct {
	EngName  string  `json:"eng_name"`
	PlatName *string `json:"plat_name"`
	BaseUrl  *string `json:"base_url"`
	Image    *string `json:"image"`
	IsValid  *bool   `json:"is_valid"`
}

type TRelationPlatform struct {
	PlatformId int     `json:"platform_id"`
	AnimeId    int     `json:"anime_id"`
	LinkUrl    *string `json:"link_url,omitempty"`
	CreatedAt  *string `json:"created_at"`
	UpdatedAt  *string `json:"updated_at,omitempty"`
}

type TRelationPlatformInput struct {
	PlatformId int    `json:"platform_id"`
	AnimeId    int    `json:"anime_id"`
	LinkUrl    string `json:"link_url,omitempty"`
}

func ListPlatform() *sql.Rows {
	db := tools.AccessDB()
	defer db.Close()
	rows, err := db.Query("Select * from platforms")
	if err != nil {
		panic(err.Error())
	}
	return rows
}

func InsertPlatform(engName string, platName string, baseUrl string, image string, isValid bool) int {
	db := tools.AccessDB()
	defer db.Close()

	stmtInsert, err := db.Prepare(
		"INSERT INTO platforms(eng_name, plat_name, base_url, image, is_valid) VALUES(?, ?, ?, ?, ?)",
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

	var p TPlatform
	err := db.QueryRow(
		"SELECT * FROM platforms WHERE id = ?", id,
	).Scan(
		&p.ID, &p.EngName, &p.PlatName, &p.BaseUrl,
		&p.Image, &p.IsValid, &p.CreatedAt, &p.UpdatedAt,
	)

	switch {
	case err == sql.ErrNoRows:
		p.ID = 0
	case err != nil:
		panic(err.Error())
	}
	return p
}

// validation by userId @domain or view
func UpdatePlatform(engName string, platName string, baseUrl string, image string, isValid bool, id int) int {
	db := tools.AccessDB()
	defer db.Close()

	exe, err := db.Exec(
		"UPDATE platforms SET eng_name = ?, plat_name = ?, base_url = ?, image = ?, is_valid = ? WHERE id = ?",
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

	exe, err := db.Exec("DELETE FROM platforms WHERE id = ?", id)
	if err != nil {
		panic(err.Error())
	}
	rowsAffect, _ := exe.RowsAffected()
	return int(rowsAffect)
}

/****************************
*          relation		    *
****************************/

func RelationPlatformByAnime(animeId int) *sql.Rows {
	db := tools.AccessDB()
	defer db.Close()
	rows, err := db.Query(
		"Select * from relation_anime_platform where anime_id = ?",
		animeId,
	)
	if err != nil {
		panic(err.Error())
	}
	return rows
}

func InsertRelation(platformId int, animeId int, linkUrl string) int {
	db := tools.AccessDB()
	defer db.Close()

	stmtInsert, err := db.Prepare(
		"INSERT INTO relation_anime_platform(platform_id, anime_id, link_url) VALUES(?, ?, ?)",
	)
	defer stmtInsert.Close()

	exe, err := stmtInsert.Exec(
		platformId, animeId, tools.NewNullString(linkUrl),
	)
	insertedId, err := exe.LastInsertId()
	if err != nil {
		fmt.Print(err)
	}
	return int(insertedId)
}

func DeleteRelationPlatform(animeId int, platformId int) int {
	db := tools.AccessDB()
	defer db.Close()
	exe, err := db.Exec(
		"DELETE FROM relation_anime_platform WHERE anime_id = ? AND platform_id = ?",
		animeId, platformId,
	)
	if err != nil {
		panic(err.Error())
	}
	rowsAffect, _ := exe.RowsAffected()
	return int(rowsAffect)
}
