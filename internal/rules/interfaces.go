package rules

import (
	"context"

	"github.com/titpetric/factory"
)

type ResourcesInterface interface {
	With(ctx context.Context, db *factory.DB) ResourcesInterface

	IsAllowed(resource string, operation string) Access

	Grant(teamID uint64, resource string, operations []string, value Access) error
	ListGrants(teamID uint64, resource string) ([]Rules, error)
}
