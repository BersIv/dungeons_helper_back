package main

import (
	"dungeons_helper_server/db"
	"dungeons_helper_server/internal/account"
	"dungeons_helper_server/internal/races"
	"dungeons_helper_server/internal/stats"
	"dungeons_helper_server/internal/subraces"
	"dungeons_helper_server/router"
	"log"
)

func main() {
	dbConn, err := db.NewDatabase()
	if err != nil {
		log.Fatalf("Could not initialize database connection: %s", err)
	}

	accountHandler := account.NewHandler(account.NewService(account.NewRepository(dbConn.GetDB())))
	racesHandler := races.NewHandler(races.NewService(races.NewRepository(dbConn.GetDB())))
	subracesHandler := subraces.NewHandler(subraces.NewService(subraces.NewRepository(dbConn.GetDB())))
	statsHandler := stats.NewHandler(stats.NewService(stats.NewRepository(dbConn.GetDB())))

	r := router.InitRouter(
		router.AccountRoutes(accountHandler),
		router.RacesRouter(racesHandler),
		router.SubracesRouter(subracesHandler),
		router.StatsRouter(statsHandler),
	)

	if err := router.Start("localhost:5000", r); err != nil {
		log.Fatalf("Failed to start server: %s", err)
	}
}
