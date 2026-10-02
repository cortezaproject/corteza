# Project section gotchas

Rulings and non-obvious behaviour of projects: lifecycle, revisions and branch copies, AI systems and FRIA, deletion and isolation.

## Publish and the status model

- The live state is `active`, not `published`. `published` is a valid, accessible `ProjectStatus` value that no server code sets.
- Publishing a revision activates it and marks its parent `deprecated`. `archived`, `suspended` and `deprecated` projects are inaccessible (`project_resolver.go`); archiving lives on its own `archived_at` field, not on status.
- `ParentRevisionID == 0` marks an original. Its first publish only flips draft to active: no migration, no namespace swap, and an empty deployment plan.
- `RootProjectID()` falls back to the project's own ID.
- The revision design doc is `server/devdocs/project-revisions.md`.
- Never hand-revert codegen output. If `make codegen` produces an unrelated change, keep it.

## Lifecycle behaviour

- Lifecycle endpoints carry their own ops, `project.revise` and `project.publish`. The graph and the deployment plan are gated on read.
- `PUT /projects/{id}` cannot change status. Create forces `draft`, and archiving goes through `/archive` and `/unarchive`.
- One draft per chain. An archived draft still counts against the gate.
- Publishing over a deleted parent revision is refused.
- Publish with no mappings fills them from the plan. `discardRecords: true` is the explicit opt-out.
- Record migration runs as one `dml.Migration` across all mappings. Record links are remapped in a second pass through an old-to-new id map. Each migrated row is stamped with its source id, so a retry clears only migrated rows and hand-typed draft records survive. `ownedBy`, `createdBy`, `createdAt` and `updatedAt` are written straight to the DAL.
- Value-level rejections land in `ValueError`, not `Error`. A run that only checks `Error` counts rejected rows as processed.
- Work items (incidents, tasks, features, privacy, backlog, reviews) are chain-wide: reads and writes both resolve to the chain root through `rootProjectID` (`system/service/project.go`).
- A compose resource belongs to its namespace's project. `rel_project` is inherited from the namespace and repointed when a revision is branched.

## Lifecycle rulings

1. Work items stay chain-wide. They are never per-revision.
2. Members are chain-wide (keyed by `RootProjectID()`). Roles are per revision, copied on branch with their RBAC rules rewritten onto the clone's resource ids.
3. A copied role's handle is `proj_<newRevisionID>_<suffix>` (`projectRoleHandle`). `diffSources` compares project roles by suffix, so the plan reports no role churn.
4. Chain identity belongs to the root: the live namespace slug is always the root's handle, and revision handles carry `-revN` without compounding.
5. Migration stays outside the flip transaction. Retries are made safe by idempotent import, not by locking the revision.
6. Archiving disables the project's namespace and unarchiving re-enables it. Both write an action-log entry.
7. Record ids change on migration. Record revisions and attachments are not carried.
8. The project list's status chip handles `active` and `deprecated`.
9. Intent docs for the project section are reconciled only through `/intent-task`.

## Publish approval

- Approval state lives on the revision row, owned by the server like status: `approval_status` (draft, submitted, approved, rejected), the plan fingerprint, and who submitted and decided it, with when. History comes from the action log.
- `CanGrantApproval` decides, and the submitter can never approve their own request. Anything going live needs two people.
- Approval stores a fingerprint of the deployment plan. Publish recomputes it and refuses on mismatch (`ProjectErrApprovalStale`), so no write hook has to expire approvals.
- The first publish needs approval too.
- Publish needs both an approved revision and `project.publish`.
- A branch does not inherit its parent's approval.
- Multi-approver quorums are out of scope. `ProjectReview` is a scheduled-review work item, not an approval record.

## Branch copy

- `CreateRevision` clones the compose namespace (modules, pages, layouts, charts) through envoy. `project_revision_clone.go` copies the project-scoped system resources by hand, because envoy has no scope resolution or matchup outside compose.
- Copy order is connections, TAQs, agents, chatbots, AI systems (with entries and FRIA scenarios), reviews, roles. TAQs are copied before agents, and their agent refs are repointed afterwards.
- Every reference falls into one of three cases:
  - Copied with the revision: remapped through `projectCloneIDs`. Modules are matched by handle, because the compose clone returns no id map.
  - Shared (knowledge bases, LLM providers, users, service accounts): left untouched.
  - Owned by the parent but not copied (classic workflows): dropped and logged.
- An agent access entry with an empty `ModuleIDs` covers every module in its namespace (`checkAllow`, `system/agentic/policy`). An entry whose modules all drop out must be removed whole.
- A shallow struct copy of an agent shares slices with the live agent, so use `types.Agent.Clone()`.
- A copied chatbot gets a fresh widget key, because a widget key is the public embed id.
- The draft namespace is created disabled.
- `WorkflowFilter.ProjectID` is not in `byValue`, so a project-filtered workflow search returns every workflow. Probe ownership per ID.

## AI systems and FRIA

1. The regulated unit is the "AI system" (Art. 3(1)), never a "group". In Art. 27, "groups" means groups of persons, which is the `vulnerableGroups` taxonomy.
2. Storage is `ProjectAiSystem` and its entries (`server/system/project_ai_system.cue`, `project_ai_system_entry.cue`). AI systems are per revision and copied on branch.
3. An AI system is a Govern-tab entity, not a resource kind: it has no `kinds.js` entry, no graph node and no permission-matrix row.
4. A FRIA scenario has exactly one AI system (`aiSystemID`, required).
5. The three deployer questions sit at revision level. Intended purpose and risk class are per AI system.
6. `risk_class` uses the Art. 6 tiers (prohibited, high, limited, minimal) and is a filterable column. Intended purpose lives in `ProjectAiSystemMeta`.
7. `resource_ref` uses the RBAC ref format (`corteza::compose:module/<id>`).
8. When a member resource is deleted, its entry stays, renders as a tombstone and triggers reassessment.
9. Oversight roles (Art. 14) are entries in the same table, told apart by ref kind.
10. Membership is edited in both a graph select mode and a checklist, over one shared state.

Not enforced yet:

- `risk_class` is not validated.
- `prohibited` does not block publish.
- The service merges only non-empty fields, so `riskClass` cannot be cleared.
- `chart` is in `MEMBER_KINDS` (`config/resourceRefs.js`) with no `KIND_SOURCES` entry in `AiSystemEditor.vue`, so it always lists empty.
- Governance step state is session-local (`governanceByProject` in `stores/projects.js`).

## Deletion and the status lock

1. A bound namespace cannot be deleted while it is bound. It goes only with its project. Standalone namespaces (`ProjectID` 0) delete normally.
2. Only a draft project can be deleted. A non-draft project is archived instead, and archive is terminal for a published project.
3. Deleting a project soft-deletes its sub-records and every built resource it owns. Each cascaded row gets the project's exact `DeletedAt`, and undelete clears only rows with that timestamp.
4. Records, attachments, sessions, conversations and translations stay out of the cascade.
5. The project-delete permission covers the whole cascade.
6. `can*` flags stay RBAC-only. They can advertise a delete that the status lock then refuses.
7. Compose create, update and delete on a non-draft project fail with `project.errors.locked`. Records stay writable.

Rules 1 to 5 are not implemented. `onDelete` (`system/service/project.go`) stamps only the project and its namespace, has no status check, and the generated `guard` returns nil. The delete confirmation in `project.yaml` still says "can't be undone", although the delete is soft.

## Isolation rulings

Status: decided, not implemented.

- RBAC is the single authority. Non-members get a 404 by default, but an explicit system rule can open parts of a project (to an auditor, say).
- Project roles are contextual-style roles. They apply only to that project's resources and only to users assigned to them. Membership means holding at least one project role, and every project gets default roles that include "End user".
- Only bypass roles (`pkg/rbac/ruleset_checks.go:25`) skip the gate. `admin` keeps its wildcard rules.
- The gate checks the resource's own `rel_project`, not the URL, so background paths are covered too.
- Ownership is a property (`rel_project`, where 0 means global), never part of a URL or RBAC path. A project uses its own resources and global ones, never another project's.
- Handles are unique across a project and global, with no shadowing. Externally addressed names (API gateway endpoint, queue name, widget key) stay instance-unique.
- Background runs carry the project of the thing being run. The run-as user must hold a role in that project. Cross-project triggers work when that identity's RBAC can read the target.
- Revisions copy role members at branch time.
- Project mail goes through an SMTP project connection only, with no system SMTP fallback. Login, identity providers and sessions stay global.
- Branding, upload limits and record and chart UI become per-project resources or properties, not a settings table. Branding reaches the project's pages only.
- In-project pickers show the project's own resources, with global ones behind a toggle. Global resources are managed in admin only.

## Projects webapp: server versus client state

- The server owns project CRUD, members, governance step state (`/projects/{id}/governance/{stepKey}`), AI systems, the resource graph and `ProjectConfig`.
- Creating a project also creates its compose namespace (slug = handle, `config.namespaceID`) and makes the creator a developer member (`system/service/project.go`).
- Wizard modules are real compose modules in the project namespace.
- Resource-management values persist as generic governance step values (`/governance/resource-management`), not through `ProjectConfig.ResourceManagement`.
- Sensitivity levels are real global `dalSensitivityLevel` resources. `ensureStandardLevels` (`stores/projects.js`) seeds the scheme from `config/sensitivity.js`, and each level needs a unique `level` int. Only fields carry sensitivity. Modules do not.
- When the server's capabilities disagree with the FE's `roles.js` flags, the server wins.
- `scopedTables()` (`store/adapters/rdbms/upgrade_tenancy.go`) is a handwritten list. It must name every model that carries `rel_tenant` or `rel_project` scope columns.
- Codegen chain: `server/system/rest.yaml` feeds `make codegen-legacy`, `system/*.cue` feeds `make codegen`, and `lib/js` codegen produces the API clients.
