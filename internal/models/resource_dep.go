package models

type ResourceDep struct {
	PackageID uint64 `db:"package_id"` // package
	DependsOn uint64 `db:"depends_on"` // package
}
