package types

import (
	"strconv"
	"strings"

	"github.com/crusttech/human/server/pkg/cast2"
	"github.com/crusttech/human/server/pkg/locale"

	"github.com/crusttech/human/server/pkg/filter"
	labelTypes "github.com/crusttech/human/server/pkg/label/types"
)

type (
	// PageLayout is the resource struct; it is generated into page_layout.gen.go.

	PageLayoutBlocks []PageLayoutBlock

	PageLayoutButtonConfig struct {
		New    PageLayoutButton `json:"new"`
		Edit   PageLayoutButton `json:"edit"`
		Submit PageLayoutButton `json:"submit"`
		Delete PageLayoutButton `json:"delete"`
		Clone  PageLayoutButton `json:"clone"`
		Back   PageLayoutButton `json:"back"`
	}

	PageLayoutFilter struct {
		PageLayoutID []string `json:"pageLayoutID"`
		TenantID     uint64   `json:"tenantID,string,omitempty"`
		ProjectID    uint64   `json:"projectID,string,omitempty"`
		NamespaceID  uint64   `json:"namespaceID,string"`
		PageID       uint64   `json:"pageID,string,omitempty"`
		ParentID     uint64   `json:"ParentID,string,omitempty"`
		Handle       string   `json:"handle"`
		Name         string   `json:"name"`
		Query        string   `json:"query"`

		LabeledIDs []uint64                         `json:"-"`
		Labels     map[string]labelTypes.LabelValue `json:"labels,omitempty"`

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

// Dict exposes page layout attributes for RBAC contextual role evaluation.
func (p PageLayout) Dict() map[string]interface{} {
	return map[string]interface{}{
		"ID":             p.ID,
		"pageLayoutID":   p.ID,
		"namespaceID":    p.NamespaceID,
		"pageID":         p.PageID,
		"parentID":       p.ParentID,
		"handle":         p.Handle,
		"ownedBy":        p.OwnedBy,
		"createdAt":      p.CreatedAt,
		"createdByAgent": p.CreatedByAgent,
		"updatedAt":      p.UpdatedAt,
		"deletedAt":      p.DeletedAt,
	}
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
