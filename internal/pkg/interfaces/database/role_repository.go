package database

import (
	"animar/v1/internal/pkg/domain"
	"log"
)

type RoleRepository struct {
	SqlHandler
}

func (ror *RoleRepository) FilterByAnime(animeId int) (rs []domain.AnimeStaffRole, err error) {
	rows, err := ror.Query(
		"SELECT r.num, r.role_name, s.family_name, s.given_name, s.eng_name "+
			"FROM anime_staff_roles AS asr "+
			"LEFT JOIN roles AS r ON r.id = asr.role_id "+
			"LEFT JOIN staffs AS s ON s.id = asr.staff_id "+
			"WHERE anime_id = ?",
		animeId,
	)
	if err != nil {
		log.Print(err)
		return
	}
	for rows.Next() {
		var r domain.AnimeStaffRole
		var givenName string
		rows.Scan(
			&r.Num, &r.Role, &r.Name, &givenName, &r.EngName,
		)
		r.Name += givenName
		rs = append(rs, r)
	}
	return
}

// admin

func (ror *RoleRepository) List() (rs []domain.Role, err error) {
	rows, err := ror.Query(
		"SELECT id, num, role_name, created_at, updated_at " +
			"FROM roles",
	)
	if err != nil {
		log.Print(err)
		return
	}
	for rows.Next() {
		var r domain.Role
		rows.Scan(
			&r.ID, &r.Num, &r.Role, &r.CreatedAt, &r.UpdatedAt,
		)
		rs = append(rs, r)
	}
	return
}

func (ror *RoleRepository) Insert(r domain.RoleInput) (inserted int, err error) {
	exe, err := ror.Execute(
		"INSERT INTO roles(num, role_name) VALUES(?, ?)",
		r.Num, r.Role,
	)
	if err != nil {
		log.Print(err)
		return
	}
	rawInserted, _ := exe.LastInsertId()
	inserted = int(rawInserted)
	return
}

func (ror *RoleRepository) InsertStaffRole(r domain.AnimeStaffRoleInput) (inserted int, err error) {
	exe, err := ror.Execute(
		"INSERT INTO anime_staff_roles(anime_id, role_id, staff_id) "+
			"VALUES (?, ?, ?)",
		r.AnimeId, r.RoleId, r.StaffId,
	)
	if err != nil {
		log.Print(err)
		return
	}
	rawInserted, _ := exe.LastInsertId()
	inserted = int(rawInserted)
	return
}
