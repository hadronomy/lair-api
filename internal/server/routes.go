package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"

	"lair-api/internal/models"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RealIP)
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Heartbeat("/ping"))
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)
	r.Use(render.SetContentType(render.ContentTypeJSON))

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		render.Render(w, r, ErrorNotFound(fmt.Errorf("the requested resource could not be found check the URL and try again")))
	})

	r.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		render.Render(w, r, ErrorInvalidRequest(fmt.Errorf("the requested method is not allowed for the requested resource")))
	})

	r.Get("/hello", s.HelloWorldHandler)

	r.Get("/health", s.healthHandler)

	r.Get("/panic", func(w http.ResponseWriter, r *http.Request) {
		panic("oh no")
	})

	r.Route("/lair", func(r chi.Router) {
		r.Get("/", s.getLairs)
		r.Post("/", s.createLair)
		r.Route("/{lairID}", func(r chi.Router) {
			r.Get("/", s.getLair)
			r.Put("/", s.updateLair)
			r.Delete("/", s.deleteLair)
		})
	})

	return r
}

type HelloWorldResponse struct {
	Message string `json:"message"`
}

func (s *Server) HelloWorldHandler(w http.ResponseWriter, r *http.Request) {
	response := HelloWorldResponse{
		Message: "Hello World",
	}

	jsonResp, err := json.Marshal(response)
	if err != nil {
		log.Fatalf("error handling JSON marshal. Err: %v", err)
	}

	_, _ = w.Write(jsonResp)
}

func (s *Server) getLairs(w http.ResponseWriter, r *http.Request) {
	var lairs []models.Lair
	s.db.GetDB().Find(&lairs)
	render.JSON(w, r, lairs)
}

func (s *Server) createLair(w http.ResponseWriter, r *http.Request) {
	var lair models.Lair
	if err := render.Bind(r, &lair); err != nil {
		render.Render(w, r, ErrorInvalidRequest(err))
		return
	}

	s.db.GetDB().Create(&lair)
	render.JSON(w, r, lair)
}

func (s *Server) getLair(w http.ResponseWriter, r *http.Request) {
	lairID := chi.URLParam(r, "lairID")
	var lair models.Lair
	if s.db.GetDB().Where("id = ?", lairID).First(&lair).Error != nil {
		render.Render(w, r, ErrorNotFound(fmt.Errorf("lair not found")))
		return
	}
	render.JSON(w, r, lair)
}

func (s *Server) updateLair(w http.ResponseWriter, r *http.Request) {
	lairID := chi.URLParam(r, "lairID")
	var lair models.Lair
	if err := render.Bind(r, &lair); err != nil {
		render.Render(w, r, ErrorInvalidRequest(err))
		return
	}
	if s.db.GetDB().Model(&models.Lair{}).Where("id = ?", lairID).Updates(lair).Error != nil {
		render.Render(w, r, ErrorNotFound(fmt.Errorf("lair not found")))
		return
	}
	render.JSON(w, r, lair)
}

func (s *Server) deleteLair(w http.ResponseWriter, r *http.Request) {
	lairID := chi.URLParam(r, "lairID")
	var lair models.Lair
	s.db.GetDB().Delete(&lair, lairID)
	render.JSON(w, r, lair)
}

func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	jsonResp, _ := json.Marshal(s.db.Health())
	_, _ = w.Write(jsonResp)
}
