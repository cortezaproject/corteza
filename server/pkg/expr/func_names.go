package expr

import "sort"

// BuiltInFunctionNames returns the registered function names, sorted.
//
// A name is the function only where an expression calls it: `split(a, b)` is
// the function, a bare `split` is a scope variable. The expression editor
// mirrors this list to colour calls and to flag an unknown one.
func BuiltInFunctionNames() []string {
	out := make([]string, 0, len(builtInFunctionNames))
	for n := range builtInFunctionNames {
		out = append(out, n)
	}
	sort.Strings(out)
	return out
}
