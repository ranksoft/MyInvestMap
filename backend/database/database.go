// /backend/database/database.go

package database

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func InitDB(filepath string) *sql.DB {
	db, err := sql.Open("sqlite3", filepath)
	if err != nil {
		log.Fatal(err)
	}

	if db == nil {
		log.Fatal("db nil")
	}

	createTable(db)
	return db
}

func createTable(db *sql.DB) {
	createTableSQL := `
    CREATE TABLE IF NOT EXISTS assets (
        id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
        stockTag TEXT NOT NULL,
        exchange TEXT NOT NULL,
		name TEXT,
        price REAL NOT NULL,
        quantity REAL NOT NULL,
		currentPrice REAL,
		isPurchase BOOLEAN NOT NULL DEFAULT true,
		createdAt DATETIME DEFAULT CURRENT_TIMESTAMP,
		updatedAt DATETIME DEFAULT CURRENT_TIMESTAMP
    );`

	_, err := db.Exec(createTableSQL)
	if err != nil {
		log.Fatal(err)
	}

}
