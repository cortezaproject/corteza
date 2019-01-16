package rules

import (
	"context"

	"github.com/titpetric/factory"
)

type ResourcesInterface interface {
	With(ctx context.Context, db *factory.DB) ResourcesInterface

	CheckAccessMulti(resource string, operation string) error
	CheckAccess(resource string, operation string) error

	Grant(resource string, teamID uint64, operations []string, value Access) error
}
