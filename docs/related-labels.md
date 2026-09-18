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

| #   | Ruling                                                                                                                                                                                                                                              |
| --- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| 1   | Scope is **records**: record list and read carry labels for Record/User/File values and the record's system users, in one implementation shared with the agent tools. Other resources wait for a client that needs them (revised from "everywhere") |
| 2   | A resolved reference is a **label**, not the related object                                                                                                                                                                                         |
| 3   | Sort means sort by what is displayed (the viewer's label), not by an arbitrary path                                                                                                                                                                 |
| 4   | Labels ride as a **typed sibling of `response`** in the envelope                                                                                                                                                                                    |
| 5   | Opt-in is a **query parameter on every endpoint**                                                                                                                                                                                                   |
| 6   | Sorting is a **query-time join first**; a stored sort key only where a module is measured too big                                                                                                                                                   |
| 7   | The user label is **name → email → handle** everywhere; `username` is left out                                                                                                                                                                      |
| 8   | A record's entry carries its label as a string **and** as its label field's raw values; the webapp formats the values itself                                                                                                                        |
| 9   | Nested record labels stay as they are: lists follow each level's own `labelField` (the server mirrors this), the picker keeps `recordLabelField`                                                                                                    |
| 10  | Every text sort — label sorts and today's String-column sorts — uses the ICU collation `und-x-icu` (revised from "as today")                                                                                                                        |
| 11  | A deleted target's label is returned, marked deleted, and the webapp shows it as deleted                                                                                                                                                            |
| 12  | The record count RecordList asks for is fixed before label sort                                                                                                                                                                                     |

---

## Verdict

- **Labels in the response: moderate work, cheap at runtime — built for
  records only.** There is one choke point to emit them (`api.encode`), and the
  agent tools already resolve labels for records; REST and the agent tools
  share one implementation. With an ID-set filter, resolving a 50-row page
  costs 8–10 ms beside a 75 ms list call. Extending it to the other ~60
  resource types is possible (their references are already machine-readable)
  but waits for a client that needs it.
- **The webapp keeps its formatting.** Record labels are drawn by the label
  field's own viewer (dates, numbers, select badges, translations), so a
  record's entry carries its label field's raw values as well as a string, and
  the refs land in a store of their own — the existing caches hand their
  entries out as full records.
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
`labels` (or no value) means every kind; a comma list limits the kinds. It is
declared on the record list and read endpoints in `server/compose/rest.yaml`
only, so neither code generator needs a common-parameters feature.

The envelope gains a sibling of `response`:

```json
{
  "response": { "set": [ … ], "filter": { … } },
  "refs": {
    "user":       { "506815840423837697": "Ann Smith" },
    "record":     { "900000000000000007": { "moduleID": "…", "label": "Acme Ltd", "values": { "name": ["Acme Ltd"] } },
                    "900000000000000011": { "moduleID": "…", "label": "Old Co", "values": { "name": ["Old Co"] }, "deleted": true } },
    "attachment": { "514227795113148417": { "name": "offer.pdf", "url": "…", "previewUrl": "…", "mimetype": "application/pdf", "size": 30211 } },
    "role":       { "…": "Sales" }
  },
  "refsTruncated": false
}
```

- The resource shape never changes; read and list carry `refs` the same way.
- Buckets are keyed by reference kind, so a client knows what an ID is.
- A reference the caller cannot read, or that does not exist at all, is
  **absent**; a client shows its fallback (the ID) exactly as it does today.
- A deleted record the caller could read is returned with `"deleted": true`.
  Listing already takes `deleted=1` under the same search and read checks
  (`server/compose/rest/record.go:112`), so including deleted targets widens
  nothing.
- `refsTruncated` is set when a cap was hit (see _Limits_).
- A record's entry carries its label twice: `label`, a plain string for a
  client that only wants text, and `values`, the raw value(s) of its label
  field, which the webapp formats itself (see _Frontend_).
- An attachment's entry is the fields the file viewer draws: its URL is signed
  per caller and a viewer needs it to show a thumbnail.

### Label rules — one rule per kind, owned by the server

| Kind                                               | Label                                                                                                                                                                                                                                                                                                                                                                                                                                                                 |
| -------------------------------------------------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| User (every user reference, including `deletedBy`) | name → email → handle, else absent. `username` is never set (no field in the admin user editor, 0 of 63 dev users). Masking applies: a masked name or email is not a label, so fall through to the next part                                                                                                                                                                                                                                                          |
| Record                                             | `labelField` if it resolves, else the module's first field — the viewer's rule. When that field is itself a Record, the nested record gets its own entry, found through that field's own `labelField`, as the viewer recurses; three levels at most; `recordLabelField` is not followed (the picker keeps it). `values` holds only the label field. Deleted targets are included and marked. `label` is its first value, Select as the option's text, dates as stored |
| Attachment                                         | `{name, url, previewUrl, mimetype, size}`                                                                                                                                                                                                                                                                                                                                                                                                                             |
| Other resources                                    | The resource's display name (`name`, or `meta.name` for the 11 resources that keep it there), else its handle — when refs reach them (later)                                                                                                                                                                                                                                                                                                                          |

The webapp does not render `label`: it formats `values` through the same field
viewers it uses today, so dates, numbers, select badges and translations stay
in the browser. The string is for other clients and is what a label sort
orders by.

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
- **Records first.** Ruling 1 builds only the record path: module fields of
  kind Record/User/File plus the system user fields, collected in the record
  controller and resolved by the same code the agent tools' `refs` uses today.
  The typed-resource walker below is the route to the rest, when it is wanted.
- **Knowing what is a reference (later).** Generated models already mark references:
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
- Nested labels stop at three levels (the viewer has no limit).

### Client

See _Frontend_ below: what the webapp does today, where the refs land, and
which screens change.

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

Deciding "uniform" is itself an RBAC question (open ruling 1).

**What "displayed order" means per label kind.** Text sorts as text (see
collation below); Number and DateTime sort by value, which matches what is
shown; a Select label is shown as its option text, so a label sort orders by
that text (a `CASE` over the field's options), not by the stored value — plain
Select columns sort by stored value today, the same mismatch; a Record label
adds one join per level; a User label joins `users`.

**Collation.** Every text sort names `COLLATE "und-x-icu"` (ruling 10) — the
label sort and today's String-column sort alike, so two columns on one list
order the same way. The dev database collates `C.UTF-8`, byte order:
`Banana < apple < cherry < zebra < Ćevapi`; ICU gives
`apple < Banana < Čas < Ćevapi < cherry < zebra`. It goes in the Postgres
dialect's text sort expression: sqlite (the integration tests) has no ICU
collation and keeps its own order, and mysql needs its own equivalent. An index
built for a text sort must name the same collation. Existing String columns
change order when this lands — visibly, as a fix.

**Deleted targets** sort by their label like any other (ruling 11): the join
does not filter on the target's `deleted_at`.

**In the webapp.** RecordList builds `sort` from column names
(`client/web/unify/src/sections/compose/lib/record-sort.js:4-11`). A Record,
User or system-user column sends `label(field)`, and `parseSortExpression` maps
`label(field)` back to the column so the header shows it as sorted.

**Stored sort key, later.** Where a module is measured too big for the join, a
per-reference label stored on the row with a `NULLS FIRST` index (in the same
ICU collation) answers in
0.1 ms at 1M. The cost moves to writes: every rename of a target rewrites every
row pointing at it (measured above), a change of a field's `labelField`
rewrites the module, and it needs a backfill. Build it only against a measured
need.

---

## Frontend

### How the webapp resolves references today

| Where                                                                                                             | How                                                                                                                                                                                                                                                                                         | Cost                                                                                                                                                                  |
| ----------------------------------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| Record field viewer, `lib/vue/src/components/field/viewers/CFieldRecordViewer.vue`                                | `recordStore.getByID(id)` (`:101`), else `resolveRecordLabels` (`:152-164`). Renders the label with `<CFieldViewer :field="labelFieldDef" :record="rec">` (`:12-18`), so the label field's own viewer formats it. Shows the ID when the target module is not in the module store (`:81-94`) | One `recordList` per (namespace, module) per tick, 100 IDs each, then a `recordRead` per ID the batch did not return (`lib/vue/src/stores/useRecordStore.js:206-298`) |
| User field viewer, `CFieldUserViewer.vue:61-73`                                                                   | `findCached(id)`, else `resolveUsers`. Label `name → handle → email → ID` (`lib/vue/src/composables/useUserResolver.ts:15-18`)                                                                                                                                                              | No in-flight dedupe across viewers; the shell preloads 500 users on mount (`client/web/unify/src/App.vue:319`)                                                        |
| File field viewer, `CFieldFileViewer.vue:175-220`                                                                 | `attachmentRead` per ID                                                                                                                                                                                                                                                                     | One call per file, one after another, no cache                                                                                                                        |
| RecordList, `RecordListBlock.vue` `resolveLabels`                                                                 | Resolves each Record and User column before the rows are swapped in (added in `b069022ca`, 2026-09-18)                                                                                                                                                                                      | One record batch per referenced module, one user batch per User column                                                                                                |
| Chart, `ChartRenderer.vue:103-166`                                                                                | Bool, Select and Record dimensions get labels; the Record label is the raw first value                                                                                                                                                                                                      | User dimensions show IDs                                                                                                                                              |
| Calendar and map feeds, `lib/js/src/compose/types/page-block/calendar/feed-record.ts:62`, `GeometryBlock.vue:173` | Title is the title field's raw first value, else the record ID                                                                                                                                                                                                                              | A Record or User title field shows an ID                                                                                                                              |
| Pickers, `CFieldRecordEditor.vue:96-139`, `CInputRecord.vue:133-160`                                              | Their own rule: raw `labelField` value, `recordLabelField` for a nested record, else the first non-empty value, else `Record <id>`                                                                                                                                                          | Fetch full records for their options                                                                                                                                  |

The browser count in _Measurements_ (29 user calls for 6 users) was taken
before `b069022ca` added the RecordList pre-pass; re-measure before quoting the
saving.

### Formatting stays in the browser

A record's label is drawn by the label field's own viewer, and every one of
them formats:

- DateTime: moment formats, relative time, date-only or time-only, in the
  user's locale (`lib/js/src/compose/types/module-field/datetime.ts:60-79`);
- Number: prefix, suffix, precision (`formatValue`);
- Bool: the field's true/false labels or the translated yes/no
  (`CFieldBoolViewer.vue:48-54`);
- Select: the option's text, optionally as a coloured badge
  (`CFieldSelectViewer.vue`);
- User and Record: another resolution, recursively.

A string from the server cannot carry any of that, so for records the server
sends the label field's raw values and the webapp keeps formatting them. Users
and attachments have no such formatting and stay a string and a payload.

### Where the refs land: a store of their own

Not in `labelCache`, `records` or the user store:

- `getByID` returns `labelCache` entries as if they were full records, and two
  callers use them that way — the comment block's reply modal
  (`CommentBlock.vue:529,941`) and the compose context builder
  (`client/web/unify/src/sections/compose/index.js:59`), which copies a record's
  values into an agent's context. A label-only entry there shows as a truncated
  comment or an agent context with most values missing.
- The user store holds whole `system.User` objects that admin views read.

A new `useRefStore` in `lib/vue/src/stores`:

- `record(id)`, `userLabel(id)`, `attachment(id)`;
- `absent(kind, id)` for a reference the response asked about and the server
  left out (unreadable, or never existed). Without it a viewer falls back to fetching —
  and the per-ID `recordRead` fallback asks again for exactly what the server
  withheld;
- seeded in one place: the generated client keeps `refs` (today `stdResolve`
  drops it, `lib/js/tools/codegen/template.js:61`) and hands it to a callback
  the shell registers. A call site only adds `refs: 'labels'`;
- an entry is evicted when the record store updates or deletes that record, and
  the store clears with `recordStore.clearAll`.

### Viewer changes

The props do not change, as `field.intent.md` requires ("keep dispatcher prop
sets stable").

- `CFieldRecordViewer`: `getByID(id) || refStore.record(id) || { recordID }`,
  and no `resolveRecordLabels` for an ID the ref store knows or knows is absent.
  An entry marked `deleted` draws its label as deleted (muted, with a
  translated "deleted" hint — a new key in `locale/en/human-webapp/`).
- `CFieldUserViewer`: the ref store's label when the user is not cached; no
  `resolveUsers` for those.
- `CFieldFileViewer`: the ref store's entry before `attachmentRead`.
- The field contract's record-shape bullet gains ref entries. Its raw-array
  guard holds only while the target module stays unloaded: the template checks
  `labelFieldDef`, not the shape (`CFieldRecordViewer.vue:13`), so a raw entry
  cached before the module loads renders blank after. Both go through
  `/intent-task`, with a new `useRefStore.intent.md`.

### Which screens change

| Screen                                  | Today                                                    | With refs                                                                    |
| --------------------------------------- | -------------------------------------------------------- | ---------------------------------------------------------------------------- |
| RecordList                              | List, then a batch per referenced module and User column | One call; the pre-pass goes                                                  |
| Record page, record organizer, comments | Viewers resolve per tick                                 | Refs from the read or list they already make                                 |
| Calendar and map feeds                  | Record/User titles show IDs                              | Titles show labels                                                           |
| Chart                                   | User dimensions show IDs                                 | Report refs (the agent tools' `dimensionRefs` already does this server-side) |
| Pickers and editors                     | Fetch full records for their options                     | Unchanged                                                                    |

### Rules the webapp converges on

- Users: `formatUser` goes from `name → handle → email` to the ruled
  `name → email → handle`. Visible: a user with no name shows their email
  rather than their handle. The comment block (`CommentBlock.vue:376-378`), the
  session view and group members have orders of their own and move too.
- Records stay as they are (ruling 9): the server follows the viewer, the
  pickers keep `recordLabelField`, and the chart keeps reading the raw first
  value — which is the same text the server's `label` carries.

---

## Plan

Step 0 left the webapp's per-ID label reads only for a batch that fails; deleted
targets now come in the batch. `useRecordStore.intent.md` still describes the
old fallback (queued for an intent audit).

Effort is an estimate from what was read, not a commitment.

| Phase | Work                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                       | Size |
| ----- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ---- |
| 0     | **Done** (2026-09-18). `RecordFilter.RecordID` → DAL `IN`, and `recordID[]` on record list, patch, bulk delete/undelete and export. The agent tools, the webapp's label batches, export labels and import's existing-record check use it. Also fixed: bulk edit and selected-row export (both sent a `recordID IN (…)` the server rejects), and export labels placed by result order rather than by ID. Bulk delete and restore now take one request, user lookups share in-flight calls, and file lookups run in parallel | S    |
| 1     | The record count RecordList asks for (`b072b2825`): stop walking the whole module in Go — count in SQL where read access allows it, or cap it                                                                                                                                                                                                                                                                                                                                                                              | M    |
| 2     | `label(field)` sort: rdbms join, virtual cursor attribute, readability gate, rank fallback, Select-as-text, deleted targets included; `und-x-icu` on every text sort; RecordList sends it for reference columns and maps it back to the header                                                                                                                                                                                                                                                                             | L    |
| 3     | Record `refs`: `refs` on record list/read in `rest.yaml`, collector and `encode()` sibling, one resolver shared with the agent tools; entries with `label` + label-field `values`, nested, deleted marked; users (ruled order) and attachments; lib/js keeps `refs`; `useRefStore`; record, user and file viewers read it; RecordList, record page, organizer, comments, calendar/map feeds and chart opt in (`/intent-task`: `field.intent.md`, new `useRefStore.intent.md`)                                              | M–L  |
| later | Refs for other resources: common `refs` param in both generators, `TypeRef` walker, other kinds, cue annotations for plain-ID and JSON-nested references — when a client needs them                                                                                                                                                                                                                                                                                                                                        | M–L  |
| later | Stored sort key for modules measured too big                                                                                                                                                                                                                                                                                                                                                                                                                                                                               | L    |

---

## Open rulings

1. **"Read is uniform"** — which RBAC evaluation decides that a caller can
   read every record of a module, given contextual roles and that `rbac.Can`
   answers false for any wildcard resource.
2. **Deleted and suspended users.** Ruling 11 covers records; do users follow
   it (label returned, marked)?

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
