package types

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/cortezaproject/corteza-server/pkg/filter"

	"github.com/pkg/errors"

	"github.com/cortezaproject/corteza-server/pkg/rbac"
)

type (
	Application struct {
		ID      uint64 `json:"applicationID,string"`
		Name    string `json:"name"`
		OwnerID uint64 `json:"ownerID"`
		Enabled bool   `json:"enabled"`
		Weight  int    `json:"weight"`

		Unify *ApplicationUnify `json:"unify,omitempty"`

		Labels map[string]string `json:"labels,omitempty"`
		Flags  []string          `json:"flags,omitempty"`

		CreatedAt time.Time  `json:"createdAt,omitempty"`
		UpdatedAt *time.Time `json:"updatedAt,omitempty"`
		DeletedAt *time.Time `json:"deletedAt,omitempty"`
	}

	ApplicationUnify struct {
		Name   string `json:"name,omitempty"`
		Listed bool   `json:"listed"`
		Icon   string `json:"icon"`
		Logo   string `json:"logo"`
		Url    string `json:"url"`
		Config string `json:"config"`
	}

	ApplicationFilter struct {
		Name  string `json:"name"`
		Query string `json:"query"`

		LabeledIDs []uint64          `json:"-"`
		Labels     map[string]string `json:"labels,omitempty"`

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

func (r *Application) DynamicRoles(userID uint64) []uint64 {
	return nil
}

// Resource returns a resource ID for this type
func (r Application) RBACResource() rbac.Resource {
	return ApplicationRBACResource.AppendID(r.ID)
}

func (au *ApplicationUnify) Scan(value interface{}) error {
	//lint:ignore S1034 This typecast is intentional, we need to get []byte out of a []uint8
	switch value.(type) {
	case nil:
		au = nil
	case []uint8:
		if err := json.Unmarshal(value.([]byte), au); err != nil {
			return errors.Wrapf(err, "Can not scan '%v' into ApplicationUnify", value)
		}
	}

	return nil
}

func (au ApplicationUnify) Value() (driver.Value, error) {
	return json.Marshal(au)
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
