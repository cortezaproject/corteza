package federation

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/crusttech/human/server/federation/types"
	"github.com/crusttech/human/server/tests/helpers"
)

// Nodes are changed only by users allowed to manage them

func TestNodeUpdateForbidden(t *testing.T) {
	h := newHelper(t)
	h.clearNodes()
	n := h.createNode()

	h.apiInit().
		Post(fmt.Sprintf("/nodes/%d", n.ID)).
		Header("Accept", "application/json").
		FormData("name", "renamed").
		FormData("baseURL", "https://example.tld/federation").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("not allowed to manage this node")).
		End()

	h.a.Equal("node", h.lookupNodeByID(n.ID).Name)
}

func TestNodeUpdate(t *testing.T) {
	h := newHelper(t)
	h.clearNodes()
	n := h.createNode()
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
	n := h.createNode()

	h.apiInit().
		Delete(fmt.Sprintf("/nodes/%d", n.ID)).
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("not allowed to manage this node")).
		End()

	h.a.Nil(h.lookupNodeByID(n.ID).DeletedAt)
}

func TestNodeDelete(t *testing.T) {
	h := newHelper(t)
	h.clearNodes()
	n := h.createNode()
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

// A new pairing URI lets its holder pair with the node, so only users allowed
// to manage the node may generate one
func TestNodeGenerateURIForbidden(t *testing.T) {
	h := newHelper(t)
	h.clearNodes()
	n := h.createNode()
	before := h.lookupNodeByID(n.ID).PairToken

	h.apiInit().
		Post(fmt.Sprintf("/nodes/%d/uri", n.ID)).
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("not allowed to manage this node")).
		End()

	h.a.Equal(before, h.lookupNodeByID(n.ID).PairToken)
}

func TestNodeGenerateURI(t *testing.T) {
	h := newHelper(t)
	h.clearNodes()
	n := h.createNode()
	before := h.lookupNodeByID(n.ID).PairToken
	helpers.AllowMe(h, types.NodeRbacResource(0), "manage")

	h.apiInit().
		Post(fmt.Sprintf("/nodes/%d/uri", n.ID)).
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()

	h.a.NotEqual(before, h.lookupNodeByID(n.ID).PairToken)
}
