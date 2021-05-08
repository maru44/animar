package anime

import "database/sql"

func ListAnimeDomain() []TAnime {
	rows := ListAnime()
	var animes []TAnime
	for rows.Next() {
		var ani TAnime
		nullContent := new(sql.NullString)
		err := rows.Scan(&ani.ID, &ani.Title, nullContent, &ani.CreatedAt, &ani.UpdatedAt)
		if err != nil {
			panic(err.Error())
		}
		ani.Content = nullContent.String
		animes = append(animes, ani)
	}

	defer rows.Close()

	return animes
}
