package models

type ResourceDep struct {
	Package   uint64 `db:"package"`    // package
	DependsOn uint64 `db:"depends_on"` // package
}
