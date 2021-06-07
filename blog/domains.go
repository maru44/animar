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
		var b TBlog
		err := rows.Scan(
			&b.ID, &b.Slug, &b.Title, &b.Abstract,
			&b.Content, &b.UserId, &b.IsPublic,
			&b.CreatedAt, &b.UpdatedAt,
		)
		if err != nil {
			panic(err.Error())
		}
		blogs = append(blogs, b)
	}
	defer rows.Close()
	return blogs
}

// List (blog + anime)
func ListBlogJoinAnimeDomain() []TBlogJoinAnimes {
	rows := ListBlog()
	var blogs []TBlogJoinAnimes
	for rows.Next() {
		var b TBlogJoinAnimes
		err := rows.Scan(
			&b.ID, &b.Slug, &b.Title, &b.Abstract,
			&b.Content, &b.UserId, &b.IsPublic,
			&b.CreatedAt, &b.UpdatedAt,
		)

		// anim info(minimum)
		blogID := b.ID
		b.Animes = RelationBlogAnimesDomain(blogID)

		if err != nil {
			fmt.Print(err)
		}
		blogs = append(blogs, b)
	}
	defer rows.Close()
	return blogs
}

// List (blog + anime + user)
// @NOTUSED
func ListBlogJoinAnimeUserDomain() []TBlogJoinAnimesUser {
	ctx := context.Background()
	rows := ListBlog()
	var blogs []TBlogJoinAnimesUser
	for rows.Next() {
		var b TBlogJoinAnimesUser
		err := rows.Scan(
			&b.ID, &b.Slug, &b.Title, &b.Abstract,
			&b.Content, &b.UserId, &b.IsPublic,
			&b.CreatedAt, &b.UpdatedAt,
		)

		// anim info(minimum)
		blogID := b.ID
		b.Animes = RelationBlogAnimesDomain(blogID)

		user := auth.GetUserFirebase(ctx, b.UserId)
		b.User = user

		if err != nil {
			fmt.Print(err)
		}
		blogs = append(blogs, b)
	}
	defer rows.Close()
	return blogs
}

// List of users (blog + anime)
func ListBlogByUserJoinAnimeDomain(uid string) []TBlogJoinAnimes {
	rows := ListUsersBlog(uid)
	var blogs []TBlogJoinAnimes
	for rows.Next() {
		var b TBlogJoinAnimes
		err := rows.Scan(
			&b.ID, &b.Slug, &b.Title, &b.Abstract,
			&b.Content, &b.IsPublic, &b.UserId,
			&b.CreatedAt, &b.UpdatedAt,
		)
		blogID := b.ID
		b.Animes = RelationBlogAnimesDomain(blogID)

		if err != nil {
			fmt.Print(err)
		}
		blogs = append(blogs, b)
	}
	defer rows.Close()
	return blogs
}

// List of users (blog + anime + user)
func ListBlogByUserJoinAnimeUserDomain(uid string) []TBlogJoinAnimesUser {
	ctx := context.Background()
	rows := ListUsersBlog(uid)
	var blogs []TBlogJoinAnimesUser
	for rows.Next() {
		var b TBlogJoinAnimesUser
		err := rows.Scan(
			&b.ID, &b.Slug, &b.Title, &b.Abstract,
			&b.Content, &b.UserId, &b.IsPublic,
			&b.CreatedAt, &b.UpdatedAt,
		)

		blogID := b.ID
		b.Animes = RelationBlogAnimesDomain(blogID)

		user := auth.GetUserFirebase(ctx, b.UserId)
		b.User = user

		if err != nil {
			fmt.Print(err)
		}
		blogs = append(blogs, b)
	}
	defer rows.Close()
	return blogs
}

func DetailBlogJoinAnimeDomain(slug string) TBlogJoinAnimes {
	blog := DetailBlogBySlug(slug)
	blog.Animes = RelationBlogAnimesDomain(blog.ID)
	return blog
}

func DetailBlogJoinAnimeFromIdDomain(id int) TBlogJoinAnimes {
	blog := DetailBlog(id)
	blog.Animes = RelationBlogAnimesDomain(id)
	return blog
}

// @NOTUSED
func DetailBlogJoinAnimeUserDomain(slug string) TBlogJoinAnimesUser {
	blog := DetailBlogWithUserBySlug(slug)
	blog.Animes = RelationBlogAnimesDomain(blog.ID)
	ctx := context.Background()
	user := auth.GetUserFirebase(ctx, blog.UserId)
	blog.User = user
	return blog
}

// @NOTUSED
func DetailBlogJoinAnimeUserFromIdDomain(id int) TBlogJoinAnimesUser {
	blog := DetailBlogWithUser(id)
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
		var a TJoinedAnime
		err := rows.Scan(
			&a.AnimeId, &a.Slug, &a.Title,
		)
		if err != nil {
			fmt.Print(err)
		}
		ret = append(ret, a)
	}
	defer rows.Close()
	return ret
}

func UpdateRelationBlogAnimesDomain(animeIds []int, blogId int) {
	rows := RelationAnimeByBlog(blogId)
	var originAnimeIds []int
	for rows.Next() {
		var a TJoinedAnime
		err := rows.Scan(
			&a.AnimeId, &a.Slug, &a.Title,
		)
		if err != nil {
			fmt.Print(err)
		}
		originAnimeIds = append(originAnimeIds, a.AnimeId)
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
