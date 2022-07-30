package rdbms

import (
	"context"
	. "github.com/cortezaproject/corteza-server/store/adapters/rdbms/ddl"
)

func fix202209_extendComposeModuleForPrivacyAndDAL(ctx context.Context, s *Store) (err error) {
	s.log(ctx).Info("extending compose_module table with config column")
	return s.SchemaAPI.AddColumn(
		ctx, s.DB,
		&Table{Name: "compose_module"},
		&Column{Type: ColumnType{Type: ColumnTypeJson}, DefaultValue: "'{}'", Name: "config"},
	)
}

func fix202209_extendComposeModuleFieldsForPrivacyAndDAL(ctx context.Context, s *Store) (err error) {
	s.log(ctx).Info("extending compose_module_field table with privacy and encoding_strategy columns")
	return s.SchemaAPI.AddColumn(
		ctx, s.DB,
		&Table{Name: "compose_module_field"},
		&Column{Type: ColumnType{Type: ColumnTypeJson}, DefaultValue: "'{}'", Name: "privacy"},
		&Column{Type: ColumnType{Type: ColumnTypeJson}, DefaultValue: "'{}'", Name: "encoding_strategy"},
	)
}
