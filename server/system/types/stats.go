package types

import "time"

// The admin dashboard's data. The store answers SystemStatsRaw for a range
// (status counts, per-day rows, recent failures); the statistics service
// rolls the days into buckets and gates each section on the caller's
// permissions, producing SystemStats.

type (
	// SystemStatsRange bounds every series: [From, To).
	SystemStatsRange struct {
		From time.Time
		To   time.Time
	}

	// SystemStatsDaily is one aggregate row: the day (YYYY-MM-DD), the key the
	// count is split on (a status, "total", "error"; "" when there is no split)
	// and the count.
	SystemStatsDaily struct {
		Day   string
		Key   string
		Count uint
	}

	// SystemStatsKeyCount is a label with a count, for rankings and status splits.
	SystemStatsKeyCount struct {
		Key   string `json:"key"`
		Count uint   `json:"count"`
	}

	// SystemStatsRawResource holds one inventory resource: status counts over
	// the whole table (deleted included) and created_at rows within the range.
	SystemStatsRawResource struct {
		Status  map[string]uint
		Created []SystemStatsDaily
	}

	// SystemStatsWorkflowFailure is a failed workflow session within the range.
	SystemStatsWorkflowFailure struct {
		SessionID      uint64    `json:"sessionID,string"`
		WorkflowID     uint64    `json:"workflowID,string"`
		WorkflowHandle string    `json:"workflowHandle"`
		WorkflowName   string    `json:"workflowName"`
		EventType      string    `json:"eventType"`
		Error          string    `json:"error"`
		CreatedAt      time.Time `json:"createdAt"`
	}

	// SystemStatsLogEntry is an action-log row the dashboard lists.
	SystemStatsLogEntry struct {
		ActionID    uint64                 `json:"actionID,string"`
		Timestamp   time.Time              `json:"timestamp"`
		Resource    string                 `json:"resource"`
		Action      string                 `json:"action"`
		Description string                 `json:"description"`
		Error       string                 `json:"error"`
		ActorID     uint64                 `json:"actorID,string"`
		Meta        map[string]interface{} `json:"meta,omitempty"`
	}

	// SystemStatsRaw is the store's answer for one range.
	SystemStatsRaw struct {
		// keyed by inventory resource name (users, roles, workflows, ...)
		Resources map[string]*SystemStatsRawResource

		// workflow sessions created within the range, per day and status name
		WorkflowRuns      []SystemStatsDaily
		WorkflowRunTotals map[string]uint
		WorkflowFailures  []*SystemStatsWorkflowFailure

		// TAQ runs read from the action log, per day and outcome (completed, failed, cancelled)
		TaqRuns      []SystemStatsDaily
		TaqRunTotals map[string]uint
		TaqFailures  []*SystemStatsLogEntry

		// action-log entries within the range, per day, keyed "total" and "error"
		Activity       []SystemStatsDaily
		ActivityTotals map[string]uint
		RecentErrors   []*SystemStatsLogEntry
		TopResources   []SystemStatsKeyCount

		// auth sessions created within the range, per day
		Signins      []SystemStatsDaily
		SigninTotal  uint
		SigninUsers  uint
		LiveSessions uint
	}
)

// Bucket granularities the service rolls days into.
const (
	SystemStatsBucketDay   = "day"
	SystemStatsBucketWeek  = "week"
	SystemStatsBucketMonth = "month"
)

// Inventory resource names; each maps to one store table and one permission.
const (
	SystemStatsUsers        = "users"
	SystemStatsRoles        = "roles"
	SystemStatsApplications = "applications"
	SystemStatsAuthClients  = "authClients"
	SystemStatsAgents       = "agents"
	SystemStatsChatbots     = "chatbots"
	SystemStatsProjects     = "projects"
	SystemStatsWorkflows    = "workflows"
	SystemStatsTaqs         = "taqs"
	SystemStatsNamespaces   = "namespaces"
	SystemStatsModules      = "modules"
)

type (
	// SystemStats is the REST payload.
	SystemStats struct {
		Range      SystemStatsRangeInfo            `json:"range"`
		Resources  map[string]*SystemStatsResource `json:"resources"`
		Automation *SystemStatsAutomation          `json:"automation,omitempty"`
		Activity   *SystemStatsActivity            `json:"activity,omitempty"`
		Signins    *SystemStatsSignins             `json:"signins,omitempty"`
	}

	// SystemStatsRangeInfo describes the buckets every series is aligned to.
	SystemStatsRangeInfo struct {
		From    time.Time `json:"from"`
		To      time.Time `json:"to"`
		Bucket  string    `json:"bucket"`
		Buckets []string  `json:"buckets"`
	}

	// SystemStatsResource is one inventory tile.
	SystemStatsResource struct {
		Total          uint            `json:"total"`
		Status         map[string]uint `json:"status"`
		Created        []uint          `json:"created"`
		CreatedInRange uint            `json:"createdInRange"`
	}

	SystemStatsAutomation struct {
		Workflows *SystemStatsRuns `json:"workflows,omitempty"`
		Taqs      *SystemStatsRuns `json:"taqs,omitempty"`
	}

	// SystemStatsRuns describes runs within the range: totals by outcome, a
	// series per outcome aligned to the buckets, and the latest failures.
	SystemStatsRuns struct {
		Total    uint              `json:"total"`
		ByStatus map[string]uint   `json:"byStatus"`
		Series   map[string][]uint `json:"series"`
		Failures interface{}       `json:"failures"`
	}

	SystemStatsActivity struct {
		Total        uint                   `json:"total"`
		Errors       uint                   `json:"errors"`
		Series       map[string][]uint      `json:"series"`
		RecentErrors []*SystemStatsLogEntry `json:"recentErrors"`
		ByResource   []SystemStatsKeyCount  `json:"byResource"`
	}

	SystemStatsSignins struct {
		Total  uint   `json:"total"`
		Users  uint   `json:"users"`
		Live   uint   `json:"live"`
		Series []uint `json:"series"`
	}
)

type (
	// SystemStatsItem is one row of an inventory resource, as the drill-down lists it.
	SystemStatsItem struct {
		ID        uint64     `json:"id,string"`
		Label     string     `json:"label"`
		Handle    string     `json:"handle"`
		Status    string     `json:"status"`
		CreatedAt time.Time  `json:"createdAt"`
		UpdatedAt *time.Time `json:"updatedAt,omitempty"`
		DeletedAt *time.Time `json:"deletedAt,omitempty"`
	}

	// SystemStatsDetailRaw is the store's answer for one resource: status
	// counts over the whole table, movement rows per day within the range
	// (keys created, updated, deleted) and the newest rows.
	SystemStatsDetailRaw struct {
		Status   map[string]uint
		Movement []SystemStatsDaily
		Recent   []*SystemStatsItem
	}

	// SystemStatsDetail is the drill-down payload.
	SystemStatsDetail struct {
		Resource string               `json:"resource"`
		Range    SystemStatsRangeInfo `json:"range"`
		Total    uint                 `json:"total"`
		Status   map[string]uint      `json:"status"`
		Series   map[string][]uint    `json:"series"`
		InRange  map[string]uint      `json:"inRange"`
		Recent   []*SystemStatsItem   `json:"recent"`
	}
)

// Movement keys of a resource detail.
const (
	SystemStatsCreated = "created"
	SystemStatsUpdated = "updated"
	SystemStatsDeleted = "deleted"
)
