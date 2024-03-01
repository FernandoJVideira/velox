package velox

import (
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func (v *Velox) MigrateUp(dsn string) error {
	m, err := migrate.New("file://"+v.RootPath+"/migrations", dsn)
	if err != nil {
		return err
	}
	defer m.Close()

	if err := m.Up(); err != nil {
		log.Println("Error applying migrations:", err)
		return err
	}
	return nil
}

func (v *Velox) MigrateDownAll(dsn string) error {
	m, err := migrate.New("file://"+v.RootPath+"/migrations", dsn)
	if err != nil {
		return err
	}
	defer m.Close()

	if err := m.Down(); err != nil {
		log.Println("Error applying migrations:", err)
		return err
	}
	return nil
}

func (v *Velox) Steps(dsn string, n int) error {
	m, err := migrate.New("file://"+v.RootPath+"/migrations", dsn)
	if err != nil {
		return err
	}
	defer m.Close()

	if err := m.Steps(n); err != nil {
		log.Println("Error applying migrations:", err)
		return err
	}

	return nil
}

func (v *Velox) MigrateForce(dsn string) error {
	m, err := migrate.New("file://"+v.RootPath+"/migrations", dsn)
	if err != nil {
		return err
	}
	defer m.Close()

	if err := m.Force(-1); err != nil {
		log.Println("Error applying migrations:", err)
		return err
	}
	return nil
}
