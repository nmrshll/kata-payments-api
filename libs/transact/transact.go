package transact

import (
	"context"
	"database/sql"

	"github.com/volatiletech/sqlboiler/boil"
)

func Transact(txFunc func(*sql.Tx) error) (err error) {
	tx, err := boil.BeginTx(context.Background(), nil)
	if err != nil {
		return err
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p) // re-throw panic after Rollback
		} else if err != nil {
			tx.Rollback() // err is non-nil; don't change it
		} else {
			err = tx.Commit() // err is nil; if Commit returns error update err
		}
	}()

	err = txFunc(tx)
	return err
}
