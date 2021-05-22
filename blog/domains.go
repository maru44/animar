package blog

func ListBlogDomain() []TBlog {
	rows := ListBlog()
	var blogs []TBlog
	for rows.Next() {
		var blog TBlog
		err := rows.Scan(
			&blog.ID, &blog.Slug, &blog.Title, &blog.Abstract,
			&blog.Content, &blog.UserId, &blog.CreatedAt, &blog.UpdatedAt,
		)
		if err != nil {
			panic(err.Error())
		}
		blogs = append(blogs, blog)
	}
	defer rows.Close()
	return blogs
}
