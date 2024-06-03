package logger

import (
	"os"
	"time"

	"github.com/charmbracelet/log"
)

var (
	Logger *log.Logger
)

func GetDefault() *log.Logger {
	return Logger
}

func SetDefault(logger *log.Logger) {
	Logger = logger
	log.SetDefault(logger)
}

func init() {
	SetDefault(
		log.NewWithOptions(os.Stdout,
			log.Options{
				ReportTimestamp: true,
				TimeFormat:      time.DateTime,
			}),
	)
}
