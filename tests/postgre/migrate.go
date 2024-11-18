package postgre

import (
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"go-hexagonal/config"
)

func GolangMigrateUp(conf *config.Config) error {

	if conf.Postgre == nil {
		return nil
	}

	m, err := migrate.New(
		"file://"+conf.MigrationDir,
		fmt.Sprintf(
			"postgresql://%s:%s@%s:%d/%s?sslmode=%s",
			conf.Postgre.Username,
			conf.Postgre.Password,
			conf.Postgre.Host,
			conf.Postgre.Port,
			conf.Postgre.DbName,
			conf.Postgre.SSLMode,
		),
	)
	if err != nil {
		return err
	}
	defer m.Close()

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}

	return nil
}

func GolangMigrateDrop(conf *config.Config) error {

	m, err := migrate.New(
		"file://"+conf.MigrationDir,
		fmt.Sprintf(
			"postgresql://%s:%s@%s:%d/%s?sslmode=%s",
			conf.Postgre.DbName,
			conf.Postgre.Password,
			conf.Postgre.Host,
			conf.Postgre.Port,
			conf.Postgre.DbName,
			conf.Postgre.SSLMode,
		),
	)
	if err != nil {
		return err
	}
	defer m.Close()

	if err := m.Drop(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}

	return nil
}
