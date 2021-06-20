package series

import "animar/v1/pkg/tools/tools"

func ListSeriesDomain() []TSeries {
	rows := ListSeries()
	var sers []TSeries
	for rows.Next() {
		var s TSeries
		err := rows.Scan(
			&s.ID, &s.EngName, &s.SeriesName,
			&s.CreatedAt, &s.UpdatedAt,
		)
		if err != nil {
			tools.ErrorLog(err)
		}
		sers = append(sers, s)
	}

	defer rows.Close()
	return sers
}
