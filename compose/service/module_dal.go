package service

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/dal"
	"github.com/cortezaproject/corteza-server/pkg/handle"
	systemTypes "github.com/cortezaproject/corteza-server/system/types"
)

type (
	dalDDL interface {
		ConnectionDefaults(ctx context.Context, connectionID uint64) (dft dal.ConnectionDefaults, err error)

		ReloadModel(ctx context.Context, models ...*dal.Model) (err error)
		AddModel(ctx context.Context, models ...*dal.Model) (err error)
		RemoveModel(ctx context.Context, models ...*dal.Model) (err error)
	}
)

const (
	// https://www.rfc-editor.org/errata/eid1690
	emailLength = 254

	// Generally the upper most limit
	urlLength = 2048

	sysID          = "ID"
	sysNamespaceID = "namespaceID"
	sysModuleID    = "moduleID"
	sysCreatedAt   = "createdAt"
	sysCreatedBy   = "createdBy"
	sysUpdatedAt   = "updatedAt"
	sysUpdatedBy   = "updatedBy"
	sysDeletedAt   = "deletedAt"
	sysDeletedBy   = "deletedBy"
	sysOwnedBy     = "ownedBy"

	colSysID          = "id"
	colSysNamespaceID = "rel_namespace"
	colSysModuleID    = "module_id"
	colSysCreatedAt   = "created_at"
	colSysCreatedBy   = "created_by"
	colSysUpdatedAt   = "updated_at"
	colSysUpdatedBy   = "updated_by"
	colSysDeletedAt   = "deleted_at"
	colSysDeletedBy   = "deleted_by"
	colSysOwnedBy     = "owned_by"
)

// ReloadDALModels reconstructs the DAL's data model based on the store.Storer
//
// Directly using store so we don't spam the action log
func (svc *module) ReloadDALModels(ctx context.Context) (err error) {
	var (
		namespaces types.NamespaceSet
		modules    types.ModuleSet
		fields     types.ModuleFieldSet
		models     dal.ModelSet
	)

	namespaces, _, err = svc.store.SearchComposeNamespaces(ctx, types.NamespaceFilter{})
	if err != nil {
		return
	}

	var model *dal.Model
	for _, ns := range namespaces {
		modules, _, err = svc.store.SearchComposeModules(ctx, types.ModuleFilter{
			NamespaceID: ns.ID,
		})
		if err != nil {
			return
		}

		for _, mod := range modules {
			fields, _, err = svc.store.SearchComposeModuleFields(ctx, types.ModuleFieldFilter{
				ModuleID: []uint64{mod.ID},
			})
			if err != nil {
				return
			}

			mod.Fields = append(mod.Fields, fields...)

			model, err = svc.moduleToModel(ctx, ns, mod)
			if err != nil {
				return
			}

			models = append(models, model)
		}
	}

	return svc.dal.ReloadModel(ctx, models...)
}

func (svc *module) moduleToModel(ctx context.Context, ns *types.Namespace, mod *types.Module) (*dal.Model, error) {
	ccfg, err := svc.dal.ConnectionDefaults(ctx, mod.ModelConfig.ConnectionID)
	if err != nil {
		return nil, err
	}
	getCodec := moduleFieldCodecBuilder(mod.ModelConfig.Partitioned, ccfg)

	// Metadata
	out := &dal.Model{
		ConnectionID: mod.ModelConfig.ConnectionID,
		Ident:        svc.formatPartitionIdent(ns, mod, ccfg),
		Label:        mod.Handle,

		Attributes: make(dal.AttributeSet, len(mod.Fields)),

		SensitivityLevel: mod.Privacy.SensitivityLevel,

		ResourceID:   mod.ID,
		ResourceType: types.ModuleResourceType,
		Resource:     mod.RbacResource(),
	}

	// Handle user-defined fields
	for i, f := range mod.Fields {
		out.Attributes[i], err = svc.moduleFieldToAttribute(getCodec, ns, mod, f)
		if err != nil {
			return nil, err
		}
	}

	// Handle system fields; either default or user defined
	if !mod.ModelConfig.Partitioned {
		// When not partitioned the default system fields should be defined along side the `values` column
		out.Attributes = append(out.Attributes, svc.moduleModelDefaultSysAttributes(getCodec)...)
	} else {
		// When partitioned, we use store codec defined on the module
		out.Attributes = append(out.Attributes, svc.moduleModelSysAttributes(mod, getCodec)...)
	}

	return out, nil
}

func (svc *module) moduleModelSysAttributes(mod *types.Module, getCodec func(f *types.ModuleField) dal.Codec) (out dal.AttributeSet) {
	sysEnc := mod.ModelConfig.SystemFieldEncoding

	if sysEnc.ID != nil {
		out = append(out, dal.PrimaryAttribute(sysID, getCodec(&types.ModuleField{Name: sysID, EncodingStrategy: *sysEnc.ID})))
	}

	if sysEnc.ModuleID != nil {
		out = append(out, dal.FullAttribute(sysModuleID, &dal.TypeID{}, getCodec(&types.ModuleField{Name: sysModuleID, EncodingStrategy: *sysEnc.ModuleID})))
	}
	if sysEnc.NamespaceID != nil {
		out = append(out, dal.FullAttribute(sysNamespaceID, &dal.TypeID{}, getCodec(&types.ModuleField{Name: sysNamespaceID, EncodingStrategy: *sysEnc.NamespaceID})))
	}

	if sysEnc.OwnedBy != nil {
		out = append(out, dal.FullAttribute(sysOwnedBy, &dal.TypeID{}, getCodec(&types.ModuleField{Name: sysOwnedBy, EncodingStrategy: *sysEnc.OwnedBy})))
	}

	if sysEnc.CreatedAt != nil {
		out = append(out, dal.FullAttribute(sysCreatedAt, &dal.TypeTimestamp{}, getCodec(&types.ModuleField{Name: sysCreatedAt, EncodingStrategy: *sysEnc.CreatedAt})))
	}
	if sysEnc.CreatedBy != nil {
		out = append(out, dal.FullAttribute(sysCreatedBy, &dal.TypeID{}, getCodec(&types.ModuleField{Name: sysCreatedBy, EncodingStrategy: *sysEnc.CreatedBy})))
	}

	if sysEnc.UpdatedAt != nil {
		out = append(out, dal.FullAttribute(sysUpdatedAt, &dal.TypeTimestamp{}, getCodec(&types.ModuleField{Name: sysUpdatedAt, EncodingStrategy: *sysEnc.UpdatedAt})))
	}
	if sysEnc.UpdatedBy != nil {
		out = append(out, dal.FullAttribute(sysUpdatedBy, &dal.TypeID{}, getCodec(&types.ModuleField{Name: sysUpdatedBy, EncodingStrategy: *sysEnc.UpdatedBy})))
	}

	if sysEnc.DeletedAt != nil {
		out = append(out, dal.FullAttribute(sysDeletedAt, &dal.TypeTimestamp{}, getCodec(&types.ModuleField{Name: sysDeletedAt, EncodingStrategy: *sysEnc.DeletedAt})))
	}
	if sysEnc.DeletedBy != nil {
		out = append(out, dal.FullAttribute(sysDeletedBy, &dal.TypeID{}, getCodec(&types.ModuleField{Name: sysDeletedBy, EncodingStrategy: *sysEnc.DeletedBy})))
	}

	return
}

func (svc *module) moduleModelDefaultSysAttributes(getCodec func(f *types.ModuleField) dal.Codec) dal.AttributeSet {
	return dal.AttributeSet{
		dal.PrimaryAttribute(sysID, &dal.CodecAlias{Ident: colSysID}),

		dal.FullAttribute(sysModuleID, &dal.TypeID{}, &dal.CodecAlias{Ident: colSysModuleID}),
		dal.FullAttribute(sysNamespaceID, &dal.TypeID{}, &dal.CodecAlias{Ident: colSysNamespaceID}),

		dal.FullAttribute(sysOwnedBy, &dal.TypeID{}, &dal.CodecAlias{Ident: colSysOwnedBy}),

		dal.FullAttribute(sysCreatedAt, &dal.TypeTimestamp{}, &dal.CodecAlias{Ident: colSysCreatedAt}),
		dal.FullAttribute(sysCreatedBy, &dal.TypeID{}, &dal.CodecAlias{Ident: colSysCreatedBy}),

		dal.FullAttribute(sysUpdatedAt, &dal.TypeTimestamp{}, &dal.CodecAlias{Ident: colSysUpdatedAt}),
		dal.FullAttribute(sysUpdatedBy, &dal.TypeID{}, &dal.CodecAlias{Ident: colSysUpdatedBy}),

		dal.FullAttribute(sysDeletedAt, &dal.TypeTimestamp{}, &dal.CodecAlias{Ident: colSysDeletedAt}),
		dal.FullAttribute(sysDeletedBy, &dal.TypeID{}, &dal.CodecAlias{Ident: colSysDeletedBy}),
	}
}

func (svc *module) moduleFieldToAttribute(getCodec func(f *types.ModuleField) dal.Codec, ns *types.Namespace, mod *types.Module, f *types.ModuleField) (out *dal.Attribute, err error) {
	kind := f.Kind
	if kind == "" {
		kind = "String"
	}

	switch strings.ToLower(kind) {
	case "bool":
		at := &dal.TypeBoolean{}
		out = dal.FullAttribute(f.Name, at, getCodec(f))
	case "datetime":
		switch {
		case f.IsDateOnly():
			at := &dal.TypeDate{}
			out = dal.FullAttribute(f.Name, at, getCodec(f))
		case f.IsTimeOnly():
			at := &dal.TypeTime{}
			out = dal.FullAttribute(f.Name, at, getCodec(f))
		default:
			at := &dal.TypeTimestamp{}
			out = dal.FullAttribute(f.Name, at, getCodec(f))
		}
	case "email":
		at := &dal.TypeText{Length: emailLength}
		out = dal.FullAttribute(f.Name, at, getCodec(f))
	case "file":
		at := &dal.TypeRef{
			RefModel:     &dal.Model{Resource: "corteza::system:attachment"},
			RefAttribute: &dal.Attribute{Ident: "id"},
		}
		out = dal.FullAttribute(f.Name, at, getCodec(f))
	case "number":
		at := &dal.TypeNumber{
			Precision: f.Options.Precision(),
		}
		out = dal.FullAttribute(f.Name, at, getCodec(f))
	case "record":
		at := &dal.TypeRef{
			RefModel: &dal.Model{
				ResourceID:   f.Options.UInt64("moduleID"),
				ResourceType: types.ModuleResourceType,
			},
			RefAttribute: &dal.Attribute{
				Ident: "id",
			},
		}
		out = dal.FullAttribute(f.Name, at, getCodec(f))
	case "select":
		at := &dal.TypeEnum{
			Values: f.SelectOptions(),
		}
		out = dal.FullAttribute(f.Name, at, getCodec(f))
	case "string":
		at := &dal.TypeText{
			Length: 0,
		}
		out = dal.FullAttribute(f.Name, at, getCodec(f))
	case "url":
		at := &dal.TypeText{
			Length: urlLength,
		}
		out = dal.FullAttribute(f.Name, at, getCodec(f))
	case "user":
		at := &dal.TypeRef{
			RefModel: &dal.Model{
				ResourceType: systemTypes.UserResourceType,
			},
			RefAttribute: &dal.Attribute{
				Ident: "id",
			},
		}
		out = dal.FullAttribute(f.Name, at, getCodec(f))

	default:
		return nil, fmt.Errorf("invalid field %s: kind %s not supported", f.Name, f.Kind)
	}

	out.SensitivityLevel = f.Privacy.SensitivityLevel
	out.Label = f.Name

	return
}

func (svc *module) formatPartitionIdent(ns *types.Namespace, mod *types.Module, cfg dal.ConnectionDefaults) string {
	if !mod.ModelConfig.Partitioned {
		return cfg.ModelIdent
	}

	pfmt := mod.ModelConfig.PartitionFormat
	if pfmt == "" {
		pfmt = cfg.PartitionFormat
	}
	if pfmt == "" {
		// @todo put in config or something
		pfmt = "compose_record_{{namespace}}_{{module}}"
	}

	// @note we must not use name here since it is translatable
	mh, _ := handle.Cast(nil, mod.Handle, strconv.FormatUint(mod.ID, 10))
	nsh, _ := handle.Cast(nil, ns.Slug, strconv.FormatUint(ns.ID, 10))
	rpl := strings.NewReplacer(
		"{{module}}", mh,
		"{{namespace}}", nsh,
	)

	return rpl.Replace(pfmt)
}

func moduleFieldCodecBuilder(partitioned bool, cfg dal.ConnectionDefaults) func(f *types.ModuleField) dal.Codec {
	return func(f *types.ModuleField) dal.Codec {
		return moduleFieldCodec(f, partitioned, cfg)
	}
}

func moduleFieldCodec(f *types.ModuleField, partitioned bool, cfg dal.ConnectionDefaults) (strat dal.Codec) {
	if partitioned {
		strat = &dal.CodecPlain{}
	} else {
		ident := cfg.AttributeIdent
		if ident == "" {
			// @todo put in configs or something
			ident = "values"
		}

		strat = &dal.CodecRecordValueSetJSON{
			Ident: ident,
		}
	}

	switch {
	case f.EncodingStrategy.EncodingStrategyAlias != nil:
		strat = &dal.CodecAlias{
			Ident: f.EncodingStrategy.EncodingStrategyAlias.Ident,
		}
	case f.EncodingStrategy.EncodingStrategyJSON != nil:
		strat = &dal.CodecRecordValueSetJSON{
			Ident: f.EncodingStrategy.EncodingStrategyJSON.Ident,
		}
	}

	return
}

// // // // // // // // // // // // // // // // // // // // // // // // //
// Utilities

func (svc *module) addModuleToDAL(ctx context.Context, ns *types.Namespace, mod *types.Module) (err error) {
	// Update DAL
	model, err := svc.moduleToModel(ctx, ns, mod)
	if err != nil {
		return
	}
	if err = svc.dal.AddModel(ctx, model); err != nil {
		return
	}

	return
}
