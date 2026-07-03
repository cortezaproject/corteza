<template>
  <div class="flex flex-col min-h-0">
    <!-- Legend (+ capability selector in the multi-role step view) -->
    <div v-if="roles.length" class="shrink-0 flex flex-wrap items-center gap-x-4 gap-y-2 mb-3">
      <SelectButton
        v-if="!hideRoleHeader"
        v-model="capability"
        :options="capOptions"
        option-label="label"
        option-value="value"
        :allow-empty="false"
      />
      <!-- Legend mirrors the cell icons exactly (same stateIcon), so the key
           always matches what's rendered in the grid. -->
      <div class="flex flex-wrap items-center gap-x-3 gap-y-1 text-xs text-muted-color">
        <span v-for="l in LEGEND" :key="l.state" class="inline-flex items-center gap-1.5">
          <i :class="stateIcon(l.state)" />
          {{ $t(`project.permissions.access.${l.labelKey}`) }}
        </span>
      </div>
    </div>

    <div v-if="!roles.length" class="grid place-items-center text-center px-6 py-10 text-muted-color">
      <div>
        <i class="pi pi-id-card text-4xl mb-2" />
        <p class="text-sm">{{ $t('project.permissions.noRoles') }}</p>
      </div>
    </div>
    <div
      v-else-if="!visibleRows.length"
      class="grid place-items-center text-center px-6 py-10 text-muted-color"
    >
      <div>
        <i class="pi pi-sitemap text-4xl mb-2" />
        <p class="text-sm">{{ $t('project.permissions.noResources') }}</p>
      </div>
    </div>

    <!-- ============ Single-role detail view: all capabilities per row ============ -->
    <div v-else-if="hideRoleHeader" class="min-h-0 overflow-auto">
      <TransitionGroup
        tag="div"
        name="rowanim"
        class="rounded-lg border border-surface divide-y divide-surface"
      >
        <div
          v-for="row in visibleRows"
          :key="row.key"
          class="group flex items-center gap-2 pr-2 min-h-[2.75rem] transition-colors hover:bg-emphasis"
          :class="row.rowBg"
        >
          <div
            class="flex items-center gap-2 min-w-0 flex-1 py-1.5"
            :class="row.children.length ? 'cursor-pointer select-none' : ''"
            :style="{ paddingLeft: 0.5 + row.level * 0.9 + 'rem' }"
            @click="row.children.length && toggle(row)"
          >
            <span
              v-if="row.icon"
              class="inline-flex items-center justify-center w-6 h-6 rounded-md ring-1 shrink-0"
              :class="[row.icon.bg, row.icon.ring]"
            >
              <i :class="[row.icon.icon, row.icon.text, 'text-xs']" />
            </span>
            <span :class="row.labelClass" class="truncate min-w-0">{{ row.label }}</span>
            <span
              v-if="row.tag"
              class="shrink-0 text-[11px] font-medium px-1.5 py-0.5 rounded bg-emphasis text-muted-color border border-surface"
            >
              {{ row.tag }}
            </span>
            <i
              v-if="row.children.length"
              class="pi shrink-0 text-xs text-muted-color"
              :class="isOpen(row) ? 'pi-chevron-down' : 'pi-chevron-right'"
            />
            <!-- Full permissions editor, next to the collapse chevron. -->
            <button
              v-if="row.advancedResource"
              type="button"
              class="shrink-0 inline-flex items-center justify-center w-6 h-6 rounded-full text-muted-color hover:bg-surface-200 dark:hover:bg-surface-700"
              :disabled="disabled"
              v-tooltip.bottom="{ value: $t('project.permissions.advancedFor', { target: row.tipTarget }), showDelay: 400 }"
              @click.stop="openAdvanced(row, soleRole.id)"
            >
              <i class="pi pi-lock text-xs" />
            </button>
          </div>

          <div class="flex items-center gap-3 shrink-0 pl-2">
            <div v-for="c in rowCaps(row)" :key="c.key" class="flex items-center gap-1">
              <span class="text-[11px] text-muted-color">{{ c.label }}</span>
              <button
                type="button"
                class="inline-flex items-center justify-center w-7 h-7 rounded-full hover:bg-surface-200 dark:hover:bg-surface-700 disabled:cursor-default"
                :disabled="disabled"
                v-tooltip.bottom="{ value: capTooltip(row, c, soleRole), showDelay: 400 }"
                @click="toggleCap(row, c, soleRole)"
              >
                <i :class="stateIcon(capState(soleRole.id, c.ops))" />
              </button>
            </div>
          </div>
        </div>
      </TransitionGroup>
    </div>

    <!-- ============ Multi-role step view: roles as columns + capability lens ============ -->
    <div v-else class="min-h-0 overflow-auto">
      <div class="inline-block min-w-full rounded-lg border border-surface overflow-hidden">
        <table class="text-sm border-collapse w-full">
          <thead>
            <tr class="bg-emphasis">
              <th class="text-left font-medium px-3 py-2 sticky left-0 z-10 bg-emphasis min-w-80" />
              <th
                v-for="role in roles"
                :key="role.id"
                class="font-medium px-3 py-2 border-l border-surface text-center whitespace-nowrap"
              >
                <span class="inline-flex items-center justify-center gap-1.5">
                  <i :class="[roleCfg.icon, roleCfg.text, 'text-xs']" />
                  {{ role.name }}
                </span>
              </th>
              <th aria-hidden="true" class="w-full p-0 border-l border-surface" />
            </tr>
          </thead>
          <TransitionGroup tag="tbody" name="rowanim">
            <tr
              v-for="row in visibleRows"
              :key="row.key"
              class="group border-t border-surface transition-colors hover:!bg-highlight"
              :class="row.rowBg"
            >
              <td
                class="pr-2 py-2 sticky left-0 z-10 bg-inherit"
                :class="row.children.length ? 'cursor-pointer select-none' : ''"
                :style="{ paddingLeft: 0.5 + row.level * 0.9 + 'rem' }"
                :aria-expanded="row.children.length ? isOpen(row) : undefined"
                @click="row.children.length && toggle(row)"
              >
                <div class="min-w-0">
                  <div class="flex items-center gap-2 min-w-0">
                    <span
                      v-if="row.icon"
                      class="inline-flex items-center justify-center w-6 h-6 rounded-md ring-1 shrink-0"
                      :class="[row.icon.bg, row.icon.ring]"
                    >
                      <i :class="[row.icon.icon, row.icon.text, 'text-xs']" />
                    </span>
                    <span :class="row.labelClass" class="truncate min-w-0">{{ row.label }}</span>
                    <span
                      v-if="row.tag"
                      class="shrink-0 text-[11px] font-medium px-1.5 py-0.5 rounded bg-emphasis text-muted-color border border-surface"
                    >
                      {{ row.tag }}
                    </span>
                    <i
                      v-if="row.children.length"
                      class="pi shrink-0 pl-1 text-xs text-muted-color"
                      :class="isOpen(row) ? 'pi-chevron-down' : 'pi-chevron-right'"
                    />
                    <!-- Full permissions editor, next to the collapse chevron. -->
                    <button
                      v-if="row.advancedResource"
                      type="button"
                      class="shrink-0 inline-flex items-center justify-center w-6 h-6 rounded-full text-muted-color hover:bg-emphasis opacity-0 focus:opacity-100 group-hover:opacity-100 transition-opacity"
                      :disabled="disabled"
                      v-tooltip.bottom="{ value: $t('project.permissions.advancedFor', { target: row.tipTarget }), showDelay: 500 }"
                      @click.stop="openAdvanced(row)"
                    >
                      <i class="pi pi-lock text-xs" />
                    </button>
                  </div>
                  <!-- What the current-lens toggle actually grants on this row. -->
                  <p
                    v-if="laneDesc(row)"
                    class="text-xs text-muted-color leading-snug mt-0.5 whitespace-nowrap"
                    :style="row.icon ? { paddingLeft: '2rem' } : {}"
                  >
                    {{ laneDesc(row) }}
                  </p>
                </div>
              </td>
              <!-- Section-header rows are pure group labels: render their role
                   columns as an empty, borderless divider strip (no faint dashes). -->
              <td
                v-for="role in roles"
                :key="role.id"
                :class="isHeaderRow(row) ? '' : 'border-l border-surface'"
              >
                <div
                  v-if="!isHeaderRow(row)"
                  class="flex items-center justify-center px-1 py-1 min-h-[2.25rem]"
                >
                  <!-- State is the (always-visible) icon; clicking it toggles.
                       No background fill, no per-cell gear — the full editor
                       lives on the resource row (left). -->
                  <button
                    type="button"
                    class="inline-flex items-center justify-center w-7 h-7 rounded-full enabled:hover:bg-emphasis disabled:cursor-default"
                    :disabled="disabled || !laneOps(row).length || capState(role.id, laneOps(row)) === 'loading'"
                    v-tooltip.bottom="{ value: laneTooltip(row, role), showDelay: 500 }"
                    @click="toggleLane(row, role)"
                  >
                    <i :class="stateIcon(capState(role.id, laneOps(row)))" />
                  </button>
                </div>
              </td>
              <td aria-hidden="true" :class="isHeaderRow(row) ? 'w-full p-0' : 'w-full border-l border-surface'" />
            </tr>
          </TransitionGroup>
        </table>
      </div>
    </div>

    <!-- Toggle confirm. When the clicked row has nested children it offers the
         cascade choice as two action buttons (only this / include everything
         nested); a leaf row gets a single confirm. -->
    <Dialog
      v-model:visible="confirmVisible"
      modal
      :draggable="false"
      :header="$t('project.permissions.access.confirm.header')"
      class="w-full max-w-md"
      :pt="{ footer: { class: 'flex justify-end gap-2' } }"
    >
      <div class="flex items-start gap-3">
        <i
          class="text-xl mt-0.5"
          :class="pending?.next === 'deny' ? 'pi pi-ban text-red-500' : 'pi pi-check-circle text-green-500'"
        />
        <p class="text-sm leading-relaxed">{{ confirmMessage }}</p>
      </div>
      <template #footer>
        <Button
          :label="$t('general.label.cancel')"
          severity="secondary"
          text
          size="small"
          @click="pending = null"
        />
        <template v-if="pending?.hasNested">
          <Button
            :label="$t('project.permissions.access.confirm.onlyThis')"
            outlined
            size="small"
            :severity="pending?.next === 'deny' ? 'danger' : 'success'"
            @click="applyToggle(false)"
          />
          <Button
            :label="$t('project.permissions.access.confirm.includeNested')"
            size="small"
            :severity="pending?.next === 'deny' ? 'danger' : 'success'"
            @click="applyToggle(true)"
          />
        </template>
        <Button
          v-else
          :label="$t('project.permissions.access.confirm.save')"
          size="small"
          :severity="pending?.next === 'deny' ? 'danger' : 'success'"
          @click="applyToggle(false)"
        />
      </template>
    </Dialog>
  </div>
</template>

<script setup>
import { fieldTypeLabelKey } from '@/sections/project/config/fieldTypes'
import { kindConfig } from '@/sections/project/config/kinds'
import { useProjectsStore } from '@/sections/project/stores/projects'
import { usePermissions } from '@planetcrust/human-vue'
import { computed, inject, reactive, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'

const props = defineProps({
  project: { type: Object, required: true },
  // Role columns. All the project's roles in the step view; a single-element
  // array for one role's detail view.
  roles: { type: Array, default: () => [] },
  // Single-role (detail) mode: no role-name columns, no capability selector —
  // every capability is shown per row for the one role.
  hideRoleHeader: { type: Boolean, default: false },
  disabled: { type: Boolean, default: false },
})

const store = useProjectsStore()
const { t } = useI18n()
const $toast = inject('$toast')
const { open: openPermissions } = usePermissions()

const ns = computed(() => props.project.namespaceID)
const soleRole = computed(() => props.roles[0] || null)

// The kinds an END USER of the deployed app interacts with. These access roles
// are for people using the app (pages + records), not building it — so only the
// runtime surface is listed: data (modules → records/fields), pages (view), and
// the AI they can use (agents/chatbots). Build-time kinds (connections,
// automations) and build-time ops (create/edit modules & pages, manage AI) are
// intentionally excluded; rare ops stay reachable via the ⚙ full editor.
const KINDS = [
  { kind: 'module', getter: 'resourcesFor', type: 'corteza::compose:module', ns: true },
  { kind: 'page', getter: 'pagesFor', type: 'corteza::compose:page', ns: true },
  { kind: 'agent', getter: 'agentsFor', type: 'corteza::system:agent', ns: false },
  { kind: 'chatbot', getter: 'chatbotsFor', type: 'corteza::system:chatbot', ns: false },
]

const HEADER_BG = 'bg-emphasis'
const ITEM_BG = 'bg-surface'

// Role kind visual config (id-card icon + violet) — shown next to each role name
// in the column headers, matching the role styling used elsewhere.
const roleCfg = kindConfig('role')

function resStr(def, id) {
  return def.ns ? `${def.type}/${ns.value}/${id}` : `${def.type}/${id}`
}

// Localized field-kind label (falls back to the raw kind for unknown types) —
// shown as a small tag on each field row, the way the module builder does.
function fieldTypeText(type) {
  const key = fieldTypeLabelKey(type)
  return key ? t(key) : type
}

// The capability order used everywhere (the step-view lens options and the
// column order in the detail view).
const CAP_ORDER = ['create', 'read', 'update', 'delete']

// Each node carries a `caps` map: capability → [{ resource, op }] it reads/writes
// (absent capability = not applicable to that row). A capability may span several
// ops on different resources — e.g. Records → read = record.read + records.search
// (view AND list), the single honest "can see records" toggle. `descGroup` +
// `descTarget` drive the plain-language tooltips so a cell never hides what it
// touches.
const tree = computed(() => {
  const nsID = ns.value
  const nodes = []

  // App access (namespace read) is intentionally NOT listed: every user added to
  // the project already reads its namespace, so a toggle here would be misleading.

  for (const def of KINDS) {
    const items = store[def.getter](props.project.id) || []
    if (!items.length) continue
    const cfg = kindConfig(def.kind)

    const section = {
      key: `k:${def.kind}`,
      label: t(cfg.labelKey),
      level: 0,
      tipTarget: t(cfg.labelKey),
      icon: cfg,
      rowBg: HEADER_BG,
      labelClass: 'font-medium',
      defaultOpen: true,
      caps: {},
      children: [],
    }

    for (const it of items) {
      if (def.kind === 'module') {
        const moduleRes = resStr(def, it.id)
        const recordRes = `corteza::compose:record/${nsID}/${it.id}/*`
        const fieldsRes = `corteza::compose:module-field/${nsID}/${it.id}/*`

        // Records row (nested under the module separator) = the data end users
        // work with. Read = "see the data": record.read + records.search (list) +
        // module.read (load the model) + record.value.read (all field values) —
        // one honest toggle covering records AND their values. Create is the
        // module-owned record.create; Update/Delete map to the record ops.
        // Individual fields nest directly below to refine per-field value access.
        const recordsNode = {
          key: `records:${it.id}`,
          label: t('project.permissions.records'),
          level: 2,
          tipTarget: t('project.permissions.recordsOf', { target: it.name }),
          icon: null,
          rowBg: ITEM_BG,
          labelClass: '',
          defaultOpen: false,
          advancedResource: recordRes,
          advancedTitle: it.name,
          advancedAllSpecific: true,
          descGroup: 'records',
          descTarget: it.name,
          caps: {
            create: [{ resource: moduleRes, op: 'record.create' }],
            read: [
              { resource: recordRes, op: 'read' },
              { resource: moduleRes, op: 'records.search' },
              { resource: fieldsRes, op: 'record.value.read' },
            ],
            update: [{ resource: recordRes, op: 'update' }],
            delete: [{ resource: recordRes, op: 'delete' }],
          },
          children: (it.fields || []).map(f => {
            const fieldRes = `corteza::compose:module-field/${nsID}/${it.id}/${f.id}`
            return {
              key: `field:${it.id}:${f.id}`,
              label: f.name,
              level: 3,
              tipTarget: f.name,
              icon: null,
              tag: fieldTypeText(f.type),
              rowBg: ITEM_BG,
              labelClass: '',
              defaultOpen: false,
              advancedResource: fieldRes,
              advancedTitle: f.name,
              descGroup: 'field',
              descTarget: f.name,
              capLabels: { update: 'project.permissions.caps.write' },
              caps: {
                read: [{ resource: fieldRes, op: 'record.value.read' }],
                update: [{ resource: fieldRes, op: 'record.value.update' }],
              },
              children: [],
            }
          }),
        }

        // Module row = a single "can this role SEE the module" toggle (module.read).
        // No icon + lighter font so it doesn't compete with the section header;
        // records/fields nest beneath. Schema create/edit/delete is build-time and
        // stays out of the matrix (⚙ full editor only).
        section.children.push({
          key: `module:${it.id}`,
          label: it.name,
          level: 1,
          tipTarget: it.name,
          icon: null,
          rowBg: ITEM_BG,
          labelClass: '',
          defaultOpen: true,
          advancedResource: moduleRes,
          advancedTitle: it.name,
          descGroup: 'module',
          descTarget: it.name,
          caps: { read: [{ resource: moduleRes, op: 'read' }] },
          children: [recordsNode],
        })
        continue
      }

      if (def.kind === 'page') {
        // Pages: end users only VIEW them (read). Edit/delete is build-time.
        const res = resStr(def, it.id)
        section.children.push({
          key: `page:${it.id}`,
          label: it.name,
          level: 1,
          tipTarget: it.name,
          icon: null,
          rowBg: ITEM_BG,
          labelClass: '',
          defaultOpen: false,
          advancedResource: res,
          advancedTitle: it.name,
          descGroup: 'page',
          descTarget: it.name,
          capLabels: { read: 'project.permissions.caps.view' },
          caps: { read: [{ resource: res, op: 'read' }] },
          children: [],
        })
        continue
      }

      // Agents / chatbots: a single "Use" (read) toggle — end users interact with
      // the embedded widget. NOTE: starting a conversation/session also needs a
      // platform-wide grant the project can't set here; the tooltip says so.
      const res = resStr(def, it.id)
      section.children.push({
        key: `${def.kind}:${it.id}`,
        label: it.name,
        level: 1,
        tipTarget: it.name,
        icon: null,
        rowBg: ITEM_BG,
        labelClass: '',
        defaultOpen: false,
        advancedResource: res,
        advancedTitle: it.name,
        descGroup: def.kind,
        descTarget: it.name,
        capLabels: { read: 'project.permissions.caps.use' },
        caps: { read: [{ resource: res, op: 'read' }] },
        children: [],
      })
    }

    nodes.push(section)
  }

  return nodes
})

// --- expand / collapse --------------------------------------------------------
const overrides = reactive(new Map())
function isOpen(node) {
  return overrides.has(node.key) ? overrides.get(node.key) : node.defaultOpen
}
function toggle(node) {
  overrides.set(node.key, !isOpen(node))
}

const visibleRows = computed(() => {
  const out = []
  const walk = list => {
    for (const n of list) {
      out.push(n)
      if (n.children.length && isOpen(n)) walk(n.children)
    }
  }
  walk(tree.value)
  return out
})

// Section headers carry no capabilities of their own (pure group labels). Their
// role columns render as an empty, borderless divider rather than faint dashes.
function isHeaderRow(row) {
  return !Object.keys(row.caps || {}).length
}

// --- effective access evaluation ----------------------------------------------
// Every resource referenced by any capability across the whole tree, deduped.
// The store traces all ops per resource, so this covers every cell (incl. the
// module-owned record.create / records.search behind the Records row).
const allResources = computed(() => {
  const set = new Set()
  const walk = list => {
    for (const n of list) {
      for (const cap of CAP_ORDER) {
        for (const { resource } of n.caps[cap] || []) set.add(resource)
      }
      if (n.advancedResource) set.add(n.advancedResource)
      if (n.children.length) walk(n.children)
    }
  }
  walk(tree.value)
  return [...set]
})

watch(
  [() => props.project.id, () => props.roles, allResources, () => store.graphVersion],
  () => {
    if (props.project.id) store.loadEffectiveAccess(props.project.id, allResources.value)
  },
  { immediate: true },
)

// --- capability lens (step view) ----------------------------------------------
const capability = ref('read')
const capOptions = computed(() =>
  CAP_ORDER.map(c => ({ value: c, label: t(`project.permissions.caps.${c}`) })),
)
// The ops a row exposes under the currently-selected lens (step view).
function laneOps(row) {
  return row.caps[capability.value] || []
}

// Plain-language description of what the current-lens toggle actually grants on
// this row — shown under the resource name. Empty when the row has no toggle for
// the selected capability (headers, or a lens the row doesn't expose).
function laneDesc(row) {
  return laneOps(row).length ? capDesc(row, capability.value) : ''
}

// The capabilities a row exposes, in order, for the detail view (all at once).
function rowCaps(node) {
  return CAP_ORDER.filter(c => (node.caps[c] || []).length).map(c => ({
    key: c,
    label: node.capLabels?.[c] ? t(node.capLabels[c]) : t(`project.permissions.caps.${c}`),
    ops: node.caps[c],
  }))
}

// --- cell state (compound-aware) ----------------------------------------------
// Combine the effective access of every op behind a capability into one honest
// state: allow/deny/unknown only when ALL ops agree; 'partial' when they differ
// (e.g. record.read allowed but records.search denied); 'na' when the capability
// doesn't apply; 'loading' while the trace is in flight; 'none' when unevaluated.
function capState(roleId, ops) {
  if (!ops || !ops.length) return 'na'
  const states = ops.map(({ resource, op }) =>
    store.effectiveAccess(props.project.id, roleId, resource, op),
  )
  if (states.some(s => !s) && store.isEffectiveAccessLoading(props.project.id)) return 'loading'
  const defined = states.filter(Boolean)
  if (!defined.length) return 'none'
  const uniq = [...new Set(defined)]
  // All ops present and unanimous → that state; anything else (disagreement, or
  // some ops unevaluated) → partial, so the cell never overstates a clean grant.
  if (defined.length === ops.length && uniq.length === 1) return uniq[0]
  return 'partial'
}

function stateIcon(state) {
  switch (state) {
    case 'allow':
      return 'pi pi-check-circle text-green-500'
    case 'deny':
      return 'pi pi-times-circle text-red-500'
    case 'partial':
      return 'pi pi-exclamation-circle text-amber-500'
    case 'unknown':
      return 'pi pi-question-circle text-muted-color'
    case 'loading':
      // In-flight trace — never render as a grey minus (that reads as
      // 'not granted'); a spinner keeps the cell honest until it resolves.
      return 'pi pi-spin pi-spinner text-muted-color'
    case 'na':
      return 'pi pi-minus text-muted-color opacity-40'
    default:
      return 'pi pi-minus text-muted-color'
  }
}

// Legend entries (in order). Rendered with the very same stateIcon() the cells
// use, so the key can never drift from what's shown. Partial has its own icon
// (a half-filled marker) distinct from allow/deny/unknown.
const LEGEND = [
  { state: 'allow', labelKey: 'allowed' },
  { state: 'deny', labelKey: 'denied' },
  { state: 'partial', labelKey: 'partial' },
  { state: 'unknown', labelKey: 'unknown' },
]

// --- honest descriptions & tooltips -------------------------------------------
// Plain-language statement of exactly what a capability controls, naming the
// target (and, for the coupled Records → Read, both underlying actions).
function capDesc(node, cap) {
  return t(`project.permissions.capDesc.${node.descGroup}.${cap}`, { target: node.descTarget })
}

// Full-sentence tooltip: current state + what a click does. Always spells out the
// underlying actions so the toggle can't deceive.
function tooltipFor(state, desc, actionable) {
  const p = 'project.permissions.capState.'
  let text
  switch (state) {
    case 'allow':
      text = t(`${p}allow`, { desc })
      break
    case 'deny':
      text = t(`${p}deny`, { desc })
      break
    case 'partial':
      return t(`${p}partial`, { desc })
    case 'unknown':
      text = t(`${p}unknown`, { desc })
      break
    case 'loading':
      return t('project.permissions.access.loading')
    case 'na':
      return t('project.permissions.access.naSimple', { desc })
    default:
      text = t(`${p}none`, { desc })
  }
  if (actionable) {
    const next = state === 'allow' ? 'deny' : 'allow'
    text += ` · ${t(`${p}clickTo.${next}`)}`
  }
  return text
}

// Detail view (one role): tooltip for a capability toggle.
function capTooltip(node, c, role) {
  if (!role) return ''
  return tooltipFor(capState(role.id, c.ops), capDesc(node, c.key), !props.disabled)
}

// Step view: tooltip for the selected-lens cell.
function laneTooltip(row, role) {
  const ops = laneOps(row)
  if (!ops.length) return t('project.permissions.access.naSimple', { desc: row.tipTarget })
  return tooltipFor(capState(role.id, ops), capDesc(row, capability.value), !props.disabled)
}

// --- toggling -----------------------------------------------------------------
// All ops for a capability across a node AND its descendants (deduped by
// resource+op). Toggling a parent cascades the same grant/deny to every nested
// child — module Read reaches Records + every field, Records Update reaches each
// field's Write, etc.
function capOpsDeep(node, cap) {
  const seen = new Set()
  const out = []
  const walk = n => {
    for (const o of n.caps[cap] || []) {
      const k = `${o.resource}::${o.op}`
      if (!seen.has(k)) {
        seen.add(k)
        out.push(o)
      }
    }
    for (const c of n.children || []) walk(c)
  }
  walk(node)
  return out
}

// A staged toggle awaiting the user's choice in the confirm dialog. Holds the
// clicked row's OWN ops and its full-subtree ops, so — when the row has nested
// children — the dialog can offer "only this" vs "this and everything nested"
// (the same two-button shape as the resource-create dialogs).
const pending = ref(null)

const confirmVisible = computed({
  get: () => !!pending.value,
  set: v => {
    if (!v) pending.value = null
  },
})

const confirmMessage = computed(() => {
  const p = pending.value
  if (!p) return ''
  return t(`project.permissions.access.confirm.${p.next}`, { role: p.role.name, desc: p.desc })
})

// Stage a toggle. Direction comes from the row's OWN state (allow → deny;
// deny/partial/unknown/none → allow). `hasNested` decides whether the dialog
// shows the cascade choice or a single confirm.
function requestToggle(role, node, capKey, ownOps, desc) {
  if (!role) return
  const state = capState(role.id, ownOps)
  if (props.disabled || !ownOps.length || state === 'loading') return
  const deepOps = capOpsDeep(node, capKey)
  pending.value = {
    role,
    desc,
    next: state === 'allow' ? 'deny' : 'allow',
    ownOps,
    deepOps,
    hasNested: deepOps.length > ownOps.length,
  }
}

// Commit the staged toggle, writing either the row's own ops or its whole
// subtree. The optimistic patch + re-trace flip every affected cell.
async function applyToggle(cascade) {
  const p = pending.value
  if (!p) return
  const ops = cascade ? p.deepOps : p.ownOps
  pending.value = null
  try {
    await store.setCapabilityAccess(props.project.id, p.role.id, ops, p.next)
    $toast.toastSuccess(t('project.permissions.toastSaved'))
  } catch (e) {
    $toast.toastErrorHandler(t('project.permissions.toastSaveFailed'))(e)
  }
}

function toggleCap(node, c, role) {
  requestToggle(role, node, c.key, c.ops, capDesc(node, c.key))
}

function toggleLane(row, role) {
  requestToggle(role, row, capability.value, laneOps(row), capDesc(row, capability.value))
}

// Open the shared full permissions editor (all real ops) for the row's primary
// resource, preselecting the role. Covers the rare ops kept out of the matrix.
function openAdvanced(row, roleID) {
  if (props.disabled || !row.advancedResource) return
  openPermissions({
    resource: row.advancedResource,
    title: row.advancedTitle,
    target: row.advancedTitle,
    allSpecific: !!row.advancedAllSpecific,
    roleID,
    onSaved: () => store.touch(),
  })
}
</script>

<style scoped>
/* Row expand/collapse: entering rows slide+fade in, leaving rows fade out, and
   the rows below reflow smoothly to their new positions. */
.rowanim-enter-active {
  transition: opacity 0.2s ease, transform 0.2s ease;
}
.rowanim-leave-active {
  transition: opacity 0.12s ease;
}
.rowanim-enter-from {
  opacity: 0;
  transform: translateY(-6px);
}
.rowanim-leave-to {
  opacity: 0;
}
.rowanim-move {
  transition: transform 0.2s ease;
}
@media (prefers-reduced-motion: reduce) {
  .rowanim-enter-active,
  .rowanim-leave-active,
  .rowanim-move {
    transition: none;
  }
}
</style>
