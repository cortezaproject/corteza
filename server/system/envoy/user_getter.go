package envoy

import (
	"context"
	"fmt"
	"strings"

	"github.com/cortezaproject/corteza/server/pkg/envoyx"
	"github.com/cortezaproject/corteza/server/store"
	"github.com/cortezaproject/corteza/server/system/types"
	"github.com/spf13/cast"
)

type (
	// UserGetter is a utility struct to resolve user references from
	// different parts of the system such as the dep graph and the database
	UserGetter struct {
		depGraph *envoyx.DepGraph

		store      store.Storer
		baseFilter types.UserFilter

		// resolved refs, including misses; imports resolve the same handful
		// of refs over and over so this saves a query per value
		cache map[string]uint64
	}
)

func MakeUserGetter(s store.Storer, tt envoyx.Traverser) (g *UserGetter) {
	g = &UserGetter{
		store: s,
		cache: make(map[string]uint64),
	}

	g.baseFilter = types.UserFilter{}

	return
}

// Resolve returns the user ID for the provided reference
//
// If the user can not be resolved, 0 is returned.
func (g *UserGetter) Resolve(ctx context.Context, ref any) (out uint64, err error) {
	// Try to get from datasource
	if g.depGraph != nil {
		out, err = g.getDS(ref)
		if err != nil {
			return
		}
	}
	if out > 0 {
		return
	}

	// Try to get from the database
	out, err = g.getDB(ctx, ref)
	return
}

func (g *UserGetter) getDS(ref any) (out uint64, err error) {
	n := g.depGraph.NodeForRef(envoyx.Ref{
		ResourceType: types.UserResourceType,
		Identifiers:  envoyx.MakeIdentifiers(ref),
	})

	if n == nil {
		return
	}

	return n.Resource.GetID(), nil
}

// @todo this can be improved by prefetching and indexing refs
func (g *UserGetter) getDB(ctx context.Context, ref any) (out uint64, err error) {
	aux := strings.TrimSpace(cast.ToString(ref))

	// Nothing to resolve; without this the query below would match any user
	if aux == "" {
		return 0, nil
	}

	// Numeric references are user IDs (record exports write IDs); keep them
	// as-is instead of reinterpreting them as a name/email query — the query
	// could substring-match an unrelated user, and unresolvable IDs (deleted
	// or purged users) round-tripped verbatim before user refs were resolved
	// through this getter.
	if id := cast.ToUint64(aux); id > 0 {
		return id, nil
	}

	if out, ok := g.cache[aux]; ok {
		return out, nil
	}
	defer func() {
		if err == nil {
			g.cache[aux] = out
		}
	}()

	f := g.baseFilter
	// @todo expand this
	f.Query = aux

	set, _, err := store.SearchUsers(ctx, g.store, f)
	if err != nil {
		return
	}

	if len(set) == 0 {
		return
	}

	if len(set) > 1 {
		err = fmt.Errorf("ambiguous reference %v: matches more then one user", ref)
		return
	}

	return set[0].ID, nil
}
