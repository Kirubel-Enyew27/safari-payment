package intiator

import (
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"go.uber.org/zap"
)

func InitiateMigration(path, conn string, log *zap.Logger) *migrate.Migrate {
	m, err := migrate.New(fmt.Sprintf("file://%s", path), conn)
	if err != nil {
		log.Fatal("could not create migrator", zap.Error(err))
	}
	return m
}

func UpMigration(m *migrate.Migrate, log *zap.Logger) {
	err := m.Up()
	if err != nil && err != migrate.ErrNoChange {
		log.Fatal("could not migrate", zap.Error(err))
	}
}
