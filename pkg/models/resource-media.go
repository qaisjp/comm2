package models

import "time"

type ResourceMedia struct {
	ID           uint64
	DateCreated  time.Time
	DateModified time.Time

	Resource    uint64 // relation
	Title       string
	Description string
	File        string // todo
	Author      uint64 // user
}
