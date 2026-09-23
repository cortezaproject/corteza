# ProjectGroup + ProjectGroupMember — Implementation Spec

## Context

Projects need a persistent "Technical Architecture" grouping feature. Named groups bundle resources (modules, pages, automations, agents, etc.) that logically belong together within a project. Groups are server-side entities — not metadata blobs — so they can be queried, listed, and independently managed.

The member of a group is a **resource ID** (uint64 pointing to any project resource by ID), not a user. The resource's `kind` is already on the resource itself so no kind column is needed on the join table.

Pattern follows `project_member` exactly throughout. Reference files:
- `server/system/project_member.cue` — CUE pattern for project-scoped member join
- `server/system/project.cue` — CUE pattern for project-scoped parent
- `server/system/service/project.go` — service pattern with member ops
- `server/system/rest/project.go` — REST controller pattern

---

## Step 1 — CUE definitions (triggers codegen)

### `server/system/project_group.cue`

```cue
package system

import (
	"github.com/crusttech/human/server/codegen/schema"
)

project_group: {
	features: {
		labels:        false
		projectScoped: true
	}

	model: {
		attributes: {
			id:         schema.IdField
			tenant_id:  schema.TenantRefField
			project_id: schema.ProjectRefField
			handle:     schema.HandleField
			name: {
				sortable: true
				dal: {}
			}
			description: {
				goType: "string"
				dal: {length: 512}
			}
			created_at: schema.SortableTimestampNowField
			updated_at: schema.SortableTimestampNilField
			deleted_at: schema.SortableTimestampNilField
		}

		indexes: {
			"primary": {attribute: "id"}
			"unique_handle_per_project": {
				fields: [{attribute: "project_id"}, {attribute: "handle", modifiers: ["LOWERCASE"]}]
				predicate: "handle != '' AND deleted_at IS NULL"
			}
		}
	}

	filter: {
		struct: {
			project_group_id: {goType: "[]uint64", ident: "projectGroupID", storeIdent: "id"}
			tenant_id:        schema.TenantFilterField
			project_id:       schema.ProjectFilterField
			handle:           {goType: "string"}
			name:             {goType: "string"}
			deleted:          {goType: "filter.State", storeIdent: "deleted_at"}
		}
		query: ["handle", "name"]
		byValue: ["project_group_id", "project_id", "handle"]
		byNilState: ["deleted"]
	}

	rbac: {
		operations: {
			read: description:             "Read project group"
			update: description:           "Update project group"
			delete: description:           "Delete project group"
			"members.manage": description: "Manage project group members"
		}
	}

	envoy: { omit: true }

	store: {
		api: {
			lookups: [
				{
					fields: ["id"]
					description: "searches for project group by ID"
				},
				{
					fields: ["project_id", "handle"]
					nullConstraint: ["deleted_at"]
					constraintCheck: true
					description: "searches for project group by project and handle; returns only non-deleted"
				},
			]
		}
	}
}
```

### `server/system/project_group_member.cue`

```cue
package system

project_group_member: {
	features: {
		labels:  false
		paging:  false
		sorting: false
		checkFn: false
	}

	model: {
		attributes: {
			project_group_id: {
				ident:      "projectGroupID"
				goType:     "uint64"
				storeIdent: "rel_project_group"
				dal: {type: "Ref", refModelResType: "corteza::system:project_group"}
			}
			resource_id: {
				ident:      "resourceID"
				goType:     "uint64"
				storeIdent: "rel_resource"
				dal: {type: "ID"}
			}
			created_at: schema.SortableTimestampNowField
		}

		indexes: {
			"primary": {
				fields: [{attribute: "project_group_id"}, {attribute: "resource_id"}]
			}
		}
	}

	filter: {
		struct: {
			project_group_id: {goType: "uint64", ident: "projectGroupID", storeIdent: "rel_project_group"}
			resource_id:      {goType: "uint64", ident: "resourceID",     storeIdent: "rel_resource"}
		}
		byValue: ["project_group_id", "resource_id"]
	}

	envoy: { omit: true }

	store: {
		api: {
			lookups: [
				{
					fields: ["project_group_id", "resource_id"]
					description: "searches for group member by group and resource"
				},
			]
		}
	}
}
```

---

## Step 2 — REST yaml additions

Add to `server/system/rest.yaml`. Follow the existing project members block as reference.

```yaml
  - title: Project Groups
    path: '/projects/{projectID}/groups'
    entrypoint: projectGroup
    authentication:
      - Client ID
      - Session ID
    imports:
      - github.com/crusttech/human/server/system/types
      - time
    apis:
      - name: list
        method: GET
        title: List project groups
        path: '/'
        parameters:
          path:
            - type: uint64
              name: projectID
              required: true
              title: Project ID
          get:
            - { type: string,   name: query,          title: Search query }
            - { type: string,   name: handle,          title: Handle filter }
            - { type: '[]string', name: projectGroupID, title: Filter by group IDs }
            - { type: uint,     name: deleted,         title: "Deleted filter (0=exclude 1=only 2=include)" }
            - { type: uint,     name: limit,           title: Limit }
            - { type: bool,     name: incTotal,        title: Include total count }
            - { type: string,   name: pageCursor,      title: Page cursor }
            - { type: string,   name: sort,            title: Sort }

      - name: create
        method: POST
        title: Create project group
        path: '/'
        parameters:
          path:
            - { type: uint64, name: projectID, required: true, title: Project ID }
          post:
            - { type: string, name: handle,      required: true, title: Handle }
            - { type: string, name: name,         title: Name }
            - { type: string, name: description,  title: Description }

      - name: read
        method: GET
        title: Read project group
        path: '/{projectGroupID}'
        parameters:
          path:
            - { type: uint64, name: projectID,      required: true, title: Project ID }
            - { type: uint64, name: projectGroupID,  required: true, title: Project Group ID }

      - name: update
        method: PUT
        title: Update project group
        path: '/{projectGroupID}'
        parameters:
          path:
            - { type: uint64, name: projectID,      required: true, title: Project ID }
            - { type: uint64, name: projectGroupID,  required: true, title: Project Group ID }
          post:
            - { type: string,     name: handle,      title: Handle }
            - { type: string,     name: name,         title: Name }
            - { type: string,     name: description,  title: Description }
            - { type: '*time.Time', name: updatedAt,  title: Last update timestamp }

      - name: delete
        method: DELETE
        title: Delete project group
        path: '/{projectGroupID}'
        parameters:
          path:
            - { type: uint64, name: projectID,      required: true, title: Project ID }
            - { type: uint64, name: projectGroupID,  required: true, title: Project Group ID }

      - name: memberList
        method: GET
        title: List group members
        path: '/{projectGroupID}/members'
        parameters:
          path:
            - { type: uint64, name: projectID,      required: true, title: Project ID }
            - { type: uint64, name: projectGroupID,  required: true, title: Project Group ID }

      - name: memberAdd
        method: POST
        title: Add resource to group
        path: '/{projectGroupID}/members/{resourceID}'
        parameters:
          path:
            - { type: uint64, name: projectID,      required: true, title: Project ID }
            - { type: uint64, name: projectGroupID,  required: true, title: Project Group ID }
            - { type: uint64, name: resourceID,      required: true, title: Resource ID }

      - name: memberRemove
        method: DELETE
        title: Remove resource from group
        path: '/{projectGroupID}/members/{resourceID}'
        parameters:
          path:
            - { type: uint64, name: projectID,      required: true, title: Project ID }
            - { type: uint64, name: projectGroupID,  required: true, title: Project Group ID }
            - { type: uint64, name: resourceID,      required: true, title: Resource ID }
```

---

## Step 3 — Run codegen

```bash
make codegen
```

Generates (do not edit these manually):

| File | What |
|---|---|
| `server/system/types/project_group.go` | `ProjectGroup`, `ProjectGroupFilter`, `ProjectGroupSet` |
| `server/system/types/project_group_member.go` | `ProjectGroupMember`, `ProjectGroupMemberFilter`, `ProjectGroupMemberSet` |
| `server/system/types/type_set.gen.go` | Updated with new Set types |
| `server/system/types/getters_setters.gen.go` | Updated |
| `server/system/types/rbac.gen.go` | Updated with ProjectGroup RBAC operations |
| `server/store/interfaces.gen.go` | `ProjectGroups` + `ProjectGroupMembers` store interfaces |
| `server/store/adapters/rdbms/rdbms.gen.go` | SQL implementations |
| `server/system/rest/request/project_group.go` | Request structs |
| `server/system/rest/handlers/project_group.go` | HTTP handler wiring |

---

## Step 4 — Service layer

### `server/system/service/project_group.go` (new, manual)

Model on `service/project.go` exactly.

**Interface:**
```go
type ProjectGroupService interface {
    FindByID(ctx context.Context, id uint64) (*types.ProjectGroup, error)
    Search(ctx context.Context, f types.ProjectGroupFilter) (types.ProjectGroupSet, types.ProjectGroupFilter, error)
    Create(ctx context.Context, g *types.ProjectGroup) (*types.ProjectGroup, error)
    Update(ctx context.Context, g *types.ProjectGroup) (*types.ProjectGroup, error)
    DeleteByID(ctx context.Context, id uint64) error

    MemberList(ctx context.Context, projectGroupID uint64) (types.ProjectGroupMemberSet, error)
    MemberAdd(ctx context.Context, projectGroupID, resourceID uint64) error
    MemberRemove(ctx context.Context, projectGroupID, resourceID uint64) error
}
```

**Create:** validate handle uniqueness within project via `LookupProjectGroupByProjectIDHandle`; set `TenantID` and `ProjectID` from context scope.

**Update:** re-validate handle uniqueness if handle changed (exclude self from check).

**DeleteByID:** soft delete — set `DeletedAt`.

**MemberAdd:**
1. Load group → verify belongs to same project as request scope
2. `LookupProjectGroupMemberByProjectGroupIDResourceID` → if found, no-op or return error
3. `CreateProjectGroupMember`

**MemberRemove:** `DeleteProjectGroupMemberByProjectGroupIDResourceID`

**MemberList:** `SearchProjectGroupMembers` filtered by `projectGroupID`

**Default singleton:**
```go
var DefaultProjectGroup ProjectGroupService
```

---

## Step 5 — REST controller

### `server/system/rest/project_group.go` (new, manual)

Model exactly on `rest/project.go`.

```go
type (
    ProjectGroup struct {
        svc projectGroupService
        ac  projectGroupAccessController
    }

    projectGroupPayload struct {
        *types.ProjectGroup
        CanUpdate        bool `json:"canUpdate"`
        CanDelete        bool `json:"canDelete"`
        CanManageMembers bool `json:"canManageMembers"`
    }

    projectGroupSetPayload struct {
        Filter types.ProjectGroupFilter `json:"filter"`
        Set    []*projectGroupPayload   `json:"set"`
    }

    projectGroupMemberPayload struct {
        ResourceID uint64 `json:"resourceID,string"`
    }

    projectGroupService interface {
        FindByID(ctx context.Context, id uint64) (*types.ProjectGroup, error)
        Search(ctx context.Context, f types.ProjectGroupFilter) (types.ProjectGroupSet, types.ProjectGroupFilter, error)
        Create(ctx context.Context, g *types.ProjectGroup) (*types.ProjectGroup, error)
        Update(ctx context.Context, g *types.ProjectGroup) (*types.ProjectGroup, error)
        DeleteByID(ctx context.Context, id uint64) error
        MemberList(ctx context.Context, projectGroupID uint64) (types.ProjectGroupMemberSet, error)
        MemberAdd(ctx context.Context, projectGroupID, resourceID uint64) error
        MemberRemove(ctx context.Context, projectGroupID, resourceID uint64) error
    }

    projectGroupAccessController interface {
        CanGrant(context.Context) bool
        CanUpdateProjectGroup(context.Context, *types.ProjectGroup) bool
        CanDeleteProjectGroup(context.Context, *types.ProjectGroup) bool
        CanManageMembersOnProjectGroup(context.Context, *types.ProjectGroup) bool
    }
)
```

`New()` wires `DefaultProjectGroup` and `DefaultAccessControl`.

Handlers: List, Create, Read, Update, Delete, MemberList, MemberAdd, MemberRemove — all follow same shape as `rest/project.go` methods.

---

## Step 6 — Wire up

### `server/system/component.go`

Instantiate alongside other default services:
```go
DefaultProjectGroup = service.NewProjectGroup(store, DefaultAccessControl)
```

### Router mount

Wherever project routes are mounted, add:
```go
r.Mount("/", rest.ProjectGroup{}.New().MountRoutes)
```

---

## Step 7 — DB migration

New migration file in the existing migrations directory.

**`project_groups` table:**
```sql
CREATE TABLE project_groups (
    id          BIGINT UNSIGNED NOT NULL,
    tenant_id   BIGINT UNSIGNED NOT NULL DEFAULT 0,
    project_id  BIGINT UNSIGNED NOT NULL,
    handle      VARCHAR(255)    NOT NULL DEFAULT '',
    name        VARCHAR(255)    NOT NULL DEFAULT '',
    description VARCHAR(512)    NOT NULL DEFAULT '',
    created_at  DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at  DATETIME,
    deleted_at  DATETIME,
    PRIMARY KEY (id)
);
-- Uniqueness of handle within project enforced at application layer (deleted_at partial)
```

**`project_group_members` table:**
```sql
CREATE TABLE project_group_members (
    rel_project_group  BIGINT UNSIGNED NOT NULL,
    rel_resource       BIGINT UNSIGNED NOT NULL,
    created_at         DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (rel_project_group, rel_resource)
);
```

---

## Verification checklist

- `POST /api/system/projects/:pid/groups` with `handle` + `name` → 200, group returned
- `GET /api/system/projects/:pid/groups` → list includes created group
- `GET /api/system/projects/:pid/groups/:gid` → single group
- `PUT /api/system/projects/:pid/groups/:gid` → updates name/description
- Duplicate handle within same project → error
- `DELETE /api/system/projects/:pid/groups/:gid` → soft-deleted, excluded from default list
- `POST /api/system/projects/:pid/groups/:gid/members/:rid` → member added
- `GET /api/system/projects/:pid/groups/:gid/members` → lists resource IDs
- `DELETE /api/system/projects/:pid/groups/:gid/members/:rid` → removed
- Adding same resource twice → no-op or clear error
