package types

import (
	"database/sql/driver"
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"github.com/cortezaproject/corteza/server/pkg/cast2"
	"github.com/cortezaproject/corteza/server/pkg/locale"
	"github.com/cortezaproject/corteza/server/pkg/sql"

	"github.com/cortezaproject/corteza/server/pkg/filter"
)

type (
	PageLayout struct {
		ID          uint64 `json:"pageLayoutID,string"`
		NamespaceID uint64 `json:"namespaceID,string"`
		PageID      uint64 `json:"pageID,string"`
		ParentID    uint64 `json:"parentID,string"`
		Handle      string `json:"handle"`
		Primary     bool   `json:"primary"`

		Weight int `json:"weight"`

		Meta PageLayoutMeta `json:"meta,omitempty"`

		Config PageLayoutConfig `json:"config"`
		Blocks PageLayoutBlocks `json:"blocks,omitempty"`

		Labels map[string]string `json:"labels,omitempty"`

		OwnedBy uint64 `json:"ownedBy,string"`

		CreatedAt time.Time  `json:"createdAt,omitempty"`
		UpdatedAt *time.Time `json:"updatedAt,omitempty"`
		DeletedAt *time.Time `json:"deletedAt,omitempty"`
	}

	PageLayoutBlocks []PageLayoutBlock

	PageLayoutBlock struct {
		BlockID uint64         `json:"blockID,string,omitempty" yaml:"blockID"`
		XYWH    [4]int         `json:"xywh" yaml:"xywh"`
		Meta    map[string]any `json:"meta,omitempty"`
	}

	PageLayoutMeta struct {
		// Warning: value of this field is now handled via resource-translation facility
		//          struct field is kept for the convenience for now since it allows us
		//          easy encoding/decoding of the outgoing/incoming values
		Title string `json:"title"`

		// Warning: value of this field is now handled via resource-translation facility
		//          struct field is kept for the convenience for now since it allows us
		//          easy encoding/decoding of the outgoing/incoming values
		Description string `json:"description"`

		Style map[string]any `json:"style,omitempty"`
	}

	PageLayoutButton struct {
		Enabled bool `json:"enabled"`

		// Warning: value of this field is now handled via resource-translation facility
		//          struct field is kept for the convenience for now since it allows us
		//          easy encoding/decoding of the outgoing/incoming values
		Label string `json:"label"`
	}

	PageLayoutButtonConfig struct {
		New    PageLayoutButton `json:"new"`
		Edit   PageLayoutButton `json:"edit"`
		Submit PageLayoutButton `json:"submit"`
		Delete PageLayoutButton `json:"delete"`
		Clone  PageLayoutButton `json:"clone"`
		Back   PageLayoutButton `json:"back"`
	}

	PageLayoutConfig struct {
		Visibility PageLayoutVisibility `json:"visibility"`

		Buttons PageLayoutButtonConfig `json:"buttons"`
		Actions []PageLayoutAction     `json:"actions,omitempty"`
	}

	PageLayoutVisibility struct {
		Expression string   `json:"expression"`
		Roles      []string `json:"roles,omitempty"`
	}

	PageLayoutAction struct {
		ActionID  uint64               `json:"actionID,string"`
		Placement string               `json:"placement"`
		Meta      PageLayoutActionMeta `json:"meta"`
		Enabled   bool                 `json:"enabled"`

		// Kind and Params specify the action's behavior and the parameters it
		// can use for execution
		Kind   string `json:"kind"`
		Params any    `json:"params"`
	}

	PageLayoutActionMeta struct {
		// Warning: value of this field is now handled via resource-translation facility
		//          struct field is kept for the convenience for now since it allows us
		//          easy encoding/decoding of the outgoing/incoming values
		Label string `json:"label"`

		Style map[string]any `json:"style,omitempty"`
	}

	PageLayoutFilter struct {
		PageLayoutID []string `json:"pageLayoutID"`
		NamespaceID  uint64   `json:"namespaceID,string"`
		PageID       uint64   `json:"pageID,string,omitempty"`
		ParentID     uint64   `json:"ParentID,string,omitempty"`
		Handle       string   `json:"handle"`
		Name         string   `json:"name"`
		Query        string   `json:"query"`

		LabeledIDs []uint64          `json:"-"`
		Labels     map[string]string `json:"labels,omitempty"`

		Deleted filter.State `json:"deleted"`

		// Check fn is called by store backend for each resource found function can
		// modify the resource and return false if store should not return it
		//
		// Store then loads additional resources to satisfy the paging parameters
		Check func(*PageLayout) (bool, error) `json:"-"`

		// Standard helpers for paging and sorting
		filter.Sorting
		filter.Paging
	}
)

func (m PageLayout) Clone() *PageLayout {
	c := &m
	return c
}

func (p *PageLayout) decodeTranslations(tt locale.ResourceTranslationIndex) {
	var aux *locale.ResourceTranslation

	// @note not doing blocks because they are simply copied from the page's index

	if aux = tt.FindByKey(LocaleKeyPageLayoutConfigButtonsNewLabel.Path); aux != nil {
		p.Config.Buttons.New.Label = aux.Msg
	}
	if aux = tt.FindByKey(LocaleKeyPageLayoutConfigButtonsEditLabel.Path); aux != nil {
		p.Config.Buttons.Edit.Label = aux.Msg
	}
	if aux = tt.FindByKey(LocaleKeyPageLayoutConfigButtonsSubmitLabel.Path); aux != nil {
		p.Config.Buttons.Submit.Label = aux.Msg
	}
	if aux = tt.FindByKey(LocaleKeyPageLayoutConfigButtonsDeleteLabel.Path); aux != nil {
		p.Config.Buttons.Delete.Label = aux.Msg
	}
	if aux = tt.FindByKey(LocaleKeyPageLayoutConfigButtonsCloneLabel.Path); aux != nil {
		p.Config.Buttons.Clone.Label = aux.Msg
	}
	if aux = tt.FindByKey(LocaleKeyPageLayoutConfigButtonsBackLabel.Path); aux != nil {
		p.Config.Buttons.Back.Label = aux.Msg
	}

	for i, action := range p.Config.Actions {
		actionID := locale.ContentID(action.ActionID, i)
		rpl := strings.NewReplacer(
			"{{actionID}}", strconv.FormatUint(actionID, 10),
		)

		if aux = tt.FindByKey(rpl.Replace(LocaleKeyPagePageBlockBlockIDTitle.Path)); aux != nil {
			p.Config.Actions[i].Meta.Label = aux.Msg
		}
	}
}

func (p *PageLayout) encodeTranslations() (out locale.ResourceTranslationSet) {
	out = make(locale.ResourceTranslationSet, 0, 8)

	// @note not doing blocks because they are simply copied from the page's index

	// Actions
	for i, action := range p.Config.Actions {
		actionID := locale.ContentID(action.ActionID, i)
		rpl := strings.NewReplacer(
			"{{actionID}}", strconv.FormatUint(actionID, 10),
		)

		out = append(out, &locale.ResourceTranslation{
			Resource: p.ResourceTranslation(),
			Key:      rpl.Replace(LocaleKeyPageLayoutConfigActionsActionIDMetaLabel.Path),
			Msg:      action.Meta.Label,
		})
	}

	return
}

// FindByHandle finds pageLayout by it's handle
func (set PageLayoutSet) FindByHandle(handle string) *PageLayout {
	for i := range set {
		if set[i].Handle == handle {
			return set[i]
		}
	}

	return nil
}

func (r *PageLayout) getValue(name string, pos uint) (any, error) {
	switch name {
	case "selfD", "SelfD":
		return r.PageID, nil
	}
	return nil, nil
}
func (r *PageLayout) setValue(name string, pos uint, value any) (err error) {
	switch name {
	case "selfID", "SelfID":
		return cast2.Uint64(value, &r.PageID)

	}
	return nil
}

func (bb *PageLayoutConfig) Scan(src any) error          { return sql.ParseJSON(src, &bb) }
func (bb PageLayoutConfig) Value() (driver.Value, error) { return json.Marshal(bb) }

func (vv *PageLayoutMeta) Scan(src any) error          { return sql.ParseJSON(src, &vv) }
func (vv PageLayoutMeta) Value() (driver.Value, error) { return json.Marshal(vv) }

func (bb *PageLayoutBlocks) Scan(src any) error          { return sql.ParseJSON(src, bb) }
func (bb PageLayoutBlocks) Value() (driver.Value, error) { return json.Marshal(bb) }
