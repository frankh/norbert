package graph

import (
	"context"

	"github.com/frankh/norbert/cmd/norbert/config"
	"github.com/frankh/norbert/cmd/norbert/models"
)

type resolver struct {
}

func (r *resolver) RootQuery() RootQueryResolver {
	return r
}

func NewResolver() ResolverRoot {
	return &resolver{}
}

func (r *resolver) Service() ServiceResolver {
	return r
}

func (r *resolver) Services(ctx context.Context) ([]models.Service, error) {
	services := make([]models.Service, 0)

	for _, service := range config.Services {
		services = append(services, *service)
	}
	return services, nil
}

func (r *resolver) Checks(ctx context.Context, svc *models.Service) ([]models.Check, error) {
	return config.Checks[svc.Name], nil
}
