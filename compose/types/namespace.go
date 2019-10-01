package types

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/pkg/errors"

	"github.com/cortezaproject/corteza-server/pkg/permissions"
)

type (
	Namespace struct {
		ID      uint64        `db:"id"        json:"namespaceID,string"`
		Name    string        `db:"name"      json:"name"`
		Slug    string        `db:"slug"      json:"slug"`
		Enabled bool          `db:"enabled"   json:"enabled"`
		Meta    NamespaceMeta `db:"meta"      json:"meta"`

		CreatedAt time.Time  `db:"created_at"  json:"createdAt,omitempty"`
		UpdatedAt *time.Time `db:"updated_at"  json:"updatedAt,omitempty"`
		DeletedAt *time.Time `db:"deleted_at"  json:"deletedAt,omitempty"`
	}

	NamespaceFilter struct {
		Query   string `json:"query"`
		Slug    string `json:"slug"`
		Page    uint   `json:"page"`
		PerPage uint   `json:"perPage"`
		Sort    string `json:"sort"`
		Count   uint   `json:"count"`
	}

	NamespaceMeta struct {
		Subtitle    string `json:"subtitle,omitempty"    yaml:",omitempty"`
		Description string `json:"description,omitempty" yaml:",omitempty"`
	}
)

const (
	NamespaceCRM uint64 = 88714882739863655
)

// Resource returns a system resource ID for this type
func (n Namespace) PermissionResource() permissions.Resource {
	return NamespacePermissionResource.AppendID(n.ID)
}

// FindByHandle finds namespace by it's handle/slug
func (set NamespaceSet) FindByHandle(handle string) *Namespace {
	for i := range set {
		if set[i].Slug == handle {
			return set[i]
		}
	}

	return nil
}

func (nm *NamespaceMeta) Scan(value interface{}) error {
	//lint:ignore S1034 This typecast is intentional, we need to get []byte out of a []uint8
	switch value.(type) {
	case nil:
		*nm = NamespaceMeta{}
	case []uint8:
		b := value.([]byte)
		if err := json.Unmarshal(b, nm); err != nil {
			return errors.Wrapf(err, "Can not scan '%v' into NamespaceMeta", string(b))
		}
	}

	return nil
}

func (nm NamespaceMeta) Value() (driver.Value, error) {
	return json.Marshal(nm)
}
