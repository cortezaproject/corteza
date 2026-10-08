package compose

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/cortezaproject/corteza/server/compose/types"
	"github.com/cortezaproject/corteza/server/tests/helpers"
)

// Reports aggregate record values, fields that can not be read must not be used in them
func TestRecordReportUnreadableField(t *testing.T) {
	h := newHelper(t)
	h.clearRecords()

	module := h.repoMakeRecordModuleWithFields(
		"record testing module",
		&types.ModuleField{Name: "public", Kind: "Number"},
		&types.ModuleField{Name: "salary", Kind: "Number"},
	)

	helpers.AllowMe(h, types.NamespaceRbacResource(0), "read")
	helpers.AllowMe(h, types.ModuleRbacResource(0, 0), "read", "records.search")
	helpers.AllowMe(h, types.RecordRbacResource(0, 0, 0), "read")
	helpers.DenyMe(h, types.ModuleFieldRbacResource(0, 0, module.Fields.FindByName("salary").ID), "record.value.read")

	for name, q := range map[string][3]string{
		"metric":      {"SUM(salary) AS total", "created_at", ""},
		"metric-case": {"SUM(SALARY) AS total", "created_at", ""},
		"filter-case": {"COUNT(ID) AS total", "created_at", "Salary > 1000"},
		"dimension":   {"COUNT(ID) AS total", "salary", ""},
		"filter":      {"COUNT(ID) AS total", "created_at", "salary > 1000"},
	} {
		t.Run(name, func(t *testing.T) {
			h.apiInit().
				Get(fmt.Sprintf("/namespace/%d/module/%d/record/report", module.NamespaceID, module.ID)).
				Query("metrics", q[0]).
				Query("dimensions", q[1]).
				Query("filter", q[2]).
				Header("Accept", "application/json").
				Expect(t).
				Status(http.StatusOK).
				Assert(helpers.AssertError("record.errors.notAllowedToFilterByField")).
				End()
		})
	}
}

// A readable field keeps working in reports; only unreadable ones are refused
func TestRecordReportReadableField(t *testing.T) {
	h := newHelper(t)
	h.clearRecords()

	module := h.repoMakeRecordModuleWithFields(
		"record report readable module",
		&types.ModuleField{Name: "public", Kind: "Number"},
		&types.ModuleField{Name: "salary", Kind: "Number"},
	)

	helpers.AllowMe(h, types.NamespaceRbacResource(0), "read")
	helpers.AllowMe(h, types.ModuleRbacResource(0, 0), "read", "records.search")
	helpers.AllowMe(h, types.RecordRbacResource(0, 0, 0), "read")
	helpers.DenyMe(h, types.ModuleFieldRbacResource(0, 0, module.Fields.FindByName("salary").ID), "record.value.read")

	h.apiInit().
		Get(fmt.Sprintf("/namespace/%d/module/%d/record/report", module.NamespaceID, module.ID)).
		Query("metrics", "SUM(public) AS total").
		Query("dimensions", "created_at").
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()
}

// A report adds up only the records the caller may read: whoever reads some of
// a module's records gets totals over those, not over the whole module.
func TestRecordReportReadableRecords(t *testing.T) {
	h := newHelper(t)
	h.clearRecords()

	module := h.repoMakeRecordModuleWithFields(
		"record report readable records",
		&types.ModuleField{Name: "g", Kind: "String"},
		&types.ModuleField{Name: "n", Kind: "Number"},
	)
	helpers.AllowMe(h, types.ModuleRbacResource(0, 0), "records.search")

	var ids []uint64
	for _, n := range []string{"1", "2", "4"} {
		r := h.makeRecord(module, &types.RecordValue{Name: "g", Value: "x"}, &types.RecordValue{Name: "n", Value: n})
		ids = append(ids, r.ID)
	}

	report := func(t *testing.T) (count, total float64) {
		var out struct {
			Response []map[string]any `json:"response"`
		}
		h.apiInit().
			Get(fmt.Sprintf("/namespace/%d/module/%d/record/report", module.NamespaceID, module.ID)).
			Query("metrics", "SUM(n) AS total").
			Query("dimensions", "g").
			Header("Accept", "application/json").
			Expect(t).
			Status(http.StatusOK).
			Assert(helpers.AssertNoErrors).
			End().
			JSON(&out)
		for _, row := range out.Response {
			c, _ := row["count"].(float64)
			s, _ := row["total"].(float64)
			count, total = count+c, total+s
		}
		return
	}

	t.Run("reading every record, all of them", func(t *testing.T) {
		count, total := report(t)
		h.a.Equal(3.0, count)
		h.a.Equal(7.0, total)
	})

	t.Run("one record denied, the rest", func(t *testing.T) {
		helpers.DenyMe(h, types.RecordRbacResource(module.NamespaceID, module.ID, ids[2]), "read")
		count, total := report(t)
		h.a.Equal(2.0, count)
		h.a.Equal(3.0, total)
	})

	t.Run("reading one record, that one", func(t *testing.T) {
		helpers.DenyMe(h, types.RecordRbacResource(0, 0, 0), "read")
		helpers.AllowMe(h, types.RecordRbacResource(module.NamespaceID, module.ID, ids[1]), "read")
		count, total := report(t)
		h.a.Equal(1.0, count)
		h.a.Equal(2.0, total)
	})

	t.Run("reading none, nothing", func(t *testing.T) {
		helpers.DenyMe(h, types.RecordRbacResource(module.NamespaceID, module.ID, ids[1]), "read")
		count, _ := report(t)
		h.a.Equal(0.0, count)
	})
}
