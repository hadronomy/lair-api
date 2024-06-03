package server

import (
	"net/http"

	_ "github.com/joho/godotenv/autoload"
	"gorm.io/gorm"
)

type Server interface {
	GetPort() int
	GetDB() *gorm.DB
	RegisterRoutes() http.Handler
}
