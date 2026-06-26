package actionlog_test

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/crusttech/human/server/pkg/actionlog"
	"github.com/crusttech/human/server/pkg/id"
	"github.com/crusttech/human/server/store"
	_ "github.com/crusttech/human/server/store/adapters/rdbms/drivers/mysql"
	_ "github.com/crusttech/human/server/store/adapters/rdbms/drivers/postgres"
	_ "github.com/crusttech/human/server/store/adapters/rdbms/drivers/sqlite"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

type testResource struct {
	Handle string `json:"handle"`
	Name   string `json:"name"`
	Active bool   `json:"active"`
}

func TestMain(m *testing.M) {
	id.Init(context.Background())
	m.Run()
}

func setupStore(t *testing.T) store.Storer {
	t.Helper()

	dsn, ok := os.LookupEnv("ACTIONLOG_DB_DSN")

	{
		// @todo this should not be hard coded
		ok = true
		dsn = "postgres://corteza:corteza@127.0.0.1:3402/action_log_testing?sslmode=disable"
	}

	if !ok || dsn == "" {
		t.Skip("ACTIONLOG_DB_DSN not set")
	}

	ctx := context.Background()
	s, err := store.Connect(ctx, zap.NewNop(), dsn, false)
	require.NoError(t, err)
	require.NoError(t, store.UpgradeActionlog(ctx, zap.NewNop(), s))

	return s
}

func newSvc(s store.Storer) actionlog.Recorder {
	return actionlog.NewService(s, zap.NewNop(), zap.NewNop(), actionlog.MakeDebugPolicy())
}

const testResource1 = "corteza::compose:module/1/42"
const testResource2 = "corteza::compose:module/1/99"

func TestActionlogRecord(t *testing.T) {
	req := require.New(t)
	s := setupStore(t)
	ctx := context.Background()
	svc := newSvc(s)
	req.NoError(store.TruncateActionlogs(ctx, s))

	t.Run("store", func(t *testing.T) {
		svc.Record(ctx, &actionlog.Action{
			Resource: testResource1,
			Action:   "lookup",
			Severity: actionlog.Notice,
		})
	})

	t.Run("retrieve", func(t *testing.T) {
		set, _, err := svc.Find(ctx, actionlog.Filter{Resource: testResource1})
		req.NoError(err)
		req.Len(set, 1)
		req.Equal(testResource1, set[0].Resource)
		req.Equal("lookup", set[0].Action)
	})
}

func TestActionlogFilters(t *testing.T) {
	req := require.New(t)
	s := setupStore(t)
	ctx := context.Background()

	t.Run("by resource", func(t *testing.T) {
		svc := newSvc(s)
		req.NoError(store.TruncateActionlogs(ctx, s))

		svc.Record(ctx, &actionlog.Action{Resource: testResource1, Action: "create", Severity: actionlog.Notice})
		svc.Record(ctx, &actionlog.Action{Resource: testResource2, Action: "create", Severity: actionlog.Notice})

		set, _, err := svc.Find(ctx, actionlog.Filter{Resource: testResource1})
		req.NoError(err)
		req.Len(set, 1)
		req.Equal(testResource1, set[0].Resource)
	})

	t.Run("by action type", func(t *testing.T) {
		svc := newSvc(s)
		req.NoError(store.TruncateActionlogs(ctx, s))

		svc.Record(ctx, &actionlog.Action{Resource: testResource1, Action: "create", Severity: actionlog.Notice})
		svc.Record(ctx, &actionlog.Action{Resource: testResource1, Action: "delete", Severity: actionlog.Notice})

		set, _, err := svc.Find(ctx, actionlog.Filter{Action: "create"})
		req.NoError(err)
		req.Len(set, 1)
		req.Equal("create", set[0].Action)
	})

	t.Run("by origin", func(t *testing.T) {
		svc := newSvc(s)
		req.NoError(store.TruncateActionlogs(ctx, s))

		svc.Record(ctx, &actionlog.Action{Resource: testResource1, Action: "create", RequestOrigin: actionlog.RequestOrigin_API_REST})
		svc.Record(ctx, &actionlog.Action{Resource: testResource1, Action: "create", RequestOrigin: actionlog.RequestOrigin_Automation})

		set, _, err := svc.Find(ctx, actionlog.Filter{Origin: actionlog.RequestOrigin_API_REST})
		req.NoError(err)
		req.Len(set, 1)
		req.Equal(actionlog.RequestOrigin_API_REST, set[0].RequestOrigin)
	})

	t.Run("limit", func(t *testing.T) {
		svc := newSvc(s)
		req.NoError(store.TruncateActionlogs(ctx, s))

		for range 5 {
			svc.Record(ctx, &actionlog.Action{Resource: testResource1, Action: "lookup", Severity: actionlog.Notice})
		}

		set, _, err := svc.Find(ctx, actionlog.Filter{Resource: testResource1, Limit: 3})
		req.NoError(err)
		req.Len(set, 3)
	})
}

func TestActionlogDelta(t *testing.T) {
	req := require.New(t)
	s := setupStore(t)
	ctx := context.Background()
	svc := newSvc(s)
	req.NoError(store.TruncateActionlogs(ctx, s))

	old := &testResource{Handle: "proj-a", Name: "Old Name", Active: false}
	updated := &testResource{Handle: "proj-a", Name: "New Name", Active: true}
	delta := actionlog.DiffResourceState(updated, old)

	t.Run("store", func(t *testing.T) {
		req.NotEmpty(delta)
		svc.Record(ctx, &actionlog.Action{
			Resource: testResource1,
			Action:   "update",
			Severity: actionlog.Notice,
			Delta:    actionlog.Delta(delta),
		})
	})

	t.Run("retrieve", func(t *testing.T) {
		set, _, err := svc.Find(ctx, actionlog.Filter{Resource: testResource1})
		req.NoError(err)
		req.Len(set, 1)
		req.NotNil(set[0].Delta)

		var found bool
		for _, ch := range set[0].Delta {
			if ch.Key == "name" {
				found = true
				req.NotEmpty(ch.Old)
				req.NotEmpty(ch.New)
			}
		}
		req.True(found, "expected 'name' change in delta")
	})
}

func TestActionlogOldState(t *testing.T) {
	req := require.New(t)
	s := setupStore(t)
	ctx := context.Background()
	svc := newSvc(s)
	req.NoError(store.TruncateActionlogs(ctx, s))

	snapshot := &testResource{Handle: "proj-b", Name: "Before Update", Active: true}
	raw, err := json.Marshal(snapshot)
	req.NoError(err)

	t.Run("store", func(t *testing.T) {
		svc.Record(ctx, &actionlog.Action{
			Resource: testResource1,
			Action:   "update",
			Severity: actionlog.Notice,
			OldState: actionlog.OldState(raw),
		})
	})

	t.Run("retrieve", func(t *testing.T) {
		set, _, err := svc.Find(ctx, actionlog.Filter{Resource: testResource1})
		req.NoError(err)
		req.Len(set, 1)
		req.NotNil(set[0].OldState)

		var retrieved testResource
		req.NoError(json.Unmarshal([]byte(set[0].OldState), &retrieved))
		req.Equal("proj-b", retrieved.Handle)
		req.Equal("Before Update", retrieved.Name)
	})
}

func TestActionlogMeta(t *testing.T) {
	req := require.New(t)
	s := setupStore(t)
	ctx := context.Background()
	svc := newSvc(s)
	req.NoError(store.TruncateActionlogs(ctx, s))

	meta := actionlog.Meta{}
	meta.Set("projectID", uint64(42), false)
	meta.Set("label", "production", false)

	t.Run("store", func(t *testing.T) {
		svc.Record(ctx, &actionlog.Action{
			Resource: testResource1,
			Action:   "deploy",
			Severity: actionlog.Notice,
			Meta:     meta,
		})
	})

	t.Run("retrieve", func(t *testing.T) {
		set, _, err := svc.Find(ctx, actionlog.Filter{Resource: testResource1})
		req.NoError(err)
		req.Len(set, 1)
		req.Equal("42", set[0].Meta["projectID"])
		req.Equal("production", set[0].Meta["label"])
	})
}

func TestActionlogPolicy(t *testing.T) {
	req := require.New(t)
	s := setupStore(t)
	ctx := context.Background()

	t.Run("disabled skips storage", func(t *testing.T) {
		svc := actionlog.NewService(s, zap.NewNop(), zap.NewNop(), actionlog.MakeDisabledPolicy())
		req.NoError(store.TruncateActionlogs(ctx, s))

		svc.Record(ctx, &actionlog.Action{
			Resource: testResource1,
			Action:   "lookup",
			Severity: actionlog.Notice,
		})

		set, _, err := svc.Find(ctx, actionlog.Filter{Resource: testResource1})
		req.NoError(err)
		req.Empty(set)
	})
}

func TestActionlogHistory(t *testing.T) {
	req := require.New(t)
	s := setupStore(t)
	ctx := context.Background()
	svc := actionlog.NewService(s, zap.NewNop(), zap.NewNop(), actionlog.MakeDebugPolicy())
	req.NoError(store.TruncateActionlogs(ctx, s))

	base := time.Now().Add(-3 * time.Second).Truncate(time.Millisecond)

	t.Run("store", func(t *testing.T) {
		for i := range 3 {
			req.NoError(store.CreateActionlog(ctx, s, &actionlog.Action{
				ID:          uint64(i + 1),
				Timestamp:   base.Add(time.Duration(i) * time.Second),
				Resource:    testResource1,
				Action:      "step",
				Severity:    actionlog.Notice,
				Description: fmt.Sprintf("step-%d", i),
			}))
		}
	})

	t.Run("retrieve oldest first", func(t *testing.T) {
		history, err := svc.History(ctx, testResource1)
		req.NoError(err)
		req.Len(history, 3)
		for i, h := range history {
			req.Equal(fmt.Sprintf("step-%d", i), h.Description, "expected oldest first at index %d", i)
		}
	})
}

func TestDiffResourceState(t *testing.T) {
	req := require.New(t)

	old := &testResource{Handle: "proj", Name: "Old", Active: false}
	updated := &testResource{Handle: "proj", Name: "New", Active: true}

	changes := actionlog.DiffResourceState(updated, old)
	req.NotEmpty(changes)

	byKey := make(map[string]*struct{ old, new any })
	for _, ch := range changes {
		byKey[ch.Key] = &struct{ old, new any }{}
		if len(ch.Old) > 0 {
			byKey[ch.Key].old = ch.Old[0]
		}
		if len(ch.New) > 0 {
			byKey[ch.Key].new = ch.New[0]
		}
	}

	req.Contains(byKey, "name")
	req.Contains(byKey, "active")
	req.NotContains(byKey, "handle", "unchanged field must not appear in diff")
}
