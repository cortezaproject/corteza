package cast2

import "github.com/spf13/cast"

func Bool(in any, out *bool) error {
	aux, err := cast.ToBoolE(in)
	if err != nil {
		return err
	}

	*out = aux
	return nil
}
