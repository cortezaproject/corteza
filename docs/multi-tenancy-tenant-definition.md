# Multi-Tenancy — Tenant Definition

Companion to `multi-tenancy-design.md`.

Covers all types, stores, and API surface needed to define and manage tenants.

---

## Types (`server/system/types/`)

### `Tenant`

```
ID          uint64
Handle      string          // unique across system, URL-safe slug
Name        string          // display name
Status      TenantStatus    // active | suspended | archived
Config      *TenantConfig
Meta        *TenantMeta
CreatedAt   time.Time
UpdatedAt   *time.Time
SuspendedAt *time.Time
DeletedAt   *time.Time      // soft delete
```

### `TenantStatus`

Enum (string or iota):
- `active` — fully operational
- `suspended` — login rejected, data preserved
- `archived` — read-only, no login

### `TenantConfig`

```
DALConnectionID   uint64          // FK to DalConnection — tenant's data store
AuthProviders     []string        // allowed auth provider handles
FeatureFlags      map[string]bool // corteza feature toggles
Quotas            TenantQuotas
Locale            string          // default locale (BCP 47)
Timezone          string          // default timezone (IANA)
```

### `TenantQuotas`

```
MaxUsers    int   // 0 = unlimited
MaxProjects int
MaxStorage  int64 // bytes; 0 = unlimited
```

### `TenantMeta`

```
Description string
LogoID      uint64   // attachment ID
Color       string   // hex color for UI
Tags        []string
```

### `TenantMembership`

```
UserID     uint64
TenantID   uint64
Role       TenantMemberRole  // admin | member
Status     TenantMemberStatus // active | suspended | invited
InvitedBy  uint64            // userID of inviter; 0 if provisioned directly
CreatedAt  time.Time
UpdatedAt  *time.Time
```

### `TenantMemberRole`

Enum:
- `admin` — can manage tenant settings, members, projects
- `member` — can access projects per project visibility rules

### `TenantMemberStatus`

Enum:
- `active`
- `suspended` — membership disabled, user cannot log in
- `invited` — pending acceptance

### `TenantFilter`

```
Handle   string
Status   TenantStatus
Query    string          // full-text on name/handle
Page     filter.Page
Sort     filter.Sort
```

### `TenantMembershipFilter`

```
TenantID uint64
UserID   uint64
Role     TenantMemberRole
Status   TenantMemberStatus
Page     filter.Page
Sort     filter.Sort
```

---

## Store interface (`server/store/`)

### Tenant store

```
CreateTenant(ctx, *Tenant) error
UpdateTenant(ctx, *Tenant) error
DeleteTenant(ctx, tenantID uint64) error

LookupTenantByID(ctx, tenantID uint64) (*Tenant, error)
LookupTenantByHandle(ctx, handle string) (*Tenant, error)
SearchTenants(ctx, TenantFilter) ([]*Tenant, filter.Meta, error)
```

### TenantMembership store

```
CreateTenantMembership(ctx, *TenantMembership) error
UpdateTenantMembership(ctx, *TenantMembership) error
DeleteTenantMembership(ctx, userID, tenantID uint64) error

LookupTenantMembershipByUser(ctx, userID uint64) (*TenantMembership, error)
LookupTenantMembershipByUserAndTenant(ctx, userID, tenantID uint64) (*TenantMembership, error)
SearchTenantMemberships(ctx, TenantMembershipFilter) ([]*TenantMembership, filter.Meta, error)
```

`LookupTenantMembershipByUser` is the hot path — called at token issuance to resolve tenantID. Must be indexed on `userID`.

---

## Service layer (`server/system/service/`)

### `TenantService`

```
Create(ctx, *Tenant) (*Tenant, error)
Update(ctx, *Tenant) (*Tenant, error)
Delete(ctx, tenantID uint64) error
Archive(ctx, tenantID uint64) error
Suspend(ctx, tenantID uint64) error
Activate(ctx, tenantID uint64) error

Read(ctx, tenantID uint64) (*Tenant, error)
Search(ctx, TenantFilter) ([]*Tenant, filter.Meta, error)
```

### `TenantMembershipService`

```
Invite(ctx, tenantID, inviteeUserID uint64, role TenantMemberRole) (*TenantMembership, error)
Accept(ctx, tenantID, userID uint64) error
Remove(ctx, tenantID, userID uint64) error
Suspend(ctx, tenantID, userID uint64) error
Activate(ctx, tenantID, userID uint64) error

ListMembers(ctx, TenantMembershipFilter) ([]*TenantMembership, filter.Meta, error)
```

Constraint enforced in `Invite`: check no active membership exists for `userID` across any tenant (one user = one tenant).

---

## REST API surface

All endpoints require system-admin scope except where noted.

### Tenant CRUD

```
POST   /api/system/tenants                    create tenant
GET    /api/system/tenants                    list/search tenants
GET    /api/system/tenants/:tenantID          read tenant
PUT    /api/system/tenants/:tenantID          update tenant
DELETE /api/system/tenants/:tenantID          delete (soft)

POST   /api/system/tenants/:tenantID/suspend  suspend
POST   /api/system/tenants/:tenantID/activate activate
POST   /api/system/tenants/:tenantID/archive  archive
```

### Tenant membership

```
GET    /api/system/tenants/:tenantID/members               list members
POST   /api/system/tenants/:tenantID/members               invite user
PUT    /api/system/tenants/:tenantID/members/:userID        update role/status
DELETE /api/system/tenants/:tenantID/members/:userID        remove member

POST   /api/system/tenants/:tenantID/members/:userID/suspend
POST   /api/system/tenants/:tenantID/members/:userID/activate
```

Invite acceptance (if required — see open questions in design doc):
```
POST   /api/system/tenant-invitations/:token/accept
```

---

## DB schema notes

Two new tables:

**`tenants`**
- Primary key: `id`
- Unique index: `handle`
- Soft delete via `deleted_at`

**`tenant_memberships`**
- Composite primary key: `(user_id, tenant_id)`
- Index on `user_id` (hot path for token issuance)
- Index on `tenant_id` (for member listing)
- Unique constraint: one active membership per `user_id` across all tenants (enforced at application layer; DB uniqueness on `user_id` alone if one-tenant constraint is hard)

---

## Files to create / modify

| File | Action |
|---|---|
| `server/system/types/tenant.go` | New — `Tenant`, `TenantStatus`, `TenantConfig`, `TenantQuotas`, `TenantMeta`, `TenantFilter` |
| `server/system/types/tenant_membership.go` | New — `TenantMembership`, `TenantMemberRole`, `TenantMemberStatus`, `TenantMembershipFilter` |
| `server/store/tenants.go` | New — store interface for tenant + membership |
| `server/store/rdbms/tenants.go` | New — SQL implementation |
| `server/system/service/tenant.go` | New — `TenantService` |
| `server/system/service/tenant_membership.go` | New — `TenantMembershipService` |
| `server/system/rest/tenant.go` | New — REST handlers |
| `server/system/rest.yaml` | Modify — add tenant endpoints |
| `server/pkg/auth/token_issuer.go` | Modify — `TenantFromToken`; `TenantID` on `TokenRequest` (see middleware+token spec) |
| DB migration (TBD path) | New — `tenants` and `tenant_memberships` tables |
