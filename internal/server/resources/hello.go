package resources

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"

	"lair-api/internal/server"
)

type HelloResource struct{}

type GreatingResponse struct {
	Body struct {
		Message string `json:"name"`
	}
}

func (l *HelloResource) Init(api huma.API, s server.Server) {
	huma.Register(api, huma.Operation{
		OperationID: "get-hello",
		Method:      http.MethodGet,
		Path:        "/hello",
		Summary:     "Greating",
		Tags:        []string{"Miscellaneous"},
	}, func(ctx context.Context, input *struct{}) (*GreatingResponse, error) {
		resp := &GreatingResponse{}
		resp.Body.Message = "Hello, World!"
		return resp, nil
	})
}
