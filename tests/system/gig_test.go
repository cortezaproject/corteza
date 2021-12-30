package system

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/cortezaproject/corteza-server/system/types"
	"github.com/cortezaproject/corteza-server/tests/helpers"
	jsonpath "github.com/steinfletcher/apitest-jsonpath"
)

func apiCreateGig(t *testing.T, h helper, pl string) uint64 {
	type auxGigRsp struct {
		Response struct {
			ID uint64 `json:"gigID,string"`
		}
	}

	r := h.apiInit().Post("/gig/").
		Header("Accept", "application/json").
		JSON(pl).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()

	g := auxGigRsp{}
	b, err := ioutil.ReadAll(r.Response.Body)
	h.a.NoError(err)
	defer r.Response.Body.Close()
	h.a.NoError(json.Unmarshal(b, &g))

	return g.Response.ID
}

func TestGigCreate(t *testing.T) {
	h := newHelper(t)
	h.clearReports()
	helpers.AllowMe(h, types.ComponentRbacResource(), "report.create")

	h.apiInit().
		Post("/gig/").
		Header("Accept", "application/json").
		JSON(`{
			"worker": "noop"
		}`).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		Assert(jsonpath.Equal("$.response.worker.ref", "noop")).
		End()
}

func TestGigCreate_invalidWorker(t *testing.T) {
	h := newHelper(t)
	h.clearReports()
	helpers.AllowMe(h, types.ComponentRbacResource(), "report.create")

	h.apiInit().
		Post("/gig/").
		Header("Accept", "application/json").
		JSON(`{
			"worker": "does_not_exist"
		}`).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("unknown worker: does_not_exist")).
		End()
}

func TestGigAddSources(t *testing.T) {
	h := newHelper(t)
	h.clearReports()
	helpers.AllowMe(h, types.ComponentRbacResource(), "report.create")

	gigID := apiCreateGig(t, h, `{ "worker": "noop" }`)

	body, contentType := h.apiInitRecordImport("srcA.txt", "secret-key.txt")
	h.apiInit().Patch(fmt.Sprintf("/gig/%d/sources", gigID)).
		Header("Accept", "application/json").
		Body(body).
		ContentType(contentType).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		Assert(jsonpath.Len("$.response.sources", 1)).
		End()

	body, contentType = h.apiInitRecordImport("srcB.txt", "secret-message.txt", "decoders", `[{ "ref": "noop" }]`)
	h.apiInit().Patch(fmt.Sprintf("/gig/%d/sources", gigID)).
		Header("Accept", "application/json").
		Body(body).
		ContentType(contentType).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		Assert(jsonpath.Len("$.response.sources", 2)).
		Assert(jsonpath.Len("$.response.sources[1].decoders", 1)).
		End()
}
