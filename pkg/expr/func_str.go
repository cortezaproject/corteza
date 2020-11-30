package expr

import (
	"github.com/PaesslerAG/gval"
	"strings"
)

func StringFunctions() []gval.Language {
	return []gval.Language{
		gval.Function("trim", strings.TrimSpace),
		gval.Function("trimLeft", strings.TrimLeft),
		gval.Function("trimRight", strings.TrimRight),
		gval.Function("length", length),
		gval.Function("toLower", strings.ToLower),
		gval.Function("toUpper", strings.ToUpper),
		gval.Function("shortest", shortest),
		gval.Function("longest", longest),
	}
}

func shortest(f string, aa ...string) string {
	for _, s := range aa {
		if len(f) > len(s) {
			f = s
		}
	}

	return f
}

func longest(f string, aa ...string) string {
	for _, s := range aa {
		if len(f) < len(s) {
			f = s
		}
	}

	return f
}

func length(s string) int {
	return len(s)
}
