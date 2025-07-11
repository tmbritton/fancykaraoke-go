package db

import (
	"database/sql"
	"os"
	"strconv"
	"strings"

	_ "github.com/knaka/go-sqlite3-fts5"
	_ "github.com/mattn/go-sqlite3" // SQLite driver
)

type Migration struct {
	Filename string
	SQL      string
	Version  int
}

func setupDb(conn *SQLiteStore) error {
	sql := `
		PRAGMA journal_mode=WAL;
    CREATE TABLE IF NOT EXISTS schema_version (
      version INTEGER PRIMARY KEY,
      applied_at DATETIME DEFAULT CURRENT_TIMESTAMP
    );
	`
	_, err := conn.db.Exec(sql)
	return err
}

func getCurrentVersion(conn *SQLiteStore) (int, error) {
	var version int
	sqlStr := `
		SELECT version FROM schema_version
		ORDER BY applied_at DESC
		LIMIT 1
	`
	if err := conn.db.QueryRow(sqlStr).Scan(&version); err != nil {
		// return 0, err
		if err == sql.ErrNoRows {
			return 0, nil
		}
		return 0, err
	}
	return version, nil
}

func getSqlFiles() ([]string, error) {
	var files []string
	entries, err := os.ReadDir("./migrations")
	if err != nil {
		return files, err
	}

	for _, e := range entries {
		files = append(files, e.Name())
	}

	return files, nil
}

func filterDoneMigrations(current int, files []string) ([]string, error) {
	var filteredFiles []string

	for _, e := range files {
		split := strings.Split(e, "_")
		num, err := strconv.Atoi(split[0])
		if err != nil {
			return filteredFiles, err
		}
		if num > current {
			filteredFiles = append(filteredFiles, e)
		}
	}

	return filteredFiles, nil
}

func getFileContents(fileNames []string) ([]Migration, error) {
	var m []Migration

	for _, name := range fileNames {
		bytes, err := os.ReadFile("./migrations/" + name)
		if err != nil {
			return m, err
		}

		split := strings.Split(name, "_")
		num, err := strconv.Atoi(split[0])
		if err != nil {
			return m, err
		}

		m = append(m, Migration{
			Filename: name,
			SQL:      string(bytes),
			Version:  int(num),
		})
	}
	return m, nil
}

func performQueries(conn *SQLiteStore, migrations []Migration) error {
	for _, m := range migrations {
		_, err := conn.db.Exec(m.SQL)
		if err != nil {
			return err
		}
		if err := updateVersion(conn, m.Version); err != nil {
			return err
		}
	}
	return nil
}

func updateVersion(conn *SQLiteStore, version int) error {
	_, err := conn.db.Exec(`
		INSERT INTO schema_version(version)
		VALUES(?)
	`, version)
	if err != nil {
		return err
	}
	return nil
}

func DoMigrations(conn *SQLiteStore) error {
	if err := setupDb(conn); err != nil {
		return err
	}

	currentVersion, err := getCurrentVersion(conn)
	if err != nil {
		return err
	}

	files, err := getSqlFiles()
	if err != nil {
		return err
	}

	ff, err := filterDoneMigrations(currentVersion, files)
	if err != nil {
		return err
	}

	if len(ff) > 0 {
		migrations, err := getFileContents(ff)
		if err != nil {
			return err
		}
		if err := performQueries(conn, migrations); err != nil {
			return err
		}
	}

	return nil
}
