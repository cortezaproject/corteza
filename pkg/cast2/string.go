package cast2

import "github.com/spf13/cast"

func String(in any, out *string) error {
	aux, err := cast.ToStringE(in)
	if err != nil {
		return err
	}

	*out = aux
	return nil
}
