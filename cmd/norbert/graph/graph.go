package graph

import (
	"context"
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
	return nil, nil
}

func (r *resolver) Checks(ctx context.Context, svc *models.Service) ([]models.Check, error) {
	return nil, nil
}
