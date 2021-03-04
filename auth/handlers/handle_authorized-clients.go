package handlers

import (
	"github.com/cortezaproject/corteza-server/auth/request"
	"github.com/cortezaproject/corteza-server/pkg/errors"
	"github.com/cortezaproject/corteza-server/system/types"
	"sort"
	"strconv"
	"strings"
	"time"
)

type (
	authClients []authClient
	authClient  struct {
		ID          uint64
		Name        string
		Description string
		ConfirmedAt time.Time
	}
)

// Will sort clients - in order of creation
var _ sort.Interface = authClients{}

func (set authClients) Len() int      { return len(set) }
func (set authClients) Swap(i, j int) { set[i], set[j] = set[j], set[i] }
func (set authClients) Less(i, j int) bool {
	return strings.Compare(set[i].Name, set[i].Name) < 0
}

func (h *AuthHandlers) clientsView(req *request.AuthReq) error {
	req.Template = TmplAuthorizedClients

	ss, err := h.getAuthorizedClients(req)
	if err != nil {
		return err
	}

	sort.Sort(ss)
	req.Data["authorizedClients"] = ss

	return nil
}

func (h *AuthHandlers) clientsProc(req *request.AuthReq) error {
	switch {
	case len(req.Request.PostFormValue("revoke")) > 0:
		clientID, err := strconv.ParseUint(req.Request.PostFormValue("revoke"), 10, 64)
		if clientID == 0 {
			return err
		}

		if err = h.TokenService.DeleteByID(req.Context(), clientID); err != nil {
			return err
		}

		if err = h.ClientService.Revoke(req.Context(), req.AuthUser.User.ID, clientID); err != nil {
			return err
		}

		req.NewAlerts = append(req.NewAlerts, request.Alert{
			Type: "primary",
			Text: "Client authorization deleted",
		})
	}

	req.RedirectTo = GetLinks().AuthorizedClients
	return nil
}

func (h *AuthHandlers) getAuthorizedClients(req *request.AuthReq) (ss authClients, err error) {
	var (
		set    types.AuthConfirmedClientSet
		client *types.AuthClient
	)
	if set, err = h.ClientService.Confirmed(req.Context(), req.AuthUser.User.ID); err != nil {
		return
	}

	ss = make(authClients, 0, len(set))

	for i := range set {
		client, err = h.ClientService.LookupByID(req.Context(), set[i].ClientID)
		if errors.IsNotFound(err) {
			continue
		}

		if err != nil {
			return
		}

		ac := authClient{
			ID:          set[i].ClientID,
			Name:        client.Handle,
			ConfirmedAt: set[i].ConfirmedAt,
		}

		if client.Meta != nil {
			ac.Name = client.Meta.Name
		}

		ss = append(ss, ac)
	}

	return
}
