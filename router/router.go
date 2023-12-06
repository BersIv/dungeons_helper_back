package router

import (
	"dungeons_helper_server/internal/account"
	"dungeons_helper_server/internal/alignment"
	"dungeons_helper_server/internal/character"
	"dungeons_helper_server/internal/class"
	"dungeons_helper_server/internal/races"
	"dungeons_helper_server/internal/skills"
	"dungeons_helper_server/internal/stats"
	"dungeons_helper_server/internal/subraces"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rr := &responseWriterWithStatus{ResponseWriter: w, statusCode: http.StatusOK}
		next.ServeHTTP(rr, r)
		status := rr.statusCode
		fmt.Printf("Received request: %s %s - Status: %d\n", r.Method, r.RequestURI, status)
	})
}

type responseWriterWithStatus struct {
	http.ResponseWriter
	statusCode int
}

func (r *responseWriterWithStatus) WriteHeader(statusCode int) {
	r.statusCode = statusCode
	r.ResponseWriter.WriteHeader(statusCode)
}

type Option func(router *mux.Router)

func AccountRoutes(accountHandler *account.Handler) Option {
	return func(r *mux.Router) {
		r.HandleFunc("/auth/registration", accountHandler.CreateAccount).Methods("POST")
		r.HandleFunc("/auth/email", accountHandler.Login).Methods("POST")
		r.HandleFunc("/logout", accountHandler.Logout).Methods("POST")
		r.HandleFunc("/auth/restore", accountHandler.RestorePassword).Methods("POST")
		r.HandleFunc("/account/nick", accountHandler.UpdateNickname).Methods("PATCH")
		r.HandleFunc("/account/password", accountHandler.UpdatePassword).Methods("PATCH")
	}
}

func RacesRouter(racesHandler *races.Handler) Option {
	return func(r *mux.Router) {
		r.HandleFunc("/getRaces", racesHandler.GetAllRaces).Methods("GET")
	}
}

func SubracesRouter(subracesHandler *subraces.Handler) Option {
	return func(r *mux.Router) {
		r.HandleFunc("/getSubraces", subracesHandler.GetAllSubraces).Methods("GET")
	}
}

func StatsRouter(statsHandler *stats.Handler) Option {
	return func(r *mux.Router) {
		r.HandleFunc("/getStatsById", statsHandler.GetStatsById).Methods("GET")
	}
}

func AlignmentRouter(alignmentHandler *alignment.Handler) Option {
	return func(r *mux.Router) {
		r.HandleFunc("/getAllAlignments", alignmentHandler.GetAllAlignments).Methods("GET")
	}
}

func ClassRouter(classHandler *class.Handler) Option {
	return func(r *mux.Router) {
		r.HandleFunc("/getAllClasses", classHandler.GetAllClasses).Methods("GET")
	}
}

func SkillsRouter(skillHandler *skills.Handler) Option {
	return func(r *mux.Router) {
		r.HandleFunc("/getAllSkills", skillHandler.GetAllSkills).Methods("GET")
	}
}

func CharacterRouter(characterHandler *character.Handler) Option {
	return func(r *mux.Router) {
		r.HandleFunc("/getAllCharactersByAccId", characterHandler.GetAllCharactersByAccId).Methods("GET")
		r.HandleFunc("/getCharacterById", characterHandler.GetCharacterById).Methods("GET")
		r.HandleFunc("/createCharacter", characterHandler.CreateCharacter).Methods("POST")
	}
}

func InitRouter(options ...Option) *mux.Router {
	r := mux.NewRouter()
	r.Use(LoggingMiddleware)
	for _, option := range options {
		option(r)
	}
	return r
}

func Start(addr string, r *mux.Router) error {
	return http.ListenAndServe(addr, r)
}
