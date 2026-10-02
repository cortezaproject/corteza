package compose

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/crusttech/human/server/compose/types"
	"github.com/crusttech/human/server/pkg/eventbus"
	"github.com/crusttech/human/server/tests/helpers"
)

// registerBeforeRecordHandler registers a handler that sets the given field
// on the record before it is created/updated, the way a workflow on a
// beforeCreate/beforeUpdate trigger does
func registerBeforeRecordHandler(eventType, field, value string) (unregister func()) {
	ptr := eventBus.Register(func(ctx context.Context, ev eventbus.Event) error {
		rec := ev.(interface{ Record() *types.Record }).Record()
		rec.Values = rec.Values.Set(&types.RecordValue{Name: field, Value: value})
		return nil
	}, eventbus.For("compose:record"), eventbus.On(eventType))

	return func() { eventBus.Unregister(ptr) }
}

// A before-handler (workflow, script) may change fields the invoking user
// is not allowed to update; that is the automation's action, not the user's
func TestRecordBeforeHandler_setsFieldUserCannotUpdate(t *testing.T) {
	h := newHelper(t)
	h.clearRecords()

	module := h.repoMakeRecordModuleWithFields("record before handler module",
		&types.ModuleField{Name: "name", Kind: "String"},
		&types.ModuleField{Name: "managed", Kind: "String"},
	)
	helpers.AllowMe(h, types.ModuleRbacResource(0, 0), "record.create")
	helpers.AllowMe(h, types.RecordRbacResource(0, 0, 0), "update")
	helpers.AllowMe(h, module.Fields[0].RbacResource(), "record.value.update")
	helpers.DenyMe(h, module.Fields[1].RbacResource(), "record.value.update")

	url := fmt.Sprintf("/namespace/%d/module/%d/record/", module.NamespaceID, module.ID)

	{
		// create
		unregister := registerBeforeRecordHandler("beforeCreate", "managed", "set on create")

		var created struct {
			Response *types.Record
		}

		h.apiInit().
			Post(url).
			JSON(`{"values": [{"name": "name", "value": "new"}]}`).
			Header("Accept", "application/json").
			Expect(t).
			Status(http.StatusOK).
			Assert(helpers.AssertNoErrors).
			End().
			JSON(&created)

		h.a.NotNil(created.Response)
		r := h.lookupRecordByID(module, created.Response.ID)
		h.a.Equal("set on create", r.Values.Get("managed", 0).Value)
		unregister()
	}

	unregister := registerBeforeRecordHandler("beforeUpdate", "managed", "set on update")
	defer unregister()

	{
		// update

		record := h.makeRecord(module,
			&types.RecordValue{Name: "name", Value: "old"},
			&types.RecordValue{Name: "managed", Value: "initial"},
		)

		h.apiInit().
			Post(fmt.Sprintf("%s%d", url, record.ID)).
			JSON(`{"values": [{"name": "name", "value": "changed"}, {"name": "managed", "value": "initial"}]}`).
			Header("Accept", "application/json").
			Expect(t).
			Status(http.StatusOK).
			Assert(helpers.AssertNoErrors).
			End()

		r := h.lookupRecordByID(module, record.ID)
		h.a.Equal("changed", r.Values.Get("name", 0).Value)
		h.a.Equal("set on update", r.Values.Get("managed", 0).Value)
	}

	{
		// the user still cannot change the field directly
		record := h.makeRecord(module,
			&types.RecordValue{Name: "name", Value: "old"},
			&types.RecordValue{Name: "managed", Value: "initial"},
		)

		h.apiInit().
			Post(fmt.Sprintf("%s%d", url, record.ID)).
			JSON(`{"values": [{"name": "name", "value": "changed"}, {"name": "managed", "value": "by user"}]}`).
			Header("Accept", "application/json").
			Expect(t).
			Status(http.StatusOK).
			Assert(helpers.AssertError("1 issue(s) found")).
			End()

		r := h.lookupRecordByID(module, record.ID)
		h.a.Equal("initial", r.Values.Get("managed", 0).Value)
	}
}
