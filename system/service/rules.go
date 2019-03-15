package service

import (
	"context"

	"github.com/crusttech/crust/internal/rules"
	"github.com/crusttech/crust/system/internal/service"
)

type (
	RulesService interface {
		List() (interface{}, error)
		Effective(filter string) ([]service.EffectiveRules, error)
		Check(resource rules.Resource, operation string, fallbacks ...rules.CheckAccessFunc) rules.Access
		Read(roleID uint64) (interface{}, error)
	}
)

var DefaultRules = service.DefaultRules

func Rules(ctx context.Context) service.RulesService {
	return DefaultRules.With(ctx)
}
