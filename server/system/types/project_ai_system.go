package types

import (
	"github.com/crusttech/human/server/pkg/filter"
)

type (
	ProjectAiSystemFilter struct {
		ProjectAiSystemID []uint64     `json:"projectAiSystemID"`
		TenantID          uint64       `json:"tenantID,string,omitempty"`
		ProjectID         uint64       `json:"projectID,string,omitempty"`
		Handle            string       `json:"handle,omitempty"`
		RiskClass         string       `json:"riskClass,omitempty"`
		Query             string       `json:"query,omitempty"`
		Deleted           filter.State `json:"deleted"`

		LabeledIDs []uint64 `json:"-"`

		Check func(*ProjectAiSystem) (bool, error) `json:"-"`

		filter.Sorting
		filter.Paging
	}

	ProjectAiSystemEntryFilter struct {
		ProjectAiSystemID uint64 `json:"projectAiSystemID,string,omitempty"`
		ResourceRef       string `json:"resourceRef,omitempty"`
		Limit             uint   `json:"-"`
	}
)
