package scim

import (
	"fmt"
	"github.com/cortezaproject/corteza/server/system/types"
	"net/http"
	"time"
)

type (
	metaResponse struct {
		ResourceType string     `json:"resourceType"`
		Created      time.Time  `json:"created"`
		LastModified *time.Time `json:"lastModified,omitempty"`
	}

	errorResponse struct {
		Schemas  []string `json:"schemas"`
		SCIMType string   `json:"scimType,omitempty"`
		Detail   string   `json:"detail,omitempty"`
		Status   int      `json:"status,string"`
	}
)

const (
	urnError = "urn:ietf:params:scim:api:messages:2.0:Error"
)

func newUserMetaResponse(u *types.User) *metaResponse {
	rsp := &metaResponse{
		ResourceType: "User",
		Created:      u.CreatedAt,
		LastModified: u.UpdatedAt,
	}

	return rsp
}

func newGroupMetaResponse(u *types.Role) *metaResponse {
	rsp := &metaResponse{
		ResourceType: "Group",
		Created:      u.CreatedAt,
		LastModified: u.UpdatedAt,
	}

	return rsp
}

func newErrorfResponse(httpStatus int, format string, aa ...interface{}) *errorResponse {
	return newErrorResponse(httpStatus, fmt.Errorf(format, aa...))
}

func newErrorResponse(httpStatus int, err error) *errorResponse {
	if httpStatus == 0 {
		httpStatus = http.StatusInternalServerError
	}

	er := &errorResponse{
		Schemas: []string{urnError},
		Status:  httpStatus,
	}

	if err != nil {
		er.Detail = err.Error()
	}

	return er
}

func (e *errorResponse) Error() string {
	return e.Detail
}
