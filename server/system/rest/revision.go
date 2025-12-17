package rest

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/cortezaproject/corteza/server/pkg/api"
	"github.com/cortezaproject/corteza/server/pkg/dal"
	"github.com/cortezaproject/corteza/server/pkg/revisions"
	"github.com/cortezaproject/corteza/server/system/rest/request"
	"github.com/cortezaproject/corteza/server/system/service"
)

type (
	Revision struct {
		svc interface {
			Search(ctx context.Context, f revisions.Filter) (dal.Iterator, error)
			Create(ctx context.Context, rev *revisions.Revision) error
			Delete(ctx context.Context, revisionID uint64) error
		}
	}
)

func (Revision) New() *Revision {
	return &Revision{
		svc: service.DefaultRevision,
	}
}

func (ctrl *Revision) List(ctx context.Context, r *request.RevisionList) (interface{}, error) {
	var (
		makeRev = func() dal.ValueSetter { return &revisions.Revision{} }
	)

	f := revisions.Filter{
		ResourceID:   r.ResourceID,
		ResourceType: r.ResourceType,
		Status:       r.Status,
		DeletedOnly:  r.DeletedOnly,
		Since:        r.Since,
	}

	iter, err := ctrl.svc.Search(ctx, f)
	if err != nil {
		return nil, err
	}

	return func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		if _, err = w.Write([]byte(`{"response":{"set":[`)); err != nil {
			return
		}

		err = dal.IteratorEncodeJSON(ctx, w, iter, makeRev)
		if err != nil {
			return
		}

		if _, err = w.Write([]byte(`]}}`)); err != nil {
			return
		}

		return
	}, err
}

func (ctrl *Revision) Create(ctx context.Context, r *request.RevisionCreate) (interface{}, error) {
	rev := revisions.Make(revisions.Updated, 0, r.ResourceID, 0)
	rev.ResourceType = r.ResourceType
	rev.Status = r.Status
	rev.Changes = r.Changes
	rev.Comment = r.Comment

	if err := ctrl.svc.Create(ctx, rev); err != nil {
		return nil, err
	}

	return rev, nil
}

func (ctrl *Revision) Delete(ctx context.Context, r *request.RevisionDelete) (interface{}, error) {
	return api.OK(), ctrl.svc.Delete(ctx, r.RevisionID)
}

// Helper to encode revisions to JSON (same pattern as record revisions)
func encodeRevisionsToJSON(ctx context.Context, w http.ResponseWriter, iter dal.Iterator, makeRev func() dal.ValueSetter) error {
	first := true
	for iter.Next(ctx) {
		rev := makeRev()
		if err := iter.Scan(rev); err != nil {
			return err
		}

		if !first {
			if _, err := w.Write([]byte(",")); err != nil {
				return err
			}
		}
		first = false

		b, err := json.Marshal(rev)
		if err != nil {
			return err
		}
		if _, err := w.Write(b); err != nil {
			return err
		}
	}

	return iter.Err()
}

