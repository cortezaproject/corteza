package system

import (
	"context"
	"fmt"
	"github.com/cortezaproject/corteza-server/pkg/id"
	"github.com/cortezaproject/corteza-server/store"
	"github.com/cortezaproject/corteza-server/system/service"
	"github.com/cortezaproject/corteza-server/system/types"
	"github.com/cortezaproject/corteza-server/tests/helpers"
	"github.com/spf13/cast"
	jsonpath "github.com/steinfletcher/apitest-jsonpath"
	"net/http"
	"testing"
	"time"
)

func (h helper) clearDataPrivacyRequestComments() {
	h.noError(store.TruncateDataPrivacyRequestComments(context.Background(), service.DefaultStore))
}

func (h helper) createDataPrivacyRequestComment(comment string, requestID uint64) *types.DataPrivacyRequestComment {
	res := &types.DataPrivacyRequestComment{
		ID:        id.Next(),
		RequestID: requestID,
		Comment:   comment,
		CreatedAt: time.Now(),
	}

	h.a.NoError(store.CreateDataPrivacyRequestComment(context.Background(), service.DefaultStore, res))

	return res
}

func (h helper) createSampleDataPrivacyRequestComment(requestID ...uint64) *types.DataPrivacyRequestComment {
	rID := id.Next()
	if len(requestID) == 1 {
		rID = requestID[0]
	}
	return h.createDataPrivacyRequestComment(rs(20), rID)
}

func TestDataPrivacyRequestCommentList(t *testing.T) {
	h := newHelper(t)
	h.clearDataPrivacyRequests()
	h.clearDataPrivacyRequestComments()

	request := h.createSampleDataPrivacyRequest()
	reqID := request.ID
	h.createSampleDataPrivacyRequestComment(reqID)
	h.createSampleDataPrivacyRequestComment(reqID)
	h.createSampleDataPrivacyRequestComment()
	h.createSampleDataPrivacyRequestComment()

	h.apiInit().
		Get(fmt.Sprintf("/data-privacy/requests/%d/comments/", reqID)).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		Assert(jsonpath.Len(`$.response.set`, 2)).
		End()
}

func TestDataPrivacyRequestCreateComment(t *testing.T) {
	h := newHelper(t)
	h.clearDataPrivacyRequestComments()

	request := h.createSampleDataPrivacyRequest()
	reqID := request.ID
	comment := rs(20)

	helpers.AllowMe(h, types.DataPrivacyRequestRbacResource(0), "read")

	h.apiInit().
		Post(fmt.Sprintf("/data-privacy/requests/%d/comments/", reqID)).
		Header("Accept", "application/json").
		FormData("comment", comment).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		Assert(jsonpath.Present(`$.response.requestID`)).
		Assert(jsonpath.Present(`$.response.createdAt`)).
		Assert(jsonpath.Present(`$.response.createdBy`)).
		Assert(jsonpath.Equal(`$.response.requestID`, cast.ToString(reqID))).
		Assert(jsonpath.Equal(`$.response.comment`, comment)).
		End()
}

func TestDataPrivacyRequestCreateCommentForbidden(t *testing.T) {
	h := newHelper(t)
	h.clearDataPrivacyRequestComments()
	comment := rs(20)

	h.apiInit().
		Post(fmt.Sprintf("/data-privacy/requests/%d/comments/", id.Next())).
		Header("Accept", "application/json").
		FormData("comment", comment).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("dataPrivacy.errors.notFound")).
		End()
}
