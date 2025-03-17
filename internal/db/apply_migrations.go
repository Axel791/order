package db

import (
	"context"
	"github.com/pressly/goose/v3"

	"github.com/Axel791/order/internal/config"
	"github.com/jmoiron/sqlx"
)

// AppleMigration - Применение миграций
func AppleMigration(dbConn *sqlx.DB, cfg *config.Config) error {
	return goose.RunContext(context.Background(), "up", dbConn.DB, cfg.MigrationsPath)
}
