package service

import (
	"context"
	"errors"
	"net/url"
	"strconv"

	"github.com/cortezaproject/corteza-server/federation/types"
	"github.com/cortezaproject/corteza-server/federation/util"
	"github.com/cortezaproject/corteza-server/pkg/actionlog"
	"github.com/cortezaproject/corteza-server/pkg/auth"
	"github.com/cortezaproject/corteza-server/pkg/id"
	"github.com/cortezaproject/corteza-server/pkg/rand"
	"github.com/cortezaproject/corteza-server/store"
	"github.com/cortezaproject/corteza-server/system/service"
	stypes "github.com/cortezaproject/corteza-server/system/types"
	"github.com/davecgh/go-spew/spew"
)

var (
	// This is temporary for local testing
	tmpNodeStore = make(types.NodeSet, 0)

	ErrorNodeNotFound   = errors.New("pair: node not found")
	ErrorInvalidNodeURI = errors.New("pair: invalid node uri provided")
)

const (
	TokenLength = 15
)

type (
	node struct {
		store   store.Storer
		sysUser service.UserService

		actionlog actionlog.Recorder

		tokenEncoder auth.TokenEncoder
	}

	NodeService interface {
		Create(ctx context.Context, n *types.Node, sharedDomain string) (*types.Node, error)
		CreateFromURI(ctx context.Context, uri string, sharedDomain string) (*types.Node, error)

		Pair(ctx context.Context, nodeID uint64) error
		HandshakeInit(ctx context.Context, sourceNodeID, nodeIDB uint64, nodeURI, token string) error
		HandshakeConfirm(ctx context.Context, nodeID uint64) error
		HandshakeComplete(ctx context.Context, nodeID uint64, token string) error
	}
)

func Node(s store.Storer, u service.UserService, al actionlog.Recorder, th auth.TokenHandler) NodeService {
	return (&node{
		store:        s,
		sysUser:      u,
		actionlog:    al,
		tokenEncoder: th,
	})
}

// @todo move myDomain to configuration
func (svc node) Create(ctx context.Context, n *types.Node, myDomain string) (*types.Node, error) {
	n.ID = id.Next()
	n.Status = types.NodeStatusPending

	ott := string(rand.Bytes(TokenLength))
	i := util.EncodeURI(ott, myDomain, n.ID)
	if n.NodeURI == "" {
		n.NodeURI = i
	}

	// @todo store an initial Node entry
	tmpNodeStore = append(tmpNodeStore, n)

	return n, nil
}

func (svc node) CreateFromURI(ctx context.Context, uri string, sharedDomain string) (*types.Node, error) {
	n := &types.Node{
		ID:     id.Next(),
		Status: types.NodeStatusPending,
	}

	uri, err := url.QueryUnescape(uri)
	if err != nil {
		return nil, err
	}
	pr, err := util.DecodeURI(uri)
	if err != nil {
		return nil, err
	}

	n.Domain = pr.Domain
	n.Name = pr.Params.Name
	n.Status = types.NodeStatusPending
	n.NodeURI = uri

	// @todo store an initial Node entry
	tmpNodeStore = append(tmpNodeStore, n)

	return n, nil
}

func (svc node) Pair(ctx context.Context, nodeID uint64) error {
	// @todo store
	n := tmpNodeStore.FindByID(nodeID)
	if n == nil {
		return ErrorNodeNotFound
	}

	// Handle fenedrated user
	u, err := svc.fetchFederatedUser(ctx, n)
	if err != nil {
		return err
	}

	// Generate JWT token for the federated user
	t := svc.tokenEncoder.Encode(u)
	// @todo remove
	spew.Dump(u, t)

	// Ping node A to request the handshake
	// @todo...

	return nil
}

func (svc node) HandshakeInit(ctx context.Context, nodeIDA, nodeIDB uint64, nodeURI, token string) error {
	// @todo store
	n := tmpNodeStore.FindByID(nodeIDA)
	if n == nil {
		return ErrorNodeNotFound
	}

	if n.NodeURI != nodeURI {
		return ErrorInvalidNodeURI
	}

	// @todo store...
	n.Token = token
	n.SharedID = nodeIDA
	n.Status = types.NodeStatusPairRequest
	spew.Dump(n)

	// Notify the node administrator about the request
	// @todo

	return nil
}

func (svc node) HandshakeConfirm(ctx context.Context, nodeID uint64) error {
	// @todo store
	n := tmpNodeStore.FindByID(nodeID)
	if n == nil {
		return ErrorNodeNotFound
	}

	// Handle fenedrated user
	u, err := svc.fetchFederatedUser(ctx, n)
	if err != nil {
		return err
	}

	// Generate JWT token for the federated user
	t := svc.tokenEncoder.Encode(u)
	// @todo remove
	spew.Dump(u, t)

	// Ping node B to complete the handshake
	// @todo...

	// @todo store
	n.Status = types.NodeStatusPairComplete

	return nil
}

func (svc node) HandshakeComplete(ctx context.Context, nodeID uint64, token string) error {
	// @todo store
	n := tmpNodeStore.FindByID(nodeID)
	if n == nil {
		return ErrorNodeNotFound
	}

	// Final update -- update the token and node status
	// @todo store
	n.Token = token
	n.Status = types.NodeStatusPairComplete

	return nil
}

func (svc node) fetchFederatedUser(ctx context.Context, n *types.Node) (*stypes.User, error) {
	// Handle fenedrated user
	uHandle := "federation_" + strconv.FormatUint(n.ID, 10)
	u, _ := svc.sysUser.With(ctx).FindByHandle(uHandle)

	if u == nil {
		var err error

		// Create a system user to service this node
		u, err = svc.sysUser.With(ctx).Create(&stypes.User{
			Email:  strconv.FormatUint(n.ID, 10) + "@federation.corteza",
			Handle: uHandle,
			Kind:   stypes.FederationUser,
		})
		if err != nil {
			return nil, err
		}
	}
	return u, nil
}
