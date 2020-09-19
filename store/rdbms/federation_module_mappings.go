package rdbms

import (
	"github.com/Masterminds/squirrel"
	"github.com/cortezaproject/corteza-server/federation/types"
)

func (s Store) convertFederationModuleMappingFilter(f types.ModuleMappingFilter) (query squirrel.SelectBuilder, err error) {
	query = s.federationModuleMappingsSelectBuilder()

	if f.ComposeModuleID > 0 {
		query = query.Where("cmd.compose_module_id = ?", f.ComposeModuleID)
	}

	if f.FederationModuleID > 0 {
		query = query.Where("cmd.federation_module_id = ?", f.FederationModuleID)
	}

	return
}
