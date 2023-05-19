package types

type Episode struct {
	Id            string `json:"id"`
	Title         string `json:"title"`
	Description   string `json:"description"`
	PodcastId     string `json:"podcast_id"`
	CoverImageUrl string `json:"cover_image_url"`
    AudioUrl      string `json:"audio_url"`
    DurationMs    int    `json:"duration_ms"`
	CreatedAt     string `json:"created_at"`
	UpdatedAt     string `json:"updated_at"`
}
