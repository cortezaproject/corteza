package rest

import (
	"context"
	"fmt"
	"strconv"
	"time"

	cs "github.com/cortezaproject/corteza-server/compose/service"
	ct "github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/federation/rest/request"
	"github.com/cortezaproject/corteza-server/federation/service"
	"github.com/cortezaproject/corteza-server/federation/types"
	"github.com/cortezaproject/corteza-server/pkg/filter"
	"github.com/davecgh/go-spew/spew"
)

type (
	SyncData struct{}

	recordPayload struct {
		*ct.Record

		Records ct.RecordSet `json:"records,omitempty"`

		CanUpdateRecord bool `json:"canUpdateRecord"`
		CanDeleteRecord bool `json:"canDeleteRecord"`
	}

	listRecordResponse struct {
		Filter *ct.RecordFilter `json:"filter,omitempty"`
		Set    *ct.RecordSet    `json:"set"`
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
	var (
		err  error
		node *types.Node
		f    = ct.RecordFilter{}
		em   *types.ExposedModule
	)

	if node, err = service.DefaultNode.FindBySharedNodeID(ctx, r.NodeID); err != nil {
		return nil, err
	}

	// use the fetched node
	if em, err = (service.ExposedModule()).FindByID(ctx, node.ID, r.ModuleID); err != nil {
		return nil, err
	}

	f.ModuleID = em.ComposeModuleID
	f.Query = r.Query

	if f.Paging, err = filter.NewPaging(r.Limit, r.PageCursor); err != nil {
		return nil, err
	}

	if f.Sorting, err = filter.NewSorting(r.Sort); err != nil {
		return nil, err
	}

	if r.LastSync != "" {
		ls, err := parseLastSync(r.LastSync)

		if err != nil {
			return nil, err
		}

		f.Query = fmt.Sprintf("updated_at >= '%s' OR created_at >= '%s'", ls, ls)
	}

	list, f, err := (cs.Record()).Find(f)

	if err != nil {
		return nil, err
	}

	// do the actual field filtering
	err = list.Walk(func(r *ct.Record) error {
		r.Values, err = r.Values.Filter(func(rv *ct.RecordValue) (bool, error) {
			return em.Fields.HasField(rv.Name)
		})

		return err
	})

	if err != nil {
		return nil, err
	}

	return listRecordResponse{
		Set:    &list,
		Filter: &f,
	}, nil
}

func parseLastSync(lastSync string) (*time.Time, error) {
	spew.Dump("LASTSYNC", lastSync)
	if i, err := strconv.ParseInt(lastSync, 10, 64); err == nil {
		t := time.Unix(i, 0)
		return &t, nil
	}

	// try different format if above fails
	if t, err := time.Parse("2006-01-02 15:04:05", lastSync); err == nil {
		return &t, nil
	}

	t, err := time.Parse("2006-01-02", lastSync)

	if err != nil {
		return nil, err
	}

	return &t, nil
}
