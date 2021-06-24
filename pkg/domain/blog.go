package domain

type TBlog struct {
	ID        int     `json:"id"`
	Slug      string  `json:"slug"`
	Title     string  `json:"title"`
	Abstract  *string `json:"abstract"`
	Content   string  `json:"content"`
	IsPublic  bool    `json:"is_public"`
	UserId    string  `json:"user_id"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt string  `json:"updated_at"`
}

type TJoinedAnime struct {
	AnimeId int    `json:"anime_id"`
	Slug    string `json:"slug"`
	Title   string `json:"title"`
}

type TJoinedBlog struct {
	BlogId    int    `json:"blog_id"`
	Title     string `json:"title"`
	Slug      string `json:"slug"`
	Abstract  string `json:"abstract"`
	CreatedAt string `json:"created_at"`
}

type TBlogJoinAnime struct {
	ID        int            `json:"id"`
	Slug      string         `json:"slug"`
	Title     string         `json:"title"`
	Abstract  *string        `json:"abstract"`
	Content   string         `json:"content"`
	UserId    string         `json:"user_id"`
	IsPublic  bool           `json:"is_public"`
	CreatedAt string         `json:"created_at"`
	UpdatedAt string         `json:"updated_at"`
	Animes    []TJoinedAnime `json:"animes"`
}

type TBlogInsert struct {
	Title    string `json:"title"`
	Abstract string `json:"abstract,omitempty"`
	Content  string `json:"content,omitempty"`
	IsPublic bool   `json:"is_public"`
	AnimeIds []int  `json:"anime_ids"`
}

type TBlogs []TBlog

type TJoinedBlogs []TJoinedBlog

type TBlogJoinAnimes []TBlogJoinAnime

type BlogInteractor interface {
	ListBlog() (TBlogJoinAnimes, error)
	ListBlogByUser(string, string) (TBlogJoinAnimes, error)
	BlogUserId(int) (string, error)
	DetailBlog(int) (TBlogJoinAnime, error)
	DetailBlogBySlug(string) (TBlogJoinAnime, error)
	InsertBlog(TBlogInsert, string) (int, error)
	UpdateBlog(TBlogInsert, int) (int, error)
	DeleteBlog(int) (int, error)
	// relation
	RelationAnimeByBlog(int) ([]TJoinedAnime, error)
	// RelationBlogByAnime
	InsertRelationAnime(int, int) (bool, error)
	DeleteRelationAnime(int, int) error
}

func (b TBlog) GetId() int {
	return b.ID
}

func (b TBlogJoinAnime) GetId() int {
	return b.ID
}
