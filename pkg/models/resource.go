package models

import "time"

type Resource struct {
	ID        uint64    `db:"id`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`

	Name        string `db:"name"`
	Title       string `db:"title"`
	Description string `db:"description"`
	Rating      int    `db:"rating"`
	Downloads   int    `db:"downloads"`
	Type        int    `db:"type"` // todo: ResourceType
}
