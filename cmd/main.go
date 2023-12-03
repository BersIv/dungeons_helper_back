package main

import (
	"dungeons_helper_server/db"
	"dungeons_helper_server/internal/account"
	"dungeons_helper_server/router"
	"log"
)

func main() {
	dbConn, err := db.NewDatabase()
	if err != nil {
		log.Fatalf("Could not initialize database connection: %s", err)
	}

	accountRep := account.NewRepository(dbConn.GetDB())
	accountSvc := account.NewService(accountRep)
	accountHandler := account.NewHandler(accountSvc)
	r := router.InitRouter(accountHandler)
	if err := router.Start("localhost:5000", r); err != nil {
		log.Fatalf("Failed to start server: %s", err)
	}
}
