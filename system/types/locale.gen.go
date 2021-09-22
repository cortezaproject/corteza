package types

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//

// Definitions file that controls how this file is generated:
// - system.report.yaml

import (
	"fmt"
	"github.com/cortezaproject/corteza-server/pkg/locale"
	"strconv"
)

type (
	LocaleKey struct {
		Name          string
		Resource      string
		Path          string
		CustomHandler string
	}
)

// Types and stuff
const (
	ReportResourceTranslationType = "system:report"
)

var (
	LocaleKeyReportName = LocaleKey{
		Name:     "name",
		Resource: ReportResourceTranslationType,
		Path:     "name",
	}
	LocaleKeyReportDescription = LocaleKey{
		Name:     "description",
		Resource: ReportResourceTranslationType,
		Path:     "description",
	}
	LocaleKeyReportProjectionTitle = LocaleKey{
		Name:     "projection title",
		Resource: ReportResourceTranslationType,
		Path:     "projection.{{projectionID}}.title",
	}
	LocaleKeyReportProjectionDescription = LocaleKey{
		Name:     "projection description",
		Resource: ReportResourceTranslationType,
		Path:     "projection.{{projectionID}}.description",
	}
)

// ResourceTranslation returns string representation of Locale resource for Report by calling ReportResourceTranslation fn
//
// Locale resource is in the system:report/... format
//
// This function is auto-generated
func (r Report) ResourceTranslation() string {
	return ReportResourceTranslation(r.ID)
}

// ReportResourceTranslation returns string representation of Locale resource for Report
//
// Locale resource is in the system:report/... format
//
// This function is auto-generated
func ReportResourceTranslation(id uint64) string {
	cpts := []interface{}{ReportResourceTranslationType}
	cpts = append(cpts, strconv.FormatUint(id, 10))

	return fmt.Sprintf(ReportResourceTranslationTpl(), cpts...)
}

// @todo template
func ReportResourceTranslationTpl() string {
	return "%s/%s"
}

func (r *Report) DecodeTranslations(tt locale.ResourceTranslationIndex) {
	var aux *locale.ResourceTranslation
	if aux = tt.FindByKey(LocaleKeyReportName.Path); aux != nil {
		r.Meta.Name = aux.Msg
	}
	if aux = tt.FindByKey(LocaleKeyReportDescription.Path); aux != nil {
		r.Meta.Description = aux.Msg
	}

	r.decodeTranslations(tt)
}

func (r *Report) EncodeTranslations() (out locale.ResourceTranslationSet) {
	out = locale.ResourceTranslationSet{}
	if r.Meta.Name != "" {
		out = append(out, &locale.ResourceTranslation{
			Resource: r.ResourceTranslation(),
			Key:      LocaleKeyReportName.Path,
			Msg:      r.Meta.Name,
		})
	}
	if r.Meta.Description != "" {
		out = append(out, &locale.ResourceTranslation{
			Resource: r.ResourceTranslation(),
			Key:      LocaleKeyReportDescription.Path,
			Msg:      r.Meta.Description,
		})
	}

	out = append(out, r.encodeTranslations()...)

	return out
}
