<template>
  <Teleport to="#topbar-title" defer>
    <span>{{ $t('system.apigw.profiler.title') }} — {{ decodedRoute }}</span>
  </Teleport>

  <div class="flex flex-col h-full">
    <div class="container mx-auto p-4 flex-1 flex flex-col min-h-0 gap-4 overflow-y-auto">
      <div class="flex items-center gap-2">
        <Button
          icon="pi pi-arrow-left"
          :label="$t('general.label.back')"
          severity="secondary"
          @click="$router.push({ name: 'system.apiGateway.profiler' })"
        />
      </div>

      <CResourceTable
        :items="items"
        :fields="hitFields"
        :loading="loading"
        :action-items="getHitActions"
        primary-key="hitID"
        :empty-message="$t('general.notFound')"
      >
        <template #body-ts="{ data }">
          {{ data.ts ? new Date(data.ts).toLocaleString() : '' }}
        </template>
      </CResourceTable>
    </div>
  </div>
</template>

<script setup>
import { inject, ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { components } from '@cortezaproject/corteza-vue-next'

const { CResourceTable } = components
const route = useRoute()
const router = useRouter()
const $SystemAPI = inject('$SystemAPI')
const $toast = inject('$toast')
const { t } = useI18n()

const items = ref([])
const loading = ref(false)

const decodedRoute = computed(() => {
  try {
    return atob(route.params.routeID)
  } catch {
    return route.params.routeID
  }
})

async function loadData() {
  loading.value = true
  try {
    const result = await $SystemAPI.apigwProfilerHitList({ routeID: route.params.routeID })
    items.value = result?.set || []
  } catch (e) {
    $toast.toastErrorHandler(t('notification.gateway.profiler.fetch.error'))(e)
  } finally {
    loading.value = false
  }
}
const hitFields = [
  { key: 'hitID', header: t('system.apigw.profiler.hit.columns.hitID') },
  { key: 'status', header: t('system.apigw.profiler.hit.columns.status'), sortable: true },
  { key: 'ts', header: t('system.apigw.profiler.hit.columns.time'), sortable: true },
]

function getHitActions(data) {
  return [
    {
      label: t('general.label.details'),
      icon: 'pi pi-info-circle',
      command: () => router.push({ name: 'system.apiGateway.profiler.hit', params: { routeID: route.params.routeID, hitID: data.hitID } }),
    },
  ]
}

onMounted(() => loadData())
</script>
