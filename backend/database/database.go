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

	createTables(db)
	return db
}

func createTables(db *sql.DB) {

	createUsersTableSQL := `
    CREATE TABLE IF NOT EXISTS users (
        id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
        username TEXT NOT NULL UNIQUE,
        email TEXT NOT NULL UNIQUE,
        password TEXT NOT NULL
    );`

	_, err := db.Exec(createUsersTableSQL)
	if err != nil {
		log.Fatal(err)
	}

	createAssetsTableSQL := `
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
	    updatedAt DATETIME DEFAULT CURRENT_TIMESTAMP,
	    user_id INTEGER,
	    FOREIGN KEY (user_id) REFERENCES users(id)
	);`

	_, err = db.Exec(createAssetsTableSQL)
	if err != nil {
		log.Fatal(err)
	}

	createApiKeysTableSQL := `
    CREATE TABLE IF NOT EXISTS api_keys (
        id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
        user_id INTEGER NOT NULL,
        api_key TEXT NOT NULL,
        FOREIGN KEY (user_id) REFERENCES users(id)
    );`

	_, err = db.Exec(createApiKeysTableSQL)
	if err != nil {
		log.Fatal(err)
	}
}
