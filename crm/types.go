package crm

import (
	"github.com/pkg/errors"
)

var _ = errors.Wrap

func (*Types) List(r *typesListRequest) (interface{}, error) {
	return nil, errors.New("Not implemented: Types.list")
}

func (*Types) Type(r *typesTypeRequest) (interface{}, error) {
	return nil, errors.New("Not implemented: Types.type")
}
