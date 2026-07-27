package expr

import "sort"

// builtInFunctionNames lists every identifier registered as a function in the
// expression language (see NumericFunctions, StringFunctions, ArrayFunctions,
// TimeFunctions, GenericFunctions and gvalFunc registrations).
//
// A scope variable that shares a name with one of these is effectively
// unreadable: a bare reference in an expression resolves to the function, not
// the variable, silently yielding the function's zero value. Consumers use
// IsBuiltInFunction to flag such shadowing.
//
// Keep this in sync with the gval.Function registrations; TestBuiltInFunctionNames
// guards against entries that are no longer reserved.
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

// IsBuiltInFunction reports whether name is a reserved expression-language
// function. A scope variable with such a name cannot be read back reliably.
func IsBuiltInFunction(name string) bool {
	_, ok := builtInFunctionNames[name]
	return ok
}

// BuiltInFunctionNames returns the reserved function names, sorted.
func BuiltInFunctionNames() []string {
	out := make([]string, 0, len(builtInFunctionNames))
	for n := range builtInFunctionNames {
		out = append(out, n)
	}
	sort.Strings(out)
	return out
}
