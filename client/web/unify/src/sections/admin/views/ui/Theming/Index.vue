<template>
  <Teleport to="#topbar-title" defer>
    <span>{{ $t('ui.settings.editor.human-studio.title') }}</span>
  </Teleport>

  <div v-if="loading" class="flex items-center justify-center h-full">
    <ProgressSpinner />
  </div>

  <div v-else class="flex flex-col h-full">
    <CViewContainer scroll gap="5">
      <Panel
        :header="$t('ui.settings.editor.human-studio.title')"
        toggleable
        :collapsed="false"
        class="shadow"
        :pt="{ content: { style: 'padding: 0 !important' } }"
      >
        <Tabs value="light">
          <TabList>
            <Tab v-for="theme in themes" :key="theme.id" :value="theme.id">
              {{ theme.title }}
            </Tab>
          </TabList>

          <TabPanels>
            <TabPanel v-for="theme in themes" :key="theme.id" :value="theme.id">
              <!-- Logos for this theme, previewed on its sidebar colour -->
              <div v-if="theme.logos" class="grid grid-cols-1 md:grid-cols-2 gap-6 mb-6">
                <CFormGroup
                  v-for="kind in logoKinds"
                  :key="kind"
                  :label="$t(`ui.settings.editor.human-studio.${kind}Logo.title`)"
                  :description="$t(`ui.settings.editor.human-studio.logo.previewSurface.${kind}`)"
                >
                  <CFileDropZone
                    accept="image/*"
                    :uploading="theme.logos[kind].uploading"
                    :error="theme.logos[kind].error"
                    :preview-url="theme.logos[kind].url || inheritedSlot(theme, kind)?.url"
                    :preview-style="{ backgroundColor: '#' + theme.variables[logoSurface[kind]] }"
                    :clearable="theme.logos[kind].custom"
                    :drop-label="
                      $t(`ui.settings.editor.human-studio.${kind}Logo.uploader.instructions`)
                    "
                    :uploading-label="
                      $t(`ui.settings.editor.human-studio.${kind}Logo.uploader.uploading`)
                    "
                    :label="$t(`ui.settings.editor.human-studio.${kind}Logo.title`)"
                    compact
                    preview-max-width="100%"
                    preview-max-height="200px"
                    @select="onLogoSelect(theme.logos[kind], $event)"
                    @clear="onLogoClear(theme.logos[kind])"
                  />
                </CFormGroup>
              </div>

              <!-- Color variables for light/dark tabs only -->
              <div v-if="theme.id !== 'general'" class="grid grid-cols-1 lg:grid-cols-2 gap-4 mb-6">
                <CFormGroup
                  v-for="key in themeVariableKeys"
                  :key="key"
                  :label="$t(`ui.settings.editor.human-studio.theme.variables.${key}.label`)"
                  :description="
                    $t(`ui.settings.editor.human-studio.theme.variables.${key}.description`)
                  "
                >
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
                </CFormGroup>
              </div>

              <!-- Custom CSS for all tabs -->
              <CFormGroup :label="$t('ui.settings.editor.human-studio.custom-css')">
                <Textarea v-model="theme.customCSS" rows="16" class="w-full font-mono text-sm" />
              </CFormGroup>
            </TabPanel>
          </TabPanels>
        </Tabs>
      </Panel>
    </CViewContainer>

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
import {
  currentTheme,
  setThemes,
  useTheme,
  components,
  useFileUpload,
} from '@planetcrust/human-vue'

const { CFileDropZone, CInputColorPicker, CViewContainer } = components

const { t } = useI18n()

const $toast = inject('$toast')
const $SystemAPI = inject('$SystemAPI')
const $Settings = inject('$Settings')

const loading = ref(false)
const saving = ref(false)

// One upload slot per theme and kind. Uploads address the setting itself,
// kebab-cased (`ui.main-logo-dark`); the structured settings payload reads it
// back JSON-cased (`ui.mainLogoDark`).
function logoSlot(key, read) {
  const { uploading, uploadError, uploadFileRaw, reset } = useFileUpload()
  return reactive({
    key,
    read,
    url: '',
    custom: false,
    uploading,
    error: uploadError,
    uploadFileRaw,
    reset,
  })
}

const logoKinds = ['main', 'icon']

// The surface each image is shown on: the logo in the sidebar, the icon in the topbar
const logoSurface = { main: 'sidebar-bg', icon: 'topbar-bg' }

// The light slot a dark slot inherits from while it has no image of its own
function inheritedSlot(theme, kind) {
  if (theme.id !== 'dark' || theme.logos[kind].url) return null
  return themes.find(t => t.id === 'light').logos[kind]
}

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
    id: 'light',
    title: t('ui.settings.editor.human-studio.tabs.light'),
    variables: { ...lightModeDefaults },
    defaultVariables: { ...lightModeDefaults },
    customCSS: '',
    logos: {
      main: logoSlot('ui.main-logo', 'ui.mainLogo'),
      icon: logoSlot('ui.icon-logo', 'ui.iconLogo'),
    },
  },
  {
    id: 'dark',
    title: t('ui.settings.editor.human-studio.tabs.dark'),
    variables: { ...darkModeDefaults },
    defaultVariables: { ...darkModeDefaults },
    customCSS: '',
    logos: {
      main: logoSlot('ui.main-logo-dark', 'ui.mainLogoDark'),
      icon: logoSlot('ui.icon-logo-dark', 'ui.iconLogoDark'),
    },
  },
  {
    id: 'general',
    title: t('ui.settings.editor.human-studio.tabs.general'),
    variables: {},
    defaultVariables: {},
    customCSS: '',
    logos: null,
  },
])

function stripHash(color) {
  if (!color) return ''
  return color.replace(/^#/, '')
}

function refreshLogoUrls() {
  for (const theme of themes) {
    for (const slot of Object.values(theme.logos || {})) {
      slot.url = $Settings.attachment(slot.read) || ''

      // Clearable only when a custom image was uploaded
      const raw = $Settings.get(slot.read, '')
      slot.custom = typeof raw === 'string' && raw.startsWith('attachment:')
    }
  }
}

async function onLogoSelect(slot, files) {
  const file = files[0]
  if (!file) return

  try {
    const endpoint = $SystemAPI.baseURL + $SystemAPI.settingsSetEndpoint({ key: slot.key })
    const token = $SystemAPI.accessTokenFn ? $SystemAPI.accessTokenFn() : ''
    await slot.uploadFileRaw(file, { url: endpoint, token })
    await $Settings.fetch()
    refreshLogoUrls()
    $toast.toastSuccess(t('notification.settings.theming.update.success'))
  } catch {
    // slot.error is set by the upload composable
  }
}

async function onLogoClear(slot) {
  try {
    await $SystemAPI.settingsUpdate({
      values: [{ name: slot.key, value: null }],
    })
    await $Settings.fetch()
    refreshLogoUrls()
    slot.reset()
    $toast.toastSuccess(t('notification.settings.theming.update.success'))
  } catch (e) {
    $toast.toastErrorHandler(t('notification.settings.theming.update.error'))(e)
  }
}

async function loadSettings() {
  loading.value = true
  try {
    // The boot snapshot may predate a server restart or another admin's upload
    await $Settings.fetch()
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
    useTheme(currentTheme.value)

    $toast.toastSuccess(t('notification.settings.theming.update.success'))
  } catch (e) {
    $toast.toastErrorHandler(t('notification.settings.theming.update.error'))(e)
  } finally {
    saving.value = false
  }
}

onMounted(() => loadSettings())
</script>
