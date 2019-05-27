package models

import (
	"time"
)

type Account struct {
	ID        uint64    `id`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`

	Username string `db:"username"`
	Password string `db:"password"`
	Email    string `db:"email"`
	// Slug           string `db:"slug" valid:"stringlength(1|255),required"`
	// Level          int    `db:"level"`
	// Banned         bool   `db:"banned"`
	Activated bool `db:"is_activated"`
	// FollowingCount int    `db:"following_count"` // calculate on fly? feature necessary?
}
