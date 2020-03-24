package types

import (
	"fmt"
	"strconv"

	"github.com/PaesslerAG/gval"
)

// generates a simple gval language to be used within the migration
func exprs() gval.Language {
	return gval.NewLanguage(
		gval.JSON(),
		gval.Arithmetic(),

		gval.Function("numFmt", func(number, format string) (string, error) {
			nn, err := strconv.ParseFloat(number, 64)
			if err != nil {
				return "", err
			}

			return fmt.Sprintf(format, nn), nil
		}),
	)
}
