package types

import (
	"github.com/crusttech/human/server/pkg/filter"
	labelTypes "github.com/crusttech/human/server/pkg/label/types"
)

type (
	ApplicationUnify struct {
		Name   string `json:"name,omitempty"`
		Listed bool   `json:"listed"`
		Url    string `json:"url"`
		Config string `json:"config"`

		// Temporary icon & logo URLs
		// @todo rework this when we rework attachment management
		Icon   string `json:"icon,omitempty"`
		IconID uint64 `json:"iconID,string"`
		Logo   string `json:"logo,omitempty"`
		LogoID uint64 `json:"logoID,string"`
	}

	ApplicationFilter struct {
		ApplicationID []string `json:"applicationID"`
		Name          string   `json:"name"`
		Query         string   `json:"query"`

		LabeledIDs []uint64                         `json:"-"`
		Labels     map[string]labelTypes.LabelValue `json:"labels,omitempty"`

		FlaggedIDs []uint64 `json:"-"`
		Flags      []string `json:"flags,omitempty"`
		IncFlags   uint     `json:"-"`

		Deleted filter.State `json:"deleted"`

		// Check fn is called by store backend for each resource found function can
		// modify the resource and return false if store should not return it
		//
		// Store then loads additional resources to satisfy the paging parameters
		Check func(*Application) (bool, error) `json:"-"`

		// Standard helpers for paging and sorting
		filter.Sorting
		filter.Paging
	}

	ApplicationMetrics struct {
		Total   uint `json:"total"`
		Deleted uint `json:"deleted"`
		Valid   uint `json:"valid"`
	}
)

func (a *Application) Valid() bool {
	return a.ID > 0 && a.DeletedAt == nil
}

// // // These will get generated later on

// SetFlags adds new label to label map
func (a *Application) SetFlags(flags []string) {
	a.Flags = flags
}

// GetFlags returns current flags on the resource
func (a *Application) GetFlags() []string {
	return a.Flags
}

// FlagResourceKind returns the resource kind for the flag
func (*Application) FlagResourceKind() string {
	return "system:application"
}

// GetLabels adds new label to label map
func (a *Application) FlagResourceID() uint64 {
	return a.ID
}
