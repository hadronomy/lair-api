package api

import (
	"fmt"
	"lair-api/internal/database"
	"net/http"
	"os"
	"strconv"
	"time"

	"gorm.io/gorm"
)

type APIServer struct {
	port int

	db database.Service
}

func NewServer() *http.Server {
	env_port := os.Getenv("PORT")
	port, _ := strconv.Atoi(env_port)
	NewServer := &APIServer{
		port: port,

		db: database.New(),
	}

	// Declare Server config
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", NewServer.port),
		Handler:      NewServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}

func (s *APIServer) GetDB() *gorm.DB {
	return s.db.GetDB()
}

func (s *APIServer) GetPort() int {
	return s.port
}
