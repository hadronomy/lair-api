package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"

	"lair-api/internal/models"
)

type GreatingOutput struct {
	Body struct {
		Message string `json:"name"`
	}
}

type GetLairsOutput struct {
	Body []models.Lair
}

type GetLairOutput struct {
	Body models.Lair
}

type UpdateLairsOutput struct {
	Body models.Lair
}

func (s *Server) RegisterRoutes() http.Handler {
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

	config := huma.DefaultConfig("Lair API", "1.0.0")
	api := humachi.New(r, config)

	huma.Register(api, huma.Operation{
		OperationID: "get-hello",
		Method:      http.MethodGet,
		Path:        "/hello",
		Summary:     "Greating",
	}, func(ctx context.Context, input *struct{}) (*GreatingOutput, error) {
		resp := &GreatingOutput{}
		resp.Body.Message = "Hello, World!"
		return resp, nil
	})

	huma.Register(api, huma.Operation{
		OperationID: "get-lairs",
		Method:      http.MethodGet,
		Path:        "/lairs",
		Summary:     "List Lairs",
	}, func(ctx context.Context, input *struct{}) (*GetLairsOutput, error) {
		var lairs []models.Lair
		s.db.GetDB().Find(&lairs)
		return &GetLairsOutput{
			Body: lairs,
		}, nil
	})

	huma.Register(api, huma.Operation{
		OperationID: "get-lair",
		Method:      http.MethodGet,
		Path:        "/lair/{lairID}",
		Summary:     "Get a Lair",
	}, func(ctx context.Context, input *struct {
		LairID string `path:"lairID"`
	}) (*GetLairOutput, error) {
		var lair models.Lair
		if s.db.GetDB().Where("id = ?", input.LairID).First(&lair).Error != nil {
			return nil, fmt.Errorf("lair not found")
		}
		return &GetLairOutput{
			Body: lair,
		}, nil
	})

	huma.Register(api, huma.Operation{
		OperationID: "update-lair",
		Method:      http.MethodPut,
		Path:        "/lair/{lairID}",
		Summary:     "Update a Lair",
	}, func(ctx context.Context, input *struct {
		ID   string             `path:"lairID"`
		Body models.LairRequest `json:"body"`
	}) (*UpdateLairsOutput, error) {
		var lair models.Lair
		if s.db.GetDB().Where("id = ?", input.ID).First(&lair).Error != nil {
			return nil, fmt.Errorf("lair not found")
		}
		s.db.GetDB().Model(&lair).Updates(input.Body)
		return &UpdateLairsOutput{
			Body: lair,
		}, nil
	})

	huma.Register(api, huma.Operation{
		OperationID: "delete-lair",
		Method:      http.MethodDelete,
		Path:        "/lair/{lairID}",
		Summary:     "Delete a Lair",
	}, func(ctx context.Context, input *struct {
		LairID string `path:"lairID"`
	}) (*UpdateLairsOutput, error) {
		var lair models.Lair
		if s.db.GetDB().Where("id = ?", input.LairID).First(&lair).Error != nil {
			return nil, fmt.Errorf("lair not found")
		}
		s.db.GetDB().Delete(&lair, input.LairID)
		return &UpdateLairsOutput{
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
