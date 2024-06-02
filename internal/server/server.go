package server

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/joho/godotenv/autoload"

	"lair-api/internal/database"
)

type Server struct {
	port int

	db database.Service
}

func NewServer() *http.Server {
	env_port := os.Getenv("PORT")
	port, _ := strconv.Atoi(env_port)
	NewServer := &Server{
		port: port,

		db: database.New(),
	}
	println("Server is running on port: ", NewServer.port)

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
