package series

func ListSeriesDomain() []TSeries {
	rows := ListSeries()
	var sers []TSeries
	for rows.Next() {
		var ser TSeries
		err := rows.Scan(
			&ser.ID, &ser.EngName, &ser.SeriesName,
			&ser.CreatedAt, &ser.UpdatedAt,
		)
		if err != nil {
			panic(err.Error())
		}
		sers = append(sers, ser)
	}

	defer rows.Close()
	return sers
}
