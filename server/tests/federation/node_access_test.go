package federation

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/cortezaproject/corteza/server/federation/service"
	"github.com/cortezaproject/corteza/server/federation/types"
	"github.com/cortezaproject/corteza/server/pkg/id"
	"github.com/cortezaproject/corteza/server/store"
	"github.com/cortezaproject/corteza/server/tests/helpers"
)

func (h helper) repoMakeNode(pairToken string) *types.Node {
	n := &types.Node{
		ID:        id.Next(),
		Name:      "node",
		BaseURL:   "https://example.tld/federation",
		Status:    types.NodeStatusPending,
		PairToken: pairToken,
		CreatedAt: time.Now(),
	}

	h.noError(store.CreateFederationNode(context.Background(), service.DefaultStore, n))
	return n
}

// Nodes can only be changed by users that are allowed to manage them

func TestNodeUpdateForbidden(t *testing.T) {
	h := newHelper(t)
	h.clearNodes()
	n := h.repoMakeNode("")

	h.apiInit().
		Post(fmt.Sprintf("/nodes/%d", n.ID)).
		Header("Accept", "application/json").
		FormData("name", "hijacked").
		FormData("baseURL", "https://evil.tld/federation").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("not allowed to manage this node")).
		End()

	h.a.Equal("https://example.tld/federation", h.lookupNodeByID(n.ID).BaseURL)
}

func TestNodeUpdate(t *testing.T) {
	h := newHelper(t)
	h.clearNodes()
	n := h.repoMakeNode("")
	helpers.AllowMe(h, types.NodeRbacResource(0), "manage")

	h.apiInit().
		Post(fmt.Sprintf("/nodes/%d", n.ID)).
		Header("Accept", "application/json").
		FormData("name", "renamed").
		FormData("baseURL", "https://example.tld/federation").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()

	h.a.Equal("renamed", h.lookupNodeByID(n.ID).Name)
}

func TestNodeDeleteForbidden(t *testing.T) {
	h := newHelper(t)
	h.clearNodes()
	n := h.repoMakeNode("")

	h.apiInit().
		Delete(fmt.Sprintf("/nodes/%d", n.ID)).
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("not allowed to manage this node")).
		End()

	h.a.Nil(h.lookupNodeByID(n.ID).DeletedAt)
}

func TestNodeUndeleteForbidden(t *testing.T) {
	h := newHelper(t)
	h.clearNodes()
	n := h.repoMakeNode("")

	h.apiInit().
		Post(fmt.Sprintf("/nodes/%d/undelete", n.ID)).
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("not allowed to manage this node")).
		End()
}

// Regenerating the URI exposes the new pairing token
func TestNodeGenerateURIForbidden(t *testing.T) {
	h := newHelper(t)
	h.clearNodes()
	n := h.repoMakeNode("original-pair-token-original-1234")

	h.apiInit().
		Post(fmt.Sprintf("/nodes/%d/uri", n.ID)).
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("not allowed to manage this node")).
		End()

	h.a.Equal("original-pair-token-original-1234", h.lookupNodeByID(n.ID).PairToken)
}

// Nodes without a generated pair token must not accept handshakes with an empty token
func TestNodeHandshakeInitEmptyPairToken(t *testing.T) {
	h := newHelper(t)
	h.clearNodes()
	n := h.repoMakeNode("")

	err := service.DefaultNode.HandshakeInit(context.Background(), n.ID, "", n.ID, "attacker-auth-token")
	h.a.Error(err)

	stored := h.lookupNodeByID(n.ID)
	h.a.Empty(stored.AuthToken)
	h.a.Equal(types.NodeStatusPending, stored.Status)
}

// Pairing with the right token still works and records the auth token
func TestNodeHandshakeInitValidPairToken(t *testing.T) {
	h := newHelper(t)
	h.clearNodes()
	token := "pair-token-0123456789-abcdefghij-xyz"
	n := h.repoMakeNode(token)

	h.noError(service.DefaultNode.HandshakeInit(context.Background(), n.ID, token, n.ID, "auth-token"))

	stored := h.lookupNodeByID(n.ID)
	h.a.Equal("auth-token", stored.AuthToken)
	h.a.Equal(types.NodeStatusPairRequested, stored.Status)
}

func TestNodeDelete(t *testing.T) {
	h := newHelper(t)
	h.clearNodes()
	n := h.repoMakeNode("")
	helpers.AllowMe(h, types.NodeRbacResource(0), "manage")

	h.apiInit().
		Delete(fmt.Sprintf("/nodes/%d", n.ID)).
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()

	h.a.NotNil(h.lookupNodeByID(n.ID).DeletedAt)
}

func TestNodeGenerateURI(t *testing.T) {
	h := newHelper(t)
	h.clearNodes()
	n := h.repoMakeNode("")
	helpers.AllowMe(h, types.NodeRbacResource(0), "manage")

	h.apiInit().
		Post(fmt.Sprintf("/nodes/%d/uri", n.ID)).
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()
}
