package rdbms

import (
	"context"
)

func (s Store) SearchRoleMembers(ctx context.Context, roleID uint64) ([]uint64, error) {
	return nil, nil
}

func (s Store) SearchUserMemberships(ctx context.Context, userID uint64) ([]uint64, error) {
	return nil, nil
}
