package service

import (
	"context"
	"fmt"
	"sort"
	"time"

	automationTypes "github.com/crusttech/human/server/automation/types"
	composeTypes "github.com/crusttech/human/server/compose/types"
	"github.com/crusttech/human/server/pkg/actionlog"
	"github.com/crusttech/human/server/pkg/rbac"
	"github.com/crusttech/human/server/store"
	"github.com/crusttech/human/server/system/types"
)

type (
	statistics struct {
		ac        statsAccessControl
		actionlog actionlog.Recorder
		store     store.Storer
	}

	statsAccessControl interface {
		CanReadActionLog(context.Context) bool
		can(ctx context.Context, op string, res rbac.Resource) bool
	}

	// StatisticsRequest is what the dashboard asks for: a range and a bucket size
	StatisticsRequest struct {
		From   *time.Time
		To     *time.Time
		Bucket string
	}

	// statsPermission names the operation that reveals one section of the payload
	statsPermission struct {
		op  string
		res rbac.Resource
	}
)

const (
	statsDefaultRange = 30 * 24 * time.Hour
	statsMaxRange     = 3 * 366 * 24 * time.Hour
	statsDayFormat    = "2006-01-02"

	// ranges up to these lengths bucket by day / week; longer ones by month
	statsDayBucketMax  = 45 * 24 * time.Hour
	statsWeekBucketMax = 200 * 24 * time.Hour
)

var (
	systemComponent     = &types.Component{}
	automationComponent = &automationTypes.Component{}
	composeComponent    = &composeTypes.Component{}

	// every inventory resource is shown to whoever may search it
	statsResourcePermissions = map[string]statsPermission{
		types.SystemStatsUsers:        {"users.search", systemComponent},
		types.SystemStatsRoles:        {"roles.search", systemComponent},
		types.SystemStatsApplications: {"applications.search", systemComponent},
		types.SystemStatsAuthClients:  {"auth-clients.search", systemComponent},
		types.SystemStatsAgents:       {"agents.search", systemComponent},
		types.SystemStatsChatbots:     {"chatbots.search", systemComponent},
		types.SystemStatsProjects:     {"projects.search", systemComponent},
		types.SystemStatsWorkflows:    {"workflows.search", automationComponent},
		types.SystemStatsTaqs:         {"ng-automations.search", automationComponent},
		types.SystemStatsNamespaces:   {"namespaces.search", composeComponent},
		types.SystemStatsModules:      {"namespaces.search", composeComponent},
	}

	statsSessionOutcomes = []string{"completed", "failed", "canceled", "started", "prompted", "suspended"}
	statsTaqOutcomes     = []string{"completed", "failed", "cancelled"}
)

func Statistics() *statistics {
	return &statistics{
		ac:        DefaultAccessControl,
		actionlog: DefaultActionlog,
		store:     DefaultStore,
	}
}

// Metrics answers the admin dashboard: one store aggregate, rolled into the
// requested buckets, with every section the caller may not see left out.
func (svc statistics) Metrics(ctx context.Context, req StatisticsRequest) (rval *types.SystemStats, err error) {
	err = func() error {
		visible := svc.visibleResources(ctx)
		canAutomation := svc.ac.can(ctx, "workflows.search", automationComponent) ||
			svc.ac.can(ctx, "ng-automations.search", automationComponent)
		canActivity := svc.ac.CanReadActionLog(ctx)
		canSignins := svc.ac.can(ctx, "users.search", systemComponent)

		if len(visible) == 0 && !canAutomation && !canActivity && !canSignins {
			return StatisticsErrNotAllowedToReadStatistics()
		}

		r, bucket, err := statsResolveRange(req, time.Now())
		if err != nil {
			return err
		}

		raw, err := store.SystemStats(ctx, svc.store, r)
		if err != nil {
			return err
		}

		buckets := statsBuckets(r, bucket)

		rval = &types.SystemStats{
			Range: types.SystemStatsRangeInfo{
				From:    r.From,
				To:      r.To,
				Bucket:  bucket,
				Buckets: buckets.labels,
			},
			Resources: make(map[string]*types.SystemStatsResource, len(visible)),
		}

		for _, name := range visible {
			rr := raw.Resources[name]
			if rr == nil {
				continue
			}

			out := &types.SystemStatsResource{
				Status:  rr.Status,
				Created: buckets.roll(rr.Created)[""],
			}
			for _, c := range rr.Status {
				out.Total += c
			}
			for _, d := range rr.Created {
				out.CreatedInRange += d.Count
			}
			if out.Created == nil {
				out.Created = buckets.empty()
			}

			rval.Resources[name] = out
		}

		if canAutomation {
			rval.Automation = &types.SystemStatsAutomation{}

			if svc.ac.can(ctx, "workflows.search", automationComponent) {
				rval.Automation.Workflows = statsRuns(buckets, raw.WorkflowRuns, raw.WorkflowRunTotals, statsSessionOutcomes, raw.WorkflowFailures)
			}

			if svc.ac.can(ctx, "ng-automations.search", automationComponent) {
				rval.Automation.Taqs = statsRuns(buckets, raw.TaqRuns, raw.TaqRunTotals, statsTaqOutcomes, raw.TaqFailures)
			}
		}

		if canActivity {
			series := buckets.roll(raw.Activity)
			rval.Activity = &types.SystemStatsActivity{
				Total:  raw.ActivityTotals["ok"] + raw.ActivityTotals["error"],
				Errors: raw.ActivityTotals["error"],
				Series: map[string][]uint{
					"total":  buckets.sum(series["ok"], series["error"]),
					"errors": buckets.orEmpty(series["error"]),
				},
				RecentErrors: raw.RecentErrors,
				ByResource:   raw.TopResources,
			}
		}

		if canSignins {
			rval.Signins = &types.SystemStatsSignins{
				Total:  raw.SigninTotal,
				Users:  raw.SigninUsers,
				Live:   raw.LiveSessions,
				Series: buckets.orEmpty(buckets.roll(raw.Signins)[""]),
			}
		}

		return nil
	}()

	return rval, svc.recordAction(ctx, &statisticsActionProps{}, StatisticsActionServe, err)
}

// Detail answers the dashboard's drill-down for one inventory resource.
func (svc statistics) Detail(ctx context.Context, resource string, req StatisticsRequest) (rval *types.SystemStatsDetail, err error) {
	err = func() error {
		p, ok := statsResourcePermissions[resource]
		if !ok {
			return fmt.Errorf("unknown inventory resource %q", resource)
		}

		if !svc.ac.can(ctx, p.op, p.res) {
			return StatisticsErrNotAllowedToReadStatistics()
		}

		r, bucket, err := statsResolveRange(req, time.Now())
		if err != nil {
			return err
		}

		raw, err := store.SystemStatsResourceDetail(ctx, svc.store, resource, r)
		if err != nil {
			return err
		}

		buckets := statsBuckets(r, bucket)
		rolled := buckets.roll(raw.Movement)

		rval = &types.SystemStatsDetail{
			Resource: resource,
			Range: types.SystemStatsRangeInfo{
				From:    r.From,
				To:      r.To,
				Bucket:  bucket,
				Buckets: buckets.labels,
			},
			Status:  raw.Status,
			Series:  make(map[string][]uint, 3),
			InRange: make(map[string]uint, 3),
			Recent:  raw.Recent,
		}

		for _, c := range raw.Status {
			rval.Total += c
		}

		for _, key := range []string{types.SystemStatsCreated, types.SystemStatsUpdated, types.SystemStatsDeleted} {
			rval.Series[key] = buckets.orEmpty(rolled[key])
			for _, n := range rval.Series[key] {
				rval.InRange[key] += n
			}
		}

		return nil
	}()

	return rval, svc.recordAction(ctx, &statisticsActionProps{}, StatisticsActionServe, err)
}

// visibleResources lists the inventory resources the caller may search, in payload order.
func (svc statistics) visibleResources(ctx context.Context) (out []string) {
	for name, p := range statsResourcePermissions {
		if svc.ac.can(ctx, p.op, p.res) {
			out = append(out, name)
		}
	}

	sort.Strings(out)
	return
}

func statsRuns(b statsBucketSet, daily []types.SystemStatsDaily, totals map[string]uint, outcomes []string, failures interface{}) *types.SystemStatsRuns {
	rolled := b.roll(daily)

	out := &types.SystemStatsRuns{
		ByStatus: make(map[string]uint, len(outcomes)),
		Series:   make(map[string][]uint, len(outcomes)),
		Failures: failures,
	}

	for _, o := range outcomes {
		out.ByStatus[o] = totals[o]
		out.Series[o] = b.orEmpty(rolled[o])
	}

	for _, c := range totals {
		out.Total += c
	}

	return out
}

// statsResolveRange fills the defaults in and picks the bucket size.
func statsResolveRange(req StatisticsRequest, now time.Time) (r types.SystemStatsRange, bucket string, err error) {
	r.To = now
	if req.To != nil {
		r.To = *req.To
	}

	r.From = r.To.Add(-statsDefaultRange)
	if req.From != nil {
		r.From = *req.From
	}

	if !r.From.Before(r.To) {
		return r, "", fmt.Errorf("statistics range is empty: from must be before to")
	}

	if r.To.Sub(r.From) > statsMaxRange {
		return r, "", fmt.Errorf("statistics range too long: at most %d days", int(statsMaxRange.Hours()/24))
	}

	bucket = req.Bucket
	switch bucket {
	case "":
		span := r.To.Sub(r.From)
		switch {
		case span <= statsDayBucketMax:
			bucket = types.SystemStatsBucketDay
		case span <= statsWeekBucketMax:
			bucket = types.SystemStatsBucketWeek
		default:
			bucket = types.SystemStatsBucketMonth
		}
	case types.SystemStatsBucketDay, types.SystemStatsBucketWeek, types.SystemStatsBucketMonth:
	default:
		return r, "", fmt.Errorf("unknown statistics bucket %q", bucket)
	}

	return r, bucket, nil
}

// statsBucketSet is the ordered list of bucket start days a range rolls into.
type statsBucketSet struct {
	labels []string
	index  map[string]int
	size   string
}

// statsBuckets lists the bucket that contains From, then every following
// bucket that starts before To. Weeks start on Monday; days are taken in
// the server's local zone, which is also what the store's date truncation uses.
func statsBuckets(r types.SystemStatsRange, size string) statsBucketSet {
	b := statsBucketSet{index: make(map[string]int), size: size}

	cur := statsBucketStart(r.From.Local(), size)
	end := r.To.Local()

	for cur.Before(end) {
		label := cur.Format(statsDayFormat)
		b.index[label] = len(b.labels)
		b.labels = append(b.labels, label)
		cur = statsBucketNext(cur, size)
	}

	return b
}

func statsBucketStart(t time.Time, size string) time.Time {
	y, m, d := t.Date()
	day := time.Date(y, m, d, 0, 0, 0, 0, t.Location())

	switch size {
	case types.SystemStatsBucketWeek:
		offset := (int(day.Weekday()) + 6) % 7
		return day.AddDate(0, 0, -offset)
	case types.SystemStatsBucketMonth:
		return time.Date(y, m, 1, 0, 0, 0, 0, t.Location())
	default:
		return day
	}
}

func statsBucketNext(t time.Time, size string) time.Time {
	switch size {
	case types.SystemStatsBucketWeek:
		return t.AddDate(0, 0, 7)
	case types.SystemStatsBucketMonth:
		return t.AddDate(0, 1, 0)
	default:
		return t.AddDate(0, 0, 1)
	}
}

// roll sums daily rows into buckets, one series per key. Days that fall
// outside the bucket list are dropped.
func (b statsBucketSet) roll(daily []types.SystemStatsDaily) map[string][]uint {
	out := make(map[string][]uint)

	for _, d := range daily {
		day, err := time.ParseInLocation(statsDayFormat, d.Day, time.Local)
		if err != nil {
			continue
		}

		i, ok := b.index[statsBucketStart(day, b.size).Format(statsDayFormat)]
		if !ok {
			continue
		}

		if out[d.Key] == nil {
			out[d.Key] = b.empty()
		}

		out[d.Key][i] += d.Count
	}

	return out
}

func (b statsBucketSet) empty() []uint {
	return make([]uint, len(b.labels))
}

func (b statsBucketSet) orEmpty(s []uint) []uint {
	if s == nil {
		return b.empty()
	}

	return s
}

func (b statsBucketSet) sum(ss ...[]uint) []uint {
	out := b.empty()
	for _, s := range ss {
		for i := range s {
			out[i] += s[i]
		}
	}

	return out
}
