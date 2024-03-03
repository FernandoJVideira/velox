package velox

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/gobuffalo/pop"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func (v *Velox) PopConnect() (*pop.Connection, error) {
	tx, err := pop.Connect("development")
	if err != nil {
		return nil, err
	}
	return tx, nil
}

func (v *Velox) CreatePopMigration(up, down []byte, migrationName, migrationType string) error {
	var migrationPath = v.RootPath + "/migrations"
	err := pop.MigrationCreate(migrationPath, migrationName, migrationType, up, down)
	if err != nil {
		return err
	}
	return nil
}

func (v *Velox) RunPopMigrations(tx *pop.Connection) error {
	var migrationPath = v.RootPath + "/migrations"

	fm, err := pop.NewFileMigrator(migrationPath, tx)
	if err != nil {
		return err
	}

	err = fm.Up()
	if err != nil {
		return err
	}
	return nil
}

func (v *Velox) PopMigrateDown(tx *pop.Connection, steps ...int) error {
	var migrationPath = v.RootPath + "/migrations"

	step := 1
	if len(steps) > 0 {
		step = steps[0]
	}

	fm, err := pop.NewFileMigrator(migrationPath, tx)
	if err != nil {
		return err
	}

	err = fm.Down(step)
	if err != nil {
		return err
	}
	return nil
}

func (v *Velox) PopMigrateReset(tx *pop.Connection) error {
	var migrationPath = v.RootPath + "/migrations"

	fm, err := pop.NewFileMigrator(migrationPath, tx)
	if err != nil {
		return err
	}

	err = fm.Reset()
	if err != nil {
		return err
	}
	return nil
}
