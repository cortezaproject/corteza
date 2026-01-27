package registry

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/cortezaproject/corteza/server/pkg/automation_exec/types"
)

type (
	ExecutableEntry struct {
		Executable   types.Executable
		Status       types.ExecutableStatus
		RegisteredAt time.Time
		DeprecatedAt *time.Time
	}

	UsageChecker interface {
		IsInUse(ctx context.Context, executableID types.ExecutableID, revision int) (bool, error)
	}

	registry struct {
		// map[ExecutableID]map[Revision]*ExecutableEntry
		entries map[types.ExecutableID]map[int]*ExecutableEntry

		usageChecker UsageChecker
		stats        stats
		mux          sync.RWMutex
	}
)

// The Registry is the catalog of executable definitions
//
// * Stores versioned executables
// * Tracks active vs deprecated revisions
// * Answers what can be executed
// * Prevents removal of in-use definitions
func Registry(usageChecker UsageChecker) *registry {
	return &registry{
		entries:      make(map[types.ExecutableID]map[int]*ExecutableEntry),
		usageChecker: usageChecker,
	}
}

// Add registers an executable revision.
// exec.Revision is treated as authoritative.
func (r *registry) Add(ctx context.Context, exec types.Executable) error {
	r.mux.Lock()
	defer r.mux.Unlock()

	revisions := r.entries[exec.ID]
	if revisions == nil {
		revisions = make(map[int]*ExecutableEntry)
		r.entries[exec.ID] = revisions
		r.statsOnNewExecutable()
	}

	if _, exists := revisions[exec.Revision]; exists {
		return ErrRevisionAlreadyExists
	}

	entry := &ExecutableEntry{
		Executable:   exec,
		Status:       types.StatusActive,
		RegisteredAt: time.Now(),
	}

	revisions[exec.Revision] = entry
	r.statsOnRegisterRevision()

	return nil
}

func (r *registry) Deprecate(ctx context.Context, executableID types.ExecutableID, revision int) error {
	r.mux.Lock()
	defer r.mux.Unlock()

	entry, err := r.getEntryLocked(executableID, revision)
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

func (r *registry) Remove(ctx context.Context, executableID types.ExecutableID, revision int) error {
	r.mux.Lock()
	defer r.mux.Unlock()

	entry, err := r.getEntryLocked(executableID, revision)
	if err != nil {
		return err
	}

	if entry.Status != types.StatusDeprecated {
		return ErrExecutableNotDeprecated
	}

	if r.usageChecker == nil {
		return ErrUsageCheckerRequired
	}

	inUse, err := r.usageChecker.IsInUse(ctx, executableID, revision)
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

func (r *registry) Get(ctx context.Context, executableID types.ExecutableID, revision int) (ExecutableEntry, error) {
	r.mux.RLock()
	defer r.mux.RUnlock()

	entry, err := r.getEntryLocked(executableID, revision)
	if err != nil {
		return ExecutableEntry{}, err
	}

	return *entry, nil
}

func (r *registry) GetLatest(ctx context.Context, executableID types.ExecutableID) (ExecutableEntry, error) {
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

func (r *registry) getEntryLocked(executableID types.ExecutableID, revision int) (*ExecutableEntry, error) {
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
