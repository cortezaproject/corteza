<template>
  <Teleport to="#topbar-title" defer>
    <span>{{ $t('system.apigw.profiler.hit.title') }}</span>
  </Teleport>

  <div class="flex flex-col h-full">
    <CViewContainer scroll>
      <div class="flex items-center gap-2">
        <Button
          icon="pi pi-arrow-left"
          :label="$t('general.label.back')"
          severity="secondary"
          @click="$router.go(-1)"
        />
      </div>

      <div v-if="loading" class="flex items-center justify-center p-8">
        <ProgressSpinner />
      </div>

      <Panel v-else-if="hit" :header="$t('system.apigw.profiler.hit.title')" class="shadow">
        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
          <CFormGroup :label="$t('system.apigw.profiler.hit.columns.hitID')">
            <span class="font-mono text-sm">{{ hit.hitID }}</span>
          </CFormGroup>
          <CFormGroup :label="$t('system.apigw.profiler.hit.columns.status')">
            <Tag
              :value="String(hit.status || '')"
              :severity="hit.status >= 200 && hit.status < 300 ? 'success' : 'warn'"
            />
          </CFormGroup>
          <CFormGroup :label="$t('system.apigw.profiler.hit.columns.time')">
            <span>{{ hit.ts ? new Date(hit.ts).toLocaleString() : '' }}</span>
          </CFormGroup>
          <CFormGroup :label="$t('system.apigw.profiler.hit.request')" class="md:col-span-2">
            <pre class="text-xs bg-emphasis p-3 rounded-lg overflow-auto max-h-64">{{
              JSON.stringify(hit.request, null, 2)
            }}</pre>
          </CFormGroup>
          <CFormGroup :label="$t('system.apigw.profiler.hit.response')" class="md:col-span-2">
            <pre class="text-xs bg-emphasis p-3 rounded-lg overflow-auto max-h-64">{{
              JSON.stringify(hit.response, null, 2)
            }}</pre>
          </CFormGroup>
        </div>
      </Panel>

      <div v-else class="text-center p-8 text-muted-color">
        {{ $t('general.notFound') }}
      </div>
    </CViewContainer>
  </div>
</template>

<script setup>
import { inject, ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { components } from '@planetcrust/human-vue'

const { CViewContainer } = components
const route = useRoute()
const $SystemAPI = inject('$SystemAPI')
const $toast = inject('$toast')
const { t } = useI18n()

const hit = ref(null)
const loading = ref(false)

async function loadHit() {
  loading.value = true
  try {
    hit.value = await $SystemAPI.apigwProfilerHit({ hitID: route.params.hitID })
  } catch (e) {
    $toast.toastErrorHandler(t('notification.gateway.profiler.fetch.error'))(e)
  } finally {
    loading.value = false
  }
}

onMounted(() => loadHit())
</script>
