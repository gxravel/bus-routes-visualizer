package database

import (
	"database/sql"

	"github.com/lopezator/migrator"
	"github.com/pkg/errors"
)

//nolint // to bypass gosec sql concat warning
func migrationInit(schema string) *migrator.Migration {
	return &migrator.Migration{
		Name: "202107051748_init",
		Func: func(tx *sql.Tx) error {
			qs := []string{
				`CREATE TABLE IF NOT EXISTS route (
					id BIGINT AUTO_INCREMENT PRIMARY KEY,
					city VARCHAR(255) NOT NULL,
					bus VARCHAR(255) NOT NULL,
					UNIQUE (city, bus)
				)`,

				`CREATE TABLE IF NOT EXISTS route_point (
					id BIGINT AUTO_INCREMENT PRIMARY KEY,
					step TINYINT NOT NULL,
					address VARCHAR(255) NOT NULL,
					route_id BIGINT NOT NULL,
					FOREIGN KEY(route_id) REFERENCES route(id) ON UPDATE CASCADE ON DELETE CASCADE
				)`,
			}

			for k, query := range qs {
				if _, err := tx.Exec(query); err != nil {
					return errors.Wrapf(err, "applying 202107051748_init migration #%d", k)
				}
			}
			return nil
		},
	}
}

/* ROLLBACK SQL
DROP TABLE IF EXISTS route
DROP TABLE IF EXISTS route_point
*/
