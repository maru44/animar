package blog

import (
	"animar/v1/auth"
	"animar/v1/tools"
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
			&blog.Content, &blog.IsPublic,
			&blog.UserId, &blog.CreatedAt, &blog.UpdatedAt,
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
			&blog.Content, &blog.IsPublic,
			&blog.UserId, &blog.CreatedAt, &blog.UpdatedAt,
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
			&blog.Content, &blog.IsPublic, &blog.UserId,
			&blog.CreatedAt, &blog.UpdatedAt,
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

func ListBlogByUserJoinAnimeUserDomain(uid string) []TBlogJoinAnimesUser {
	ctx := context.Background()
	rows := ListUsersBlog(uid)
	var blogs []TBlogJoinAnimesUser
	for rows.Next() {
		var blog TBlogJoinAnimesUser
		err := rows.Scan(
			&blog.ID, &blog.Slug, &blog.Title, &blog.Abstract,
			&blog.Content, &blog.IsPublic, &blog.UserId,
			&blog.CreatedAt, &blog.UpdatedAt,
		)

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

func DetailBlogJoinAnimeUserFromIdDomain(id int) TBlogJoinAnimesUser {
	blog := DetailBlog(id)
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

func UpdateRelationBlogAnimesDomain(animeIds []int, blogId int) {
	rows := RelationAnimeByBlog(blogId)
	var originAnimeIds []int
	for rows.Next() {
		var ani TJoinedAnime
		err := rows.Scan(
			&ani.AnimeId, &ani.Slug, &ani.Title,
		)
		if err != nil {
			fmt.Print(err)
		}
		originAnimeIds = append(originAnimeIds, ani.AnimeId)
	}
	defer rows.Close()

	for _, animeId := range animeIds {
		containIdx := tools.IsContainInt(originAnimeIds, animeId)
		if containIdx == -1 {
			// もともとのに無ければ新規作成
			InsertRelationAnimeBlog(animeId, blogId)
		} else {
			// もともとのにあれば何もしないでもともとのリストから削除
			originAnimeIds = tools.SliceRemove(originAnimeIds, containIdx)
		}
	}
	// originで余ってるものはdelete対象
	for _, originId := range originAnimeIds {
		DeleteRelationAnimeBlog(originId, blogId)
	}
}
