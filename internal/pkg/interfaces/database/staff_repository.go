package database

import (
	"animar/v1/internal/pkg/domain"

	"github.com/maru44/perr"
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
		return sfs, perr.Wrap(err, perr.InternalServerErrorWithUrgency)
	}
	defer rows.Close()

	for rows.Next() {
		var sf domain.Staff
		err = rows.Scan(
			&sf.ID, &sf.EngName, &sf.FamilyName, &sf.GivenName, &sf.CreatedAt, &sf.UpdatedAt,
		)
		if err != nil {
			return sfs, perr.Wrap(err, perr.NotFound)
		}
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
		return inserted, perr.Wrap(err, perr.InternalServerErrorWithUrgency)
	}

	rawInserted, err := exe.LastInsertId()
	if err != nil {
		return inserted, perr.Wrap(err, perr.BadRequest)
	}
	inserted = int(rawInserted)
	return
}
