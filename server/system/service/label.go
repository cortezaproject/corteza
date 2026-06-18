package service

import (
	"context"
	"strconv"

	automationTypes "github.com/crusttech/human/server/automation/types"
	composeTypes "github.com/crusttech/human/server/compose/types"
	"github.com/crusttech/human/server/pkg/filter"
	"github.com/crusttech/human/server/pkg/label/types"
	"github.com/crusttech/human/server/store"
	systemTypes "github.com/crusttech/human/server/system/types"
)

type (
	labelSvc struct {
		store store.Storer
	}
	LabelService interface {
		List(ctx context.Context, filter types.LabelFilter) (types.LabelSet, types.LabelFilter, error)
		Delete(ctx context.Context, name, kind string) error
	}
)

func Label() LabelService {
	return &labelSvc{
		store: DefaultStore,
	}
}

// Delete removes every label entry matching the given name (optionally scoped
// to a resource kind), across all resources.
func (svc labelSvc) Delete(ctx context.Context, name, kind string) (err error) {
	set, _, err := store.SearchLabels(ctx, svc.store, types.LabelFilter{Name: name, Kind: kind})
	if err != nil {
		return err
	}

	if len(set) == 0 {
		return nil
	}

	return store.DeleteLabel(ctx, svc.store, set...)
}

func (svc labelSvc) List(ctx context.Context, f types.LabelFilter) (set types.LabelSet, outF types.LabelFilter, err error) {
	set, outF, err = store.SearchLabels(ctx, svc.store, f)
	if err != nil {
		return nil, outF, err
	}

	set, err = svc.filterActiveResources(ctx, set)
	return set, outF, err
}

// Delete removes every label entry with the given name, optionally restricted
// to a single resource kind. The route is auth-protected; like List, no
// per-resource access control is applied (labels are cross-resource metadata).
func (svc labelSvc) Delete(ctx context.Context, name, kind string) error {
	if name == "" {
		return nil
	}

	set, _, err := store.SearchLabels(ctx, svc.store, types.LabelFilter{Name: name, Kind: kind})
	if err != nil {
		return err
	}
	if len(set) == 0 {
		return nil
	}

	return store.DeleteLabel(ctx, svc.store, set...)
}

// filterActiveResources removes label entries whose resource has been deleted.
// Unknown kinds are passed through unchanged.
func (svc labelSvc) filterActiveResources(ctx context.Context, set types.LabelSet) (types.LabelSet, error) {
	// Group resource IDs by kind
	kindIDs := make(map[string][]uint64)
	for _, l := range set {
		kindIDs[l.Kind] = append(kindIDs[l.Kind], l.ResourceID)
	}

	// Resolve active IDs per kind
	active := make(map[string]map[uint64]bool, len(kindIDs))
	for kind, ids := range kindIDs {
		m, err := svc.activeIDsForKind(ctx, kind, ids)
		if err != nil {
			return nil, err
		}
		active[kind] = m
	}

	// Keep only labels whose resource is active
	out := set[:0]
	for _, l := range set {
		if active[l.Kind][l.ResourceID] {
			out = append(out, l)
		}
	}
	return out, nil
}

func uint64sToStrings(ids []uint64) []string {
	ss := make([]string, len(ids))
	for i, id := range ids {
		ss[i] = strconv.FormatUint(id, 10)
	}
	return ss
}

func (svc labelSvc) activeIDsForKind(ctx context.Context, kind string, ids []uint64) (map[uint64]bool, error) {
	result := make(map[uint64]bool, len(ids))
	strIDs := uint64sToStrings(ids)

	switch kind {
	case "agent":
		rr, _, err := store.SearchAgents(ctx, svc.store, systemTypes.AgentFilter{
			AgentID: strIDs,
			Deleted: filter.StateExcluded,
		})
		if err != nil {
			return nil, err
		}
		for _, r := range rr {
			result[r.ID] = true
		}

	case "executable":
		rr, _, err := store.SearchAutomationNgAutomations(ctx, svc.store, automationTypes.NgAutomationFilter{
			AutomationID: strIDs,
			Deleted:      filter.StateExcluded,
			Disabled:     filter.StateInclusive,
		})
		if err != nil {
			return nil, err
		}
		for _, r := range rr {
			result[r.ID] = true
		}

	case "namespace":
		rr, _, err := store.SearchComposeNamespaces(ctx, svc.store, composeTypes.NamespaceFilter{
			NamespaceID: strIDs,
			Deleted:     filter.StateExcluded,
		})
		if err != nil {
			return nil, err
		}
		for _, r := range rr {
			result[r.ID] = true
		}

	default:
		// Unknown kind — don't filter, include all
		for _, id := range ids {
			result[id] = true
		}
	}

	return result, nil
}