package rest

import (
	"context"

	cs "github.com/cortezaproject/corteza-server/compose/service"
	"github.com/cortezaproject/corteza-server/compose/types"
	ct "github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/federation/rest/request"
	"github.com/cortezaproject/corteza-server/federation/service"
	"github.com/davecgh/go-spew/spew"
)

type (
	SyncData struct{}

	recordPayload struct {
		*types.Record

		Records ct.RecordSet `json:"records,omitempty"`

		CanUpdateRecord bool `json:"canUpdateRecord"`
		CanDeleteRecord bool `json:"canDeleteRecord"`
	}

	recordSetPayload struct {
		Filter *ct.RecordFilter `json:"filter,omitempty"`
		Set    []*recordPayload `json:"set"`
	}
)

func (SyncData) New() *SyncData {
	return &SyncData{}
}

func (ctrl SyncData) ReadExposedAll(ctx context.Context, r *request.SyncDataReadExposedAll) (interface{}, error) {
	// /exposed/records/
	// fetch the node info here so we can use the node ID in the filtering

	return nil, nil
}

func (ctrl SyncData) ReadExposed(ctx context.Context, r *request.SyncDataReadExposed) (interface{}, error) {
	// /exposed/modules/{moduleID}/records/
	// fetch the node info here so we can use the node ID in the filtering

	var (
		err error
		f   = ct.RecordFilter{}
	)

	// use the fetched node
	em, err := (service.ExposedModule()).FindByID(ctx, 276342359342989444, r.ModuleID)

	if err != nil {
		return nil, err
	}

	filter := ct.RecordFilter{ModuleID: em.ComposeModuleID}

	// find the compose module records
	list, f, err := (cs.Record()).Find(filter)

	if err != nil {
		return nil, err
	}

	spew.Dump(list, f, err)

	return nil, nil
}
