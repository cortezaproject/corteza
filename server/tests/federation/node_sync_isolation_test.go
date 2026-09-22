package federation

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/crusttech/human/server/federation/service"
	"github.com/crusttech/human/server/federation/types"
	"github.com/crusttech/human/server/pkg/auth"
	"github.com/crusttech/human/server/pkg/id"
	"github.com/crusttech/human/server/store"
	st "github.com/crusttech/human/server/system/types"
	"github.com/crusttech/human/server/tests/helpers"
)

// All federated users share one role that can manage every node,
// a paired node must still only get to the data that is shared with it
func TestNodeSyncIsolation(t *testing.T) {
	h := newHelper(t)
	h.clearNodes()
	helpers.AllowMe(h, types.NodeRbacResource(0), "manage")

	var (
		ctx = context.Background()

		makeNode = func(sharedNodeID uint64) (*types.Node, *st.User) {
			n := &types.Node{ID: id.Next(), SharedNodeID: sharedNodeID, Name: "node", BaseURL: "https://example.tld/federation", Status: types.NodeStatusPaired, CreatedAt: time.Now()}
			h.noError(store.CreateFederationNode(ctx, service.DefaultStore, n))

			u := &st.User{ID: id.Next(), Handle: fmt.Sprintf("federation_%d", n.ID), Email: fmt.Sprintf("%d@federation.corteza", n.ID), CreatedAt: time.Now()}
			h.noError(store.CreateUser(ctx, service.DefaultStore, u))
			u.SetRoles(h.roleID)
			return n, u
		}

		nodeA, userA = makeNode(id.Next())
		nodeB, userB = makeNode(id.Next())

		asNodeA = auth.SetIdentityToContext(ctx, userA)
	)

	n, err := service.DefaultNode.FindBySharedNodeID(asNodeA, nodeA.SharedNodeID)
	h.noError(err)
	h.a.Equal(nodeA.ID, n.ID)

	_, err = service.DefaultNode.FindBySharedNodeID(asNodeA, nodeB.SharedNodeID)
	h.a.Error(err, "node A must not be able to act as node B")

	// users can rename themselves, this must not help
	userB.Handle = "renamed"
	h.noError(store.UpdateUser(ctx, service.DefaultStore, userB))
	userA.Handle = fmt.Sprintf("federation_%d", nodeB.ID)
	h.noError(store.UpdateUser(ctx, service.DefaultStore, userA))
	_, err = service.DefaultNode.FindBySharedNodeID(asNodeA, nodeB.SharedNodeID)
	h.a.Error(err, "node A must not be able to act as node B after renaming")

	// users that are not federated (administrators) are limited by permissions only
	n, err = service.DefaultNode.FindBySharedNodeID(h.secCtx(), nodeB.SharedNodeID)
	h.noError(err)
	h.a.Equal(nodeB.ID, n.ID)
}
