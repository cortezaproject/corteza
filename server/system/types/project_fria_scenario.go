package types

import (
	"github.com/crusttech/human/server/pkg/filter"
)

type (
	// ProjectFriaScenarioFilter is the hand-written counterpart of the
	// filter.struct block in system/project_fria_scenario.cue — the store's
	// generated query builder reads exactly these fields.
	ProjectFriaScenarioFilter struct {
		ProjectFriaScenarioID []uint64 `json:"projectFriaScenarioID"`
		TenantID              uint64   `json:"tenantID,string,omitempty"`
		ProjectID             uint64   `json:"projectID,string,omitempty"`
		AiSystemID            uint64   `json:"aiSystemID,string,omitempty"`
		Severity              string   `json:"severity,omitempty"`
		Query                 string   `json:"query,omitempty"`

		Deleted filter.State `json:"deleted"`

		// Check fn is called by store backend for each resource found; it can
		// modify the resource and return false if store should not return it.
		Check func(*ProjectFriaScenario) (bool, error) `json:"-"`

		filter.Sorting
		filter.Paging
	}
)
