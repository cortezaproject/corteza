package system

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/cortezaproject/corteza-server/system/types"
	"github.com/cortezaproject/corteza-server/tests/helpers"
	jsonpath "github.com/steinfletcher/apitest-jsonpath"
)

func TestGigCreate(t *testing.T) {
	h := newHelper(t)
	h.clearReports()
	helpers.AllowMe(h, types.ComponentRbacResource(), "gig.create")

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

func TestGigCreateForbiden(t *testing.T) {
	h := newHelper(t)
	h.clearReports()

	h.apiInit().
		Post("/gig/").
		Header("Accept", "application/json").
		JSON(`{
			"worker": "noop"
		}`).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("gigService.errors.notAllowedToCreate")).
		End()
}

func TestGigCreate_invalidWorker(t *testing.T) {
	h := newHelper(t)
	h.clearReports()
	helpers.AllowMe(h, types.ComponentRbacResource(), "gig.create")

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

func TestGigRead(t *testing.T) {
	h := newHelper(t)
	h.clearReports()
	helpers.AllowMe(h, types.ComponentRbacResource(), "gig.create")
	gigID := apiGigCreate(t, h, `{ "worker": "noop" }`)
	helpers.AllowMe(h, types.GigRbacResource(gigID), "read")

	h.apiInit().
		Get(fmt.Sprintf("/gig/%d", gigID)).
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		Assert(jsonpath.Equal("$.response.worker.ref", "noop")).
		End()
}

func TestGigReadForbiden(t *testing.T) {
	h := newHelper(t)
	h.clearReports()
	helpers.AllowMe(h, types.ComponentRbacResource(), "gig.create")
	gigID := apiGigCreate(t, h, `{ "worker": "noop" }`)

	h.apiInit().
		Get(fmt.Sprintf("/gig/%d", gigID)).
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("gigService.errors.notAllowedToRead")).
		End()
}

func TestGigUpdate(t *testing.T) {
	h := newHelper(t)
	h.clearReports()
	helpers.AllowMe(h, types.ComponentRbacResource(), "gig.create")
	gigID := apiGigCreate(t, h, `{ "worker": "noop" }`)
	helpers.AllowMe(h, types.GigRbacResource(gigID), "update")

	h.apiInit().
		Put(fmt.Sprintf("/gig/%d", gigID)).
		Header("Accept", "application/json").
		JSON(`{ "postprocessors": [{ "ref": "noop" }] }`).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		Assert(jsonpath.Len("$.response.postprocess", 1)).
		End()
}

func TestGigUpdateForbiden(t *testing.T) {
	h := newHelper(t)
	h.clearReports()
	helpers.AllowMe(h, types.ComponentRbacResource(), "gig.create")
	gigID := apiGigCreate(t, h, `{ "worker": "noop" }`)

	h.apiInit().
		Put(fmt.Sprintf("/gig/%d", gigID)).
		Header("Accept", "application/json").
		JSON(`{ "postprocessors": [{ "ref": "noop" }] }`).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("gigService.errors.notAllowedToUpdate")).
		End()
}

func TestGigDelete(t *testing.T) {
	h := newHelper(t)
	h.clearReports()
	helpers.AllowMe(h, types.ComponentRbacResource(), "gig.create")
	gigID := apiGigCreate(t, h, `{ "worker": "noop" }`)
	helpers.AllowMe(h, types.GigRbacResource(gigID), "delete")

	h.apiInit().
		Delete(fmt.Sprintf("/gig/%d", gigID)).
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()
}

func TestGigDeleteForbiden(t *testing.T) {
	h := newHelper(t)
	h.clearReports()
	helpers.AllowMe(h, types.ComponentRbacResource(), "gig.create")
	gigID := apiGigCreate(t, h, `{ "worker": "noop" }`)

	h.apiInit().
		Delete(fmt.Sprintf("/gig/%d", gigID)).
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("gigService.errors.notAllowedToDelete")).
		End()
}

func TestGigUndelete(t *testing.T) {
	h := newHelper(t)
	h.clearReports()
	helpers.AllowMe(h, types.ComponentRbacResource(), "gig.create")
	gigID := apiGigCreate(t, h, `{ "worker": "noop" }`)
	helpers.AllowMe(h, types.GigRbacResource(gigID), "undelete")

	h.apiInit().
		Post(fmt.Sprintf("/gig/%d/undelete", gigID)).
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()
}

func TestGigUndeleteForbiden(t *testing.T) {
	h := newHelper(t)
	h.clearReports()
	helpers.AllowMe(h, types.ComponentRbacResource(), "gig.create")
	gigID := apiGigCreate(t, h, `{ "worker": "noop" }`)

	h.apiInit().
		Post(fmt.Sprintf("/gig/%d/undelete", gigID)).
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("gigService.errors.notAllowedToUndelete")).
		End()
}

func TestGigAddSource(t *testing.T) {
	h := newHelper(t)
	h.clearReports()
	helpers.AllowMe(h, types.ComponentRbacResource(), "gig.create")
	gigID := apiGigCreate(t, h, `{ "worker": "noop" }`)
	helpers.AllowMe(h, types.GigRbacResource(gigID), "update")

	body, contentType := h.apiInitGigSource("srcA.txt", "secret-key.txt")
	h.apiInit().Patch(fmt.Sprintf("/gig/%d/sources", gigID)).
		Header("Accept", "application/json").
		Body(body).
		ContentType(contentType).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		Assert(jsonpath.Len("$.response.sources", 1)).
		End()
}

func TestGigAddSourceForbiden(t *testing.T) {
	h := newHelper(t)
	h.clearReports()
	helpers.AllowMe(h, types.ComponentRbacResource(), "gig.create")
	gigID := apiGigCreate(t, h, `{ "worker": "noop" }`)

	body, contentType := h.apiInitGigSource("srcA.txt", "secret-key.txt")
	h.apiInit().Patch(fmt.Sprintf("/gig/%d/sources", gigID)).
		Header("Accept", "application/json").
		Body(body).
		ContentType(contentType).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("gigService.errors.notAllowedToUpdate")).
		End()
}

func TestGigRemoveSource(t *testing.T) {
	h := newHelper(t)
	h.clearReports()
	helpers.AllowMe(h, types.ComponentRbacResource(), "gig.create")
	gigID := apiGigCreate(t, h, `{ "worker": "noop" }`)
	helpers.AllowMe(h, types.GigRbacResource(gigID), "update")
	sourceID := apiGigAddSource(t, h, gigID, "srcA.txt", "secret-key.txt")

	h.apiInit().
		Delete(fmt.Sprintf("/gig/%d/sources/%d", gigID, sourceID)).
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()
}

func TestGigRemoveSourceForbiden(t *testing.T) {
	h := newHelper(t)
	h.clearReports()
	helpers.AllowMe(h, types.ComponentRbacResource(), "gig.create")
	gigID := apiGigCreate(t, h, `{ "worker": "noop" }`)
	helpers.AllowMe(h, types.GigRbacResource(gigID), "update")
	sourceID := apiGigAddSource(t, h, gigID, "srcA.txt", "secret-key.txt")
	helpers.DenyMe(h, types.GigRbacResource(gigID), "update")

	h.apiInit().
		Delete(fmt.Sprintf("/gig/%d/sources/%d", gigID, sourceID)).
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("gigService.errors.notAllowedToUpdate")).
		End()
}

func TestGigOutput(t *testing.T) {
	h := newHelper(t)
	h.clearReports()
	helpers.AllowMe(h, types.ComponentRbacResource(), "gig.create")
	gigID := apiGigCreate(t, h, `{ "worker": "noop" }`)
	helpers.AllowMe(h, types.GigRbacResource(gigID), "read")

	h.apiInit().
		Get(fmt.Sprintf("/gig/%d/output", gigID)).
		Header("Accept", "application/json").
		Expect(t).
		End()
}

func TestGigOutputForbiden(t *testing.T) {
	h := newHelper(t)
	h.clearReports()
	helpers.AllowMe(h, types.ComponentRbacResource(), "gig.create")
	gigID := apiGigCreate(t, h, `{ "worker": "noop" }`)

	h.apiInit().
		Get(fmt.Sprintf("/gig/%d/output", gigID)).
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("gigService.errors.notAllowedToRead")).
		End()
}

func TestGigState(t *testing.T) {
	h := newHelper(t)
	h.clearReports()
	helpers.AllowMe(h, types.ComponentRbacResource(), "gig.create")
	gigID := apiGigCreate(t, h, `{ "worker": "noop" }`)
	helpers.AllowMe(h, types.GigRbacResource(gigID), "read")

	h.apiInit().
		Get(fmt.Sprintf("/gig/%d/state", gigID)).
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()
}

func TestGigStateForbiden(t *testing.T) {
	h := newHelper(t)
	h.clearReports()
	helpers.AllowMe(h, types.ComponentRbacResource(), "gig.create")
	gigID := apiGigCreate(t, h, `{ "worker": "noop" }`)

	h.apiInit().
		Get(fmt.Sprintf("/gig/%d/state", gigID)).
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("gigService.errors.notAllowedToRead")).
		End()
}

func TestGigPrepare(t *testing.T) {
	h := newHelper(t)
	h.clearReports()
	helpers.AllowMe(h, types.ComponentRbacResource(), "gig.create")
	gigID := apiGigCreate(t, h, `{ "worker": "noop" }`)
	helpers.AllowMe(h, types.GigRbacResource(gigID), "exec")

	h.apiInit().
		Put(fmt.Sprintf("/gig/%d/prepare", gigID)).
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()
}

func TestGigPrepareForbiden(t *testing.T) {
	h := newHelper(t)
	h.clearReports()
	helpers.AllowMe(h, types.ComponentRbacResource(), "gig.create")
	gigID := apiGigCreate(t, h, `{ "worker": "noop" }`)

	h.apiInit().
		Put(fmt.Sprintf("/gig/%d/prepare", gigID)).
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("gigService.errors.notAllowedToExec")).
		End()
}

func TestGigExec(t *testing.T) {
	h := newHelper(t)
	h.clearReports()
	helpers.AllowMe(h, types.ComponentRbacResource(), "gig.create")
	gigID := apiGigCreate(t, h, `{ "worker": "noop" }`)
	helpers.AllowMe(h, types.GigRbacResource(gigID), "exec")

	h.apiInit().
		Put(fmt.Sprintf("/gig/%d/exec", gigID)).
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()
}

func TestGigExecForbiden(t *testing.T) {
	h := newHelper(t)
	h.clearReports()
	helpers.AllowMe(h, types.ComponentRbacResource(), "gig.create")
	gigID := apiGigCreate(t, h, `{ "worker": "noop" }`)

	h.apiInit().
		Put(fmt.Sprintf("/gig/%d/exec", gigID)).
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("gigService.errors.notAllowedToExec")).
		End()
}

func TestGigComplete(t *testing.T) {
	h := newHelper(t)
	h.clearReports()
	helpers.AllowMe(h, types.ComponentRbacResource(), "gig.create")
	gigID := apiGigCreate(t, h, `{ "worker": "noop" }`)
	helpers.AllowMe(h, types.GigRbacResource(gigID), "update")

	h.apiInit().
		Patch(fmt.Sprintf("/gig/%d/complete", gigID)).
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()
}

func TestGigCompleteForbiden(t *testing.T) {
	h := newHelper(t)
	h.clearReports()
	helpers.AllowMe(h, types.ComponentRbacResource(), "gig.create")
	gigID := apiGigCreate(t, h, `{ "worker": "noop" }`)

	h.apiInit().
		Patch(fmt.Sprintf("/gig/%d/complete", gigID)).
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("gigService.errors.notAllowedToUpdate")).
		End()
}
