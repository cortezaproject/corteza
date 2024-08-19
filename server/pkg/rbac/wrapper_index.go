package rbac

import "sync"

type (
	wrapperIndex struct {
		mux   sync.RWMutex
		rules *ruleIndex
	}
)

func (svc *wrapperIndex) get(role uint64, op string, res string) (out []*Rule) {
	svc.mux.RLock()
	defer svc.mux.RUnlock()

	return svc.rules.get(role, op, res)
}

func (svc *wrapperIndex) hasRole(role uint64) (ok bool) {
	svc.mux.RLock()
	defer svc.mux.RUnlock()

	_, ok = svc.rules.children[role]
	return
}

// @todo since it's like so, we might not need the trie to have deletable elements
func (svc *wrapperIndex) remove(roles ...uint64) {
	svc.mux.Lock()
	defer svc.mux.Unlock()
	for _, r := range roles {
		delete(svc.rules.children, r)
	}
}

func (svc *wrapperIndex) add(rules ...*Rule) {
	svc.mux.Lock()
	defer svc.mux.Unlock()
	svc.rules.add(rules...)
}
