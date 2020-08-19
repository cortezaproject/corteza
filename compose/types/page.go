package types

import (
	"database/sql/driver"
	"encoding/json"
	"github.com/cortezaproject/corteza-server/store"
	"time"

	"github.com/pkg/errors"

	"github.com/cortezaproject/corteza-server/pkg/permissions"
	"github.com/cortezaproject/corteza-server/pkg/rh"
)

type (
	// Page - page structure
	Page struct {
		ID     uint64 `json:"pageID,string"`
		SelfID uint64 `json:"selfID,string"`

		NamespaceID uint64 `json:"namespaceID,string"`

		ModuleID uint64 `json:"moduleID,string"`

		Handle      string `json:"handle"`
		Title       string `json:"title"`
		Description string `json:"description"`

		Blocks PageBlocks `json:"blocks"`

		Children PageSet `json:"children,omitempty" db:"-"`

		Visible bool `json:"visible"`
		Weight  int  `json:"weight"`

		CreatedAt time.Time  `json:"createdAt,omitempty"`
		UpdatedAt *time.Time `json:"updatedAt,omitempty"`
		DeletedAt *time.Time `json:"deletedAt,omitempty"`
	}

	PageBlocks []PageBlock

	PageBlock struct {
		Title       string                 `json:"title,omitempty"        yaml:",omitempty"`
		Description string                 `json:"description,omitempty"  yaml:",omitempty"`
		Options     map[string]interface{} `json:"options,omitempty"      yaml:",omitempty"`
		Style       PageBlockStyle         `json:"style,omitempty"        yaml:",omitempty"`
		Kind        string                 `json:"kind"`
		XYWH        [4]int                 `json:"xywh"                    yaml:"xywh,flow"` // x,y,w,h
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

		Deleted rh.FilterState `json:"deleted"`

		// Check fn is called by store backend for each resource found function can
		// modify the resource and return false if store should not return it
		//
		// Store then loads additional resources to satisfy the paging parameters
		Check func(*Page) (bool, error) `json:"-"`

		// Standard helpers for paging and sorting
		store.Sorting
		store.Paging
	}
)

// Resource returns a system resource ID for this type
func (p Page) PermissionResource() permissions.Resource {
	return PagePermissionResource.AppendID(p.ID)
}

func (p Page) DynamicRoles(userID uint64) []uint64 {
	return nil
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

// Helper to extract old encoding to new one
func (b *PageBlock) UnmarshalJSON(data []byte) (err error) {
	type internalPageBlock PageBlock
	i := struct {
		internalPageBlock
		X      int `json:"x,omitempty"            yaml:"-"`
		Y      int `json:"y,omitempty"            yaml:"-"`
		Width  int `json:"width,omitempty"        yaml:"-"`
		Height int `json:"height,omitempty"       yaml:"-"`
	}{}

	if err = json.Unmarshal(data, &i); err != nil {
		return
	}

	*b = PageBlock(i.internalPageBlock)
	if i.XYWH[0]+i.XYWH[1]+i.XYWH[2]+i.XYWH[3] > 0 {
		return nil
	}

	if i.X+i.Y+i.Width+i.Height > 0 {
		// Cast old x,y,w,h structure to this:
		b.XYWH = [4]int{i.X, i.Y, i.Width, i.Height}
	}

	return nil
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
