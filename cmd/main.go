package main

import (
	"dungeons_helper_server/db"
	"log"
)

func main() {
	dbConn, err := db.NewDatabase()
	if err != nil {
		log.Fatalf("Could not initialize database connection: %s", err)
	}
	print(dbConn)
	dbConn.Close()
}
