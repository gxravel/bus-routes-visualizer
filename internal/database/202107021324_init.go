package database

import (
	"database/sql"

	"github.com/lopezator/migrator"
	"github.com/pkg/errors"
)

//nolint // to bypass gosec sql concat warning
func migrationInit(schema string) *migrator.Migration {
	return &migrator.Migration{
		Name: "202107021324_init",
		Func: func(tx *sql.Tx) error {
			qs := []string{}

			for k, query := range qs {
				if _, err := tx.Exec(query); err != nil {
					return errors.Wrapf(err, "applying 202107021324_init migration #%d", k)
				}
			}
			return nil
		},

	}
}

/* ROLLBACK SQL

*/
