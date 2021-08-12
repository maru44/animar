package domain

type Article struct {
	ID          int     `json:"id"`
	Slug        string  `json:"slug"`
	ArticleType string  `json:"article_type"`
	Abstract    *string `json:"abstract,omitempty"`
	Content     *string `json:"content,omitempty"`
	Image       *string `json:"image,omitempty"`
	Author      *string `json:"author,omitempty"`
	IsPublic    bool    `json:"is_public"`
	UserId      string  `json:"user_id"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
}

type ArticleCharacter struct {
	ID        int    `json:"id"`
	ArticleId *int   `json:"article_id,omitempty"`
	Name      string `json:"chara_name,omitempty"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type InterviewQuote struct {
	ID        int    `json:"id"`
	ArticleId *int   `json:"article_id,omitempty"`
	CharaId   *int   `json:"chara_id,omitempty"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
