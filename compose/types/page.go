package types

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/pkg/errors"

	"github.com/cortezaproject/corteza-server/internal/permissions"
)

type (
	// Page - page structure
	Page struct {
		ID     uint64 `json:"pageID,string" db:"id"`
		SelfID uint64 `json:"selfID,string" db:"self_id"`

		NamespaceID uint64 `json:"namespaceID,string" db:"rel_namespace"`

		ModuleID uint64 `json:"moduleID,string" db:"rel_module"`

		Handle      string `json:"handle" db:"handle"`
		Title       string `json:"title" db:"title"`
		Description string `json:"description" db:"description"`

		Blocks PageBlocks `json:"blocks" db:"blocks"`

		Children PageSet `json:"children,omitempty" db:"-"`

		Visible bool `json:"visible" db:"visible"`
		Weight  int  `json:"-" db:"weight"`

		CreatedAt time.Time  `db:"created_at" json:"createdAt,omitempty"`
		UpdatedAt *time.Time `db:"updated_at" json:"updatedAt,omitempty"`
		DeletedAt *time.Time `db:"deleted_at" json:"deletedAt,omitempty"`
	}

	PageBlocks []PageBlock

	PageBlock struct {
		Title       string                 `json:"title,omitempty"        yaml:",omitempty"`
		Description string                 `json:"description,omitempty"  yaml:",omitempty"`
		Options     map[string]interface{} `json:"options,omitempty"      yaml:",omitempty"`
		Style       PageBlockStyle         `json:"style,omitempty"        yaml:",omitempty"`
		Kind        string                 `json:"kind"`
		X           int                    `json:"x"`
		Y           int                    `json:"y"`
		Width       int                    `json:"width"`
		Height      int                    `json:"height"`
	}

	PageBlockStyle struct {
		Variants map[string]string `json:"variants,omitempty" yaml:",omitempty,flow"`
	}

	PageFilter struct {
		NamespaceID uint64 `json:"namespaceID,string"`
		ParentID    uint64 `json:"parentID,string,omitempty"`
		Root        bool   `json:"root,omitempty"`
		Handle      string `json:"handle"`
		Query       string `json:"query"`
		Page        uint   `json:"page"`
		PerPage     uint   `json:"perPage"`
		Count       uint   `json:"count"`
	}
)

// Resource returns a system resource ID for this type
func (p Page) PermissionResource() permissions.Resource {
	return PagePermissionResource.AppendID(p.ID)
}

// FindByHandle finds page by it's handle
func (set PageSet) FindByHandle(handle string) *Page {
	for i := range set {
		if set[i].Handle == handle {
			return set[i]
		}
	}

	return nil
}

func (bb *PageBlocks) Scan(value interface{}) error {
	//lint:ignore S1034 This typecast is intentional, we need to get []byte out of a []uint8
	switch value.(type) {
	case nil:
		*bb = PageBlocks{}
	case []uint8:
		b := value.([]byte)
		if err := json.Unmarshal(b, bb); err != nil {
			return errors.Wrapf(err, "Can not scan '%v' into PageBlocks", string(b))
		}
	}

	return nil
}

func (bb PageBlocks) Value() (driver.Value, error) {
	return json.Marshal(bb)
}

func (set PageSet) FindByParent(parentID uint64) (out PageSet) {
	out = PageSet{}
	for i := range set {
		if set[i].SelfID == parentID {
			out = append(out, set[i])
		}
	}

	return
}
