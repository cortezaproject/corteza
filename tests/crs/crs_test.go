package crs

import (
	"context"
	"testing"
	"time"

	"github.com/cortezaproject/corteza-server/compose/crs"
	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/id"

	"github.com/stretchr/testify/require"

	"github.com/cortezaproject/corteza-server/compose/crs/capabilities"
	_ "github.com/cortezaproject/corteza-server/store/adapters/rdbms/drivers/mysql"
)

/*
See here for my table definitions

CREATE TABLE `the_cake` (
  `id` bigint unsigned NOT NULL,
  `name` VARCHAR(45),
  `want` BOOLEAN,
  `ownedBy` bigint unsigned NOT NULL,
  `createdAt` datetime NOT NULL,
  `updatedAt` datetime DEFAULT NULL,
  `deletedAt` datetime DEFAULT NULL,
  `createdBy` bigint unsigned NOT NULL,
  `updatedBy` bigint unsigned NOT NULL DEFAULT '0',
  `deletedBy` bigint unsigned NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`),
  KEY `compose_record_owner` (`ownedBy`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb3;

CREATE TABLE `the_cookie` (
  `id` bigint unsigned NOT NULL,
  `name` VARCHAR(45),
  `is_good` BOOLEAN,
  `ownedBy` bigint unsigned NOT NULL,
  `createdAt` datetime NOT NULL,
  `updatedAt` datetime DEFAULT NULL,
  `deletedAt` datetime DEFAULT NULL,
  `createdBy` bigint unsigned NOT NULL,
  `updatedBy` bigint unsigned NOT NULL DEFAULT '0',
  `deletedBy` bigint unsigned NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`),
  KEY `compose_record_owner` (`ownedBy`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb3;
*/

func TestHello(t *testing.T) {
	ctx := context.Background()

	c, err := crs.ComposeRecordStore(
		ctx,
		nil,
		false,
		// Primary...
		crs.CRSConnectionWrap(0, "mysql://envoy:envoy@tcp(localhost:3306)/crs?collation=utf8mb4_general_ci", capabilities.FullCapabilities()...),

		// Others...
		crs.CRSConnectionWrap(1, "mysql://envoy:envoy@tcp(localhost:3306)/crs?collation=utf8mb4_general_ci", capabilities.FullCapabilities()...),
	)
	require.NoError(t, err)
	_ = c

	// ---

	cookieModule := &types.Module{
		ID:     10001,
		Handle: "cookie",
		Store: types.CRSDef{
			ComposeRecordStoreID: 0,
			Capabilities:         capabilities.FullCapabilities(),
			Partitioned:          true,
			PartitionFormat:      "the_{{module}}",
		},
		Fields: types.ModuleFieldSet{&types.ModuleField{
			ModuleID: 10001,
			ID:       20001,
			Name:     "name",
			Encoding: types.EncodingStrategy{
				EncodingStrategyAlias: &types.EncodingStrategyAlias{
					Ident: "name",
				},
			},
		}, &types.ModuleField{
			ModuleID: 10001,
			ID:       20002,
			Name:     "is_good",
			Kind:     "Bool",
			Encoding: types.EncodingStrategy{
				EncodingStrategyAlias: &types.EncodingStrategyAlias{
					Ident: "is_good",
				},
			},
		}},
	}

	cakeModule := &types.Module{
		ID:     10002,
		Handle: "cake",
		Store: types.CRSDef{
			ComposeRecordStoreID: 1,
			Capabilities:         capabilities.FullCapabilities(),
			Partitioned:          true,
			PartitionFormat:      "the_{{module}}",
		},
		Fields: types.ModuleFieldSet{&types.ModuleField{
			ModuleID: 10002,
			ID:       20003,
			Name:     "name",
			Encoding: types.EncodingStrategy{
				EncodingStrategyAlias: &types.EncodingStrategyAlias{
					// this doesn't work
					Ident: "name",
				},
			},
		}, &types.ModuleField{
			ModuleID: 10002,
			ID:       20004,
			Name:     "want",
			Kind:     "Bool",
			Encoding: types.EncodingStrategy{
				EncodingStrategyAlias: &types.EncodingStrategyAlias{
					Ident: "want",
				},
			},
		}},
	}

	err = c.ReloadModules(ctx, cookieModule, cakeModule)
	require.NoError(t, err)

	// ---

	a := time.Now()

	cookieRecord := &types.Record{
		ID:       id.Next(),
		ModuleID: 10001,
		Values: types.RecordValueSet{{
			Name:  "name",
			Value: "SOME COOKIE HERE",
		}, {
			Name:  "is_good",
			Value: "1",
		}},
		CreatedAt: time.Now(),
		UpdatedAt: &a,
		DeletedAt: &a,
		OwnedBy:   20003,
		CreatedBy: 20003,
		UpdatedBy: 20003,
		DeletedBy: 20003,
	}

	cakeRecord := &types.Record{
		ID:       id.Next(),
		ModuleID: 10002,
		Values: types.RecordValueSet{{
			Name:  "name",
			Value: "DOME CAKE HERE",
		}, {
			Name:  "want",
			Value: "1",
		}},
		CreatedAt: time.Now(),
		UpdatedAt: &a,
		DeletedAt: &a,
		OwnedBy:   20003,
		CreatedBy: 20003,
		UpdatedBy: 20003,
		DeletedBy: 20003,
	}

	err = c.ComposeRecordCreate(ctx, cookieModule, cookieRecord)
	require.NoError(t, err)

	err = c.ComposeRecordCreate(ctx, cakeModule, cakeRecord)
	require.NoError(t, err)

	require.FailNow(t, "the test dies a failNow so I can see logs!!!")
}
