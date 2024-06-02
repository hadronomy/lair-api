package main

import (
	"context"
	"log"
	"time"

	"github.com/danielgtaylor/huma/v2/humacli"

	"lair-api/internal/server"
)

type Options struct {
}

func main() {
	cli := humacli.New(func(hooks humacli.Hooks, options *Options) {
		server := server.NewServer()

		hooks.OnStart(func() {
			if err := server.ListenAndServe(); err != nil {
				log.Printf("failed to start server: %v", err)
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
