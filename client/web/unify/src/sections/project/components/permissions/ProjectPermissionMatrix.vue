<template>
  <div class="flex flex-col min-h-0">
    <!-- Legend above the matrix. Mirrors the cell icons exactly (same stateIcon),
         so the key always matches what's rendered in the grid. Only the multi-role
         step view; the single-role detail view (role dialog) shows none. -->
    <div
      v-if="(roles.length || evalUserId) && !hideRoleHeader"
      class="shrink-0 flex flex-wrap items-center gap-x-3 gap-y-1 mb-3 text-xs text-muted-color"
    >
      <span v-for="l in LEGEND" :key="l.state" class="inline-flex items-center gap-1.5">
        <i :class="stateIcon(l.state)" />
        {{ $t(`project.permissions.access.${l.labelKey}`) }}
      </span>
    </div>

    <div
      v-if="!roles.length && !evalUserId"
      class="grid place-items-center text-center px-6 py-10 text-muted-color"
    >
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
        name="rowcollapse"
        class="rounded-lg border border-surface text-sm overflow-hidden"
      >
        <div v-for="row in visibleRows" :key="row.key" class="rowcollapse-item">
          <div class="overflow-hidden min-h-0">
            <div
              class="group flex items-center gap-2 pr-2 min-h-[2.75rem] border-t border-surface transition-colors hover:!bg-highlight"
              :class="row.rowBg"
            >
              <div
                class="flex items-center gap-2 min-w-0 flex-1 py-1.5"
                :class="row.children.length ? 'cursor-pointer select-none' : ''"
                :style="{ paddingLeft: indentRem(row) }"
                @click="row.children.length && toggle(row)"
              >
                <div class="min-w-0 flex-1">
                  <div class="flex items-center gap-2 min-w-0">
                    <KindIcon v-if="row.icon" :config="row.icon" size="md" />
                    <span :class="row.labelClass" class="truncate min-w-0">{{ row.label }}</span>
                    <FieldKindTag v-if="row.tag" :type="row.tag.type" />
                    <i
                      v-if="row.children.length"
                      class="pi shrink-0 pl-1 text-xs text-muted-color"
                      :class="isOpen(row) ? 'pi-chevron-down' : 'pi-chevron-right'"
                    />
                  </div>
                  <!-- What this row's capabilities grant, in plain language — matches
                   the step view's per-row description. -->
                  <p
                    v-if="rowDesc(row)"
                    class="text-xs text-muted-color leading-snug mt-0.5 whitespace-nowrap first-letter:uppercase"
                    :style="row.icon ? { paddingLeft: '2rem' } : {}"
                  >
                    {{ rowDesc(row) }}
                  </p>
                </div>
                <!-- Full permissions editor: same outlined secondary button as
                 CPermissionsButton elsewhere; right-aligned and hover-revealed.
                 Opaque highlight background so it masks the (non-truncated) row
                 description text it overlaps when revealed on hover. -->
                <Button
                  v-if="row.advancedResource"
                  icon="pi pi-lock"
                  severity="secondary"
                  outlined
                  size="small"
                  class="shrink-0 ml-auto !bg-highlight opacity-0 focus:opacity-100 group-hover:opacity-100 transition-opacity"
                  :disabled="disabled"
                  v-tooltip.bottom="{
                    value: $t('project.permissions.advancedFor', { target: row.tipTarget }),
                    showDelay: 400,
                  }"
                  @click.stop="openAdvanced(row, soleRole.id)"
                />
              </div>

              <div class="flex items-center gap-2 shrink-0 pl-2">
                <!-- One chip per capability, styled + behaving like the role toggle
                 chips in the user dialog: the state icon + label together, full
                 when allowed and dimmed when not; clicking toggles the grant. -->
                <button
                  v-for="c in rowCaps(row)"
                  :key="c.key"
                  type="button"
                  class="inline-flex items-center gap-1.5 pl-1.5 pr-3 py-1.5 rounded-lg border border-surface text-sm font-medium transition-all enabled:cursor-pointer disabled:cursor-default"
                  :class="
                    capState(soleRole.id, c.ops) === 'allow'
                      ? 'text-color'
                      : 'text-muted-color opacity-50 enabled:hover:opacity-100'
                  "
                  :disabled="disabled"
                  v-tooltip.bottom="{ value: capTooltip(row, c, soleRole), showDelay: 400 }"
                  @click="toggleCap(row, c, soleRole)"
                >
                  <span class="inline-flex items-center justify-center w-5 h-5 shrink-0">
                    <i :class="stateIcon(capState(soleRole.id, c.ops))" />
                  </span>
                  {{ c.label }}
                </button>
              </div>
            </div>
          </div>
        </div>
      </TransitionGroup>
    </div>

    <!-- ============ Multi-role step view: roles as columns + capability lens ============ -->
    <!-- Not a <table>: each row is its own fixed-column grid so it can animate its
         own height on expand/collapse (like the module card). Columns line up
         because every row shares the same fixed template (gridStyle). -->
    <div v-else class="min-h-0 overflow-auto">
      <!-- overflow-clip (not hidden) clips the rounded corners without becoming a
           scroll container, so the sticky name column below still anchors to the
           outer horizontal scroller. -->
      <div class="w-max min-w-full rounded-lg border border-surface overflow-clip">
        <!-- Header row: same column template as every data row. -->
        <div class="grid min-w-max bg-surface" :style="gridStyle">
          <!-- Corner cell: the capability lens lives here (instead of a separate
               header strip) so the matrix carries its own capability selector. -->
          <div class="sticky left-0 z-10 bg-surface px-3 py-2 flex items-center">
            <SelectButton
              v-model="capability"
              :options="capOptions"
              option-label="label"
              option-value="value"
              :allow-empty="false"
              size="small"
            />
          </div>
          <!-- Read-only "Evaluated" column: the user's resolved access. The tinted
               band + lock icon mark it non-editable (changed via the role columns). -->
          <div
            v-if="evalUserId"
            class="flex items-center justify-center gap-1.5 px-3 py-2 min-w-0 font-medium bg-black/[0.04] dark:bg-white/[0.04]"
          >
            <KindIcon kind="user" />
            <span class="truncate">
              {{ evalLabel || $t('project.userDetail.permissions.evaluated') }}
            </span>
            <i
              class="pi pi-lock text-[10px] text-muted-color shrink-0"
              v-tooltip.bottom="{
                value: $t('project.userDetail.permissions.evaluatedHint'),
                showDelay: 400,
              }"
            />
          </div>
          <!-- Fixed-width role columns: long names truncate; the header tooltip
               shows the full name. -->
          <div
            v-for="role in roles"
            :key="role.id"
            class="flex items-center justify-center gap-1.5 px-2 py-2 min-w-0 font-medium border-l border-surface"
            v-tooltip.bottom="{ value: role.name, showDelay: 400 }"
          >
            <KindIcon kind="role" />
            <span class="truncate">{{ role.name }}</span>
          </div>
          <div aria-hidden="true" class="border-l border-surface" />
        </div>

        <!-- Data rows. Each sits in a height-animating collapse box, so expanding
             or collapsing grows/shrinks it smoothly instead of popping in/out. -->
        <TransitionGroup tag="div" name="rowcollapse">
          <div v-for="row in visibleRows" :key="row.key" class="rowcollapse-item min-w-max">
            <!-- overflow-clip clips the collapsing height without creating a scroll
                 container (which would break the sticky name column). -->
            <div class="overflow-clip min-h-0 min-w-max">
              <div
                class="grid group min-w-max border-t border-surface transition-colors hover:!bg-highlight"
                :class="row.rowBg"
                :style="gridStyle"
              >
                <!-- Resource name (sticky while the role columns scroll). -->
                <div
                  class="sticky left-0 z-10 bg-inherit pr-2 py-2 min-h-[2.75rem] flex items-center min-w-0 overflow-hidden"
                  :class="row.children.length ? 'cursor-pointer select-none' : ''"
                  :style="{ paddingLeft: indentRem(row) }"
                  :aria-expanded="row.children.length ? isOpen(row) : undefined"
                  @click="row.children.length && toggle(row)"
                >
                  <div class="flex items-center gap-2 min-w-0 flex-1">
                    <div class="min-w-0 flex-1">
                      <div class="flex items-center gap-2 min-w-0">
                        <KindIcon v-if="row.icon" :config="row.icon" size="md" />
                        <span :class="row.labelClass" class="truncate min-w-0">
                          {{ row.label }}
                        </span>
                        <FieldKindTag v-if="row.tag" :type="row.tag.type" />
                        <i
                          v-if="row.children.length"
                          class="pi shrink-0 pl-1 text-xs text-muted-color"
                          :class="isOpen(row) ? 'pi-chevron-down' : 'pi-chevron-right'"
                        />
                      </div>
                      <!-- What the current-lens toggle actually grants on this row. -->
                      <p
                        v-if="laneDesc(row)"
                        class="text-xs text-muted-color leading-snug mt-0.5 truncate first-letter:uppercase"
                        :style="row.icon ? { paddingLeft: '2rem' } : {}"
                      >
                        {{ laneDesc(row) }}
                      </p>
                    </div>
                    <!-- Full permissions editor: same outlined secondary button as
                         CPermissionsButton elsewhere; hover-revealed to keep the
                         dense grid calm. -->
                    <Button
                      v-if="row.advancedResource"
                      icon="pi pi-lock"
                      severity="secondary"
                      outlined
                      size="small"
                      class="shrink-0 ml-auto opacity-0 focus:opacity-100 group-hover:opacity-100 transition-opacity"
                      :disabled="disabled"
                      v-tooltip.bottom="{
                        value: $t('project.permissions.advancedFor', { target: row.tipTarget }),
                        showDelay: 500,
                      }"
                      @click.stop="openAdvanced(row)"
                    />
                  </div>
                </div>

                <!-- Read-only evaluated cell (resolved access under the lens). -->
                <div
                  v-if="evalUserId"
                  class="flex items-center justify-center px-1"
                  :class="
                    isHeaderRow(row)
                      ? ''
                      : 'border-l border-surface bg-black/[0.04] dark:bg-white/[0.04]'
                  "
                >
                  <span
                    v-if="!isHeaderRow(row)"
                    class="inline-flex items-center justify-center w-7 h-7 cursor-default"
                    v-tooltip.bottom="{ value: evalTooltip(row), showDelay: 500 }"
                  >
                    <i :class="stateIcon(userCapState(laneOps(row)))" />
                  </span>
                </div>

                <!-- One toggle per role. Section-header rows render an empty,
                     borderless divider strip (no faint dashes). The whole cell is
                     the toggle target (button fills it edge-to-edge). -->
                <div
                  v-for="role in roles"
                  :key="role.id"
                  class="relative"
                  :class="isHeaderRow(row) ? '' : 'border-l border-surface'"
                >
                  <button
                    v-if="!isHeaderRow(row)"
                    type="button"
                    class="absolute inset-0 w-full h-full flex items-center justify-center transition-colors enabled:cursor-pointer enabled:hover:bg-black/10 dark:enabled:hover:bg-white/10 disabled:cursor-default"
                    :disabled="
                      disabled ||
                      !laneOps(row).length ||
                      capState(role.id, laneOps(row)) === 'loading'
                    "
                    v-tooltip.bottom="{ value: laneTooltip(row, role), showDelay: 500 }"
                    @click="toggleLane(row, role)"
                  >
                    <i :class="stateIcon(capState(role.id, laneOps(row)))" />
                  </button>
                </div>

                <!-- Filler column: absorbs slack so the roles stay left-adjacent. -->
                <div
                  aria-hidden="true"
                  :class="isHeaderRow(row) ? '' : 'border-l border-surface'"
                />
              </div>
            </div>
          </div>
        </TransitionGroup>
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
          :class="
            pending?.next === 'deny'
              ? 'pi pi-ban text-red-500'
              : 'pi pi-check-circle text-green-500'
          "
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
import KindIcon from '@/sections/project/components/KindIcon.vue'
import FieldKindTag from '@/sections/project/components/datamodel/FieldKindTag.vue'
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
  // Per-user (evaluation) mode for the step view: when set, a leading read-only
  // column shows this user's *resolved* effective access across all their roles
  // (traced by userID). Cells can't be toggled — access is edited via the role
  // columns beside it. `evalLabel` names the column header (falls back to a
  // generic "Evaluated" label).
  evalUserId: { type: String, default: null },
  evalLabel: { type: String, default: '' },
  // Single-resource mode: `{ kind, id }`. When set, the tree renders only that
  // one resource's own row + nested subtree (records/fields, layouts, …) with no
  // section-header row — for embedding in a resource detail dialog, roles still
  // the columns. Null = the full project matrix (all kinds).
  scope: { type: Object, default: null },
})

const store = useProjectsStore()
const { t } = useI18n()
const $toast = inject('$toast')
const { open: openPermissions } = usePermissions()

const ns = computed(() => props.project.namespaceID)
const soleRole = computed(() => props.roles[0] || null)

// The kinds an END USER of the deployed app interacts with. These access roles
// are for people using the app (pages + records), not building it — so only the
// runtime surface is listed: data (modules → records/fields), pages (view), the
// AI they can use (agents/chatbots), and automations they can RUN. Build-time ops
// (create/edit modules & pages, manage AI, edit automations) stay out; rare ops
// remain reachable via the ⚙ full editor. Connections are pure build-time infra
// (no end-user runtime op) and are intentionally absent.
const KINDS = [
  { kind: 'module', getter: 'resourcesFor', type: 'corteza::compose:module', ns: true },
  { kind: 'page', getter: 'pagesFor', type: 'corteza::compose:page', ns: true },
  { kind: 'agent', getter: 'agentsFor', type: 'corteza::system:agent', ns: false },
  { kind: 'chatbot', getter: 'chatbotsFor', type: 'corteza::system:chatbot', ns: false },
  {
    kind: 'automation',
    getter: 'automationsFor',
    type: 'corteza::automation:ng-automation',
    ns: false,
  },
]

const HEADER_BG = 'bg-emphasis'
const ITEM_BG = 'bg-surface'

// Step-view columns. Each row is its own grid (so rows can animate their height),
// so the columns need shared, fixed widths to line up down the matrix: a fixed
// title/description column, an optional evaluated column, fixed uniform role
// columns (long names truncate), then a filler that keeps the roles left-adjacent.
const gridStyle = computed(() => {
  const cols = ['24rem'] // title + description
  if (props.evalUserId) cols.push('8rem') // evaluated
  for (let i = 0; i < props.roles.length; i++) cols.push('7rem') // one per role
  cols.push('1fr') // filler
  return { gridTemplateColumns: cols.join(' ') }
})

function resStr(def, id) {
  return def.ns ? `${def.type}/${ns.value}/${id}` : `${def.type}/${id}`
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
    // Single-resource scope: only the matching kind, only the matching id.
    if (props.scope && def.kind !== props.scope.kind) continue
    const items = (store[def.getter](props.project.id) || []).filter(
      it => !props.scope || String(it.id) === String(props.scope.id),
    )
    if (!items.length) continue
    const cfg = kindConfig(def.kind)

    const section = {
      key: `k:${def.kind}`,
      label: t(cfg.labelKey),
      level: 0,
      tipTarget: t(cfg.labelKey),
      icon: cfg,
      rowBg: HEADER_BG,
      labelClass: 'font-normal',
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
              tag: { type: f.type },
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
          labelClass: 'font-medium',
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
        // Layouts nest beneath, each a single "view" toggle — a page only renders
        // once the role can read one of its layouts, so the page's View cascades
        // (capOpsDeep) to every layout read, completing the grant. Layout
        // edit/delete stays build-time (⚙ full editor only).
        const res = resStr(def, it.id)
        section.children.push({
          key: `page:${it.id}`,
          label: it.name,
          level: 1,
          tipTarget: it.name,
          icon: null,
          rowBg: ITEM_BG,
          labelClass: 'font-medium',
          defaultOpen: false,
          advancedResource: res,
          advancedTitle: it.name,
          descGroup: 'page',
          descTarget: it.name,
          capLabels: { read: 'project.permissions.caps.view' },
          caps: { read: [{ resource: res, op: 'read' }] },
          children: (it.layouts || []).map(l => {
            const layoutRes = `corteza::compose:page-layout/${nsID}/${it.id}/${l.id}`
            return {
              key: `layout:${it.id}:${l.id}`,
              label: l.name,
              level: 2,
              tipTarget: l.name,
              icon: null,
              rowBg: ITEM_BG,
              labelClass: '',
              defaultOpen: false,
              advancedResource: layoutRes,
              advancedTitle: l.name,
              descGroup: 'pageLayout',
              descTarget: l.name,
              capLabels: { read: 'project.permissions.caps.view' },
              caps: { read: [{ resource: layoutRes, op: 'read' }] },
              children: [],
            }
          }),
        })
        continue
      }

      if (def.kind === 'automation') {
        // Automations: a single "Run" toggle — end users trigger them at runtime,
        // they don't build them. Folded ops read + execute: `read` is enforced
        // (list/lookup the automation) so the app can surface it; `execute` is the
        // declared run capability (and the flag the UI uses). Edit/delete/undelete
        // stay build-time (⚙ full editor only). Uses the 'read' lane so it shows
        // under the Read lens, exactly like the agent/chatbot "Use" toggle.
        const res = resStr(def, it.id)
        section.children.push({
          key: `automation:${it.id}`,
          label: it.name,
          level: 1,
          tipTarget: it.name,
          icon: null,
          rowBg: ITEM_BG,
          labelClass: 'font-medium',
          defaultOpen: false,
          advancedResource: res,
          advancedTitle: it.name,
          descGroup: 'automation',
          descTarget: it.name,
          capLabels: { read: 'project.permissions.caps.run' },
          caps: {
            read: [
              { resource: res, op: 'read' },
              { resource: res, op: 'execute' },
            ],
          },
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
        labelClass: 'font-medium',
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

    // Scoped to one resource: drop the level-0 kind header, surface just the
    // resource's own subtree so the dialog shows only that resource's rows.
    if (props.scope) nodes.push(...section.children)
    else nodes.push(section)
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

// Row indentation: a 0.5rem base gutter plus a compact 0.7rem step per nesting
// level, measured from the resource level (level 1), NOT the section-header level
// (level 0). So a section header and its resource rows share the base gutter, and
// a resource row lands at the same indent whether it sits under a section (full
// matrix) or is the scoped root of a resource dialog (no section). Records/fields
// /layouts still step in beneath their resource.
function indentRem(row) {
  return 0.5 + Math.max(0, row.level - 1) * 0.7 + 'rem'
}

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
  [
    () => props.project.id,
    () => props.roles,
    allResources,
    () => store.graphVersion,
    () => props.evalUserId,
  ],
  () => {
    if (!props.project.id) return
    store.loadEffectiveAccess(props.project.id, allResources.value)
    if (props.evalUserId)
      store.loadUserEffectiveAccess(props.project.id, props.evalUserId, allResources.value)
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

// Plain-language summary of a row for the single-role detail view (which has no
// lens — all caps show at once). Uses the row's most representative capability,
// preferring 'read' (the "see it" base) so e.g. a Records row reads "view records
// in …" while its create/update/delete chips sit alongside. Empty for headers.
const ROW_DESC_ORDER = ['read', 'create', 'update', 'delete']
function rowDesc(row) {
  const cap = ROW_DESC_ORDER.find(c => (row.caps[c] || []).length)
  return cap ? capDesc(row, cap) : ''
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

// Same compound-state logic as capState, but for the evaluated user column —
// reads the user's resolved effective access (across all their roles) instead
// of a single role's.
function userCapState(ops) {
  if (!ops || !ops.length) return 'na'
  const states = ops.map(({ resource, op }) =>
    store.userEffectiveAccess(props.project.id, props.evalUserId, resource, op),
  )
  if (
    states.some(s => !s) &&
    store.isUserEffectiveAccessLoading(props.project.id, props.evalUserId)
  )
    return 'loading'
  const defined = states.filter(Boolean)
  if (!defined.length) return 'none'
  const uniq = [...new Set(defined)]
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

// Evaluated column: read-only tooltip for the user's resolved access (never
// actionable — the last arg is false, so no "click to…" hint is appended).
function evalTooltip(row) {
  const ops = laneOps(row)
  if (!ops.length) return t('project.permissions.access.naSimple', { desc: row.tipTarget })
  return tooltipFor(userCapState(ops), capDesc(row, capability.value), false)
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
/* Row expand/collapse. Each row animates its OWN height via grid-template-rows
   0fr<->1fr (the same trick the module card uses) plus a fade, so rows grow and
   shrink smoothly instead of popping. Because the height animates in normal flow,
   the rows below just follow along — no FLIP/reflow jump, in either direction. */
.rowcollapse-item {
  display: grid;
  grid-template-rows: 1fr;
}
.rowcollapse-enter-active,
.rowcollapse-leave-active {
  transition:
    grid-template-rows 0.2s ease,
    opacity 0.2s ease;
}
.rowcollapse-enter-from,
.rowcollapse-leave-to {
  grid-template-rows: 0fr;
  opacity: 0;
}
@media (prefers-reduced-motion: reduce) {
  .rowcollapse-enter-active,
  .rowcollapse-leave-active {
    transition: none;
  }
}
</style>
