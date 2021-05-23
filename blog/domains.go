package blog

import "fmt"

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

func ListBlogJoinAnimeDomain() []TBlogJoinAnimes {
	rows := ListBlog()
	var blogs []TBlogJoinAnimes
	for rows.Next() {
		var blog TBlogJoinAnimes
		err := rows.Scan(
			&blog.ID, &blog.Slug, &blog.Title, &blog.Abstract,
			&blog.Content, &blog.UserId, &blog.CreatedAt, &blog.UpdatedAt,
		)

		// anim info(minimum)
		blogID := blog.ID
		blog.Animes = RelationBlogAnimesDomain(blogID)

		if err != nil {
			fmt.Print(err)
		}
		blogs = append(blogs, blog)
	}
	defer rows.Close()
	return blogs
}

func DetailBlogJoinAnimeDomain(slug string) TBlogJoinAnimes {
	blog := DetailBlogBySlug(slug)
	fmt.Print(blog)
	blog.Animes = RelationBlogAnimesDomain(blog.ID)
	return blog
}

func RelationBlogAnimesDomain(blogId int) []TJoinedAnime {
	rows := RelationBlogAnimes(blogId)
	var ret []TJoinedAnime

	for rows.Next() {
		var ani TJoinedAnime
		err := rows.Scan(
			&ani.AnimeId, &ani.Slug, &ani.Title,
		)
		if err != nil {
			fmt.Print(err)
		}
		ret = append(ret, ani)
	}
	defer rows.Close()
	return ret
}
