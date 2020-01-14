package models

type ResourceRating struct {
	UserID     uint64 `db:"user_id" json:"user_id"`
	ResourceID uint64 `db:"resource_id" json:"resource_id"`
	Positive   bool   `db:"positive" json:"positive"`
}
