package server

import (
	"github.com/gofiber/fiber/v2"
	_ "github.com/joho/godotenv/autoload"
	"gorm.io/gorm"

	"github.com/hadronomy/lair-api/internal/database"
)

type Server interface {
	GetPort() int
	GetDB() *gorm.DB
	GetDBService() database.Service
	RegisterRoutes() *fiber.App
}
