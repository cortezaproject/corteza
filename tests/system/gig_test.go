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

func apiGigCreate(t *testing.T, h helper, pl string) uint64 {
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

func apiGigAddSource(t *testing.T, h helper, gigID uint64, path, name string, meta ...string) uint64 {
	type auxGigRsp struct {
		Response struct {
			ID      uint64 `json:"gigID,string"`
			Sources []struct {
				ID uint64 `json:"sourceID,string"`
			} `json:"sources"`
		}
	}

	body, contentType := h.apiInitGigSource(path, name, meta...)
	r := h.apiInit().Patch(fmt.Sprintf("/gig/%d/sources", gigID)).
		Header("Accept", "application/json").
		Body(body).
		ContentType(contentType).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()

	g := auxGigRsp{}
	b, err := ioutil.ReadAll(r.Response.Body)
	h.a.NoError(err)
	defer r.Response.Body.Close()
	h.a.NoError(json.Unmarshal(b, &g))

	for _, s := range g.Response.Sources {
		return s.ID
	}

	h.a.FailNow("no sources in the response")
	return 0
}

func TestGigAddSources(t *testing.T) {
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

	body, contentType = h.apiInitGigSource("srcB.txt", "secret-message.txt", "decoders", `[{ "ref": "noop" }]`)
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
