<template>
  <!-- Real event log for the published-project dashboard. Ports the admin
       ActionLog list (filters, cursor pagination, drill-down, expansion rows)
       but wears the CategoryView layout idiom: a fixed badge/title header, a
       fixed filter row, then the list fills the rest and scrolls internally.
       DashboardLayout owns the topbar, so there is NO Teleport here. -->
  <div class="flex flex-col h-full min-h-0 min-w-0 overflow-hidden">
    <!-- Title bar — wizard-style leading badge + title + description. -->
    <header class="shrink-0 border-b border-surface px-4 py-3 flex items-center gap-3">
      <span
        class="inline-flex items-center justify-center w-9 h-9 rounded-md ring-1 shrink-0 bg-emphasis ring-surface"
      >
        <i class="pi pi-list text-primary" />
      </span>
      <div class="min-w-0">
        <h2 class="text-xl font-semibold text-color truncate">
          {{ $t('project.dashboard.allEvents.title') }}
        </h2>
        <p class="text-sm text-muted-color">{{ $t('project.dashboard.allEvents.desc') }}</p>
      </div>
    </header>

    <!-- Filters — fixed above the list; same grid as the admin page. -->
    <div class="shrink-0 p-4">
      <!-- Each control is wrapped in its <label> so the caption is
           programmatically associated with the input (a11y) without per-component
           inputId juggling (CInputUser doesn't forward one). -->
      <div class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 xl:grid-cols-6 gap-3 w-full">
        <label class="flex flex-col gap-1">
          <span class="text-sm font-medium text-muted-color">
            {{ $t('system.actionlog.list.filter.from') }}
          </span>
          <DatePicker
            v-model="filter.from"
            showTime
            hourFormat="24"
            showButtonBar
            size="small"
            fluid
            @update:modelValue="reload"
          />
        </label>

        <label class="flex flex-col gap-1">
          <span class="text-sm font-medium text-muted-color">
            {{ $t('system.actionlog.list.filter.to') }}
          </span>
          <DatePicker
            v-model="filter.to"
            showTime
            hourFormat="24"
            showButtonBar
            size="small"
            fluid
            @update:modelValue="reload"
          />
        </label>

        <label class="flex flex-col gap-1">
          <span class="text-sm font-medium text-muted-color">
            {{ $t('system.actionlog.list.filter.actor') }}
          </span>
          <CInputUser v-model="filter.actorID" size="small" @update:modelValue="reload" />
        </label>

        <label class="flex flex-col gap-1">
          <span class="text-sm font-medium text-muted-color">
            {{ $t('system.actionlog.list.filter.origin') }}
          </span>
          <Select
            v-model="filter.origin"
            :options="originOptions"
            option-label="label"
            option-value="value"
            filter
            show-clear
            size="small"
            @update:modelValue="reload"
          />
        </label>

        <label class="flex flex-col gap-1">
          <span class="text-sm font-medium text-muted-color">
            {{ $t('system.actionlog.list.filter.resource') }}
          </span>
          <Select
            v-model="filter.resource"
            :options="resourceOptions"
            option-label="label"
            option-value="value"
            filter
            show-clear
            size="small"
            @update:modelValue="reload"
          />
        </label>

        <label class="flex flex-col gap-1">
          <span class="text-sm font-medium text-muted-color">
            {{ $t('system.actionlog.list.filter.action') }}
          </span>
          <Select
            v-model="filter.action"
            :options="actionOptions"
            option-label="label"
            option-value="value"
            filter
            show-clear
            size="small"
            @update:modelValue="reload"
          />
        </label>
      </div>
    </div>

    <!-- The log — fills remaining space; CResourceList owns its internal scroll. -->
    <div class="flex-1 min-h-0 px-4 pb-4">
      <CResourceList
        class="h-full"
        primary-key="actionID"
        :fields="fields"
        :items="items"
        :filter="{}"
        :sorting="{}"
        :pagination="{}"
        :loading="loading && items.length === 0"
        :expandable="true"
        hide-search
        hide-pagination
        :translations="{ noItems: $t('project.dashboard.allEvents.empty') }"
      >
        <template #body-timestamp="{ data }">
          {{ locFullDateTime(data.timestamp) }}
        </template>

        <template #body-actor="{ data }">
          <span
            v-if="data.actor || data.actorID"
            v-tooltip.top="data.actorID"
            role="button"
            tabindex="0"
            :class="filterLinkClass(filter.actorID, data.actorID)"
            @click.stop="drillDown('actorID', data.actorID)"
            @keydown.enter.prevent="drillDown('actorID', data.actorID)"
            @keydown.space.prevent="drillDown('actorID', data.actorID)"
          >
            {{ actorLabel(data) }}
          </span>
        </template>

        <template #body-resource="{ data }">
          <span
            v-if="data.resource"
            v-tooltip.top="data.resource"
            role="button"
            tabindex="0"
            :class="filterLinkClass(filter.resource, data.resource)"
            @click.stop="drillDown('resource', data.resource)"
            @keydown.enter.prevent="drillDown('resource', data.resource)"
            @keydown.space.prevent="drillDown('resource', data.resource)"
          >
            {{ resourceLabel(data.resource) }}
          </span>
        </template>

        <template #body-action="{ data }">
          <span
            v-if="data.action"
            v-tooltip.top="data.action"
            role="button"
            tabindex="0"
            :class="filterLinkClass(filter.action, data.action)"
            @click.stop="drillDown('action', data.action)"
            @keydown.enter.prevent="drillDown('action', data.action)"
            @keydown.space.prevent="drillDown('action', data.action)"
          >
            {{ actionLabel(data.action) }}
          </span>
        </template>

        <template #body-requestOrigin="{ data }">
          <span v-if="data.requestOrigin" v-tooltip.top="data.requestOrigin">
            {{ originLabel(data.requestOrigin) }}
          </span>
        </template>

        <template #body-severity="{ data }">
          <Tag
            :value="
              $t('system.actionlog.list.severity.' + (severityMap[data.severity]?.label ?? 'info'))
            "
            :severity="severityMap[data.severity]?.severity ?? 'info'"
          />
        </template>

        <template v-if="items.length" #footer>
          <div class="flex justify-center px-3 py-2">
            <Button
              :label="$t('system.actionlog.list.loadOlder')"
              :loading="loading"
              severity="secondary"
              size="small"
              @click="loadOlder"
            />
          </div>
        </template>

        <template #expansion="{ data }">
          <div class="flex flex-wrap gap-4 p-4">
            <div class="flex-1 min-w-64">
              <p class="font-semibold mb-2">{{ $t('system.actionlog.list.details.header') }}</p>
              <table class="text-sm w-full">
                <tbody>
                  <tr>
                    <td class="text-muted-color pr-4 py-0.5 whitespace-nowrap">
                      {{ $t('system.actionlog.list.details.id') }}
                    </td>
                    <td class="py-0.5">{{ data.actionID }}</td>
                  </tr>
                  <tr>
                    <td class="text-muted-color pr-4 py-0.5 whitespace-nowrap">
                      {{ $t('system.actionlog.list.details.timestamp') }}
                    </td>
                    <td class="py-0.5">{{ locFullDateTime(data.timestamp) }}</td>
                  </tr>
                  <tr>
                    <td class="text-muted-color pr-4 py-0.5 whitespace-nowrap">
                      {{ $t('system.actionlog.list.details.requestOrigin') }}
                    </td>
                    <td class="py-0.5" v-tooltip.top="data.requestOrigin">
                      {{ originLabel(data.requestOrigin) }}
                    </td>
                  </tr>
                  <tr>
                    <td class="text-muted-color pr-4 py-0.5 whitespace-nowrap">
                      {{ $t('system.actionlog.list.details.requestID') }}
                    </td>
                    <td class="py-0.5">{{ data.requestID }}</td>
                  </tr>
                  <tr>
                    <td class="text-muted-color pr-4 py-0.5 whitespace-nowrap">
                      {{ $t('system.actionlog.list.details.actorIPAddr') }}
                    </td>
                    <td class="py-0.5">{{ data.actorIPAddr }}</td>
                  </tr>
                  <tr>
                    <td class="text-muted-color pr-4 py-0.5 whitespace-nowrap">
                      {{ $t('system.actionlog.list.details.actor') }}
                    </td>
                    <td class="py-0.5">{{ actorLabel(data) }}</td>
                  </tr>
                  <tr>
                    <td class="text-muted-color pr-4 py-0.5 whitespace-nowrap">
                      {{ $t('system.actionlog.list.details.actorID') }}
                    </td>
                    <td class="py-0.5">{{ data.actorID }}</td>
                  </tr>
                  <tr>
                    <td class="text-muted-color pr-4 py-0.5 whitespace-nowrap">
                      {{ $t('system.actionlog.list.details.severity') }}
                    </td>
                    <td class="py-0.5">
                      {{
                        $t(
                          'system.actionlog.list.severity.' +
                            (severityMap[data.severity]?.label ?? 'info'),
                        )
                      }}
                    </td>
                  </tr>
                  <tr>
                    <td class="text-muted-color pr-4 py-0.5 whitespace-nowrap">
                      {{ $t('system.actionlog.list.details.resource') }}
                    </td>
                    <td class="py-0.5" v-tooltip.top="data.resource">
                      {{ resourceLabel(data.resource) }}
                    </td>
                  </tr>
                  <tr>
                    <td class="text-muted-color pr-4 py-0.5 whitespace-nowrap">
                      {{ $t('system.actionlog.list.details.action') }}
                    </td>
                    <td class="py-0.5" v-tooltip.top="data.action">
                      {{ actionLabel(data.action) }}
                    </td>
                  </tr>
                </tbody>
              </table>
            </div>

            <div class="flex-1 min-w-64">
              <p class="font-semibold mb-2">
                {{ $t('system.actionlog.list.details.headerAdditional') }}
              </p>
              <table class="text-sm w-full">
                <tbody>
                  <tr>
                    <td class="text-muted-color pr-4 py-0.5 whitespace-nowrap">
                      {{ $t('system.actionlog.list.details.description') }}
                    </td>
                    <td class="py-0.5">{{ data.description }}</td>
                  </tr>
                  <tr v-if="data.error">
                    <td class="text-muted-color pr-4 py-0.5 whitespace-nowrap">
                      {{ $t('system.actionlog.list.details.error') }}
                    </td>
                    <td class="py-0.5 text-red-500">{{ data.error }}</td>
                  </tr>
                </tbody>
              </table>

              <template v-if="data.meta && Object.keys(data.meta).length">
                <Divider />
                <p class="font-semibold mb-2">{{ $t('system.actionlog.list.details.meta') }}</p>
                <table class="text-sm w-full font-mono">
                  <tbody>
                    <tr v-for="(val, key) in data.meta" :key="key">
                      <td class="text-muted-color pr-4 py-0.5 whitespace-nowrap">{{ key }}</td>
                      <td class="py-0.5">{{ val }}</td>
                    </tr>
                  </tbody>
                </table>
              </template>
            </div>
          </div>
        </template>
      </CResourceList>
    </div>
  </div>
</template>

<script setup>
import { inject, reactive, ref, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRoute } from 'vue-router'
import { components, filters, useUserStore } from '@planetcrust/human-vue'
// Cross-section reuse: the ActionLog vocab (resource/action/origin option lists
// and label helpers) is the single source of truth — do NOT duplicate it here.
import {
  RESOURCE_TYPES,
  COMMON_ACTIONS,
  ORIGINS,
  SEVERITY_MAP,
  actionLabel,
  resourceLabel,
  originLabel,
} from '@/sections/admin/views/system/ActionLog/vocab'

const { CResourceList, CInputUser } = components
const { locFullDateTime } = filters

const { t } = useI18n()
const route = useRoute()
const $SystemAPI = inject('$SystemAPI')

// Shared user store: batch-resolve actor IDs and look up labels (no bespoke cache).
const userStore = useUserStore()

// The project this dashboard is for; scopes the log to its events (see
// buildParams). The backend tags each event with the active project from the
// request scope, so every resource touched in the project (compose, agents,
// chatbots, roles, users…) is included.
const projectID = computed(() => route.params.projectId || undefined)

const items = ref([])
const loading = ref(false)

const filter = reactive({
  from: null,
  to: null,
  resource: '',
  action: '',
  origin: '',
  actorID: '',
})

// Severity (uint8 0–7) → Tag severity + label key; shared with the admin log.
const severityMap = SEVERITY_MAP

const fields = [
  {
    key: 'timestamp',
    sortable: false,
    header: t('system.actionlog.list.columns.timestamp'),
  },
  {
    key: 'actor',
    sortable: false,
    header: t('system.actionlog.list.columns.actor'),
  },
  {
    key: 'requestOrigin',
    sortable: false,
    header: t('system.actionlog.list.columns.requestOrigin'),
  },
  {
    key: 'resource',
    sortable: false,
    header: t('system.actionlog.list.columns.resource'),
  },
  {
    key: 'action',
    sortable: false,
    header: t('system.actionlog.list.columns.action'),
  },
  {
    key: 'description',
    sortable: false,
    header: t('system.actionlog.list.columns.description'),
  },
  {
    key: 'severity',
    sortable: false,
    header: t('system.actionlog.list.columns.severity'),
    class: 'text-right',
    pt: { columnHeaderContent: 'justify-end' },
  },
]

// Authoritative option lists from server enums; actions also merge any unusual
// values present in loaded items so service-specific names stay discoverable.
// If user drills down by clicking a row, filter.resource may carry an ID
// suffix that's not in RESOURCE_TYPES. Append it as a virtual option so the
// Select can render it instead of going blank with only the clear button.
const resourceOptions = computed(() => {
  const v = filter.resource
  if (!v) return RESOURCE_TYPES
  if (RESOURCE_TYPES.some(o => o.value === v)) return RESOURCE_TYPES
  return [...RESOURCE_TYPES, { value: v, label: resourceLabel(v) }]
})

const originOptions = ORIGINS

const actionOptions = computed(() => {
  const known = new Set(COMMON_ACTIONS.map(o => o.value))
  const extras = [...new Set(items.value.map(i => i.action).filter(v => v && !known.has(v)))].map(
    value => ({ value, label: actionLabel(value) }),
  )
  return [...COMMON_ACTIONS, ...extras].sort((a, b) => a.label.localeCompare(b.label))
})

// Click a cell value to filter by it; click the active filter value to clear it
function drillDown(field, value) {
  filter[field] = filter[field] === value ? '' : value
  reload()
}

function filterLinkClass(currentFilterValue, cellValue) {
  const isActive = currentFilterValue === cellValue
  return [
    'cursor-pointer hover:underline',
    isActive ? 'text-primary font-medium' : 'text-primary-500',
  ]
}

function buildParams(beforeActionID) {
  return {
    from: filter.from ? filter.from.toISOString() : undefined,
    to: filter.to ? filter.to.toISOString() : undefined,
    resource: filter.resource || undefined,
    // Scope to this project so the log shows only its events.
    projectID: projectID.value || undefined,
    action: filter.action || undefined,
    origin: filter.origin || undefined,
    // actorID can be null when cleared via Select's clear button
    actorID: filter.actorID ? [filter.actorID] : undefined,
    beforeActionID: beforeActionID || undefined,
    limit: 50,
  }
}

// Sequence token: discard responses from filter changes that have been superseded.
let loadSeq = 0

async function load(reset = false) {
  // Allow reset to bypass the in-flight guard (filter change should always reload).
  if (loading.value && !reset) return

  if (reset) {
    items.value = []
  }

  const beforeActionID = reset
    ? undefined
    : items.value.length > 0
      ? items.value[items.value.length - 1].actionID
      : undefined

  loading.value = true
  const mySeq = ++loadSeq
  try {
    const { set } = await $SystemAPI.actionlogList(buildParams(beforeActionID))
    if (mySeq !== loadSeq) return // stale response
    if (reset) {
      items.value = set ?? []
    } else {
      items.value = [...items.value, ...(set ?? [])]
    }
    resolveActors(set ?? [])
  } finally {
    if (mySeq === loadSeq) loading.value = false
  }
}

function reload() {
  load(true)
}

// Actor resolution via the shared user store — batch-resolve encountered IDs.
function resolveActors(rows) {
  const ids = rows.map(r => r.actorID).filter(id => id && id !== '0')
  // Fire-and-forget; resolveUsers re-throws on API failure, so swallow it — a
  // failed name lookup must not surface as an unhandled rejection.
  userStore.resolveUsers(ids).catch(() => {})
}

function actorLabel(data) {
  const user = data.actorID && userStore.findByID(data.actorID)
  if (user) return user.name || user.handle || user.username || user.email || user.userID
  return data.actor || data.actorID || ''
}

function loadOlder() {
  load(false)
}

onMounted(() => load(true))
</script>
