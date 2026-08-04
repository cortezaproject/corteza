# L1 — the developer MCP

The developer layer. Standardises how AI-assisted work on this repository is
done: investigate, verify, review, commit, and exercise a running Human.

Not to be confused with the configurator MCP (`server/system/agentic/mcp`),
which is product code that configures a running Human for its users. This one
configures nothing; it operates on the repository and on a local dev server.

## Why it is a separate binary

L1 tools shell out to `git`, `go`, `npx`, `prettier` and `.intent/intent.mjs`.
That capability must not exist inside the binary customers run, and a build tag
is a weaker guarantee than a separate module — a tag has to be correct in every
build path, whereas a module that nothing imports cannot be linked in by
accident.

L1 also has to work with the Human server down. Investigating, reviewing and
committing do not need an API, and making them depend on a booted server with a
database would be a self-inflicted wound.

`dev/mcp/` is its own Go module with `replace github.com/crusttech/human/server
=> ../../server`, following the precedent already set by
`extra/server-discovery/`.

## What it shares with L2

Everything structural, nothing domain. The registry, the toolkit, group/risk
tagging, progressive disclosure, scope, the structural test and the generated
`TOOLS.md` are factored into `server/pkg/mcpkit/` and imported by both servers.

**`mcpkit` may not import any Human domain package.** That is the constraint
that keeps the vendor tree of this module small and the boundary honest. If a
piece of the current MCP package resists the move because it knows about
Human's types, it belongs on the L2 side of the line.

## What is a tool and what is a skill

Skills own procedure; tools own primitives.

A skill is markdown, changes in seconds, and encodes how this team works —
`/intent-task`'s five phases are exactly the kind of thing that should stay
editable without a rebuild. A tool is the expensive primitive underneath: the
one that runs a test suite and returns three failures instead of four thousand
lines, or resolves which contracts govern a file without six greps.

The token saving comes from the primitives. The repeatability comes from the
skills. Neither replaces the other.

## Tools

Naming follows `RESOURCES.md`: `dev_<resource>_<verb>`, snake_case, singular.
All in group `development`.

### Investigate

| Tool | Risk | Does |
|---|---|---|
| `dev_intent_governing` | read | Which intent docs govern these files, with the locked-contract and WIP blocks verbatim |
| `dev_intent_affected` | read | Which e2e specs a change affects |
| `dev_convention_lookup` | read | Which conventions apply to these paths (formatting, codegen, i18n, frozen routes) |

### Verify

| Tool | Risk | Does |
|---|---|---|
| `dev_test_run` | read | Runs the touched package's suite, Go or vitest, and returns **only failures** with `file:line` |
| `dev_lint_run` | read | Lints changed files only |
| `dev_intent_check` | read | Drift check; returns the docs that need reconciling |
| `dev_format_run` | write | gofmt/prettier on changed files only — never a directory, which has caused reverted churn twice |

### Review

| Tool | Risk | Does |
|---|---|---|
| `dev_diff_survey` | read | The branch's changed surface: files, hunk summary, governing contracts, conventions in play |

### Commit

| Tool | Risk | Does |
|---|---|---|
| `dev_commit_create` | write | Stages and commits with the discipline enforced: imperative short subject, no AI trailers, formatted before staging, refuses a commit that spans unrelated logical changes |
| `dev_branch_status` | read | Branch, base, ahead/behind, working tree state |

The commit tool writes to history, which is worth being explicit about: it does
not change *who decides* to commit — that remains the human, every time — it
changes whether the conventions are followed by an agent remembering CLAUDE.md
or by code that refuses. Enforcement, not autonomy.

### Exercise a running Human

The purpose is narrow and worth stating exactly: **when working on a feature or
an issue, be able to test that feature.** That usually means standing up a
namespace and a configuration built for the thing under test — not a generic
demo — and then either running e2e specs against it or driving it by hand.

These tools drive the dev server and the frontend. They do **not** configure
Human: the L2 tools already create namespaces, modules, pages, records and
automation, and both servers are connected at once, so the configuration is
composed from L2 and the environment from L1. Duplicating the configurator here
would be two implementations of one thing, drifting.

| Tool | Risk | Does |
|---|---|---|
| `dev_server_status` | read | Is it up, which build, how stale |
| `dev_server_logs` | read | Filtered tail |
| `dev_fixture_seed` | write | Versioned fixtures from `dev/fixtures/` |
| `dev_fixture_cleanup` | destructive | Removes `agent-` prefixed data only |
| `dev_scratch_build` | write | Stands up a scratch environment for the feature under test — pages and charts via `pagebuild.py`, never envoy YAML, because block refs do not resolve that way |
| `dev_ui_verify` | read | Drives the frontend to check what a change actually looks like |
| `dev_e2e_run` | read | Runs named e2e specs, returns failures and trace paths |

## Rules that carry over from L2

- Declarations in `*_tools.go`, handlers in `*_handler.go`.
- `WithRisk` is the sole writer of the protocol annotation hints.
- Every ID-shaped parameter is a JSON string.
- Anything created on the dev server is `agent-` prefixed and disposable;
  unprefixed data is off-limits.
- The structural test and generated `TOOLS.md` cover this server too.

## Order of work

1. Factor `mcpkit` out of `server/system/agentic/mcp`, with the no-domain-imports
   rule enforced by a test. L2 keeps working throughout.
2. Stand up `dev/mcp` — module, vendor, boot, `.mcp.json` registration, one
   trivial tool end to end.
3. Verify family first: it is the largest token saving and pays off on the work
   already in flight.
4. Investigate, then review, then commit.
5. The dev-server/FE family last, since it is the one with a dependency on a
   booted server.
