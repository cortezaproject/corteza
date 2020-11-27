package scim

import (
	"github.com/cortezaproject/corteza-server/system/types"
	"time"
)

type (
	metaResponse struct {
		ResourceType string     `json:"resourceType"`
		Created      time.Time  `json:"created"`
		LastModified *time.Time `json:"lastModified,omitempty"`
	}
)

func newUserMetaResponse(u *types.User) *metaResponse {
	rsp := &metaResponse{
		ResourceType: "User",
		Created:      u.CreatedAt,
		LastModified: u.UpdatedAt,
	}

	return rsp
}
