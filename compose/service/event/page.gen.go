package event

// This file is auto-generated.
//
// YAML event definitions:
//   compose/service/event/events.yaml
//
// Regenerate with:
//   go run codegen/v2/events.go --service compose
//

import (
	"encoding/json"

	"github.com/cortezaproject/corteza-server/compose/types"

	"github.com/cortezaproject/corteza-server/pkg/auth"
)

type (
	// pageBase
	//
	// This type is auto-generated.
	pageBase struct {
		immutable bool
		page      *types.Page
		oldPage   *types.Page
		namespace *types.Namespace
		invoker   auth.Identifiable
	}

	// pageOnManual
	//
	// This type is auto-generated.
	pageOnManual struct {
		*pageBase
	}

	// pageBeforeCreate
	//
	// This type is auto-generated.
	pageBeforeCreate struct {
		*pageBase
	}

	// pageBeforeUpdate
	//
	// This type is auto-generated.
	pageBeforeUpdate struct {
		*pageBase
	}

	// pageBeforeDelete
	//
	// This type is auto-generated.
	pageBeforeDelete struct {
		*pageBase
	}

	// pageAfterCreate
	//
	// This type is auto-generated.
	pageAfterCreate struct {
		*pageBase
	}

	// pageAfterUpdate
	//
	// This type is auto-generated.
	pageAfterUpdate struct {
		*pageBase
	}

	// pageAfterDelete
	//
	// This type is auto-generated.
	pageAfterDelete struct {
		*pageBase
	}
)

// ResourceType returns "compose:page"
//
// This function is auto-generated.
func (pageBase) ResourceType() string {
	return "compose:page"
}

// EventType on pageOnManual returns "onManual"
//
// This function is auto-generated.
func (pageOnManual) EventType() string {
	return "onManual"
}

// EventType on pageBeforeCreate returns "beforeCreate"
//
// This function is auto-generated.
func (pageBeforeCreate) EventType() string {
	return "beforeCreate"
}

// EventType on pageBeforeUpdate returns "beforeUpdate"
//
// This function is auto-generated.
func (pageBeforeUpdate) EventType() string {
	return "beforeUpdate"
}

// EventType on pageBeforeDelete returns "beforeDelete"
//
// This function is auto-generated.
func (pageBeforeDelete) EventType() string {
	return "beforeDelete"
}

// EventType on pageAfterCreate returns "afterCreate"
//
// This function is auto-generated.
func (pageAfterCreate) EventType() string {
	return "afterCreate"
}

// EventType on pageAfterUpdate returns "afterUpdate"
//
// This function is auto-generated.
func (pageAfterUpdate) EventType() string {
	return "afterUpdate"
}

// EventType on pageAfterDelete returns "afterDelete"
//
// This function is auto-generated.
func (pageAfterDelete) EventType() string {
	return "afterDelete"
}

// PageOnManual creates onManual for compose:page resource
//
// This function is auto-generated.
func PageOnManual(
	argPage *types.Page,
	argOldPage *types.Page,
	argNamespace *types.Namespace,
) *pageOnManual {
	return &pageOnManual{
		pageBase: &pageBase{
			immutable: false,
			page:      argPage,
			oldPage:   argOldPage,
			namespace: argNamespace,
		},
	}
}

// PageOnManualImmutable creates onManual for compose:page resource
//
// None of the arguments will be mutable!
//
// This function is auto-generated.
func PageOnManualImmutable(
	argPage *types.Page,
	argOldPage *types.Page,
	argNamespace *types.Namespace,
) *pageOnManual {
	return &pageOnManual{
		pageBase: &pageBase{
			immutable: true,
			page:      argPage,
			oldPage:   argOldPage,
			namespace: argNamespace,
		},
	}
}

// PageBeforeCreate creates beforeCreate for compose:page resource
//
// This function is auto-generated.
func PageBeforeCreate(
	argPage *types.Page,
	argOldPage *types.Page,
	argNamespace *types.Namespace,
) *pageBeforeCreate {
	return &pageBeforeCreate{
		pageBase: &pageBase{
			immutable: false,
			page:      argPage,
			oldPage:   argOldPage,
			namespace: argNamespace,
		},
	}
}

// PageBeforeCreateImmutable creates beforeCreate for compose:page resource
//
// None of the arguments will be mutable!
//
// This function is auto-generated.
func PageBeforeCreateImmutable(
	argPage *types.Page,
	argOldPage *types.Page,
	argNamespace *types.Namespace,
) *pageBeforeCreate {
	return &pageBeforeCreate{
		pageBase: &pageBase{
			immutable: true,
			page:      argPage,
			oldPage:   argOldPage,
			namespace: argNamespace,
		},
	}
}

// PageBeforeUpdate creates beforeUpdate for compose:page resource
//
// This function is auto-generated.
func PageBeforeUpdate(
	argPage *types.Page,
	argOldPage *types.Page,
	argNamespace *types.Namespace,
) *pageBeforeUpdate {
	return &pageBeforeUpdate{
		pageBase: &pageBase{
			immutable: false,
			page:      argPage,
			oldPage:   argOldPage,
			namespace: argNamespace,
		},
	}
}

// PageBeforeUpdateImmutable creates beforeUpdate for compose:page resource
//
// None of the arguments will be mutable!
//
// This function is auto-generated.
func PageBeforeUpdateImmutable(
	argPage *types.Page,
	argOldPage *types.Page,
	argNamespace *types.Namespace,
) *pageBeforeUpdate {
	return &pageBeforeUpdate{
		pageBase: &pageBase{
			immutable: true,
			page:      argPage,
			oldPage:   argOldPage,
			namespace: argNamespace,
		},
	}
}

// PageBeforeDelete creates beforeDelete for compose:page resource
//
// This function is auto-generated.
func PageBeforeDelete(
	argPage *types.Page,
	argOldPage *types.Page,
	argNamespace *types.Namespace,
) *pageBeforeDelete {
	return &pageBeforeDelete{
		pageBase: &pageBase{
			immutable: false,
			page:      argPage,
			oldPage:   argOldPage,
			namespace: argNamespace,
		},
	}
}

// PageBeforeDeleteImmutable creates beforeDelete for compose:page resource
//
// None of the arguments will be mutable!
//
// This function is auto-generated.
func PageBeforeDeleteImmutable(
	argPage *types.Page,
	argOldPage *types.Page,
	argNamespace *types.Namespace,
) *pageBeforeDelete {
	return &pageBeforeDelete{
		pageBase: &pageBase{
			immutable: true,
			page:      argPage,
			oldPage:   argOldPage,
			namespace: argNamespace,
		},
	}
}

// PageAfterCreate creates afterCreate for compose:page resource
//
// This function is auto-generated.
func PageAfterCreate(
	argPage *types.Page,
	argOldPage *types.Page,
	argNamespace *types.Namespace,
) *pageAfterCreate {
	return &pageAfterCreate{
		pageBase: &pageBase{
			immutable: false,
			page:      argPage,
			oldPage:   argOldPage,
			namespace: argNamespace,
		},
	}
}

// PageAfterCreateImmutable creates afterCreate for compose:page resource
//
// None of the arguments will be mutable!
//
// This function is auto-generated.
func PageAfterCreateImmutable(
	argPage *types.Page,
	argOldPage *types.Page,
	argNamespace *types.Namespace,
) *pageAfterCreate {
	return &pageAfterCreate{
		pageBase: &pageBase{
			immutable: true,
			page:      argPage,
			oldPage:   argOldPage,
			namespace: argNamespace,
		},
	}
}

// PageAfterUpdate creates afterUpdate for compose:page resource
//
// This function is auto-generated.
func PageAfterUpdate(
	argPage *types.Page,
	argOldPage *types.Page,
	argNamespace *types.Namespace,
) *pageAfterUpdate {
	return &pageAfterUpdate{
		pageBase: &pageBase{
			immutable: false,
			page:      argPage,
			oldPage:   argOldPage,
			namespace: argNamespace,
		},
	}
}

// PageAfterUpdateImmutable creates afterUpdate for compose:page resource
//
// None of the arguments will be mutable!
//
// This function is auto-generated.
func PageAfterUpdateImmutable(
	argPage *types.Page,
	argOldPage *types.Page,
	argNamespace *types.Namespace,
) *pageAfterUpdate {
	return &pageAfterUpdate{
		pageBase: &pageBase{
			immutable: true,
			page:      argPage,
			oldPage:   argOldPage,
			namespace: argNamespace,
		},
	}
}

// PageAfterDelete creates afterDelete for compose:page resource
//
// This function is auto-generated.
func PageAfterDelete(
	argPage *types.Page,
	argOldPage *types.Page,
	argNamespace *types.Namespace,
) *pageAfterDelete {
	return &pageAfterDelete{
		pageBase: &pageBase{
			immutable: false,
			page:      argPage,
			oldPage:   argOldPage,
			namespace: argNamespace,
		},
	}
}

// PageAfterDeleteImmutable creates afterDelete for compose:page resource
//
// None of the arguments will be mutable!
//
// This function is auto-generated.
func PageAfterDeleteImmutable(
	argPage *types.Page,
	argOldPage *types.Page,
	argNamespace *types.Namespace,
) *pageAfterDelete {
	return &pageAfterDelete{
		pageBase: &pageBase{
			immutable: true,
			page:      argPage,
			oldPage:   argOldPage,
			namespace: argNamespace,
		},
	}
}

// SetPage sets new page value
//
// This function is auto-generated.
func (res *pageBase) SetPage(argPage *types.Page) {
	res.page = argPage
}

// Page returns page
//
// This function is auto-generated.
func (res pageBase) Page() *types.Page {
	return res.page
}

// OldPage returns oldPage
//
// This function is auto-generated.
func (res pageBase) OldPage() *types.Page {
	return res.oldPage
}

// Namespace returns namespace
//
// This function is auto-generated.
func (res pageBase) Namespace() *types.Namespace {
	return res.namespace
}

// SetInvoker sets new invoker value
//
// This function is auto-generated.
func (res *pageBase) SetInvoker(argInvoker auth.Identifiable) {
	res.invoker = argInvoker
}

// Invoker returns invoker
//
// This function is auto-generated.
func (res pageBase) Invoker() auth.Identifiable {
	return res.invoker
}

// Encode internal data to be passed as event params & arguments to triggered Corredor script
func (res pageBase) Encode() (args map[string][]byte, err error) {
	args = make(map[string][]byte)

	if args["page"], err = json.Marshal(res.page); err != nil {
		return nil, err
	}

	if args["oldPage"], err = json.Marshal(res.oldPage); err != nil {
		return nil, err
	}

	if args["namespace"], err = json.Marshal(res.namespace); err != nil {
		return nil, err
	}

	if args["invoker"], err = json.Marshal(res.invoker); err != nil {
		return nil, err
	}

	return
}

// Decode return values from Corredor script into struct props
func (res *pageBase) Decode(results map[string][]byte) (err error) {
	if res.immutable {
		// Respect immutability
		return
	}
	if res.page != nil {
		if r, ok := results["result"]; ok && len(results) == 1 {
			if err = json.Unmarshal(r, res.page); err != nil {
				return
			}
		}
	}

	if res.page != nil {
		if r, ok := results["page"]; ok {
			if err = json.Unmarshal(r, res.page); err != nil {
				return
			}
		}
	}

	if res.invoker != nil {
		if r, ok := results["invoker"]; ok {
			if err = json.Unmarshal(r, res.invoker); err != nil {
				return
			}
		}
	}
	return
}
