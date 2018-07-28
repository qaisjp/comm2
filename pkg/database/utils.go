package database

import (
	"github.com/jmoiron/sqlx"
)

// RunInTransaction runs a function in a transaction. If function
// returns an error transaction is rollbacked, otherwise transaction
// is committed.
func RunInTransaction(db *sqlx.DB, fn func(*sqlx.Tx) error) (err error) {
	tx, err := db.Beginx()
	if err != nil {
		return err
	}

	defer func() {
		if localErr := recover(); localErr != nil {
			err = localErr.(error)
			tx.Rollback()
		}
	}()

	if err := fn(tx); err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit()

	return err
}
