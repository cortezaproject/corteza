package service

import (
	"github.com/crusttech/crust/system/internal/service"
)

func Init() {
	service.Init()
	DefaultRules = service.DefaultRules
	DefaultUser = service.DefaultUser
}
