package service

import (
	"github.com/crusttech/crust/crm/types"
)

func (s *page) preloadAll(pages types.PageSet) (err error) {
	var modules types.ModuleSet
	modules, err = s.moduleRepo.Find()
	if err != nil {
		return err
	}
	return pages.Walk(func(page *types.Page) error {
		page.Module = modules.FindByID(page.ModuleID)
		return nil
	})
}

func (s *page) preload(page *types.Page) (err error) {
	if page.ModuleID > 0 {
		page.Module, err = s.moduleRepo.FindByID(page.ModuleID)
		return
	}
	return
}
