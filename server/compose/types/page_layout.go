package types

import (
	"database/sql/driver"
	"encoding/json"
	"strconv"
	"strings"
	"time"

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

		Meta *PageLayoutMeta `json:"meta"`

		Config PageLayoutConfig `json:"config"`
		Blocks PageBlocks       `json:"blocks"`

		Labels map[string]string `json:"labels,omitempty"`

		OwnedBy uint64 `json:"ownedBy,string"`

		CreatedAt time.Time  `json:"createdAt,omitempty"`
		UpdatedAt *time.Time `json:"updatedAt,omitempty"`
		DeletedAt *time.Time `json:"deletedAt,omitempty"`
	}

	PageLayoutMeta struct {
		// Warning: value of this field is now handled via resource-translation facility
		//          struct field is kept for the convenience for now since it allows us
		//          easy encoding/decoding of the outgoing/incoming values
		Name string `json:"name"`

		// Warning: value of this field is now handled via resource-translation facility
		//          struct field is kept for the convenience for now since it allows us
		//          easy encoding/decoding of the outgoing/incoming values
		Description string `json:"description"`
	}

	PageLayoutButton struct {
		Label   string `json:"label"`
		Enabled bool   `json:"enabled"`
	}

	PageLayoutConfig struct {
		Visibility PageLayoutVisibility `json:"visibility"`
		Actions    []PageLayoutAction   `json:"actions"`
	}

	PageLayoutVisibility struct {
		Expression string   `json:"expression"`
		Roles      []string `json:"roles"`
	}

	PageLayoutAction struct {
		ActionID  uint64 `json:"actionID,string"`
		Placement string
		Meta      any

		// Warning: value of this field is now handled via resource-translation facility
		//          struct field is kept for the convenience for now since it allows us
		//          easy encoding/decoding of the outgoing/incoming values
		Label string `json:"label"`

		// Kind and Params specify the action's behavior and the parameters it
		// can use for execution
		Kind   string
		Params any
	}

	PageLayoutFilter struct {
		NamespaceID uint64 `json:"namespaceID,string"`
		PageID      uint64 `json:"pageID,string,omitempty"`
		Primary     bool   `json:"primary,omitempty"`
		Handle      string `json:"handle"`
		Name        string `json:"name"`
		Query       string `json:"query"`

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

	for i, action := range p.Config.Actions {
		actionID := locale.ContentID(action.ActionID, i)
		rpl := strings.NewReplacer(
			"{{actionID}}", strconv.FormatUint(actionID, 10),
		)

		if aux = tt.FindByKey(rpl.Replace(LocaleKeyPagePageBlockBlockIDTitle.Path)); aux != nil {
			p.Config.Actions[i].Label = aux.Msg
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
			Key:      rpl.Replace(LocaleKeyPageLayoutConfigActionsActionIDLabel.Path),
			Msg:      action.Label,
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

func (bb *PageLayoutConfig) Scan(src any) error          { return sql.ParseJSON(src, bb) }
func (bb PageLayoutConfig) Value() (driver.Value, error) { return json.Marshal(bb) }

func (vv *PageLayoutMeta) Scan(src any) error           { return sql.ParseJSON(src, vv) }
func (vv *PageLayoutMeta) Value() (driver.Value, error) { return json.Marshal(vv) }
