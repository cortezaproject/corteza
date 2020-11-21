package expr

import (
	"fmt"
	"github.com/PaesslerAG/gval"
)

func Parser(ll ...gval.Language) gval.Language {
	return gval.Full(append(AllFunctions(), ll...)...)
}

func AllFunctions() []gval.Language {
	ff := make([]gval.Language, 0, 100)

	//ff = append(ff, GenericFunctions()...)
	ff = append(ff, StringFunctions()...)
	ff = append(ff, NumericFunctions()...)
	ff = append(ff, TimeFunctions()...)

	return ff
}

// utility function for examples
func eval(e string, p interface{}) {
	result, err := Parser().Evaluate(e, p)
	if err != nil {
		fmt.Printf("error: %v", err)
	} else {
		fmt.Printf("%v", result)
	}
}
