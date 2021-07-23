package database

import (
	"animar/v1/internal/pkg/domain"
	"log"
)

type StaffRepository struct {
	SqlHandler
}

func (sfr *StaffRepository) List() (sfs []domain.Staff, err error) {
	rows, err := sfr.Query(
		"SELECT id, eng_name, family_name, given_name, created_at, updated_at " +
			"FROM staffs",
	)
	if err != nil {
		log.Print(err)
	}
	for rows.Next() {
		var sf domain.Staff
		rows.Scan(
			&sf.ID, &sf.EngName, &sf.FamilyName, &sf.GivenName, &sf.CreatedAt, &sf.UpdatedAt,
		)
		sfs = append(sfs, sf)
	}
	return
}

// admin

func (sfr *StaffRepository) Insert(sf domain.StaffInput) (inserted int, err error) {
	exe, err := sfr.Execute(
		"INSERT INTO staffs(eng_name, family_name, given_name) VALUES(?, ?, ?)",
		sf.EngName, sf.FamilyName, sf.GivenName,
	)
	if err != nil {
		log.Print(err)
		return
	}
	rawInserted, _ := exe.LastInsertId()
	inserted = int(rawInserted)
	return
}
