package db

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

func StartDatabase(DBPath string) (*sql.DB, error) {
	dbDir := filepath.Dir(DBPath)
	os.MkdirAll(dbDir, 0755)
	database, err := sql.Open("sqlite3", DBPath)
	if err != nil {
		// fmt.Println("Error opening database:", err)
		return nil, err
	}

	_, err = database.Exec("PRAGMA foreign_keys = ON;")
	if err != nil {
		database.Close()
		return nil, err
	}

	database.SetMaxOpenConns(1)

	err = database.Ping()
	if err != nil {
		database.Close()
		return nil, err
	}

	err = RunMigrations(database)
	if err != nil {
		database.Close()
		fmt.Println("Error running migrations:", err)
		return nil, err
	}

	fmt.Println("db is open now")

	return database, nil
}

func RunMigrations(database *sql.DB) error { //todo: add dir string in the input of function
	fileNames, err := GetFilesNames("backend/db/migration")
	if err != nil {
		return err
	}

	for _, f := range fileNames {
		fullPath := filepath.Join("backend/db/migration/" + f)

		query, err := os.ReadFile(fullPath)
		if err != nil {
			return err
		}

		_, err = database.Exec(string(query))
		if err != nil {
			return err
		}

		fmt.Printf("Migrated: %s\n", f)
	}

	return nil
}

func GetFilesNames(path string) ([]string, error) { //?      ./db/migration

	files, err := os.ReadDir(path)
	if err != nil {
		fmt.Printf("Error reading directory %s: %v\n", path, err)
		return nil, err
	}

	var filesNames []string

	for _, file := range files {
		if !file.IsDir() && filepath.Ext(file.Name()) == ".sql" { //Extention
			filesNames = append(filesNames, file.Name())
		}
	}

	return filesNames, nil
}
