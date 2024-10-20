package main

import (
	"context"
	"net/http"
	"time"

	"github.com/charmbracelet/log"
	"github.com/danielgtaylor/huma/v2/humacli"

	"github.com/hadronomy/lair-api/internal/server/api"
)

type Options struct{}

func main() {
	cli := humacli.New(func(hooks humacli.Hooks, options *Options) {
		server := api.NewServer()
		api := server.GetApp()

		hooks.OnStart(func() {
			if err := api.Listen(":3000"); err != nil {
				if err != http.ErrServerClosed {
					log.Fatalf("Server error: %v", err)
				}
			}
		})

		hooks.OnStop(func() {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			api.ShutdownWithContext(ctx)
		})
	})
	cli.Run()
}
