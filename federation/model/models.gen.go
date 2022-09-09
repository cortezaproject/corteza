package model

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//

import (
	"github.com/cortezaproject/corteza-server/federation/types"
	"github.com/cortezaproject/corteza-server/pkg/dal"
)

var ExposedModule = &dal.Model{
	Ident:        "federation_module_exposed",
	ResourceType: types.ExposedModuleResourceType,

	Attributes: dal.AttributeSet{
		&dal.Attribute{
			Ident: "ID",
			Type:  &dal.TypeID{},
			Store: &dal.CodecAlias{Ident: "id"},
		},

		&dal.Attribute{
			Ident: "Handle",
			Type:  &dal.TypeText{Length: 64},
			Store: &dal.CodecAlias{Ident: "handle"},
		},

		&dal.Attribute{
			Ident: "Name", Sortable: true,
			Type:  &dal.TypeText{},
			Store: &dal.CodecAlias{Ident: "name"},
		},

		&dal.Attribute{
			Ident: "NodeID", Sortable: true,
			Type: &dal.TypeRef{HasDefault: true,
				DefaultValue: 0,

				RefAttribute: "id",
				RefModel: &dal.ModelRef{
					ResourceType: "corteza::federation:node",
				},
			},
			Store: &dal.CodecAlias{Ident: "rel_node"},
		},

		&dal.Attribute{
			Ident: "ComposeModuleID",
			Type:  &dal.TypeID{},
			Store: &dal.CodecAlias{Ident: "rel_compose_module"},
		},

		&dal.Attribute{
			Ident: "ComposeNamespaceID",
			Type:  &dal.TypeID{},
			Store: &dal.CodecAlias{Ident: "rel_compose_namespace"},
		},

		&dal.Attribute{
			Ident: "Fields",
			Type: &dal.TypeJSON{
				DefaultValue: "{}",
			},
			Store: &dal.CodecAlias{Ident: "fields"},
		},

		&dal.Attribute{
			Ident: "CreatedAt", Sortable: true,
			Type: &dal.TypeTimestamp{
				DefaultCurrentTimestamp: true, Timezone: true, Precision: -1,
			},
			Store: &dal.CodecAlias{Ident: "created_at"},
		},

		&dal.Attribute{
			Ident: "UpdatedAt", Sortable: true,
			Type:  &dal.TypeTimestamp{Nullable: true, Timezone: true, Precision: -1},
			Store: &dal.CodecAlias{Ident: "updated_at"},
		},

		&dal.Attribute{
			Ident: "DeletedAt", Sortable: true,
			Type:  &dal.TypeTimestamp{Nullable: true, Timezone: true, Precision: -1},
			Store: &dal.CodecAlias{Ident: "deleted_at"},
		},

		&dal.Attribute{
			Ident: "CreatedBy",
			Type: &dal.TypeRef{HasDefault: true,
				DefaultValue: 0,

				RefAttribute: "id",
				RefModel: &dal.ModelRef{
					ResourceType: "corteza::system:user",
				},
			},
			Store: &dal.CodecAlias{Ident: "created_by"},
		},

		&dal.Attribute{
			Ident: "UpdatedBy",
			Type: &dal.TypeRef{HasDefault: true,
				DefaultValue: 0,

				RefAttribute: "id",
				RefModel: &dal.ModelRef{
					ResourceType: "corteza::system:user",
				},
			},
			Store: &dal.CodecAlias{Ident: "updated_by"},
		},

		&dal.Attribute{
			Ident: "DeletedBy",
			Type: &dal.TypeRef{HasDefault: true,
				DefaultValue: 0,

				RefAttribute: "id",
				RefModel: &dal.ModelRef{
					ResourceType: "corteza::system:user",
				},
			},
			Store: &dal.CodecAlias{Ident: "deleted_by"},
		},
	},

	Indexes: dal.IndexSet{
		&dal.Index{
			Ident: "PRIMARY",
			Type:  "BTREE",

			Fields: []*dal.IndexField{
				{
					AttributeIdent: "ID",
				},
			},
		},
	},
}

var ModuleMapping = &dal.Model{
	Ident:        "federation_module_mapping",
	ResourceType: types.ModuleMappingResourceType,

	Attributes: dal.AttributeSet{
		&dal.Attribute{
			Ident: "FederationModuleID", Sortable: true,
			Type:  &dal.TypeID{},
			Store: &dal.CodecAlias{Ident: "federation_module_id"},
		},

		&dal.Attribute{
			Ident: "ComposeModuleID", Sortable: true,
			Type:  &dal.TypeID{},
			Store: &dal.CodecAlias{Ident: "compose_module_id"},
		},

		&dal.Attribute{
			Ident: "ComposeNamespaceID", Sortable: true,
			Type:  &dal.TypeID{},
			Store: &dal.CodecAlias{Ident: "compose_namespace_id"},
		},

		&dal.Attribute{
			Ident: "FieldMapping",
			Type: &dal.TypeJSON{
				DefaultValue: "{}",
			},
			Store: &dal.CodecAlias{Ident: "field_mapping"},
		},
	},

	Indexes: dal.IndexSet{
		&dal.Index{
			Ident: "uniqueModuleComposeModule",
			Type:  "BTREE",

			Fields: []*dal.IndexField{
				{
					AttributeIdent: "FederationModuleID",
				},

				{
					AttributeIdent: "ComposeModuleID",
				},

				{
					AttributeIdent: "ComposeNamespaceID",
				},
			},
		},
	},
}

var Node = &dal.Model{
	Ident:        "federation_nodes",
	ResourceType: types.NodeResourceType,

	Attributes: dal.AttributeSet{
		&dal.Attribute{
			Ident: "ID",
			Type:  &dal.TypeID{},
			Store: &dal.CodecAlias{Ident: "id"},
		},

		&dal.Attribute{
			Ident: "SharedNodeID", Sortable: true,
			Type:  &dal.TypeID{},
			Store: &dal.CodecAlias{Ident: "shared_node_id"},
		},

		&dal.Attribute{
			Ident: "Name", Sortable: true,
			Type:  &dal.TypeText{},
			Store: &dal.CodecAlias{Ident: "name"},
		},

		&dal.Attribute{
			Ident: "BaseURL", Sortable: true,
			Type:  &dal.TypeText{},
			Store: &dal.CodecAlias{Ident: "base_url"},
		},

		&dal.Attribute{
			Ident: "Status", Sortable: true,
			Type:  &dal.TypeText{},
			Store: &dal.CodecAlias{Ident: "status"},
		},

		&dal.Attribute{
			Ident: "Contact", Sortable: true,
			Type:  &dal.TypeText{},
			Store: &dal.CodecAlias{Ident: "contact"},
		},

		&dal.Attribute{
			Ident: "PairToken",
			Type:  &dal.TypeText{},
			Store: &dal.CodecAlias{Ident: "pair_token"},
		},

		&dal.Attribute{
			Ident: "AuthToken",
			Type:  &dal.TypeText{},
			Store: &dal.CodecAlias{Ident: "auth_token"},
		},

		&dal.Attribute{
			Ident: "CreatedAt", Sortable: true,
			Type: &dal.TypeTimestamp{
				DefaultCurrentTimestamp: true, Timezone: true, Precision: -1,
			},
			Store: &dal.CodecAlias{Ident: "created_at"},
		},

		&dal.Attribute{
			Ident: "UpdatedAt", Sortable: true,
			Type:  &dal.TypeTimestamp{Nullable: true, Timezone: true, Precision: -1},
			Store: &dal.CodecAlias{Ident: "updated_at"},
		},

		&dal.Attribute{
			Ident: "DeletedAt", Sortable: true,
			Type:  &dal.TypeTimestamp{Nullable: true, Timezone: true, Precision: -1},
			Store: &dal.CodecAlias{Ident: "deleted_at"},
		},

		&dal.Attribute{
			Ident: "CreatedBy",
			Type: &dal.TypeRef{HasDefault: true,
				DefaultValue: 0,

				RefAttribute: "id",
				RefModel: &dal.ModelRef{
					ResourceType: "corteza::system:user",
				},
			},
			Store: &dal.CodecAlias{Ident: "created_by"},
		},

		&dal.Attribute{
			Ident: "UpdatedBy",
			Type: &dal.TypeRef{HasDefault: true,
				DefaultValue: 0,

				RefAttribute: "id",
				RefModel: &dal.ModelRef{
					ResourceType: "corteza::system:user",
				},
			},
			Store: &dal.CodecAlias{Ident: "updated_by"},
		},

		&dal.Attribute{
			Ident: "DeletedBy",
			Type: &dal.TypeRef{HasDefault: true,
				DefaultValue: 0,

				RefAttribute: "id",
				RefModel: &dal.ModelRef{
					ResourceType: "corteza::system:user",
				},
			},
			Store: &dal.CodecAlias{Ident: "deleted_by"},
		},
	},

	Indexes: dal.IndexSet{
		&dal.Index{
			Ident: "PRIMARY",
			Type:  "BTREE",

			Fields: []*dal.IndexField{
				{
					AttributeIdent: "ID",
				},
			},
		},
	},
}

var NodeSync = &dal.Model{
	Ident:        "federation_nodes_sync",
	ResourceType: types.NodeSyncResourceType,

	Attributes: dal.AttributeSet{
		&dal.Attribute{
			Ident: "NodeID", Sortable: true,
			Type:  &dal.TypeID{},
			Store: &dal.CodecAlias{Ident: "node_id"},
		},

		&dal.Attribute{
			Ident: "ModuleID", Sortable: true,
			Type:  &dal.TypeID{},
			Store: &dal.CodecAlias{Ident: "module_id"},
		},

		&dal.Attribute{
			Ident: "SyncType", Sortable: true,
			Type:  &dal.TypeText{},
			Store: &dal.CodecAlias{Ident: "sync_type"},
		},

		&dal.Attribute{
			Ident: "SyncStatus", Sortable: true,
			Type:  &dal.TypeText{},
			Store: &dal.CodecAlias{Ident: "sync_status"},
		},

		&dal.Attribute{
			Ident: "TimeOfAction", Sortable: true,
			Type:  &dal.TypeTimestamp{Timezone: true, Precision: -1},
			Store: &dal.CodecAlias{Ident: "time_of_action"},
		},
	},

	Indexes: dal.IndexSet{},
}

var SharedModule = &dal.Model{
	Ident:        "federation_module_shared",
	ResourceType: types.SharedModuleResourceType,

	Attributes: dal.AttributeSet{
		&dal.Attribute{
			Ident: "ID",
			Type:  &dal.TypeID{},
			Store: &dal.CodecAlias{Ident: "id"},
		},

		&dal.Attribute{
			Ident: "Handle",
			Type:  &dal.TypeText{Length: 64},
			Store: &dal.CodecAlias{Ident: "handle"},
		},

		&dal.Attribute{
			Ident: "NodeID", Sortable: true,
			Type:  &dal.TypeID{},
			Store: &dal.CodecAlias{Ident: "rel_node"},
		},

		&dal.Attribute{
			Ident: "Name", Sortable: true,
			Type:  &dal.TypeText{},
			Store: &dal.CodecAlias{Ident: "name"},
		},

		&dal.Attribute{
			Ident: "ExternalFederationModuleID", Sortable: true,
			Type:  &dal.TypeID{},
			Store: &dal.CodecAlias{Ident: "xref_module"},
		},

		&dal.Attribute{
			Ident: "Fields",
			Type: &dal.TypeJSON{
				DefaultValue: "{}",
			},
			Store: &dal.CodecAlias{Ident: "fields"},
		},

		&dal.Attribute{
			Ident: "CreatedAt", Sortable: true,
			Type: &dal.TypeTimestamp{
				DefaultCurrentTimestamp: true, Timezone: true, Precision: -1,
			},
			Store: &dal.CodecAlias{Ident: "created_at"},
		},

		&dal.Attribute{
			Ident: "UpdatedAt", Sortable: true,
			Type:  &dal.TypeTimestamp{Nullable: true, Timezone: true, Precision: -1},
			Store: &dal.CodecAlias{Ident: "updated_at"},
		},

		&dal.Attribute{
			Ident: "DeletedAt", Sortable: true,
			Type:  &dal.TypeTimestamp{Nullable: true, Timezone: true, Precision: -1},
			Store: &dal.CodecAlias{Ident: "deleted_at"},
		},

		&dal.Attribute{
			Ident: "CreatedBy",
			Type: &dal.TypeRef{HasDefault: true,
				DefaultValue: 0,

				RefAttribute: "id",
				RefModel: &dal.ModelRef{
					ResourceType: "corteza::system:user",
				},
			},
			Store: &dal.CodecAlias{Ident: "created_by"},
		},

		&dal.Attribute{
			Ident: "UpdatedBy",
			Type: &dal.TypeRef{HasDefault: true,
				DefaultValue: 0,

				RefAttribute: "id",
				RefModel: &dal.ModelRef{
					ResourceType: "corteza::system:user",
				},
			},
			Store: &dal.CodecAlias{Ident: "updated_by"},
		},

		&dal.Attribute{
			Ident: "DeletedBy",
			Type: &dal.TypeRef{HasDefault: true,
				DefaultValue: 0,

				RefAttribute: "id",
				RefModel: &dal.ModelRef{
					ResourceType: "corteza::system:user",
				},
			},
			Store: &dal.CodecAlias{Ident: "deleted_by"},
		},
	},

	Indexes: dal.IndexSet{
		&dal.Index{
			Ident: "PRIMARY",
			Type:  "BTREE",

			Fields: []*dal.IndexField{
				{
					AttributeIdent: "ID",
				},
			},
		},
	},
}

func init() {
	models = append(
		models,
		ExposedModule,
		ModuleMapping,
		Node,
		NodeSync,
		SharedModule,
	)
}
