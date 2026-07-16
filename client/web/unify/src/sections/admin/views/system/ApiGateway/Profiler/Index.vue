<template>
  <Teleport to="#topbar-title" defer>
    <span>{{ $t('system.apigw.profiler.title') }}</span>
  </Teleport>

  <CViewContainer>
    <CResourceList
      primary-key="path"
      :fields="profilerFields"
      :items="items"
      :filter="filter"
      :sorting="sorting"
      :pagination="pagination"
      :loading="loading"
      :action-items="getProfilerActions"
      :translations="{
        resourceSingle: $t('system.apigw.profiler.title'),
        resourcePlural: $t('system.apigw.profiler.title'),
      }"
      hide-search
      hide-pagination
      clickable
      @row-click="handleRowClick"
    >
      <template #header>
        <div class="flex items-center gap-2 flex-wrap w-full">
          <Button
            v-if="items.length"
            :label="$t('system.apigw.profiler.purgeAll')"
            icon="pi pi-trash"
            severity="danger"
            size="small"
            outlined
            :loading="purging"
            @click="purgeAll"
          />
          <Button
            :label="countdown > 0
              ? $t('system.apigw.profiler.refreshingIn', { seconds: countdown })
              : $t('general.label.refresh')"
            icon="pi pi-refresh"
            severity="secondary"
            size="small"
            outlined
            :loading="loading"
            class="ml-auto w-44 justify-center"
            @click="loadData"
          />
        </div>
      </template>

      <template #body-size_min="{ data }">{{ ((data.size_min || 0) / 1000).toFixed(3) }} kB</template>
      <template #body-size_max="{ data }">{{ ((data.size_max || 0) / 1000).toFixed(3) }} kB</template>
      <template #body-size_avg="{ data }">{{ ((data.size_avg || 0) / 1000).toFixed(3) }} kB</template>
      <template #body-time_min="{ data }">{{ (data.time_min || 0).toFixed(2) }} ms</template>
      <template #body-time_max="{ data }">{{ (data.time_max || 0).toFixed(2) }} ms</template>
      <template #body-time_avg="{ data }">{{ (data.time_avg || 0).toFixed(2) }} ms</template>
    </CResourceList>
  </CViewContainer>
</template>

<script setup>
import { inject, ref, onMounted, onUnmounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRouter } from 'vue-router'
import { components } from '@planetcrust/human-vue'

const { CResourceList, CViewContainer } = components
const router = useRouter()
const $SystemAPI = inject('$SystemAPI')
const $toast = inject('$toast')
const { t } = useI18n()

const items = ref([])
const loading = ref(false)
const purging = ref(false)
const countdown = ref(0)
let timer = null

const filter = ref({})
const sorting = ref({ sortBy: 'count', sortDesc: true })
const pagination = ref({ total: 0, limit: 0, page: 1 })

async function loadData() {
  clearTimer()
  loading.value = true
  try {
    const result = await $SystemAPI.apigwProfilerAggregation({})
    items.value = (result?.set || []).map(i => ({ ...i, routeID: btoa(i.path) }))
    pagination.value.total = items.value.length
  } catch (e) {
    $toast.toastErrorHandler(t('notification.gateway.profiler.fetch.error'))(e)
  } finally {
    loading.value = false
    startTimer()
  }
}

async function purgeAll() {
  purging.value = true
  try {
    await $SystemAPI.apigwProfilerPurgeAll()
    $toast.toastSuccess(t('notification.gateway.profiler.purge.success'))
    await loadData()
  } catch (e) {
    $toast.toastErrorHandler(t('notification.gateway.profiler.purge.error'))(e)
  } finally {
    purging.value = false
  }
}

const profilerFields = [
  { key: 'path', header: t('system.apigw.profiler.columns.path'), sortable: true },
  { key: 'count', header: t('system.apigw.profiler.columns.count'), sortable: true, class: 'text-right', pt: { columnHeaderContent: 'justify-end' } },
  { key: 'size_min', header: t('system.apigw.profiler.columns.sizeMin'), sortable: true, class: 'text-right', pt: { columnHeaderContent: 'justify-end' } },
  { key: 'size_max', header: t('system.apigw.profiler.columns.sizeMax'), sortable: true, class: 'text-right', pt: { columnHeaderContent: 'justify-end' } },
  { key: 'size_avg', header: t('system.apigw.profiler.columns.sizeAvg'), sortable: true, class: 'text-right', pt: { columnHeaderContent: 'justify-end' } },
  { key: 'time_min', header: t('system.apigw.profiler.columns.timeMin'), sortable: true, class: 'text-right', pt: { columnHeaderContent: 'justify-end' } },
  { key: 'time_max', header: t('system.apigw.profiler.columns.timeMax'), sortable: true, class: 'text-right', pt: { columnHeaderContent: 'justify-end' } },
  { key: 'time_avg', header: t('system.apigw.profiler.columns.timeAvg'), sortable: true, class: 'text-right', pt: { columnHeaderContent: 'justify-end' } },
]

function handleRowClick({ data }) {
  router.push({ name: 'system.apiGateway.profiler.route', params: { routeID: btoa(data.path) } })
}

function getProfilerActions(data) {
  return [
    {
      label: t('general.label.details'),
      icon: 'pi pi-info-circle',
      command: () => router.push({ name: 'system.apiGateway.profiler.route', params: { routeID: btoa(data.path) } }),
    },
  ]
}

function startTimer() {
  countdown.value = 10
  function tick() {
    countdown.value--
    if (countdown.value <= 0) {
      loadData()
    } else {
      timer = setTimeout(tick, 1000)
    }
  }
  timer = setTimeout(tick, 1000)
}

function clearTimer() {
  if (timer) {
    clearTimeout(timer)
    timer = null
  }
  countdown.value = 0
}

onMounted(() => loadData())
onUnmounted(() => clearTimer())
</script>
