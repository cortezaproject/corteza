# Related Labels — Design Study

## Context

A record's values hold bare IDs for its Record, User and File fields, and every
resource carries user IDs in `ownedBy` / `createdBy` / `updatedBy` /
`deletedBy`. A client that wants to show them has to look each kind up itself.
This study asks what it would take for the API to return **display labels for
every ID reference** in the same response — opt-in, because it costs — and for
a list to **sort by the label a person sees** instead of by the ID.

Status: research. Nothing here is built. Every number below was measured on a
throwaway worktree (Postgres 16, 2026-09-18); the method is at the end.

### Rulings taken

| #   | Ruling                                                                                                                |
| --- | --------------------------------------------------------------------------------------------------------------------- |
| 1   | Scope is everywhere: record Record/User/File values, record system users, and user references on every other resource |
| 2   | A resolved reference is a **label**, not the related object                                                           |
| 3   | Sort means sort by what is displayed (the viewer's label), not by an arbitrary path                                   |
| 4   | Labels ride as a **typed sibling of `response`** in the envelope                                                      |
| 5   | Opt-in is a **query parameter on every endpoint**                                                                     |
| 6   | Sorting is a **query-time join first**; a stored sort key only where a module is measured too big                     |
| 7   | The user label is **name → email → handle** everywhere; `username` is left out                                        |

---

## Verdict

- **Labels in the response: moderate work, cheap at runtime.** There is one
  choke point to emit them (`api.encode`), machine-readable reference metadata
  for ~60 types already exists, and the agent tools already resolve labels for
  records. With an ID-set filter, resolving a 50-row page costs 8–10 ms beside
  a 75 ms list call.
- **Sorting by label: the hard part.** The query layer has no join in its list
  path, the page cursor can only carry values that exist on the row, and
  record-level read permission cannot be expressed in SQL. It is buildable,
  and the query cost is acceptable up to about a million rows (430 ms at 1M),
  but it touches the DAL, the cursor and RBAC at once.

---

## What exists today

| Where                                                                           | What it does                                                                                                                                                                                                                   | Gap                                                                                                                                                           |
| ------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| Agent tools, `server/compose/agentic/record_labels.go:46`                       | `refs` dictionary (flat `id → label`) for Record/User values, `ownedBy`/`createdBy`/`updatedBy`, report group keys. One search per referenced module plus one user lookup, two levels deep, 500 cap, via services (RBAC holds) | Records only; flat map; batches as a `recordID = a OR …` chain; skips `deletedBy` and File values                                                             |
| Record export `resolveRefs`, `server/compose/envoy/record_datasource.go:227`    | Replaces reference values with labels in CSV/XLSX                                                                                                                                                                              | One query per row; reads users with `store.SearchUsers` (no read check, no masking, `:456`); labels placed by result position, not matched by ID (`:361-373`) |
| Action log `userPreloader`, `server/system/rest/actionlog.go:215`               | Preloads actor users                                                                                                                                                                                                           | `store.SearchUsers` directly — the pattern to avoid                                                                                                           |
| Webapp stores, `lib/vue/src/stores/useRecordStore.js:206`, `useUserStore.js:36` | Record labels batched per (namespace, module) per microtask, 100 per call; users per call                                                                                                                                      | Users have no in-flight dedupe; File values fetched one `attachmentRead` at a time; calendar and chart user dimensions never resolved                         |

The three label rules for users disagree with each other:

| Path                               | User label order                 |
| ---------------------------------- | -------------------------------- |
| Webapp `useUserResolver.ts:15`     | name → handle → email → ID       |
| Agent tools `record_labels.go:274` | Name → Username → Handle → Email |
| Export `record_datasource.go:488`  | Handle → Email → Name → ID       |

For records, when no `labelField` is set the viewer and agent tools use the
module's first field, the picker uses the first non-empty value, and export
shows nothing. The viewer formats the label with the field's own viewer
(select option text, dates); both server paths return the raw first value.

---

## Measurements

Worktree `refs-research`, a Contact module (name, Record → Company, User,
multi-value Record, File) with 100k rows (1M where stated), Company with 2,000
rows, 63 users. Medians of 15 calls unless noted.

### Lists and labels

| What                                                                                           | Result                                                                                      |
| ---------------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------- |
| Record list, 50 rows / 500 rows                                                                | 75 ms / 141 ms                                                                              |
| Same, `incTotal=true`, 100k rows / 1M rows                                                     | **4.7 s / 211 s** (0.43 s without at 1M)                                                    |
| Label resolution in process, today's OR chain, 50 rows (138 refs) / 200 rows (442 refs)        | 12–14 ms / 53–68 ms                                                                         |
| Same with an ID-set filter (`id IN (…)`), spiked                                               | **8–10 ms / 17–19 ms**                                                                      |
| OR-chain batch through REST, 885 IDs                                                           | 324 ms, of which Postgres spends 130 ms planning the next-page probe query                  |
| Same 885 IDs as `id = ANY($1)` in SQL                                                          | 0.64 ms                                                                                     |
| Webapp, RecordList with 50 rows and 6 columns (name, company, owner, tags, ownedBy, createdBy) | 1 list + 2 company label batches + **29 single-user `userList` calls for 6 distinct users** |

Every record list issues two queries: the page, and a probe for the next page.

### Sorting (SQL, `EXPLAIN ANALYZE`, LIMIT 50)

| Sort                                                          | 100k rows | 1M rows    |
| ------------------------------------------------------------- | --------- | ---------- |
| Own String field (today)                                      | 22 ms     | 127 ms     |
| Record field by ID (today — what a Record column sorts by)    | 31 ms     | 205 ms     |
| Record field by referenced label (join)                       | 52 ms     | 430 ms     |
| User field by user name (join)                                | 52 ms     | 431 ms     |
| `ownedBy` by user name (join)                                 | 19 ms     | 157 ms     |
| Page N by label (keyset on label + id)                        | —         | 407 ms     |
| Record field by rank array computed in Go (RBAC-safe variant) | —         | 719 ms     |
| Stored sort key on the row, expression index                  | —         | **0.1 ms** |
| Stored sort key, same index, query shape the DAL emits        | —         | 111 ms     |

The last two rows differ only in `NULLS FIRST`: the DAL orders `ASC NULLS
FIRST`, which a default (`NULLS LAST`) index cannot serve. An index declared
`NULLS FIRST` brings it back to 0.1 ms.

A stored sort key has to be rewritten when its target is renamed: renaming a
company rewrote 55 contacts in 24 ms; renaming the busiest owner rewrote 3,159
in 103 ms (100k rows, no index on the reference).

---

## Design

### API contract

Opt in per request:

```
GET /compose/namespace/{ns}/module/{m}/record/?refs=labels
GET /compose/namespace/{ns}/module/{m}/record/?refs=user,record
```

`labels` is already a parameter name on 53 endpoints, so the switch is `refs`.
`labels` (or no value) means every kind; a comma list limits the kinds.

The envelope gains a sibling of `response`:

```json
{
  "response": { "set": [ … ], "filter": { … } },
  "refs": {
    "user":       { "506815840423837697": "Ann Smith" },
    "record":     { "900000000000000007": "Acme Ltd" },
    "attachment": { "514227795113148417": { "name": "offer.pdf", "url": "…", "previewUrl": "…", "mimetype": "application/pdf", "size": 30211 } },
    "role":       { "…": "Sales" }
  },
  "refsTruncated": false
}
```

- The resource shape never changes; read and list carry `refs` the same way.
- Buckets are keyed by reference kind, so a client knows what an ID is.
- A reference the caller cannot read, or that no longer exists, is **absent**.
  A client shows its fallback (the ID) exactly as it does today.
- `refsTruncated` is set when a cap was hit (see _Limits_).
- An attachment's entry is the one bucket that is not a string: its URL is
  signed per caller and a viewer needs it to show a thumbnail.

### Label rules — one rule per kind, owned by the server

| Kind                                               | Label                                                                                                                                                                                                                                                          |
| -------------------------------------------------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| User (every user reference, including `deletedBy`) | name → email → handle, else absent. `username` is never set (no field in the admin user editor, 0 of 63 dev users). Masking applies: a masked name or email is not a label, so fall through to the next part                                                   |
| Record                                             | `labelField` if it resolves, else the module's first field; when that field is itself a Record, follow `recordLabelField` one more level (two levels, as the agent tools do). Multi-value label field → first value. Select → the option's text, not its value |
| Attachment                                         | `{name, url, previewUrl, mimetype, size}`                                                                                                                                                                                                                      |
| Other resources                                    | The resource's display name (`name`, or `meta.name` for the 11 resources that keep it there), else its handle                                                                                                                                                  |

The webapp viewers switch to the server's label when it is present, so every
screen reads one rule. This changes `CFieldRecordViewer`'s formatting for
non-string label fields; `lib/vue/src/components/field/field.intent.md`
governs the viewers, so that change goes through `/intent-task`.

### Architecture (server)

```
request ──► refs middleware (reads ?refs=, puts a Collector in ctx)
              │
controller ──► refs.Collect(ctx, x)      generated List/Read + hand-written payload builders
              │   typed resources: walk dal.Model attributes of Type TypeRef
              │   records: module fields of kind Record/User/File + system user fields
              ▼
api.encode ──► collector.Resolve(ctx) ──► resolver registry, keyed by ResourceType
                                           corteza::system:user      → system user service
                                           corteza::compose:record   → record service, per (ns, module)
                                           compose attachment        → attachment service
                                           corteza::system:role, …   → their services
           ◄── {"response": …, "refs": …}
```

- **Collector.** A per-request set of IDs by kind (and by module for records).
  Nothing is collected unless `?refs` is present, so the default path costs a
  context lookup.
- **Knowing what is a reference.** Generated models already mark references:
  `&dal.Attribute{Ident: "CreatedBy", Type: &dal.TypeRef{RefModel:
&dal.ModelRef{ResourceType: "corteza::system:user"}}}`
  (`server/compose/model/models.gen.go:1258`), reachable through each
  component's `model.Models()`. A walker keyed on the Go field name covers
  about 60 types with no new codegen. Not covered yet: references stored as a
  plain `ID` type (tenant/project membership, project approval fields), IDs
  nested in JSON config (agent, chatbot, application), and one wrong target
  (`DataPrivacyRequestComment.RequestID` is declared a user reference). Those
  need annotating in cue before the walker sees them.
- **Emitting.** Every generated handler ends in `api.Send`, and `encode()`
  builds the envelope in one place (`server/pkg/api/response.go:87`). That is
  where `refs` is added; it already has the request.
- **Resolving.** Always through the services — never `store.*` directly —
  so record read, field-level value read, user read and privacy masking all
  hold. One batch per kind (per module for records), fetched by **ID set**, not
  by an OR chain in the query language.
- **Prerequisite:** `RecordFilter` gains an ID-set filter that becomes a DAL
  constraint (`id IN (…)`). Spiked in the worktree with a dozen lines
  (`dalutils.prepFilter`); it is what took 200-row resolution from ~60 ms to
  ~18 ms. The agent tools, the webapp's label batches and export all benefit.

### Limits

- A cap on distinct references per response (the agent tools use 500), per
  kind, with `refsTruncated` set when hit.
- No cross-request cache: labels depend on the caller's permissions and on
  masking.
- Nested labels stop at two levels.

### Client

- **lib/js:** the generated client's `stdResolve` returns only
  `response.data.response`, so `refs` is dropped today
  (`lib/js/tools/codegen/template.js:61`). The template keeps it — returned
  alongside the response when the call asked for `refs`. The parameter needs a
  "common parameters" feature in `rest.yaml` and both generators
  (`server/pkg/codegen/rest.go`, `lib/js/tools/codegen/human-api-client.js`);
  neither has file-level parameters today.
- **lib/vue:** a small label store keyed by `(kind, id)`, seeded from any
  response that carries `refs`. Viewers read it first and fall back to today's
  resolution. The RecordList block asks for `refs=labels` and the per-row label
  calls disappear.

### Sorting by the displayed label

Today `sort=company` on a Record field compiles to
`CAST(CASE WHEN values->'company'->>0 ~ '^[0-9]+$' THEN … END AS BIGINT)` — the
list is ordered by the referenced record's ID, so the labels on screen look
unsorted. A multi-value field sorts by its first value, so a multi-value
reference sorts by its first ID. A dotted path fails with `unknown attribute`.

**Syntax.** `sort=label(company) ASC`, `sort=label(ownedBy) DESC`. A field
name cannot contain parentheses, so it never collides, and `sort=company`
keeps meaning the ID.

**Allowed when:**

- the field is a single-value Record or User field, or a system user field;
- the target lives in the same database connection as the module (a module on
  its own connection cannot join `compose_record` or `users`);
- the caller can search the target module and read its label field. Field-level
  read is decided per field, not per record (`ComposeRecordFilterAC`,
  `server/compose/service/record.go:2536`), so it is one check per request.

**Query.** The rdbms driver adds a `LEFT JOIN` on the extracted reference and
orders by the target's label expression, then by `id`, as it does today. The
user label is the same `COALESCE` the label rule defines.

**Cursor.** The page cursor carries the sort values read off the last row
(`server/pkg/filter/pagination.go:222`), and a joined label is not on the row.
Three places read it as a row attribute and fail
(`server/store/adapters/rdbms/dal/iterator.go:147,304`, the table codec lookup). The label has
to be scanned as an extra column and registered as a virtual sort attribute.

**Record-level read.** Record read is decided per record in Go after the fetch
(`server/compose/dalutils/records.go:216`), because contextual roles read record
values. A plain join would order rows by a label the caller cannot see, so the
order would leak it. Two paths:

- **Read is uniform** for the caller on the target module → the join.
- **Otherwise** → resolve the readable targets' labels through the service,
  rank them in Go, and sort by that rank. An unreadable target sorts as empty.
  Measured at 719 ms for 2,000 targets over 1M rows; capped by target count,
  above which the column is not sortable by label.

Deciding "uniform" is itself an RBAC question (open ruling 2).

**Stored sort key, later.** Where a module is measured too big for the join, a
per-reference label stored on the row with a `NULLS FIRST` index answers in
0.1 ms at 1M. The cost moves to writes: every rename of a target rewrites every
row pointing at it (measured above), a change of a field's `labelField`
rewrites the module, and it needs a backfill. Build it only against a measured
need.

---

## Plan

Effort is an estimate from what was read, not a commitment.

| Phase | Work                                                                                                                                                                                                                                                              | Size              |
| ----- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ----------------- |
| 0     | ID-set filter on `RecordFilter` → DAL `IN`; switch agent tools, webapp label batches and export to it                                                                                                                                                             | S                 |
| 1     | `refs` middleware, collector, `encode()` sibling, user resolver; record system users and User values; `TypeRef` walker for typed resources; common `refs` param in both generators; lib/js keeps `refs`; lib/vue label store; RecordList and user viewers read it | M–L               |
| 2     | Record labels (two levels, Select text), attachments, the ruled user label (agent tools, export and webapp all change), viewers switched to server labels (`/intent-task`, `field.intent.md`)                                                                     | M                 |
| 3     | Other resource kinds (role, namespace, module, workflow, …); annotate plain-ID and JSON-nested references in cue                                                                                                                                                  | M                 |
| 4     | `label(field)` sort: rdbms join, virtual cursor attribute, readability gate, rank fallback; RecordList enables it for single-value reference columns                                                                                                              | L                 |
| 5     | Stored sort key for modules measured too big                                                                                                                                                                                                                      | L, only if needed |

---

## Open rulings

1. **Record label formatting.** Server returns a string; the viewer currently
   formats dates, numbers and select text. Which kinds does the server format?
2. **"Read is uniform"** — which RBAC evaluation decides that a caller can
   read every record of a module, given contextual roles and that `rbac.Can`
   answers false for any wildcard resource.
3. **Deleted targets.** Both server paths exclude deleted records; the
   webapp's per-ID fallback shows them. Should a label of a deleted target be
   returned?

---

## Found on the way (filed separately)

| Backlog     | Finding                                                                                                                                                                                                                                                                            |
| ----------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `be1a3f34e` | **Security**: a record attachment downloads with no login and no signature if `record` is changed to `page` in its URL — `isAccessible` treats page/icon/namespace kinds as public and serving loads by ID regardless of kind (`server/compose/rest/attachment.go:95`). Reproduced |
| `b7f1681bf` | Attachment `FindByID`/`DeleteByID` have no RBAC; `Search`'s page and record checks are inverted (`server/compose/service/attachment.go:95,104,111`)                                                                                                                                |
| `babe761eb` | Record export `resolveRefs` bypasses record/user read and masking, queries per row, can misplace labels                                                                                                                                                                            |
| `b072b2825` | RecordList always sends `incTotal=true`; the count walks the whole module in Go — 211 s at 1M rows                                                                                                                                                                                 |
| `b3ac01c23` | `useUserStore.resolveUsers` has no in-flight dedupe — 29 calls for 6 users on one page                                                                                                                                                                                             |

---

## Method

Seed (worktree database, after creating the Company and Contact modules
through the API):

```sql
insert into compose_record (id, rel_module, rel_namespace, values, meta, created_at, owned_by, created_by)
select 900000000000000000 + g, :company, :ns,
  jsonb_build_object('name', jsonb_build_array(initcap(substr(md5(g::text),1,8)) || ' Ltd')),
  '{}', now(), :user, :user
from generate_series(1, 2000) g;

with uu as (select array_agg(id) a from users where deleted_at is null)
insert into compose_record (id, rel_module, rel_namespace, values, meta, created_at, owned_by, created_by)
select 910000000000000000 + g, :contact, :ns,
  jsonb_build_object(
    'name',    jsonb_build_array('Contact ' || substr(md5('c'||g),1,10)),
    'company', jsonb_build_array((900000000000000001 + (hashint4(g)::bigint & 2147483647) % 2000)::text),
    'owner',   jsonb_build_array((uu.a[1 + (hashint4(g+7)::bigint & 2147483647) % array_length(uu.a,1)])::text)),
  '{}', now() - (g || ' seconds')::interval,
  uu.a[1 + (hashint4(g+3)::bigint & 2147483647) % array_length(uu.a,1)],
  uu.a[1 + (hashint4(g+5)::bigint & 2147483647) % array_length(uu.a,1)]
from generate_series(1, 100000) g, uu;
```

Sort by label (the join measured above):

```sql
select r.id from compose_record r
left join compose_record c
  on c.id = (case when r.values->'company'->>0 ~ '^[0-9]+$' then (r.values->'company'->>0)::numeric end)
 and c.rel_module = :company and c.deleted_at is null
where r.rel_module = :contact and r.rel_namespace = :ns and r.deleted_at is null
order by c.values->'name'->>0 asc nulls first, r.id
limit 50;
```

API timings: median of 15 authenticated `GET`s after one warm-up. In-process
label timings: a timer around `refLabels` in a scratch build, called through
`compose_record_lookup`. Query text: `log_min_duration_statement = 0` on the
worktree database. Browser request counts: `dev/agent/drive.mjs` recording
every `/api/` request on a page built with `pagebuild.py`.
