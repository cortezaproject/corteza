package values

import (
	"fmt"
	"github.com/cortezaproject/corteza-server/pkg/expr"
	"github.com/cortezaproject/corteza-server/pkg/logger"
	"go.uber.org/zap"
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
	var (
		exprParser = expr.Parser()
	)

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

		log = logger.Default().
			WithOptions(zap.AddStacktrace(zap.PanicLevel)).
			With(zap.Uint64("module", m.ID))
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

		if len(f.Expressions.Sanitizers) > 0 {
			for _, expr := range f.Expressions.Sanitizers {
				rval, err := exprParser.Evaluate(expr, map[string]interface{}{"value": v.Value})
				if err != nil {
					log.Error(
						"failed to evaluate sanitizer expression",
						zap.String("field", f.Name),
						zap.String("expr", expr),
						zap.Error(err),
					)
					continue
				}
				v.Value = sanitize(f, rval)
			}
		}

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
			v.Value = sBool(v.Value)

		case "datetime":
			v.Value = sDatetime(v.Value, f.Options.Bool("onlyDate"), f.Options.Bool("onlyTime"))

		case "number":
			v.Value = sNumber(v.Value, f.Options.Precision())

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

func sBool(v interface{}) string {
	switch c := v.(type) {
	case bool:
		if c {
			return strBoolTrue
		}

	case string:
		if truthy.MatchString(strings.ToLower(c)) {
			return strBoolTrue
		}
	}

	return strBoolFalse
}

func sDatetime(v interface{}, onlyDate, onlyTime bool) string {
	var (
		// input format set
		inputFormats []string

		// output format
		internalFormat string

		datetime = fmt.Sprintf("%v", v)
	)

	if onlyDate {
		internalFormat = datetimeInternalFormatDate
		inputFormats = []string{
			datetimeInternalFormatDate,
			"02 Jan 06",
			"Monday, 02-Jan-06",
			"Mon, 02 Jan 2006",
			"2006/_1/_2",
		}
	} else if onlyTime {
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
		if isoDaty.MatchString(datetime) && !hasTimezone.MatchString(datetime) {
			// No timezone, add Z to satisfy parser
			datetime = datetime + "Z"

			// Simplifiy list of rules
			inputFormats = []string{time.RFC3339}
		}
	}

	for _, format := range inputFormats {
		parsed, err := time.Parse(format, datetime)
		if err == nil {
			return parsed.UTC().Format(internalFormat)
		}
	}

	return ""
}

func sNumber(num interface{}, p uint) string {
	base, err := strconv.ParseFloat(fmt.Sprintf("%v", num), 64)
	if err != nil {
		return "0"
	}

	// Format the value to the desired precision
	str := strconv.FormatFloat(base, 'f', int(p), 64)

	// In case of fractures, remove trailing 0's
	if strings.Contains(str, ".") {
		str = strings.TrimRight(str, "0")
		str = strings.TrimRight(str, ".")
	}

	return str
}

// sanitize casts value to field kind format
func sanitize(f *types.ModuleField, v interface{}) string {
	switch strings.ToLower(f.Kind) {
	case "bool":
		return sBool(v)
	case "datetime":
		v = sDatetime(v, f.Options.Bool("onlyDate"), f.Options.Bool("onlyTime"))
	case "number":
		v = sNumber(v, f.Options.Precision())
	}

	return fmt.Sprintf("%v", v)
}
