package rest

import (
	"context"
	"fmt"
	"strings"
	"time"

	cs "github.com/cortezaproject/corteza-server/compose/service"
	ct "github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/federation/rest/request"
	"github.com/cortezaproject/corteza-server/federation/service"
	"github.com/cortezaproject/corteza-server/federation/types"
	"github.com/cortezaproject/corteza-server/pkg/filter"
	ss "github.com/cortezaproject/corteza-server/system/service"
	st "github.com/cortezaproject/corteza-server/system/types"
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

	listResponse struct {
		Set *responseSet `json:"set"`
	}

	responseSet []*moduleResponse

	moduleResponse struct {
		Type string `json:"type"`
		Rel  string `json:"rel"`
		Href string `json:"href"`
	}
)

func (SyncData) New() *SyncData {
	return &SyncData{}
}

func (ctrl SyncData) ReadExposedAll(ctx context.Context, r *request.SyncDataReadExposedAll) (interface{}, error) {
	// todo - handle paging
	var (
		err  error
		node *types.Node
	)

	if node, err = service.DefaultNode.FindBySharedNodeID(ctx, r.NodeID); err != nil {
		return nil, err
	}

	s, _, err := service.DefaultExposedModule.Find(ctx, types.ExposedModuleFilter{NodeID: node.ID})

	if err != nil {
		return nil, err
	}

	composeModuleList := make(map[uint64][]uint64, len(s))

	err = s.Walk(func(em *types.ExposedModule) error {
		composeModuleList[em.ComposeModuleID] = []uint64{em.ComposeNamespaceID, em.ID}
		return nil
	})

	if err != nil || len(composeModuleList) == 0 {
		return nil, err
	}

	recordQuery := buildLastSyncQuery(r.LastSync)
	responseSet := responseSet{}

	for composeModuleID, idList := range composeModuleList {
		rf := ct.RecordFilter{
			NamespaceID: idList[0],
			ModuleID:    composeModuleID,
			Query:       recordQuery,
			Paging:      filter.Paging{Limit: 1},
		}

		// todo - handle error properly
		if list, _, err := (cs.Record()).Find(ctx, rf); err != nil || len(list) == 0 {
			continue
		}

		// generate url
		moduleResponse := &moduleResponse{
			Type: "GET",
			Rel:  "Federation Module",
			Href: fmt.Sprintf("/nodes/%d/modules/%d/records/", node.ID, idList[1]),
		}

		responseSet = append(responseSet, moduleResponse)
	}

	return listResponse{&responseSet}, nil
}

func (ctrl SyncData) ReadExposed(ctx context.Context, r *request.SyncDataReadExposed) (interface{}, error) {
	var (
		err   error
		em    *types.ExposedModule
		users st.UserSet
	)

	if _, err = service.DefaultNode.FindBySharedNodeID(ctx, r.NodeID); err != nil {
		return nil, err
	}

	if em, err = service.DefaultExposedModule.FindByID(ctx, r.NodeID, r.ModuleID); err != nil {
		return nil, err
	}

	if users, _, err = ss.DefaultUser.Find(ctx, st.UserFilter{}); err != nil {
		return nil, err
	}

	users, _ = users.Filter(func(u *st.User) (bool, error) {
		return strings.Contains(u.Handle, "federation_"), nil
	})

	// TODO - remove once the federation users are doing the persist
	//        on the data sync, for now, it's superuser
	users = append(users, &st.User{ID: 10000000000000000})

	query := buildLastSyncQuery(r.LastSync)
	ignoredUsersQuery := buildIgnoredUsersQuery(users)

	if ignoredUsersQuery != "" {
		if query != "" {
			query += fmt.Sprintf(" AND %s", ignoredUsersQuery)
		} else {
			query = ignoredUsersQuery
		}
	}

	f := ct.RecordFilter{
		ModuleID: em.ComposeModuleID,
		Query:    query,
		Deleted:  filter.StateInclusive,
	}

	if f.Paging, err = filter.NewPaging(r.Limit, r.PageCursor); err != nil {
		return nil, err
	}

	if f.Sorting, err = filter.NewSorting(r.Sort); err != nil {
		return nil, err
	}

	list, f, err := (cs.Record()).Find(ctx, f)

	if err != nil {
		return nil, err
	}

	// do the actual field filtering
	err = list.Walk(filterExposedFields(em))

	if err != nil {
		return nil, err
	}

	return listRecordResponse{
		Set:    &list,
		Filter: &f,
	}, nil
}

func buildLastSyncQuery(ts uint64) string {
	if ts == 0 {
		return ""
	}

	t := time.Unix(int64(ts), 0)

	if t.IsZero() {
		return ""
	}

	return fmt.Sprintf(
		"(updated_at >= '%s' OR created_at >= '%s' OR deleted_at >= '%s')",
		t.UTC().Format(time.RFC3339),
		t.UTC().Format(time.RFC3339),
		t.UTC().Format(time.RFC3339))
}

func buildIgnoredUsersQuery(users st.UserSet) string {
	query := ""

	if len(users) == 0 {
		return query
	}

	for i, u := range users {
		query += fmt.Sprintf("createdBy != %d", u.ID)

		if i < len(users)-1 {
			query += " AND "
		}
	}

	return fmt.Sprintf("(%s)", query)
}

// filterExposedFields omits the fields that are not exposed as defined
// in the exposed module definition in store
func filterExposedFields(em *types.ExposedModule) func(r *ct.Record) error {
	return func(r *ct.Record) error {
		var err error
		r.Values, err = r.Values.Filter(func(rv *ct.RecordValue) (bool, error) {
			return em.Fields.HasField(rv.Name)
		})

		return err
	}
}
