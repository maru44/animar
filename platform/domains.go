package platform

func ListPlatformDomain() []TPlatform {
	rows := listPlatform()
	var plats []TPlatform
	for rows.Next() {
		var p TPlatform
		err := rows.Scan(
			&p.ID, &p.EngName, &p.PlatName,
			&p.BaseUrl, &p.Image, &p.IsValid,
			&p.CreatedAt, &p.UpdatedAt,
		)
		if err != nil {
			panic(err.Error())
		}
		plats = append(plats, p)
	}

	defer rows.Close()
	return plats
}

/****************************
*          relation		    *
****************************/

func ListRelationPlatformDomain(animeId int) []TRelationPlatform {
	rows := relationPlatformByAnime(animeId)
	var relations []TRelationPlatform
	for rows.Next() {
		var r TRelationPlatform
		err := rows.Scan(
			&r.PlatformId, &r.AnimeId, &r.LinkUrl,
			&r.CreatedAt, &r.UpdatedAt, &r.PlatName,
		)
		if err != nil {
			panic(err.Error())
		}
		relations = append(relations, r)
	}
	defer rows.Close()
	return relations
}
