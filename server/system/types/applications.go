package types

import (
	"github.com/crusttech/human/server/pkg/filter"
	labelTypes "github.com/crusttech/human/server/pkg/label/types"
)

type (
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
)

func (a *Application) Valid() bool {
	return a.ID > 0 && a.DeletedAt == nil
}

// Dict exposes the fields a contextual role's expression may test, so a rule
// such as "resource.ownerID == userID" can scope editing to the author.
func (a Application) Dict() map[string]interface{} {
	dict := map[string]interface{}{
		"ID":        a.ID,
		"name":      a.Name,
		"enabled":   a.Enabled,
		"ownerID":   a.OwnerID,
		"createdAt": a.CreatedAt,
		"updatedAt": a.UpdatedAt,
		"deletedAt": a.DeletedAt,
	}

	if a.Unify != nil {
		dict["kind"] = a.Unify.Kind
	}

	return dict
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
