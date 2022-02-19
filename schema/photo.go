package schema

import "time"

// swagger:model Photo
type Photo struct {
	ID        int       `json:"id,omitempty" sql:"primary"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
	Size      int64     `json:"size,omitempty"`
	FileName  string    `json:"file_name,omitempty"`
	FileHash  string    `json:"file_hash,omitempty"`
	UserID    int       `json:"user_id,omitempty"`
	Extension string    `json:"extension,omitempty"`
}