package models

import (
	"time"
)

// User represents a user account
type User struct {
	ID        uint64    `db:"id"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`

	Username string `db:"username"`
	Password string `db:"password"`
	Email    string `db:"email"`
	// Slug           string `db:"slug" valid:"stringlength(1|255),required"`
	Level     int  `db:"level"`
	Activated bool `db:"is_activated"`
	Banned    bool `db:"is_banned"`
	// FollowingCount int    `db:"following_count"` // calculate on fly? feature necessary?
}
