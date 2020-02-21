package values

import (
	"github.com/cortezaproject/corteza-server/compose/types"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	boolTrue  = "1"
	boolFalse = "0"
)

type (
	sanitizer struct {
	}
)

var (
	// value resembles something that can be true
	truthy = regexp.MustCompile(`^(t(rue)?|y(es)?|1)$`)

	// valeu resembles something that can be a reference
	refy = regexp.MustCompile(`^[1-9](\d*)$`)
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
			v = s.sanitizeBool(v, f, m)
		case "datetime":
			v = s.sanitizeDatetime(v, f, m)
		case "email":
			v = s.sanitizeEmail(v, f, m)
		case "file":
			v = s.sanitizeFile(v, f, m)
		case "number":
			v = s.sanitizeNumber(v, f, m)
		case "record":
			v = s.sanitizeRecord(v, f, m)
		case "select":
			v = s.sanitizeSelect(v, f, m)
		case "string":
			v = s.sanitizeString(v, f, m)
		case "url":
			v = s.sanitizeUrl(v, f, m)
		case "user":
			v = s.sanitizeUser(v, f, m)
		}
	}

	return
}

func (sanitizer) sanitizeBool(v *types.RecordValue, f *types.ModuleField, m *types.Module) *types.RecordValue {
	if truthy.MatchString(strings.ToLower(v.Value)) {
		v.Value = boolTrue
	} else {
		v.Value = boolFalse
	}

	return v
}

func (sanitizer) sanitizeDatetime(v *types.RecordValue, f *types.ModuleField, m *types.Module) *types.RecordValue {
	var (
		// input format set
		inputFormats []string

		// output format
		of string
	)

	if f.Options.Bool("onlyDate") {
		of = "2006-01-02"
		inputFormats = []string{
			"2006-01-02",
			"02 Jan 06",
			"Monday, 02-Jan-06",
			"Mon, 02 Jan 2006",
			"2019/_1/_2",
		}
	} else if f.Options.Bool("onlyTime") {
		of = "15:04:05"
		inputFormats = []string{
			"15:04:05",
			"15:04",
			time.Kitchen,
		}
	} else {
		of = time.RFC3339
		// date & time
		inputFormats = []string{
			time.RFC3339,
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

func (sanitizer) sanitizeEmail(v *types.RecordValue, f *types.ModuleField, m *types.Module) *types.RecordValue {
	return v
}

func (sanitizer) sanitizeFile(v *types.RecordValue, f *types.ModuleField, m *types.Module) *types.RecordValue {
	return v
}

func (sanitizer) sanitizeNumber(v *types.RecordValue, f *types.ModuleField, m *types.Module) *types.RecordValue {
	return v
}

func (sanitizer) sanitizeRecord(v *types.RecordValue, f *types.ModuleField, m *types.Module) *types.RecordValue {
	return v
}

func (sanitizer) sanitizeSelect(v *types.RecordValue, f *types.ModuleField, m *types.Module) *types.RecordValue {
	return v
}

func (sanitizer) sanitizeString(v *types.RecordValue, f *types.ModuleField, m *types.Module) *types.RecordValue {
	return v
}

func (sanitizer) sanitizeUrl(v *types.RecordValue, f *types.ModuleField, m *types.Module) *types.RecordValue {
	return v
}

func (sanitizer) sanitizeUser(v *types.RecordValue, f *types.ModuleField, m *types.Module) *types.RecordValue {
	return v
}
