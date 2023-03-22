package cast2

import "github.com/spf13/cast"

func Uint64(in any, out *uint64) error {
	aux, err := cast.ToUint64E(in)
	if err != nil {
		return err
	}

	*out = aux
	return nil
}

func Uint(in any, out *uint) error {
	aux, err := cast.ToUintE(in)
	if err != nil {
		return err
	}

	*out = aux
	return nil
}
