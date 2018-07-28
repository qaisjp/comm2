package models

import "time"

type ResourcePackage struct {
	ID           uint64
	DateCreated  time.Time
	DateModified time.Time

	Resource  uint64 // relation
	Filename  string // download filename(?)
	URL       string // url to download at
	Version   string
	Downloads int
}
