package graph

import (
	"context"

	"github.com/frankh/norbert/cmd/norbert/config"
	"github.com/frankh/norbert/cmd/norbert/models"
	"github.com/frankh/norbert/cmd/norbert/repository"
	"github.com/frankh/norbert/pkg/check"
	"github.com/nats-io/go-nats"
)

type resolver struct {
	db repository.Repository
	nc *nats.Conn
}

func (r *resolver) RootQuery() RootQueryResolver {
	return r
}

func NewResolver(db repository.Repository, nc *nats.Conn) ResolverRoot {
	return &resolver{db, nc}
}

func (r *resolver) Service() ServiceResolver {
	return r
}

func (r *resolver) Check() CheckResolver {
	return r
}

func (r *resolver) Subscription() SubscriptionResolver {
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

	for _, c := range config.Checks[svc.Name] {
		checks = append(checks, *c)
	}
	return checks, nil
}

func (r *resolver) Results(ctx context.Context, c *models.Check) ([]*models.CheckResult, error) {
	return r.db.CheckResults(c.Id())
}

func (r *resolver) Status(ctx context.Context, c *models.Check) (models.CheckStatus, error) {
	results, err := r.db.CheckResults(c.Id())
	if err != nil {
		return 0, nil
	}

	if len(results) == 0 {
		return models.Ok, nil
	}

	if results[0].ResultCode != check.Success {
		return models.Failed, nil
	}

	return models.Ok, nil
}
