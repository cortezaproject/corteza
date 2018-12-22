package service

import (
	"encoding/json"

	"github.com/crusttech/crust/crm/types"
)

func (r *content) preloadAll(contents []*types.Record, fields ...string) error {
	for _, content := range contents {
		if err := r.preload(content, fields...); err != nil {
			return err
		}
	}
	return nil
}

func (r *content) preload(content *types.Record, fields ...string) (err error) {
	for _, field := range fields {
		switch field {
		case "fields":
			fields, err := r.Fields(content)
			if err != nil {
				return err
			}
			json, err := json.Marshal(fields)
			if err != nil {
				return err
			}
			if err := (&content.Fields).Scan(json); err != nil {
				return err
			}
		case "page":
			if content.Page, err = r.pageRepo.FindByModuleID(content.ModuleID); err != nil {
				return
			}
		case "user":
			if content.UserID > 0 {
				if content.User, err = r.userSvc.FindByID(content.UserID); err != nil {
					return
				}
			}
		}
	}
	return
}
