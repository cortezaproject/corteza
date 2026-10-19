package compose

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/cortezaproject/corteza/server/compose/types"
	"github.com/cortezaproject/corteza/server/tests/helpers"
	"github.com/steinfletcher/apitest"
)

// Filtering or sorting by a field the caller can not read is refused however
// the field name is capitalised; the DAL resolves identifiers
// case-insensitively, so a case variant would otherwise filter or order the
// records by the hidden value and reveal it.
func TestRecordListUnreadableFieldFilterSort(t *testing.T) {
	h := newHelper(t)
	h.clearRecords()

	module := h.repoMakeRecordModuleWithFields(
		"record list access module",
		&types.ModuleField{Name: "Title", Kind: "String"},
		&types.ModuleField{Name: "Secret", Kind: "String"},
	)

	helpers.AllowMe(h, types.NamespaceRbacResource(0), "read")
	helpers.AllowMe(h, types.ModuleRbacResource(0, 0), "read", "records.search")
	helpers.AllowMe(h, types.RecordRbacResource(0, 0, 0), "read")
	helpers.DenyMe(h, types.ModuleFieldRbacResource(0, 0, module.Fields.FindByName("Secret").ID), "record.value.read")

	h.makeRecord(module, &types.RecordValue{Name: "Title", Value: "Alpha"}, &types.RecordValue{Name: "Secret", Value: "vqhl-om60"})

	list := func(t *testing.T, query, sort string) *apitest.Response {
		return h.apiInit().
			Get(fmt.Sprintf("/namespace/%d/module/%d/record/", module.NamespaceID, module.ID)).
			Query("query", query).
			Query("sort", sort).
			Header("Accept", "application/json").
			Expect(t).
			Status(http.StatusOK)
	}

	for name, q := range map[string]string{
		"exact":      "Secret = 'vqhl-om60'",
		"lower":      "secret = 'vqhl-om60'",
		"upper":      "SECRET LIKE 'v%'",
		"and-prefix": "Title = 'Alpha' AND secret LIKE 'vq%'",
	} {
		t.Run("filter "+name, func(t *testing.T) {
			list(t, q, "").Assert(helpers.AssertError("record.errors.notAllowedToFilterByField")).End()
		})
	}

	for name, s := range map[string]string{
		"exact": "Secret",
		"lower": "secret",
		"desc":  "Title, secret DESC",
	} {
		t.Run("sort "+name, func(t *testing.T) {
			list(t, "", s).Assert(helpers.AssertError("record.errors.notAllowedToSortByField")).End()
		})
	}

	t.Run("readable field", func(t *testing.T) {
		list(t, "title = 'Alpha'", "TITLE DESC").Assert(helpers.AssertNoErrors).End()
	})
}
