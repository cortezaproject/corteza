<template>
  <Teleport to="#topbar-title" defer>
    <span>{{ $t('system.apigw.profiler.title') }}</span>
  </Teleport>

  <div class="flex flex-col h-full">
    <div class="container mx-auto p-4 flex-1 flex flex-col min-h-0 gap-4 overflow-y-auto">
      <div class="flex items-center gap-2">
        <Button
          :label="$t('general.label.refresh')"
          icon="pi pi-refresh"
          severity="secondary"
          outlined
          @click="loadData"
          :loading="loading"
        />
        <span v-if="countdown > 0" class="text-sm text-muted-color">
          {{ $t('system.apigw.profiler.refreshingIn', { seconds: countdown }) }}
        </span>
        <Button
          v-if="items.length"
          :label="$t('system.apigw.profiler.purgeAll')"
          icon="pi pi-trash"
          severity="danger"
          outlined
          size="small"
          @click="purgeAll"
          :loading="purging"
          class="ml-auto"
        />
      </div>

      <CResourceTable
        :items="items"
        :fields="profilerFields"
        :loading="loading"
        :action-items="getProfilerActions"
        primary-key="path"
        :empty-message="$t('general.notFound')"
      >
        <template #body-size_min="{ data }">{{ ((data.size_min || 0) / 1000).toFixed(3) }} kB</template>
        <template #body-size_max="{ data }">{{ ((data.size_max || 0) / 1000).toFixed(3) }} kB</template>
        <template #body-size_avg="{ data }">{{ ((data.size_avg || 0) / 1000).toFixed(3) }} kB</template>
        <template #body-time_min="{ data }">{{ (data.time_min || 0).toFixed(2) }} ms</template>
        <template #body-time_max="{ data }">{{ (data.time_max || 0).toFixed(2) }} ms</template>
        <template #body-time_avg="{ data }">{{ (data.time_avg || 0).toFixed(2) }} ms</template>
      </CResourceTable>
    </div>
  </div>
</template>

<script setup>
import { inject, ref, onMounted, onUnmounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRouter } from 'vue-router'
import { components } from '@cortezaproject/corteza-vue-next'

const { CResourceTable } = components
const router = useRouter()
const $SystemAPI = inject('$SystemAPI')
const $toast = inject('$toast')
const { t } = useI18n()

const items = ref([])
const loading = ref(false)
const purging = ref(false)
const countdown = ref(0)
let timer = null

async function loadData() {
  clearTimer()
  loading.value = true
  try {
    const result = await $SystemAPI.apigwProfilerAggregation({})
    items.value = (result?.set || []).map(i => ({ ...i, routeID: btoa(i.path) }))
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
  { key: 'count', header: t('system.apigw.profiler.columns.count'), sortable: true, headerClass: 'text-right', bodyClass: 'text-right' },
  { key: 'size_min', header: t('system.apigw.profiler.columns.sizeMin'), sortable: true, headerClass: 'text-right', bodyClass: 'text-right' },
  { key: 'size_max', header: t('system.apigw.profiler.columns.sizeMax'), sortable: true, headerClass: 'text-right', bodyClass: 'text-right' },
  { key: 'size_avg', header: t('system.apigw.profiler.columns.sizeAvg'), sortable: true, headerClass: 'text-right', bodyClass: 'text-right' },
  { key: 'time_min', header: t('system.apigw.profiler.columns.timeMin'), sortable: true, headerClass: 'text-right', bodyClass: 'text-right' },
  { key: 'time_max', header: t('system.apigw.profiler.columns.timeMax'), sortable: true, headerClass: 'text-right', bodyClass: 'text-right' },
  { key: 'time_avg', header: t('system.apigw.profiler.columns.timeAvg'), sortable: true, headerClass: 'text-right', bodyClass: 'text-right' },
]

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
