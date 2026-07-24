---
name: intent-audit
description: Semantic drift audit for an intent-covered area — verify docs against code beyond what hash-sync can see, report contradictions with evidence, interview the human on rulings, fix the wrong side. Run per area on a cadence and before releases.
---

# /intent-audit [path]

The semantic layer of the intent system: hashes prove docs were _acknowledged_;
this audit proves they are _true_.

## Procedure

1. **Scope**: use the given path; without one, pick the oldest area from
   `node .intent/intent.mjs status` drift-age. Confirm scope with the human if
   it exceeds one section/package.
2. **Gather**: list the area's intent docs and covered files. For large areas
   fan out read-only subagents (model: sonnet, one per subfolder), all
   launched in one message; each returns findings with `file:line` evidence
   and verbatim doc excerpts, never prose summaries of code. Don't barrier:
   start verifying (step 3) a subfolder's docs as soon as its gather returns,
   while the other gathers still run.
3. **Verify, per doc** — main-agent work; subagent evidence feeds it but the
   verdicts are never delegated:
   - Every stated contract against actual code behavior.
   - **Locked contracts** (area constitutions): does code still honor them?
   - **WIP / DRIFT notes**: still accurate? (Resolved drift must lose its note;
     new drift must gain one.)
   - Frontmatter `depends-on` / `touched-by` / `tests` paths exist and are
     still the real relationships.
   - Dead references: docs describing files/exports that no longer exist, and
     orphaned code no doc claims.
4. **Classify findings**: code-wrong (drift from recorded intent) /
   doc-wrong (stale doc) / ruling-needed (docs and code disagree and the truth
   is a human call).
5. **Resolve**: interview ruling-needed items via AskUserQuestion (concrete
   options, evidence quoted). Fix doc-wrong directly. Fix code-wrong only with
   explicit approval — record a `> **DRIFT:**` note instead if deferred.
6. **Close**: prettier changed docs, `sync`, `check` green (read exit codes),
   commit (audit fixes separate from any code fixes), update `.intent/TODO.md`,
   report per-doc verdicts and what changed.
