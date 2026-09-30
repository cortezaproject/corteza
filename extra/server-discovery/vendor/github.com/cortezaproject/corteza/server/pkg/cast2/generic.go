package cast2

import (
	"github.com/cortezaproject/corteza/server/pkg/id"
	"github.com/spf13/cast"
)

func ID(in any, out *id.ID) error {
	aux, err := cast.ToStringE(in)
	if err != nil {
		return err
	}

	*out, err = id.ByteID([]byte(aux))
	if err != nil {
		return err
	}

	return nil
}
