<template>
  <Teleport to="#topbar-title" defer>
    <span>{{ $t('automation.scripts.list.title') }}</span>
  </Teleport>

  <div v-if="loading" class="flex items-center justify-center h-full">
    <ProgressSpinner />
  </div>

  <div v-else class="flex flex-col h-full">
    <CViewContainer scroll>
      <Panel
        :header="$t('automation.scripts.list.title')"
        toggleable
        :collapsed="false"
        class="shadow"
      >
        <!-- Corredor connection, as the script list reported it -->
        <Message v-if="banner" :severity="banner.severity" :closable="false" class="mb-4">
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

        <!-- Search + Filters -->
        <div class="flex flex-col gap-4 mb-4">
          <CFormGroup :label="$t('automation.scripts.list.filter.searchQuery')">
            <InputText v-model="filter.query" class="w-full md:w-1/2" />
          </CFormGroup>

          <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
            <CInputSwitch
              v-model="filter.incScriptsWithErrors"
              :label="
                $t('automation.scripts.list.filter.incScriptsWithErrors', {
                  count: totalScriptsWithErrors,
                })
              "
            />
            <CInputSwitch
              v-model="filter.incScriptsWithTriggers"
              :label="
                $t('automation.scripts.list.filter.incScriptsWithTriggers', {
                  count: totalScriptsWithTriggers,
                })
              "
            />
            <CInputSwitch
              v-model="filter.incScriptsWithIterator"
              :label="
                $t('automation.scripts.list.filter.incScriptsWithIterator', {
                  count: totalScriptsWithIterator,
                })
              "
            />
            <CInputSwitch
              v-model="filter.incScriptsWithSecurity"
              :label="
                $t('automation.scripts.list.filter.incScriptsWithSecurity', {
                  count: totalScriptsWithSecurity,
                })
              "
            />
            <CInputSwitch
              v-model="filter.incServerScripts"
              :label="
                $t('automation.scripts.list.filter.incServerScripts', {
                  count: totalServerScripts,
                })
              "
            />
            <CInputSwitch
              v-model="filter.incClientScripts"
              :label="
                $t('automation.scripts.list.filter.incClientScripts', {
                  count: totalClientScripts,
                })
              "
            />
          </div>
        </div>

        <Divider />

        <div v-if="grouped.length" class="flex flex-col gap-6">
          <CResourceTable
            v-for="group in grouped"
            :key="group.extension || '-'"
            :items="group.items"
            :fields="scriptFields"
            primary-key="name"
            :empty-message="$t('general.resourceList.noItems')"
          >
            <template #header>
              <div class="flex items-center gap-2">
                <i class="pi pi-folder text-muted-color" />
                <span class="font-medium">
                  {{ group.extension || $t('automation.scripts.list.groups.root') }}
                </span>
                <Tag :value="`${group.items.length}`" severity="secondary" class="text-xs" />
              </div>
            </template>

            <template #body-name="{ data }">
              <div class="flex flex-col gap-1">
                <div class="flex items-center gap-2 flex-wrap">
                  <span v-if="data.label" class="font-medium">{{ data.label }}</span>
                  <span v-else class="text-muted-color italic">
                    {{ $t('automation.scripts.list.labelMissing') }}
                  </span>

                  <Tag :value="kindLabel(data)" severity="contrast" class="text-xs" />
                </div>

                <span v-if="data.description" class="text-xs text-muted-color max-w-md">
                  {{ data.description }}
                </span>

                <code class="text-xs text-muted-color truncate max-w-md">{{ data.name }}</code>

                <div
                  v-for="(trigger, i) in triggerRows(data)"
                  :key="`trigger-${i}`"
                  class="flex items-center gap-1 flex-wrap mt-1"
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

                <div v-if="data.iterator" class="flex items-center gap-1 flex-wrap mt-1">
                  <Tag :value="iteratorChip(data.iterator)" severity="warn" class="text-xs" />
                  <Tag
                    v-for="(constraint, j) in iteratorFilterChips(data.iterator)"
                    :key="`iterator-filter-${j}`"
                    :value="constraint"
                    severity="secondary"
                    class="text-xs"
                  />
                </div>

                <div v-if="data.security" class="flex items-center gap-1 flex-wrap mt-1">
                  <Tag
                    v-for="(chip, i) in securityChips(data.security)"
                    :key="`security-${i}`"
                    :value="chip"
                    severity="contrast"
                    class="text-xs"
                  />
                </div>

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

            <template #body-changedAt="{ data }">
              <span v-if="changedAt(data)" class="text-sm text-muted-color">
                {{ formatDate(changedAt(data)) }}
              </span>
            </template>
          </CResourceTable>
        </div>

        <div v-else class="flex items-center justify-center p-4 text-muted-color">
          {{ $t('general.resourceList.noItems') }}
        </div>
      </Panel>
    </CViewContainer>
  </div>
</template>

<script setup>
import { computed, inject, onMounted, reactive, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { changedAt, changedAtField, components, constraintChips } from '@planetcrust/human-vue'
import {
  corredorBanner,
  groupScripts,
  matchesKindFilter,
  relativeTime,
  scriptBundle,
  scriptKind,
  triggerRows,
} from './script-inventory'

const { CResourceTable, CViewContainer } = components
const { t, locale } = useI18n()

const $toast = inject('$toast')
const $SystemAPI = inject('$SystemAPI')

const loading = ref(false)
const refreshing = ref(false)
const items = ref([])
const status = ref({})

const filter = reactive({
  query: '',
  incScriptsWithErrors: false,
  incScriptsWithTriggers: false,
  incScriptsWithIterator: false,
  incScriptsWithSecurity: false,
  incServerScripts: false,
  incClientScripts: false,
})

// Corredor scripts are listed in memory from the server's script bundle, so the
// column cannot be sorted the way a stored resource's can. CResourceTable takes
// header/body classes rather than the shared `class`.
const scriptFields = [
  { key: 'name', header: t('automation.scripts.list.columns.name') },
  changedAtField(t('general.columns.changedAt'), {
    sortable: false,
    headerStyle: 'width: 12rem',
    headerClass: 'text-right',
    bodyClass: 'text-right',
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
      ({ name, label }) =>
        lcQuery.length === 0 || (name + ' ' + (label || '')).toLocaleLowerCase().includes(lcQuery),
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
})

const grouped = computed(() => groupScripts(filtered.value))

const totalScriptsWithErrors = computed(
  () => items.value.filter(({ errors }) => errors && errors.length > 0).length,
)
const totalScriptsWithTriggers = computed(
  () => items.value.filter(({ triggers }) => !!triggers).length,
)
const totalScriptsWithIterator = computed(
  () => items.value.filter(({ iterator }) => !!iterator).length,
)
const totalScriptsWithSecurity = computed(
  () => items.value.filter(({ security }) => !!security).length,
)
const totalServerScripts = computed(
  () => items.value.filter(script => scriptKind(script) === 'server').length,
)
const totalClientScripts = computed(
  () => items.value.filter(script => scriptKind(script) === 'client').length,
)

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

function formatDate(dateStr) {
  if (!dateStr) return ''
  const d = new Date(dateStr)
  return d.toLocaleDateString(undefined, { year: 'numeric', month: 'short', day: 'numeric' })
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
