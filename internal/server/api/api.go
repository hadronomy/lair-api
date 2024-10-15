package api

import (
	"os"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/hadronomy/lair-api/internal/database"

	"gorm.io/gorm"
)

type APIServer struct {
	port int

	db  database.Service
	app *fiber.App
}

func NewServer() *APIServer {
	env_port := os.Getenv("PORT")
	port, _ := strconv.Atoi(env_port)

	new_server := &APIServer{
		port: port,
		db:   database.New(),
	}

	new_server.app = new_server.RegisterRoutes()

	return new_server
}

func (s *APIServer) GetDBService() database.Service {
	return s.db
}

func (s *APIServer) GetDB() *gorm.DB {
	return s.db.GetDB()
}

func (s *APIServer) GetApp() *fiber.App {
	return s.app
}

func (s *APIServer) GetPort() int {
	return s.port
}
