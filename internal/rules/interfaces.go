package rules

import (
	"context"

	"github.com/titpetric/factory"
)

type ResourcesInterface interface {
	With(ctx context.Context, db *factory.DB) ResourcesInterface

	Check(resource string, operation string, fallbacks ...CheckAccessFunc) Access

	Grant(roleID uint64, rules []Rule) error
	Read(roleID uint64) ([]Rule, error)
	Delete(roleID uint64) error
}
