package resources

import (
	"context"

	"github.com/danielgtaylor/huma/v2"

	"github.com/hadronomy/lair-api/internal/database"
	"github.com/hadronomy/lair-api/internal/server"
)

type HealthResource struct{}

type HealthResponse struct {
	Body database.DatabaseHealth
}

func (l *HealthResource) Init(api huma.API, s server.Server) {
	huma.Register(api, huma.Operation{
		OperationID: "get-health",
		Method:      "GET",
		Path:        "/health",
		Summary:     "Health Check",
		Tags:        []string{"Health"},
	}, func(ctx context.Context, input *struct{}) (*HealthResponse, error) {
		return &HealthResponse{
			Body: s.GetDBService().Health().(database.DatabaseHealth),
		}, nil
	})
}
