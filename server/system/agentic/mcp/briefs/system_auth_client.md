# Brief: `system_auth_client`

Tool file: `server/system/agentic/auth_client_tools.go`
Handler:   `server/system/agentic/auth_client_handler.go`
Group: `configuring`

## 1. Service

`system/service/auth_client.go` — concrete `*authClient`,
`sysService.DefaultAuthClient` (not an interface, unlike `DefaultUser`).

| Method | Signature | Becomes |
|---|---|---|
| `LookupByID` | `(ctx, ID uint64) (*types.AuthClient, error)` | single-fetch — **blanks the secret** |
| `Search` | `(ctx, types.AuthClientFilter)` | list mode |
| `Update` | `(ctx, *types.AuthClient) (*types.AuthClient, error)` | `update` |
| `DeleteByID` | `(ctx, ID uint64) error` | `delete` |
| `ExposeSecret` | `(ctx, ID uint64) (string, error)` | **do not write** — see traps |
| `RegenerateSecret` | via `onRegenerateSecret` | **do not write** — see traps |
| `IsDefaultClient` | predicate | not a tool; use it in the projection |

**Corrected after implementation.** `Create` *does* exist in the standard
shape — `auth_client.gen.go:27`, authorized by `CanCreateAuthClient`. It lives
in the generated file, which is the same trap this document warns about for
`UndeleteByID`. It is still not written: `beforeCreate` (`auth_client.go:65`)
generates `new.Secret` and `Create` returns that same object with the field
populated, so any `JSONResult` of it hands over a live credential. Filed under
§8.6b.

`UndeleteByID` also exists (`auth_client.gen.go:67`, guarded by
`CanDeleteAuthClient`), so §4.2 applies and `undelete` is written. An earlier
draft of this brief claimed neither existed; both claims were wrong.

Risks: `lookup` read; `update` write; `delete` destructive.

## 2. Authorization (§8.6) — PASSES, but the seed deny-list was imprecise

5 `ac.Can*` calls. `lookupByID` (`auth_client.go:137`) calls
`CanReadAuthClient`, and `ExposeSecret` goes through it — so `ExposeSecret` is
**authorized**, contrary to the seed deny-list in CONVENTIONS §8.6, which
listed it as unauthorized.

The real objection is different and worse: **read permission yields the
secret**. `LookupByID` deliberately blanks `client.Secret` before returning;
`ExposeSecret` returns it. Anyone who can read an auth client can therefore
obtain its credential, and §4's note that `read` is about state change and not
disclosure is exactly this case.

Correct the seed list when you touch it: ExposeSecret is excluded because it
discloses a credential, not because it skips a check.

## 3. Identifier strategy

**No `FindByAny`, no `FindByHandle`** — only `FindByID`/`LookupByID`. But
`AuthClientFilter` has a `Handle` field, so a handle can be resolved through
`Search`.

Write a local resolver: numeric → `LookupByID`; otherwise `Search` with
`Handle` set and take the single result, erroring clearly on none or many.
Follow `taq_handler.go`'s `resolve` for the shape.

## 4. Filter fields → params

`AuthClientFilter` (`system/types/auth_client.go:13`):

| Field | Param? | Note |
|---|---|---|
| `Handle` | yes | |
| `Deleted` | yes, `includeDeleted` boolean | |
| `AuthClientID []string` | **no** | the ref param covers single fetch |
| `LabeledIDs`, `Labels`, `Check` | **no** | internal |

There is no `Query` field — do not declare one. Plus `limit` / `pageCursor`.

## 5. List projection

`{authClientID, handle, enabled, validFrom, expiresAt, deletedAt}` plus
`isDefault` from `IsDefaultClient`. **Never include `secret`** in any
projection, even though the raw type carries the field — `LookupByID` blanks
it, so a single-fetch that returns the raw type is already safe, but a
hand-built projection could reintroduce it. Do not build one from the store.

## 6. Traps

**Do not write a tool that returns a secret.** `ExposeSecret` and any
regenerate operation hand a working credential to whatever is reading the
tool's output — a model's context, a transcript, a log. There is no legitimate
agent workflow that needs it, and the disclosure is irreversible. File it as a
deliberate gap in the coverage matrix with this reasoning.

**`Update` returns the secret unblanked.** Found during implementation, and the
most dangerous thing in this resource. `Update` (`auth_client.go:239`) loads the
row with `loadAuthClient` — straight from the store, nothing blanked — copies
the changed fields on and returns it. Only `LookupByID` and `Search` blank
`Secret`. Marshalling `Update`'s return value verbatim, which is exactly the
shape the reminder exemplar invites, discloses a working credential on every
update. Clear `res.Secret` before `JSONResult`.

Safe in the other direction: `Update` copies a fixed field list off the object
it is given and `Secret` is not among them, so handing back a blanked client
does not wipe the stored credential.

**Security fields are privilege assignment, not configuration.**
`impersonateUser`, `userGroup`, `permittedRoles`, `prohibitedRoles` and
`forcedRoles` on update are how a client becomes a log-in-as-anyone key. Keep
them out of the tool, and say in the description that they are configured by a
person.

**`DefaultAuthClient` is a concrete pointer**, not an interface, so unlike
`ModuleService` there is no interface to extend if a method is missing. If the
method you need is unexported, that is a service change, not a tool change —
raise it.
