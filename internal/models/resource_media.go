package models

import "time"

type ResourceMedia struct {
	ID        uint64    `db:"id"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`

	ResourceID  uint64 `db:"resource_id"` // relation
	Title       string `db:"title"`
	Description string `db:"description"`
	AuthorID    uint64 `db:"author_id"` // user - the person who uploaded the media

	Type int `db:"type"` // image or video?

	// Do we need some sort of abstract file storage here?
	// todo
	Filename string `db:"filename"`
	FileURL  string `db:"url"`
}
