package database

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/charmbracelet/log"
	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/joho/godotenv/autoload"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/hadronomy/lair-api/internal/models"
)

// Service represents a service that interacts with a database.
type Service interface {
	// Health returns a map of health status information.
	// The keys and values in the map are service-specific.
	Health() interface{}

	// GetDB returns the database connection.
	GetDB() *gorm.DB

	// // Close terminates the database connection.
	// // It returns an error if the connection cannot be closed.
	// Close() error
}

type dbService struct {
	db *gorm.DB
}

var (
	database   = os.Getenv("DB_DATABASE")
	password   = os.Getenv("DB_PASSWORD")
	username   = os.Getenv("DB_USERNAME")
	port       = os.Getenv("DB_PORT")
	host       = os.Getenv("DB_HOST")
	dbInstance *dbService
)

func New() Service {
	// Reuse Connection
	if dbInstance != nil {
		return dbInstance
	}
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", username, password, host, port, database)
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN: connStr,
	}), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("failed to get db instance: %v", err)
	}
	sqlDB.SetConnMaxIdleTime(time.Hour)
	sqlDB.SetConnMaxLifetime(24 * time.Hour)
	sqlDB.SetMaxIdleConns(100)
	sqlDB.SetMaxOpenConns(200)

	// db.Use(
	// 	dbresolver.Register(dbresolver.Config{
	// 		Sources: []gorm.Dialector{postgres.Open(connStr)},
	// 	}).
	// 		SetConnMaxIdleTime(time.Hour).
	// 		SetConnMaxLifetime(24 * time.Hour).
	// 		SetMaxIdleConns(100).
	// 		SetMaxOpenConns(200),
	// )

	db.AutoMigrate(&models.Lair{})

	dbInstance = &dbService{
		db: db,
	}

	return dbInstance
}

func (s *dbService) GetDB() *gorm.DB {
	return s.db
}

type Status string

const (
	StatusUp   Status = "up"
	StatusDown Status = "down"
)

type DatabaseHealth struct {
	Status            Status `json:"status"`
	Message           string `json:"message,omitempty" required:"false"`
	Error             string `json:"error,omitempty" required:"false"`
	OpenConnections   int    `json:"open_connections,omitempty" required:"false"`
	InUse             int    `json:"in_use,omitempty" required:"false"`
	Idle              int    `json:"idle,omitempty" required:"false"`
	WaitCount         int64  `json:"wait_count,omitempty" required:"false"`
	WaitDuration      string `json:"wait_duration,omitempty" required:"false"`
	MaxIdleClosed     int64  `json:"max_idle_closed,omitempty" required:"false"`
	MaxLifetimeClosed int64  `json:"max_lifetime_closed,omitempty" required:"false"`
}

// Health checks the health of the database connection by pinging the database.
// It returns a map with keys indicating various health statistics.
func (s *dbService) Health() interface{} {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	stats := DatabaseHealth{}

	// Ping the database
	db, err := s.db.DB()
	if err != nil {
		stats.Status = StatusDown
		stats.Error = fmt.Sprintf("db down: %v", err)
		log.Error(fmt.Sprintf("db down: %v", err)) // Log the error and terminate the program
		return stats
	}

	err = db.PingContext(ctx)
	if err != nil {
		stats.Status = StatusDown
		stats.Error = fmt.Sprintf("db down: %v", err)
		log.Errorf(fmt.Sprintf("db down: %v", err)) // Log the error and terminate the program
		return stats
	}

	// Database is up, add more statistics
	stats.Status = StatusUp
	stats.Message = "It's healthy"

	// Get database stats (like open connections, in use, idle, etc.)
	dbStats := db.Stats()
	stats.OpenConnections = dbStats.OpenConnections
	stats.InUse = dbStats.InUse
	stats.Idle = dbStats.Idle
	stats.WaitCount = dbStats.WaitCount
	if dbStats.WaitDuration > 0 {
		stats.WaitDuration = dbStats.WaitDuration.String()
	}
	stats.MaxIdleClosed = dbStats.MaxIdleClosed
	stats.MaxLifetimeClosed = dbStats.MaxLifetimeClosed

	// Evaluate stats to provide a health message
	if dbStats.OpenConnections > 40 { // Assuming 50 is the max for this example
		stats.Message = "The database is experiencing heavy load."
	}

	if dbStats.WaitCount > 1000 {
		stats.Message = "The database has a high number of wait events, indicating potential bottlenecks."
	}

	if dbStats.MaxIdleClosed > int64(dbStats.OpenConnections)/2 {
		stats.Message = "Many idle connections are being closed, consider revising the connection pool settings."
	}

	if dbStats.MaxLifetimeClosed > int64(dbStats.OpenConnections)/2 {
		stats.Message = "Many connections are being closed due to max lifetime, consider increasing max lifetime or revising the connection usage pattern."
	}

	return stats
}

// // Close closes the database connection.
// // It logs a message indicating the disconnection from the specific database.
// // If the connection is successfully closed, it returns nil.
// // If an error occurs while closing the connection, it returns the error.
// func (s *service) Close() error {
// 	log.Printf("Disconnected from database: %s", database)
// 	return s.db.Close()
// }
