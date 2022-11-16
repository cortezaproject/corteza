package scim

import (
	"encoding/json"
	"fmt"
	"github.com/cortezaproject/corteza/server/system/types"
	"io"
	"strconv"
)

const (
	urnGroup                   = "urn:ietf:params:scim:schemas:core:2.0:Group"
	groupLabel_SCIM_externalId = "SCIM_externalId"
)

type (
	groupResourceResponse struct {
		Schemas    []string      `json:"schemas"`
		Meta       *metaResponse `json:"meta,omitempty"`
		ID         string        `json:"id,omitempty"`
		ExternalId string        `json:"externalId,omitempty"`
		Name       string        `json:"displayName"`
	}

	groupResourceRequest struct {
		Schemas    []string      `json:"schemas"`
		Meta       *metaResponse `json:"meta,omitempty"`
		ExternalId *string       `json:"externalId,omitempty"`
		Name       *string       `json:"displayName"`
	}
)

func newGroupResourceResponse(u *types.Role) *groupResourceResponse {
	rsp := &groupResourceResponse{
		Schemas:    []string{urnGroup},
		Meta:       newGroupMetaResponse(u),
		ID:         strconv.FormatUint(u.ID, 10),
		ExternalId: u.Labels[groupLabel_SCIM_externalId],
		Name:       u.Name,
	}

	return rsp
}

func (req *groupResourceRequest) decodeJSON(r io.Reader) error {
	if err := json.NewDecoder(r).Decode(req); err != nil {
		return fmt.Errorf("could not decode group payload: %w", err)
	}

	return nil
}

func (req *groupResourceRequest) applyTo(u *types.Role) {
	if req.Name != nil {
		u.Name = *req.Name
	}

	if req.ExternalId != nil {
		u.SetLabel("SCIM_externalId", *req.ExternalId)
	}
}
