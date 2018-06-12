package rbac

// Permissions is a stateful object
type Permissions interface /* for Session, User, Roles, Resource */ {
	// Scoped for [Resource]
	Grant(permission string) error
	Revoke(permission string) error
	List() ([]string, error)

	// Check permission of stateful object (Session, User, Roles)
	CheckAccess(permission string) (bool, error)
}

// Roles is a stateful object
type Roles interface /* for Session, User */ {
	// Scoped to User
	Add(role string) error
	Delete(role string) error

	// Scoped to Session, User
	List() ([]string, error)
	ListAuthorized() ([]string, error)

	// Scoped to Session
	GrantRole(role string) error
	RevokeRole(role string) error

	// Permissions are scoped to [Session, User]
	Permissions(role string) Permissions
}

// Session object holds session state (Create, Load)
type Session interface {
	// Unscoped functions
	Create(userID string, roles ...string) error
	Load(sessionID string) error
	Delete() error

	// User returns User scoped object with global roles/permissions
	User() (User, error)

	// Roles and Permissions return session scoped objects
	Roles() Roles
	Permissions() Permissions
}

// Resource is a static object
type Resource interface {
	Load(resource string) error
	Create(resource string) error
	Delete(resource string) error

	RolePermissions(resource string, role string) Permissions
	UserPermissions(resource string, user string) Permissions
}

// Users is a static object
type User interface {
	Load(user string) error
	Create(user string) error
	Delete(user string) error

	// Roles and Permissions return User scoped objects
	Roles(user string) Roles
	Permissions(user string) Permissions
}
