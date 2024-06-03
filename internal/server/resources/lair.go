package resources

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"

	"lair-api/internal/crypto"
	"lair-api/internal/models"
	"lair-api/internal/server"
)

type GetLairsResponse struct {
	Body []models.Lair
}

type GetLairResponse struct {
	Body models.Lair
}

type UpdateLairsResponse struct {
	Body models.Lair
}

type LairResource struct{}

func (l *LairResource) Init(api huma.API, s server.Server) {
	huma.Register(api, huma.Operation{
		OperationID: "get-lairs",
		Method:      http.MethodGet,
		Path:        "/lairs",
		Summary:     "List Lairs",
		Tags:        []string{"Lairs"},
	}, func(ctx context.Context, input *struct{}) (*GetLairsResponse, error) {
		var lairs []models.Lair
		s.GetDB().Find(&lairs)
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
		id, err := crypto.GenerateID()
		if err != nil {
			return nil, huma.Error500InternalServerError("failed to generate ID")
		}
		lair := models.Lair{
			Model: models.Model{
				ID: id,
			},
			Name:    input.Body.Name,
			Owner:   input.Body.Owner,
			Private: input.Body.Private,
		}
		if err := s.GetDB().Create(&lair).Error; err != nil {
			return nil, huma.Error500InternalServerError("failed to create lair")
		}
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
		if s.GetDB().Where("id = ?", input.LairID).First(&lair).Error != nil {
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
		if s.GetDB().Where("id = ?", input.ID).First(&lair).Error != nil {
			return nil, huma.Error404NotFound("lair not found")
		}
		s.GetDB().Model(&lair).Updates(input.Body)
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
		if s.GetDB().Where("id = ?", input.LairID).First(&lair).Error != nil {
			return nil, huma.Error404NotFound("lair not found")
		}
		s.GetDB().Delete(&lair, input.LairID)
		return &UpdateLairsResponse{
			Body: lair,
		}, nil
	})
}
