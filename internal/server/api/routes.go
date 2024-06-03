package api

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"

	"lair-api/internal/logger"
	"lair-api/internal/server/resources"
)

// ConfigureApi configures the API routes and returns a huma.API instance.
// It takes a chi.Mux instance as a parameter and sets up the necessary routes and configurations.
// The function returns the configured huma.API instance.
func (s *APIServer) configureApi(r *chi.Mux) huma.API {
	config := huma.DefaultConfig("Lair API", "1.0.0")
	config.DocsPath = ""
	api := humachi.New(r, config)

	r.Get("/docs", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(`<!doctype html>
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
	</html>`))
	})

	return api
}

// ConfigureRouter configures the router for the server.
// It sets up various middlewares and handles not found and method not allowed routes.
// Returns the configured router.
func (s *APIServer) configureRouter() *chi.Mux {
	r := chi.NewMux()

	r.Use(middleware.RealIP)
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Heartbeat("/ping"))
	r.Use(middleware.Recoverer)
	r.Use(render.SetContentType(render.ContentTypeJSON))

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		render.JSON(w, r, &huma.ErrorModel{
			Title:    "Resource not found",
			Status:   http.StatusNotFound,
			Detail:   "The requested resource was not found, check the URL and try again.",
			Instance: r.URL.Path,
		})
	})

	r.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		render.JSON(w, r, &huma.ErrorModel{
			Title:    "Method not allowed",
			Status:   http.StatusMethodNotAllowed,
			Detail:   "The requested method is not allowed for the resource, check the documentation and try again.",
			Instance: r.URL.Path,
		})
	})

	return r
}

// RegisterRoutes registers the routes for the server.
// It returns an http.Handler that can be used to handle HTTP requests.
func (s *APIServer) RegisterRoutes() http.Handler {
	r := s.configureRouter()
	api := s.configureApi(r)

	resources.Register(api, s, &resources.LairResource{})
	resources.Register(api, s, &resources.HelloResource{})

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
	stdlog := logger.GetDefault().StandardLog()
	middleware.DefaultLogger =
		middleware.RequestLogger(
			&middleware.DefaultLogFormatter{
				Logger:  stdlog,
				NoColor: false,
			},
		)
}
