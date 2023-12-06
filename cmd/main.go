package main

import (
	"dungeons_helper_server/db"
	"dungeons_helper_server/internal/account"
	"dungeons_helper_server/internal/alignment"
	"dungeons_helper_server/internal/character"
	"dungeons_helper_server/internal/class"
	"dungeons_helper_server/internal/lobby"
	"dungeons_helper_server/internal/races"
	"dungeons_helper_server/internal/skills"
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
	alignmentHandler := alignment.NewHandler(alignment.NewService(alignment.NewRepository(dbConn.GetDB())))
	classHandler := class.NewHandler(class.NewService(class.NewRepository(dbConn.GetDB())))
	skillHandler := skills.NewHandler(skills.NewService(skills.NewRepository(dbConn.GetDB())))
	characterHandler := character.NewHandler(character.NewService(character.NewRepository(dbConn.GetDB())))
	lobbyHandler := lobby.NewHandler(lobby.NewService(lobby.NewRepository(dbConn.GetDB())))

	r := router.InitRouter(
		router.AccountRoutes(accountHandler),
		router.RacesRouter(racesHandler),
		router.SubracesRouter(subracesHandler),
		router.StatsRouter(statsHandler),
		router.AlignmentRouter(alignmentHandler),
		router.ClassRouter(classHandler),
		router.SkillsRouter(skillHandler),
		router.CharacterRouter(characterHandler),
		router.LobbyRouter(lobbyHandler),
	)

	if err := router.Start("localhost:5000", r); err != nil {
		log.Fatalf("Failed to start server: %s", err)
	}
}
