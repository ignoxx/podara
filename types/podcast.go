package types

type Podcast struct {
	Id            string `json:"id"`
	Title         string `json:"title"`
	Description   string `json:"description"`
	UserId        string `json:"user_id"`
	CoverImageUrl string `json:"cover_image_url"`
	CreatedAt     string `json:"created_at,omitempty"`
	UpdatedAt     string `json:"updated_at,omitempty"`
}
