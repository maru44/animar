package database

import (
	"animar/v1/pkg/domain"
	"animar/v1/pkg/tools/tools"
)

type BlogRepository struct {
	SqlHandler
}

func (repo *BlogRepository) ListBlog() (blogs domain.TBlogs, err error) {
	rows, err := repo.Query(
		"SELECT * FROM blogs WHERE is_public = true",
	)
	defer rows.Close()

	if err != nil {
		tools.ErrorLog(err)
		return
	}
	for rows.Next() {
		var b domain.TBlog
		err := rows.Scan(
			&b.ID, &b.Slug, &b.Title, &b.Abstract, &b.Content,
			&b.IsPublic, &b.UserId, &b.CreatedAt, &b.UpdatedAt,
		)
		if err != nil {
			tools.ErrorLog(err)
		}
		blogs = append(blogs, b)
	}
	return
}

func (repo *BlogRepository) ListBlogByUser(accessUserId string, blogUserId string) {}
