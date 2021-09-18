package resource

import (
	"fmt"
	"strconv"

	"github.com/cortezaproject/corteza-server/system/types"
)

type (
	Report struct {
		*base
		Res *types.Report

		Sources     []*ReportSource
		Projections []*ReportProjection
	}

	ReportSource struct {
		*base
		Res *types.ReportDataSource
	}

	ReportProjection struct {
		*base
		Res *types.ReportProjection
	}
)

func NewReport(res *types.Report) *Report {
	r := &Report{
		base: &base{},
	}
	r.SetResourceType(types.ReportResourceType)
	r.Res = res

	r.AddIdentifier(identifiers(res.Handle, res.ID)...)

	// Initial stamps
	r.SetTimestamps(MakeTimestampsCUDA(&res.CreatedAt, res.UpdatedAt, res.DeletedAt, nil))
	r.SetUserstamps(MakeUserstampsCUDO(res.CreatedBy, res.UpdatedBy, res.DeletedBy, res.OwnedBy))

	return r
}

func (r *Report) AddReportSource(res *types.ReportDataSource) *ReportSource {
	s := &ReportSource{
		base: &base{},
	}

	s.Res = res
	r.Sources = append(r.Sources, s)

	return s
}

func (r *Report) RBACParts() (resource string, ref *Ref, path []*Ref) {
	ref = r.Ref()
	path = nil
	resource = fmt.Sprintf(types.ReportRbacResourceTpl(), types.ReportResourceType, firstOkString(strconv.FormatUint(r.Res.ID, 10), r.Res.Handle))

	return
}

func (r *Report) ResourceTranslationParts() (resource string, ref *Ref, path []*Ref) {
	ref = r.Ref()
	path = nil
	resource = fmt.Sprintf(types.ReportResourceTranslationTpl(), types.ReportResourceTranslationType, firstOkString(strconv.FormatUint(r.Res.ID, 10), r.Res.Handle))

	return
}

func (r *Report) AddReportProjection(res *types.ReportProjection) *ReportProjection {
	p := &ReportProjection{
		base: &base{},
	}

	p.Res = res
	r.Projections = append(r.Projections, p)

	return p
}

func (r *Report) SysID() uint64 {
	return r.Res.ID
}

func (r *Report) encodeTranslations() ([]*ResourceTranslation, error) {
	return nil, nil
}

// FindReport looks for the workflow in the resource set
func FindReport(rr InterfaceSet, ii Identifiers) (ns *types.Report) {
	var wfRes *Report

	rr.Walk(func(r Interface) error {
		wr, ok := r.(*Report)
		if !ok {
			return nil
		}

		if wr.Identifiers().HasAny(ii) {
			wfRes = wr
		}
		return nil
	})

	// Found it
	if wfRes != nil {
		return wfRes.Res
	}
	return nil
}

func ReportErrUnresolved(ii Identifiers) error {
	return fmt.Errorf("report unresolved %v", ii.StringSlice())
}
