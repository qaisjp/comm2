package models

import "time"

type Resource struct {
	ID        uint64    `db:"id"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
	AuthorID  uint64    `db:"author_id"`

	Name        string `db:"name"`
	Title       string `db:"title"`
	Description string `db:"description"`
	// Rating      int    `db:"rating"`
	// Downloads   int    `db:"downloads"`
	// Type        int    `db:"type"` // todo: ResourceType
}

type ResourcePackage struct {
	ID        uint64    `db:"id"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`

	ResourceID  uint64 `db:"resource_id"` // relation
	AuthorID    uint64 `db:"author_id"`   // relation
	Version     string `db:"version"`
	Description string `db:"description"`

	Filename string `db:"filename"`

	Draft bool `db:"draft"`
}
