package locale

var (
	// Global locales service
	global *Languages
)

// Global returns global RBAC service
func Global() *Languages {
	return global
}

// SetGlobal re-sets global service
func SetGlobal(ll *Languages) {
	global = ll
}
