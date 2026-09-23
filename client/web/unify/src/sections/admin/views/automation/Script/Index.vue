<template>
  <Teleport to="#topbar-title" defer>
    <span>{{ $t('automation.scripts.list.title') }}</span>
  </Teleport>

  <CViewContainer>
    <div class="flex flex-col gap-4 h-full min-h-0">
      <!-- Corredor connection, as the script list reported it -->
      <Message v-if="banner" :severity="banner.severity" :closable="false" class="shrink-0">
        <div class="flex items-center justify-between gap-3 flex-wrap">
          <span>{{ bannerText }}</span>
          <Button
            v-if="banner.state === 'connected'"
            :label="$t('automation.scripts.list.status.refresh')"
            :title="$t('automation.scripts.list.status.refreshTooltip')"
            icon="pi pi-refresh"
            size="small"
            severity="secondary"
            :loading="refreshing"
            @click="refreshScripts"
          />
        </div>
      </Message>

      <CResourceList
        primary-key="name"
        :fields="fields"
        :items="paged"
        :filter="filter"
        :sorting="sorting"
        :pagination="paginationState"
        :loading="loading"
        :translations="{
          searchPlaceholder: $t('automation.scripts.list.filterForm.query.placeholder'),
          showingPagination: 'general.resourceList.pagination.showing',
          singlePluralPagination: 'general.resourceList.pagination.single',
          prevPagination: $t('general.resourceList.pagination.prev'),
          nextPagination: $t('general.resourceList.pagination.next'),
          recordsPerPage: $t('general.resourceList.pagination.recordsPerPage'),
          resourceSingle: $t('automation.scripts.list.resource.single'),
          resourcePlural: $t('automation.scripts.list.resource.plural'),
        }"
        class="flex-1 min-h-0"
        @update:filter="Object.assign(filter, $event)"
        @sort="handleSort"
        @page-change="handlePageChange"
      >
        <template #filter>
          <Button
            v-tooltip.bottom="$t('automation.scripts.list.filterForm.title')"
            :aria-label="$t('automation.scripts.list.filterForm.title')"
            icon="pi pi-filter"
            severity="secondary"
            size="small"
            text
            @click="toggleFilterMenu"
          />
        </template>

        <template #body-name="{ data }">
          <div class="flex flex-col gap-1 max-w-sm">
            <span v-if="data.label" class="font-medium truncate">{{ data.label }}</span>
            <span v-else class="text-muted-color italic">
              {{ $t('automation.scripts.list.labelMissing') }}
            </span>

            <span v-if="data.description" class="text-xs text-muted-color truncate">
              {{ data.description }}
            </span>

            <code class="text-xs text-muted-color truncate">{{ data.name }}</code>

            <div v-if="data.errors && data.errors.length" class="flex flex-col gap-1 mt-1">
              <Message
                v-for="(error, i) in data.errors"
                :key="i"
                severity="warn"
                :closable="false"
                class="text-sm"
              >
                {{ error }}
              </Message>
            </div>
          </div>
        </template>

        <template #body-extension="{ data }">
          <span v-if="scriptExtension(data)" class="flex items-center gap-2">
            <i class="pi pi-folder text-muted-color" />
            {{ scriptExtension(data) }}
          </span>
          <span v-else class="text-muted-color">
            {{ $t('automation.scripts.list.extensionNone') }}
          </span>
        </template>

        <template #body-kind="{ data }">
          <Tag :value="kindLabel(data)" severity="contrast" class="text-xs" />
        </template>

        <template #body-triggers="{ data }">
          <div class="flex flex-col gap-1 max-w-md">
            <div
              v-for="(trigger, i) in triggerRows(data)"
              :key="`trigger-${i}`"
              class="flex items-center gap-1 flex-wrap"
            >
              <Tag :value="trigger.label" severity="info" class="text-xs" />
              <Tag
                v-for="(constraint, j) in trigger.constraints"
                :key="`constraint-${j}`"
                :value="constraint"
                severity="secondary"
                class="text-xs"
              />
            </div>

            <div v-if="data.iterator" class="flex items-center gap-1 flex-wrap">
              <Tag :value="iteratorChip(data.iterator)" severity="warn" class="text-xs" />
              <Tag
                v-for="(constraint, j) in iteratorFilterChips(data.iterator)"
                :key="`iterator-filter-${j}`"
                :value="constraint"
                severity="secondary"
                class="text-xs"
              />
            </div>

            <div v-if="data.security" class="flex items-center gap-1 flex-wrap">
              <Tag
                v-for="(chip, i) in securityChips(data.security)"
                :key="`security-${i}`"
                :value="chip"
                severity="contrast"
                class="text-xs"
              />
            </div>
          </div>
        </template>

        <template #body-changedAt="{ data }">
          {{ changedAtText(data) }}
        </template>
      </CResourceList>

      <Popover ref="filterMenu">
        <div class="flex flex-col gap-4 p-2 w-72">
          <div class="flex flex-col gap-2">
            <span class="font-medium text-sm text-primary">
              {{ $t('automation.scripts.list.filterForm.has.label') }}
            </span>
            <div
              v-for="option in capabilityOptions"
              :key="option.key"
              class="flex items-center gap-2"
            >
              <Checkbox v-model="filter[option.key]" :inputId="`filter-${option.key}`" binary />
              <label :for="`filter-${option.key}`" class="text-sm cursor-pointer">
                {{ option.label }}
              </label>
            </div>
          </div>

          <div class="flex flex-col gap-2">
            <span class="font-medium text-sm text-primary">
              {{ $t('automation.scripts.list.filterForm.kind.label') }}
            </span>
            <div v-for="option in kindOptions" :key="option.key" class="flex items-center gap-2">
              <Checkbox v-model="filter[option.key]" :inputId="`filter-${option.key}`" binary />
              <label :for="`filter-${option.key}`" class="text-sm cursor-pointer">
                {{ option.label }}
              </label>
            </div>
          </div>

          <div class="flex flex-col gap-2">
            <span class="font-medium text-sm text-primary">
              {{ $t('automation.scripts.list.filterForm.extension.label') }}
            </span>
            <div
              v-for="option in extensionOptions"
              :key="option.id"
              class="flex items-center gap-2"
            >
              <RadioButton
                v-model="filter.extension"
                :inputId="`filter-extension-${option.id}`"
                :value="option.value"
              />
              <label :for="`filter-extension-${option.id}`" class="text-sm cursor-pointer">
                {{ option.label }}
              </label>
            </div>
          </div>
        </div>
      </Popover>
    </div>
  </CViewContainer>
</template>

<script setup>
import { computed, inject, onMounted, reactive, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { changedAtField, changedAtText, components, constraintChips } from '@planetcrust/human-vue'
import {
  ANY_EXTENSION,
  corredorBanner,
  groupScripts,
  matchesExtensionFilter,
  matchesKindFilter,
  relativeTime,
  scriptBundle,
  scriptExtension,
  scriptKind,
  sortScripts,
  triggerRows,
} from './script-inventory'

const { CResourceList, CViewContainer } = components
const { t, locale } = useI18n()

const $toast = inject('$toast')
const $SystemAPI = inject('$SystemAPI')

const loading = ref(false)
const refreshing = ref(false)
const items = ref([])
const status = ref({})

const filterMenu = ref()

const filter = reactive({
  query: '',
  extension: ANY_EXTENSION,
  incScriptsWithErrors: false,
  incScriptsWithTriggers: false,
  incScriptsWithIterator: false,
  incScriptsWithSecurity: false,
  incServerScripts: false,
  incClientScripts: false,
})

const sorting = reactive({ sortBy: 'name', sortDesc: false })
const pagination = reactive({ page: 1, limit: 50 })

// A row is as tall as its chips, and a cell centred against that floats away
// from the script it belongs to; every cell starts at the top of its row.
const topAligned = { bodyCell: { class: 'align-top' } }

const fields = [
  {
    key: 'name',
    header: t('automation.scripts.list.columns.name'),
    sortable: true,
    style: 'min-width: 16rem',
    pt: topAligned,
  },
  {
    key: 'extension',
    header: t('automation.scripts.list.columns.extension'),
    sortable: true,
    style: 'width: 10rem',
    pt: topAligned,
  },
  {
    key: 'kind',
    header: t('automation.scripts.list.columns.kind'),
    sortable: true,
    style: 'width: 9rem',
    pt: topAligned,
  },
  {
    key: 'triggers',
    header: t('automation.scripts.list.columns.triggers'),
    sortable: false,
    style: 'min-width: 16rem',
    pt: topAligned,
  },
  changedAtField(t('general.columns.changedAt'), {
    style: 'width: 10rem',
    pt: { columnHeaderContent: 'justify-end', ...topAligned },
  }),
]

const banner = computed(() => corredorBanner(status.value))

const bannerText = computed(() => {
  const state = banner.value?.state

  if (state === 'disabled') return t('automation.scripts.list.status.disabled')
  if (state === 'unreachable') return t('automation.scripts.list.status.unreachable')
  if (state !== 'connected') return ''

  const refreshedAt = status.value.refreshedAt

  if (!refreshedAt) return t('automation.scripts.list.status.connectedWithoutRefreshTime')

  return t('automation.scripts.list.status.connected', {
    relativeTime: relativeTime(refreshedAt, locale.value),
    absoluteTime: new Date(refreshedAt).toLocaleString(),
  })
})

const filtered = computed(() => {
  const lcQuery = filter.query.toLocaleLowerCase()
  return items.value
    .filter(
      ({ name, label, description }) =>
        lcQuery.length === 0 ||
        [name, label, description].join(' ').toLocaleLowerCase().includes(lcQuery),
    )
    .filter(({ errors }) => !filter.incScriptsWithErrors || (errors && errors.length > 0))
    .filter(({ triggers }) => !filter.incScriptsWithTriggers || !!triggers)
    .filter(({ iterator }) => !filter.incScriptsWithIterator || !!iterator)
    .filter(({ security }) => !filter.incScriptsWithSecurity || !!security)
    .filter(script =>
      matchesKindFilter(scriptKind(script), {
        server: filter.incServerScripts,
        client: filter.incClientScripts,
      }),
    )
    .filter(script => matchesExtensionFilter(script, filter.extension))
})

const sorted = computed(() => sortScripts(filtered.value, sorting.sortBy, sorting.sortDesc))

const paged = computed(() => {
  const from = (pagination.page - 1) * pagination.limit
  return sorted.value.slice(from, from + pagination.limit)
})

// The shell gates its pager on cursors. A set held in memory pages by index, so
// the cursors here are only the flags that light the two buttons.
const paginationState = computed(() => ({
  page: pagination.page,
  limit: pagination.limit,
  total: sorted.value.length,
  prevPage: pagination.page > 1 ? 'prev' : '',
  nextPage: pagination.page * pagination.limit < sorted.value.length ? 'next' : '',
}))

// Counted over everything fetched, so a count says how much turning the filter
// on would leave, not how much the current filter already left.
const counts = computed(() => ({
  errors: items.value.filter(({ errors }) => errors && errors.length > 0).length,
  triggers: items.value.filter(({ triggers }) => !!triggers).length,
  iterator: items.value.filter(({ iterator }) => !!iterator).length,
  security: items.value.filter(({ security }) => !!security).length,
  server: items.value.filter(script => scriptKind(script) === 'server').length,
  client: items.value.filter(script => scriptKind(script) === 'client').length,
}))

const capabilityOptions = computed(() => [
  {
    key: 'incScriptsWithErrors',
    label: t('automation.scripts.list.filterForm.incScriptsWithErrors', {
      count: counts.value.errors,
    }),
  },
  {
    key: 'incScriptsWithTriggers',
    label: t('automation.scripts.list.filterForm.incScriptsWithTriggers', {
      count: counts.value.triggers,
    }),
  },
  {
    key: 'incScriptsWithIterator',
    label: t('automation.scripts.list.filterForm.incScriptsWithIterator', {
      count: counts.value.iterator,
    }),
  },
  {
    key: 'incScriptsWithSecurity',
    label: t('automation.scripts.list.filterForm.incScriptsWithSecurity', {
      count: counts.value.security,
    }),
  },
])

const kindOptions = computed(() => [
  {
    key: 'incServerScripts',
    label: t('automation.scripts.list.filterForm.incServerScripts', {
      count: counts.value.server,
    }),
  },
  {
    key: 'incClientScripts',
    label: t('automation.scripts.list.filterForm.incClientScripts', {
      count: counts.value.client,
    }),
  },
])

const extensionOptions = computed(() => [
  {
    id: 'any',
    value: ANY_EXTENSION,
    label: t('automation.scripts.list.filterForm.extension.any', { count: items.value.length }),
  },
  ...groupScripts(items.value).map(group => ({
    id: group.extension || 'none',
    value: group.extension,
    label: group.extension
      ? t('automation.scripts.list.filterForm.extension.named', {
          extension: group.extension,
          count: group.items.length,
        })
      : t('automation.scripts.list.filterForm.extension.none', { count: group.items.length }),
  })),
])

function toggleFilterMenu(event) {
  filterMenu.value.toggle(event)
}

function handleSort({ sortField, sortOrder }) {
  sorting.sortBy = sortField
  sorting.sortDesc = sortOrder === -1
}

function handlePageChange({ page, limit }) {
  pagination.page = page || 1
  if (limit) pagination.limit = limit
}

// A narrowed or reordered set makes the page the pager is on meaningless.
watch([filter, sorting], () => (pagination.page = 1), { deep: true })

function kindLabel(script) {
  if (scriptKind(script) === 'server') return t('automation.scripts.list.kind.server')

  const bundle = scriptBundle(script)

  return bundle
    ? t('automation.scripts.list.kind.client', { bundle })
    : t('automation.scripts.list.kind.clientWithoutBundle')
}

function iteratorChip(iterator) {
  const schedule = (iterator?.deferred || []).join(', ')
  const params = {
    eventType: iterator?.eventType || '',
    resourceType: iterator?.resourceType || '',
    action: iterator?.action || '',
    schedule,
  }

  return schedule
    ? t('automation.scripts.list.iterator.chipWithSchedule', params)
    : t('automation.scripts.list.iterator.chip', params)
}

// An iterator's filter names the records it walks, in the same shape a trigger
// constraint does.
function iteratorFilterChips(iterator) {
  const iteratorFilter = iterator?.filter || {}

  return constraintChips(
    Object.keys(iteratorFilter).map(name => ({ name, value: [String(iteratorFilter[name])] })),
  )
}

function securityChips(security) {
  const chips = []

  if (security?.runAs) {
    chips.push(t('automation.scripts.list.security.runAs', { runAs: security.runAs }))
  }
  if (security?.allow?.length) {
    chips.push(t('automation.scripts.list.security.allow', { roles: security.allow.join(', ') }))
  }
  if (security?.deny?.length) {
    chips.push(t('automation.scripts.list.security.deny', { roles: security.deny.join(', ') }))
  }

  return chips
}

async function fetchScripts(indicator) {
  indicator.value = true
  try {
    const result = await $SystemAPI.automationList({})
    const set = Array.isArray(result) ? result : result?.set || []
    items.value = Array.isArray(set) ? set : []
    status.value = Array.isArray(result) ? {} : result || {}
  } catch (e) {
    $toast.toastErrorHandler(t('notification.settings.fetch.error'))(e)
  } finally {
    indicator.value = false
  }
}

const loadScripts = () => fetchScripts(loading)
const refreshScripts = () => fetchScripts(refreshing)

onMounted(() => loadScripts())
</script>
