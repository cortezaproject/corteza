# Frontend gotchas

Facts about the unify webapp, lib/vue, PrimeVue, Vue 3, Tailwind, CodeMirror and i18n that the code does not make obvious.

## Vue

### An immediate watcher runs before later consts exist

In `<script setup>`, `watch(..., { immediate: true })` fires at the line where it is written, so a `const`/`let` declared further down is still in its temporal dead zone and throws `ReferenceError`. A surrounding `try/catch` meant for one environmental failure (such as no localStorage) swallows it silently, and unit tests stay green.

**Why:** a hoisted `function` exists from the top of setup, a `const … = computed(...)` only from its line; `RecordListBlock.vue` loads its stored filter from such an immediate watcher.

**How to apply:** helpers reached from an immediate watcher are `function` declarations or sit above the watcher; when something inside a tolerant `try/catch` stops working, suspect the catch first and verify persistence in a browser.

### A record picker clears itself when its namespace blinks

`CInputRecord` (`lib/vue/src/components/input/CInputRecord.vue`) emits `null` whenever its namespace or module ID changes from one value to another, `undefined` included. A parent that resets the namespace to `{}` and refills it after an async lookup therefore wipes the selection on every re-run, and the option the user then picks never reaches `onSelect`.

**How to apply:** resolve first and assign field and namespace together, as `ConditionRow.vue` does; guarded by "never shows the editor an empty namespace while re-resolving".

### v-for'd inline elements render with no space between them

Vue's default `whitespace: 'condense'` removes a whitespace-only node containing a newline between two elements, but condenses it to one space next to a text node. So a `v-for` over inline elements renders as one run-on token (`${recordID}${ownerID}`) while the surrounding prose keeps its spaces. Adding `{{ ' ' }}` then yields two spaces (harmless once rendered).

**How to apply:** assert on `textContent.replace(/\s+/g, ' ')`; prefer a real space over CSS `gap` for snippets a user copies, as `CExpressionHint.vue` does.

### Kebab attributes camelize to lowerCamel

`:role-id` becomes `roleId`, never `roleID`. `CInputUser` declares only `roleID`, so `:role-id` matches no prop and falls through as an attribute, silently dropping the filter.

**How to apply:** bind props with an acronym tail in camelCase (`:roleID="roles"`); Vue matches the raw attribute name.

### Fallthrough attrs override the root component's bound prop

An attribute a parent passes that a component does not declare falls through to its root; if that root is a component, it arrives as a prop and overrides what the template binds explicitly. Vue 2 only set such attrs on the DOM, so a Corteza port can turn a dead attr live (a `supported-metrics` passed to `RadarChart.vue` would beat its own `:supported-metrics="-1"` on `ReportEdit`).

**How to apply:** when porting a Vue 2 wrapper whose root is a component, check what the parent passes that the wrapper does not declare.

### Untyped-prop wrappers turn bare booleans into ""

`VueDraggable` (vue-draggable-plus) declares `props: [names]` with no types, so a bare `force-fallback` reaches SortableJS as `""`, which is falsy. To see what SortableJS received, read `el[Object.keys(el).find(k => k.startsWith('Sortable'))].options` in the browser.

**How to apply:** bind booleans explicitly (`:force-fallback="true"`) on any untyped-prop wrapper.

### A line-wrapped :deep() selector matches every descendant

When prettier breaks `.a :deep(\n long > selector\n)::before` over lines, the scoped-CSS compiler keeps the whitespace before `::before`, emitting a descendant `::before` selector. `content: ''` there wipes every PrimeIcons glyph in the row.

**How to apply:** keep `:deep()` selectors short enough for one line (put a short class on the element via `pt`); if icons vanish, read the compiled selector from `document.styleSheets`.

### A label transform breaks two-part labels

Bulk-converting `<label>{{ $t('k') }}</label>` to `<CFormGroup :label="$t('k')">` turns `{{ $t('k') }} ({{ n }})` into `$t('k')(n)`, a call on a string that no test or typecheck catches.

**How to apply:** after such a sweep, read every `^+.*:label=` line in the diff; two-part labels need CFormGroup's `#label` slot. `curl -o /dev/null -w '%{http_code}' http://localhost:5173/@fs/<abs-path>.vue` returns 500 on a template compile error.

### Pinia $onAction misses actions called inside the store

In a setup store, `$onAction` fires only for calls made through the store instance. One action calling another directly (`handleRealtime` → `addNotification` in `useNotificationsStore`) is a raw function call and is never seen.

**How to apply:** listen for the outer action and inspect its `args` (`name === 'handleRealtime' && args[0]['@type'] === 'notification'`).

## PrimeVue

### A shared popup TieredMenu never re-anchors

PrimeVue 4.5 `TieredMenu.show()` keeps the first `target` it saw until unmount, so one popup shared across a list anchors to the first trigger forever, and to (0,0) once that element detaches or is hidden by `v-show`. `Menu` clears `target` on hide and takes `show(event, target)`; `Popover` and `ContextMenu` are also safe. `Menu` cannot do tiers: a nested `items` array renders as a header with inline children, and a third level is silently dropped.

**How to apply:** never share one popup TieredMenu across triggers; use `Menu` with `show(event, event.currentTarget)`, and if per-row tiers are needed, re-anchor by hand and test it like `CResourceList.anchor.test.ts`.

### Overlays close on any scroll of the trigger's ancestors

PrimeVue popups bind a scroll handler to every scrollable ancestor of their target and `hide()` on any scroll. A trigger partly outside a horizontally scrollable container gets scrolled into view on click, which closes the menu that just opened, so the button looks dead. It depends on viewport width (`CFormList` is `min-w-max` and overflows below ~1600px).

**How to apply:** when a popup will not open, patch `hide()` to log a stack or watch `document` scroll events in capture mode; fix the container so it does not scroll.

### Select overlays size to their widest option

`Select`/`MultiSelect` overlays get an inline `min-width: <trigger>px` and `width: auto`, so one long option makes the panel far wider than its control. Option content cannot bound it (the `li` is flex; `w-0 min-w-full` and `min-w-full max-w-0` variants all fail), and `li.p-select-option` is `white-space: nowrap`, which inherits.

**How to apply:** cap the overlay with `:pt="{ overlay: { class: 'max-w-lg' } }"` (inline min-width wins over it, so wide triggers keep their width) and put `w-full min-w-0 whitespace-normal break-words` on the option wrapper; `CInputPickers.descriptions.test.ts` asserts this for the `CInput*` pickers.

### A filtered Select shows the placeholder

PrimeVue 4.5 `Select` computes its label from `visibleOptions`, the filtered list, so while the filter hides the current option the closed box shows the placeholder, and it stays blank without `reset-filter-on-hide`. Typing in the filter resets `focusedOptionIndex` to -1, so Enter right after typing closes without selecting (ArrowDown then Enter works). Without `auto-filter-focus`, each keystroke on the closed Select replaces the filter instead of appending.

**How to apply:** render the label through a `#value="{ value, placeholder }"` slot (`value?.name || placeholder`), as the sidebar namespace switcher does; test via `[data-pc-section="label"]`, since tests mount PrimeVue with `unstyled: true`, which leaves no `p-*` classes.

### An empty-string option is never selected inside a Form

BaseEditableHolder's immediate `$formDefaultValue` watcher treats `''` as empty, so inside a `<Form>` it overwrites `d_value` with `undefined`, and `modelValue` never changes to restore it. A SelectButton option with `value: ''` never shows as selected.

**How to apply:** inside a Form, map a non-empty stand-in value to the stored `''` through a computed.

### A manual :invalid inside a resolver Form is a second rule

`<Form :resolver>` + `CFormGroup name="x"` already drives the input's red state and renders the message. A hand-bound `:invalid` on the input is an independent rule nothing checks at submit, since `handleSubmit({ valid })` sees only the resolver: red border, no message, no effect on save.

**Why:** the visual state and the save gate need one source of truth, and the resolver decides.

**How to apply:** put the rule in the resolver and bind no `:invalid` (reference: `sections/admin/views/system/UserGroup/Editor.vue`); outside a resolver Form, gate a manual `:invalid` on a `submitted` ref, as `project/components/datamodel/ModuleCreateDialog.vue` does.

### DataTable header pt and multi-sort

In PrimeVue 4.5, table-level `pt: { headerCell }` is a no-op: header cells merge `ptm('column.headerCell')`. Use `pt: { column: { headerCell } }` for every column, or each Column's own `:pt="{ headerCell }"` (classes merge; `CResourceList` passes `headerCellClass` to each). `multiSortMeta` is aliased: ctrl-click mutates the bound array in place before `@sort`, and in multiple mode the event carries no clicked column. PrimeVue's loading mask fades over ~150–300ms, so a raw `:loading` washes even a 36ms request; `useTableBusy` defers it.

**How to apply:** read the clicked header from `event.originalEvent.target.closest('th')` with a `data-field` set through column pt; check computed `backgroundColor` per `th` rather than trusting a screenshot.

### Record-list cells fill in after the first render

Every field viewer in `lib/vue/src/components/field/registry.ts` is `defineAsyncComponent`, so a record list's first render has empty cells and a post-flush layout measurement sees header-only widths; `useGrowOnlyColumns` also measures in a `flush: 'pre'` watcher. `resolveRecordLabels` shares in-flight IDs, so overlapping fetches finish in order; to test `fetchSeq` in RecordListBlock, race a slow-label sort against a search that returns nothing. `useTheme.ts` makes resizable headers clip their content, so the resize grip straddles the boundary; frozen headers sit above it at z-index 2, which only matters once the table scrolls past its first column (test at 700px).

### Tree drop zones are the row's height bands

PrimeVue 4 Tree binds dragover on `.p-tree-node-content` only: the top 25% drops before, the bottom 25% after, the middle inside. Space between rows is no target, so widening gaps makes dropping harder, and a `w-fit` row leaves its right side dead. The marker is an empty `.p-tree-node-drop-point` sibling with a 1px outline.

**How to apply:** make rows taller, not gaps wider.

## Drag and drop

### CDraggableList is the one flat-list drag primitive

`lib/vue/src/components/drag/CDraggableList.vue` wraps vue-draggable-plus (SortableJS) with fallback mode, placeholder, autoscroll and touch. Items are found by `itemSelector` (default `[data-drag-item]`), and a cross-list move arrives as `fromKey`/`toKey` from `dragKey`. vue-draggable-plus clones a crossing item with `JSON.parse(JSON.stringify(item))`, which honours `compose.Record.values.toJSON` (`lib/js/src/compose/types/record.ts`) and lands a blank card, so `clone` defaults to identity. Never wrap the `v-for` in a `v-if` on length (dragging the last item out leaves a ghost node), and never put `a` in `filter` (CAppList rows are anchors). Trees are not its job.

**How to apply:** override `clone` only for lists that pull copies; give `<a>`/`<img>` items `draggable="false"`, as CAppList does.

### SortableJS lessons for flat lists

- Drop areas that appear at drag start must take no room; in-flow slots shift content under a still pointer.
- SortableJS re-evaluates only when the pointer moves, so a dwell rule never completes on a still pointer; in fallback mode it also polls every 50ms at the last position, so content can slide under it.
- A refused `onMove` bubbles to the parent list with the ancestor as target; use `document.elementFromPoint(x, y).closest(...)`.
- Wrapping `options.onMove` before a drag is undone at drag start; wrap after it starts.
- A prop-backed `defineModel` lags the DOM: read the reactive data after `nextTick` in a drop handler.
- Playwright `mouse.move(..., { steps })` fires in one tick; pace moves (8px, 25ms).
- Fallback mode needs `body.c-dragging { user-select: none }` and `.c-drag-image { list-style: none }`; CDraggableList has both.

### A drag fires click on the common ancestor

A `click` is dispatched on the nearest common ancestor of the mousedown and mouseup targets, and PrimeVue `Slider` binds release to `document` mouseup, so dragging a slider and releasing off it clicks the wrapper. `@click.stop` on the Slider is never in the path. `CInputToggleCard` records the `pointerdown` target and toggles only when the click target matches.

**How to apply:** reproduce in jsdom by dispatching `pointerdown` on the inner control, then a bubbling `click` on the wrapper.

## Expressions, CodeMirror and dates

### CInputExpression and its six scopes

`lib/vue/src/components/expression/` is the one input for expressions, registered globally by `PrimeVueComponentsPlugin`. `catalog.ts` builds scopes and `resolvePath` returns found, unknown or unverifiable (the third stops lint on an unloaded module); `syntax.ts` holds pure lint/completion functions; `EXPR_FUNCTIONS` is generated into `functions.gen.ts` from `server/pkg/expr/expr_functions.yaml`. `useExpressionScope.ts` states the two-module rule: `${record...}` resolves against the page's module, bare QL identifiers against the queried one. Each call site has its own scope: templates (`record`, `user`, `recordID`, `ownerID`, `userID`), visibility (`record`, `user`, `screen`, `isView`/`isCreate`/`isEdit`), sanitizer (`value`, `values/sanitizer.go:92`), validator (`value`, `oldValue`, `values.<field>`, `values/validator.go:232`), value expression (bare field names, `new`, `old`, `values/expr.go`).

**How to apply:** an unknown root or function in `expr` is an error (visibility then reads every expression false, hiding the block); an unknown member resolves to null, so it is a warning.

### CodeMirror inside a PrimeVue dialog

CodeMirror's keymap calls `preventDefault` but not `stopPropagation`, so Escape on the completion popup also closes the dialog. The fix is a `Prec.highest(EditorView.domEventHandlers({ keydown }))` that stops propagation when `completionStatus` is active; at lower precedence `completionKeymap` closes first and the guard silently declines. CodeMirror filters options against the text from `from` to `pos`, so a range starting at punctuation (`$` → `${recordID}`) discards every option.

**How to apply:** set `filter: false` and sort options yourself; jsdom with `EditorView.findFromDOM` + `startCompletion`/`currentCompletions` reproduces both.

### Expression box height beside small fields

`CInputExpression` `size="small"` with `:min-lines="1"` reads `--p-form-field-sm-*` tokens (there is no `--p-form-field-font-size`) and measures 32.94px against 33.25px for a small Select; the gap comes from Poppins' `normal` line-height and is not worth chasing. Inside an `InputGroup`, the addon (no small variant, 40px) sets the height; `class="px-2 py-0 text-sm"` on it brings the group to 32.94px.

**How to apply:** change height through `.cm-content`'s `minHeight`, never the `.c-expression` wrapper.

### Dates come from Intl

Date format, day/month names and week start come from `Intl` (`lib/vue/src/plugins/primevue-locale.ts`). PrimeVue's tokens are jQuery-UI's (`dd`/`mm`, `yy` is the four-digit year), and its `firstDayOfWeek` counts Sunday..Saturday as 0..6 while Intl counts Monday..Sunday as 1..7. Week data has two spellings, `getWeekInfo()` and `weekInfo`, and Node 22 has only `weekInfo`. Week start must use the raw locale tag, since `resolvedOptions().locale` collapses an unshipped language to en-US. The date locale is the browser's, not the UI language (which falls back to `en`).

**How to apply:** probe both week-data spellings; never round-trip the tag through `DateTimeFormat`.

## Tailwind and theming

### dark: with ! loses to the light class

Tailwind utilities sit in `@layer tailwind-utilities` (`client/web/unify/src/assets/styles.css`), but variant rules like `dark:` are emitted unlayered. Normal declarations favour unlayered rules, so `bg-x dark:bg-y` works; `!important` flips that, so `!bg-x dark:!bg-y` shows the light colour in dark mode too.

**How to apply:** skip `!` when recolouring PrimeVue (layer order `tailwind-base, primevue, tailwind-utilities` already wins); otherwise use `!` on neither class, or put both in scoped CSS.

### Tailwind turns generic class names into rules

The JIT scans `./src/**/*.{vue,js,ts,jsx,tsx}` and `lib/vue/src` (`tailwind.config.shared.js`), so a literal semantic class that matches a utility becomes a real rule: Corteza's `block` emitted `.block { display: block }` and broke every PrimeVue Card. Only literals count; dynamic PascalCase kind classes are safe. The rule lives in an injected layer, so searching `document.styleSheets` for it finds nothing.

**How to apply:** check hook class names against Tailwind's utilities (the compose block root is `page-block`); assert on `getComputedStyle`, not the rule.

### Theme tokens are not Tailwind's palette

`var(--p-green-500)` computes to `rgb(67, 170, 139)` here, not Tailwind's `rgb(34, 197, 94)`. A browser check that asserts a colour from memory reports failure on a correct page. The root font is 15px, so rem values are 15px-based.

**How to apply:** read the computed value once with the feature working, paste it into the check, and name its token in a comment.

## Layout and text

### A table in a flex column stretches

A `<table>` inside `flex flex-col` is a flex item, so `align-items: stretch` sizes it to the container and `w-full`, `w-auto` or no class render identically, pushing numbers far from their labels.

**How to apply:** add `self-start` so the table sizes to its content.

### #topbar-title is the page heading hook

`CTopbar.vue` renders `<div id="topbar-title">` and views `<Teleport>` their heading into it, so it answers "which page is this" without per-route plumbing; `client/web/unify/src/utils/documentTitle.js` mirrors it into `document.title`, and chrome beside the heading opts out with `data-title-exclude`. Teleport plants empty text nodes as anchors, so the first text node is `""`.

**How to apply:** walk to the first text node with non-whitespace content.

## i18n

### The filename is the i18n namespace

Each `locale/en/human-webapp/<name>.yaml` mounts under `<name>`, so its content starts at root keys: `general.yaml` has a root `label:` read as `$t('general.label.cancel')`. The FE loads it with one `localeGet({ lang, application: 'human-webapp' })` (`client/web/unify/src/plugins/index.js`). vue-i18n uses `a | b` plurals with a numeric choice, not `_one`/`_other` suffixes; `numberOfResults` stays non-plural because its call site passes a named `{count}`. Edits are live in development: the server overlays `LOCALE_PATH` and refreshes per request, so re-query `GET /system/locale/en/human-webapp` instead of rebuilding.

**How to apply:** never wrap a file in its own name; parse locale YAML with a YAML 1.2 loader, because PyYAML's 1.1 default turns `yes/no/on/off` into booleans.

### Unused locale keys mark unported features

`locale/en/human-webapp/` was backfilled from Corteza, so it carries strings for features whose code never landed. A key nothing in `client/web/unify/src` or `lib/vue/src` references signals a half-ported feature and gives its intended wording. The YAMLs nest, so the dotted path cannot be read off a grep line, and a wrong path renders the raw key with no error.

**How to apply:** grep the locale for the feature, grep the code for the key, and resolve the full path with a YAML parser.

## Webapp structure

### Section migration hidden issues

Classes of bugs that compile and lint clean when a section moves into unify: its `public/` assets must be copied into `unify/public/` (workflow icons load via `getIcon()` from `${BASE_URL}icons/`); section-specific global CSS must merge into `unify/src/assets/styles.css`; literal path navigation (`push('/')`, ``push(`/builder/...`)``, `to="/..."`) breaks because paths are section-prefixed and `/` is home; prompt kinds filter on `meta.webapps` in `lib/vue/src/components/prompts/definitions.ts`, so the shell passes `route.meta.section` as the current webapp.

**How to apply:** grep each section for path literals and convert to named routes or prefixed paths.

### Section sidebars mount lazily

`CSidebar` is a PrimeVue `Drawer`, which renders its body only when opened, so a collapsed section sidebar never mounts. Data loaded in a sidebar's `onMounted` never arrives for a view that needs it (the TAQ builder's catalog loads in `Builder.vue`'s `onMounted` for this reason).

**How to apply:** data a view needs loads in the view or a non-lazy host; only sidebar-only nav data belongs in the sidebar.

### Two RBAC mechanisms

`useRBACStore.can('system/', 'project.create')` loads `permissionsEffective({})` without a resource, so it answers only component ops (`*.create`, `*.search`, `grant`), deny by default. Rows and editors gate on per-resource `can*` flags in the REST payload. A hand-written lib/js model that omits a flag the server ships makes a view wrapping rows in it read `undefined`, hiding the control from everyone, admins included.

**How to apply:** before gating, confirm the server ships the flag and that any wrapping model declares and `Apply()`s it; verify a permitted user still sees the control.

### Permission controls sweep

`client/web/unify/src/sections/permission-controls.test.js` asserts that every create control resolves to a `<x>.create` op or `canCreate*`, every submit in an edit/create view to `canUpdate*`/`canManage*`/`canSave`/`canEdit`, and every `<CPermissionsButton>` to `canGrant`. `sections/control-gates.js` resolves the control's actual condition (own `v-if`/`:disabled`, open ancestor `v-if`s, enclosing `if` blocks for menu entries, named computeds up to three hops) and is shared with `restore-controls.test.js`. The `project/` section is skipped: it gates on the project role preset.

**How to apply:** adding a create, save or permissions control means adding its gate, or the test names the file.

### The restore control

Restoring a soft-deleted resource uses label `general.label.restore` ("Restore", never "Undelete" where a user reads it), icon `pi pi-replay` (`pi-undo` and `pi-refresh` already mean other things), a standalone `<Button severity="warn">` with its own `restoring` ref or a menu entry with label and icon only, gated on deleted state and a permission, in delete's slot. An editor's restore handler re-runs its `load<X>()`, because reassigning the model makes `useDraftGuard` read the cleared `deletedAt` as an unsaved edit.

**How to apply:** `client/web/unify/src/sections/restore-controls.test.js` enforces this; `dev/agent/checks/restore-controls.mjs` covers it in a browser.

### lib/js clients require fields REST does not

The generated `lib/js/src/api-clients/*.ts` guard required fields more strictly than the server: `namespaceCreate` and `moduleCreate` throw `field meta is empty` without `meta` (`{}` passes), `userCreate` needs `userGroupID` (`'0'` passes), `workflowCreate` needs `runAs` and `ownedBy` (`'0'` passes). A payload proven with `api.sh` can still throw through the browser's `$ComposeAPI`.

**How to apply:** read the `if (!x) { throw }` block at the top of the generated method instead of guessing.

### OS notifications in a test browser

`lib/vue/src/composables/useSystemNotifications.ts` cannot show a real notification in headless Chromium, and CDP `Browser.setPermission` does not model the permission. Replace `window.Notification` with a recorder class via `addInitScript` and override `document.hasFocus`. Notifications are scoped to their recipient (another user's reads as "notification not found" beside HTTP 200), the human's browser runs as their own account rather than the agent user, and permission is per origin including port. The e2e login snapshot captures localStorage, including `notificationsFocusedTab`, so per-tab state leaks into every spec.

**How to apply:** look the human's account up in the checkout's DB before sending them a test notification; only the human can confirm a real OS notification.

### A reactive object posted to a frame throws DataCloneError

`postMessage` structured-clones what it sends, and a Vue reactive proxy cannot be cloned: `Failed to execute 'postMessage' … could not be cloned`. It shows only where the value is reactive. A compose block's options are a reactive draft in the page builder and plain on the viewed page, so the same block works on the page and fails in the builder.

**How to apply:** hand a frame a plain copy (`JSON.parse(JSON.stringify(v))` or `toRaw` deep), and check frame-backed blocks in the builder as well as on the page.
