# Frontend Safety Audit — 2026-04-17

Scope: all seven frontend webapps (`admin`, `agentic`, `compose`, `home`, `one`, `taq`, `workflow`) plus shared libs (`lib/vue`, `lib/js`). Dimensions: security (XSS/injection), runtime safety, dependencies, lint-debt shape. Out of scope: authz/RBAC logic, type-safety coverage, accessibility, performance.

## Executive summary

| Severity | Count |
|---|---|
| Critical | 3 |
| High | 5 |
| Medium (patterns) | 3 |
| Low / lint debt | ~430 across apps |

**Top 5 fixes this week:**

1. **Remove `eval()` from `evaluatePrefilter`** — `client/web/compose/src/lib/record-filter.js:309-313`. Any author who can edit a ContentBlock body or prefilter ships arbitrary JS to every viewer of that page.
2. **Stop rendering page-block content with `v-html`** — `client/web/compose/src/components/PageBlocks/Blocks/ContentBlock.vue:3`. Couples with #1 for full stored-XSS + RCE.
3. **Replace `new Function()` in chart metric/transform code** — `lib/js/src/compose/types/chart/base.ts:66`, `lib/js/src/compose/types/page-block/metric.ts:145`. Chart config is user-authored; this is RCE on every dashboard viewer.
4. **Upgrade `axios`, `lodash-es`, `picomatch`** — actionable direct-dep CVEs at high severity.
5. **Add security headers in `client/web/nginx.conf`** — currently no CSP, no X-Frame-Options, no X-Content-Type-Options, no Referrer-Policy. A CSP alone mitigates most of #1–#3 as defense-in-depth.

---

## Critical findings

### C1. `eval()` on user-authored template strings in `evaluatePrefilter`

**File:** `client/web/compose/src/lib/record-filter.js:309-313`

```js
export function evaluatePrefilter(prefilter, { record, user, recordID, ownerID, userID }) {
  return (function (prefilter) {
    return eval('`' + prefilter + '`')
  })(prefilter)
}
```

**Risk:** The `prefilter` argument is user content from page-block configuration (e.g. ContentBlock body, record-list prefilters). Anyone who can author/edit these blocks can smuggle arbitrary JS through template-literal `${…}` interpolation — `${alert(document.cookie)}` is a one-liner RCE. Since page configs are persisted server-side and then executed in every viewer's browser, this is **stored XSS + client-side RCE**. Token theft via DOM is trivial once JS runs; a CSP (see H2) is not in place.

**Fix:** Replace `eval` with a safe interpolator. Build a small parser that only substitutes from a whitelist of named placeholders (`recordID`, `ownerID`, `userID`, `record.values.X`, `user.name`, etc.):
```js
export function evaluatePrefilter(tpl, ctx) {
  return tpl.replace(/\$\{([a-zA-Z_][\w.]*)\}/g, (_, path) => {
    const v = path.split('.').reduce((o, k) => (o == null ? undefined : o[k]), ctx)
    return v == null ? '' : String(v)
  })
}
```

### C2. `v-html` on eval-produced output in ContentBlock

**File:** `client/web/compose/src/components/PageBlocks/Blocks/ContentBlock.vue:3`

```html
<div class="p-3" style="white-space: pre-wrap" v-html="contentBody" />
```
`contentBody` is `evaluatePrefilter(body, …)` (see C1). This is both the rendering sink for C1's RCE and an independent stored-XSS vector for any author who writes raw HTML into the block body.

**Fix:** Two layers: (1) migrate `evaluatePrefilter` per C1, and (2) decide whether Content blocks should render HTML at all. If they must (legacy), pipe through DOMPurify with an allow-list of safe tags; otherwise switch to `{{ contentBody }}`.

### C3. `new Function()` on user-authored chart formulas

**Files:**
- `lib/js/src/compose/types/chart/base.ts:66` — `const fx = new Function('n', 'm', 'r', fxRaw)` with `fxRaw = m.fx || defaultFx`
- `lib/js/src/compose/types/page-block/metric.ts:145` — `rtr = new Function('v', \`return ${m.transformFx}\`)(rtr)`

**Risk:** `m.fx` and `m.transformFx` come from chart/metric block configuration edited in the UI (Admin → Charts, metric blocks). Any chart author can write a formula that runs arbitrary code in every dashboard viewer's browser. Same blast radius as C1: stored XSS + client RCE across all users who view the chart.

**Fix:** Replace with a restricted expression evaluator (e.g. `expr-eval`, `mathjs` with a safe scope) that only exposes the intended numeric operations. Reject or sandbox existing configurations at load; add server-side validation when saving.

---

## High findings

### H1. `v-html` of raw record string values

**File:** `lib/vue/src/components/field/viewers/CFieldStringViewer.vue:6`

```html
<p v-if="formatted" :class="viewerClasses" v-html="formatted" />
```
`formatted` comes directly from `record.values[field.name]`. Every string field renders its value as HTML regardless of whether the field's `useRichTextEditor` option is set (that option only toggles CSS classes, not the sink). **Any user who can write a string field value can XSS any viewer of that record.**

**Fix:** Split the viewer: when `field.options?.useRichTextEditor` is true, sanitize via DOMPurify with a tight allow-list before `v-html`. Otherwise render with `{{ formatted }}` (preserving `multiline` via `white-space: pre-line`).

### H2. Missing security headers in nginx

**File:** `client/web/nginx.conf`

No `add_header` directives anywhere. Missing: `Content-Security-Policy`, `X-Frame-Options` (clickjacking), `X-Content-Type-Options: nosniff`, `Referrer-Policy`, `Strict-Transport-Security` (if TLS-terminated upstream). A restrictive CSP with `script-src 'self'` would neutralize most of C1-C3 and H1/H3 as defense-in-depth.

**Fix:** Add to the `server` block:
```
add_header Content-Security-Policy "default-src 'self'; script-src 'self'; style-src 'self' 'unsafe-inline'; img-src 'self' data: blob: https:; connect-src 'self' https:; frame-ancestors 'none'; base-uri 'self';" always;
add_header X-Content-Type-Options "nosniff" always;
add_header X-Frame-Options "DENY" always;
add_header Referrer-Policy "strict-origin-when-cross-origin" always;
```
Verify before deploying — `unsafe-inline` style may need to stay until PrimeVue inline styles are handled; `script-src` must not need `'unsafe-inline'` or `'unsafe-eval'` (another reason to fix C1/C3 first).

### H3. `v-html` of prompt messages from automation/workflow configs

**Files:**
- `lib/vue/src/components/prompts/kinds/CPromptAlert.vue:3`
- `lib/vue/src/components/prompts/kinds/CPromptInput.vue:3`
- `lib/vue/src/components/prompts/kinds/CPromptOptions.vue:3`
- `lib/vue/src/components/prompts/kinds/CPromptChoice.vue:3`
- `lib/vue/src/components/prompts/kinds/CPromptComposeRecordPicker.vue:3`
- `lib/vue/src/components/prompts/kinds/CPromptNotification.vue:2`

All render `message` (a workflow/automation-authored string) with `v-html`. Any workflow editor can stage XSS against the end users who receive the prompt.

**Fix:** Either sanitize (DOMPurify) or switch to `{{ message }}` with `white-space: pre-wrap` already in place. Prompts are user-facing notifications — they rarely need raw HTML.

### H4. Direct-dep CVEs (actionable)

From `pnpm audit --prod` (19 total advisories):

| Module | Severity | Advisory | Path |
|---|---|---|---|
| `axios` | high | DoS via `__proto__` in form data | root |
| `axios` | moderate | NO_PROXY hostname bypass (SSRF) | root |
| `axios` | moderate | Cloud metadata exfiltration via headers | root |
| `lodash-es` | high | Code injection via `_.template` imports | root |
| `lodash-es` | moderate | Prototype pollution via `_.unset` | root |
| `lodash-es` | moderate | Prototype pollution via array path bypass | root |
| `lodash` | high/moderate | Same lodash issues | `client/web/taq` (via `dagre`) |
| `picomatch` | high | ReDoS via extglob quantifiers | root |
| `picomatch` | moderate | POSIX character class method injection | root |

**Fix order:** `axios` and `lodash-es` are used in runtime code paths and are the priority. Run `pnpm up axios lodash-es picomatch` at the root and verify builds/tests. `lodash` under `dagre` is transitive — see M3.

### H5. `JSON.parse` on localStorage deserialization without try/catch

**Files (selected):**
- `client/web/compose/src/components/PageBlocks/Blocks/RecordListBlock.vue:1755,1770`
- `client/web/admin/src/views/system/Settings/Index.vue:579,794`
- `client/web/compose/src/components/ModuleFields/Configurator/index.vue:102`

localStorage can be corrupted, truncated, or manipulated. A bare `JSON.parse(stored)` crashes the component tree on malformed input. Not a security issue by itself (local attacker already has access), but the resulting "white screen of death" after a botched extension / dev-tools tweak is a real support cost.

**Fix pattern:** Wrap all localStorage reads in `try/catch` and return a defined default (`CPermissionGrid.vue:307` already does this — replicate).

---

## Medium patterns (aggregated)

### M1. Bare `JSON.parse` on clipboard / drag-drop payloads

**Examples:**
- `client/web/workflow/src/composables/useWorkflowClipboard.js:72` — parses clipboard text
- `client/web/workflow/src/composables/useWorkflowDnD.js:62` — parses `dataTransfer`
- `client/web/compose/src/components/Modules/ModuleImporter.vue:163` — parses imported file text
- `client/web/admin/src/components/Template/CTemplatePreview.vue:110`

**Count:** ~8 sites across the client.

**Fix pattern:** Wrap in `try/catch`, surface a user-visible "invalid clipboard content" toast rather than crashing the composable.

### M2. Unhandled promise chains

**Signal:** 74 `.then(` uses vs 50 `.catch(` uses across the client. That's 24 net `.then` without obvious rejection handling — many live in UI event handlers where a rejected fetch will currently silently swallow.

**Representative sites:**
- `client/web/compose/src/components/PageBlocks/Blocks/MetricBlock.vue` — one chain, no catch
- `client/web/admin/src/components/Permissions/CPermissionGrid.vue` — 3 chains, 1 catch
- `client/web/workflow/src/components/Export.vue` — 2 chains, 2 catches (OK)
- `client/web/compose/src/views/Namespace/Manage.vue` — 2 chains, 2 catches (OK)

**Fix pattern:** Each `.then` in an event handler or lifecycle hook should have a `.catch(e => toast(e))` or equivalent. Prefer `async/await` + `try/catch` for new code — easier to enforce.

### M3. Transitive-only dep vulnerabilities

From `pnpm audit --prod`, these have no direct action:

| Module | Severity | Source |
|---|---|---|
| `minimatch` (3 advisories, ReDoS) | high | `lib/js` via `glob` (dev/build path) |
| `@isaacs/brace-expansion` | high | `lib/js` |
| `follow-redirects` | moderate | via `axios` |
| `qs` | low | root transitive |

**Fix:** Monitor. `follow-redirects` will resolve when `axios` is upgraded (H4). `minimatch` / `brace-expansion` in `lib/js` dev deps have low runtime impact. `qs` is low severity and awaits upstream.

---

## Low / lint debt (summary)

| Workspace | Total | Errors | Warnings | Notes |
|---|---|---|---|---|
| `client/web/admin` | 34 | 23 | 11 | |
| `client/web/compose` | 218 | 156 | 62 | largest debt |
| `client/web/agentic` | 5 | 0 | 5 | cleanest |
| `client/web/taq` | 26 | 1 | 25 | |
| `client/web/workflow` | 22 | 2 | 20 | |
| `client/web/one` | 0 (quiet) | 0 | 0 | likely clean |
| `client/web/home` | — | — | — | **no `eslint.config.js` — lint never runs** |
| `lib/vue` | 92 | 4 | 88 | |
| `lib/js` | 30 | 7 | 23 | |
| **Total** | **~427** | **~193** | **~234** | |

**Top rule categories in `compose` (largest bucket):** `vue/no-mutating-props` (errors), `no-unused-vars` (warnings), `vue/multi-word-component-names`. None security-adjacent; this is style + correctness debt.

**Action item (meta-finding):** `client/web/home` has no eslint config — lint is silently a no-op there. Copy an existing `eslint.config.js` (e.g. from `one`) as a starting point.

---

## Dependency findings (summary)

- **Actionable (direct):** see H4. Run `pnpm up axios lodash-es picomatch` at root; verify `client/web/taq` → `dagre@latest` or override to remove `lodash@<4.18.1`.
- **Upstream-blocked / monitor:** see M3. No immediate action.
- **Outdated critical libs:** not surveyed in depth — run `pnpm outdated -r` as a follow-up if you want an update cadence view.

---

## Appendix

### Methodology

1. **Tooling sweep.** `pnpm audit --prod --json` at the monorepo root; `pnpm exec eslint src` per workspace; ripgrep patterns for known unsafe sinks and missing error handling.
2. **Hotspot manual review.** Content rendering (v-html sites), `eval` / `Function` sites, auth template (`server/auth/assets/templates/login.html.tpl`), nginx config, localStorage reads, store injection patterns.
3. **Triage.** Severity rubric: Critical = exploitable RCE/XSS reachable by normal users; High = likely-exploitable or high-severity CVE with fix; Medium = repeated patterns aggregated; Low = lint/style debt summarized only.

### Grep patterns used

- `v-html` — found 11 live sites; 3 are safe (`WorkflowNode.vue`, `TriggerNode.vue`, `CChatMessages.vue` via pre-escaped markdown), 8 are unsafe (see C2, H1, H3).
- `innerHTML\s*=` — 2 sites in `CEmojiPicker.vue` — safe, no user input.
- `\beval\s*\(|new Function\s*\(` — 4 sites; 1 critical (record-filter), 2 critical (chart/metric), 1 is build-time (vue.config-builder.js, not shipped).
- `href="javascript:` — none found.
- `localStorage\.(set|get)Item` — surveyed for untrusted deserialization (see H5).
- `.then(` vs `.catch(` counts — see M2.

### Tools used

- `pnpm audit --prod --json` (pnpm 9.x advisory DB, snapshot 2026-04-17)
- `pnpm exec eslint src` per workspace (eslint 9.39.2)
- `ripgrep` via Grep tool

### Not covered

- Authorization / RBAC logic in route guards and conditional rendering — deliberately deferred per scope.
- Type-safety (`any` usage, unsafe casts) — deferred.
- Accessibility, performance, i18n.
- Full `pnpm outdated` sweep.
- CI security posture (GitHub Actions permissions, secret handling in `release.yml`).
- Server-side routes; this audit is FE-only.

### Limitations

- Pattern matching misses logic-level flaws (e.g. SSRF by URL composition that doesn't hit obvious sinks).
- Hotspot review is not exhaustive; low-traffic surfaces may harbor untriaged issues.
- Advisory DB is a snapshot — re-run `pnpm audit` before shipping.
