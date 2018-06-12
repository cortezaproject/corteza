package rbac

// Permissions is a stateful object (Session, User, Roles, Resource)
type Permissions interface {
	Grant(permission string) error
	Revoke(permission string) error

	// User may modify own permissions (ie, "enter moderator view", "exit moderator view" or similar scenario);
	GrantAuthorize(permission string) error
	RevokeAuthorize(permission string) error

	// List active permissions
	List() ([]string, error)

	// List authorized permissions
	ListAuthorized() ([]string, error)

	// Check permission of stateful object (Session, User, Roles)
	CheckAccess(permission string) (bool, error)
}

// Roles is a stateful object (Session, User)
type Roles interface {
	Create(role string) error
	Delete(role string) error
	List() ([]string, error)
	Permissions(role string) Permissions
}

// Session object holds session state (Create, Load)
type Session interface {
	Create(userID string, roles ...string) error
	Load(sessionID string) error

	Delete() error
	User() (string, error)

	// Roles and Permissions return session scoped objects
	Roles() Roles
	Permissions() Permissions
}

// Resource
type Resource interface {
	Create(resource string) error
	Delete(resource string) error

	RolePermissions(resource string, role string) Permissions
	UserPermissions(resource string, user string) Permissions
}

// User is a static object
type User interface {
	Create(user string) error
	Delete(user string) error

	Roles(user string) Roles
	Permissions(user string) Roles
}
