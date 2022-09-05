package rdbms

import (
	"context"
	"fmt"
	automationModels "github.com/cortezaproject/corteza-server/automation/model"
	composeModels "github.com/cortezaproject/corteza-server/compose/model"
	federationModels "github.com/cortezaproject/corteza-server/federation/model"
	"github.com/cortezaproject/corteza-server/pkg/dal"
	"github.com/cortezaproject/corteza-server/store/adapters/rdbms/ddl"
	systemModels "github.com/cortezaproject/corteza-server/system/model"
	"go.uber.org/zap"
)

func (s *Store) Upgrade(ctx context.Context) (err error) {

	err = createTablesFromModels(
		ctx,
		s.log(ctx),
		s.DataDefiner,
		systemModels.All(),
		composeModels.All(),
		automationModels.All(),
		federationModels.All(),
	)

	if err != nil {
		return err
	}

	//var (
	//	tableExists bool
	//)
	//
	//for _, t := range Tables() {
	//	tableExists, err = s.SchemaAPI.TableExists(ctx, s.DB, t.Name)
	//	if err != nil {
	//		return fmt.Errorf("could not check table %q existance: %w", t.Name, err)
	//	}
	//
	//	if !tableExists {
	//		s.log(ctx).Debug("creating table", zap.String("table", t.Name))
	//		if err = s.SchemaAPI.CreateTable(ctx, s.DB, t); err != nil {
	//			return fmt.Errorf("could not create table %q: %w", t.Name, err)
	//		}
	//	}
	//}
	//
	//fixes := []func(context.Context, *Store) error{
	//	fix202209_extendComposeModuleForPrivacyAndDAL,
	//	fix202209_extendComposeModuleFieldsForPrivacyAndDAL,
	//	fix202209_dropObsoleteComposeModuleFields,
	//	fix202209_composeRecordRevisions,
	//	fix202209_extendDalConnectionsForMeta,
	//	fix202209_renameModuleColOnComposeRecords,
	//	fix202209_addMetaOnComposeRecords,
	//}
	//
	//for _, fix := range fixes {
	//	if err = fix(ctx, s); err != nil {
	//		return
	//	}
	//}
	//
	return
}

func createTablesFromModels(ctx context.Context, log *zap.Logger, dd ddl.DataDefiner, sets ...dal.ModelSet) (err error) {
	var (
		tbl *ddl.Table
	)

	for _, mm := range sets {
		for _, m := range mm {
			log.Debug("creating table", zap.String("table", m.Ident))

			if tbl, err = dd.ConvertModel(m); err != nil {
				return fmt.Errorf("can not convert model %q to table: %w", m.Ident, err)
			}

			if err = dd.TableCreate(ctx, tbl); err != nil {
				return fmt.Errorf("can not create table from mdoel model %q: %w", m.Ident, err)
			}
		}
	}

	return nil
}
