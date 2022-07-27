package cast2

import (
	"github.com/modern-go/reflect2"
	"github.com/spf13/cast"
	"time"
)

func Time(in any, out *time.Time) error {
	if reflect2.IsNil(in) {
		if out != nil {
			*out = time.Time{}
		}

		return nil
	}

	aux, err := cast.ToTimeE(in)
	if err != nil {
		return err
	}

	*out = aux
	return nil
}

func TimePtr(in any, out **time.Time) error {
	if reflect2.IsNil(in) {
		*out = nil
		return nil
	}

	aux, err := cast.ToTimeE(in)
	if err != nil {
		return err
	}

	*out = &aux
	return nil
}
