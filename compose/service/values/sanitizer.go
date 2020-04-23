package values

import (
	"strconv"
	"strings"
	"time"

	"github.com/cortezaproject/corteza-server/compose/types"
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
		var i = 0
		for _, v := range vv.FilterByName(f.Name) {
			if v.IsDeleted() {
				continue
			}

			c := v.Clone()
			c.Place = uint(i)
			out = append(out, c)
			i++
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

		if v.IsDeleted() || !v.Updated {
			// Ignore unchanged and deleted
			continue
		}

		kind = strings.ToLower(f.Kind)

		if kind != "string" {
			// Trim all but string
			v.Value = strings.TrimSpace(v.Value)
		}

		if f.IsRef() {
			if refy.MatchString(v.Value) {
				v.Ref, _ = strconv.ParseUint(v.Value, 10, 64)
			}

			if v.Ref == 0 {
				v.Value = ""
			}
		}

		// Per field type validators
		switch strings.ToLower(f.Kind) {
		case "bool":
			v = s.sBool(v, f, m)
		case "datetime":
			v = s.sDatetime(v, f, m)
		case "number":
			v = s.sNumber(v, f, m)

			// Uncomment when they become relevant for sanitization
			//case "email":
			//	v = s.sEmail(v, f, m)
			//case "file":
			//	v = s.sFile(v, f, m)
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
		internalFormat string
	)

	if f.Options.Bool("onlyDate") {
		internalFormat = datetimeInternalFormatDate
		inputFormats = []string{
			datetimeInternalFormatDate,
			"02 Jan 06",
			"Monday, 02-Jan-06",
			"Mon, 02 Jan 2006",
			"2006/_1/_2",
		}
	} else if f.Options.Bool("onlyTime") {
		internalFormat = datetimeIntenralFormatTime
		inputFormats = []string{
			datetimeIntenralFormatTime,
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
		internalFormat = datetimeInternalFormatFull
		// date & time
		inputFormats = []string{
			datetimeInternalFormatFull,
			time.RFC1123Z,
			time.RFC1123,
			time.RFC850,
			time.RFC822Z,
			time.RFC822,
			time.RubyDate,
			time.UnixDate,
			time.ANSIC,
			"2006/_1/_2 15:04:05",
			"2006/_1/_2 15:04",
		}

		// if string looks like a RFC 3330 (ISO 8601), see if we need to suffix it with Z
		if isoDaty.MatchString(v.Value) && !hasTimezone.MatchString(v.Value) {
			// No timezone, add Z to satisfy parser
			v.Value = v.Value + "Z"

			// Simplifiy list of rules
			inputFormats = []string{time.RFC3339}
		}
	}

	for _, format := range inputFormats {
		parsed, err := time.Parse(format, v.Value)
		if err == nil {
			v.Value = parsed.UTC().Format(internalFormat)
			return v
		}
	}

	v.Value = ""
	return v
}

func (sanitizer) sNumber(v *types.RecordValue, f *types.ModuleField, m *types.Module) *types.RecordValue {
	// No point in continuing
	if v.Value == "" || f.Options == nil {
		return v
	}

	// Default to 0 for consistency
	if f.Options["precision"] == nil || f.Options["precision"] == "" {
		f.Options["precision"] = 0
	}

	// Since Options are not structured, there would appear that there can be a bit of a mess
	// when it comes to types,so this is needed.
	var prec float64
	unk := f.Options["precision"]
	switch i := unk.(type) {
	case float64:
		prec = i
	case int:
		prec = float64(i)
	case int64:
		prec = float64(i)
	case string:
		pp, err := strconv.ParseFloat(i, 64)
		if err != nil {
			prec = 0
			break
		}
		prec = pp
	}

	// Clamp between 0 and 6; this was originally done in corteza-js so we keep it here.
	if prec < 0 {
		prec = 0
	}
	if prec > 6 {
		prec = 6
	}

	base, err := strconv.ParseFloat(v.Value, 64)
	if err != nil {
		return v
	}

	// 1. Format the value to the desired precision
	// 2. In case of fractures, remove trailing 0's
	v.Value = strconv.FormatFloat(base, 'f', int(prec), 64)
	if strings.Contains(v.Value, ".") {
		v.Value = strings.TrimRight(v.Value, "0")
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
