package dalutils

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/dal"
	"github.com/cortezaproject/corteza-server/pkg/handle"
	"github.com/cortezaproject/corteza-server/store"
	systemTypes "github.com/cortezaproject/corteza-server/system/types"
)

type (
	identFormatter interface {
		ModelIdentFormatter(connectionID uint64) (f *dal.IdentFormatter, err error)
	}

	modelReloader interface {
		identFormatter
		ReloadModel(ctx context.Context, models ...*dal.Model) (err error)
	}

	modelCreator interface {
		identFormatter
		CreateModel(ctx context.Context, models ...*dal.Model) (err error)
	}

	modelUpdater interface {
		identFormatter
		UpdateModel(ctx context.Context, old, new *dal.Model) (err error)
	}

	attributeUpdater interface {
		identFormatter
		UpdateModelAttribute(ctx context.Context, model *dal.Model, old, new *dal.Attribute, trans ...dal.TransformationFunction) (err error)
	}

	modelDeleter interface {
		identFormatter
		DeleteModel(ctx context.Context, models ...*dal.Model) (err error)
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

func ComposeModulesReload(ctx context.Context, s store.Storer, r modelReloader) (err error) {
	var (
		namespaces types.NamespaceSet
		modules    types.ModuleSet
		fields     types.ModuleFieldSet
		models     dal.ModelSet
	)

	namespaces, _, err = s.SearchComposeNamespaces(ctx, types.NamespaceFilter{})
	if err != nil {
		return
	}

	var model *dal.Model
	for _, ns := range namespaces {
		modules, _, err = s.SearchComposeModules(ctx, types.ModuleFilter{
			NamespaceID: ns.ID,
		})
		if err != nil {
			return
		}

		for _, mod := range modules {
			fields, _, err = s.SearchComposeModuleFields(ctx, types.ModuleFieldFilter{
				ModuleID: []uint64{mod.ID},
			})
			if err != nil {
				return
			}

			mod.Fields = append(mod.Fields, fields...)

			model, err = ModuleToModel(ctx, r, ns, mod)
			if err != nil {
				return
			}

			models = append(models, model)
		}
	}

	return r.ReloadModel(ctx, models...)
}

func ComposeModuleCreate(ctx context.Context, c modelCreator, ns *types.Namespace, mod *types.Module) (err error) {
	model, err := ModuleToModel(ctx, c, ns, mod)
	if err != nil {
		return
	}
	if err = c.CreateModel(ctx, model); err != nil {
		return
	}

	return
}

func ComposeModuleUpdate(ctx context.Context, u modelUpdater, ns *types.Namespace, old, new *types.Module) (err error) {
	oldModel, err := ModuleToModel(ctx, u, ns, old)
	if err != nil {
		return
	}
	newModel, err := ModuleToModel(ctx, u, ns, new)
	if err != nil {
		return
	}

	if err = u.UpdateModel(ctx, oldModel, newModel); err != nil {
		return
	}

	return
}

func ComposeModuleFieldsUpdate(ctx context.Context, u attributeUpdater, ns *types.Namespace, old, new *types.Module) (err error) {
	oldModel, err := ModuleToModel(ctx, u, ns, old)
	if err != nil {
		return
	}
	newModel, err := ModuleToModel(ctx, u, ns, new)
	if err != nil {
		return
	}

	diff := oldModel.Diff(newModel)
	for _, d := range diff {
		if err = u.UpdateModelAttribute(ctx, oldModel, d.Original, d.Asserted); err != nil {
			return
		}
	}

	return
}

func ComposeModuleDelete(ctx context.Context, d modelDeleter, ns *types.Namespace, mod *types.Module) (err error) {
	model, err := ModuleToModel(ctx, d, ns, mod)
	if err != nil {
		return
	}
	if err = d.DeleteModel(ctx, model); err != nil {
		return
	}

	return
}

func ComposeModuleModelFormatter(f identFormatter, ns *types.Namespace, mod *types.Module) (formatter *dal.IdentFormatter, tplParts []string, err error) {
	formatter, err = f.ModelIdentFormatter(mod.ModelConfig.ConnectionID)
	if err != nil {
		return
	}

	modHandle, _ := handle.Cast(nil, mod.Handle, strconv.FormatUint(mod.ID, 10))
	nsHandle, _ := handle.Cast(nil, ns.Slug, strconv.FormatUint(ns.ID, 10))
	tplParts = []string{
		"module", modHandle,
		"namespace", nsHandle,
	}

	return
}

// // // // // // // // // // // // // // // // // // // // // // // // //
// Utilities

func ModuleToModel(ctx context.Context, f identFormatter, ns *types.Namespace, mod *types.Module) (*dal.Model, error) {
	formatter, tplParts, err := ComposeModuleModelFormatter(f, ns, mod)
	if err != nil {
		return nil, err
	}

	// Metadata
	out := &dal.Model{
		ConnectionID: mod.ModelConfig.ConnectionID,
		Label:        mod.Handle,

		Attributes: make(dal.AttributeSet, len(mod.Fields)),

		SensitivityLevel: mod.Privacy.SensitivityLevel,

		ResourceID:   mod.ID,
		ResourceType: types.ModuleResourceType,
		Resource:     mod.RbacResource(),
	}

	var ok bool
	out.Ident, ok = formatter.ModelIdent(ctx, mod.ModelConfig.Partitioned, mod.ModelConfig.PartitionFormat, tplParts...)
	if !ok {
		return nil, fmt.Errorf("invalid model identifier generated: %s", out.Ident)
	}

	getCodec := moduleFieldCodecBuilder(mod.ModelConfig.Partitioned, formatter)

	// Handle user-defined fields
	for i, f := range mod.Fields {
		out.Attributes[i], err = moduleFieldToAttribute(getCodec, ns, mod, f)
		if err != nil {
			return nil, err
		}
	}

	// Handle system fields; either default or user defined
	if !mod.ModelConfig.Partitioned {
		// When not partitioned the default system fields should be defined along side the `values` column
		out.Attributes = append(out.Attributes, moduleModelDefaultSysAttributes()...)
	} else {
		// When partitioned, we use store codec defined on the module
		out.Attributes = append(out.Attributes, moduleModelSysAttributes(mod, getCodec)...)
	}

	return out, nil
}

func moduleModelSysAttributes(mod *types.Module, getCodec func(f *types.ModuleField) dal.Codec) (out dal.AttributeSet) {
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
		out = append(out, dal.FullAttribute(sysUpdatedAt, &dal.TypeTimestamp{Nullable: true}, getCodec(&types.ModuleField{Name: sysUpdatedAt, EncodingStrategy: *sysEnc.UpdatedAt})))
	}
	if sysEnc.UpdatedBy != nil {
		out = append(out, dal.FullAttribute(sysUpdatedBy, &dal.TypeID{Nullable: true}, getCodec(&types.ModuleField{Name: sysUpdatedBy, EncodingStrategy: *sysEnc.UpdatedBy})))
	}

	if sysEnc.DeletedAt != nil {
		out = append(out, dal.FullAttribute(sysDeletedAt, &dal.TypeTimestamp{Nullable: true}, getCodec(&types.ModuleField{Name: sysDeletedAt, EncodingStrategy: *sysEnc.DeletedAt})))
	}
	if sysEnc.DeletedBy != nil {
		out = append(out, dal.FullAttribute(sysDeletedBy, &dal.TypeID{Nullable: true}, getCodec(&types.ModuleField{Name: sysDeletedBy, EncodingStrategy: *sysEnc.DeletedBy})))
	}

	return
}

func moduleModelDefaultSysAttributes() dal.AttributeSet {
	return dal.AttributeSet{
		dal.PrimaryAttribute(sysID, &dal.CodecAlias{Ident: colSysID}),

		dal.FullAttribute(sysModuleID, &dal.TypeID{}, &dal.CodecAlias{Ident: colSysModuleID}),
		dal.FullAttribute(sysNamespaceID, &dal.TypeID{}, &dal.CodecAlias{Ident: colSysNamespaceID}),

		dal.FullAttribute(sysOwnedBy, &dal.TypeID{}, &dal.CodecAlias{Ident: colSysOwnedBy}),

		dal.FullAttribute(sysCreatedAt, &dal.TypeTimestamp{}, &dal.CodecAlias{Ident: colSysCreatedAt}),
		dal.FullAttribute(sysCreatedBy, &dal.TypeID{}, &dal.CodecAlias{Ident: colSysCreatedBy}),

		dal.FullAttribute(sysUpdatedAt, &dal.TypeTimestamp{Nullable: true}, &dal.CodecAlias{Ident: colSysUpdatedAt}),
		dal.FullAttribute(sysUpdatedBy, &dal.TypeID{Nullable: true}, &dal.CodecAlias{Ident: colSysUpdatedBy}),

		dal.FullAttribute(sysDeletedAt, &dal.TypeTimestamp{Nullable: true}, &dal.CodecAlias{Ident: colSysDeletedAt}),
		dal.FullAttribute(sysDeletedBy, &dal.TypeID{Nullable: true}, &dal.CodecAlias{Ident: colSysDeletedBy}),
	}
}

func moduleFieldToAttribute(getCodec func(f *types.ModuleField) dal.Codec, ns *types.Namespace, mod *types.Module, f *types.ModuleField) (out *dal.Attribute, err error) {
	kind := f.Kind
	if kind == "" {
		kind = "String"
	}

	switch strings.ToLower(kind) {
	case "bool", "boolean":
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

func moduleFieldCodecBuilder(partitioned bool, formatter *dal.IdentFormatter) func(f *types.ModuleField) dal.Codec {
	return func(f *types.ModuleField) dal.Codec {
		return moduleFieldCodec(f, partitioned, formatter)
	}
}

func moduleFieldCodec(f *types.ModuleField, partitioned bool, formatter *dal.IdentFormatter) (strat dal.Codec) {
	if partitioned {
		strat = &dal.CodecPlain{}
	} else {
		ident, ok := formatter.AttributeIdent(partitioned, f.Name)
		if !ok {
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
