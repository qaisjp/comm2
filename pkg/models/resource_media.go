package models

import "time"

type ResourceMedia struct {
	ID        uint64    `db:"id"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`

	Resource    uint64 `db:"resource"` // relation
	Title       string `db:"title"`
	Description string `db:"description"`
	Author      uint64 `db:"author"` // user - the person who uploaded the media

	Type int `db:"type"` // image or video?

	// Do we need some sort of abstract file storage here?
	// todo
	Filename string `db:"filename"`
	FileURL  string `db:"url"`
}
