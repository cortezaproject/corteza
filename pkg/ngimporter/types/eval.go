package types

//import (
//	"fmt"
//	"strconv"
//	"time"
//
//	"github.com/PaesslerAG/gval"
//)
//
//// Glang generates a gval language, that can be used for expression evaluation
//func GLang() gval.Language {
//	return gval.NewLanguage(
//		gval.Full(),
//
//		gval.Function("f64", func(v string) (float64, error) {
//			nn, err := strconv.ParseFloat(v, 64)
//			if err != nil {
//				return 0, err
//			}
//			return nn, nil
//		}),
//
//		gval.Function("concat", func(vv ...string) (string, error) {
//			out := ""
//			for _, v := range vv {
//				out += v
//			}
//			return out, nil
//		}),
//
//		gval.Function("numFmt", func(number, format string) (string, error) {
//			nn, err := strconv.ParseFloat(number, 64)
//			if err != nil {
//				return "", err
//			}
//
//			return fmt.Sprintf(format, nn), nil
//		}),
//
//		gval.Function("fFmt", func(number float64, format string) (string, error) {
//			return fmt.Sprintf(format, number), nil
//		}),
//
//		// diff between two dates in seconds
//		gval.Function("dateDiff", func(d1, d2 string) (float64, error) {
//			t1, err := time.Parse(SfDateTimeLayout, d1)
//			if err != nil {
//				return 0, err
//			}
//
//			t2, err := time.Parse(SfDateTimeLayout, d2)
//			if err != nil {
//				return 0, err
//			}
//
//			dr := t2.Sub(t1)
//			return dr.Seconds(), nil
//		}),
//	)
//}
