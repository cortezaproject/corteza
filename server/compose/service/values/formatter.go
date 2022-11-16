package values

import (
	"strings"
	"time"

	"github.com/cortezaproject/corteza/server/compose/types"
	"github.com/cortezaproject/corteza/server/pkg/expr"
	"github.com/cortezaproject/corteza/server/pkg/logger"
	"go.uber.org/zap"
)

type (
	formatter struct{}
)

// Formatter initializes formatter
//
// Not really needed, following pattern in the package
func Formatter() *formatter {
	return &formatter{}
}

func (f formatter) Run(m *types.Module, vv types.RecordValueSet) types.RecordValueSet {
	var (
		exprParser = expr.Parser()

		log = logger.Default().
			With(zap.Uint64("module", m.ID))
	)

	for _, v := range vv {
		fld := m.Fields.FindByName(v.Name)
		if fld == nil {
			// Unknown field,
			// if it is not handled before,
			// formatter does not care about it
			continue
		}

		if fld.Expressions.ValueExpr != "" {
			// do not do any validation if field has value expression!
			continue
		}

		if !(fld.Expressions.DisableDefaultFormatters && len(fld.Expressions.Formatters) > 0) {
			// Per field type validators
			switch strings.ToLower(fld.Kind) {
			case "bool":
				v = f.fBool(v)
			case "datetime":
				v = f.fDatetime(v, fld)
			}
		}

		if len(fld.Expressions.Formatters) > 0 {
			for _, expr := range fld.Expressions.Formatters {
				rval, err := exprParser.Evaluate(expr, map[string]interface{}{"value": v.Value})
				if err != nil {
					log.Error(
						"failed to evaluate sanitizer expression",
						zap.String("field", fld.Name),
						zap.String("expr", expr),
						zap.Error(err),
					)
					continue
				}
				v.Value = sanitize(fld, rval)
			}
		}
	}

	return vv
}

// Boolean values are outputed as "1" (true) and "" (false)
func (formatter) fBool(v *types.RecordValue) *types.RecordValue {
	if v.Value != strBoolTrue {
		v.Value = ""
	}

	return v
}

func (formatter) fDatetime(v *types.RecordValue, f *types.ModuleField) *types.RecordValue {
	var (
		// database formats
		dbFormats []string

		// output format
		internalFormat string
	)

	if f.Options.Bool("onlyDate") {
		internalFormat = datetimeInternalFormatDate
		dbFormats = []string{
			datetimeInternalFormatDate,
		}
	} else if f.Options.Bool("onlyTime") {
		internalFormat = datetimeIntenralFormatTime
		dbFormats = []string{
			datetimeIntenralFormatTime,
			"15:04",
		}
	} else {
		internalFormat = datetimeInternalFormatFull
		// date & time
		dbFormats = []string{
			datetimeInternalFormatFull,
			"2006-01-02 15:04:05",
			"2006-01-02 15:04",
		}
	}

	for _, format := range dbFormats {
		parsed, err := time.Parse(format, v.Value)
		if err == nil {
			v.Value = parsed.UTC().Format(internalFormat)
			return v
		}
	}

	v.Value = ""
	return v

}
