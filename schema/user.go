package schema

import "time"

// swagger:model User
type User struct {
	ID             int       `json:"id,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      *time.Time `json:"updated_at"`
	Email          string    `json:"username,omitempty"`
	Password       string    `json:"password,omitempty"`
	StoragePath    string    `json:"storage_path,omitempty"`
}