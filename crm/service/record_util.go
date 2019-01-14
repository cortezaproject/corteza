package service

import (
	"github.com/crusttech/crust/crm/types"
)

func (s *record) preloadAll(module *types.Module, records []*types.Record, fields ...string) (err error) {
	if len(records) == 0 {
		return nil
	}

	if module == nil {
		if module, err = s.moduleRepo.FindByID(records[0].ID); err != nil {
			// Assuming all records are from the same module
			return
		}
	}

	for _, record := range records {
		if err = s.preload(module, record, fields...); err != nil {
			return
		}
	}
	return
}

func (s *record) preload(module *types.Module, record *types.Record, fields ...string) (err error) {
	if module == nil {
		if module, err = s.moduleRepo.FindByID(record.ModuleID); err != nil {
			return err
		}
	}

	for _, field := range fields {
		switch field {
		case "fields":
			// fields, err := s.Fields(module, record)
			// if err != nil {
			// 	return err
			// }
			// json, err := json.Marshal(fields)
			// if err != nil {
			// 	return err
			// }
			// if err := (&record.Values).Scan(json); err != nil {
			// 	return err
			// }
		case "page":
			if record.Page, err = s.pageRepo.FindByModuleID(record.ModuleID); err != nil {
				return
			}
		case "user":
			if record.UserID > 0 {
				if record.User, err = s.userSvc.FindByID(record.UserID); err != nil {
					return
				}
			}
		}
	}
	return
}
