# Field Type Audit: Human vs Corteza

> Comprehensive comparison of all 11 field types across editors, viewers, and configurators.
> Only **Human** files should be changed. `lib/js` field type definitions are the source of truth for options.

---

## Summary Matrix

| Field Type | Editor | Viewer | Configurator | Status |
|-----------|--------|--------|-------------|--------|
| **String** | ⚠️ Missing rich text editor | ✅ OK | ✅ OK | Needs work |
| **Number** | ⚠️ Missing progress bar mode | ⚠️ Missing progress bar rendering | ✅ OK | Needs work |
| **Bool** | ✅ OK | ✅ OK | ✅ OK | Good |
| **DateTime** | ✅ OK | ✅ OK | ✅ OK | Good |
| **Select** | ✅ OK | ✅ OK | ✅ OK | Good |
| **Email** | ⚠️ Falls back to String (no `type="email"`) | ✅ OK | ⚠️ Minimal (only multiDelimiter) | Minor |
| **Url** | ⚠️ Falls back to String (no URL processing) | ✅ OK | ⚠️ Missing trim/secure options in configurator | Needs work |
| **User** | 🔴 Minimal (no roles, selectType, search) | ✅ OK | ✅ OK | Needs work |
| **Record** | ⚠️ Missing queryFields, prefilter, selectType | ✅ OK | ✅ OK | Needs work |
| **File** | ⚠️ Missing allowImages/allowDocuments filter | ⚠️ Partial | ✅ OK | Needs work |
| **Geometry** | 🔴 Entirely missing | 🔴 Entirely missing | ⚠️ Configurator exists but no viewer/editor | Needs work |

---

## Detailed Findings Per Field Type

### 1. String

**lib/js options**: `multiLine`, `useRichTextEditor`, `multiDelimiter`

#### Editor ([CFieldStringEditor.vue](file:///home/fajfa/Human/human/lib/vue/src/components/field/editors/CFieldStringEditor.vue))
- ✅ `multiLine` → renders `Textarea` vs `InputText`
- 🔴 **`useRichTextEditor` is NOT implemented** — the configurator lets you enable it, but the editor always renders plain `InputText`/`Textarea`
- Corteza uses `CRichTextInput` (Tiptap-based) when `useRichTextEditor` is true

#### Viewer ([CFieldStringViewer.vue](file:///home/fajfa/Human/human/lib/vue/src/components/field/viewers/CFieldStringViewer.vue))
- ✅ Appears functional

#### Configurator ([String.vue](file:///home/fajfa/Human/human/client/web/compose/src/components/ModuleFields/Configurator/kinds/String.vue))
- ✅ `multiLine` checkbox
- ✅ `useRichTextEditor` checkbox
- ✅ Multi-delimiter component

> [!WARNING]
> The configurator lets users enable `useRichTextEditor`, but the editor ignores it. Users will configure it and see no effect.

---

### 2. Number

**lib/js options**: `presetFormat`, `format`, `prefix`, `suffix`, `precision`, `multiDelimiter`, `display` (number/progress), `min`, `max`, `step`, `showValue`, `showRelative`, `showProgress`, `animated`, `variant`, `thresholds`

#### Editor ([CFieldNumberEditor.vue](file:///home/fajfa/Human/human/lib/vue/src/components/field/editors/CFieldNumberEditor.vue))
- ✅ `precision`, `min`, `max`, `step`, `prefix`, `suffix` → passed to `InputNumber`
- ⚠️ Always shows `InputNumber` regardless of `display` option — no progress bar editing

#### Viewer ([CFieldNumberViewer.vue](file:///home/fajfa/Human/human/lib/vue/src/components/field/viewers/CFieldNumberViewer.vue))
- need to verify — likely missing progress bar rendering

#### Configurator ([Number.vue](file:///home/fajfa/Human/human/client/web/compose/src/components/ModuleFields/Configurator/kinds/Number.vue))
- ✅ Display type toggle (number vs progress)
- ✅ Number: format, prefix, suffix, precision
- ✅ Progress: min, max, showValue, animated
- ⚠️ Missing from Corteza: `showRelative`, `showProgress`, `variant`, `thresholds`

> [!WARNING]
> Configurator lets users choose progress display mode, but the editor and viewer may not render a progress bar.

---

### 3. Bool ✅

**lib/js options**: `trueLabel`, `falseLabel`, `switch`

#### Editor ([CFieldBoolEditor.vue](file:///home/fajfa/Human/human/lib/vue/src/components/field/editors/CFieldBoolEditor.vue))
- ✅ `switch` → `ToggleSwitch` vs `Checkbox`
- ✅ `trueLabel`, `falseLabel` with fallback to translated "Yes"/"No"

#### Viewer + Configurator — ✅ All good

---

### 4. DateTime ✅

**lib/js options**: `format`, `onlyDate`, `onlyTime`, `onlyPastValues`, `onlyFutureValues`, `outputRelative`, `multiDelimiter`

#### Editor ([CFieldDateTimeEditor.vue](file:///home/fajfa/Human/human/lib/vue/src/components/field/editors/CFieldDateTimeEditor.vue))
- ✅ `onlyDate`, `onlyTime`, `onlyFutureValues`, `onlyPastValues`

#### Viewer + Configurator — ✅ All good (uses `formatValue` from lib/js)

---

### 5. Select ✅

**lib/js options**: `options[]` (value, text, style {textColor, backgroundColor}), `selectType`, `displayType`, `multiDelimiter`, `isUniqueMultiValue`

#### Editor, Viewer, Configurator — ✅ All options fully implemented

---

### 6. Email

**lib/js options**: `outputPlain`, `multiDelimiter`

#### Editor (falls back to [CFieldStringEditor.vue](file:///home/fajfa/Human/human/lib/vue/src/components/field/editors/CFieldStringEditor.vue))
- ⚠️ Uses plain `InputText` — no `type="email"` HTML attribute
- Corteza uses `<b-form-input type="email" />`

#### Viewer ([CFieldEmailViewer.vue](file:///home/fajfa/Human/human/lib/vue/src/components/field/viewers/CFieldEmailViewer.vue))
- ✅ Appears to handle `outputPlain` option

#### Configurator ([Email.vue](file:///home/fajfa/Human/human/client/web/compose/src/components/ModuleFields/Configurator/kinds/Email.vue))
- need to verify `outputPlain` checkbox

---

### 7. Url

**lib/js options**: `trimFragment`, `trimQuery`, `trimPath`, `onlySecure`, `outputPlain`, `multiDelimiter`

#### Editor (falls back to [CFieldStringEditor.vue](file:///home/fajfa/Human/human/lib/vue/src/components/field/editors/CFieldStringEditor.vue))
- 🔴 **No URL processing at all** — Corteza applies `trimFragment`, `trimQuery`, `trimPath`, `onlySecure` via a `fixUrl` formatter on the input
- No `type="url"` HTML attribute
- No placeholder hint

#### Viewer ([CFieldUrlViewer.vue](file:///home/fajfa/Human/human/lib/vue/src/components/field/viewers/CFieldUrlViewer.vue))
- ✅ Appears functional

#### Configurator ([Url.vue](file:///home/fajfa/Human/human/client/web/compose/src/components/ModuleFields/Configurator/kinds/Url.vue))
- need to verify all trim/secure options are configurable

> [!IMPORTANT]
> The Url editor needs its own component that applies URL processing options (`trimFragment`, `trimQuery`, `trimPath`, `onlySecure`). The url utility functions already exist in [url.ts](file:///home/fajfa/Human/human/lib/vue/src/components/field/url.ts).

---

### 8. User

**lib/js options**: `roles[]`, `presetWithAuthenticated`, `selectType`, `multiDelimiter`, `isUniqueMultiValue`

#### Editor ([CFieldUserEditor.vue](file:///home/fajfa/Human/human/lib/vue/src/components/field/editors/CFieldUserEditor.vue))
- 🔴 **Only 29 lines** — just a thin wrapper around `CInputUser`
- Missing features vs Corteza (390-line editor):
  - **No `roles` filter** — users should be filtered by role when specified
  - **No `selectType` modes** — Corteza supports `default`, `multiple`, `each` with different UI
  - **No `presetWithAuthenticated`** — Corteza auto-fills current user on new records
  - **No `isUniqueMultiValue`** — Corteza prevents duplicates in multi-select
  - **No pagination** — Corteza has prev/next page for user search results
  - **No debounced search** — Corteza has debounced user search

#### Configurator ([User.vue](file:///home/fajfa/Human/human/client/web/compose/src/components/ModuleFields/Configurator/kinds/User.vue))
- ✅ Full: `presetWithAuthenticated`, `roles`, `selectType`, `isUniqueMultiValue`
- All options are configurable but the editor doesn't use them

> [!CAUTION]
> This is the biggest gap. The configurator is well-built with all options, but the editor is minimal and ignores almost all of them.

---

### 9. Record

**lib/js options**: `moduleID`, `labelField`, `recordLabelField`, `queryFields[]`, `selectType`, `multiDelimiter`, `isUniqueMultiValue`, `prefilter`

#### Editor ([CFieldRecordEditor.vue](file:///home/fajfa/Human/human/lib/vue/src/components/field/editors/CFieldRecordEditor.vue))
- ✅ `moduleID`, `labelField` — passed to `CInputRecord`
- ⚠️ **`queryFields` not passed** — search should query these fields
- ⚠️ **`prefilter` not passed** — records should be filtered
- ⚠️ **`selectType` not handled** — no `default`/`multiple`/`each` modes
- ⚠️ **`isUniqueMultiValue` not handled**
- Missing: pagination for record search results

#### Configurator ([Record.vue](file:///home/fajfa/Human/human/client/web/compose/src/components/ModuleFields/Configurator/kinds/Record.vue))
- ✅ Full: `moduleID`, `labelField`, `queryFields`, `prefilter`, `selectType`, `isUniqueMultiValue`

> [!WARNING]
> Same pattern as User: configurator is comprehensive, but editor is a thin wrapper that ignores the configured options.

---

### 10. File

**lib/js options**: `allowImages`, `allowDocuments`, `maxSize`, `mode` (list/gallery), `inline`, `hideFileName`, `mimetypes`, styling props, `clickToView`, `enableDownload`, `multiDelimiter`, `enableWebcam`

#### Editor ([CFieldFileEditor.vue](file:///home/fajfa/Human/human/lib/vue/src/components/field/editors/CFieldFileEditor.vue))
- ✅ `mimetypes` → accept filter
- ✅ `maxSize`
- ✅ `isMulti` → multiple files
- ⚠️ **`allowImages`/`allowDocuments` not used directly** — should construct mimetype filter
- ⚠️ **`enableWebcam` not used**

#### Viewer ([CFieldFileViewer.vue](file:///home/fajfa/Human/human/lib/vue/src/components/field/viewers/CFieldFileViewer.vue))
- Need to verify: `mode` (list vs gallery), `inline`, `hideFileName`, styling, `clickToView`, `enableDownload`

#### Configurator ([File.vue](file:///home/fajfa/Human/human/client/web/compose/src/components/ModuleFields/Configurator/kinds/File.vue))
- Need to verify against Corteza

---

### 11. Geometry 🔴

**lib/js options**: `center[]`, `zoom`, `multiDelimiter`, `prefillWithCurrentLocation`, `hideCurrentLocationButton`, `hideGeoSearch`

#### Editor — Missing entirely
- No `Geometry` entry in [registry.ts](file:///home/fajfa/Human/human/lib/vue/src/components/field/registry.ts) — falls back to String editor
- Corteza has a full map-based editor (10.5 KB) using Leaflet

#### Viewer — Missing entirely
- No `CFieldGeometryViewer.vue`
- Corteza has a map-based viewer

#### Configurator ([Geometry.vue](file:///home/fajfa/Human/human/client/web/compose/src/components/ModuleFields/Configurator/kinds/Geometry.vue))
- Exists but need to verify coverage

---

## Priority Ranking

| Priority | Field | Issue | Effort |
|----------|-------|-------|--------|
| **P1** | String | Add rich text editor support (editor ignores `useRichTextEditor`) | Medium |
| **P1** | User | Build proper user editor with roles, selectType, search, presetWithAuthenticated | High |
| **P1** | Record | Pass `queryFields`, `prefilter`, `selectType` to `CInputRecord` | Medium |
| **P2** | Url | Create dedicated URL editor with `fixUrl` processing | Medium |
| **P2** | Number | Implement progress bar rendering in editor and viewer | Medium |
| **P2** | Email | Create dedicated email editor with `type="email"` input | Low |
| **P2** | File | Wire `allowImages`/`allowDocuments` and verify viewer options | Medium |
| **P3** | Geometry | Entirely missing editor + viewer (requires map library) | High |

---

## Files Referenced

### Human — Editors (`lib/vue/src/components/field/editors/`)
- [CFieldStringEditor.vue](file:///home/fajfa/Human/human/lib/vue/src/components/field/editors/CFieldStringEditor.vue)
- [CFieldNumberEditor.vue](file:///home/fajfa/Human/human/lib/vue/src/components/field/editors/CFieldNumberEditor.vue)
- [CFieldBoolEditor.vue](file:///home/fajfa/Human/human/lib/vue/src/components/field/editors/CFieldBoolEditor.vue)
- [CFieldDateTimeEditor.vue](file:///home/fajfa/Human/human/lib/vue/src/components/field/editors/CFieldDateTimeEditor.vue)
- [CFieldSelectEditor.vue](file:///home/fajfa/Human/human/lib/vue/src/components/field/editors/CFieldSelectEditor.vue)
- [CFieldUserEditor.vue](file:///home/fajfa/Human/human/lib/vue/src/components/field/editors/CFieldUserEditor.vue)
- [CFieldRecordEditor.vue](file:///home/fajfa/Human/human/lib/vue/src/components/field/editors/CFieldRecordEditor.vue)
- [CFieldFileEditor.vue](file:///home/fajfa/Human/human/lib/vue/src/components/field/editors/CFieldFileEditor.vue)

### Human — Configurators (`compose/src/components/ModuleFields/Configurator/kinds/`)
- [String.vue](file:///home/fajfa/Human/human/client/web/compose/src/components/ModuleFields/Configurator/kinds/String.vue)
- [Number.vue](file:///home/fajfa/Human/human/client/web/compose/src/components/ModuleFields/Configurator/kinds/Number.vue)
- [User.vue](file:///home/fajfa/Human/human/client/web/compose/src/components/ModuleFields/Configurator/kinds/User.vue)
- [Record.vue](file:///home/fajfa/Human/human/client/web/compose/src/components/ModuleFields/Configurator/kinds/Record.vue)
- [Select.vue](file:///home/fajfa/Human/human/client/web/compose/src/components/ModuleFields/Configurator/kinds/Select.vue)
- [DateTime.vue](file:///home/fajfa/Human/human/client/web/compose/src/components/ModuleFields/Configurator/kinds/DateTime.vue)

### Corteza — Reference editors (`compose/src/components/ModuleFields/Editor/`)
- [String.vue](file:///home/fajfa/Corteza/corteza/client/web/compose/src/components/ModuleFields/Editor/String.vue) — has rich text editor
- [User.vue](file:///home/fajfa/Corteza/corteza/client/web/compose/src/components/ModuleFields/Editor/User.vue) — 390 lines, full implementation
- [Url.vue](file:///home/fajfa/Corteza/corteza/client/web/compose/src/components/ModuleFields/Editor/Url.vue) — URL processing
- [Email.vue](file:///home/fajfa/Corteza/corteza/client/web/compose/src/components/ModuleFields/Editor/Email.vue) — type="email"

### Shared — lib/js type definitions (`lib/js/src/compose/types/module-field/`)
- [base.ts](file:///home/fajfa/Human/human/lib/js/src/compose/types/module-field/base.ts)
- [string.ts](file:///home/fajfa/Human/human/lib/js/src/compose/types/module-field/string.ts)
- [bool.ts](file:///home/fajfa/Human/human/lib/js/src/compose/types/module-field/bool.ts)
- [datetime.ts](file:///home/fajfa/Human/human/lib/js/src/compose/types/module-field/datetime.ts)
- [number.ts](file:///home/fajfa/Human/human/lib/js/src/compose/types/module-field/number.ts)
- [select.ts](file:///home/fajfa/Human/human/lib/js/src/compose/types/module-field/select.ts)
- [user.ts](file:///home/fajfa/Human/human/lib/js/src/compose/types/module-field/user.ts)
- [record.ts](file:///home/fajfa/Human/human/lib/js/src/compose/types/module-field/record.ts)
- [file.ts](file:///home/fajfa/Human/human/lib/js/src/compose/types/module-field/file.ts)
- [email.ts](file:///home/fajfa/Human/human/lib/js/src/compose/types/module-field/email.ts)
- [url.ts](file:///home/fajfa/Human/human/lib/js/src/compose/types/module-field/url.ts)
- [geometry.ts](file:///home/fajfa/Human/human/lib/js/src/compose/types/module-field/geometry.ts)
