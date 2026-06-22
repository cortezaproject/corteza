package service

import (
	"context"
	"encoding/json"
	"fmt"
	"hash/fnv"
	"regexp"
	"sort"
	"strings"

	"go.uber.org/zap"

	"github.com/crusttech/human/server/pkg/actionlog"
	a "github.com/crusttech/human/server/pkg/auth"
	"github.com/crusttech/human/server/pkg/errors"
	"github.com/crusttech/human/server/pkg/filter"
	"github.com/crusttech/human/server/pkg/label"
	"github.com/crusttech/human/server/store"
	"github.com/crusttech/human/server/system/service/appstore"
	"github.com/crusttech/human/server/system/types"
)

type (
	connection struct {
		actionlog            actionlog.Recorder
		store                store.Storer
		ac                   connectionAccessController
		configuredConnection *configuredConnection
		catalog              appstore.Client
		logger               *zap.Logger
	}

	connectionAccessController interface {
		CanSearchConnections(ctx context.Context) bool

		CanCreateConnection(context.Context) bool
		CanReadConnection(context.Context, *types.Connection) bool
		CanUpdateConnection(context.Context, *types.Connection) bool
		CanDeleteConnection(context.Context, *types.Connection) bool
	}

	ConnectionService interface {
		FindByID(ctx context.Context, ID uint64) (*types.Connection, error)
		Create(ctx context.Context, new *types.Connection) (*types.Connection, error)
		Update(ctx context.Context, upd *types.Connection) (*types.Connection, error)
		DeleteByID(ctx context.Context, ID uint64) error
		UndeleteByID(ctx context.Context, ID uint64) error
		Search(ctx context.Context, filter types.ConnectionFilter) (types.ConnectionSet, types.ConnectionFilter, error)
		Import(ctx context.Context, catalogID string) (*types.Connection, error)
		Enable(ctx context.Context, ID uint64) (*types.Connection, error)
	}

	ConfiguredConnectionService interface {
		FindByID(ctx context.Context, ID uint64) (*types.ConfiguredConnection, error)
		Search(ctx context.Context, filter types.ConfiguredConnectionFilter) (types.ConfiguredConnectionSet, types.ConfiguredConnectionFilter, error)
		DeleteByID(ctx context.Context, ID uint64) error
		Enable(ctx context.Context, ID uint64) (*types.ConfiguredConnection, error)
	}

	dispatchRsp struct {
		dalConnection *types.DalConnection
	}

	// catalogArb holds per-page catalog pagination state encoded in PagingCursor.Arb.
	// We fetch ALL catalog summaries on every request and filter by the sort window
	// defined by the last DB record of the previous page, so catalog entries always
	// appear exactly once on the page where their sort position falls.
	catalogArb struct {
		LastSortKey string `json:"lastSortKey,omitempty"`
		SortCol     string `json:"sortCol,omitempty"`
		SortDesc    bool   `json:"sortDesc,omitempty"`
	}
)

// catalogIDToSyntheticID converts a catalog string ID to a stable uint64
// by hashing with FNV-64a and setting bit 63 to avoid collision with DB IDs.
func catalogIDToSyntheticID(catalogID string) uint64 {
	h := fnv.New64a()
	h.Write([]byte(catalogID))
	return h.Sum64() | (1 << 63)
}

// catalogConnectionToLocal converts a full appstore connection to a local Connection struct
// with a synthetic ID, source=catalog, and full config populated.
func catalogConnectionToLocal(conn *appstore.Connection) *types.Connection {
	c := &types.Connection{
		ID:        catalogIDToSyntheticID(conn.ID),
		Handle:    conn.Handle,
		CatalogID: conn.ID,
		Source:    "catalog",
		Status:    "draft",
		Meta: types.ConnectionMeta{
			Short:       conn.Meta.Short,
			Description: conn.Meta.Description,
			Icon:        conn.Meta.Icon,
			Tags:        conn.Meta.Tags,
		},
	}
	if conn.Service != nil {
		raw, _ := json.Marshal(conn.Service)
		_ = json.Unmarshal(raw, &c.Service)
	}
	if conn.Operations != nil {
		raw, _ := json.Marshal(conn.Operations)
		_ = json.Unmarshal(raw, &c.Operations)
	}
	if conn.Resources != nil {
		raw, _ := json.Marshal(conn.Resources)
		_ = json.Unmarshal(raw, &c.Resources)
	}
	return c
}

func Connection() *connection {
	return &connection{
		actionlog: DefaultActionlog,
		store:     DefaultStore,
		ac:        DefaultAccessControl,
		logger:    DefaultLogger.Named("connection"),
	}
}

func (svc *connection) WithConfiguredConnection(cc *configuredConnection) *connection {
	svc.configuredConnection = cc
	return svc
}

func (svc *connection) WithCatalog(c appstore.Client) *connection {
	svc.catalog = c
	return svc
}

func (svc *connection) FindByID(ctx context.Context, ID uint64) (res *types.Connection, err error) {
	var (
		aProps = &connectionActionProps{connection: &types.Connection{ID: ID}}
	)

	err = func() error {
		res, err = loadConnection(ctx, svc.store, ID)
		if err != nil {
			// High bit set → synthetic catalog ID. Try catalog lookup.
			if svc.catalog != nil && ID&(1<<63) != 0 {
				res, err = svc.findByCatalogSyntheticID(ctx, ID)
				if err != nil {
					return err
				}
				aProps.setConnection(res)
				svc.deriveParams(res)
				return nil
			}
			return err
		}

		aProps.setConnection(res)

		if !svc.ac.CanReadConnection(ctx, res) {
			return ConnectionErrNotAllowedToRead()
		}

		if err = label.Load(ctx, svc.store, res); err != nil {
			return err
		}

		if res.Source == "catalog" {
			res.Status = "active"
		}

		svc.deriveParams(res)
		return nil
	}()

	return res, svc.recordAction(ctx, aProps, ConnectionActionLookup, err)
}

// findByCatalogSyntheticID pages through the catalog until it finds the entry
// whose catalogIDToSyntheticID hash matches id, then fetches the full config.
func (svc *connection) findByCatalogSyntheticID(ctx context.Context, id uint64) (*types.Connection, error) {
	const pageSize = 50
	for page := 1; ; page++ {
		summaries, err := svc.catalog.ListPage(ctx, page, pageSize)
		if err != nil {
			return nil, err
		}
		for _, s := range summaries {
			if catalogIDToSyntheticID(s.ID) == id {
				conn, err := svc.catalog.GetConnection(ctx, s.ID)
				if err != nil {
					return nil, err
				}
				return catalogConnectionToLocal(conn), nil
			}
		}
		if len(summaries) < pageSize {
			break
		}
	}
	return nil, ConnectionErrNotFound()
}

// beforeCreate runs after the create RBAC check and before id/timestamps are
// assigned. It normalises + validates the incoming connection, enforces handle
// uniqueness and fills in the audit/default fields the generated body does not.
func (svc *connection) beforeCreate(ctx context.Context, new *types.Connection) (err error) {
	svc.preprocess(new)

	if err = svc.validateConnection(new); err != nil {
		return err
	}

	if err = svc.uniqueHandleCheck(ctx, new.Handle, 0); err != nil {
		return err
	}

	new.CreatedBy = a.GetIdentityFromContext(ctx).Identity()
	new.Revision = 1

	if new.Status == "" {
		new.Status = "draft"
	}

	if new.Source == "" {
		new.Source = "local"
	}

	return nil
}

// afterCreate derives the read-only DerivedParams on the freshly stored record.
func (svc *connection) afterCreate(ctx context.Context, res *types.Connection) error {
	svc.deriveParams(res)
	return nil
}

// beforeDelete rejects deletion while active ConfiguredConnections reference the
// connection.
func (svc *connection) beforeDelete(ctx context.Context, res *types.Connection) error {
	cc, _, err := store.SearchConfiguredConnections(ctx, svc.store, types.ConfiguredConnectionFilter{
		ConnectionID: res.ID,
	})
	if err != nil {
		return err
	}
	if len(cc) > 0 {
		return ConnectionErrHasActiveConnections()
	}

	// stamp the deleter (the generated soft-delete only sets DeletedAt)
	res.DeletedBy = a.GetIdentityFromContext(ctx).Identity()

	return nil
}

// beforeUndelete clears the deleter; the generated undelete only resets DeletedAt.
func (svc *connection) beforeUndelete(ctx context.Context, res *types.Connection) error {
	res.DeletedBy = 0
	return nil
}

func (svc *connection) Update(ctx context.Context, upd *types.Connection) (res *types.Connection, err error) {
	var (
		aProps = &connectionActionProps{update: upd}
		old    *types.Connection
	)

	err = func() (err error) {
		if res, err = loadConnection(ctx, svc.store, upd.ID); err != nil {
			return err
		}
		old = res.Clone()

		aProps.setConnection(res)

		if !svc.ac.CanUpdateConnection(ctx, res) {
			return ConnectionErrNotAllowedToUpdate()
		}

		svc.preprocess(upd)

		if err = svc.validateConnection(upd); err != nil {
			return err
		}

		if upd.Handle != res.Handle {
			if err = svc.uniqueHandleCheck(ctx, upd.Handle, upd.ID); err != nil {
				return err
			}
		}

		res.Handle = upd.Handle
		res.Meta = upd.Meta
		res.Service = upd.Service
		res.Resources = upd.Resources
		res.Operations = upd.Operations
		res.Revision++

		n := now()
		res.UpdatedAt = n
		res.UpdatedBy = a.GetIdentityFromContext(ctx).Identity()

		if err = store.UpdateConnection(ctx, svc.store, res); err != nil {
			return err
		}

		if label.Changed(res.Labels, upd.Labels) {
			if err = label.Update(ctx, svc.store, upd); err != nil {
				return err
			}
			res.Labels = upd.Labels
		}

		svc.deriveParams(res)

		return nil
	}()

	return res, svc.recordAction(ctx, aProps, ConnectionActionUpdate, err, old, res)
}

func (svc *connection) Enable(ctx context.Context, ID uint64) (res *types.Connection, err error) {
	var (
		aProps = &connectionActionProps{connection: &types.Connection{ID: ID}}
		old    *types.Connection
	)

	err = func() (err error) {
		if res, err = loadConnection(ctx, svc.store, ID); err != nil {
			return err
		}
		old = res.Clone()

		aProps.setConnection(res)

		if !svc.ac.CanUpdateConnection(ctx, res) {
			return ConnectionErrNotAllowedToUpdate()
		}

		if err = svc.validateConnectionEnable(res); err != nil {
			return err
		}

		res.Status = "active"
		n := now()
		res.UpdatedAt = n
		res.UpdatedBy = a.GetIdentityFromContext(ctx).Identity()

		if err = store.UpdateConnection(ctx, svc.store, res); err != nil {
			return err
		}

		return nil
	}()

	return res, svc.recordAction(ctx, aProps, ConnectionActionUpdate, err, old, res)
}

// matchesQuery checks if name matches query (case-insensitive substring).
func matchesQuery(query, name string) bool {
	return strings.Contains(strings.ToLower(name), strings.ToLower(query))
}

func (svc *connection) Search(ctx context.Context, filter types.ConnectionFilter) (set types.ConnectionSet, f types.ConnectionFilter, err error) {
	var (
		aProps = &connectionActionProps{filter: &filter}
	)

	filter.Check = func(res *types.Connection) (bool, error) {
		if !svc.ac.CanReadConnection(ctx, res) {
			return false, nil
		}
		if filter.Query != "" && !matchesQuery(filter.Query, res.Meta.Short) {
			return false, nil
		}
		return true, nil
	}

	err = func() error {
		if !svc.ac.CanSearchConnections(ctx) {
			return ConnectionErrNotAllowedToSearch()
		}

		set, f, err = store.SearchConnections(ctx, svc.store, filter)
		if err != nil {
			return err
		}

		// Merge with appstore catalog when configured.
		// DB records keep their stored Source value ("catalog" or "local").
		// Catalog-only entries (not yet imported) are injected into the correct
		// sort position by fetching ALL catalog summaries and filtering to the
		// sort window defined by the current DB page.
		var catalogOnly []*types.Connection
		if svc.catalog != nil {
			var arb catalogArb
			if filter.PageCursor != nil {
				_ = filter.PageCursor.GetArb(&arb)
			}

			// Determine primary sort column and direction.
			sortCol := "name"
			sortDesc := false
			if len(filter.Sort) > 0 {
				sortCol = strings.ToLower(filter.Sort[0].Column)
				sortDesc = filter.Sort[0].Descending
			}

			// If sort changed, discard stale lower bound.
			if arb.SortCol != "" && arb.SortCol != sortCol {
				arb.LastSortKey = ""
			}

			// Fetch ALL catalog summaries (bounded to 1000 entries).
			const maxPages = 20
			const pageSize = 50
			var allSummaries []appstore.ConnectionSummary
			catFetchOK := true
			for page := 1; page <= maxPages; page++ {
				summaries, catErr := svc.catalog.ListPage(ctx, page, pageSize)
				if catErr != nil {
					svc.logger.Warn("appstore unavailable, serving DB-only results", zap.Error(catErr))
					catFetchOK = false
					break
				}
				allSummaries = append(allSummaries, summaries...)
				if len(summaries) < pageSize {
					break
				}
			}

			if catFetchOK {
				// Index DB records by handle; link catalog IDs for matching DB entries.
				dbByHandle := make(map[string]*types.Connection, len(set))
				for _, c := range set {
					dbByHandle[c.Handle] = c
				}
				for _, s := range allSummaries {
					if dbConn, ok := dbByHandle[s.Handle]; ok {
						dbConn.CatalogID = s.ID
					}
				}

				// Build catalog-only entries (handle not in DB).
				for _, s := range allSummaries {
					if _, ok := dbByHandle[s.Handle]; !ok {
						// Filter by query if specified
						if filter.Query != "" && !matchesQuery(filter.Query, s.Short) {
							continue
						}
						catalogOnly = append(catalogOnly, &types.Connection{
							ID:     catalogIDToSyntheticID(s.ID),
							Handle: s.Handle,
							Status: "draft",
							Meta: types.ConnectionMeta{
								Short:       s.Short,
								Description: s.Description,
							},
							Source:    "catalog",
							CatalogID: s.ID,
						})
					}
				}

				// Sort window: include catalog entries whose sort key falls
				// between the previous page's last DB key and this page's last DB key.
				//   ascending:  prevKey < entryKey <= lastDBKey  (or no upper bound on last page)
				//   descending: prevKey > entryKey >= lastDBKey  (or no upper bound on last page)
				isLastDBPage := f.NextPage == nil
				var lastDBSortKey string
				if len(set) > 0 {
					lastDBSortKey = connectionSortKey(set[len(set)-1], sortCol)
				}

				isBefore := func(a, b string) bool {
					if sortDesc {
						return a > b
					}
					return a < b
				}

				for _, c := range catalogOnly {
					sk := connectionSortKey(c, sortCol)
					afterPrev := arb.LastSortKey == "" || isBefore(arb.LastSortKey, sk)
					withinBound := isLastDBPage || lastDBSortKey == "" || !isBefore(lastDBSortKey, sk)
					if afterPrev && withinBound {
						set = append(set, c)
					}
				}

				// Persist last DB sort key so the next page knows its lower bound.
				if f.NextPage != nil {
					nextArb := catalogArb{
						LastSortKey: lastDBSortKey,
						SortCol:     sortCol,
						SortDesc:    sortDesc,
					}
					_ = f.NextPage.SetArb(nextArb)
				}
			}
		}

		// Apply source filter after merge.
		if filter.Source != "" {
			filtered := set[:0]
			for _, c := range set {
				if c.Source == filter.Source {
					filtered = append(filtered, c)
				}
			}
			set = filtered
		}

		// Adjust total to include catalog-only entries not yet imported into DB.
		// The store counts only DB records; catalog-only entries must be added.
		// When filtering by source="local", catalog entries are excluded so no adjustment needed.
		if filter.IncTotal && svc.catalog != nil && filter.Source != "local" {
			f.Total += uint(len(catalogOnly))
		}

		// Re-sort the merged set so catalog-appended entries land in the right position.
		if svc.catalog != nil {
			sortConnectionSet(set, filter.Sort)
		}

		// Count installed ConfiguredConnections per DB record.
		var dbIDs []uint64
		for _, c := range set {
			if c.ID != 0 && c.ID&(1<<63) == 0 {
				dbIDs = append(dbIDs, c.ID)
			}
		}
		if len(dbIDs) > 0 {
			counts, countErr := svc.countConfiguredByIDs(ctx, dbIDs)
			if countErr != nil {
				svc.logger.Warn("could not count configured connections", zap.Error(countErr))
			} else {
				for _, c := range set {
					c.InstalledCount = counts[c.ID]
				}
			}
		}

		// Load labels only for records that exist in DB.
		if err = label.Load(ctx, svc.store, toLabeledConnections(filterDBConnections(set))...); err != nil {
			return err
		}

		for _, c := range set {
			svc.deriveParams(c)
		}
		return nil
	}()

	return set, f, svc.recordAction(ctx, aProps, ConnectionActionSearch, err)
}

// Import fetches a connection definition from the appstore by its catalog ID and
// creates it locally. If a connection with the same handle already exists the
// existing record is returned without error.
func (svc *connection) Import(ctx context.Context, catalogID string) (res *types.Connection, err error) {
	if svc.catalog == nil {
		return nil, fmt.Errorf("appstore catalog is not configured")
	}

	catalogConn, err := svc.catalog.GetConnection(ctx, catalogID)
	if err != nil {
		return nil, err
	}

	// Check for existing record with the same handle — idempotent.
	existing, _, lookupErr := store.SearchConnections(ctx, svc.store, types.ConnectionFilter{Handle: catalogConn.Handle})
	if lookupErr != nil {
		return nil, lookupErr
	}
	if len(existing) > 0 {
		return existing[0], nil
	}

	conn := &types.Connection{
		Handle:    catalogConn.Handle,
		CatalogID: catalogConn.ID,
		Source:    "catalog",
		Status:    "active",
		Meta: types.ConnectionMeta{
			Short:       catalogConn.Meta.Short,
			Description: catalogConn.Meta.Description,
			Icon:        catalogConn.Meta.Icon,
			Tags:        catalogConn.Meta.Tags,
		},
	}

	if catalogConn.Service != nil {
		raw, _ := json.Marshal(catalogConn.Service)
		_ = json.Unmarshal(raw, &conn.Service)
	}
	if catalogConn.Operations != nil {
		raw, _ := json.Marshal(catalogConn.Operations)
		_ = json.Unmarshal(raw, &conn.Operations)
	}
	if catalogConn.Resources != nil {
		raw, _ := json.Marshal(catalogConn.Resources)
		_ = json.Unmarshal(raw, &conn.Resources)
	}

	return svc.Create(ctx, conn)
}

// countConfiguredByIDs returns a map of connectionID → count of ConfiguredConnections.
func (svc *connection) countConfiguredByIDs(ctx context.Context, ids []uint64) (map[uint64]int, error) {
	idSet := make(map[uint64]bool, len(ids))
	for _, id := range ids {
		idSet[id] = true
	}

	cc, _, err := store.SearchConfiguredConnections(ctx, svc.store, types.ConfiguredConnectionFilter{})
	if err != nil {
		return nil, err
	}

	counts := make(map[uint64]int, len(ids))
	for _, c := range cc {
		if idSet[c.ConnectionID] {
			counts[c.ConnectionID]++
		}
	}
	return counts, nil
}

// filterDBConnections returns only connections with a real DB record
// (ID != 0 and no synthetic high bit from catalog).
func filterDBConnections(set types.ConnectionSet) types.ConnectionSet {
	out := make(types.ConnectionSet, 0, len(set))
	for _, c := range set {
		if c.ID != 0 && c.ID&(1<<63) == 0 {
			out = append(out, c)
		}
	}
	return out
}

// Configure creates a new ConfiguredConnection in draft status.
// No provisioning is done yet; config can be freely edited.
func (svc *connection) Configure(ctx context.Context, new *types.ConfiguredConnection) (res *types.ConfiguredConnection, err error) {
	return svc.configuredConnection.Create(ctx, new)
}

func (svc *connection) UpdateConfiguration(ctx context.Context, upd *types.ConfiguredConnection) (res *types.ConfiguredConnection, err error) {
	res, err = svc.configuredConnection.Update(ctx, upd)
	if err != nil {
		return nil, err
	}

	return
}

// -- preprocessing -----------------------------------------------------------

func (svc *connection) preprocess(c *types.Connection) {
	svc.normaliseOperationTypes(c)
}

func (svc *connection) normaliseOperationTypes(c *types.Connection) {
	for i := range c.Operations {
		for j := range c.Operations[i].Input {
			if strings.EqualFold(c.Operations[i].Input[j].Type, "Object") {
				c.Operations[i].Input[j].Type = "Vars"
			}
		}

		for j := range c.Operations[i].Output {
			if strings.EqualFold(c.Operations[i].Output[j].Type, "Object") {
				c.Operations[i].Output[j].Type = "Vars"
			}
		}
	}
}

// -- validation ---------------------------------------------------------------

func (svc *connection) validateConnection(c *types.Connection) error {
	if c.Meta.Short == "" {
		return ConnectionErrMissingShortName()
	}
	if c.Service.BaseURL.Value == "" {
		return ConnectionErrMissingBaseURL()
	}
	if c.Service.Auth.Method == "" {
		return ConnectionErrMissingAuthMethod()
	}

	for _, op := range c.Operations {
		if len(op.Steps) == 0 {
			return fmt.Errorf("operation must have at least one step")
		}
	}

	return nil
}

func (svc *connection) validateConnectionEnable(c *types.Connection) error {
	if c.Source != "local" {
		return errors.InvalidData("only locally-built connections can be explicitly enabled")
	}
	if c.Status != "draft" {
		return errors.InvalidData("only draft connections can be enabled")
	}
	return nil
}

func (svc *connection) uniqueHandleCheck(ctx context.Context, handle string, excludeID uint64) error {
	if handle == "" {
		return nil
	}
	set, _, err := store.SearchConnections(ctx, svc.store, types.ConnectionFilter{Handle: handle})
	if err != nil {
		return err
	}
	for _, c := range set {
		if c.ID != excludeID {
			return ConnectionErrHandleNotUnique()
		}
	}
	return nil
}

// -- parameter derivation -----------------------------------------------------

var placeholderRe = regexp.MustCompile(`\{\{\s*(\w+)\s*\}\}`)

// deriveParams scans all Template fields in a Connection, extracts {{placeholder}}
// variables, and populates the read-only DerivedParams field.
func (svc *connection) deriveParams(c *types.Connection) {
	type scopedTemplate struct {
		scope []string
		tpl   types.ConnectionTemplate
	}

	var tt []scopedTemplate

	// Service-level templates
	tt = append(tt, scopedTemplate{[]string{"service", "baseURL"}, c.Service.BaseURL})
	for _, h := range c.Service.Headers {
		tt = append(tt, scopedTemplate{[]string{"service", "headers"}, h})
	}
	for _, p := range c.Service.Auth.Params {
		tt = append(tt, scopedTemplate{[]string{"service", "auth"}, p})
	}

	// Extract unique params (keyed by scope+name)
	seen := make(map[string]bool)
	var params []types.ConnectionDerivedParam

	for _, st := range tt {
		// Build a lookup of inline placeholder metadata
		meta := make(map[string]types.ConnectionPlaceholder, len(st.tpl.Placeholders))
		for _, p := range st.tpl.Placeholders {
			meta[p.Name] = p
		}

		for _, match := range placeholderRe.FindAllStringSubmatch(st.tpl.Value, -1) {
			name := match[1]
			key := name + "|" + joinScope(st.scope)
			if seen[key] {
				continue
			}
			seen[key] = true

			p := types.ConnectionDerivedParam{
				Name:     name,
				Label:    labelFromName(name),
				Scope:    st.scope,
				Type:     "string",
				Required: true,
			}
			if m, ok := meta[name]; ok {
				if m.Type != "" {
					p.Type = m.Type
				}
				if m.Label != "" {
					p.Label = m.Label
				}
				p.Description = m.Description
				p.Required = m.Required
				p.Default = m.Default
				p.Options = m.Options
			}
			params = append(params, p)
		}
	}

	// service.params are explicit declarations (not template references).
	// Append them directly with scope ["service"] so the UI exposes them
	// as connection-level configuration inputs.
	svcScope := []string{"service"}
	for _, sp := range c.Service.Params {
		key := sp.Name + "|" + joinScope(svcScope)
		if seen[key] {
			continue
		}
		seen[key] = true
		dp := types.ConnectionDerivedParam{
			Name:        sp.Name,
			Label:       sp.Label,
			Scope:       svcScope,
			Type:        sp.Type,
			Description: sp.Description,
			Required:    sp.Required,
			Default:     sp.Default,
			Options:     sp.Options,
		}
		if dp.Label == "" {
			dp.Label = labelFromName(sp.Name)
		}
		params = append(params, dp)
	}

	c.DerivedParams = params
}

func joinScope(ss []string) string {
	out := ""
	for i, s := range ss {
		if i > 0 {
			out += "."
		}
		out += s
	}
	return out
}

// labelFromName converts any common identifier style into a human-friendly
// title-cased string. Supported input styles:
//
//	camelCase        → "Camel Case"
//	PascalCase       → "Pascal Case"
//	snake_case       → "Snake Case"
//	kebab-case       → "Kebab Case"
//	SCREAMING_SNAKE  → "Screaming Snake"
//	APIKey / apiID   → "API Key" / "Api ID"  (consecutive-uppercase runs kept as words)
func labelFromName(name string) string {
	runes := []rune(name)
	n := len(runes)
	if n == 0 {
		return ""
	}

	// Split into words first
	var words []string
	start := 0

	upper := func(r rune) bool { return r >= 'A' && r <= 'Z' }
	lower := func(r rune) bool { return r >= 'a' && r <= 'z' }
	sep   := func(r rune) bool { return r == '_' || r == '-' }

	flush := func(end int) {
		if end > start {
			words = append(words, string(runes[start:end]))
		}
	}

	for i := 1; i < n; i++ {
		prev, cur := runes[i-1], runes[i]
		switch {
		case sep(cur):
			flush(i)
			start = i + 1
		case sep(prev):
			start = i
		case upper(cur) && lower(prev):
			// fooBar → foo | Bar
			flush(i)
			start = i
		case upper(cur) && i+1 < n && lower(runes[i+1]) && upper(prev):
			// APIKey → API | Key
			flush(i)
			start = i
		}
	}
	flush(n)

	// Title-case each word; preserve all-uppercase words (acronyms) as-is
	var b strings.Builder
	for i, w := range words {
		if i > 0 {
			b.WriteByte(' ')
		}
		rr := []rune(w)
		allUpper := true
		for _, r := range rr {
			if r >= 'a' && r <= 'z' {
				allUpper = false
				break
			}
		}
		if allUpper && len(rr) > 1 {
			// acronym — keep as-is (e.g. API, ID)
			b.WriteString(w)
		} else {
			// title-case: uppercase first, lowercase rest
			for j, r := range rr {
				if j == 0 && r >= 'a' && r <= 'z' {
					r -= 32
				} else if j > 0 && r >= 'A' && r <= 'Z' {
					r += 32
				}
				b.WriteRune(r)
			}
		}
	}
	return b.String()
}

// connectionSortKey returns the string sort key for a connection given a column name.
func connectionSortKey(c *types.Connection, col string) string {
	switch col {
	case "handle":
		return strings.ToLower(c.Handle)
	case "status":
		return c.Status
	case "source":
		return c.Source
	default: // "name"
		return strings.ToLower(c.Meta.Short)
	}
}

// sortConnectionSet sorts the merged set in-memory using the same columns as the
// DB query so catalog-appended entries land in the correct position.
// Only name, handle, status, and source are supported; unknown columns are ignored.
// When no sort is specified it falls back to name ascending.
func sortConnectionSet(set types.ConnectionSet, ss filter.SortExprSet) {
	if len(ss) == 0 {
		ss = filter.SortExprSet{{Column: "name", Descending: false}}
	}

	sort.SliceStable(set, func(i, j int) bool {
		a, b := set[i], set[j]
		for _, s := range ss {
			var va, vb string
			switch strings.ToLower(s.Column) {
			case "name":
				va, vb = strings.ToLower(a.Meta.Short), strings.ToLower(b.Meta.Short)
			case "handle":
				va, vb = strings.ToLower(a.Handle), strings.ToLower(b.Handle)
			case "status":
				va, vb = a.Status, b.Status
			case "source":
				va, vb = a.Source, b.Source
			default:
				continue
			}
			if va == vb {
				continue
			}
			if s.Descending {
				return va > vb
			}
			return va < vb
		}
		return false
	})
}
