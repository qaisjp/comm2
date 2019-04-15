package models

type ResourceRating struct {
	Account  uint64 `db:"account"`
	Resource uint64 `db:"resource"`
	Positive bool   `db:"positive"`
}
