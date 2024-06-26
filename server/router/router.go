package router

import (
	"dungeons_helper/internal/account"
	"dungeons_helper/internal/alignment"
	"dungeons_helper/internal/character"
	"dungeons_helper/internal/class"
	"dungeons_helper/internal/races"
	"dungeons_helper/internal/skills"
	"dungeons_helper/internal/stats"
	"dungeons_helper/internal/subraces"
	"dungeons_helper/internal/websocket"
	"dungeons_helper/util"
	"net/http"

	"github.com/gorilla/mux"
)

type Option func(router *mux.Router)

func AccountRoutes(accountHandler *account.Handler) Option {
	return func(r *mux.Router) {
		r.HandleFunc("/auth/registration", accountHandler.CreateAccount).Methods("POST")
		r.HandleFunc("/auth/byEmail", accountHandler.Login).Methods("POST")
		r.HandleFunc("/logout", accountHandler.Logout).Methods("POST")
		r.HandleFunc("/auth/restorePassword", accountHandler.RestorePassword).Methods("POST")
		r.HandleFunc("/account/change/nickname", accountHandler.UpdateNickname).Methods("PATCH")
		r.HandleFunc("/account/change/password", accountHandler.UpdatePassword).Methods("PATCH")
		r.HandleFunc("/auth/google/login", accountHandler.LoginGoogle).Methods("GET")
	}
}

func RacesRouter(racesHandler *races.Handler) Option {
	return func(r *mux.Router) {
		r.HandleFunc("/race/getRaces", racesHandler.GetAllRaces).Methods("GET")
	}
}

func SubracesRouter(subracesHandler *subraces.Handler) Option {
	return func(r *mux.Router) {
		r.HandleFunc("/subrace/getSubraces", subracesHandler.GetAllSubraces).Methods("GET")
	}
}

func StatsRouter(statsHandler *stats.Handler) Option {
	return func(r *mux.Router) {
		r.HandleFunc("/getStatsById", statsHandler.GetStatsById).Methods("GET")
	}
}

func AlignmentRouter(alignmentHandler *alignment.Handler) Option {
	return func(r *mux.Router) {
		r.HandleFunc("/getAlignments", alignmentHandler.GetAllAlignments).Methods("GET")
	}
}

func ClassRouter(classHandler *class.Handler) Option {
	return func(r *mux.Router) {
		r.HandleFunc("/getClasses", classHandler.GetAllClasses).Methods("GET")
	}
}

func SkillsRouter(skillHandler *skills.Handler) Option {
	return func(r *mux.Router) {
		r.HandleFunc("/getSkills", skillHandler.GetAllSkills).Methods("GET")
	}
}

func CharacterRouter(characterHandler *character.Handler) Option {
	return func(r *mux.Router) {
		r.HandleFunc("/getAllCharactersByAccId", characterHandler.GetAllCharactersByAccId).Methods("GET")
		r.HandleFunc("/getCharacterById", characterHandler.GetCharacterById).Methods("GET")
		r.HandleFunc("/createCharacter", characterHandler.CreateCharacter).Methods("POST")
	}
}

// func LobbyRouter(lobbyHandler *lobby.Handler) Option {
// 	return func(r *mux.Router) {
// 		r.HandleFunc("/createLobby", lobbyHandler.CreateLobby).Methods("POST")
// 		r.HandleFunc("/getAllLobby", lobbyHandler.GetAllLobby).Methods("POST")
// 		r.HandleFunc("/joinLobby", lobbyHandler.JoinLobby)
// 	}
// }

func WebsocketRouter(wsHandler *websocket.Handler) Option {
	return func(r *mux.Router) {
		r.HandleFunc("/ws/join", wsHandler.JoinLobby)
		r.HandleFunc("/ws/create", wsHandler.CreateLobby)
	}
}

func InitRouter(options ...Option) *mux.Router {
	r := mux.NewRouter()
	r.Use(util.LoggingMiddleware)

	for _, option := range options {
		option(r)
	}
	return r
}

func Start(addr string, r *mux.Router) error {
	return http.ListenAndServe(addr, r)
}
