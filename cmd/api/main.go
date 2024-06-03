package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/charmbracelet/lipgloss"
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
			printBanner(server)
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

func printBanner(server *http.Server) {
	style := lipgloss.NewStyle().
		Padding(1, 10).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("63")).
		BorderTop(true).
		BorderBottom(true).
		BorderLeft(true).
		BorderRight(true)

	title := lipgloss.NewStyle().
		PaddingBottom(1).
		Foreground(lipgloss.Color("255")).
		Bold(true).
		Render(
			"LAIR API",
		)
	url := lipgloss.NewStyle().
		Foreground(lipgloss.Color("86")).
		Render(
			fmt.Sprintf("http://localhost%s", server.Addr),
		)
	descriptionStyle := lipgloss.NewStyle().Faint(true)
	description := descriptionStyle.Render(
		fmt.Sprintf("(bound on host 0.0.0.0 and port %s)", server.Addr),
	)
	block := lipgloss.JoinVertical(lipgloss.Center, title, url, description)
	fmt.Println()
	fmt.Println(style.Render(block))
}
