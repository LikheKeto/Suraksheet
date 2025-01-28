package main

import (
	"log"
	"os"

	"github.com/LikheKeto/Suraksheet/config"
	"github.com/LikheKeto/Suraksheet/db"
	"github.com/golang-migrate/migrate/v4"
	pgsql "github.com/golang-migrate/migrate/v4/database/pgx"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	db := db.NewSQLStorage(db.DBConfig{
		User:     config.Envs.DBUser,
		Host:     config.Envs.DBHost,
		Port:     config.Envs.DBPort,
		Password: config.Envs.DBPassword,
		DBname:   config.Envs.DBName,
	})

	driver, err := pgsql.WithInstance(db, &pgsql.Config{})
	if err != nil {
		log.Fatal(err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://cmd/migrate/migrations",
		"mysql",
		driver,
	)

	if err != nil {
		log.Fatal(err)
	}

	cmd := os.Args[(len(os.Args) - 1)]
	if cmd == "up" {
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			log.Fatal(err)
		}
	}
	if cmd == "down" {
		if err := m.Down(); err != nil && err != migrate.ErrNoChange {
			log.Fatal(err)
		}
	}
}
