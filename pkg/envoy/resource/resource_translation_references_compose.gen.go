package resource

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//

import (
	"github.com/cortezaproject/corteza-server/compose/types"
)

// ComposeChartResourceTranslationReferences generates Locale references
//
// This function is auto-generated
func ComposeChartResourceTranslationReferences(namespaceID string, self string) (res *Ref, pp []*Ref, err error) {
	res = &Ref{ResourceType: types.ChartResourceType, Identifiers: MakeIdentifiers(self)}
	pp = append(pp, &Ref{ResourceType: types.NamespaceResourceType, Identifiers: MakeIdentifiers(namespaceID)})

	return
}

// ComposeModuleResourceTranslationReferences generates Locale references
//
// This function is auto-generated
func ComposeModuleResourceTranslationReferences(namespaceID string, self string) (res *Ref, pp []*Ref, err error) {
	res = &Ref{ResourceType: types.ModuleResourceType, Identifiers: MakeIdentifiers(self)}
	pp = append(pp, &Ref{ResourceType: types.NamespaceResourceType, Identifiers: MakeIdentifiers(namespaceID)})

	return
}

// ComposeModuleFieldResourceTranslationReferences generates Locale references
//
// This function is auto-generated
func ComposeModuleFieldResourceTranslationReferences(namespaceID string, moduleID string, self string) (res *Ref, pp []*Ref, err error) {
	res = &Ref{ResourceType: types.ModuleFieldResourceType, Identifiers: MakeIdentifiers(self)}
	pp = append(pp, &Ref{ResourceType: types.NamespaceResourceType, Identifiers: MakeIdentifiers(namespaceID)})
	pp = append(pp, &Ref{ResourceType: types.ModuleResourceType, Identifiers: MakeIdentifiers(moduleID)})

	return
}

// ComposeNamespaceResourceTranslationReferences generates Locale references
//
// This function is auto-generated
func ComposeNamespaceResourceTranslationReferences(self string) (res *Ref, pp []*Ref, err error) {
	res = &Ref{ResourceType: types.NamespaceResourceType, Identifiers: MakeIdentifiers(self)}

	return
}

// ComposePageResourceTranslationReferences generates Locale references
//
// This function is auto-generated
func ComposePageResourceTranslationReferences(namespaceID string, self string) (res *Ref, pp []*Ref, err error) {
	res = &Ref{ResourceType: types.PageResourceType, Identifiers: MakeIdentifiers(self)}
	pp = append(pp, &Ref{ResourceType: types.NamespaceResourceType, Identifiers: MakeIdentifiers(namespaceID)})

	return
}
