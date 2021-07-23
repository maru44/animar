package database

import (
	"animar/v1/internal/pkg/domain"
	"log"
)

type CompanyRepository struct {
	SqlHandler
}

// this is for admin
func (cr *CompanyRepository) Insert(c domain.CompanyInput) (inserted int, err error) {
	exe, err := cr.Execute(
		"INSERT INTO companies(name, eng_name, official_url) VALUES(?, ?, ?)",
		c.Name, c.EngName, c.OfficialUrl,
	)
	if err != nil {
		log.Print(err)
		return
	}
	rawInserted, _ := exe.LastInsertId()
	inserted = int(rawInserted)
	return
}

func (cr *CompanyRepository) List() (cs []domain.Company, err error) {
	rows, err := cr.Query(
		"SELECT id, name, eng_name, official_url, created_at, updated_at " +
			"FROM companies",
	)
	if err != nil {
		log.Print(err)
		return
	}
	for rows.Next() {
		var c domain.Company
		if err := rows.Scan(
			&c.ID, &c.Name, &c.EngName, &c.OfficialUrl, &c.CreatedAt, &c.UpdatedAt,
		); err != nil {
			log.Print(err)
		}
		cs = append(cs, c)
	}
	return
}
