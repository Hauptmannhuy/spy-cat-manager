package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type sqlDB struct {
	db *sql.DB
}

func openAndMigrateDB() *sqlDB {
	godotenv.Load(".env")
	dataSourceName := os.Getenv("DATABASE_CREDS")

	db, err := sql.Open("postgres", dataSourceName)
	if err != nil {
		log.Fatal("Error opening DB", err)
	}

	driver, error := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatal(error)
	}

	m, err := migrate.NewWithDatabaseInstance("file://migrations", "postgres", driver)
	if err != nil {
		log.Fatalf("Could not start migration: %v", err)
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		log.Fatalf("Migration failed: %v", err)
	}
	return &sqlDB{
		db: db,
	}
}

func (sql *sqlDB) createSpy(data *Spy) error {
	_, err := sql.db.Exec("INSERT INTO spies (name, breed, experience, salary) VALUES ($1, $2, $3, $4)", data.Name, data.Breed, data.Experience, data.Salary)
	if err != nil {
		return err
	}
	return nil
}

func (sql *sqlDB) getSpy(id int) (*Spy, error) {
	res := &Spy{}
	row := sql.db.QueryRow("SELECT * FROM spies WHERE id = $1", id)
	err := row.Scan(&res.Id, &res.Name, &res.Breed, &res.Experience, &res.Salary)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (sql *sqlDB) updateSpy(id int, newData *Spy) (*Spy, error) {
	row := sql.db.QueryRow(`
		UPDATE spies
		SET salary = $1
		WHERE id = $2
		RETURNING *
	`, newData.Salary, id)

	err := row.Scan(&newData.Id, &newData.Name, &newData.Breed, &newData.Salary)
	if err != nil {
		return nil, err
	}
	return newData, nil
}

func (sql *sqlDB) deleteSpy(id int) error {
	_, err := sql.db.Exec("DELETE FROM spies WHERE id = $1", id)
	if err != nil {
		return err
	}
	return nil
}

func (sql *sqlDB) createMission(data *Mission) error {
	query := "INSERT INTO missions (spy_id, "
	closeQuery := "VALUES ($1, "
	args := []interface{}{data.SpyID}

	if len(data.Targets) > 0 {
		for i, target := range data.Targets {
			query += fmt.Sprintf("target_%d, ", i+1)
			closeQuery += fmt.Sprintf("$%d, ", i+2)
			args = append(args, target)
		}

		query = query[:len(query)-2] + ") "
		closeQuery = closeQuery[:len(closeQuery)-2] + ")"
		query += closeQuery
	} else {
		query = "INSERT INTO missions (spy_id) VALUES ($1)"
	}

	_, err := sql.db.Exec(query, args...)
	return err
}

func (sql *sqlDB) getMission(id int) (*Mission, error) {
	missionDest := &Mission{
		Targets: make([]Target, 3),
	}
	var dest struct {
		spy_id  int
		targets []int
	}
	row := sql.db.QueryRow(`
	SELECT m.spy_id, m.
  FROM missions as m
	JOIN targets as t
	ON m.id = t.mission_id

	`, id)
	err := row.Scan(&dest.spy_id, &dest.targets[0], &dest.targets[1], &dest.targets[2])
	if err != nil {
		return nil, err
	}

	for i, id := range dest.targets {
		row := sql.db.QueryRow("SELECT * FROM targets WHERE id = $1", id)
		row.Scan(&missionDest.Targets[i])
		if err != nil {
			return nil, err
		}
	}
	return missionDest, nil
}

func (sql *sqlDB) updateMission(id int, data *Mission) (*Mission, error) {}

func (sql *sqlDB) deleteMission(id int) error {}

func (sql *sqlDB) createTarget(data *Target) (*Target, error) {}

func (sql *sqlDB) updateTarget(id int, data *Target) (*Target, error) {}

func (sql *sqlDB) getTarget(id int) (*Target, error) {}

func (sql *sqlDB) deleteTarget(id *Target, data *Target) error {}
