package types

import (
	"github.com/pkg/errors"
	"time"

	"github.com/cortezaproject/corteza-server/pkg/filter"
)

type (
	DataPrivacyRequest struct {
		ID uint64 `json:"requestID,string"`

		Kind   RequestKind   `json:"kind"`
		Status RequestStatus `json:"status"`

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

func (k RequestKind) IsValid() error {
	switch k {
	case RequestKindCorrect, RequestKindDelete, RequestKindExport:
		return nil
	}
	return errors.New("invalid request kind")
}

func (s RequestStatus) IsValid() error {
	switch s {
	case RequestStatusPending, RequestStatusCanceled, RequestStatusApproved, RequestStatusRejected:
		return nil
	}
	return errors.New("invalid request status")
}
