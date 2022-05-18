package types

import (
	"time"

	"github.com/cortezaproject/corteza-server/pkg/filter"
)

type (
	DataPrivacyRequest struct {
		ID   uint64 `json:"requestID,string"`
		Name string `json:"name"`

		RequestType RequestType   `json:"requestType"`
		Status      RequestStatus `json:"status"`

		RequestedAt time.Time  `json:"requestedAt,omitempty"`
		RequestedBy uint64     `json:"requestedBy,string"`
		CompletedAt *time.Time `json:"completedAt,omitempty"`
		CompletedBy uint64     `json:"completedBy,string,omitempty" `

		CreatedAt time.Time  `json:"createdAt,omitempty"`
		CreatedBy uint64     `json:"createdBy,string" `
		UpdatedAt *time.Time `json:"updatedAt,omitempty"`
		UpdatedBy uint64     `json:"updatedBy,string,omitempty"`
		DeletedAt *time.Time `json:"deletedAt,omitempty"`
		DeletedBy uint64     `json:"deletedBy,string,omitempty"`
	}

	DataPrivacyRequestFilter struct {
		RequestID []uint64 `json:"requestID"`

		Name   string          `json:"name"`
		Status []RequestStatus `json:"status"`

		Deleted filter.State `json:"deleted"`

		// Check fn is called by store backend for each resource found function can
		// modify the resource and return false if store should not return it
		//
		// Store then loads additional resources to satisfy the paging parameters
		Check func(request *DataPrivacyRequest) (bool, error) `json:"-"`

		// Standard helpers for paging and sorting
		filter.Sorting
		filter.Paging
	}

	RequestStatus int
	RequestType   int
)

const (
	// RequestStatusPending initially request will be in pending status
	RequestStatusPending RequestStatus = 1
	// RequestStatusCancel owner of request has cancelled the request
	RequestStatusCancel RequestStatus = 2
	// RequestStatusApprove data officer has of request has cancelled the request
	RequestStatusApprove RequestStatus = 3
	// RequestStatusReject data officer has denied the request
	RequestStatusReject RequestStatus = 4

	// RequestTypeCorrection
	RequestTypeCorrection RequestType = 1
	// RequestTypeDeletion
	RequestTypeDeletion RequestType = 2
	// RequestTypeExport
	RequestTypeExport RequestType = 3
)
