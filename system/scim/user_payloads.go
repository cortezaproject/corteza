package scim

import (
	"encoding/json"
	"fmt"
	"github.com/cortezaproject/corteza-server/pkg/handle"
	"github.com/cortezaproject/corteza-server/system/types"
	"io"
	"strconv"
)

const (
	urnUser                   = "urn:ietf:params:scim:schemas:core:2.0:User"
	userLabel_SCIM_externalId = "SCIM_externalId"
)

type (
	emailResponse struct {
		Value   string `json:"value"`
		Primary bool   `json:"primary,omitempty"`
	}

	emailsResponse []*emailResponse

	userNameResponse struct {
		Formatted string `json:"formatted"`
	}

	userResourceResponse struct {
		Schemas    []string          `json:"schemas"`
		Meta       *metaResponse     `json:"meta,omitempty"`
		ID         string            `json:"id,omitempty"`
		ExternalId string            `json:"externalId,omitempty"`
		UserName   string            `json:"userName,omitempty"`
		NickName   string            `json:"nickName,omitempty"`
		Name       *userNameResponse `json:"displayName"`
		Emails     emailsResponse    `json:"emails,omitempty"`
	}

	userResourceRequest struct {
		Schemas    []string          `json:"schemas"`
		Meta       *metaResponse     `json:"meta,omitempty"`
		ExternalId *string           `json:"externalId,omitempty"`
		UserName   *string           `json:"userName,omitempty"`
		NickName   *string           `json:"nickName,omitempty"`
		Password   *string           `json:"password,omitempty"`
		Name       *userNameResponse `json:"name"`
		Emails     emailsResponse    `json:"emails,omitempty"`
	}
)

func newUserResourceResponse(u *types.User) *userResourceResponse {
	rsp := &userResourceResponse{
		Schemas:    []string{urnUser},
		Meta:       newUserMetaResponse(u),
		ID:         strconv.FormatUint(u.ID, 10),
		ExternalId: u.Labels[userLabel_SCIM_externalId],
		UserName:   u.Username,
		NickName:   u.Handle,
		Emails:     emailsResponse{{u.Email, true}},
	}

	if u.Name != "" {
		rsp.Name = &userNameResponse{Formatted: u.Name}
	}

	return rsp
}

// returns first (primary) email
func (ee emailsResponse) getFirst() string {
	if len(ee) == 0 {
		return ""
	}

	var match int
	for i, e := range ee {
		if e.Primary {
			match = i
			break
		}
	}

	return ee[match].Value
}

func (req *userResourceRequest) decodeJSON(r io.Reader) error {
	if err := json.NewDecoder(r).Decode(req); err != nil {
		return fmt.Errorf("could not decode user payload: %w", err)
	}

	return nil
}

func (req *userResourceRequest) applyTo(u *types.User) {
	if v := req.Emails.getFirst(); len(v) > 0 {
		u.Email = v
	}

	if req.Name != nil {
		u.Name = req.Name.Formatted
	}

	if req.UserName != nil {
		u.Username = *req.UserName
	}

	if req.NickName != nil && handle.IsValid(*req.NickName) {
		u.Handle = *req.NickName
	}

	if req.ExternalId != nil {
		u.SetLabel("SCIM_externalId", *req.ExternalId)
	}
}
