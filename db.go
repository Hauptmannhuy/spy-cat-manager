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

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatal(err)
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

func (sql *sqlDB) createSpy(data Spy) error {
	_, err := sql.db.Exec("INSERT INTO spies (name, breed, experience, salary) VALUES ($1, $2, $3, $4)", data.Name, data.Breed, data.Experience, data.Salary)
	if err != nil {
		return err
	}
	return nil
}

func (sql *sqlDB) getSpies() ([]Spy, error) {
	data := []Spy{}
	rows, err := sql.db.Query("SELECT * FROM spies")
	if err != nil {
		return data, err
	}
	for rows.Next() {
		s := Spy{}
		rows.Scan(&s.Id, &s.Name, &s.Breed, &s.Salary, &s.Experience)
		data = append(data, s)
	}
	return data, nil
}

func (sql *sqlDB) getSpy(id int) (Spy, error) {
	res := Spy{}
	row := sql.db.QueryRow("SELECT * FROM spies WHERE id = $1", id)
	err := row.Scan(&res.Id, &res.Name, &res.Breed, &res.Experience, &res.Salary)
	if err != nil {
		return res, err
	}
	return res, nil
}

func (sql *sqlDB) updateSpy(id int, newData Spy) (Spy, error) {
	fmt.Println("salary", newData)
	row := sql.db.QueryRow(`
		UPDATE spies
		SET salary = $1
		WHERE id = $2
		RETURNING *
	`, newData.Salary, id)

	err := row.Scan(&newData.Id, &newData.Name, &newData.Breed, &newData.Salary, &newData.Experience)
	if err != nil {
		return newData, err
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

func (sql *sqlDB) createMission(data Mission) error {
	var query string
	var err error
	var id int
	if data.SpyID > 0 {
		query = "INSERT INTO missions (spy_id, completed) VALUES ($1, $2) RETURNING id"
		err = sql.db.QueryRow(query, data.SpyID, false).Scan(&id)

	} else {
		query = "INSERT INTO missions (spy_id, completed) VALUES ($1, $2) RETURNING id"
		err = sql.db.QueryRow(query, nil, false).Scan(&id)
	}
	if err != nil {
		return err
	}

	if len(data.Targets) > 0 {
		for _, target := range data.Targets {
			_, err = sql.db.Exec("INSERT INTO targets (name, mission_id, country, completed) VALUES ($1, $2, $3, $4)", target.Name, id, target.Country, target.CompleteState)
			if err != nil {
				return err
			}
		}
	}

	return err
}

func (sql *sqlDB) getMissions() ([]Mission, error) {
	data := []Mission{}
	rows, err := sql.db.Query("SELECT * FROM missions")
	if err != nil {
		return data, err
	}
	defer rows.Close()

	for rows.Next() {
		m := Mission{
			Targets: []Target{},
		}
		err := rows.Scan(&m.Id, &m.SpyID, &m.CompleteState)
		if err != nil {
			return nil, err
		}

		targetRows, err := sql.db.Query("SELECT * FROM targets WHERE mission_id = $1", m.Id)
		if err != nil {
			return nil, err
		}
		defer targetRows.Close()
		for targetRows.Next() {
			var t Target
			err = targetRows.Scan(&t.Id, &t.Name, &t.MissionId, &t.Country, &t.CompleteState)
			if err != nil {
				return nil, err
			}
			m.Targets = append(m.Targets, t)
		}

		data = append(data, m)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	fmt.Println(data)
	return data, nil
}

func (sql *sqlDB) getMission(id int) (Mission, error) {
	missionDest := Mission{
		Targets: []Target{},
	}
	var spyID interface{}

	err := sql.db.QueryRow(`
	SELECT * 
	FROM missions
	WHERE id = $1

	`, id).Scan(&missionDest.Id, &spyID, &missionDest.CompleteState)
	if err != nil {
		return missionDest, err
	}

	assertedID, ok := spyID.(int)
	if ok {
		missionDest.SpyID = assertedID
	}

	rows, err := sql.db.Query("SELECT * FROM targets WHERE mission_id = $1", id)
	if err != nil {
		return missionDest, err
	}
	for rows.Next() {
		var target Target
		err := rows.Scan(&target.Id, &target.Name, &target.MissionId, &target.Country, &target.CompleteState)
		if err != nil {
			return missionDest, err
		}
		missionDest.Targets = append(missionDest.Targets, target)
	}
	return missionDest, nil
}

func (sql *sqlDB) updateMission(id int, data *Mission) error {
	if data.CompleteState {
		_, err := sql.db.Exec("UPDATE missions SET completed = $2 WHERE id = $1", id, data.CompleteState)
		return err
	}
	if data.SpyID > 0 {
		fmt.Println("updating to ", data.SpyID)
		_, err := sql.db.Exec("UPDATE missions SET spy_id = $2 WHERE id = $1 ", id, data.SpyID)
		return err
	}

	return nil
}

func (sql *sqlDB) deleteMission(id int) error {
	_, err := sql.db.Exec("DELETE FROM targets WHERE mission_id = $1", id)
	if err != nil {
		return err
	}
	_, err = sql.db.Exec("DELETE FROM missions WHERE id = $1", id)
	return err
}

func (sql *sqlDB) addTargetToMission(id int, data Target) error {

	_, err := sql.db.Exec("INSERT INTO targets (name, country, mission_id) VALUES ($1, $2, $3)", data.Name, data.Country, id)
	return err
}

func (sql *sqlDB) updateTarget(missionID, targetID int, data Target) (Target, error) {
	var newTarget Target
	err := sql.db.QueryRow("UPDATE targets SET completed = $1 WHERE mission_id = $2 AND id = $3 RETURNING *",
		data.CompleteState, missionID, targetID).Scan(&newTarget.Id, &newTarget.Name, &newTarget.MissionId, &newTarget.Country, &newTarget.CompleteState)
	return newTarget, err
}

func (sql *sqlDB) deleteTarget(missionID, targetID int) error {
	_, err := sql.db.Exec("DELETE FROM targets WHERE id = $1 AND mission_id = $2", missionID, targetID)
	return err
}
