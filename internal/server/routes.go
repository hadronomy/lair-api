package server

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"

	"lair-api/internal/models"
)

type GreatingResponse struct {
	Body struct {
		Message string `json:"name"`
	}
}

type GetLairsResponse struct {
	Body []models.Lair
}

type GetLairResponse struct {
	Body models.Lair
}

type UpdateLairsResponse struct {
	Body models.Lair
}

// ConfigureApi configures the API routes and returns a huma.API instance.
// It takes a chi.Mux instance as a parameter and sets up the necessary routes and configurations.
// The function returns the configured huma.API instance.
func (s *Server) configureApi(r *chi.Mux) huma.API {
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
func (s *Server) configureRouter() *chi.Mux {
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
func (s *Server) RegisterRoutes() http.Handler {
	r := s.configureRouter()
	api := s.configureApi(r)

	huma.Register(api, huma.Operation{
		OperationID: "get-hello",
		Method:      http.MethodGet,
		Path:        "/hello",
		Summary:     "Greating",
	}, func(ctx context.Context, input *struct{}) (*GreatingResponse, error) {
		resp := &GreatingResponse{}
		resp.Body.Message = "Hello, World!"
		return resp, nil
	})

	huma.Register(api, huma.Operation{
		OperationID: "get-lairs",
		Method:      http.MethodGet,
		Path:        "/lairs",
		Summary:     "List Lairs",
		Tags:        []string{"Lairs"},
	}, func(ctx context.Context, input *struct{}) (*GetLairsResponse, error) {
		var lairs []models.Lair
		s.db.GetDB().Find(&lairs)
		return &GetLairsResponse{
			Body: lairs,
		}, nil
	})

	huma.Register(api, huma.Operation{
		OperationID: "create-lair",
		Method:      http.MethodPost,
		Path:        "/lair",
		Summary:     "Create a Lair",
		Tags:        []string{"Lairs"},
	}, func(ctx context.Context, input *struct {
		Body models.LairRequest `json:"body"`
	}) (*UpdateLairsResponse, error) {
		lair := models.Lair{
			Name:    input.Body.Name,
			Owner:   input.Body.Owner,
			Private: input.Body.Private,
		}
		s.db.GetDB().Create(&lair)
		return &UpdateLairsResponse{
			Body: lair,
		}, nil
	})

	huma.Register(api, huma.Operation{
		OperationID: "get-lair",
		Method:      http.MethodGet,
		Path:        "/lair/{lairID}",
		Summary:     "Get a Lair",
		Tags:        []string{"Lairs"},
	}, func(ctx context.Context, input *struct {
		LairID string `path:"lairID"`
	}) (*GetLairResponse, error) {
		var lair models.Lair
		if s.db.GetDB().Where("id = ?", input.LairID).First(&lair).Error != nil {
			return nil, huma.Error404NotFound("lair not found")
		}
		return &GetLairResponse{
			Body: lair,
		}, nil
	})

	huma.Register(api, huma.Operation{
		OperationID: "update-lair",
		Method:      http.MethodPut,
		Path:        "/lair/{lairID}",
		Summary:     "Update a Lair",
		Tags:        []string{"Lairs"},
	}, func(ctx context.Context, input *struct {
		ID   string             `path:"lairID"`
		Body models.LairRequest `json:"body"`
	}) (*UpdateLairsResponse, error) {
		var lair models.Lair
		if s.db.GetDB().Where("id = ?", input.ID).First(&lair).Error != nil {
			return nil, huma.Error404NotFound("lair not found")
		}
		s.db.GetDB().Model(&lair).Updates(input.Body)
		return &UpdateLairsResponse{
			Body: lair,
		}, nil
	})

	huma.Register(api, huma.Operation{
		OperationID: "delete-lair",
		Method:      http.MethodDelete,
		Path:        "/lair/{lairID}",
		Summary:     "Delete a Lair",
		Tags:        []string{"Lairs"},
	}, func(ctx context.Context, input *struct {
		LairID string `path:"lairID"`
	}) (*UpdateLairsResponse, error) {
		var lair models.Lair
		if s.db.GetDB().Where("id = ?", input.LairID).First(&lair).Error != nil {
			return nil, huma.Error404NotFound("lair not found")
		}
		s.db.GetDB().Delete(&lair, input.LairID)
		return &UpdateLairsResponse{
			Body: lair,
		}, nil
	})

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
