package service

import (
	"context"
	"testing"

	"github.com/crusttech/crust/system/internal/service"
	"github.com/crusttech/crust/system/types"
)

type (
	RoleService interface {
		FindByID(roleID uint64) (*types.Role, error)
		Find(filter *types.RoleFilter) ([]*types.Role, error)
	}
)

var DefaultRole = service.DefaultRole

func Role(ctx context.Context) RoleService {
	return DefaultRole.With(ctx)
}

// Expose the full Role API for testing
func TestRole(_ *testing.T, ctx context.Context) service.RoleService {
	return DefaultRole.With(ctx)
}
