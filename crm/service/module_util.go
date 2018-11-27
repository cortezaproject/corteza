package service

import (
	"github.com/crusttech/crust/crm/types"
)

func (r *module) preload(mod *types.Module) (err error) {
	mod.Page, err = r.pageRepo.FindByModuleID(mod.ID)
	return
}
