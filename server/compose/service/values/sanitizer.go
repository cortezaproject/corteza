package values

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/cortezaproject/corteza/server/pkg/expr"
	"github.com/cortezaproject/corteza/server/pkg/logger"
	"github.com/cortezaproject/corteza/server/pkg/xss"
	"github.com/spf13/cast"
	"go.uber.org/zap"

	"github.com/cortezaproject/corteza/server/compose/types"
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
//   - fix multi-value order/place index
//   - trim all the strings!
//   - parse & format input values to match field specific -- nullify/falsify invalid
//   - field kind specific, no errors raised, data is modified
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

			v.Place = uint(i)
			out = append(out, v)
			i++
		}
	}

	var (
		log = logger.Default().
			With(logger.Uint64("module", m.ID))
	)

	wg := sync.WaitGroup{}
	for i := range out {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()

			f := m.Fields.FindByName(out[i].Name)
			if f == nil {
				// Unknown field,
				// if it is not handled before,
				// sanitizer does not care about it
				return
			}

			if f.Expressions.ValueExpr != "" {
				// do not do any sanitization if field has value expression!
				return
			}

			if out[i].IsDeleted() || !out[i].Updated {
				// Ignore unchanged and deleted
				return
			}

			kind := strings.ToLower(f.Kind)

			if len(f.Expressions.Sanitizers) > 0 {
				for _, expr := range f.Expressions.Sanitizers {
					rval, err := exprParser.Evaluate(expr, map[string]interface{}{"value": out[i].Value})
					if err != nil {
						log.Error(
							"failed to evaluate sanitizer expression",
							zap.String("field", f.Name),
							zap.String("expr", expr),
							zap.Error(err),
						)
						return
					}
					out[i].Value = sanitize(f, rval)
				}
			}

			if kind != "string" {
				// Trim all but string
				out[i].Value = strings.TrimSpace(out[i].Value)
			}

			if f.IsRef() {
				if refy.MatchString(out[i].Value) {
					out[i].Ref, _ = strconv.ParseUint(out[i].Value, 10, 64)
				}

				if out[i].Ref == 0 {
					out[i].Value = ""
				}
			}

			out[i].Value = sanitize(f, out[i].Value)
		}(i)
	}
	wg.Wait()

	return
}

func (s sanitizer) RunXSS(m *types.Module, vv types.RecordValueSet) types.RecordValueSet {
	var (
		f *types.ModuleField
	)

	wg := sync.WaitGroup{}
	for i := range vv {
		f = m.Fields.FindByName(vv[i].Name)
		if f == nil {
			// Unknown field,
			// if it is not handled before,
			// sanitizer does not care about it
			continue
		}

		switch strings.ToLower(f.Kind) {
		case "string":
			wg.Add(1)
			go func(i int) {
				vv[i].Value = sString(vv[i].Value)
				wg.Done()
			}(i)
		}
	}

	wg.Wait()

	return vv
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

	// Returning empty string here to align false value with everything else
	return strBoolFalseAlt
}

func sDatetime(v interface{}, onlyDate, onlyTime bool) string {
	var (
		// input format set
		inputFormats []string

		// output format
		internalFormat string

		datetime = fmt.Sprintf("%v", v)
	)

	if onlyTime {
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
		if onlyDate {
			// In case only date is used, make sure we format it properly
			internalFormat = datetimeInternalFormatDate
		} else {
			internalFormat = datetimeInternalFormatFull
		}

		// date & time
		inputFormats = []string{
			datetimeInternalFormatFull,
			"2006-01-02T15:04:05", // iso8601 without timezone
			time.RFC1123Z,
			time.RFC1123,
			time.RFC822Z,
			time.RFC822,
			time.RFC850,
			time.ANSIC,
			time.UnixDate,
			time.RubyDate,
			"2006-01-02 15:04:05.999999999 -0700 MST", // Time.String()
			"2006-01-02",
			"02 Jan 2006",
			"2006-01-02T15:04:05-0700", // RFC3339 without timezone hh:mm colon
			"2006-01-02 15:04:05 -07:00",
			"2006-01-02 15:04:05 -0700",
			"2006-01-02 15:04:05Z07:00", // RFC3339 without T
			"2006-01-02 15:04:05Z0700",  // RFC3339 without T or timezone hh:mm colon
			"2006-01-02 15:04:05",
			time.Kitchen,
			time.Stamp,
			time.StampMilli,
			time.StampMicro,
			time.StampNano,
			datetimeInternalFormatDate,
			"02 Jan 06",
			"Monday, 02-Jan-06",
			"Mon, 02 Jan 2006",
			"2006/_1/_2",
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

func sString(str interface{}) string {
	base, err := cast.ToStringE(str)
	if err != nil {
		return ""
	}

	return xss.RichText(base)
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
	case "string":
		v = sString(v)
	}

	return fmt.Sprintf("%v", v)
}
