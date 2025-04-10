package utils

import (
	"context"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/Slicky-Team/slickfame/sqlengine"
	"github.com/jackc/pgx/v4"
)

func getAppliedMigrations(ctx context.Context, dbEngine sqlengine.DBEngine) (map[string]struct{}, error) {
	rows, err := dbEngine.QueryContext(ctx, "SELECT migration_name FROM migrations")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	appliedMigrations := make(map[string]struct{})
	for rows.Next() {
		var migrationName string
		if err := rows.Scan(&migrationName); err != nil {
			return nil, err
		}
		appliedMigrations[migrationName] = struct{}{}
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return appliedMigrations, nil
}

func RunMigrations(ctx context.Context, dbEngine sqlengine.DBEngine) error {
	slog.Info("Starting to run migrations...")

	migrations, err := getMigrationsFromFiles("migrations")
	if err != nil {
		return err
	}

	appliedMigrations, err := getAppliedMigrations(ctx, dbEngine)
	if err != nil {
		_, err := dbEngine.ExecuteCmd(ctx, "CREATE TABLE IF NOT EXISTS migrations (migration_name TEXT PRIMARY KEY, created_at TIMESTAMPTZ DEFAULT NOW())")
		if err != nil {
			return err
		}
		appliedMigrations = make(map[string]struct{})
	}

	newMigrationsCount := 0
	err = dbEngine.EnableTx(ctx, func(tx pgx.Tx) error {
		// Perform your database operations here
		for _, migration := range migrations {
			if _, exists := appliedMigrations[migration]; !exists {
				content, err := os.ReadFile(filepath.Join("migrations", migration))
				if err != nil {
					return err
				}
				_, err = tx.Exec(ctx, string(content))
				if err != nil {
					return err
				}
				_, err = tx.Exec(ctx, "INSERT INTO migrations (migration_name) VALUES ($1)", migration)
				if err != nil {
					return err
				}
				newMigrationsCount++
			}
		}

		return nil
	})

	if err != nil {
		return err
	}

	slog.Info("Applied new migration(s): ", slog.Int("count", newMigrationsCount))
	return nil
}

func getMigrationsFromFiles(dir string) ([]string, error) {
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var migrations []string
	for _, file := range files {
		if !file.IsDir() {
			migrations = append(migrations, file.Name())
		}
	}

	return migrations, nil
}
