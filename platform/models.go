package platform

import (
	"animar/v1/tools/connector"
	"animar/v1/tools/tools"
	"database/sql"
)

type TPlatform struct {
	ID        int     `json:"id"`
	EngName   string  `json:"eng_name"`
	PlatName  *string `json:"plat_name"`
	BaseUrl   *string `json:"base_url"`
	Image     *string `json:"image"`
	IsValid   bool    `json:"is_valid"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt string  `json:"updated_at"`
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
	LinkUrl    *string `json:"link_url"`
	CreatedAt  *string `json:"created_at"`
	UpdatedAt  *string `json:"updated_at"`
	PlatName   *string `json:"plat_name"`
}

type TRelationPlatformInput struct {
	PlatformId int    `json:"platform_id"`
	AnimeId    int    `json:"anime_id"`
	LinkUrl    string `json:"link_url"`
}

func listPlatform() *sql.Rows {
	db := connector.AccessDB()
	defer db.Close()
	rows, err := db.Query("Select * from platforms")
	if err != nil {
		tools.ErrorLog(err)
	}
	return rows
}

func insertPlatform(engName string, platName string, baseUrl string, image string, isValid bool) int {
	db := connector.AccessDB()
	defer db.Close()

	stmt, err := db.Prepare(
		"INSERT INTO platforms(eng_name, plat_name, base_url, image, is_valid) VALUES(?, ?, ?, ?, ?)",
	)
	defer stmt.Close()
	exe, err := stmt.Exec(
		engName, tools.NewNullString(platName),
		tools.NewNullString(baseUrl), tools.NewNullString(image),
		isValid,
	)
	insertedId, err := exe.LastInsertId()
	if err != nil {
		tools.ErrorLog(err)
	}
	return int(insertedId)
}

func detailPlatfrom(id int) TPlatform {
	db := connector.AccessDB()
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
		tools.ErrorLog(err)
	}
	return p
}

// validation by userId @domain or view
func updatePlatform(engName string, platName string, baseUrl string, image string, isValid bool, id int) int {
	db := connector.AccessDB()
	defer db.Close()

	stmt, err := db.Prepare("UPDATE platforms SET eng_name = ?, plat_name = ?, base_url = ?, image = ?, is_valid = ? WHERE id = ?")
	defer stmt.Close()
	exe, err := stmt.Exec(
		engName, tools.NewNullString(platName),
		tools.NewNullString(baseUrl), tools.NewNullString(image),
		isValid, id,
	)

	if err != nil {
		tools.ErrorLog(err)
	}
	updatedId, _ := exe.RowsAffected()
	return int(updatedId)
}

func deletePlatform(id int) int {
	db := connector.AccessDB()
	defer db.Close()

	stmt, err := db.Prepare("DELETE FROM platforms WHERE id = ?")
	exe, err := stmt.Exec(id)
	defer stmt.Close()
	if err != nil {
		tools.ErrorLog(err)
	}
	rowsAffect, _ := exe.RowsAffected()
	return int(rowsAffect)
}

/****************************
*          relation		    *
****************************/

func relationPlatformByAnime(animeId int) *sql.Rows {
	db := connector.AccessDB()
	defer db.Close()
	rows, err := db.Query(
		"Select relation_anime_platform.*, platforms.plat_name FROM relation_anime_platform "+
			"LEFT JOIN platforms ON relation_anime_platform.platform_id = platforms.id "+
			"WHERE anime_id = ?", animeId,
	)
	if err != nil {
		tools.ErrorLog(err)
	}
	return rows
}

func insertRelation(platformId int, animeId int, linkUrl string) int {
	db := connector.AccessDB()
	defer db.Close()

	stmt, err := db.Prepare(
		"INSERT INTO relation_anime_platform(platform_id, anime_id, link_url) VALUES(?, ?, ?)",
	)
	defer stmt.Close()

	exe, err := stmt.Exec(
		platformId, animeId, tools.NewNullString(linkUrl),
	)
	insertedId, err := exe.LastInsertId()
	if err != nil {
		tools.ErrorLog(err)
	}
	return int(insertedId)
}

func deleteRelationPlatform(animeId int, platformId int) int {
	db := connector.AccessDB()
	defer db.Close()
	stmt, err := db.Prepare("DELETE FROM relation_anime_platform WHERE anime_id = ? AND platform_id = ?")
	defer stmt.Close()
	exe, err := stmt.Exec(
		animeId, platformId,
	)
	if err != nil {
		tools.ErrorLog(err)
	}
	rowsAffect, _ := exe.RowsAffected()
	return int(rowsAffect)
}
