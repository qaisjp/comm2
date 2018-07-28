package models

import "time"

type Resource struct {
	ID           uint64
	DateCreated  time.Time
	DateModified time.Time

	Name        string
	LongName    string
	Description string
	Rating      int
	Downloads   int
	Type        int // todo: ResourceType
}
