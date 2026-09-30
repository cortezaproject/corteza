package federation

import (
	"context"
	"testing"
	"time"

	ct "github.com/crusttech/human/server/compose/types"
	"github.com/crusttech/human/server/federation/service"
	"github.com/crusttech/human/server/federation/types"
	"github.com/crusttech/human/server/pkg/id"
	"github.com/crusttech/human/server/store"
	"github.com/crusttech/human/server/tests/helpers"
)

func (h helper) createNode() *types.Node {
	n := &types.Node{ID: id.Next(), Name: "node", BaseURL: "https://example.tld/federation", Status: types.NodeStatusPaired, CreatedAt: time.Now()}
	h.noError(store.CreateFederationNode(context.Background(), service.DefaultStore, n))
	return n
}

func TestSharedModuleUpdateKeepsFields(t *testing.T) {
	h := newHelper(t)
	h.clearNodes()
	helpers.AllowMe(h, types.NodeRbacResource(0), "manage")

	var (
		ctx  = context.Background()
		node = h.createNode()
		sm   = &types.SharedModule{
			ID:                         id.Next(),
			Handle:                     "shared",
			Name:                       "Shared",
			NodeID:                     node.ID,
			ExternalFederationModuleID: id.Next(),
			Fields:                     types.ModuleFieldSet{{Kind: "String", Name: "before"}},
			CreatedAt:                  time.Now(),
			CreatedBy:                  h.cUser.ID,
		}
	)
	h.noError(store.CreateFederationSharedModule(ctx, service.DefaultStore, sm))

	upd := sm.Clone()
	upd.Name = "Shared edited"
	upd.Fields = types.ModuleFieldSet{{Kind: "String", Name: "after"}, {Kind: "Number", Name: "added"}}

	_, err := service.DefaultSharedModule.Update(h.secCtx(), upd)
	h.noError(err)

	stored, err := store.LookupFederationSharedModuleByID(ctx, service.DefaultStore, sm.ID)
	h.noError(err)
	h.a.Equal("Shared edited", stored.Name)
	h.a.Len(stored.Fields, 2)
	h.a.Equal("after", stored.Fields[0].Name)
	h.a.Equal(sm.CreatedBy, stored.CreatedBy)
}

func TestExposedModuleUpdateKeepsFieldsAndModule(t *testing.T) {
	h := newHelper(t)
	h.clearNodes()
	helpers.AllowMe(h, types.NodeRbacResource(0), "manage")
	helpers.AllowMe(h, types.ExposedModuleRbacResource(0, 0), "manage")
	helpers.AllowMe(h, ct.NamespaceRbacResource(0), "read")
	helpers.AllowMe(h, ct.ModuleRbacResource(0, 0), "read", "update")

	var (
		ctx  = context.Background()
		node = h.createNode()
		ns   = &ct.Namespace{ID: id.Next(), Slug: "fed_ns", Name: "Fed NS", Enabled: true, CreatedAt: time.Now()}

		makeModule = func(handle string) *ct.Module {
			m := &ct.Module{ID: id.Next(), NamespaceID: ns.ID, Handle: handle, Name: handle, CreatedAt: time.Now()}
			h.noError(store.CreateComposeModule(ctx, service.DefaultStore, m))
			return m
		}
	)
	h.noError(store.CreateComposeNamespace(ctx, service.DefaultStore, ns))

	var (
		before = makeModule("fed_before")
		after  = makeModule("fed_after")
		em     = &types.ExposedModule{
			ID:                 id.Next(),
			Handle:             "exposed",
			Name:               "Exposed",
			NodeID:             node.ID,
			ComposeModuleID:    before.ID,
			ComposeNamespaceID: ns.ID,
			Fields:             types.ModuleFieldSet{{Kind: "String", Name: "before"}},
			CreatedAt:          time.Now(),
			CreatedBy:          h.cUser.ID,
		}
	)
	h.noError(store.CreateFederationExposedModule(ctx, service.DefaultStore, em))

	upd := em.Clone()
	upd.ComposeModuleID = after.ID
	upd.Fields = types.ModuleFieldSet{{Kind: "String", Name: "after"}}

	_, err := service.DefaultExposedModule.Update(h.secCtx(), upd)
	h.noError(err)

	stored, err := store.LookupFederationExposedModuleByID(ctx, service.DefaultStore, em.ID)
	h.noError(err)
	h.a.Equal(after.ID, stored.ComposeModuleID)
	h.a.Len(stored.Fields, 1)
	h.a.Equal("after", stored.Fields[0].Name)
	h.a.Equal(em.CreatedBy, stored.CreatedBy)
}
