package gvalfnc

import (
	"github.com/modern-go/reflect2"
	"github.com/spf13/cast"
)

func IsNil(i any) bool {
	return reflect2.IsNil(i)
}

func CastFloat(i any) (float64, error) {
	return cast.ToFloat64E(i)
}
func CastInt(i any) (int, error) {
	return cast.ToIntE(i)
}
func CastString(i any) (string, error) {
	return cast.ToStringE(i)
}
