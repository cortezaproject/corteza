package rest

import (
	"context"

	"github.com/cortezaproject/corteza/server/pkg/label/types"
	"github.com/cortezaproject/corteza/server/system/rest/request"
	"github.com/cortezaproject/corteza/server/system/service"
)

type (
	Label struct {
		label service.LabelService
	}
	LabelSetPayload struct {
		Filter types.LabelFilter `json:"filter"`
		Set types.LabelSet `json:"set"`
	}
)

func (Label) New() *Label {
	return &Label{
		label: service.Label(),
	}
}

func (ctrl Label) List(ctx context.Context, r *request.LabelList) (interface{}, error) {
	var (
		err error
		set types.LabelSet
		f = types.LabelFilter{
			Kind: r.Kind,
			Limit: uint(r.Limit),
		}
	)

	if r.Name != "" {
		f.Name = r.Name
	}

	set, f, err = ctrl.label.List(ctx, f)
	if err != nil {
		return nil, err
	}

	return &LabelSetPayload{
		Filter: f,
		Set: set,
	}, nil
}

