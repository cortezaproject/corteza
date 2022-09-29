package rdbms

import (
	"context"
	"encoding/json"
	"github.com/cortezaproject/corteza-server/pkg/dal"
	labelsType "github.com/cortezaproject/corteza-server/pkg/label/types"
	"github.com/cortezaproject/corteza-server/store"
	"go.uber.org/zap"
)

// RDBMS database fixes
//
// Schema changes that can not be automatically applied or complex changes
// that require some logic to be applied are handled here.
//
// Function names should start with "fix" and version.
// This does not have any effect on how fixes are executed, only for organisation purposes.

var (
	// all enabled fix function need to be listed here
	fixes = []func(context.Context, *Store) error{
		fix_2022_09_00_extendComposeModuleForPrivacyAndDAL,
		fix_2022_09_00_extendComposeModuleFieldsForPrivacyAndDAL,
		fix_2022_09_00_dropObsoleteComposeModuleFields,
		fix_2022_09_00_extendDalConnectionsForMeta,
		fix_2022_09_00_renameModuleColOnComposeRecords,
		fix_2022_09_00_addMetaOnComposeRecords,
	}
)

func fix_2022_09_00_extendComposeModuleForPrivacyAndDAL(ctx context.Context, s *Store) (err error) {
	return addColumn(ctx, s,
		"compose_module",
		&dal.Attribute{Ident: "config", Type: &dal.TypeJSON{DefaultValue: "{}"}},
	)
}

func fix_2022_09_00_extendComposeModuleFieldsForPrivacyAndDAL(ctx context.Context, s *Store) (err error) {
	return addColumn(ctx, s,
		"compose_module",
		&dal.Attribute{Ident: "config", Type: &dal.TypeJSON{DefaultValue: "{}"}},
	)
}

func fix_2022_09_00_dropObsoleteComposeModuleFields(ctx context.Context, s *Store) (err error) {
	return dropColumns(ctx, s,
		"compose_module_field",
		"is_private",
		"is_visible",
	)
}

func fix_2022_09_00_extendDalConnectionsForMeta(ctx context.Context, s *Store) (err error) {
	return addColumn(ctx, s,
		"dal_connections",
		&dal.Attribute{Ident: "meta", Type: &dal.TypeJSON{DefaultValue: "{}"}},
	)
}

func fix_2022_09_00_renameModuleColOnComposeRecords(ctx context.Context, s *Store) (err error) {
	return renameColumn(ctx, s, "compose_record", "module_id", "rel_module")
}

func fix_2022_09_00_addMetaOnComposeRecords(ctx context.Context, s *Store) (err error) {
	var (
		log = s.log(ctx)

		groupedMeta = make(map[uint64]map[string]any)
		packed      []byte
	)

	err = addColumn(ctx, s,
		"compose_record",
		&dal.Attribute{Ident: "meta", Type: &dal.TypeJSON{DefaultValue: "{}"}},
	)

	if err != nil {
		return
	}

	return s.Tx(ctx, func(ctx context.Context, s store.Storer) (err error) {
		log.Info("collecting record labels")
		ll, _, err := store.SearchLabels(ctx, s, labelsType.LabelFilter{Kind: "compose:record"})
		if err != nil {
			return
		}

		log.Info("grouping labels", zap.Int("count", len(ll)))
		for _, l := range ll {
			if _, has := groupedMeta[l.ResourceID]; !has {
				groupedMeta[l.ResourceID] = make(map[string]any)
			}

			groupedMeta[l.ResourceID][l.Name] = l.Value
			if err = store.DeleteLabel(ctx, s, l); err != nil {
				return
			}
		}

		log.Info("updating records with meta", zap.Int("count", len(ll)))
		for recordID, labels := range groupedMeta {
			packed, err = json.Marshal(labels)
			_, err = s.(*Store).DB.ExecContext(ctx, "UPDATE compose_record SET meta = $1 WHERE id = $2", packed, recordID)
			if err != nil {
				return
			}
		}

		return

	})

}
