package rbac

type (
	statLogger interface {
		CacheHit(roles []uint64, resource string, op string)
		CacheMiss(roles []uint64, resource string, op string)
		CacheUpdate(*Rule)
	}
)
