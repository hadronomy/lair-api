package server

import (
	"net/http"

	_ "github.com/joho/godotenv/autoload"
	"gorm.io/gorm"

	"github.com/hadronomy/lair-api/internal/database"
)

type Server interface {
	GetPort() int
	GetDB() *gorm.DB
	GetDBService() database.Service
	RegisterRoutes() http.Handler
}
