<template>
  <Teleport to="#topbar-title" defer>
    <span>{{ $t('ui.settings.editor.title') }}</span>
  </Teleport>

  <div v-if="loading" class="flex items-center justify-center h-full">
    <ProgressSpinner />
  </div>

  <div v-else class="flex flex-col h-full">
    <div class="container mx-auto p-4 flex-1 flex flex-col min-h-0 gap-4 overflow-y-auto">
      <!-- Navigation (Topbar) -->
      <Panel
        :header="$t('ui.settings.editor.topbar.title')"
        toggleable
        :collapsed="false"
        class="shadow"
      >
        <h5 class="text-sm font-semibold mb-2">
          {{ $t('ui.settings.editor.topbar.general') }}
        </h5>
        <div class="grid grid-cols-1 md:grid-cols-2 gap-4 mb-4">
          <CInputSwitch
            v-model="topbar.hideAppSelector"
            :label="$t('ui.settings.editor.topbar.app-selector.hide')"
          />
          <CInputSwitch
            v-model="topbar.hideNotifications"
            :label="$t('ui.settings.editor.topbar.notifications.hide')"
          />
          <CInputSwitch
            v-model="topbar.hideHelp"
            :label="$t('ui.settings.editor.topbar.help.hide')"
          />
          <CInputSwitch
            v-model="topbar.hideProfile"
            :label="$t('ui.settings.editor.topbar.profile.hide')"
          />
          <CInputSwitch v-model="hideDrafts" :label="$t('ui.settings.editor.topbar.drafts.hide')" />
        </div>

        <Divider />

        <!-- Help sub-section -->
        <h5 class="text-sm font-semibold mb-2">
          {{ $t('ui.settings.editor.topbar.help.title') }}
        </h5>
        <div class="grid grid-cols-1 md:grid-cols-2 gap-4 mb-4">
          <CInputSwitch
            v-model="topbar.hideForumLink"
            :label="$t('ui.settings.editor.topbar.help.hide-forum-link')"
          />
          <CInputSwitch
            v-model="topbar.hideDocumentationLink"
            :label="$t('ui.settings.editor.topbar.help.hide-documentation-link')"
          />
          <CInputSwitch
            v-model="topbar.hideFeedbackLink"
            :label="$t('ui.settings.editor.topbar.help.hide-feedback-link')"
          />
        </div>

        <!-- Help custom links -->
        <div class="flex flex-col gap-2 mb-4">
          <label class="font-medium text-sm">
            {{ $t('ui.settings.editor.topbar.links.title') }}
          </label>
          <DataTable :value="topbar.helpLinks" class="border rounded" size="small">
            <Column :header="$t('ui.settings.editor.topbar.links.handle')" class="w-1/3">
              <template #body="{ data }">
                <InputText v-model="data.handle" size="small" class="w-full" />
              </template>
            </Column>
            <Column :header="$t('ui.settings.editor.topbar.links.url')" class="w-1/2">
              <template #body="{ data }">
                <InputText v-model="data.url" size="small" class="w-full" />
              </template>
            </Column>
            <Column
              :header="$t('ui.settings.editor.topbar.links.new-tab')"
              class="w-20 text-center"
            >
              <template #body="{ data }">
                <Checkbox v-model="data.newTab" :binary="true" />
              </template>
            </Column>
            <Column class="w-16 text-right">
              <template #body="{ index }">
                <Button
                  icon="pi pi-trash"
                  severity="danger"
                  text
                  rounded
                  size="small"
                  @click="topbar.helpLinks.splice(index, 1)"
                />
              </template>
            </Column>
          </DataTable>
          <div>
            <Button
              :label="$t('general.label.add')"
              icon="pi pi-plus"
              size="small"
              severity="secondary"
              @click="topbar.helpLinks.push({ handle: '', url: '', newTab: true })"
            />
          </div>
        </div>

        <Divider />

        <!-- Profile sub-section -->
        <h5 class="text-sm font-semibold mb-2">
          {{ $t('ui.settings.editor.topbar.profile.title') }}
        </h5>
        <div class="grid grid-cols-1 md:grid-cols-2 gap-4 mb-4">
          <CInputSwitch
            v-model="topbar.hideProfileLink"
            :label="$t('ui.settings.editor.topbar.profile.hide-profile-link')"
          />
          <CInputSwitch
            v-model="topbar.hideChangePasswordLink"
            :label="$t('ui.settings.editor.topbar.profile.hide-change-password-link')"
          />
          <CInputSwitch
            v-model="topbar.hideThemeSelector"
            :label="$t('ui.settings.editor.topbar.profile.hide-theme-selector')"
          />
        </div>

        <!-- Profile custom links -->
        <div class="flex flex-col gap-2">
          <label class="font-medium text-sm">
            {{ $t('ui.settings.editor.topbar.links.title') }}
          </label>
          <DataTable :value="topbar.profileLinks" class="border rounded" size="small">
            <Column :header="$t('ui.settings.editor.topbar.links.handle')" class="w-1/3">
              <template #body="{ data }">
                <InputText v-model="data.handle" size="small" class="w-full" />
              </template>
            </Column>
            <Column :header="$t('ui.settings.editor.topbar.links.url')" class="w-1/2">
              <template #body="{ data }">
                <InputText v-model="data.url" size="small" class="w-full" />
              </template>
            </Column>
            <Column
              :header="$t('ui.settings.editor.topbar.links.new-tab')"
              class="w-20 text-center"
            >
              <template #body="{ data }">
                <Checkbox v-model="data.newTab" :binary="true" />
              </template>
            </Column>
            <Column class="w-16 text-right">
              <template #body="{ index }">
                <Button
                  icon="pi pi-trash"
                  severity="danger"
                  text
                  rounded
                  size="small"
                  @click="topbar.profileLinks.splice(index, 1)"
                />
              </template>
            </Column>
          </DataTable>
          <div>
            <Button
              :label="$t('general.label.add')"
              icon="pi pi-plus"
              size="small"
              severity="secondary"
              @click="topbar.profileLinks.push({ handle: '', url: '', newTab: true })"
            />
          </div>
        </div>
      </Panel>

      <!-- Location (Geosearch) -->
      <Panel
        :header="$t('ui.settings.editor.location.title')"
        toggleable
        :collapsed="false"
        class="shadow"
      >
        <h5 class="text-sm font-semibold mb-2">
          {{ $t('ui.settings.editor.location.geosearch.title') }}
        </h5>
        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div class="flex flex-col gap-1">
            <label class="font-medium text-sm">
              {{ $t('ui.settings.editor.location.geosearch.provider.label') }}
            </label>
            <span class="text-xs text-surface-500">
              {{ $t('ui.settings.editor.location.geosearch.provider.description') }}
            </span>
            <Select
              v-model="location.geoSearchProvider"
              :options="providerOptions"
              option-label="text"
              option-value="value"
              class="w-full"
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
              class="w-full"
            />
          </div>
        </div>
      </Panel>

      <!-- Theming (simplified) -->
      <Panel
        :header="$t('ui.settings.editor.theming.title')"
        toggleable
        :collapsed="false"
        class="shadow"
      >
        <div class="flex flex-col gap-4">
          <div v-for="theme in themes" :key="theme.id" class="flex flex-col gap-1">
            <label class="font-medium text-sm">
              {{ $t('ui.settings.editor.corteza-studio.custom-css') }} — {{ theme.title }}
            </label>
            <Textarea v-model="theme.customCSS" rows="8" class="w-full font-mono text-sm" />
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
const rawSettings = reactive({})

// Topbar settings
const topbar = reactive({
  hideAppSelector: false,
  hideNotifications: false,
  hideHelp: false,
  hideProfile: false,
  showDrafts: true,
  showSearch: true,
  hideForumLink: false,
  hideDocumentationLink: false,
  hideFeedbackLink: false,
  hideProfileLink: false,
  hideChangePasswordLink: false,
  hideThemeSelector: false,
  helpLinks: [],
  profileLinks: [],
})

const hideDrafts = computed({
  get: () => topbar.showDrafts !== true,
  set: val => {
    topbar.showDrafts = !val
  },
})

// Location settings
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

// Theming
const themeTabs = ['general', 'light', 'dark']
const themes = reactive(
  themeTabs.map(id => ({
    id,
    title: t(`ui.settings.editor.corteza-studio.tabs.${id}`),
    customCSS: '',
  })),
)

async function loadSettings() {
  loading.value = true
  try {
    const result = await $SystemAPI.settingsList({ prefix: 'ui' })
    for (const s of result || []) {
      rawSettings[s.name] = s.value
    }

    // Hydrate topbar
    const topbarData = rawSettings['ui.topbar'] || {}
    Object.assign(topbar, {
      ...topbarData,
      helpLinks: topbarData.helpLinks || [],
      profileLinks: topbarData.profileLinks || [],
    })

    // Hydrate location
    const locationData = rawSettings['ui.location'] || {}
    Object.assign(location, locationData)
    if (!location.geoSearchProvider) location.geoSearchProvider = 'openstreetmap'

    // Hydrate theming
    const customCSS = rawSettings['ui.studio.custom-css'] || []
    for (const theme of themes) {
      const existing = customCSS.find(c => c.id === theme.id)
      if (existing) theme.customCSS = existing.values || ''
    }
  } catch (e) {
    $toast.toastErrorHandler(t('notification.settings.fetch.error'))(e)
  } finally {
    loading.value = false
  }
}

async function handleSave() {
  saving.value = true
  try {
    const values = [
      { name: 'ui.topbar', value: { ...topbar } },
      { name: 'ui.location', value: { ...location } },
      {
        name: 'ui.studio.custom-css',
        value: themes.map(theme => ({
          id: theme.id,
          title: theme.title,
          values: theme.customCSS,
        })),
      },
    ]

    await $SystemAPI.settingsUpdate({ values })
    $toast.toastSuccess(t('notification.settings.update.success'))
  } catch (e) {
    $toast.toastErrorHandler(t('notification.settings.update.error'))(e)
  } finally {
    saving.value = false
  }
}

onMounted(() => loadSettings())
</script>
