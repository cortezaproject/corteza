package rules

import (
	"context"

	"github.com/titpetric/factory"
)

type ResourcesInterface interface {
	With(ctx context.Context, db *factory.DB) ResourcesInterface

	IsAllowed(resource string, operation string) Access

	Grant(resource string, teamID uint64, operations []string, value Access) error
	ListGrants(resource string, teamID uint64) ([]Rules, error)
}
