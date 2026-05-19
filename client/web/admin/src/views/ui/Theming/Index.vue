<template>
  <Teleport to="#topbar-title" defer>
    <span>{{ $t('ui.settings.editor.human-studio.title') }}</span>
  </Teleport>

  <div v-if="loading" class="flex items-center justify-center h-full">
    <ProgressSpinner />
  </div>

  <div v-else class="flex flex-col h-full">
    <div class="container mx-auto p-4 flex-1 flex flex-col min-h-0 gap-5 overflow-y-auto">
      <!-- Branding section -->
      <Panel
        :header="$t('ui.settings.editor.human-studio.branding.title')"
        toggleable
        class="shadow mb-5"
      >
        <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
          <!-- Main Logo -->
          <div class="flex flex-col gap-2">
            <label class="font-medium text-sm text-primary">
              {{ $t('ui.settings.editor.human-studio.mainLogo.title') }}
            </label>

            <CFileDropZone
              accept="image/*"
              :uploading="mainLogoUploading"
              :error="mainLogoError"
              :preview-url="mainLogoUrl"
              :clearable="mainLogoIsCustom"
              :drop-label="$t('ui.settings.editor.human-studio.mainLogo.uploader.instructions')"
              :uploading-label="$t('ui.settings.editor.human-studio.mainLogo.uploader.uploading')"
              :label="$t('ui.settings.editor.human-studio.mainLogo.title')"
              compact
              preview-max-width="100%"
              preview-max-height="200px"
              @select="onMainLogoSelect"
              @clear="onMainLogoClear"
            />
          </div>

          <!-- Icon Logo -->
          <div class="flex flex-col gap-2">
            <label class="font-medium text-sm text-primary">
              {{ $t('ui.settings.editor.human-studio.iconLogo.title') }}
            </label>

            <CFileDropZone
              accept="image/*"
              :uploading="iconLogoUploading"
              :error="iconLogoError"
              :preview-url="iconLogoUrl"
              :clearable="iconLogoIsCustom"
              :drop-label="$t('ui.settings.editor.human-studio.iconLogo.uploader.instructions')"
              :uploading-label="$t('ui.settings.editor.human-studio.iconLogo.uploader.uploading')"
              :label="$t('ui.settings.editor.human-studio.iconLogo.title')"
              compact
              preview-max-width="100%"
              preview-max-height="200px"
              @select="onIconLogoSelect"
              @clear="onIconLogoClear"
            />
          </div>
        </div>
      </Panel>

      <Panel
        :header="$t('ui.settings.editor.human-studio.title')"
        toggleable
        :collapsed="false"
        class="shadow"
        :pt="{ content: { style: 'padding: 0 !important' } }"
      >
        <Tabs value="general">
          <TabList>
            <Tab v-for="theme in themes" :key="theme.id" :value="theme.id">
              {{ theme.title }}
            </Tab>
          </TabList>

          <TabPanels>
            <TabPanel v-for="theme in themes" :key="theme.id" :value="theme.id">
              <!-- Color variables for light/dark tabs only -->
              <div v-if="theme.id !== 'general'" class="grid grid-cols-1 lg:grid-cols-2 gap-4 mb-6">
                <div v-for="key in themeVariableKeys" :key="key" class="flex flex-col gap-1">
                  <label class="font-medium text-sm text-primary">
                    {{ $t(`ui.settings.editor.human-studio.theme.variables.${key}.label`) }}
                  </label>
                  <span class="text-xs text-muted-color">
                    {{ $t(`ui.settings.editor.human-studio.theme.variables.${key}.description`) }}
                  </span>
                  <div class="flex items-center gap-2">
                    <CInputColorPicker
                      :model-value="'#' + (theme.variables[key] || '')"
                      :default-value="'#' + (theme.defaultVariables[key] || '')"
                      show-text
                      @update:model-value="
                        theme.variables[key] = $event.replace(/^#/, '').substring(0, 6)
                      "
                    />
                    <Button
                      icon="pi pi-undo"
                      severity="secondary"
                      text
                      rounded
                      size="small"
                      :title="$t('ui.settings.editor.human-studio.label.default')"
                      @click="theme.variables[key] = theme.defaultVariables[key]"
                    />
                  </div>
                </div>
              </div>

              <!-- Custom CSS for all tabs -->
              <div class="flex flex-col gap-1">
                <label class="font-medium text-sm text-primary">
                  {{ $t('ui.settings.editor.human-studio.custom-css') }}
                </label>
                <Textarea v-model="theme.customCSS" rows="16" class="w-full font-mono text-sm" />
              </div>
            </TabPanel>
          </TabPanels>
        </Tabs>
      </Panel>
    </div>

    <CEditorActions>
      <Button
        :label="$t('general.label.save')"
        icon="pi pi-save"
        :loading="saving"
        @click="handleSave"
      />
    </CEditorActions>
  </div>
</template>

<script setup>
import { inject, onMounted, reactive, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { setThemes, useTheme, components, useFileUpload } from '@planetcrust/human-vue'

const { CInputColorPicker, CFileDropZone } = components

const { t } = useI18n()

const $toast = inject('$toast')
const $SystemAPI = inject('$SystemAPI')
const $Settings = inject('$Settings')

const loading = ref(false)
const saving = ref(false)

// Logo upload state
const {
  uploading: mainLogoUploading,
  uploadError: mainLogoError,
  uploadFileRaw: uploadMainLogoRaw,
  reset: resetMainLogoUpload,
} = useFileUpload()

const {
  uploading: iconLogoUploading,
  uploadError: iconLogoError,
  uploadFileRaw: uploadIconLogoRaw,
  reset: resetIconLogoUpload,
} = useFileUpload()

const mainLogoUrl = ref('')
const iconLogoUrl = ref('')
const mainLogoIsCustom = ref(false)
const iconLogoIsCustom = ref(false)

const themeVariableKeys = [
  'primary',
  'success',
  'warning',
  'danger',
  'body-bg',
  'sidebar-bg',
  'topbar-bg',
]

const lightModeDefaults = {
  primary: '09344E',
  success: '43AA8B',
  warning: 'E27646',
  danger: 'E54122',
  'body-bg': 'f4f4f5',
  'sidebar-bg': 'ffffff',
  'topbar-bg': 'f4f4f5',
}

const darkModeDefaults = {
  primary: 'E56B5B',
  success: '43AA8B',
  warning: 'E27646',
  danger: 'E54122',
  'body-bg': '27272a',
  'sidebar-bg': '18181b',
  'topbar-bg': '27272a',
}

const themes = reactive([
  {
    id: 'general',
    title: t('ui.settings.editor.human-studio.tabs.general'),
    variables: {},
    defaultVariables: {},
    customCSS: '',
  },
  {
    id: 'light',
    title: t('ui.settings.editor.human-studio.tabs.light'),
    variables: { ...lightModeDefaults },
    defaultVariables: { ...lightModeDefaults },
    customCSS: '',
  },
  {
    id: 'dark',
    title: t('ui.settings.editor.human-studio.tabs.dark'),
    variables: { ...darkModeDefaults },
    defaultVariables: { ...darkModeDefaults },
    customCSS: '',
  },
])

function stripHash(color) {
  if (!color) return ''
  return color.replace(/^#/, '')
}

function refreshLogoUrls() {
  mainLogoUrl.value = $Settings.attachment('ui.mainLogo') || ''
  iconLogoUrl.value = $Settings.attachment('ui.iconLogo') || ''

  // Only show delete button when a custom logo was uploaded (setting starts with 'attachment:')
  const mainLogoRaw = $Settings.get('ui.mainLogo', '')
  const iconLogoRaw = $Settings.get('ui.iconLogo', '')
  mainLogoIsCustom.value = typeof mainLogoRaw === 'string' && mainLogoRaw.startsWith('attachment:')
  iconLogoIsCustom.value = typeof iconLogoRaw === 'string' && iconLogoRaw.startsWith('attachment:')
}

async function onMainLogoSelect(files) {
  const file = files[0]
  if (!file) return

  try {
    const endpoint = $SystemAPI.baseURL + $SystemAPI.settingsSetEndpoint({ key: 'ui.mainLogo' })
    const token = $SystemAPI.accessTokenFn ? $SystemAPI.accessTokenFn() : ''
    await uploadMainLogoRaw(file, { url: endpoint, token })
    await $Settings.fetch()
    refreshLogoUrls()
    $toast.toastSuccess(t('notification.settings.theming.update.success'))
  } catch {
    // uploadError is set by composable
  }
}

async function onIconLogoSelect(files) {
  const file = files[0]
  if (!file) return

  try {
    const endpoint = $SystemAPI.baseURL + $SystemAPI.settingsSetEndpoint({ key: 'ui.iconLogo' })
    const token = $SystemAPI.accessTokenFn ? $SystemAPI.accessTokenFn() : ''
    await uploadIconLogoRaw(file, { url: endpoint, token })
    await $Settings.fetch()
    refreshLogoUrls()
    $toast.toastSuccess(t('notification.settings.theming.update.success'))
  } catch {
    // uploadError is set by composable
  }
}

async function onMainLogoClear() {
  try {
    await $SystemAPI.settingsUpdate({
      values: [{ name: 'ui.mainLogo', value: null }],
    })
    await $Settings.fetch()
    refreshLogoUrls()
    resetMainLogoUpload()
    $toast.toastSuccess(t('notification.settings.theming.update.success'))
  } catch (e) {
    $toast.toastErrorHandler(t('notification.settings.theming.update.error'))(e)
  }
}

async function onIconLogoClear() {
  try {
    await $SystemAPI.settingsUpdate({
      values: [{ name: 'ui.iconLogo', value: null }],
    })
    await $Settings.fetch()
    refreshLogoUrls()
    resetIconLogoUpload()
    $toast.toastSuccess(t('notification.settings.theming.update.success'))
  } catch (e) {
    $toast.toastErrorHandler(t('notification.settings.theming.update.error'))(e)
  }
}

async function loadSettings() {
  loading.value = true
  try {
    refreshLogoUrls()

    const result = await $SystemAPI.settingsList({ prefix: 'ui.studio' })
    for (const s of result || []) {
      if (s.name === 'ui.studio.themes' && Array.isArray(s.value)) {
        for (const theme of themes) {
          const existing = s.value.find(c => c.id === theme.id)
          if (existing) {
            const parsed = existing.values ? JSON.parse(existing.values) : {}
            // Strip # prefix from stored values for ColorPicker (hex format without #)
            const stripped = {}
            for (const [k, v] of Object.entries(parsed)) {
              stripped[k] = stripHash(v)
            }
            if (theme.id !== 'general' && Object.keys(stripped).length > 0) {
              theme.variables = stripped
            }
          }
        }
      }
      if (s.name === 'ui.studio.custom-css' && Array.isArray(s.value)) {
        for (const theme of themes) {
          const existing = s.value.find(c => c.id === theme.id)
          if (existing) theme.customCSS = existing.values || ''
        }
      }
    }
  } catch (e) {
    $toast.toastErrorHandler(t('notification.settings.theming.fetch.error'))(e)
  } finally {
    loading.value = false
  }
}

async function handleSave() {
  saving.value = true
  try {
    // Add # prefix for stored values
    const themeValues = themes.map(theme => {
      const variables = {}
      for (const [k, v] of Object.entries(theme.variables)) {
        variables[k] = v && !v.startsWith('#') ? `#${v}` : v
      }
      return {
        id: theme.id,
        title: theme.title,
        values: JSON.stringify(variables),
      }
    })

    const customCssValues = themes.map(theme => ({
      id: theme.id,
      title: theme.title,
      values: theme.customCSS,
    }))

    await $SystemAPI.settingsUpdate({
      values: [
        { name: 'ui.studio.themes', value: themeValues },
        { name: 'ui.studio.custom-css', value: customCssValues },
      ],
    })

    // Apply theme immediately so admin sees changes without refreshing
    setThemes(themeValues)
    const isDark = document.documentElement.classList.contains('dark')
    useTheme(isDark ? 'dark' : 'light')

    $toast.toastSuccess(t('notification.settings.theming.update.success'))
  } catch (e) {
    $toast.toastErrorHandler(t('notification.settings.theming.update.error'))(e)
  } finally {
    saving.value = false
  }
}

onMounted(() => loadSettings())
</script>
