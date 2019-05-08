package permissions

import (
	"context"
	"sync"
)

type (
	service struct {
		l sync.Locker

		rules      RuleSet
		repository *repository
	}
)

// Service initializes service{} struct
//
// service{} struct preloads, checks, grants and flushes privileges to and from repository
// It acts as a caching layer
func Service(repository *repository) *service {
	return &service{
		repository: repository,
	}
}

func (svc *service) Preload(ctx context.Context) (err error) {
	svc.l.Lock()
	defer svc.l.Unlock()

	svc.rules, err = svc.repository.With(ctx).Load()
	if err != nil {
		return
	}

	return nil
}

// Check verifies if role has access to perform an operation on a resource
//
// See RuleSet's Check() func for details
func (svc service) Check(res Resource, op Operation, roles ...uint64) (v Access) {
	svc.l.Lock()
	defer svc.l.Unlock()

	return svc.rules.Check(res, op, roles...)
}

// Grant appends and/or overwrites internal rules slice
//
// All rules with Inherit are removed
func (svc service) Grant(ctx context.Context, rules ...*Rule) error {
	svc.l.Lock()
	defer svc.l.Unlock()

	// @todo update svc.rules

	return nil
}

func (svc service) watcher() {
	// @todo will listen to chan and load new stuff every time it gets a ping
}
