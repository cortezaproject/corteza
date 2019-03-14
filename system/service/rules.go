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
		Check(resource string, operation string, fallbacks ...rules.CheckAccessFunc) rules.Access
		Read(roleID uint64) (interface{}, error)
	}
)

var DefaultRules = service.DefaultRules

func Rules(ctx context.Context) RulesService {
	return DefaultRules.With(ctx)
}
