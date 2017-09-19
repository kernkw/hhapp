package schema

type Frequency int

const (
	Daily Frequency = iota
	Weekly
	Monthly
)

type Notification struct {
	ID          int       `json:"id"`
	ListID      int       `json:"list_id, omitempty"`
	FavoritesID int       `json:"favorites_id, omitempty"`
	Name        string    `json:"name"`
	Frequency   Frequency `json:"frequency"`
}
