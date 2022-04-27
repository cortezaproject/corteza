package resource

import (
	"fmt"
	"strconv"

	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/minions"
	systemTypes "github.com/cortezaproject/corteza-server/system/types"
	"github.com/spf13/cast"
)

type (
	ComposeModuleField struct {
		*base
		Res *types.ModuleField

		RefNs  *Ref
		RefMod *Ref
	}

	ComposeModule struct {
		*base
		Res       *types.Module
		ResFields []*ComposeModuleField

		// Might keep track of related NS
		RefNs    *Ref
		RefMods  RefSet
		RefRoles RefSet
	}
)

func NewComposeModule(res *types.Module, nsRef *Ref) *ComposeModule {
	r := &ComposeModule{
		base:    &base{},
		RefMods: make(RefSet, 0, len(res.Fields)),
	}
	r.SetResourceType(types.ModuleResourceType)
	r.Res = res

	r.AddIdentifier(identifiers(res.Handle, res.Name, res.ID)...)

	r.RefNs = r.addRef(nsRef)

	// Field deps
	for _, f := range res.Fields {
		switch f.Kind {
		case "Record":
			refMod := f.Options.String("module")
			if refMod == "" {
				refMod = f.Options.String("moduleID")
			}
			if refMod != "" && refMod != "0" {
				r.RefMods = append(r.RefMods, r.AddRef(types.ModuleResourceType, refMod).Constraint(r.RefNs))
			}

		case "User":
			refRoles := ComposeModuleFieldExtractUserFieldRoles(f.Options["roles"])
			if len(refRoles) == 0 {
				refRoles = ComposeModuleFieldExtractUserFieldRoles(f.Options["role"])
			}
			if len(refRoles) == 0 {
				refRoles = ComposeModuleFieldExtractUserFieldRoles(f.Options["roleID"])
			}

			for _, refRole := range refRoles {
				r.RefRoles = append(r.RefRoles, r.AddRef(systemTypes.RoleResourceType, refRole))
			}
		}
	}

	// Initial timestamps
	r.SetTimestamps(MakeTimestampsCUDA(&res.CreatedAt, res.UpdatedAt, res.DeletedAt, nil))

	return r
}

func (r *ComposeModule) Resource() interface{} {
	return r.Res
}

func (r *ComposeModule) ReRef(old RefSet, new RefSet) {
	r.base.ReRef(old, new)

	// Additional references
	for _, n := range new {
		if n.ResourceType == types.NamespaceResourceType {
			r.RefNs = MakeRef(types.NamespaceResourceType, n.Identifiers)
		}
	}

	for i, o := range old {
		if o.ResourceType == types.ModuleResourceType {
			r.RefMods = r.RefMods.replaceRef(o, new[i])
		}
	}
}

func (r *ComposeModule) SysID() uint64 {
	return r.Res.ID
}

func (r *ComposeModule) resourceParts(tpl string) (resource string, ref *Ref, path []*Ref) {
	ref = r.Ref().Constraint(r.RefNs)
	path = []*Ref{r.RefNs}
	resource = fmt.Sprintf(tpl, types.ModuleResourceType, r.RefNs.Identifiers.First(), firstOkString(strconv.FormatUint(r.Res.ID, 10), r.Res.Handle))

	return
}

func (r *ComposeModule) RBACParts() (resource string, ref *Ref, path []*Ref) {
	return r.resourceParts(types.ModuleRbacResourceTpl())
}

func (r *ComposeModule) ResourceTranslationParts() (resource string, ref *Ref, path []*Ref) {
	return r.resourceParts(types.ModuleResourceTranslationTpl())
}

func (r *ComposeModule) encodeTranslations() ([]*ResourceTranslation, error) {
	out := make([]*ResourceTranslation, 0, len(r.ResFields))

	for _, f := range r.ResFields {
		rr := f.Res.EncodeTranslations()
		rr.SetLanguage(defaultLanguage)
		res, ref, pp := f.ResourceTranslationParts()
		out = append(out, NewResourceTranslation(systemTypes.FromLocale(rr), res, ref, pp...))
	}

	return out, nil
}

// FindComposeModule looks for the module in the resource set
func FindComposeModule(rr InterfaceSet, ii Identifiers) (ns *types.Module) {
	var modRes *ComposeModule

	rr.Walk(func(r Interface) error {
		mr, ok := r.(*ComposeModule)
		if !ok {
			return nil
		}

		if mr.Identifiers().HasAny(ii) {
			modRes = mr
		}
		return nil
	})

	// Found it
	if modRes != nil {
		return modRes.Res
	}
	return nil
}

func (r *ComposeModule) AddField(f *ComposeModuleField) {
	r.ResFields = append(r.ResFields, f)
}

// FindComposeModule looks for the module in the resource set
func FindComposeModuleField(rr InterfaceSet, mod, ii Identifiers) (field *types.ModuleField) {
	var fieldRes *ComposeModuleField

	rr.Walk(func(r Interface) error {
		mr, ok := r.(*ComposeModule)
		if !ok {
			return nil
		}

		if !mr.Identifiers().HasAny(mod) {
			return nil
		}

		for _, rf := range mr.ResFields {
			if rf.Identifiers().HasAny(ii) {
				fieldRes = rf
			}
		}
		return nil
	})

	// Found it
	if fieldRes != nil {
		return fieldRes.Res
	}
	return nil
}

func (r *ComposeModuleField) resourceParts(tpl string) (resource string, ref *Ref, path []*Ref) {
	ref = r.Ref().Constraint(r.RefNs).Constraint(r.RefMod)
	path = []*Ref{r.RefNs, r.RefMod}
	resource = fmt.Sprintf(types.ModuleFieldResourceTranslationTpl(), types.ModuleFieldResourceTranslationType, r.RefNs.Identifiers.First(), r.RefMod.Identifiers.First(), firstOkString(strconv.FormatUint(r.Res.ID, 10), r.Res.Name))

	return
}

func (r *ComposeModuleField) RBACParts() (resource string, ref *Ref, path []*Ref) {
	return r.resourceParts(types.ModuleFieldRbacResourceTpl())
}

func (r *ComposeModuleField) ResourceTranslationParts() (resource string, ref *Ref, path []*Ref) {
	return r.resourceParts(types.ModuleFieldResourceTranslationTpl())
}

func (r *ComposeModuleField) ResourceTranslations() ([]*ResourceTranslation, error) {
	out := make([]*ResourceTranslation, 0, 10)

	rr := r.Res.EncodeTranslations()
	rr.SetLanguage(defaultLanguage)
	res, ref, pp := r.ResourceTranslationParts()
	out = append(out, NewResourceTranslation(systemTypes.FromLocale(rr), res, ref, pp...))

	return out, nil
}

func ComposeModuleErrUnresolved(ii Identifiers) error {
	return fmt.Errorf("compose module unresolved %v", ii.StringSlice())
}

func ComposeModuleFieldErrUnresolved(ii Identifiers) error {
	return fmt.Errorf("compose module field unresolved %v", ii.StringSlice())
}

func NewComposeModuleField(res *types.ModuleField, nsRef, modRef *Ref) *ComposeModuleField {
	r := &ComposeModuleField{
		base: &base{},
	}
	r.SetResourceType(types.ModuleFieldResourceType)
	r.Res = res

	r.AddIdentifier(identifiers(res.Name, res.ID)...)

	r.RefNs = r.addRef(nsRef)
	r.RefMod = r.addRef(modRef).Constraint(r.RefNs)

	// Initial timestamps
	r.SetTimestamps(MakeTimestampsCUDA(&res.CreatedAt, res.UpdatedAt, res.DeletedAt, nil))

	return r
}

func (r *ComposeModuleField) Resource() interface{} {
	return r.Res
}

// ComposeModuleFieldExtractUserFieldRoles is a helper to extract roles
// from the given filer options.
func ComposeModuleFieldExtractUserFieldRoles(i interface{}) []string {
	if minions.IsNil(i) {
		return nil
	}

	out := make([]string, 0, 1)

	isOk := func(v string) bool {
		return v != "" && v != "0"
	}

	switch v := i.(type) {
	case uint64:
		aux := strconv.FormatUint(v, 10)
		if !isOk(aux) {
			return nil
		}
		return []string{aux}
	case []uint64:
		for _, i := range v {
			aux := strconv.FormatUint(i, 10)
			if !isOk(aux) {
				continue
			}
			out = append(out, aux)
		}
		return out

	case string:
		if !isOk(v) {
			return nil
		}
		return []string{v}
	case []string:
		for _, aux := range v {
			if !isOk(aux) {
				continue
			}
			out = append(out, aux)
		}
		return out

	case []interface{}:
		for _, i := range v {
			aux := cast.ToString(i)
			if !isOk(aux) {
				continue
			}
			out = append(out, aux)
		}
		return out
	case interface{}:
		aux := cast.ToString(v)
		if !isOk(aux) {
			return nil
		}
		return []string{aux}
	}

	return nil
}
