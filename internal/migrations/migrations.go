package migrations

import (
	"database/sql"
	"log/slog"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func RunMigrations(db *sql.DB) error {
    driver, err := postgres.WithInstance(db, &postgres.Config{})
    if err != nil {
        return err
    }

    m, err := migrate.NewWithDatabaseInstance(
        "file://./migrations",
        "postgres", 
		driver,
	)
    if err != nil {
        return err
    }

    err = m.Up()
    if err != nil && err != migrate.ErrNoChange {
        return err
    }

	slog.Info("Migrations applied successfully")
    return nil
}