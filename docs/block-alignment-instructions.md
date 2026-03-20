# Page Block Alignment: Audit & Fix Instructions

## Goal
Ensure Human's page block viewers/editors/configurators behave identically to Corteza's. Same methodology we used for field types.

---

## File Locations

### Corteza (reference — DO NOT modify)
| Component | Path |
|---|---|
| Block viewers | `corteza/client/web/compose/src/components/PageBlocks/*Base.vue` |
| Block configurators | `corteza/client/web/compose/src/components/PageBlocks/*Configurator.vue` |
| Block editor | `corteza/client/web/compose/src/components/PageBlocks/RecordEditor.vue` |
| lib/js block types | `corteza/lib/js/src/compose/types/page-block/` |

### Human (modify these)
| Component | Path |
|---|---|
| Block viewers | `human/client/web/compose/src/components/PageBlocks/Blocks/*Block.vue` |
| Block configurators | `human/client/web/compose/src/components/PageBlocks/Configurators/*Configurator.vue` |
| lib/js block types | `human/lib/js/src/compose/types/page-block/` |

---

## Block Types to Audit

| Block | Corteza | Human | Notes |
|---|---|---|---|
| Record | RecordBase + RecordEditor | RecordBlock | Single record view/edit |
| RecordList | RecordListBase | RecordListBlock | Table of records |
| Chart | ChartBase | ChartBlock | Chart display |
| Calendar | CalendarBase | CalendarBlock | Calendar view |
| Content | ContentBase | ContentBlock | Rich text/HTML |
| File | FileBase | FileBlock | File attachments |
| IFrame | IFrameBase | IFrameBlock | Embedded iframe |
| Metric | MetricBase | MetricBlock | KPI/metric display |
| Comment | Comment/ | CommentBlock | Discussion thread |
| Automation | AutomationBase | AutomationBlock | Button triggers |
| **Missing in Human** | GeometryBase | ❌ | Map block |
| **Missing in Human** | ProgressBase | ❌ | Progress display |
| **Missing in Human** | RecordOrganizer | ❌ | Kanban board |
| **Missing in Human** | RecordRevisions | ❌ | Revision history |
| **Missing in Human** | SocialFeed | ❌ | Social feed |
| **Missing in Human** | Tabs | ❌ | Tab container |
| **Missing in Human** | Navigation | ❌ | Nav links |
| **Missing in Human** | Report | ❌ | Report display |

---

## Audit Methodology (per block type)

### Step 1: Read the lib/js class (ground truth)
```
human/lib/js/src/compose/types/page-block/<type>.ts
```
List ALL options/properties. These are what the configurator should expose and the viewer should use.

### Step 2: Compare Corteza viewer → Human viewer
Open both side-by-side. Check:
- [ ] Every option from lib/js is consumed correctly
- [ ] Options that **only apply in certain modes** are not applied globally (like the Number min/max bug)
- [ ] Default values match lib/js defaults
- [ ] Multi-value handling is correct
- [ ] Format/display logic matches

### Step 3: Compare Corteza configurator → Human configurator
- [ ] Same options are exposed
- [ ] Same layout order (users shouldn't feel lost)
- [ ] Conditional visibility matches (e.g., "show X only when Y is checked")
- [ ] Live previews / examples work and are reactive
- [ ] Descriptions/tooltips/footnotes are present

### Step 4: Check locale keys
```
human/locale/en/corteza-webapp-compose/block.yaml
```
- [ ] All keys used in templates exist (no fallback translations!)
- [ ] Check `old-locale/` if keys are missing — reuse from there

### Step 5: Test reactivity
- [ ] Change options in the configurator → viewer should update immediately
- [ ] Toggling modes should show/hide the correct sections
- [ ] Switching between options shouldn't leave stale values

---

## Common Pitfalls (from field type audit)

### Don't apply mode-specific defaults globally
Example: Number field had `min: 0, max: 100` as defaults for progress mode, but they were constraining the number input too. Always check if an option is mode-specific before passing it to the viewer/editor.

### Use the block's own methods for formatting
Don't reimplement formatting logic — use the class's `formatValue()` or equivalent. Otherwise live examples and viewers will diverge.

### PrimeVue ≠ Bootstrap
When replicating Corteza's behavior:
- Bootstrap's `variant="success"` → use PrimeVue `severity` or CSS variable coloring
- Bootstrap's `:selectable` → PrimeVue filters options instead
- Bootstrap's `b-form-radio-group` → PrimeVue `RadioButton` loop
- Always use PrimeVue style classes (`text-color`, `text-muted-color`, `bg-emphasis`, etc.)
- Never use `dark:` classes or hardcoded surface values

### Always use translations
Check `locale/` first, then `old-locale/`. If missing from both, add new keys to `locale/`. Never use hardcoded strings or fallback translations.

### Register new components
If you use a new PrimeVue component, check `primevue-components.ts` first — it may already be registered. If not, register it there.

---

## Checklist Template (copy per block)

```markdown
## [Block Name] Block

### lib/js Options
- [ ] List all options from `page-block/<type>.ts`

### Viewer
- [ ] All options consumed correctly
- [ ] No mode-specific options applied globally
- [ ] Display logic matches Corteza

### Configurator
- [ ] All options exposed
- [ ] Layout order matches Corteza
- [ ] Conditional sections work
- [ ] Descriptions/labels present
- [ ] Reactive (changes reflect immediately)

### Locale
- [ ] All keys exist in `block.yaml`
```

---

## When done
Update this doc with your findings per block. Mark ✅ for aligned, ⚠️ for gaps found + fixed, ❌ for missing blocks that need full implementation.
