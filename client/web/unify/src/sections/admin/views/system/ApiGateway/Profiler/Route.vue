<template>
  <Teleport to="#topbar-title" defer>
    <span>{{ $t('system.apigw.profiler.title') }} — {{ decodedRoute }}</span>
  </Teleport>

  <CViewContainer>
    <CResourceList
      primary-key="hitID"
      :fields="hitFields"
      :items="items"
      :filter="filter"
      :sorting="sorting"
      :pagination="pagination"
      :loading="loading && items.length === 0"
      :action-items="getHitActions"
      :translations="{
        resourceSingle: $t('system.apigw.profiler.hit.title'),
        resourcePlural: $t('system.apigw.profiler.hit.title'),
      }"
      hide-search
      hide-pagination
      clickable
      @row-click="handleRowClick"
    >
      <template #header>
        <div class="flex items-center gap-2">
          <Button
            icon="pi pi-arrow-left"
            :label="$t('general.label.back')"
            severity="secondary"
            size="small"
            outlined
            @click="$router.push({ name: 'system.apiGateway.profiler' })"
          />
        </div>
      </template>

      <template #body-ts="{ data }">
        {{ data.ts ? new Date(data.ts).toLocaleString() : '' }}
      </template>

      <template v-if="items.length && hasMore" #footer>
        <div class="flex justify-center px-3 py-2">
          <Button
            :label="$t('general.label.loadOlder')"
            :loading="loading"
            severity="secondary"
            size="small"
            @click="loadOlder"
          />
        </div>
      </template>
    </CResourceList>
  </CViewContainer>
</template>

<script setup>
import { inject, ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { components } from '@planetcrust/human-vue'

const { CResourceList, CViewContainer } = components
const route = useRoute()
const router = useRouter()
const $SystemAPI = inject('$SystemAPI')
const $toast = inject('$toast')
const { t } = useI18n()

const PAGE_LIMIT = 50

const items = ref([])
const loading = ref(false)
const hasMore = ref(true)

const filter = ref({})
const sorting = ref({ sortBy: 'ts', sortDesc: true })
const pagination = ref({ total: 0, limit: 0, page: 1 })

const decodedRoute = computed(() => {
  try {
    return atob(route.params.routeID)
  } catch {
    return route.params.routeID
  }
})

async function load(reset = false) {
  if (loading.value) return
  loading.value = true
  try {
    const before =
      reset || items.value.length === 0 ? undefined : items.value[items.value.length - 1].hitID
    const result = await $SystemAPI.apigwProfilerRoute({
      routeID: route.params.routeID,
      before,
      limit: PAGE_LIMIT,
    })
    const set = result?.set || []
    items.value = reset ? set : [...items.value, ...set]
    hasMore.value = set.length >= PAGE_LIMIT
    pagination.value.total = items.value.length
  } catch (e) {
    $toast.toastErrorHandler(t('notification.gateway.profiler.fetch.error'))(e)
  } finally {
    loading.value = false
  }
}

function loadOlder() {
  load(false)
}

const hitFields = [
  { key: 'hitID', header: t('system.apigw.profiler.hit.columns.hitID') },
  { key: 'status', header: t('system.apigw.profiler.hit.columns.status'), sortable: true },
  { key: 'ts', header: t('system.apigw.profiler.hit.columns.time'), sortable: true },
]

function handleRowClick({ data }) {
  router.push({
    name: 'system.apiGateway.profiler.hit',
    params: { routeID: route.params.routeID, hitID: data.hitID },
  })
}

function getHitActions(data) {
  return [
    {
      label: t('general.label.details'),
      icon: 'pi pi-info-circle',
      command: () =>
        router.push({
          name: 'system.apiGateway.profiler.hit',
          params: { routeID: route.params.routeID, hitID: data.hitID },
        }),
    },
  ]
}

onMounted(() => load(true))
</script>
