package rest

import (
	"context"
	"net/http"

	"github.com/cortezaproject/corteza/server/federation/rest/request"
	"github.com/cortezaproject/corteza/server/federation/service"
	"github.com/cortezaproject/corteza/server/federation/types"
	"github.com/cortezaproject/corteza/server/pkg/errors"
	"github.com/cortezaproject/corteza/server/pkg/federation"
	"github.com/cortezaproject/corteza/server/pkg/filter"
)

type (
	SyncStructure struct{}
)

func (SyncStructure) New() *SyncStructure {
	return &SyncStructure{}
}

// ReadExposedInternal gets the exposed module info and serves
// the internal Corteza format of the structure
func (ctrl SyncStructure) ReadExposedInternal(ctx context.Context, r *request.SyncStructureReadExposedInternal) (interface{}, error) {
	return func(w http.ResponseWriter, req *http.Request) {
		var (
			err error
			ef  federation.EncodingFormat = federation.CortezaInternalStructure
		)

		w.Header().Add("Content-Type", "application/json")

		fEncoder := federation.NewEncoder(w, service.DefaultOptions)

		payload, err := ctrl.readExposedAll(ctx, r)

		if err != nil {
			errors.ServeHTTP(w, req, err, false)
			return
		}

		err = fEncoder.Encode(*payload, ef)

		if err != nil {
			errors.ServeHTTP(w, req, err, false)
			return
		}

		return
	}, nil
}

// ReadExposedSocial gets the exposed module info and serves
// the activity streams format of the structure
func (ctrl SyncStructure) ReadExposedSocial(ctx context.Context, r *request.SyncStructureReadExposedSocial) (interface{}, error) {
	return func(w http.ResponseWriter, req *http.Request) {
		var (
			err error
			ef  federation.EncodingFormat = federation.ActivityStreamsStructure
		)

		w.Header().Add("Content-Type", "application/json")

		fEncoder := federation.NewEncoder(w, service.DefaultOptions)

		rr := request.SyncStructureReadExposedInternal{
			NodeID:     r.NodeID,
			LastSync:   r.LastSync,
			Query:      r.Query,
			Limit:      r.Limit,
			PageCursor: r.PageCursor,
			Sort:       r.Sort,
		}

		payload, err := ctrl.readExposedAll(ctx, &rr)

		if err != nil {
			errors.ServeHTTP(w, req, err, false)
			return
		}

		err = fEncoder.Encode(*payload, ef)

		if err != nil {
			errors.ServeHTTP(w, req, err, false)
			return
		}

		return
	}, nil
}

// readExposedAll fetches the exposed modules for the specific node
func (ctrl SyncStructure) readExposedAll(ctx context.Context, r *request.SyncStructureReadExposedInternal) (*federation.ListStructurePayload, error) {
	var (
		err  error
		node *types.Node
	)

	if node, err = service.DefaultNode.FindBySharedNodeID(ctx, r.NodeID); err != nil {
		return nil, err
	}

	f := types.ExposedModuleFilter{
		NodeID:   node.ID,
		LastSync: r.LastSync,
	}

	if f.Paging, err = filter.NewPaging(r.Limit, r.PageCursor); err != nil {
		return nil, err
	}

	if f.Sorting, err = filter.NewSorting(r.Sort); err != nil {
		return nil, err
	}

	list, f, err := (service.ExposedModule()).Find(ctx, f)

	if err != nil {
		return nil, err
	}

	return &federation.ListStructurePayload{
		NodeID: node.ID,
		Filter: &f,
		Set:    &list,
	}, nil
}
