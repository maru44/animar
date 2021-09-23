package domain

type ArticleType string

const (
	ArticleTypeInterview = "interview"
	ArticleTypeArticle   = "article"
)

type Article struct {
	ID          int         `json:"id"`
	Slug        string      `json:"slug"`
	ArticleType ArticleType `json:"article_type"`
	Abstract    *string     `json:"abstract,omitempty"`
	Content     *string     `json:"content,omitempty"`
	Image       *string     `json:"image,omitempty"`
	Author      *string     `json:"author,omitempty"`
	IsPublic    bool        `json:"is_public"`
	UserId      string      `json:"user_id"`
	CreatedAt   string      `json:"created_at"`
	UpdatedAt   string      `json:"updated_at"`
	// Characters  []ArticleCharacter `json:"characters"`
}

type ArticleCharacter struct {
	ID        int     `json:"id"`
	Name      string  `json:"chara_name,omitempty"`
	Image     *string `json:"chara_image,omitempty"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt string  `json:"updated_at"`
}

type InterviewQuote struct {
	ID        int    `json:"id"`
	ArticleId *int   `json:"article_id,omitempty"`
	CharaId   *int   `json:"chara_id,omitempty"`
	Content   string `json:"content"`
	UserId    string `json:"user_id"`
	Sequence  int    `json:"sequence"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type ArticleInput struct {
	ArticleType string  `json:"article_type"`
	Abstract    *string `json:"abstract,omitempty"`
	Content     *string `json:"content,omitempty"`
	Image       *string `json:"image,omitempty"`
	Author      *string `json:"author,omitempty"`
	IsPublic    bool    `json:"is_public"`
	// UserId      string  `json:"user_id"`
}

type ArticleCharacterInput struct {
	Name  string  `json:"chara_name"`
	Image *string `json:"image,omitempty"`
}

type InterviewQuoteInput struct {
	ArticleId *int   `json:"article_id,omitempty"`
	CharaId   *int   `json:"chara_id,omitempty"`
	Content   string `json:"content"`
	Sequence  int    `json:"sequence"`
}

type RelationArticleCharacterInput struct {
	ArticleId int `json:"article_id"`
	CharaId   int `json:"chara_id"`
}

type RelationArticleAnimeInput struct {
	AnimeId   int `json:"anime_id"`
	ArticleId int `json:"article_id"`
}

type ArticleInteractor interface {
	FetchArticles() ([]Article, error)
	GetArticleById(id int) (Article, error)
	GetArticleBySlug(slug string) (Article, error)
	InsertArticle(articleInput ArticleInput, userId string) (int, error)
	UpdateArticle(articleInput ArticleInput, id int) (int, error)
	DeleteArticle(id int) (int, error)
	FetchArticleByAnime(animeId int) ([]Article, error)
	// character
	FetchArticleCharas(articleId int) ([]ArticleCharacter, error)
	FetchArticleCharasByUser(userId string) ([]ArticleCharacter, error)
	InsertArticleChara(charaInput ArticleCharacterInput, userId string) (int, error)
	UpdateArticleChara(charaInput ArticleCharacterInput, id int) (int, error)
	DeleteArticleChara(id int) (int, error)
	// interview
	FetchInterviewQuotes(articleId int) ([]InterviewQuote, error)
	InsertInterviewQuote(interviewInput InterviewQuoteInput, userId string) (int, error)
	UpdateInterviewQuote(interviewInput InterviewQuoteInput, id int) (int, error)
	DeleteInterviewQuote(id int) (int, error)
	// input relation
	InsertRelationArticleCharacter(in RelationArticleCharacterInput) (int, error)
	DeleteRelationArticleCharacter(in RelationArticleCharacterInput) (int, error)
	InsertRelationArticleAnime(in RelationArticleAnimeInput) (int, error)
	DeleteRelationArticleAnime(in RelationArticleAnimeInput) (int, error)
}
