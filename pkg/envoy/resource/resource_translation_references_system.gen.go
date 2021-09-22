package resource

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//

// Definitions file that controls how this file is generated:
// - system.report.yaml

import (
	"github.com/cortezaproject/corteza-server/system/types"
)

// SystemReportResourceTranslationReferences generates Locale references
//
// Resources with "envoy: false" are skipped
//
// This function is auto-generated
func SystemReportResourceTranslationReferences(report string) (res *Ref, pp []*Ref, err error) {
	res = &Ref{ResourceType: types.ReportResourceType, Identifiers: MakeIdentifiers(report)}

	return
}
