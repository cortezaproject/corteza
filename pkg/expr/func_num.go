package expr

import (
	"math"

	"github.com/PaesslerAG/gval"
)

func NumericFunctions() []gval.Language {
	return []gval.Language{
		gval.Function("min", min),
		gval.Function("max", max),
		gval.Function("round", round),
		gval.Function("floor", floor),
		gval.Function("ceil", ceil),
		gval.Function("abs", math.Abs),
		gval.Function("log", math.Log10),
		gval.Function("pow", math.Pow),
		gval.Function("sqrt", math.Sqrt),
		gval.Function("sum", sum),
		gval.Function("average", average),
	}
}

func min(aa ...interface{}) (min float64) {
	return findMinMax(-1, aa...)
}

func max(aa ...interface{}) (min float64) {
	return findMinMax(1, aa...)
}

func findMinMax(dir int, aa ...interface{}) (mm float64) {
	var (
		set bool
		flt float64
	)
	for i := range aa {
		switch conv := aa[i].(type) {
		case int:
			flt = float64(conv)
		case int64:
			flt = float64(conv)
		case uint:
			flt = float64(conv)
		case uint64:
			flt = float64(conv)
		case float32:
			flt = float64(conv)
		case float64:
			flt = conv
		default:
			continue
		}

		if !set {
			set = true
			mm = flt
		} else if dir < 0 {
			mm = math.Min(mm, flt)
		} else if dir > 0 {
			mm = math.Max(mm, flt)
		}
	}

	return mm
}

func round(f float64, d float64) float64 {
	p := math.Pow(10, d)
	return math.Round(f*p) / p
}

func floor(f float64) float64 {
	return math.Floor(f)
}

func ceil(f float64) float64 {
	return math.Ceil(f)
}

func sum(v ...interface{}) float64 {
	var sum float64

	for _, vv := range v {
		c, err := CastToFloat(vv)

		if err != nil {
			continue
		}

		sum += c
	}

	return sum
}

func average(v ...interface{}) float64 {
	var i int
	var j, sum float64

	for ; i < len(v); i++ {
		c, err := CastToFloat(v[i])

		if err != nil {
			continue
		}

		j++
		sum += c
	}

	return sum / j
}
