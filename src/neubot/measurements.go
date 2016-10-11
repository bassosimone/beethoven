package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

func MeasurementsAppend(runner *Runner) error {

	database, err := sql.Open("sqlite3", DefaultMeasurementsDb())
	if err != nil {
		return err
	}
	defer database.Close()

	log.Printf("database %s openned\n", DefaultMeasurementsDb())

	_, err = database.Exec(`
		CREATE TABLE IF NOT EXISTS measurements(
			id INTEGER PRIMARY KEY,
			timestamp NUMBER,
			test_name TEXT,
			status TEXT,
			test_id TEXT,
			stderr_path TEXT,
			stdout_path TEXT,
			workdir TEXT,
			cmd_line TEXT);`)
	if err != nil {
		log.Printf("cannot create table: %s\n", err)
		return err
	}

	log.Printf("table 'measurements' created if not exists\n")

	_, err = database.Exec(`
		INSERT INTO measurements(id, timestamp, test_name, status,
			test_id, stderr_path, stdout_path, workdir, cmd_line) VALUES(
			NULL, ?, ?, ?, ?, ?, ?, ?, ?);`,
		runner.Timestamp.Unix(), runner.TestName,
		runner.Status, runner.TestId, runner.StderrPath, runner.StdoutPath,
		runner.Workdir, runner.CmdLine)
	if err != nil {
		log.Printf("cannot append to database: %s\n", err)
		return err
	}

	log.Printf("measurement appended to database\n")
	return nil
}
