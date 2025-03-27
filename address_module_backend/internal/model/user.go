package model

import "time"

type User struct {
	ID             int64      `json:"id"`
	Username       string     `json:"username"`
	Email          string     `json:"email"`
	HashedPassword string     `json:"hashed_password"`
	CreatedAt      time.Time  `json:"created_at,omitempty"`
	CreatedBy      *int64     `json:"created_by,omitempty"`
	LastLogin      *time.Time `json:"last_login,omitempty"`
}
