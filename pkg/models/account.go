package models

import (
	"time"
)

type Account struct {
	ID           uint64
	DateCreated  time.Time
	DateModified time.Time

	Email          string
	Username       string
	Slug           string
	Password       string
	Level          int
	Banned         bool
	Activated      bool
	FollowingCount int // calculate on fly? feature necessary?
}
