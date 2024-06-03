package resources

import (
	"fmt"
	"lair-api/internal/server"

	"github.com/charmbracelet/log"
	"github.com/danielgtaylor/huma/v2"
)

type Resource interface {
	Init(api huma.API, s server.Server)
}

func Register(api huma.API, s server.Server, r Resource) {
	r.Init(api, s)
	log.Debug("Registered", "resource", fmt.Sprintf("%T", r))
}
