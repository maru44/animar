package domain

type TReview struct {
	ID        int     `json:"id"`
	Content   *string `json:"content,omitempty"`
	Rating    *int    `json:"rating,omitempty"`
	AnimeId   int     `json:"anime_id"`
	UserId    *string `json:"user_id"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt string  `json:"updated_at,omitempty"`
}

type TReviewJoinAnime struct {
	ID           int     `json:"id"`
	Content      *string `json:"content,omitempty"`
	Rating       *int    `json:"rating,omitempty"`
	AnimeId      int     `json:"anime_id"`
	UserId       *string `json:"user_id"`
	CreatedAt    string  `json:"created_at"`
	UpdatedAt    string  `json:"updated_at,omitempty"`
	Title        string  `json:"title"`
	Slug         string  `json:"slug"`
	AnimeContent *string `json:"anime_content,omitempty"`
	AState       *string `json:"anime_state,omitempty"`
}

type TReviewInput struct {
	AnimeId int     `json:"anime_id"`
	Content *string `json:"content,omitempty"`
	Rating  *int    `json:"rating,string,omitempty"` // text/plainのpostに対応
	UserId  string  `json:"user_id"`
}

type TReviews []TReview

type TReviewJoinAnimes []TReviewJoinAnime

func (r TReview) GetId() int {
	return r.ID
}
