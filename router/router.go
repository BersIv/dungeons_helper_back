package router

import (
	"dungeons_helper_server/internal/account"
	"fmt"
	"net/http"
)
import "github.com/gorilla/mux"

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

func InitRouter(accountHandler *account.Handler) *mux.Router {
	r := mux.NewRouter()
	r.Use(LoggingMiddleware)
	r.HandleFunc("/auth/registration", accountHandler.CreateAccount).Methods("POST")
	r.HandleFunc("/auth/email", accountHandler.Login).Methods("POST")
	r.HandleFunc("/logout", accountHandler.Logout).Methods("POST")
	r.HandleFunc("/auth/restore", accountHandler.RestorePassword).Methods("POST")
	r.HandleFunc("/account/nick", accountHandler.UpdateNickname).Methods("PATCH")
	return r
}

func Start(addr string, r *mux.Router) error {
	return http.ListenAndServe(addr, r)
}
