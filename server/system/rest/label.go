package rest

import (
	"context"

	"github.com/crusttech/human/server/pkg/api"
	"github.com/crusttech/human/server/pkg/label/types"
	"github.com/crusttech/human/server/system/rest/request"
	"github.com/crusttech/human/server/system/service"
)

type (
	Label struct {
		label service.LabelService
	}

	LabelListEntry struct {
		Name          string `json:"name"`
		ResourceCount int    `json:"resourceCount"`
	}

	LabelSetPayload struct {
		Filter types.LabelFilter `json:"filter"`
		Set    []LabelListEntry  `json:"set"`
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
		f   = types.LabelFilter{
			Kind:  r.Kind,
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

	// Deduplicate by label name and count resources per label
	counts := make(map[string]int)
	order := make([]string, 0)
	for _, label := range set {
		if counts[label.Name] == 0 {
			order = append(order, label.Name)
		}
		counts[label.Name]++
	}

	unique := make([]LabelListEntry, 0, len(order))
	for _, name := range order {
		unique = append(unique, LabelListEntry{
			Name:          name,
			ResourceCount: counts[name],
		})
	}

	return &LabelSetPayload{
		Filter: f,
		Set:    unique,
	}, nil
}

func (ctrl Label) Delete(ctx context.Context, r *request.LabelDelete) (interface{}, error) {
	return api.OK(), ctrl.label.Delete(ctx, r.Name, r.Kind)
}
