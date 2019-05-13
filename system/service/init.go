package service

import (
	"context"

	"github.com/crusttech/crust/system/internal/service"
)

func Init(ctx context.Context) error {
	err := service.Init(ctx)
	DefaultUser = service.DefaultUser
	return err
}

func Watchers(ctx context.Context) {
	service.Watchers(ctx)
}
