package messaging

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/cortezaproject/corteza-server/tests/helpers"
	jsonpath "github.com/steinfletcher/apitest-jsonpath"
)

func TestMessageSearch(t *testing.T) {
	h := newHelper(t)
	ch := h.repoMakePublicCh()

	h.repoMakeMessage("searchTestMessageA", ch, h.cUser)
	h.repoMakeMessage("searchTestMessageB", ch, h.cUser)

	h.apiInit().
		Get("/search/messages").
		Query("query", "searchTestMessageA").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		Assert(jsonpath.Len(`$.response`, 1)).
		End()
}

func TestMessageSearchAfterID(t *testing.T) {
	h := newHelper(t)
	ch := h.repoMakePublicCh()

	h.repoMakeMessage("searchTestMessageA", ch, h.cUser)
	h.repoMakeMessage("searchTestMessageB", ch, h.cUser)
	m := h.repoMakeMessage("searchTestMessageC", ch, h.cUser)
	h.repoMakeMessage("searchTestMessageD", ch, h.cUser)
	h.repoMakeMessage("searchTestMessageE", ch, h.cUser)

	h.apiInit().
		Get("/search/messages").
		Query("afterMessageID", fmt.Sprintf("%d", m.ID)).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		Assert(jsonpath.Len(`$.response`, 2)).
		Assert(jsonpath.Equal(`$.response[0].message`, "searchTestMessageD")).
		Assert(jsonpath.Equal(`$.response[1].message`, "searchTestMessageE")).
		End()
}

func TestMessageSearchFromID(t *testing.T) {
	h := newHelper(t)
	ch := h.repoMakePublicCh()

	h.repoMakeMessage("searchTestMessageA", ch, h.cUser)
	h.repoMakeMessage("searchTestMessageB", ch, h.cUser)
	m := h.repoMakeMessage("searchTestMessageC", ch, h.cUser)
	h.repoMakeMessage("searchTestMessageD", ch, h.cUser)
	h.repoMakeMessage("searchTestMessageE", ch, h.cUser)

	h.apiInit().
		Get("/search/messages").
		Query("fromMessageID", fmt.Sprintf("%d", m.ID)).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		Assert(jsonpath.Len(`$.response`, 3)).
		Assert(jsonpath.Equal(`$.response[0].message`, "searchTestMessageC")).
		Assert(jsonpath.Equal(`$.response[1].message`, "searchTestMessageD")).
		Assert(jsonpath.Equal(`$.response[2].message`, "searchTestMessageE")).
		End()
}

func TestMessageSearchBeforeID(t *testing.T) {
	h := newHelper(t)
	ch := h.repoMakePublicCh()

	h.repoMakeMessage("searchTestMessageA", ch, h.cUser)
	h.repoMakeMessage("searchTestMessageB", ch, h.cUser)
	m := h.repoMakeMessage("searchTestMessageC", ch, h.cUser)
	h.repoMakeMessage("searchTestMessageD", ch, h.cUser)
	h.repoMakeMessage("searchTestMessageE", ch, h.cUser)

	h.apiInit().
		Get("/search/messages").
		Query("beforeMessageID", fmt.Sprintf("%d", m.ID)).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		Assert(jsonpath.Equal(`$.response[0].message`, "searchTestMessageB")).
		Assert(jsonpath.Equal(`$.response[1].message`, "searchTestMessageA")).
		End()
}

func TestMessageSearchToID(t *testing.T) {
	h := newHelper(t)
	ch := h.repoMakePublicCh()

	h.repoMakeMessage("searchTestMessageA", ch, h.cUser)
	h.repoMakeMessage("searchTestMessageB", ch, h.cUser)
	m := h.repoMakeMessage("searchTestMessageC", ch, h.cUser)
	h.repoMakeMessage("searchTestMessageD", ch, h.cUser)
	h.repoMakeMessage("searchTestMessageE", ch, h.cUser)

	h.apiInit().
		Get("/search/messages").
		Query("toMessageID", fmt.Sprintf("%d", m.ID)).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		Assert(jsonpath.Equal(`$.response[0].message`, "searchTestMessageC")).
		Assert(jsonpath.Equal(`$.response[1].message`, "searchTestMessageB")).
		Assert(jsonpath.Equal(`$.response[2].message`, "searchTestMessageA")).
		End()
}

func TestMessageThreadSearch(t *testing.T) {
	h := newHelper(t)
	ch := h.repoMakePublicCh()

	msgA := h.repoMakeMessage("searchTestMessageThreadA", ch, h.cUser)
	h.apiMessageCreateReply("thrA", msgA)

	msgB := h.repoMakeMessage("searchTestMessageThreadB", ch, h.cUser)
	h.apiMessageCreateReply("thrB", msgB)

	h.apiInit().
		Get("/search/threads").
		Query("query", "searchTestMessageThreadA").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		Assert(jsonpath.Len(`$.response`, 1)).
		End()
}
