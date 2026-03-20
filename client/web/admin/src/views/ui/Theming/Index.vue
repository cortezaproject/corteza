<template>
  <Teleport to="#topbar-title" defer>
    <span>{{ $t('ui.settings.editor.corteza-studio.title') }}</span>
  </Teleport>

  <div v-if="loading" class="flex items-center justify-center h-full">
    <ProgressSpinner />
  </div>

  <div v-else class="flex flex-col h-full">
    <div class="container mx-auto p-4 flex-1 flex flex-col min-h-0 gap-5 overflow-y-auto">
      <Panel
        :header="$t('ui.settings.editor.corteza-studio.title')"
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
                    {{ $t(`ui.settings.editor.corteza-studio.theme.variables.${key}.label`) }}
                  </label>
                  <span class="text-xs text-muted-color">
                    {{ $t(`ui.settings.editor.corteza-studio.theme.variables.${key}.description`) }}
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
                      :title="$t('ui.settings.editor.corteza-studio.label.default')"
                      @click="theme.variables[key] = theme.defaultVariables[key]"
                    />
                  </div>
                </div>
              </div>

              <!-- Custom CSS for all tabs -->
              <div class="flex flex-col gap-1">
                <label class="font-medium text-sm text-primary">
                  {{ $t('ui.settings.editor.corteza-studio.custom-css') }}
                </label>
                <Textarea v-model="theme.customCSS" rows="16" class="w-full font-mono text-sm" />
              </div>
            </TabPanel>
          </TabPanels>
        </Tabs>
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
import { inject, onMounted, reactive, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { setThemes, useTheme, components } from '@cortezaproject/corteza-vue-next'

const { CInputColorPicker } = components

const { t } = useI18n()

const $toast = inject('$toast')
const $SystemAPI = inject('$SystemAPI')

const loading = ref(false)
const saving = ref(false)

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
    title: t('ui.settings.editor.corteza-studio.tabs.general'),
    variables: {},
    defaultVariables: {},
    customCSS: '',
  },
  {
    id: 'light',
    title: t('ui.settings.editor.corteza-studio.tabs.light'),
    variables: { ...lightModeDefaults },
    defaultVariables: { ...lightModeDefaults },
    customCSS: '',
  },
  {
    id: 'dark',
    title: t('ui.settings.editor.corteza-studio.tabs.dark'),
    variables: { ...darkModeDefaults },
    defaultVariables: { ...darkModeDefaults },
    customCSS: '',
  },
])

function stripHash(color) {
  if (!color) return ''
  return color.replace(/^#/, '')
}

async function loadSettings() {
  loading.value = true
  try {
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
