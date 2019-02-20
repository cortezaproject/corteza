package rules

import (
	"context"

	"github.com/titpetric/factory"
)

type ResourcesInterface interface {
	With(ctx context.Context, db *factory.DB) ResourcesInterface

	IsAllowed(resource string, operation string) Access

	GrantByResource(roleID uint64, resource string, operations []string, value Access) error
	ListByResource(roleID uint64, resource string) ([]Rule, error)

	Grant(roleID uint64, rules []Rule) error
	List(roleID uint64) ([]Rule, error)
	Delete(roleID uint64) error
}
