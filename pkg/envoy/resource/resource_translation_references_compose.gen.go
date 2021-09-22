package resource

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//

// Definitions file that controls how this file is generated:
// - compose.module-field.yaml
// - compose.module.yaml
// - compose.namespace.yaml
// - compose.page.yaml

import (
	"github.com/cortezaproject/corteza-server/compose/types"
)

// ComposeModuleFieldResourceTranslationReferences generates Locale references
//
// Resources with "envoy: false" are skipped
//
// This function is auto-generated
func ComposeModuleFieldResourceTranslationReferences(namespace string, module string, moduleField string) (res *Ref, pp []*Ref, err error) {
	pp = append(pp, &Ref{ResourceType: types.NamespaceResourceType, Identifiers: MakeIdentifiers(namespace)})
	pp = append(pp, &Ref{ResourceType: types.ModuleResourceType, Identifiers: MakeIdentifiers(module)})
	res = &Ref{ResourceType: types.ModuleFieldResourceType, Identifiers: MakeIdentifiers(moduleField)}

	return
}

// ComposeModuleResourceTranslationReferences generates Locale references
//
// Resources with "envoy: false" are skipped
//
// This function is auto-generated
func ComposeModuleResourceTranslationReferences(namespace string, module string) (res *Ref, pp []*Ref, err error) {
	pp = append(pp, &Ref{ResourceType: types.NamespaceResourceType, Identifiers: MakeIdentifiers(namespace)})
	res = &Ref{ResourceType: types.ModuleResourceType, Identifiers: MakeIdentifiers(module)}

	return
}

// ComposeNamespaceResourceTranslationReferences generates Locale references
//
// Resources with "envoy: false" are skipped
//
// This function is auto-generated
func ComposeNamespaceResourceTranslationReferences(namespace string) (res *Ref, pp []*Ref, err error) {
	res = &Ref{ResourceType: types.NamespaceResourceType, Identifiers: MakeIdentifiers(namespace)}

	return
}

// ComposePageResourceTranslationReferences generates Locale references
//
// Resources with "envoy: false" are skipped
//
// This function is auto-generated
func ComposePageResourceTranslationReferences(namespace string, page string) (res *Ref, pp []*Ref, err error) {
	pp = append(pp, &Ref{ResourceType: types.NamespaceResourceType, Identifiers: MakeIdentifiers(namespace)})
	res = &Ref{ResourceType: types.PageResourceType, Identifiers: MakeIdentifiers(page)}

	return
}
