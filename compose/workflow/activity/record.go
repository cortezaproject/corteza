package activity

import (
	"context"
	"github.com/cortezaproject/corteza-server/compose/service"
	wf "github.com/cortezaproject/corteza-server/pkg/workflow"
)

type (
	Record struct {
		service service.RecordService
	}
)

func (r *Record) LookupByID(ctx context.Context, args wf.Variables) (wf.Variables, error) {
	var (
		namespaceID = args.Uint64("namespaceID", 0)
		moduleID    = args.Uint64("moduleID", 0)
		recordID    = args.Uint64("recordID", 0)
	)

	svc := r.service.With(ctx)
	rec, err := svc.FindByID(namespaceID, moduleID, recordID)
	if err != nil {
		return nil, err
	}

	return wf.Variables{"record": rec}, nil
}

func (r *Record) ValueMapper(ctx context.Context, args wf.Variables) (wf.Variables, error) {
	// ....
	return nil, nil
}

func (r *Record) Find(ctx context.Context, args wf.Variables) (wf.Variables, error) {
	// ....
	return nil, nil
}
