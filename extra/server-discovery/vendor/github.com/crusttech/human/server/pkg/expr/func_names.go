package expr

import "sort"

// builtInFunctionNames lists every identifier registered as a function in the
// expression language (see NumericFunctions, StringFunctions, ArrayFunctions,
// TimeFunctions, GenericFunctions and gvalFunc registrations).
//
// A name is the function only where an expression calls it: `split(a, b)` is
// the function, a bare `split` is a scope variable. The expression editor
// mirrors this list to colour calls and to flag an unknown one, and
// TestBuiltInFunctionNames holds it to what the language registers.
var builtInFunctionNames = map[string]struct{}{
	"abs": {}, "average": {}, "base64encode": {}, "camelize": {}, "ceil": {},
	"coalesce": {}, "count": {}, "earliest": {}, "filter": {}, "find": {},
	"float": {}, "floor": {}, "format": {}, "has": {}, "hasAll": {},
	"hasPrefix": {}, "hasSubstring": {}, "hasSuffix": {}, "int": {}, "isEmail": {},
	"isEmpty": {}, "isLeapYear": {}, "isNil": {}, "isUrl": {}, "isWeekDay": {},
	"join": {}, "latest": {}, "length": {}, "log": {}, "longest": {},
	"match": {}, "max": {}, "merge": {}, "min": {}, "modDate": {},
	"modMonth": {}, "modTime": {}, "modWeek": {}, "modYear": {}, "now": {},
	"omit": {}, "parseDuration": {}, "parseISOTime": {}, "pop": {}, "pow": {},
	"push": {}, "random": {}, "round": {}, "set": {}, "shift": {},
	"shortest": {}, "snakify": {}, "sort": {}, "splice": {}, "split": {},
	"sqrt": {}, "strftime": {}, "sub": {}, "sum": {}, "title": {},
	"toJSON": {}, "toLower": {}, "toUpper": {}, "trim": {}, "trimLeft": {},
	"trimRight": {}, "untitle": {},
}

// BuiltInFunctionNames returns the registered function names, sorted.
func BuiltInFunctionNames() []string {
	out := make([]string, 0, len(builtInFunctionNames))
	for n := range builtInFunctionNames {
		out = append(out, n)
	}
	sort.Strings(out)
	return out
}
