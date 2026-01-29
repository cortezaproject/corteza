package registry

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/cortezaproject/corteza/server/pkg/automation_exec/types"
	"github.com/cortezaproject/corteza/server/pkg/id"
	"go.uber.org/zap"
)

type (
	ExecutableEntry struct {
		Executable   types.Executable
		Status       types.ExecutableStatus
		RegisteredAt time.Time
		DeprecatedAt *time.Time
	}

	usageChecker interface {
		IsExecutableInUse(ctx context.Context, executableID id.ID, revision int) (bool, error)
	}

	registry struct {
		log *zap.Logger

		// map[ExecutableID]map[Revision]*ExecutableEntry
		entries map[id.ID]map[int]*ExecutableEntry

		usageChecker usageChecker
		stats        stats
		mux          sync.RWMutex
	}
)

// Registry holds all of the executables
func Registry(log *zap.Logger, usageChecker usageChecker) *registry {
	return &registry{
		log:          log,
		entries:      make(map[id.ID]map[int]*ExecutableEntry),
		usageChecker: usageChecker,
	}
}

// Add registers an executable revision
func (r *registry) Add(ctx context.Context, exec types.Executable) error {
	r.mux.Lock()
	defer r.mux.Unlock()

	r.log.Debug("adding executable to registry", zap.String("id", exec.ID.Value()), zap.Int("revision", exec.Revision))

	revisions := r.entries[exec.ID]
	if revisions == nil {
		revisions = make(map[int]*ExecutableEntry)
		r.entries[exec.ID] = revisions
		r.statsOnNewExecutable()
	}

	_, exists := revisions[exec.Revision]

	entry := &ExecutableEntry{
		Executable:   exec,
		Status:       types.StatusActive,
		RegisteredAt: time.Now(),
	}

	revisions[exec.Revision] = entry

	if exists {
		r.log.Info("replaced existing revision", zap.Int("revision", exec.Revision))
	} else {
		r.statsOnRegisterRevision()
	}

	return nil
}

func (r *registry) Deprecate(ctx context.Context, executableID id.ID, revision int) error {
	r.mux.Lock()
	defer r.mux.Unlock()

	r.log.Debug("deprecating executable from registry", zap.String("id", executableID.Value()), zap.Int("revision", revision))

	entry, err := r.getEntry(executableID, revision)
	if err != nil {
		return err
	}

	if entry.Status == types.StatusDeprecated {
		return nil
	}

	now := time.Now()
	entry.Status = types.StatusDeprecated
	entry.DeprecatedAt = &now

	r.statsOnDeprecateRevision()
	return nil
}

func (r *registry) Remove(ctx context.Context, executableID id.ID, revision int) error {
	r.mux.Lock()
	defer r.mux.Unlock()

	r.log.Debug("removing executable from registry", zap.String("id", executableID.Value()), zap.Int("revision", revision))

	entry, err := r.getEntry(executableID, revision)
	if err != nil {
		return err
	}

	if entry.Status != types.StatusDeprecated {
		return ErrExecutableNotDeprecated
	}

	if r.usageChecker == nil {
		return ErrUsageCheckerRequired
	}

	inUse, err := r.usageChecker.IsExecutableInUse(ctx, executableID, revision)
	if err != nil {
		return fmt.Errorf("failed to check usage: %w", err)
	}
	if inUse {
		return ErrExecutableInUse
	}

	delete(r.entries[executableID], revision)
	r.statsOnRemoveDeprecatedRevision()

	if len(r.entries[executableID]) == 0 {
		delete(r.entries, executableID)
		r.statsOnRemoveExecutable()
	}

	return nil
}

func (r *registry) Get(ctx context.Context, executableID id.ID, revision int) (ExecutableEntry, error) {
	r.mux.RLock()
	defer r.mux.RUnlock()

	entry, err := r.getEntry(executableID, revision)
	if err != nil {
		return ExecutableEntry{}, err
	}

	return *entry, nil
}

func (r *registry) GetExecutable(ctx context.Context, executableID id.ID, revision int) (types.Executable, error) {
	e, err := r.Get(ctx, executableID, revision)
	if err != nil {
		return types.Executable{}, err
	}

	return e.Executable, nil
}

func (r *registry) GetLatest(ctx context.Context, executableID id.ID) (ExecutableEntry, error) {
	r.mux.RLock()
	defer r.mux.RUnlock()

	revisions, ok := r.entries[executableID]
	if !ok {
		return ExecutableEntry{}, ErrExecutableNotFound
	}

	var latest *ExecutableEntry
	latestRev := 0
	found := false

	for rev, entry := range revisions {
		if entry.Status != types.StatusActive {
			continue
		}
		if !found || rev > latestRev {
			latest = entry
			latestRev = rev
			found = true
		}
	}

	if !found {
		return ExecutableEntry{}, ErrNoActiveRevisions
	}

	return *latest, nil
}

func (r *registry) Stats(ctx context.Context) stats {
	r.mux.RLock()
	defer r.mux.RUnlock()

	return r.stats
}

func (r *registry) getEntry(executableID id.ID, revision int) (*ExecutableEntry, error) {
	revisions, ok := r.entries[executableID]
	if !ok {
		return nil, ErrExecutableNotFound
	}

	entry, ok := revisions[revision]
	if !ok {
		return nil, ErrRevisionNotFound
	}

	return entry, nil
}
