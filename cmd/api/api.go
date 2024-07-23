package api

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/ganthology/go-ecom-api/service/user"
	"github.com/gorilla/mux"
)

type APIServer struct {
	address  string
	database *sql.DB
}

func NewAPIServer(address string, database *sql.DB) *APIServer {
	return &APIServer{
		address:  address,
		database: database,
	}
}

func (s *APIServer) Run() error {
	router := mux.NewRouter()
	subrouter := router.PathPrefix("/api/v1").Subrouter()

	// this is dependency injection
	userStore := user.NewStore(s.database)
	userHandler := user.NewHandler(userStore)
	userHandler.RegisterRoutes(subrouter)

	log.Println("Starting server on", s.address)

	return http.ListenAndServe(s.address, router)
}
