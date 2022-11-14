package service

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/cortezaproject/corteza-server/federation/types"
	"github.com/cortezaproject/corteza-server/pkg/actionlog"
	"github.com/cortezaproject/corteza-server/pkg/auth"
	"github.com/cortezaproject/corteza-server/pkg/options"
	"github.com/cortezaproject/corteza-server/pkg/rand"
	"github.com/cortezaproject/corteza-server/store"
	"github.com/cortezaproject/corteza-server/system/service"
	sysTypes "github.com/cortezaproject/corteza-server/system/types"
)

const (
	TokenLength = 32
)

type (
	tokenIssuer func(context.Context, auth.Identifiable) (token []byte, err error)

	node struct {
		store   store.Storer
		sysUser service.UserService

		actionlog actionlog.Recorder

		tokenIssuer tokenIssuer

		name    string
		host    string
		baseURL string

		handshaker nodeHandshaker

		ac nodeAccessController
	}

	nodeAccessController interface {
		CanPair(ctx context.Context) bool
		CanSearchNodes(ctx context.Context) bool
		CanCreateNode(ctx context.Context) bool
		CanManageNode(ctx context.Context, r *types.Node) bool
	}

	nodeUpdateHandler func(ctx context.Context, n *types.Node) error

	nodeHandshaker interface {
		Init(ctx context.Context, n *types.Node, authToken string) error
		Complete(ctx context.Context, n *types.Node, authToken string) error
	}
)

func Node(s store.Storer, u service.UserService, al actionlog.Recorder, th tokenIssuer, options options.FederationOpt, sopt options.HttpServerOpt, ac nodeAccessController) *node {
	return &node{
		store:       s,
		sysUser:     u,
		actionlog:   al,
		tokenIssuer: th,
		ac:          ac,
		name:        options.Label,
		host:        options.Host,

		baseURL: fmt.Sprintf("%s/federation", strings.TrimRight(sopt.ApiBaseUrl, "/")),

		handshaker: HttpHandshake(http.DefaultClient),
	}
}

func (svc *node) SetHandshaker(h nodeHandshaker) {
	svc.handshaker = h
}

func (svc node) Search(ctx context.Context, filter types.NodeFilter) (set types.NodeSet, f types.NodeFilter, err error) {
	var (
		aProps = &nodeActionProps{filter: &filter}
	)

	filter.Check = func(res *types.Node) (bool, error) {
		if !svc.ac.CanManageNode(ctx, res) {
			return false, NodeErrNotAllowedToManage()
		}

		return true, nil
	}

	err = func() error {
		if !svc.ac.CanSearchNodes(ctx) {
			return NodeErrNotAllowedToSearch()
		}

		if set, f, err = store.SearchFederationNodes(ctx, svc.store, filter); err != nil {
			return err
		}

		return nil
	}()

	return set, f, svc.recordAction(ctx, aProps, NodeActionSearch, err)
}

// Create is used on server A to create federation with server B
func (svc node) Create(ctx context.Context, new *types.Node) (*types.Node, error) {
	var (
		aProps = &nodeActionProps{}
		n      *types.Node
		err    error
	)

	// @todo verify new.baseURL; if it's set it must be valid URL that returns 200
	err = store.Tx(ctx, svc.store, func(ctx context.Context, s store.Storer) (err error) {
		n = &types.Node{
			ID:        nextID(),
			Name:      new.Name,
			BaseURL:   new.BaseURL,
			Contact:   new.Contact,
			Status:    types.NodeStatusPending,
			CreatedAt: *now(),
		}

		if !svc.ac.CanCreateNode(ctx) {
			return NodeErrNotAllowedToCreate()
		}

		return store.CreateFederationNode(ctx, s, n)
	})

	return n, svc.recordAction(ctx, aProps, NodeActionCreate, err)
}

// Read is used mainly in UI, when retrieving details about the node
func (svc node) Read(ctx context.Context, ID uint64) (*types.Node, error) {
	var (
		n, err = store.LookupFederationNodeByID(ctx, svc.store, ID)
		aProps = &nodeActionProps{node: n}
	)

	// permission takes precedence over any db error
	if n != nil && !svc.ac.CanManageNode(ctx, n) {
		err = NodeErrNotAllowedToManage()
	}

	return n, svc.recordAction(ctx, aProps, NodeActionCreate, err)
}

// CreateFromURI is used on server B to create federation with server A
func (svc node) CreateFromPairingURI(ctx context.Context, uri string) (n *types.Node, err error) {
	var (
		existing *types.Node
		aProps   = &nodeActionProps{pairingURI: uri}
	)

	err = func(ctx context.Context) error {
		if !svc.ac.CanPair(ctx) {
			return NodeErrNotAllowedToPair()
		}

		n, err = svc.decodePairingURI(uri)

		return err
	}(ctx)

	if err != nil {
		return nil, svc.recordAction(ctx, aProps, NodeActionCreateFromPairingURI, err)
	}

	existing, err = store.LookupFederationNodeByBaseURLSharedNodeID(ctx, svc.store, n.BaseURL, n.SharedNodeID)
	aProps.setNode(existing)

	if errors.Is(err, store.ErrNotFound) {
		if !svc.ac.CanCreateNode(ctx) {
			return n, svc.recordAction(ctx, aProps, NodeActionCreateFromPairingURI, NodeErrNotAllowedToCreate())
		}

		n.ID = nextID()
		n.CreatedAt = *now()
		n.UpdatedBy = auth.GetIdentityFromContext(ctx).Identity()
		n.Status = types.NodeStatusPending

		err = store.CreateFederationNode(ctx, svc.store, n)
		return n, svc.recordAction(ctx, aProps, NodeActionCreateFromPairingURI, err)

	} else if err != nil {
		return nil, svc.recordAction(ctx, aProps, NodeActionCreateFromPairingURI, err)
	} else {
		if !svc.ac.CanManageNode(ctx, n) {
			return n, svc.recordAction(ctx, aProps, NodeActionCreateFromPairingURI, NodeErrNotAllowedToManage())
		}

		// Node with the same sharing ID and domain already exists
		// so we'll reset status and tokens

		existing.Status = types.NodeStatusPending

		// Remove existing auth token
		existing.AuthToken = ""

		// Reset pairing token
		existing.PairToken = n.PairToken
		existing.UpdatedBy = auth.GetIdentityFromContext(ctx).Identity()
		existing.UpdatedAt = now()
		err = store.UpdateFederationNode(ctx, svc.store, n)
		return n, svc.recordAction(ctx, aProps, NodeActionRecreateFromPairingURI, err)
	}
}

// RegenerateNodeURI loads and updates node with the new OTT and returns sharable link
func (svc node) RegenerateNodeURI(ctx context.Context, nodeID uint64) (string, error) {
	var (
		uri string
	)

	_, err := svc.updater(
		ctx,
		nodeID,
		NodeActionOttRegenerated,
		func(ctx context.Context, n *types.Node) error {
			n.PairToken = string(rand.Bytes(TokenLength))
			return nil
		},
		func(ctx context.Context, n *types.Node) error {
			uri = svc.makePairingURI(n)
			return nil
		},
	)

	if err != nil {
		return "", err
	}

	return uri, nil
}

func (svc node) Update(ctx context.Context, upd *types.Node) (*types.Node, error) {
	return svc.updater(ctx, upd.ID, NodeActionUpdate, func(ctx context.Context, n *types.Node) error {
		n.Name = upd.Name
		n.BaseURL = upd.BaseURL
		n.Contact = upd.Contact

		n.UpdatedBy = auth.GetIdentityFromContext(ctx).Identity()
		n.UpdatedAt = now()

		return nil
	}, nil)
}

func (svc node) DeleteByID(ctx context.Context, ID uint64) error {
	_, err := svc.updater(ctx, ID, NodeActionDelete, func(ctx context.Context, n *types.Node) error {
		n.DeletedAt = now()
		n.DeletedBy = auth.GetIdentityFromContext(ctx).Identity()

		return nil
	}, nil)
	return err
}

func (svc node) UndeleteByID(ctx context.Context, ID uint64) error {
	_, err := svc.updater(ctx, ID, NodeActionUndelete, func(ctx context.Context, n *types.Node) error {
		n.DeletedAt = nil
		n.DeletedBy = 0

		return nil
	}, nil)
	return err
}

// Pair is used on server B to send request to server A
func (svc node) Pair(ctx context.Context, nodeID uint64) error {
	if !svc.ac.CanPair(ctx) {
		return NodeErrNotAllowedToPair()
	}

	// elevate permissions for user lookup & creation!
	ctx = auth.SetIdentityToContext(ctx, auth.FederationUser())

	_, err := svc.updater(
		ctx,
		nodeID,
		NodeActionPair,
		func(ctx context.Context, n *types.Node) error {
			// Handle federated user
			u, err := svc.fetchFederatedUser(ctx, n)
			if err != nil {
				return err
			}

			var accessToken []byte
			// Generate JWT token for the federated user
			accessToken, err = svc.tokenIssuer(ctx, u)
			if err != nil {
				return err
			}

			n.UpdatedBy = auth.GetIdentityFromContext(ctx).Identity()
			n.UpdatedAt = now()

			// Start handshake initialization remote node
			if err = svc.handshaker.Init(ctx, n, string(accessToken)); err != nil {
				return err
			}

			n.Status = types.NodeStatusPairRequested
			return nil
		},
		nil,
	)

	return err
}

// HandshakeInit is used on server A to handle pairing request (see Pair fn above) from server B
func (svc node) HandshakeInit(ctx context.Context, nodeID uint64, pairToken string, sharedNodeID uint64, authToken string) error {
	// if !svc.ac.CanPair(ctx) {
	// 	return NodeErrNotAllowedToPair()
	// }

	_, err := svc.updater(
		ctx,
		sharedNodeID,
		NodeActionHandshakeInit,
		func(ctx context.Context, n *types.Node) error {
			// @todo need to check node status before we can proceed with initialization
			if n.PairToken != pairToken {
				return NodeErrPairingTokenInvalid()
			}

			n.SharedNodeID = sharedNodeID
			n.AuthToken = authToken

			n.UpdatedBy = auth.GetIdentityFromContext(ctx).Identity()
			n.UpdatedAt = now()

			n.Status = types.NodeStatusPairRequested
			return nil
		},
		func(ctx context.Context, n *types.Node) error {
			// @todo notify the node administrator about the request
			return nil
		},
	)

	return err
}

// HandshakeConfirm is used by server A to manually confirm the handshake
func (svc node) HandshakeConfirm(ctx context.Context, nodeID uint64) error {
	if !svc.ac.CanPair(ctx) {
		return NodeErrNotAllowedToPair()
	}

	_, err := svc.updater(ctx, nodeID, NodeActionHandshakeConfirm, func(ctx context.Context, n *types.Node) error {
		// Handle federated user
		u, err := svc.fetchFederatedUser(ctx, n)
		if err != nil {
			return err
		}

		// Generate JWT token for the federated user
		var accessToken []byte
		if accessToken, err = svc.tokenIssuer(ctx, u); err != nil {
			return fmt.Errorf("could not confirm handshake: %w", err)
		}

		n.UpdatedBy = auth.GetIdentityFromContext(ctx).Identity()
		n.UpdatedAt = now()

		// Complete handshake on remote node
		if err = svc.handshaker.Complete(ctx, n, string(accessToken)); err != nil {
			return err
		}

		n.Status = types.NodeStatusPaired
		return nil
	}, nil)
	return err
}

// HandshakeComplete is used by server B to handle handshake confirmation
func (svc node) HandshakeComplete(ctx context.Context, sharedNodeID uint64, token string) error {
	var (
		n   *types.Node
		err error
	)

	if !svc.ac.CanPair(ctx) {
		return NodeErrNotAllowedToPair()
	}

	if n, err = store.LookupFederationNodeBySharedNodeID(ctx, svc.store, sharedNodeID); err != nil {
		return err
	}

	_, err = svc.updater(ctx, n.ID, NodeActionHandshakeComplete, func(ctx context.Context, n *types.Node) error {
		n.AuthToken = token
		n.Status = types.NodeStatusPaired
		n.UpdatedBy = auth.GetIdentityFromContext(ctx).Identity()
		n.UpdatedAt = now()

		return nil
	}, nil)

	return err
}

func (svc node) updater(ctx context.Context, nodeID uint64, action func(...*nodeActionProps) *nodeAction, fn, afterFn nodeUpdateHandler) (*types.Node, error) {
	var (
		err    error
		n      *types.Node
		aProps = &nodeActionProps{node: &types.Node{ID: nodeID}}
	)

	err = func() error {
		n, err = store.LookupFederationNodeByID(ctx, svc.store, nodeID)

		if errors.Is(err, store.ErrNotFound) {
			return NodeErrNotFound()
		}

		aProps.setNode(n)

		if err = fn(ctx, n); err != nil {
			return err
		}

		// any other type of access check
		// should be done in fn() before
		// if !svc.ac.CanManageNode(ctx, n) {
		// 	return NodeErrNotAllowedToManage()
		// }

		if err = store.UpdateFederationNode(ctx, svc.store, n); err != nil {
			return err
		}

		if afterFn != nil {
			if err = afterFn(ctx, n); err != nil {
				return err
			}
		}

		return nil
	}()

	return n, svc.recordAction(ctx, aProps, action, err)
}

func (svc node) FindBySharedNodeID(ctx context.Context, sharedNodeID uint64) (*types.Node, error) {
	n, err := svc.store.LookupFederationNodeBySharedNodeID(ctx, sharedNodeID)

	if n != nil && !svc.ac.CanManageNode(ctx, n) {
		return nil, NodeErrNotAllowedToManage()
	}

	return n, err
}

func (svc node) FindByID(ctx context.Context, nodeID uint64) (*types.Node, error) {
	n, err := svc.store.LookupFederationNodeByID(ctx, nodeID)

	if n != nil && !svc.ac.CanManageNode(ctx, n) {
		return nil, NodeErrNotAllowedToManage()
	}

	return n, err
}

// Looks for existing user or crates a new one
func (svc node) fetchFederatedUser(ctx context.Context, n *types.Node) (*sysTypes.User, error) {
	// Generate handle for user that se this node
	uHandle := fmt.Sprintf("federation_%d", n.ID)

	u, err := svc.sysUser.FindByHandle(ctx, uHandle)
	if err == nil {
		// Reuse existing user
		return u, nil
	}

	if service.UserErrNotFound().Is(err) {
		user := &sysTypes.User{
			Email:  strconv.FormatUint(n.ID, 10) + "@federation.corteza",
			Handle: uHandle,
		}

		AddFederationLabel(user, "federation", n.BaseURL)

		// Create a user to service this node
		r, err := service.DefaultRole.FindByHandle(ctx, "federation")

		if err != nil {
			return nil, err
		}

		// context with fed. user credentials
		ctxfed := auth.SetIdentityToContext(ctx, auth.FederationUser())

		u, err = svc.sysUser.Create(ctxfed, user)

		if err != nil {
			return nil, err
		}

		if err = service.DefaultRole.MemberAdd(ctxfed, r.ID, u.ID); err != nil {
			return nil, err
		}

		u.SetRoles(append(u.Roles(), r.ID)...)

		return u, nil
	}

	return nil, err
}

// decodePairingURI decodes URI (string) to federation node
//
// Four parts are collected from the given URI:
//  1) node host from URI's host
//  2) shared node ID from URI's username
//  3) shared token from URI's password
//  4) name of the node from query string param "name" (optional)
func (node) decodePairingURI(uri string) (*types.Node, error) {
	var (
		n = &types.Node{}
	)

	return n, func() error {
		parsedURI, err := url.Parse(uri)
		if err != nil {
			return NodeErrPairingURIInvalid().Wrap(err)
		}

		n.PairToken, _ = parsedURI.User.Password()
		if len(n.PairToken) != TokenLength {
			return NodeErrPairingURITokenInvalid()
		}

		n.SharedNodeID, err = strconv.ParseUint(parsedURI.User.Username(), 10, 64)
		if err != nil || n.SharedNodeID == 0 {
			return NodeErrPairingURISourceIDInvalid().Wrap(err)
		}

		n.Name = parsedURI.Query().Get("name")
		n.BaseURL = fmt.Sprintf("%s://%s/%s", parsedURI.Scheme, parsedURI.Host, strings.Trim(parsedURI.Path, "/"))
		return nil
	}()
}

// makePairingURI encodes details about this deployment and pairing token into sharable URI
//
// This URI contains info about this server (name, host, fed. api base url) and node created here
// that is used to identify remote federation server ID with pairing token
//
// URI structure:
// corteza+federation://<node ID>:<pairing token>@<this-host-where-the-api-is></path-to-federation-api>?qs-meta-data
func (svc node) makePairingURI(n *types.Node) string {
	uri := url.URL{
		Scheme: "https",
		User:   url.UserPassword(strconv.FormatUint(n.ID, 10), n.PairToken),
		Host:   svc.host,
		Path:   svc.baseURL,
	}

	qs := url.Values{}
	if len(svc.name) > 0 {
		qs.Add("name", svc.name)
	}

	if len(qs) > 0 {
		uri.RawQuery = qs.Encode()
	}

	return uri.String()
}
