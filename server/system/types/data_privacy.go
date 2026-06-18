package types

import (
	"encoding/json"
	"time"

	"github.com/crusttech/human/server/pkg/geolocation"

	"github.com/crusttech/human/server/pkg/filter"
)

type (
	DataPrivacyRequestFilter struct {
		RequestID   []string `json:"requestID"`
		RequestedBy []string `json:"requestedBy"`

		Query string `json:"query"`

		Kind   []string `json:"kind"`
		Status []string `json:"status"`

		// Check fn is called by store backend for each resource found function can
		// modify the resource and return false if store should not return it
		//
		// Store then loads additional resources to satisfy the paging parameters
		Check func(request *DataPrivacyRequest) (bool, error) `json:"-"`

		// Standard helpers for paging and sorting
		filter.Sorting
		filter.Paging
	}

	// @todo temporarily using this; revert to proper struct after dev release
	DataPrivacyRequestPayloadSet []map[string]any
	// DataPrivacyRequestPayloadSet []DataPrivacyRequestPayload
	DataPrivacyRequestPayload struct {
		ConnectionID uint64                     `json:"connectionID,string"`
		NamespaceID  uint64                     `json:"namespaceID,string"`
		ModuleID     uint64                     `json:"moduleID,string"`
		Records      []DataPrivacyRequestRecord `json:"records"`
	}

	DataPrivacyRequestRecord struct {
		RecordID uint64                          `json:"recordID,string"`
		Values   []DataPrivacyRequestRecordValue `json:"values"`
	}

	DataPrivacyRequestRecordValue struct {
		Name     string `json:"name"`
		Value    string `json:"value"`
		Position int    `json:"position"`
	}

	DataPrivacyRequestCommentFilter struct {
		RequestID []string `json:"requestID"`

		// Check fn is called by store backend for each resource found function can
		// modify the resource and return false if store should not return it
		//
		// Store then loads additional resources to satisfy the paging parameters
		Check func(request *DataPrivacyRequestComment) (bool, error) `json:"-"`

		// Standard helpers for paging and sorting
		filter.Sorting
		filter.Paging
	}

	PrivacyDalConnection struct {
		ID     uint64 `json:"connectionID,string"`
		Handle string `json:"handle"`
		Type   string `json:"type"`

		// descriptions, notes, and other user-provided meta-data
		Meta PrivacyDalConnectionMeta `json:"meta"`

		// collection of configurations for various subsystems that
		// use this connection and how it affects their behaviour
		Config PrivacyDalConnectionConfig `json:"config"`

		CreatedAt time.Time  `json:"createdAt,omitempty"`
		CreatedBy uint64     `json:"createdBy,string" `
		UpdatedAt *time.Time `json:"updatedAt,omitempty"`
		UpdatedBy uint64     `json:"updatedBy,string,omitempty" `
		DeletedAt *time.Time `json:"deletedAt,omitempty"`
		DeletedBy uint64     `json:"deletedBy,string,omitempty" `
	}

	PrivacyDalConnectionMeta struct {
		Name      string           `json:"name"`
		Ownership string           `json:"ownership"`
		Location  geolocation.Full `json:"location"`
	}

	PrivacyDalConnectionConfig struct {
		Privacy DalConnectionConfigPrivacy `json:"privacy"`
	}

	RequestStatus string
	RequestKind   string
)

const (
	// RequestKindCorrect to correct module fields
	RequestKindCorrect RequestKind = "correct"
	// RequestKindDelete to delete module fields
	RequestKindDelete RequestKind = "delete"
	// RequestKindExport to export module fields
	RequestKindExport RequestKind = "export"

	// RequestStatusPending initially request will be in pending status
	RequestStatusPending RequestStatus = "pending"
	// RequestStatusCanceled owner of request has cancelled the request
	RequestStatusCanceled RequestStatus = "canceled"
	// RequestStatusApproved data officer has of request has cancelled the request
	RequestStatusApproved RequestStatus = "approved"
	// RequestStatusRejected data officer has denied the request
	RequestStatusRejected RequestStatus = "rejected"
)

func CastToRequestKind(s string) RequestKind {
	switch s {
	case "correct":
		return RequestKindCorrect
	case "delete":
		return RequestKindDelete
	case "export":
		return RequestKindExport
	default:
		return ""
	}
}

func (k RequestKind) String() string {
	return string(k)
}

func CastToRequestStatus(s string) RequestStatus {
	switch s {
	case "pending":
		return RequestStatusPending
	case "canceled":
		return RequestStatusCanceled
	case "approved":
		return RequestStatusApproved
	case "rejected":
		return RequestStatusRejected
	default:
		return ""
	}
}

func (s RequestStatus) String() string {
	return string(s)
}

func ParseDataPrivacyRequestPayload(ii []string) (out DataPrivacyRequestPayloadSet, err error) {
	for _, i := range ii {
		aux := make(map[string]any)
		err = json.Unmarshal([]byte(i), &aux)
		if err != nil {
			return nil, err
		}

		out = append(out, aux)
	}

	return out, err
}
