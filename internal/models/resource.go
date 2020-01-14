package models

import "time"

type Resource struct {
	ID        uint64    `db:"id" json:"id"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
	AuthorID  uint64    `db:"author_id" json:"author_id"`

	Name        string `db:"name" json:"name"`
	Title       string `db:"title" json:"title"`
	Description string `db:"description" json:"description"`
	// Rating      int    `db:"rating"`
	// Downloads   int    `db:"downloads"`
	// Type        int    `db:"type"` // todo: ResourceType
}

type ResourcePackage struct {
	ID        uint64    `db:"id" json:"id"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`

	ResourceID  uint64 `db:"resource_id" json:"resource_id"` // relation
	AuthorID    uint64 `db:"author_id" json:"author_id"`     // relation
	Version     string `db:"version" json:"version"`
	Description string `db:"description" json:"description"`

	Filename string `db:"filename" json:"filename"`

	Draft bool `db:"draft" json:"draft"`
}
