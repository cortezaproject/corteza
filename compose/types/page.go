package types

import (
	"database/sql/driver"
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"github.com/cortezaproject/corteza-server/pkg/filter"
	"github.com/cortezaproject/corteza-server/pkg/locale"
	"github.com/spf13/cast"

	"github.com/pkg/errors"
)

type (
	Page struct {
		ID     uint64 `json:"pageID,string"`
		SelfID uint64 `json:"selfID,string"`

		NamespaceID uint64 `json:"namespaceID,string"`

		ModuleID uint64 `json:"moduleID,string"`

		Handle string `json:"handle"`

		Config PageConfig `json:"config"`
		Blocks PageBlocks `json:"blocks"`

		Children PageSet `json:"children,omitempty"`

		Labels map[string]string `json:"labels,omitempty"`

		Visible bool `json:"visible"`
		Weight  int  `json:"weight"`

		CreatedAt time.Time  `json:"createdAt,omitempty"`
		UpdatedAt *time.Time `json:"updatedAt,omitempty"`
		DeletedAt *time.Time `json:"deletedAt,omitempty"`

		// Warning: value of this field is now handled via resource-translation facility
		//          struct field is kept for the convenience for now since it allows us
		//          easy encoding/decoding of the outgoing/incoming values
		Title string `json:"title"`

		// Warning: value of this field is now handled via resource-translation facility
		//          struct field is kept for the convenience for now since it allows us
		//          easy encoding/decoding of the outgoing/incoming values
		Description string `json:"description"`
	}

	PageBlocks []PageBlock

	PageBlock struct {
		BlockID uint64 `json:"blockID,string,omitempty"`

		Options map[string]interface{} `json:"options,omitempty"`
		Style   PageBlockStyle         `json:"style,omitempty"`
		Kind    string                 `json:"kind"`
		XYWH    [4]int                 `json:"xywh"` // x,y,w,h

		// Warning: value of this field is now handled via resource-translation facility
		//          struct field is kept for the convenience for now since it allows us
		//          easy encoding/decoding of the outgoing/incoming values
		Title string `json:"title,omitempty"`

		// Warning: value of this field is now handled via resource-translation facility
		//          struct field is kept for the convenience for now since it allows us
		//          easy encoding/decoding of the outgoing/incoming values
		Description string `json:"description,omitempty"`
	}

	PageBlockStyle struct {
		Variants map[string]string `json:"variants,omitempty"`
	}

	PageButton struct {
		Label   string `json:"label"`
		Enabled bool   `json:"enabled"`
	}

	PageButtonConfig struct {
		New    PageButton `json:"new"`
		Edit   PageButton `json:"edit"`
		Submit PageButton `json:"submit"`
		Delete PageButton `json:"delete"`
		Clone  PageButton `json:"clone"`
		Back   PageButton `json:"back"`
	}

	PageConfig struct {
		// How page is presented in the navigation
		NavItem struct {
			Icon *PageConfigIcon `json:"icon,omitempty"`
		} `json:"navItem"`

		Buttons *PageButtonConfig `json:"buttons,omitempty"`

		//// Example how page-config structure can evolve in the future
		//Views []struct {
		//	// what kind of output is this view intended for (screen, mobile...?)
		//	Output string
		//
		//	// Migrated page blocks, might be replaced someday with a more complex structure
		//	Blocks []PageBlock
		//}
	}

	PageConfigIcon struct {
		// Icon types and sources
		//
		// Note that backed does not enforce or validate all src value (types due to a limited
		// awareness of capabilities and
		//
		// Type: empty or "link" (default):
		// Indicate that src will contain an absolute or relative link to an icon.
		// Can also be used for inline images (storing "base64:" prefixed string in source).
		// This type and reference is not validated by the backend.
		//
		// Type: "library"
		// Source references an icon from a library. Ref's value should be in the following
		// notation: "font-awesome://<icon-identifier>".
		// This type and source is not validated by the backend.
		//
		// Type: "svg"
		// SRC contains raw SVG document

		////////////////////////////////////////////////////////////////////////////////////////////////////////
		// Other types that might be implemented in the future:
		// "attachment"
		// Reference (ID) to an existing attachment in local Corteza instance is expected
		// This type and reference must be validated by the backend.

		Type string `json:"type,omitempty"`
		Src  string `json:"src"`

		// Any custom styling that should be applied to the icon
		Style map[string]string `json:"style,omitempty"`
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

	PageChildrenDeleteStrategy string
)

const (
	PageChildrenOnDeleteAbort   PageChildrenDeleteStrategy = "abort"
	PageChildrenOnDeleteForce   PageChildrenDeleteStrategy = "force"
	PageChildrenOnDeleteRebase  PageChildrenDeleteStrategy = "rebase"
	PageChildrenOnDeleteCascade PageChildrenDeleteStrategy = "cascade"
)

func (m Page) Clone() *Page {
	c := &m
	return c
}

func (p *Page) decodeTranslations(tt locale.ResourceTranslationIndex) {
	var aux *locale.ResourceTranslation

	// Buttons here
	if p.Config.Buttons != nil {
		if aux = tt.FindByKey(LocaleKeyPageRecordToolbarNewLabel.Path); aux != nil {
			p.Config.Buttons.New.Label = aux.Msg
		}
		if aux = tt.FindByKey(LocaleKeyPageRecordToolbarEditLabel.Path); aux != nil {
			p.Config.Buttons.Edit.Label = aux.Msg
		}
		if aux = tt.FindByKey(LocaleKeyPageRecordToolbarSubmitLabel.Path); aux != nil {
			p.Config.Buttons.Submit.Label = aux.Msg
		}
		if aux = tt.FindByKey(LocaleKeyPageRecordToolbarDeleteLabel.Path); aux != nil {
			p.Config.Buttons.Delete.Label = aux.Msg
		}
		if aux = tt.FindByKey(LocaleKeyPageRecordToolbarCloneLabel.Path); aux != nil {
			p.Config.Buttons.Clone.Label = aux.Msg
		}
		if aux = tt.FindByKey(LocaleKeyPageRecordToolbarBackLabel.Path); aux != nil {
			p.Config.Buttons.Back.Label = aux.Msg
		}
	}

	for i, block := range p.Blocks {
		blockID := locale.ContentID(block.BlockID, i)
		rpl := strings.NewReplacer(
			"{{blockID}}", strconv.FormatUint(blockID, 10),
		)

		// - generic page block stuff
		if aux = tt.FindByKey(rpl.Replace(LocaleKeyPagePageBlockBlockIDTitle.Path)); aux != nil {
			p.Blocks[i].Title = aux.Msg
		}
		if aux = tt.FindByKey(rpl.Replace(LocaleKeyPagePageBlockBlockIDDescription.Path)); aux != nil {
			p.Blocks[i].Description = aux.Msg
		}

		// - automation page block stuff
		switch block.Kind {
		case "Automation":
			bb, _ := block.Options["buttons"].([]interface{})
			for j, auxBtn := range bb {
				btn := auxBtn.(map[string]interface{})

				buttonID := uint64(0)
				if aux, ok := btn["buttonID"]; ok {
					buttonID = cast.ToUint64(aux)
				}
				buttonID = locale.ContentID(buttonID, j)

				rpl := strings.NewReplacer(
					"{{blockID}}", strconv.FormatUint(blockID, 10),
					"{{buttonID}}", strconv.FormatUint(buttonID, 10),
				)

				if aux = tt.FindByKey(rpl.Replace(LocaleKeyPagePageBlockBlockIDButtonButtonIDLabel.Path)); aux != nil {
					btn["label"] = aux.Msg
				}
			}

		case "Content":
			if aux = tt.FindByKey(rpl.Replace(LocaleKeyPagePageBlockBlockIDContentBody.Path)); aux != nil {
				block.Options["body"] = aux.Msg
			}
		}

	}
}

func (p *Page) encodeTranslations() (out locale.ResourceTranslationSet) {
	out = make(locale.ResourceTranslationSet, 0, 12)

	// Button translations don't need to happen here as we don't do anything with buttons at the moment
	// @todo when we add custom buttons this should change

	// Page blocks
	for i, block := range p.Blocks {
		blockID := locale.ContentID(block.BlockID, i)
		rpl := strings.NewReplacer(
			"{{blockID}}", strconv.FormatUint(uint64(blockID), 10),
		)

		// - generic page block stuff
		out = append(out, &locale.ResourceTranslation{
			Resource: p.ResourceTranslation(),
			Key:      rpl.Replace(LocaleKeyPagePageBlockBlockIDTitle.Path),
			Msg:      block.Title,
		})

		out = append(out, &locale.ResourceTranslation{
			Resource: p.ResourceTranslation(),
			Key:      rpl.Replace(LocaleKeyPagePageBlockBlockIDDescription.Path),
			Msg:      block.Description,
		})

		// - automation page block stuff

		switch block.Kind {
		case "Automation":
			bb, _ := block.Options["buttons"].([]interface{})
			for j, auxBtn := range bb {
				btn := auxBtn.(map[string]interface{})

				if _, ok := btn["label"]; !ok {
					continue
				}

				buttonID := uint64(0)
				if aux, ok := btn["buttonID"]; ok {
					buttonID = cast.ToUint64(aux)
				}
				buttonID = locale.ContentID(buttonID, j)

				rpl := strings.NewReplacer(
					"{{blockID}}", strconv.FormatUint(blockID, 10),
					"{{buttonID}}", strconv.FormatUint(buttonID, 10),
				)

				out = append(out, &locale.ResourceTranslation{
					Resource: p.ResourceTranslation(),
					Key:      rpl.Replace(LocaleKeyPagePageBlockBlockIDButtonButtonIDLabel.Path),
					Msg:      btn["label"].(string),
				})

			}
		case "Content":
			body, _ := block.Options["body"].(string)
			out = append(out, &locale.ResourceTranslation{
				Resource: p.ResourceTranslation(),
				Key:      rpl.Replace(LocaleKeyPagePageBlockBlockIDContentBody.Path),
				Msg:      body,
			})
		}
	}

	return
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

// RecursiveWalk through all child pages
func (set PageSet) RecursiveWalk(parent *Page, fn func(c *Page, parent *Page) error) (err error) {
	if parent == nil {
		return
	}

	for _, page := range set {
		if page.SelfID != parent.ID {
			continue
		}

		if err = fn(page, parent); err != nil {
			return
		}

		if err = set.RecursiveWalk(page, fn); err != nil {
			return
		}
	}

	return
}

func (bb *PageConfig) Scan(value interface{}) error {
	//lint:ignore S1034 This typecast is intentional, we need to get []byte out of a []uint8
	switch value.(type) {
	case nil:
		*bb = PageConfig{}
	case []uint8:
		b := value.([]byte)
		if err := json.Unmarshal(b, bb); err != nil {
			return errors.Wrapf(err, "cannot scan '%v' into PageConfig", string(b))
		}
	}

	return nil
}

func (bb PageConfig) Value() (driver.Value, error) {
	// We're not saving button config to the DB; no need for it
	bb.Buttons = nil

	return json.Marshal(bb)
}
