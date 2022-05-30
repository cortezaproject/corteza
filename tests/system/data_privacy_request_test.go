package system

import (
	"context"
	"fmt"
	"github.com/cortezaproject/corteza-server/pkg/id"
	"github.com/cortezaproject/corteza-server/store"
	"github.com/cortezaproject/corteza-server/system/service"
	"github.com/cortezaproject/corteza-server/system/types"
	"github.com/cortezaproject/corteza-server/tests/helpers"
	jsonpath "github.com/steinfletcher/apitest-jsonpath"
	"net/http"
	"testing"
	"time"
)

func (h helper) clearDataPrivacyRequests() {
	h.noError(store.TruncateDataPrivacyRequests(context.Background(), service.DefaultStore))
}

func (h helper) createDataPrivacyRequest(kind types.RequestKind, status types.RequestStatus) *types.DataPrivacyRequest {
	res := &types.DataPrivacyRequest{
		ID:          id.Next(),
		Kind:        types.RequestKindCorrect,
		Status:      types.RequestStatusPending,
		RequestedAt: time.Now(),
		RequestedBy: 0,
		CreatedAt:   time.Now(),
	}

	if len(kind) > 0 {
		res.Kind = kind
	}

	if len(status) > 0 {
		res.Status = status
	}

	h.a.NoError(store.CreateDataPrivacyRequest(context.Background(), service.DefaultStore, res))

	return res
}

func (h helper) createSampleDataPrivacyRequest() *types.DataPrivacyRequest {
	return h.createDataPrivacyRequest("", "")
}

func TestDataPrivacyRequestList(t *testing.T) {
	h := newHelper(t)
	h.clearDataPrivacyRequests()

	h.createSampleDataPrivacyRequest()
	h.createSampleDataPrivacyRequest()

	helpers.AllowMe(h, types.ComponentRbacResource(), "data-privacy-requests.search")
	helpers.AllowMe(h, types.DataPrivacyRequestRbacResource(0), "read")

	h.apiInit().
		Get("/data-privacy/requests/").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		Assert(jsonpath.Len(`$.response.set`, 2)).
		End()
}

func TestDataPrivacyRequestListWithPaging(t *testing.T) {
	h := newHelper(t)
	h.clearDataPrivacyRequests()

	seedCount := 40
	for i := 0; i < seedCount; i++ {
		h.createSampleDataPrivacyRequest()
	}

	helpers.AllowMe(h, types.ComponentRbacResource(), "data-privacy-requests.search")
	helpers.AllowMe(h, types.DataPrivacyRequestRbacResource(0), "read")

	var aux = struct {
		Response struct {
			Filter struct {
				NextPage *string
				PrevPage *string
			}
		}
	}{}

	h.apiInit().
		Get("/data-privacy/requests/").
		Query("limit", "10").
		Query("sort", "kind,createdAt+DESC").
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		Assert(jsonpath.Present(`$.response.filter`)).
		Assert(jsonpath.Present(`$.response.set`)).
		Assert(jsonpath.Len(`$.response.set`, 10)).
		Assert(jsonpath.Present(`$.response.filter.nextPage`)).
		End().
		JSON(&aux)

	h.a.NotNil(aux.Response.Filter.NextPage)

	h.apiInit().
		Get("/data-privacy/requests/").
		Header("Accept", "application/json").
		Query("limit", "10").
		Query("sort", "kind,createdAt+DESC").
		Query("pageCursor", *aux.Response.Filter.NextPage).
		Expect(t).
		Status(http.StatusOK).
		Assert(jsonpath.Len(`$.response.set`, 10)).
		Assert(jsonpath.Present(`$.response.filter.prevPage`)).
		Assert(jsonpath.Present(`$.response.filter.nextPage`)).
		End().
		JSON(&aux)

	h.a.NotNil(aux.Response.Filter.PrevPage)
	h.a.NotNil(aux.Response.Filter.NextPage)
}

func TestDataPrivacyRequestListFilters(t *testing.T) {
	h := newHelper(t)
	h.clearDataPrivacyRequests()

	h.createSampleDataPrivacyRequest()
	h.createSampleDataPrivacyRequest()
	h.createDataPrivacyRequest("", types.RequestStatusApproved)
	h.createDataPrivacyRequest(types.RequestKindExport, types.RequestStatusApproved)

	helpers.AllowMe(h, types.ComponentRbacResource(), "data-privacy-requests.search")
	helpers.AllowMe(h, types.DataPrivacyRequestRbacResource(0), "read")

	h.apiInit().
		Get("/data-privacy/requests/").
		Query("query", types.RequestStatusApproved.String()).
		Query("kind", types.RequestKindExport.String()).
		Query("incTotal", "true").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		Assert(jsonpath.Len(`$.response.set`, 1)).
		End()
}

func TestDataPrivacyRequestRead(t *testing.T) {
	h := newHelper(t)
	h.clearDataPrivacyRequests()

	client := h.createSampleDataPrivacyRequest()

	helpers.AllowMe(h, types.DataPrivacyRequestRbacResource(0), "read")

	h.apiInit().
		Get(fmt.Sprintf("/data-privacy/requests/%d", client.ID)).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()
}

func TestDataPrivacyRequestCreate(t *testing.T) {
	h := newHelper(t)
	h.clearDataPrivacyRequests()

	helpers.AllowMe(h, types.ComponentRbacResource(), "data-privacy-request.create")

	h.apiInit().
		Post("/data-privacy/requests/").
		Header("Accept", "application/json").
		FormData("kind", "correct").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		Assert(jsonpath.Present(`$.response.requestID`)).
		Assert(jsonpath.Equal(`$.response.kind`, types.RequestKindCorrect.String())).
		Assert(jsonpath.Equal(`$.response.status`, types.RequestStatusPending.String())).
		End()
}

func TestDataPrivacyRequestUpdateStatus(t *testing.T) {
	h := newHelper(t)
	h.clearDataPrivacyRequests()

	req := h.createSampleDataPrivacyRequest()

	helpers.AllowMe(h, types.DataPrivacyRequestRbacResource(0), "approve")

	h.apiInit().
		Put(fmt.Sprintf("/data-privacy/requests/%d/status/approved", req.ID)).
		Header("Accept", "application/json").
		JSON(helpers.JSON(req)).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		Assert(jsonpath.Equal(`$.response.status`, types.RequestStatusApproved.String())).
		End()
}

func TestDataPrivacyRequestUpdateStatusForbidden(t *testing.T) {
	h := newHelper(t)
	h.clearDataPrivacyRequests()

	req := h.createSampleDataPrivacyRequest()

	h.apiInit().
		Put(fmt.Sprintf("/data-privacy/requests/%d/status/approved", req.ID)).
		Header("Accept", "application/json").
		JSON(helpers.JSON(req)).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("dataPrivacy.errors.notAllowedToApprove")).
		End()
}
