<template>
  <Teleport to="#topbar-title" defer>
    <span>{{ $t('ui.settings.editor.location.title') }}</span>
  </Teleport>

  <div v-if="loading" class="flex items-center justify-center h-full">
    <ProgressSpinner />
  </div>

  <div v-else class="flex flex-col h-full">
    <div class="container mx-auto p-4 flex-1 flex flex-col min-h-0 gap-5 overflow-y-auto">
      <Panel
        :header="$t('ui.settings.editor.location.geosearch.title')"
        toggleable
        :collapsed="false"
        class="shadow"
      >
        <div class="flex flex-col gap-4">
          <div class="flex flex-col gap-1">
            <label class="font-medium text-sm">
              {{ $t('ui.settings.editor.location.geosearch.provider.label') }}
            </label>
            <span class="text-xs text-muted-color">
              {{ $t('ui.settings.editor.location.geosearch.provider.description') }}
            </span>
            <Select
              v-model="location.geoSearchProvider"
              :options="providerOptions"
              option-label="text"
              option-value="value"
              class="w-full md:w-1/2"
              @change="location.geoSearchApiKey = ''"
            />
          </div>

          <div v-if="requiresApiKey" class="flex flex-col gap-1">
            <label class="font-medium text-sm">
              {{ $t('ui.settings.editor.location.geosearch.apiKey.label') }}
            </label>
            <InputText
              v-model="location.geoSearchApiKey"
              :placeholder="$t('ui.settings.editor.location.geosearch.apiKey.placeholder')"
              class="w-full md:w-1/2"
            />
          </div>
        </div>
      </Panel>
    </div>

    <div class="shrink-0 border-t border-surface bg-surface">
      <div class="p-3 flex items-center justify-end">
        <Button
          :label="$t('general.label.save')"
          icon="pi pi-save"
          :loading="saving"
          @click="handleSave"
        />
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed, inject, onMounted, reactive, ref } from 'vue'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()

const $toast = inject('$toast')
const $SystemAPI = inject('$SystemAPI')

const loading = ref(false)
const saving = ref(false)

const location = reactive({
  geoSearchProvider: 'openstreetmap',
  geoSearchApiKey: '',
})

const providerOptions = [
  { value: 'openstreetmap', text: 'OpenStreetMap', requiresApiKey: false },
  { value: 'opencage', text: 'OpenCage', requiresApiKey: true },
  { value: 'esri', text: 'Esri', requiresApiKey: false },
  { value: 'geoapify', text: 'Geoapify', requiresApiKey: true },
  { value: 'geocodeearth', text: 'Geocode Earth', requiresApiKey: true },
  { value: 'google', text: 'Google Maps', requiresApiKey: true },
  { value: 'locationiq', text: 'LocationIQ', requiresApiKey: true },
  { value: 'mapbox', text: 'Mapbox', requiresApiKey: true },
  { value: 'pelias', text: 'Pelias', requiresApiKey: true },
]

const requiresApiKey = computed(() => {
  const p = providerOptions.find(o => o.value === location.geoSearchProvider)
  return p?.requiresApiKey || false
})

async function loadSettings() {
  loading.value = true
  try {
    const result = await $SystemAPI.settingsList({ prefix: 'ui.location' })
    for (const s of result || []) {
      if (s.name === 'ui.location' && s.value) {
        Object.assign(location, s.value)
      }
    }
    if (!location.geoSearchProvider) location.geoSearchProvider = 'openstreetmap'
  } catch (e) {
    $toast.toastErrorHandler(t('notification.settings.location.fetch.error'))(e)
  } finally {
    loading.value = false
  }
}

async function handleSave() {
  saving.value = true
  try {
    await $SystemAPI.settingsUpdate({
      values: [{ name: 'ui.location', value: { ...location } }],
    })
    $toast.toastSuccess(t('notification.settings.location.update.success'))
  } catch (e) {
    $toast.toastErrorHandler(t('notification.settings.location.update.error'))(e)
  } finally {
    saving.value = false
  }
}

onMounted(() => loadSettings())
</script>
