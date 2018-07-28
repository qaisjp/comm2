package models

import "time"

type ResourcePackage struct {
	ID        uint64    `db:"id"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`

	Resource      uint64 `db:"resource"` // relation
	Version       string `db:"version"`
	DownloadCount int    `db:"download_count"`

	// Do we need some sort of abstract file storage here?
	// todo
	Filename string `db:"filename"`
	FileURL  string `db:"url"`
}
