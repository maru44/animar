package blog

import (
	"animar/v1/auth"
	"context"
	"fmt"
)

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

func ListBlogJoinAnimeUserDomain() []TBlogJoinAnimesUser {
	ctx := context.Background()
	rows := ListBlog()
	var blogs []TBlogJoinAnimesUser
	for rows.Next() {
		var blog TBlogJoinAnimesUser
		err := rows.Scan(
			&blog.ID, &blog.Slug, &blog.Title, &blog.Abstract,
			&blog.Content, &blog.UserId, &blog.CreatedAt, &blog.UpdatedAt,
		)

		// anim info(minimum)
		blogID := blog.ID
		blog.Animes = RelationBlogAnimesDomain(blogID)

		user := auth.GetUserFirebase(ctx, blog.UserId)
		blog.User = user

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

func DetailBlogJoinAnimeUserDomain(slug string) TBlogJoinAnimesUser {
	blog := DetailBlogWithUserBySlug(slug)
	blog.Animes = RelationBlogAnimesDomain(blog.ID)
	ctx := context.Background()
	user := auth.GetUserFirebase(ctx, blog.UserId)
	blog.User = user
	return blog
}

func RelationBlogAnimesDomain(blogId int) []TJoinedAnime {
	rows := RelationAnimeByBlog(blogId)
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
