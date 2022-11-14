package plugin

import (
	"context"

	"github.com/cortezaproject/corteza-server/automation/types"
	"go.uber.org/zap"
)

type (
	Setup interface {
		Setup(log *zap.Logger) error
	}

	Initialize interface {
		Initialize(ctx context.Context, log *zap.Logger) error
	}

	AutomationFunctionsProvider interface {
		AutomationFunctions() []*types.Function
	}

	//AutomationTypesProvider interface {
	//	AutomationTypes() []*expr.Type
	//}
)
