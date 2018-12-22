package service

import (
	"encoding/json"

	"github.com/crusttech/crust/crm/types"
)

func (r *record) preloadAll(records []*types.Record, fields ...string) error {
	for _, record := range records {
		if err := r.preload(record, fields...); err != nil {
			return err
		}
	}
	return nil
}

func (r *record) preload(record *types.Record, fields ...string) (err error) {
	for _, field := range fields {
		switch field {
		case "fields":
			fields, err := r.Fields(record)
			if err != nil {
				return err
			}
			json, err := json.Marshal(fields)
			if err != nil {
				return err
			}
			if err := (&record.Fields).Scan(json); err != nil {
				return err
			}
		case "page":
			if record.Page, err = r.pageRepo.FindByModuleID(record.ModuleID); err != nil {
				return
			}
		case "user":
			if record.UserID > 0 {
				if record.User, err = r.userSvc.FindByID(record.UserID); err != nil {
					return
				}
			}
		}
	}
	return
}
