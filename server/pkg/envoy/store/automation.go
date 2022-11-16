package store

import (
	"context"
	"strings"

	"github.com/cortezaproject/corteza/server/automation/types"
	"github.com/cortezaproject/corteza/server/pkg/envoy"
	"github.com/cortezaproject/corteza/server/pkg/envoy/resource"
	"github.com/cortezaproject/corteza/server/pkg/filter"
	"github.com/cortezaproject/corteza/server/store"
)

type (
	automationWorkflowFilter types.WorkflowFilter

	automationStore interface {
		store.AutomationWorkflows
		store.AutomationTriggers
	}

	automationDecoder struct {
		resourceID []uint64
		ux         *userIndex
	}
)

func newAutomationDecoder(ux *userIndex) *automationDecoder {
	return &automationDecoder{
		resourceID: make([]uint64, 0, 200),
		ux:         ux,
	}
}

func (d *automationDecoder) decodeWorkflows(ctx context.Context, s automationStore, ff []*automationWorkflowFilter) *auxRsp {
	mm := make([]envoy.Marshaller, 0, 100)
	if ff == nil {
		return &auxRsp{
			mm: mm,
		}
	}

	var nn types.WorkflowSet
	var fn types.WorkflowFilter
	var err error

	for _, f := range ff {
		aux := *f

		if aux.Limit == 0 {
			aux.Limit = 1000
		}

		for {
			nn, fn, err = s.SearchAutomationWorkflows(ctx, types.WorkflowFilter(aux))
			if err != nil {
				return &auxRsp{
					err: err,
				}
			}

			for _, n := range nn {
				// Index users
				err = d.ux.add(
					ctx,
					n.RunAs,
					n.OwnedBy,
					n.CreatedBy,
					n.UpdatedBy,
					n.DeletedBy,
				)
				if err != nil {
					return &auxRsp{
						err: err,
					}
				}

				tt, _, err := s.SearchAutomationTriggers(ctx, types.TriggerFilter{
					WorkflowID: []uint64{n.ID},
					Disabled:   filter.StateInclusive,
				})
				if err != nil {
					return &auxRsp{
						err: err,
					}
				}

				for _, t := range tt {
					// Index users
					err = d.ux.add(
						ctx,
						t.OwnedBy,
						t.CreatedBy,
						t.UpdatedBy,
						t.DeletedBy,
					)
				}
				if err != nil {
					return &auxRsp{
						err: err,
					}
				}

				mm = append(mm, newAutomationWorkflow(n, tt, d.ux))
				d.resourceID = append(d.resourceID, n.ID)
			}

			if fn.NextPage != nil {
				aux.PageCursor = fn.NextPage
			} else {
				break
			}
		}
	}

	return &auxRsp{
		mm: mm,
	}
}

func (df *DecodeFilter) automationFromResource(rr ...string) *DecodeFilter {
	for _, r := range rr {
		if !strings.HasPrefix(r, "automation") {
			continue
		}

		id := ""
		if strings.Count(r, ":") == 2 && !strings.HasSuffix(r, "*") {
			// There is an identifier
			aux := strings.Split(r, ":")

			id = aux[len(aux)-1]
			r = strings.Join(aux[:len(aux)-1], ":")
		}

		switch strings.ToLower(r) {
		case "automation:workflow":
			df = df.AutomationWorkflows(&types.WorkflowFilter{
				Query:    id,
				Disabled: filter.StateInclusive,
			})
		}
	}

	return df
}

func (df *DecodeFilter) automationFromRef(rr ...*resource.Ref) *DecodeFilter {
	for _, r := range rr {
		if strings.Index(r.ResourceType, "automation") < 0 {
			continue
		}

		switch r.ResourceType {
		case types.WorkflowResourceType:
			for _, i := range r.Identifiers.StringSlice() {
				df = df.AutomationWorkflows(&types.WorkflowFilter{
					Query:    i,
					Disabled: filter.StateInclusive,
				})
			}
		}
	}

	return df
}

// AutomationWorkflows adds a new WorkflowFilter
func (df *DecodeFilter) AutomationWorkflows(f *types.WorkflowFilter) *DecodeFilter {
	if df.automationWorkflow == nil {
		df.automationWorkflow = make([]*automationWorkflowFilter, 0, 1)
	}
	df.automationWorkflow = append(df.automationWorkflow, (*automationWorkflowFilter)(f))
	return df
}
