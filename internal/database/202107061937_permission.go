package database

import (
	"database/sql"

	"github.com/lopezator/migrator"
	"github.com/pkg/errors"
)

//nolint // to bypass gosec sql concat warning
func migrationPermission(schema string) *migrator.Migration {
	return &migrator.Migration{
		Name: "202107061916_permission",
		Func: func(tx *sql.Tx) error {
			qs := []string{
				`CREATE TABLE IF NOT EXISTS permission(
					user_id BIGINT PRIMARY KEY NOT NULL,
					actions JSON DEFAULT NULL
				)`,
			}

			for k, query := range qs {
				if _, err := tx.Exec(query); err != nil {
					return errors.Wrapf(err, "applying 202107061916_permission migration #%d", k)
				}
			}
			return nil
		},
	}
}

/* ROLLBACK SQL
DROP TABLE IF EXISTS permission
*/
