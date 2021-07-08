package types

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/cortezaproject/corteza-server/pkg/filter"

	"github.com/pkg/errors"
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

		Children PageSet `json:"children,omitempty"`

		Labels map[string]string `json:"labels,omitempty"`

		Visible bool `json:"visible"`
		Weight  int  `json:"weight"`

		CreatedAt time.Time  `json:"createdAt,omitempty"`
		UpdatedAt *time.Time `json:"updatedAt,omitempty"`
		DeletedAt *time.Time `json:"deletedAt,omitempty"`
	}

	PageBlocks []PageBlock

	PageBlock struct {
		Title       string                 `json:"title,omitempty"`
		Description string                 `json:"description,omitempty"`
		Options     map[string]interface{} `json:"options,omitempty"`
		Style       PageBlockStyle         `json:"style,omitempty"`
		Kind        string                 `json:"kind"`
		XYWH        [4]int                 `json:"xywh"` // x,y,w,h
	}

	PageBlockStyle struct {
		Variants map[string]string `json:"variants,omitempty"`
	}

	PageFilter struct {
		NamespaceID uint64 `json:"namespaceID,string"`
		ParentID    uint64 `json:"parentID,string,omitempty"`
		ModuleID    uint64 `json:"moduleID,string,omitempty"`
		Root        bool   `json:"root,omitempty"`
		Handle      string `json:"handle"`
		Title       string `json:"title"`
		Query       string `json:"query"`

		LabeledIDs []uint64          `json:"-"`
		Labels     map[string]string `json:"labels,omitempty"`

		Deleted filter.State `json:"deleted"`

		// Check fn is called by store backend for each resource found function can
		// modify the resource and return false if store should not return it
		//
		// Store then loads additional resources to satisfy the paging parameters
		Check func(*Page) (bool, error) `json:"-"`

		// Standard helpers for paging and sorting
		filter.Sorting
		filter.Paging
	}
)

func (m Page) Clone() *Page {
	c := &m
	return c
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
			return errors.Wrapf(err, "cannot scan '%v' into PageBlocks", string(b))
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
		X      int `json:"x,omitempty"`
		Y      int `json:"y,omitempty"`
		Width  int `json:"width,omitempty"`
		Height int `json:"height,omitempty"`
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
