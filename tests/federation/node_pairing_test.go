package federation

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"testing"

	"github.com/cortezaproject/corteza-server/federation/service"
	"github.com/cortezaproject/corteza-server/federation/types"
	"github.com/cortezaproject/corteza-server/store"
	"github.com/cortezaproject/corteza-server/tests/helpers"
	jsonpath "github.com/steinfletcher/apitest-jsonpath"
)

type (
	// node handshake mocker, allows custom functions
	mockNodeHandshake struct {
		init     func(context.Context, *types.Node, string) error
		complete func(context.Context, *types.Node, string) error
	}
)

func (h mockNodeHandshake) Init(ctx context.Context, n *types.Node, t string) error {
	return h.init(ctx, n, t)
}

func (h mockNodeHandshake) Complete(ctx context.Context, n *types.Node, t string) error {
	return h.complete(ctx, n, t)
}

func (h helper) clearNodes() {
	h.noError(store.TruncateFederationNodes(context.Background(), service.DefaultStore))
}

func (h helper) lookupNodeByID(ID uint64) *types.Node {
	n, err := service.DefaultStore.LookupFederationNodeByID(context.Background(), ID)
	h.a.NoError(err)
	return n
}

func TestSuccessfulNodePairing(t *testing.T) {
	var (
		h        = newHelper(t)
		aNodeID  uint64
		aNodeURI string
		bNodeID  uint64

		rspWithNode struct{ Response *types.Node }
		rspWithURI  struct{ Response string }

		checkNodeStatus = func(ID uint64, status string) {
			n := h.lookupNodeByID(aNodeID)
			h.a.NotNil(n)
			h.a.Equal(status, n.Status)
		}

		getNodeAuthToken = func(ID uint64) string {
			n := h.lookupNodeByID(ID)
			return n.AuthToken
		}
	)

	service.DefaultNode.SetHandshaker(nil)

	h.clearNodes()

	{
		// #############################################################################################################
		t.Log("Step #1, create node on server A")

		h.apiInit().
			// Debug().
			Post("/nodes/").
			FormData("name", "Server B").
			FormData("baseURL", "https://api.server-b.tld/federation").
			Expect(t).
			Status(http.StatusOK).
			Assert(helpers.AssertNoErrors).
			End().
			JSON(&rspWithNode)

		// response ok?
		h.a.NotNil(rspWithNode.Response)
		aNodeID = rspWithNode.Response.ID
		checkNodeStatus(aNodeID, types.NodeStatusPending)
	}

	{
		// #############################################################################################################
		t.Log("Step #2, request pairing URI")

		h.apiInit().
			//Debug().
			Post(fmt.Sprintf("/nodes/%d/uri", aNodeID)).
			Expect(t).
			Status(http.StatusOK).
			Assert(helpers.AssertNoErrors).
			End().
			JSON(&rspWithURI)

		// response ok?
		h.a.NotNil(rspWithNode.Response)
		aNodeURI = rspWithURI.Response

		checkNodeStatus(aNodeID, types.NodeStatusPending)
	}

	{
		// #############################################################################################################
		t.Log("Step #3, use pairing URI to create node on the 2nd server")

		h.apiInit().
			//Debug().
			Post("/nodes/").
			FormData("pairingURI", aNodeURI).
			Expect(t).
			Status(http.StatusOK).
			Assert(helpers.AssertNoErrors).
			End().
			JSON(&rspWithNode)

		// Node stored.
		h.a.NotNil(rspWithNode.Response)
		bNodeID = rspWithNode.Response.ID

		checkNodeStatus(aNodeID, types.NodeStatusPending)
		checkNodeStatus(bNodeID, types.NodeStatusPending)
	}

	{
		// #############################################################################################################
		t.Log("Step #4, admin of 2nd server requests list of nodes")

		h.apiInit().
			//Debug().
			Get("/nodes/").
			Query("status", types.NodeStatusPending).
			Expect(t).
			Status(http.StatusOK).
			Assert(helpers.AssertNoErrors).
			Assert(jsonpath.Present(fmt.Sprintf(`$.response.set[? @.sharedNodeID=="%d"]`, aNodeID))).
			Assert(jsonpath.Present(fmt.Sprintf(`$.response.set[? @.nodeID=="%d"]`, bNodeID))).
			End()
	}

	{
		// #############################################################################################################
		t.Log("Step #5, request pairing procedure on 2nd server")

		// fake remote call with local/direct change
		service.DefaultNode.SetHandshaker(&mockNodeHandshake{
			init: func(ctx context.Context, n *types.Node, authToken string) error {
				h.apiInit().
					//Debug().
					// make sure we do not use test auth-token for authentication but
					// we do it with pairing token
					Intercept(helpers.ReqHeaderAuthBearer(nil)).
					Post(fmt.Sprintf("/nodes/%d/handshake", n.SharedNodeID)).
					FormData("pairToken", n.PairToken).
					FormData("authToken", authToken).
					FormData("sharedNodeID", strconv.FormatUint(n.SharedNodeID, 10)).
					Expect(t).
					Status(http.StatusOK).
					Assert(helpers.AssertNoErrors).
					End()
				return nil
			},
		})

		h.apiInit().
			//Debug().
			Post(fmt.Sprintf("/nodes/%d/pair", bNodeID)).
			Expect(t).
			Status(http.StatusOK).
			Assert(helpers.AssertNoErrors).
			End()

		checkNodeStatus(aNodeID, types.NodeStatusPairRequested)
		checkNodeStatus(bNodeID, types.NodeStatusPairRequested)
	}

	{
		// #############################################################################################################
		t.Log("Step #6, handshake confirmation on 1st server")

		// fake remote call with local/direct change
		service.DefaultNode.SetHandshaker(&mockNodeHandshake{
			complete: func(ctx context.Context, n *types.Node, authToken string) error {
				h.apiInit().
					// Debug().
					// make sure we do not use test auth-token but
					// one provided to us in the initial handshake step
					Intercept(helpers.ReqHeaderRawAuthBearer(n.AuthToken)).
					Post(fmt.Sprintf("/nodes/%d/handshake-complete", n.SharedNodeID)).
					FormData("authToken", authToken).
					Expect(t).
					Status(http.StatusOK).
					Assert(helpers.AssertNoErrors).
					End()

				return nil
			},
		})

		h.apiInit().
			// Debug().
			Intercept(helpers.ReqHeaderRawAuthBearer(getNodeAuthToken(aNodeID))).
			Post(fmt.Sprintf("/nodes/%d/handshake-confirm", aNodeID)).
			Expect(t).
			Status(http.StatusOK).
			Assert(helpers.AssertNoErrors).
			End()

		checkNodeStatus(aNodeID, types.NodeStatusPaired)
		checkNodeStatus(bNodeID, types.NodeStatusPaired)
	}
}
