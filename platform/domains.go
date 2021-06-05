package platform

func ListPlatformDomain() []TPlatform {
	rows := ListPlatform()
	var plats []TPlatform
	for rows.Next() {
		var plat TPlatform
		err := rows.Scan(
			&plat.ID, &plat.EngName, &plat.PlatName,
			&plat.BaseUrl, &plat.Image, &plat.IsValid,
			&plat.CreatedAt, &plat.UpdatedAt,
		)
		if err != nil {
			panic(err.Error())
		}
		plats = append(plats, plat)
	}

	defer rows.Close()
	return plats
}
