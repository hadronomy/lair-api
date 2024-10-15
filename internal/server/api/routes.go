package api

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humafiber"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"

	// "github.com/hadronomy/lair-api/internal/logger"
	"github.com/hadronomy/lair-api/internal/server/resources"
)

// ConfigureApi configures the API routes and returns a huma.API instance.
// It takes a chi.Mux instance as a parameter and sets up the necessary routes and configurations.
// The function returns the configured huma.API instance.
func (s *APIServer) configureApi(r *fiber.App) huma.API {
	config := huma.DefaultConfig("Lair API", "1.0.0")
	config.DocsPath = ""
	api := humafiber.New(r, config)

	r.Get("/docs", func(c *fiber.Ctx) error {
		c.Set("Content-Type", "text/html")
		c.SendString(`<!doctype html>
	<html>
		<head>
			<title>API Reference</title>
			<meta charset="utf-8" />
			<meta
				name="viewport"
				content="width=device-width, initial-scale=1" />
		</head>
		<body>
			<script
			id="api-reference"
			data-url="/openapi.json"
			data-proxy-url="https://proxy.scalar.com"></script>

			<!-- Optional: You can set a full configuration object like this: -->
			<script>
				var configuration = {
					theme: 'bluePlanet',
				}

				document.getElementById('api-reference').dataset.configuration =
					JSON.stringify(configuration)
			</script>
			<script src="https://cdn.jsdelivr.net/npm/@scalar/api-reference"></script>
		</body>
	</html>`)
		return nil
	})

	return api
}

// ConfigureRouter configures the router for the server.
// It sets up various middlewares and handles not found and method not allowed routes.
// Returns the configured router.
func (s *APIServer) configureRouter() *fiber.App {
	r := fiber.New(fiber.Config{
		Prefork:       true,
		CaseSensitive: true,
		StrictRouting: true,
		ServerHeader:  "LAIR-API",
		AppName:       "LAIR-API v1.0.0",
	})

	// r.Use(middleware.RealIP)
	// r.Use(middleware.RequestID)
	// r.Use(middleware.Logger)
	// r.Use(middleware.Heartbeat("/ping"))
	// r.Use(middleware.Recoverer)
	// r.Use(render.SetContentType(render.ContentTypeJSON))
	r.Use(recover.New())
	r.Use(requestid.New())
	r.Use(logger.New())

	r.Use(func(c *fiber.Ctx) error {
		if c.Response().StatusCode() == http.StatusNotFound {
			return c.Status(http.StatusNotFound).JSON(&huma.ErrorModel{
				Title:    "Resource not found",
				Status:   http.StatusNotFound,
				Detail:   "The requested resource was not found, check the URL and try again.",
				Instance: c.OriginalURL(),
			})
		}
		return c.Next()
	})

	r.Use(func(c *fiber.Ctx) error {
		if c.Response().StatusCode() == http.StatusMethodNotAllowed {
			return c.Status(http.StatusMethodNotAllowed).JSON(&huma.ErrorModel{
				Title:    "Method not allowed",
				Status:   http.StatusMethodNotAllowed,
				Detail:   "The requested method is not allowed for the resource, check the documentation and try again.",
				Instance: c.OriginalURL(),
			})
		}
		return c.Next()
	})

	return r
}

// RegisterRoutes registers the routes for the server.
// It returns an http.Handler that can be used to handle HTTP requests.
func (s *APIServer) RegisterRoutes() *fiber.App {
	r := s.configureRouter()
	api := s.configureApi(r)

	resources.Register(api, s, &resources.LairResource{})
	resources.Register(api, s, &resources.HelloResource{})
	resources.Register(api, s, &resources.HealthResource{})

	// TODO: Add health check route
	// huma.Register(api, huma.Operation{
	// 	OperationID: "check-health",
	// 	Method:      http.MethodGet,
	// 	Path:        "/health",
	// 	Summary:     "Check Health",
	// }, func(ctx context.Context, input *struct{}) (map[string]string, error) {
	// 	return s.db.Health(), nil
	// })
	// r.Get("/health", s.healthHandler)
	return r
}

func init() {
	// TODO: Abtract away the logger
	// stdlog := logger.GetDefault().StandardLog()
	// middleware.DefaultLogger = middleware.RequestLogger(
	// 	&middleware.DefaultLogFormatter{
	// 		Logger:  stdlog,
	// 		NoColor: false,
	// 	},
	// )
}
