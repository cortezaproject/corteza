# Multi-Tenancy — Middleware & Access Token Changes

Companion to `multi-tenancy-design.md`. Covers only token changes and middleware implementation.

---

## Access Token Changes

### Background — current token

Token is a signed JWT produced by `makeToken` in `server/pkg/auth/token_issuer.go:157`.

Current claims: `jti`, `sub` (userID), `exp`, `aud`, `iss`, `iat`, `clientID`, `scope`, `roles`.

Built from `TokenRequest` struct (`server/pkg/auth/token_issuer.go:40`):
```
AccessToken, RefreshToken, Expiration, Audience, Issuer, IssuedAt, ClientID, UserID, Roles, Scope
```

Identity is decoded back from token via `IdentityFromToken` (`token_issuer.go:224`) — reads `sub` and `roles` claims only.

### Required changes

**1. Add `TenantID` to `TokenRequest`**

Add `TenantID uint64` field to `TokenRequest`. No other fields needed.

**2. Stamp `tenantID` claim in `makeToken`**

In `makeToken`, after existing claim sets, add:
```go
token.Set("tenantID", toString(req.TenantID))
```

Only set when `TenantID != 0` (system-level tokens have no tenant).

**3. Resolve tenant at token issuance time**

At the point where a token is issued after successful login, before calling the issuer:

1. Load `TenantMembership` by `UserID`
2. If found and active → set `req.TenantID`
3. If not found → token issued without `tenantID` (system user or unassigned)
4. If found but suspended → reject login with appropriate error

This happens in the OAuth2 token issuance path — locate where `IssueOptFn` options are assembled for the token request and add a tenant resolution step there.

**4. Add `TenantFromToken` helper**

Add `TenantFromToken(token jwt.Token) uint64` — reads the `tenantID` claim, returns 0 if absent.

Do not merge into `IdentityFromToken` — identity and scope are separate concerns.

---

## Required Middlewares

Three middlewares, applied in order on every authenticated request.

### Middleware 1 — Existing: JWT validation

Already exists at `server/pkg/auth/token_middleware.go:26`. Validates JWT, extracts identity, sets it in context. **No changes needed.**

### Middleware 2 — New: Tenant scope resolver

**Location:** `server/pkg/scope/middleware.go` (new file)

**Runs after:** JWT middleware (identity must be in context)

**Responsibility:**

1. Call `TenantFromToken(token)` to extract `tenantID` from claims
2. If `tenantID == 0` and route requires tenant scope → 403
3. Validate tenant is active — use short-lived in-process cache (TTL ~30s) keyed by `tenantID`; do not hit DB on every request
4. Call `scope.SetScopeToContext(ctx, scope.Scope{TenantID: tenantID})`
5. Initialize tenant runtime via `ScopeRegistry.Tenant(tenantID)` — lazy-inits DAL connection for the tenant

**Subdomain validation (if subdomain strategy is adopted):**

Extract tenant handle from `Host` header → resolve to `tenantID` → compare with token's `tenantID`. Mismatch → 401. Prevents a valid token from one tenant being replayed against another tenant's subdomain.

**Applied to:** all routes under `/api/` except system-admin routes.

### Middleware 3 — New: Project scope resolver

**Location:** `server/pkg/scope/middleware.go` (same file, separate function)

**Runs after:** Tenant scope middleware

**Responsibility:**

1. Extract project handle or ID from URL path (e.g. `/api/project/:projectHandle/...`)
2. Resolve project handle → `projectID` (tenant-scoped lookup; project must belong to current tenant)
3. Check `ProjectMember` record exists for `userID + projectID`
   - If no direct membership, check if user has tenant-level membership and project visibility is `open`
   - If neither → 403
4. Load `ProjectCapabilities` from member's `RolePreset` — attach to context (define `projectCtxKey` analogous to `scopeCtxKey`)
5. Set `ProjectID` on the scope: read current scope from context, set `ProjectID`, write back
6. Initialize project runtime via `ScopeRegistry.Project(tenantID, projectID)`

**Applied to:** only routes that carry a project path segment. Tenant-level routes skip this middleware entirely.

### Middleware ordering on router

```
JWT validation  →  Tenant scope resolver  →  [Project scope resolver]  →  handler
```

Project scope resolver is conditional — mounted only on project-scoped route groups, not globally.

### Capability enforcement

Capability checks (`canWrite`, `canGrantApproval`, etc.) are **not** middleware. They are enforced in the service layer or REST handler per operation. Middleware only resolves and attaches capabilities to context. Handlers read from context and reject if the required capability is absent.

---

## Summary of file changes

| File | Change |
|---|---|
| `server/pkg/auth/token_issuer.go` | Add `TenantID` to `TokenRequest`; stamp `tenantID` claim in `makeToken`; add `TenantFromToken` helper |
| OAuth2 token issuance path (TBD) | Resolve `TenantMembership` before token issuance; set `req.TenantID` |
| `server/pkg/scope/middleware.go` | New file — `TenantScopeMiddleware` and `ProjectScopeMiddleware` |
| Router setup (TBD) | Mount tenant middleware globally on `/api/`; mount project middleware on project route groups |
