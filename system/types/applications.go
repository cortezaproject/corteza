package types

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/pkg/errors"

	"github.com/cortezaproject/corteza-server/pkg/permissions"
	"github.com/cortezaproject/corteza-server/pkg/rh"
)

type (
	Application struct {
		ID      uint64 `json:"applicationID,string" db:"id"`
		Name    string `json:"name" db:"name"`
		OwnerID uint64 `json:"ownerID" db:"rel_owner"`
		Enabled bool   `json:"enabled" db:"enabled"`

		Unify *ApplicationUnify `json:"unify,omitempty" db:"unify"`

		CreatedAt time.Time  `json:"createdAt,omitempty" db:"created_at"`
		UpdatedAt *time.Time `json:"updatedAt,omitempty" db:"updated_at"`
		DeletedAt *time.Time `json:"deletedAt,omitempty" db:"deleted_at"`
	}

	ApplicationUnify struct {
		Name   string `json:"name,omitempty"`
		Listed bool   `json:"listed"`
		Icon   string `json:"icon"`
		Logo   string `json:"logo"`
		Url    string `json:"url"`
		Config string `json:"config"`
		Order  uint   `json:"order"`
	}

	ApplicationFilter struct {
		Name  string `json:"name"`
		Query string `json:"query"`

		Sort string `json:"sort"`

		// Standard paging fields & helpers
		rh.PageFilter

		// Resource permission check filter
		IsReadable *permissions.ResourceFilter `json:"-"`

		Deleted rh.FilterState `json:"deleted"`
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
func (r Application) PermissionResource() permissions.Resource {
	return ApplicationPermissionResource.AppendID(r.ID)
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
