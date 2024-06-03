package main

import (
	"context"
	"net/http"
	"time"

	"github.com/charmbracelet/log"
	"github.com/danielgtaylor/huma/v2/humacli"

	"lair-api/internal/server"
)

type Options struct {
}

func main() {
	cli := humacli.New(func(hooks humacli.Hooks, options *Options) {
		server := server.NewServer()

		hooks.OnStart(func() {
			log.Infof("Server is running on http://localhost%s", server.Addr)
			if err := server.ListenAndServe(); err != nil {
				if err != http.ErrServerClosed {
					log.Fatalf("Server error: %v", err)
				}
			}
		})

		hooks.OnStop(func() {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			server.Shutdown(ctx)
		})
	})
	cli.Run()
}
