package handler

import (
	"guthub.com/gopeshwark/go-boilerplate/internal/server"
	"guthub.com/gopeshwark/go-boilerplate/internal/service"
)

type Handlers struct {
	Health  *HealthHandler
	OpenAPI *OpenAPIHandler
}

func NewHandlers(s *server.Server, services *service.Services) *Handlers {
	return &Handlers{
		Health:  NewHealthHandler(s),
		OpenAPI: NewOpenAPIHandler(s),
	}
}
