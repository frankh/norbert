package graph

import (
	"context"

	"github.com/frankh/norbert/cmd/norbert/config"
	"github.com/frankh/norbert/cmd/norbert/models"
	"github.com/frankh/norbert/cmd/norbert/repository"
)

type resolver struct {
	db repository.Repository
}

func (r *resolver) RootQuery() RootQueryResolver {
	return r
}

func NewResolver(db repository.Repository) ResolverRoot {
	return &resolver{db}
}

func (r *resolver) Service() ServiceResolver {
	return r
}

func (r *resolver) Check() CheckResolver {
	return r
}

func (r *resolver) Services(ctx context.Context) ([]models.Service, error) {
	services := make([]models.Service, 0)

	for _, service := range config.Services {
		services = append(services, *service)
	}
	return services, nil
}

func (r *resolver) GetCheck(ctx context.Context, checkId string) (*models.Check, error) {
	return config.ChecksById[checkId], nil
}

func (r *resolver) Checks(ctx context.Context, svc *models.Service) ([]models.Check, error) {
	checks := make([]models.Check, 0)

	for _, check := range config.Checks[svc.Name] {
		checks = append(checks, *check)
	}
	return checks, nil
}

func (r *resolver) Results(ctx context.Context, check *models.Check) ([]*models.CheckResult, error) {
	return r.db.CheckResults(check.Id())
}
