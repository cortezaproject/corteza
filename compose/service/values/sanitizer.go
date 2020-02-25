package values

import (
	"github.com/cortezaproject/corteza-server/compose/types"
	"strconv"
	"strings"
	"time"
)

type (
	sanitizer struct{}
)

// Sanitizer initializes sanitizer
//
// Not really needed, following pattern in the package
func Sanitizer() *sanitizer {
	return &sanitizer{}
}

// Run cleans up input data
//  - fix multi-value order/place index
//  - trim all the strings!
//  - parse & format input values to match field specific -- nullify/falsify invalid
//  - field kind specific, no errors raised, data is modified
//
// Existing data (when updating record) is not yet loaded at this point
func (s sanitizer) Run(m *types.Module, vv types.RecordValueSet) (out types.RecordValueSet) {
	out = make([]*types.RecordValue, 0, len(vv))

	for _, f := range m.Fields {
		// Reorder and sanitize place value (no gaps)
		//
		// Values are ordered when received so we treat them like it
		// and assign the appropriate place no.
		for i, v := range vv.FilterByName(f.Name) {
			out = append(out, &types.RecordValue{
				Name:  f.Name,
				Value: v.Value,
				Ref:   v.Ref,
				Place: uint(i),
			})
		}
	}

	var (
		f    *types.ModuleField
		kind string
	)

	for _, v := range out {
		f = m.Fields.FindByName(v.Name)
		if f == nil {
			// Unknown field,
			// if it is not handled before,
			// sanitizer does not care about it
			continue
		}

		kind = strings.ToLower(f.Kind)

		if kind != "string" {
			// Trim all but string
			v.Value = strings.TrimSpace(v.Value)
		}

		if f.IsRef() && refy.MatchString(v.Value) {
			v.Ref, _ = strconv.ParseUint(v.Value, 10, 64)
		}

		// Per field type validators
		switch strings.ToLower(f.Kind) {
		case "bool":
			v = s.sBool(v, f, m)
		case "datetime":
			v = s.sDatetime(v, f, m)
			//case "email":
			//	v = s.sEmail(v, f, m)
			//case "file":
			//	v = s.sFile(v, f, m)
			//case "number":
			//	v = s.sNumber(v, f, m)
			//case "record":
			//	v = s.sRecord(v, f, m)
			//case "select":
			//	v = s.sSelect(v, f, m)
			//case "string":
			//	v = s.sString(v, f, m)
			//case "url":
			//	v = s.sUrl(v, f, m)
			//case "user":
			//	v = s.sUser(v, f, m)
		}
	}

	return
}

func (sanitizer) sBool(v *types.RecordValue, f *types.ModuleField, m *types.Module) *types.RecordValue {
	if truthy.MatchString(strings.ToLower(v.Value)) {
		v.Value = strBoolTrue
	} else {
		v.Value = strBoolFalse
	}

	return v
}

func (sanitizer) sDatetime(v *types.RecordValue, f *types.ModuleField, m *types.Module) *types.RecordValue {
	var (
		// input format set
		inputFormats []string

		// output format
		of string
	)

	if f.Options.Bool("onlyDate") {
		of = datetimeInputFormatDate
		inputFormats = []string{
			datetimeInputFormatDate,
			"02 Jan 06",
			"Monday, 02-Jan-06",
			"Mon, 02 Jan 2006",
			"2019/_1/_2",
		}
	} else if f.Options.Bool("onlyTime") {
		of = datetimeInputFormatTime
		inputFormats = []string{
			datetimeInputFormatTime,
			"15:04",
			"15:04:05Z07:00",
			"15:04:05 MST",
			"15:04:05 -0700",
			"15:04 MST",
			"15:04Z07:00",
			"15:04 -0700",
			time.Kitchen,
		}
	} else {
		of = datetimeInputFormatFull
		// date & time
		inputFormats = []string{
			datetimeInputFormatFull,
			time.RFC1123Z,
			time.RFC1123,
			time.RFC850,
			time.RFC822Z,
			time.RFC822,
			time.RubyDate,
			time.UnixDate,
			time.ANSIC,
			"2019/_1/_2 15:04:05",
			"2019/_1/_2 15:04",
		}
	}

	for _, format := range inputFormats {
		parsed, err := time.Parse(format, v.Value)

		if err == nil {
			v.Value = parsed.UTC().Format(of)
			break
		}
	}

	return v
}

//
//func (sanitizer) sEmail(v *types.RecordValue, f *types.ModuleField, m *types.Module) *types.RecordValue {
//  // @todo extract from "name" <email> format
//	return v
//}
//
//func (sanitizer) sFile(v *types.RecordValue, f *types.ModuleField, m *types.Module) *types.RecordValue {
//	return v
//}
//
//func (sanitizer) sNumber(v *types.RecordValue, f *types.ModuleField, m *types.Module) *types.RecordValue {
//	// @todo cut off the decimals / round up
//	// @todo sanitize decimal/thousands dot/comma
//	return v
//}
//
//func (sanitizer) sRecord(v *types.RecordValue, f *types.ModuleField, m *types.Module) *types.RecordValue {
//	return v
//}
//
//func (sanitizer) sSelect(v *types.RecordValue, f *types.ModuleField, m *types.Module) *types.RecordValue {
//	return v
//}
//
//func (sanitizer) sString(v *types.RecordValue, f *types.ModuleField, m *types.Module) *types.RecordValue {
//	return v
//}
//
//func (sanitizer) sUrl(v *types.RecordValue, f *types.ModuleField, m *types.Module) *types.RecordValue {
//	return v
//}
//
//func (sanitizer) sUser(v *types.RecordValue, f *types.ModuleField, m *types.Module) *types.RecordValue {
//	return v
//}
