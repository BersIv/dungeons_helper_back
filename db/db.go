package db

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"os"
)

type Database struct {
	db *sql.DB
}

func NewDatabase() (*Database, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

	dbUsername := os.Getenv("DB_USERNAME")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUsername, dbPassword, dbHost, dbPort, dbName)

	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		return nil, err
	}
	return &Database{db: db}, nil
}

func (d *Database) Close() {
	err := d.db.Close()
	if err != nil {
		return
	}
}

func (d *Database) GetDB() *sql.DB {
	return d.db
}
