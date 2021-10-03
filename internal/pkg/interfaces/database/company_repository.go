package database

import (
	"animar/v1/internal/pkg/domain"

	"github.com/maru44/perr"
)

type CompanyRepository struct {
	SqlHandler
}

func (cr *CompanyRepository) List() (cs []domain.Company, err error) {
	rows, err := cr.Query(
		"SELECT id, name, eng_name, official_url, created_at, updated_at " +
			"FROM companies",
	)
	if err != nil {
		return cs, perr.Wrap(err, perr.InternalServerErrorWithUrgency)
	}
	defer rows.Close()

	for rows.Next() {
		var c domain.Company
		if err := rows.Scan(
			&c.ID, &c.Name, &c.EngName, &c.OfficialUrl, &c.CreatedAt, &c.UpdatedAt,
		); err != nil {
			return cs, perr.Wrap(err, perr.NotFound)
		}
		cs = append(cs, c)
	}
	return
}

func (cr *CompanyRepository) DetailByEng(engName string) (c domain.CompanyDetail, err error) {
	rows, err := cr.Query(
		"SELECT id, name, eng_name, official_url, explanation, twitter_account, created_at, updated_at "+
			"FROM companies "+
			"WHERE eng_name = ?",
		engName,
	)
	if err != nil {
		return c, perr.Wrap(err, perr.InternalServerErrorWithUrgency)
	}
	defer rows.Close()

	rows.Next()
	err = rows.Scan(
		&c.ID, &c.Name, &c.EngName, &c.OfficialUrl, &c.Explanation, &c.TwitterAccount, &c.CreatedAt, &c.UpdatedAt,
	)
	return c, perr.Wrap(err, perr.NotFound)
}

// this is for admin
func (cr *CompanyRepository) Insert(c domain.CompanyInput) (inserted int, err error) {
	exe, err := cr.Execute(
		"INSERT INTO companies(name, eng_name, official_url, explanation, twitter_account) "+
			"VALUES(?, ?, ?, ?, ?)",
		c.Name, c.EngName, c.OfficialUrl, c.Explanation, c.TwitterAccount,
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

func (cr *CompanyRepository) Update(c domain.CompanyInput, engName string) (affected int, err error) {
	exe, err := cr.Execute(
		"UPDATE companies SET name = ?, eng_name = ?, official_url = ?, explanation = ?, twitter_account = ? "+
			"WHERE eng_name = ?",
		c.Name, c.EngName, c.OfficialUrl, c.Explanation, c.TwitterAccount, engName,
	)
	if err != nil {
		return affected, perr.Wrap(err, perr.InternalServerErrorWithUrgency)
	}
	rawAffected, err := exe.RowsAffected()
	if err != nil {
		return affected, perr.Wrap(err, perr.BadRequest)
	}
	affected = int(rawAffected)
	return
}
