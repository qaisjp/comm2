package api

import (
	"context"
	"crypto/md5"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
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

type AuthenticatedUser struct {
	PublicUserInfo
	Level int `json:"level"`
}

func (u User) PrivateInfo() AuthenticatedUser {
	return AuthenticatedUser{
		PublicUserInfo: u.PublicInfo(),
		Level:          u.Level,
	}
}

func (u User) GetFollowers(ctx context.Context, db *sqlx.DB) (rows []User, err error) {
	err = db.SelectContext(ctx, &rows, "select u.* from users u, user_followings f where f.target_user_id=$1 and f.source_user_id=u.id", u.ID)
	return
}

func (u User) GetFollowing(ctx context.Context, db *sqlx.DB) (rows []User, err error) {
	err = db.SelectContext(ctx, &rows, "select u.* from users u, user_followings f where f.source_user_id=$1 and f.target_user_id=u.id", u.ID)
	return
}

// User represents a public user object
type PublicUserInfo struct {
	ID        uint64    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Username string `json:"username"`
	Gravatar string `json:"gravatar"`

	FollowsYou *bool `json:"follows_you,omitempty"`
}

type UserSlice []User

func (s UserSlice) PublicInfo() []PublicUserInfo {
	output := []PublicUserInfo{}
	for _, u := range s {
		output = append(output, u.PublicInfo())
	}
	return output
}
