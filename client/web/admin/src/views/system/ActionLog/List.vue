<template>
  <Teleport to="#topbar-title" defer>
    <span>{{ $t('system.actionlog.list.title') }}</span>
  </Teleport>

  <div class="p-4 h-full overflow-hidden min-w-0 flex flex-col gap-3">
    <CResourceList
      primary-key="actionID"
      :fields="fields"
      :items="items"
      :filter="filter"
      :sorting="{}"
      :pagination="{}"
      :loading="loading && items.length === 0"
      :expandable="true"
      hide-search
      hide-pagination
      class="flex-1 min-h-0"
    >
      <template #header>
        <div class="flex flex-wrap items-center gap-3 w-full">
          <div class="flex items-center gap-2">
            <label class="text-sm font-medium text-muted-color">
              {{ $t('system.actionlog.list.filter.from') }}
            </label>
            <DatePicker
              v-model="filter.from"
              showTime
              hourFormat="24"
              size="small"
              @update:modelValue="reload"
            />
            <label class="text-sm font-medium text-muted-color">
              {{ $t('system.actionlog.list.filter.to') }}
            </label>
            <DatePicker
              v-model="filter.to"
              showTime
              hourFormat="24"
              size="small"
              @update:modelValue="reload"
            />
          </div>

          <div class="flex items-center gap-2">
            <label class="text-sm font-medium text-muted-color">
              {{ $t('system.actionlog.list.filter.resource') }}
            </label>
            <Select
              v-model="filter.resource"
              :options="resourceOptions"
              show-clear
              size="small"
              class="min-w-40"
              @update:modelValue="reload"
            />
          </div>

          <div class="flex items-center gap-2">
            <label class="text-sm font-medium text-muted-color">
              {{ $t('system.actionlog.list.filter.action') }}
            </label>
            <Select
              v-model="filter.action"
              :options="actionOptions"
              show-clear
              size="small"
              class="min-w-40"
              @update:modelValue="reload"
            />
          </div>

          <div class="flex items-center gap-2">
            <label class="text-sm font-medium text-muted-color">
              {{ $t('system.actionlog.list.filter.actor') }}
            </label>
            <Select
              v-model="filter.actorID"
              :options="actorOptions"
              option-label="label"
              option-value="value"
              show-clear
              size="small"
              class="min-w-40"
              @update:modelValue="reload"
            />
          </div>
        </div>
      </template>

      <template #body-timestamp="{ data }">
        {{ locFullDateTime(data.timestamp) }}
      </template>

      <template #body-actor="{ data }">
        <span
          v-if="data.actor || data.actorID"
          :class="filterLinkClass(filter.actorID, data.actorID)"
          @click.stop="drillDown('actorID', data.actorID)"
        >{{ data.actor || data.actorID }}</span>
      </template>

      <template #body-resource="{ data }">
        <span
          v-if="data.resource"
          :class="filterLinkClass(filter.resource, data.resource)"
          @click.stop="drillDown('resource', data.resource)"
        >{{ data.resource }}</span>
      </template>

      <template #body-action="{ data }">
        <span
          v-if="data.action"
          :class="filterLinkClass(filter.action, data.action)"
          @click.stop="drillDown('action', data.action)"
        >{{ data.action }}</span>
      </template>

      <template #body-severity="{ data }">
        <Tag
          :value="$t('system.actionlog.list.severity.' + (severityMap[data.severity]?.label ?? 'info'))"
          :severity="severityMap[data.severity]?.severity ?? 'info'"
        />
      </template>

      <template #expansion="{ data }">
        <div class="flex flex-wrap gap-4 p-4">
          <div class="flex-1 min-w-64">
            <p class="font-semibold mb-2">{{ $t('system.actionlog.list.details.header') }}</p>
            <table class="text-sm w-full">
              <tbody>
                <tr>
                  <td class="text-muted-color pr-4 py-0.5 whitespace-nowrap">{{ $t('system.actionlog.list.details.id') }}</td>
                  <td class="py-0.5">{{ data.actionID }}</td>
                </tr>
                <tr>
                  <td class="text-muted-color pr-4 py-0.5 whitespace-nowrap">{{ $t('system.actionlog.list.details.timestamp') }}</td>
                  <td class="py-0.5">{{ locFullDateTime(data.timestamp) }}</td>
                </tr>
                <tr>
                  <td class="text-muted-color pr-4 py-0.5 whitespace-nowrap">{{ $t('system.actionlog.list.details.requestOrigin') }}</td>
                  <td class="py-0.5">{{ data.requestOrigin }}</td>
                </tr>
                <tr>
                  <td class="text-muted-color pr-4 py-0.5 whitespace-nowrap">{{ $t('system.actionlog.list.details.requestID') }}</td>
                  <td class="py-0.5">{{ data.requestID }}</td>
                </tr>
                <tr>
                  <td class="text-muted-color pr-4 py-0.5 whitespace-nowrap">{{ $t('system.actionlog.list.details.actorIPAddr') }}</td>
                  <td class="py-0.5">{{ data.actorIPAddr }}</td>
                </tr>
                <tr>
                  <td class="text-muted-color pr-4 py-0.5 whitespace-nowrap">{{ $t('system.actionlog.list.details.actor') }}</td>
                  <td class="py-0.5">{{ data.actor }}</td>
                </tr>
                <tr>
                  <td class="text-muted-color pr-4 py-0.5 whitespace-nowrap">{{ $t('system.actionlog.list.details.actorID') }}</td>
                  <td class="py-0.5">{{ data.actorID }}</td>
                </tr>
                <tr>
                  <td class="text-muted-color pr-4 py-0.5 whitespace-nowrap">{{ $t('system.actionlog.list.details.severity') }}</td>
                  <td class="py-0.5">{{ $t('system.actionlog.list.severity.' + (severityMap[data.severity]?.label ?? 'info')) }}</td>
                </tr>
                <tr>
                  <td class="text-muted-color pr-4 py-0.5 whitespace-nowrap">{{ $t('system.actionlog.list.details.resource') }}</td>
                  <td class="py-0.5">{{ data.resource }}</td>
                </tr>
                <tr>
                  <td class="text-muted-color pr-4 py-0.5 whitespace-nowrap">{{ $t('system.actionlog.list.details.action') }}</td>
                  <td class="py-0.5">{{ data.action }}</td>
                </tr>
              </tbody>
            </table>
          </div>

          <div class="flex-1 min-w-64">
            <p class="font-semibold mb-2">{{ $t('system.actionlog.list.details.headerAdditional') }}</p>
            <table class="text-sm w-full">
              <tbody>
                <tr>
                  <td class="text-muted-color pr-4 py-0.5 whitespace-nowrap">{{ $t('system.actionlog.list.details.description') }}</td>
                  <td class="py-0.5">{{ data.description }}</td>
                </tr>
                <tr v-if="data.error">
                  <td class="text-muted-color pr-4 py-0.5 whitespace-nowrap">{{ $t('system.actionlog.list.details.error') }}</td>
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

    <div v-if="items.length" class="flex justify-center shrink-0">
      <Button
        :label="$t('system.actionlog.list.loadOlder')"
        :loading="loading"
        severity="secondary"
        size="small"
        @click="loadOlder"
      />
    </div>
  </div>
</template>

<script setup>
import { inject, reactive, ref, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { components, filters } from '@planetcrust/human-vue'

const { CResourceList } = components
const { locFullDateTime } = filters

const { t } = useI18n()
const $SystemAPI = inject('$SystemAPI')

const items = ref([])
const loading = ref(false)

const filter = reactive({
  from: null,
  to: null,
  resource: '',
  action: '',
  actorID: '',
})

// Severity is returned as uint8 (0–7) from the backend.
// 0=Emergency, 1=Alert, 2=Critical, 3=Error, 4=Warning, 5=Notice, 6=Info, 7=Debug
const severityMap = {
  0: { severity: 'danger', label: 'emergency' },
  1: { severity: 'danger', label: 'alert' },
  2: { severity: 'danger', label: 'critical' },
  3: { severity: 'danger', label: 'error' },
  4: { severity: 'warn', label: 'warning' },
  5: { severity: 'success', label: 'notice' },
  6: { severity: 'info', label: 'info' },
  7: { severity: 'secondary', label: 'debug' },
}

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

// Unique option lists derived from all loaded items
const resourceOptions = computed(() =>
  [...new Set(items.value.map(i => i.resource).filter(Boolean))].sort(),
)

const actionOptions = computed(() =>
  [...new Set(items.value.map(i => i.action).filter(Boolean))].sort(),
)

const actorOptions = computed(() => {
  const seen = new Set()
  return items.value
    .filter(i => i.actorID && i.actorID !== '0')
    .reduce((acc, i) => {
      if (!seen.has(i.actorID)) {
        seen.add(i.actorID)
        acc.push({ value: i.actorID, label: i.actor || i.actorID })
      }
      return acc
    }, [])
    .sort((a, b) => a.label.localeCompare(b.label))
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
    action: filter.action || undefined,
    // actorID can be null when cleared via Select's clear button
    actorID: filter.actorID ? [filter.actorID] : undefined,
    beforeActionID: beforeActionID || undefined,
    limit: 50,
  }
}

async function load(reset = false) {
  if (loading.value) return

  const beforeActionID = reset
    ? undefined
    : items.value.length > 0
      ? items.value[items.value.length - 1].actionID
      : undefined

  loading.value = true
  try {
    const { set } = await $SystemAPI.actionlogList(buildParams(beforeActionID))
    if (reset) {
      items.value = set ?? []
    } else {
      items.value = [...items.value, ...(set ?? [])]
    }
  } finally {
    loading.value = false
  }
}

function reload() {
  load(true)
}

function loadOlder() {
  load(false)
}

onMounted(() => load(true))
</script>
