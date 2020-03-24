package types

import (
	"fmt"
	"strconv"
	"time"

	"github.com/PaesslerAG/gval"
)

// generates a simple gval language to be used within the migration
func Exprs() gval.Language {
	return gval.NewLanguage(
		gval.JSON(),
		gval.Arithmetic(),
		gval.PropositionalLogic(),

		gval.Function("numFmt", func(number, format string) (string, error) {
			nn, err := strconv.ParseFloat(number, 64)
			if err != nil {
				return "", err
			}

			return fmt.Sprintf(format, nn), nil
		}),

		gval.Function("fFmt", func(number float64, format string) (string, error) {
			return fmt.Sprintf(format, number), nil
		}),

		// diff between two dates in seconds
		gval.Function("dateDiff", func(d1, d2 string) (float64, error) {
			t1, err := time.Parse(SfDateTime, d1)
			if err != nil {
				return 0, err
			}

			t2, err := time.Parse(SfDateTime, d2)
			if err != nil {
				return 0, err
			}

			dr := t2.Sub(t1)
			return dr.Seconds(), nil
		}),
	)
}
