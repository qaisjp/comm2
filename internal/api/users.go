package api

import (
	"crypto/md5"
	"fmt"
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

func (u User) PublicInfo() PublicUserInfo {
	return PublicUserInfo{
		ID:        u.ID,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,

		Username: u.Username,
		Gravatar: fmt.Sprintf("https://www.gravatar.com/avatar/%x", md5.Sum([]byte(u.Email))),
	}
}

// User represents a public user object
type PublicUserInfo struct {
	ID        uint64    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Username string `json:"username"`
	Gravatar string `json:"gravatar"`
}
