package service

import (
	"github.com/crusttech/crust/system/internal/service"
)

func Init() error {
	err := service.Init()
	DefaultRole = service.DefaultRole
	DefaultRules = service.DefaultRules
	DefaultUser = service.DefaultUser
	return err
}
