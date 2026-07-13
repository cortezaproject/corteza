package service

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/crusttech/human/server/federation/types"
	"github.com/crusttech/human/server/pkg/actionlog"
	"github.com/crusttech/human/server/pkg/auth"
	"github.com/crusttech/human/server/pkg/options"
	"github.com/crusttech/human/server/pkg/rand"
	"github.com/crusttech/human/server/store"
	"github.com/crusttech/human/server/system/service"
	sysTypes "github.com/crusttech/human/server/system/types"
)

const (
	TokenLength = 32
)

type (
	tokenIssuer func(context.Context, auth.Identifiable) (token []byte, err error)

	nodeServices struct {
		sysUser     service.UserService
		tokenIssuer tokenIssuer
		name        string
		host        string
		baseURL     string
		handshaker  nodeHandshaker
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
		store:     s,
		actionlog: al,
		ac:        ac,
		services: &nodeServices{
			sysUser:     u,
			tokenIssuer: th,
			name:        options.Label,
			host:        options.Host,
			baseURL:     fmt.Sprintf("%s/federation", strings.TrimRight(sopt.ApiBaseUrl, "/")),
			handshaker:  HttpHandshake(http.DefaultClient),
		},
	}
}

func (svc *node) SetHandshaker(h nodeHandshaker) {
	svc.services.handshaker = h
}

// onSearch implements the actual search logic
func (svc *node) onSearch(ctx context.Context, filter types.NodeFilter, aProps *nodeActionProps) (types.NodeSet, types.NodeFilter, error) {
	filter.Check = func(res *types.Node) (bool, error) {
		if !svc.ac.CanManageNode(ctx, res) {
			return false, NodeErrNotAllowedToManage()
		}
		return true, nil
	}

	return store.SearchFederationNodes(ctx, svc.store, filter)
}

// onCreate implements the actual create logic
func (svc *node) onCreate(ctx context.Context, new *types.Node) error {
	return store.Tx(ctx, svc.store, func(ctx context.Context, s store.Storer) error {
		n := &types.Node{
			ID:        nextID(),
			Name:      new.Name,
			BaseURL:   new.BaseURL,
			Contact:   new.Contact,
			Status:    types.NodeStatusPending,
			CreatedAt: *now(),
		}

		*new = *n

		return store.CreateFederationNode(ctx, s, new)
	})
}

// onUpdate hook: called inside the gen Update after guard/staleness checks, before field copy and store.Update
func (svc *node) onUpdate(ctx context.Context, s store.Storer, upd *types.Node, res *types.Node, aProps *nodeActionProps, before func() error, after func() error) error {
	res.Name = upd.Name
	res.BaseURL = upd.BaseURL
	res.Contact = upd.Contact
	res.UpdatedBy = auth.GetIdentityFromContext(ctx).Identity()
	return nil
}

// onDelete soft-deletes the node
func (svc *node) onDelete(ctx context.Context, s store.Storer, res *types.Node, aProps *nodeActionProps) error {
	res.DeletedAt = now()
	res.DeletedBy = auth.GetIdentityFromContext(ctx).Identity()
	return store.UpdateFederationNode(ctx, s, res)
}

// onUndelete restores a soft-deleted node
func (svc *node) onUndelete(ctx context.Context, s store.Storer, res *types.Node, aProps *nodeActionProps) error {
	res.DeletedAt = nil
	res.DeletedBy = 0
	return store.UpdateFederationNode(ctx, s, res)
}

// onRead loads and checks access to a node
func (svc *node) onRead(ctx context.Context, aProps *nodeActionProps, ID uint64) (*types.Node, error) {
	n, err := loadNode(ctx, svc.store, ID)
	if err != nil {
		return nil, err
	}

	if !svc.ac.CanManageNode(ctx, n) {
		return nil, NodeErrNotAllowedToManage()
	}

	aProps.setNode(n)
	return n, nil
}

// onCreateFromPairingURI decodes a pairing URI and creates/updates a node
func (svc *node) onCreateFromPairingURI(ctx context.Context, aProps *nodeActionProps, uri string) (*types.Node, error) {
	n, err := svc.decodePairingURI(uri)
	if err != nil {
		return nil, err
	}

	aProps.setPairingURI(uri)

	existing, err := store.LookupFederationNodeByBaseURLSharedNodeID(ctx, svc.store, n.BaseURL, n.SharedNodeID)
	if err == nil {
		aProps.setNode(existing)
	}

	if err != nil && !isStoreNotFound(err) {
		return nil, err
	}

	if isStoreNotFound(err) {
		if !svc.ac.CanCreateNode(ctx) {
			return n, NodeErrNotAllowedToCreate()
		}

		n.ID = nextID()
		n.CreatedAt = *now()
		n.UpdatedBy = auth.GetIdentityFromContext(ctx).Identity()
		n.Status = types.NodeStatusPending

		return n, store.CreateFederationNode(ctx, svc.store, n)
	}

	if !svc.ac.CanManageNode(ctx, n) {
		return n, NodeErrNotAllowedToManage()
	}

	existing.Status = types.NodeStatusPending
	existing.AuthToken = ""
	existing.PairToken = n.PairToken
	existing.UpdatedBy = auth.GetIdentityFromContext(ctx).Identity()
	existing.UpdatedAt = now()
	return existing, store.UpdateFederationNode(ctx, svc.store, existing)
}

// onRegenerateNodeURI regenerates the pairing token and returns the new URI
func (svc *node) onRegenerateNodeURI(ctx context.Context, aProps *nodeActionProps, nodeID uint64) (string, error) {
	var uri string

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

	return uri, err
}

// onPair sends a pairing request to the remote node
func (svc *node) onPair(ctx context.Context, aProps *nodeActionProps, nodeID uint64) error {
	ctx = auth.SetIdentityToContext(ctx, auth.FederationUser())

	_, err := svc.updater(
		ctx,
		nodeID,
		NodeActionPair,
		func(ctx context.Context, n *types.Node) error {
			u, err := svc.fetchFederatedUser(ctx, n)
			if err != nil {
				return err
			}

			accessToken, err := svc.services.tokenIssuer(ctx, u)
			if err != nil {
				return err
			}

			n.UpdatedBy = auth.GetIdentityFromContext(ctx).Identity()
			n.UpdatedAt = now()

			if err = svc.services.handshaker.Init(ctx, n, string(accessToken)); err != nil {
				return err
			}

			n.Status = types.NodeStatusPairRequested
			return nil
		},
		nil,
	)

	return err
}

// onHandshakeInit handles the incoming pairing request on server A
func (svc *node) onHandshakeInit(ctx context.Context, aProps *nodeActionProps, nodeID uint64, pairToken string, sharedNodeID uint64, authToken string) error {
	_, err := svc.updater(
		ctx,
		sharedNodeID,
		NodeActionHandshakeInit,
		func(ctx context.Context, n *types.Node) error {
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
			return nil
		},
	)

	return err
}

// onHandshakeConfirm confirms the handshake on server A
func (svc *node) onHandshakeConfirm(ctx context.Context, aProps *nodeActionProps, nodeID uint64) error {
	_, err := svc.updater(ctx, nodeID, NodeActionHandshakeConfirm, func(ctx context.Context, n *types.Node) error {
		u, err := svc.fetchFederatedUser(ctx, n)
		if err != nil {
			return err
		}

		var accessToken []byte
		if accessToken, err = svc.services.tokenIssuer(ctx, u); err != nil {
			return fmt.Errorf("could not confirm handshake: %w", err)
		}

		n.UpdatedBy = auth.GetIdentityFromContext(ctx).Identity()
		n.UpdatedAt = now()

		if err = svc.services.handshaker.Complete(ctx, n, string(accessToken)); err != nil {
			return err
		}

		n.Status = types.NodeStatusPaired
		return nil
	}, nil)

	return err
}

// onHandshakeComplete handles the handshake completion on server B
func (svc *node) onHandshakeComplete(ctx context.Context, aProps *nodeActionProps, sharedNodeID uint64, token string) error {
	n, err := store.LookupFederationNodeBySharedNodeID(ctx, svc.store, sharedNodeID)
	if err != nil {
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
		n, err = loadNode(ctx, svc.store, nodeID)

		if err != nil {
			return err
		}

		aProps.setNode(n)

		if err = fn(ctx, n); err != nil {
			return err
		}

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
	n, err := store.LookupFederationNodeBySharedNodeID(ctx, svc.store, sharedNodeID)

	if n != nil && !svc.ac.CanManageNode(ctx, n) {
		return nil, NodeErrNotAllowedToManage()
	}

	return n, err
}

func (svc node) FindByID(ctx context.Context, nodeID uint64) (n *types.Node, err error) {
	if n, err = loadNode(ctx, svc.store, nodeID); err != nil {
		return
	}

	if !svc.ac.CanManageNode(ctx, n) {
		return nil, NodeErrNotAllowedToManage()
	}

	return
}

// Looks for existing user or creates a new one
func (svc node) fetchFederatedUser(ctx context.Context, n *types.Node) (*sysTypes.User, error) {
	uHandle := fmt.Sprintf("federation_%d", n.ID)

	u, err := svc.services.sysUser.FindByHandle(ctx, uHandle)
	if err == nil {
		return u, nil
	}

	if service.UserErrNotFound().Is(err) {
		user := &sysTypes.User{
			Email:  strconv.FormatUint(n.ID, 10) + "@federation.corteza",
			Handle: uHandle,
		}

		AddFederationLabel(user, "federation", n.BaseURL)

		r, err := service.DefaultRole.FindByHandle(ctx, "federation")
		if err != nil {
			return nil, err
		}

		ctxfed := auth.SetIdentityToContext(ctx, auth.FederationUser())

		u, err = svc.services.sysUser.Create(ctxfed, user)
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
func (svc node) makePairingURI(n *types.Node) string {
	uri := url.URL{
		Scheme: "https",
		User:   url.UserPassword(strconv.FormatUint(n.ID, 10), n.PairToken),
		Host:   svc.services.host,
		Path:   svc.services.baseURL,
	}

	qs := url.Values{}
	if len(svc.services.name) > 0 {
		qs.Add("name", svc.services.name)
	}

	if len(qs) > 0 {
		uri.RawQuery = qs.Encode()
	}

	return uri.String()
}

// isStoreNotFound checks if error is a not-found sentinel from the store
func isStoreNotFound(err error) bool {
	return err == store.ErrNotFound
}
