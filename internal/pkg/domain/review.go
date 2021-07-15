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

type ReviewWithAnimeSlug struct {
	TReview
	AnimeSlug  string `json:"anime_slug,omitempty"`
	AnimeTitle string `json:"anime_title,omitempty"`
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
	AnimeId int    `json:"anime_id"`
	Content string `json:"content,omitempty"`
	Rating  int    `json:"rating,omitempty"` // text/plainのpostに対応
}

type TReviews []TReview

type TReviewJoinAnimes []TReviewJoinAnime

func (r TReview) GetId() int {
	return r.ID
}

type ReviewInteractor interface {
	GetAllReviewIds() ([]int, error)
	GetReviewById(int) (ReviewWithAnimeSlug, error)
	GetOnesReviewByAnime(int, string) (TReview, error)
	GetAnimeReviews(int, string) (TReviews, error)
	GetOnesReviews(string) (TReviewJoinAnimes, error)
	PostReviewContent(TReviewInput, string) (string, error)
	UpsertReviewContent(TReviewInput, string) (string, error)
	PostReviewRating(TReviewInput, string) (int, error)
	UpsertReviewRating(TReviewInput, string) (int, error)
	GetRatingAverage(int) (string, error)
}
