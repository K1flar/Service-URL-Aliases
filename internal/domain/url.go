package domain

type URL struct {
	ID     uint32 `json:"id,omitempty"`
	URL    string `json:"url"`
	Alias  string `json:"alias,omitempty"`
	UserID uint32 `json:"user_id,omitempty"`
}
