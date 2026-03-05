<template>
  <Teleport to="#topbar-title" defer>
    <span>{{ $t('compose.settings.editor.title', 'Compose Settings') }}</span>
  </Teleport>

  <div v-if="loading" class="flex items-center justify-center h-full">
    <ProgressSpinner />
  </div>

  <div v-else class="flex flex-col h-full">
    <div class="container mx-auto p-4 flex-1 flex flex-col min-h-0 gap-4 overflow-y-auto">
      <Panel
        v-for="group in settingGroups"
        :key="group.key"
        :header="group.label"
        toggleable
        :collapsed="false"
        class="shadow"
      >
        <div class="flex flex-col divide-y">
          <div
            v-for="setting in group.settings"
            :key="setting.name"
            class="flex items-center justify-between py-3 gap-4"
          >
            <div class="flex flex-col flex-1 min-w-0">
              <span class="font-medium text-sm">{{ setting.label }}</span>
              <span class="text-xs text-surface-500 font-mono">{{ setting.name }}</span>
            </div>

            <div class="flex-shrink-0 w-64">
              <ToggleSwitch
                v-if="typeof setting.value === 'boolean'"
                v-model="setting.value"
                @update:modelValue="markDirty(setting)"
              />
              <InputNumber
                v-else-if="typeof setting.value === 'number'"
                v-model="setting.value"
                size="small"
                class="w-full"
                @update:modelValue="markDirty(setting)"
              />
              <InputText
                v-else
                v-model="setting.value"
                size="small"
                class="w-full"
                @update:modelValue="markDirty(setting)"
              />
            </div>
          </div>

          <div v-if="group.settings.length === 0" class="py-6 text-center text-surface-500">
            {{ $t('compose.settings.empty', 'No settings in this group.') }}
          </div>
        </div>
      </Panel>

      <div v-if="settingGroups.length === 0 && !loading" class="p-6 text-center text-surface-500">
        {{ $t('compose.settings.noSettings', 'No compose settings found.') }}
      </div>
    </div>

    <div class="shrink-0 border-t border-surface bg-surface">
      <div class="p-3 flex items-center justify-end">
        <Button
          :label="$t('general.label.save')"
          icon="pi pi-save"
          :loading="saving"
          :disabled="dirtySettings.size === 0"
          @click="handleSave"
        />
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed, inject, onMounted, ref } from 'vue'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()

const $toast = inject('$toast')
const $SystemAPI = inject('$SystemAPI')

const loading = ref(false)
const saving = ref(false)
const rawSettings = ref([])
const dirtySettings = ref(new Set())

// Human-readable labels for known compose setting keys
const settingLabels = computed(() => ({
  'compose.ui.record-toolbar-buttons.edit': t(
    'compose.settings.editor.label.toolbar-edit',
    'Show Edit button in record toolbar',
  ),
  'compose.ui.record-toolbar-buttons.delete': t(
    'compose.settings.editor.label.toolbar-delete',
    'Show Delete button in record toolbar',
  ),
  'compose.ui.record-toolbar-buttons.new': t(
    'compose.settings.editor.label.toolbar-new',
    'Show New button in record toolbar',
  ),
  'compose.ui.record-toolbar-buttons.clone': t(
    'compose.settings.editor.label.toolbar-clone',
    'Show Clone button in record toolbar',
  ),
  'compose.ui.record-toolbar-buttons.back': t(
    'compose.settings.editor.label.toolbar-back',
    'Show Back button in record toolbar',
  ),
  'compose.namespace-switch-enabled': t(
    'compose.settings.editor.label.namespace-switch',
    'Namespace switcher enabled',
  ),
  'compose.page-layout-enabled': t(
    'compose.settings.editor.label.page-layout',
    'Page layout enabled',
  ),
  'compose.record-revisions-enabled': t(
    'compose.settings.editor.label.record-revisions',
    'Record revisions enabled',
  ),
  'compose.default-namespace': t(
    'compose.settings.editor.label.default-namespace',
    'Default namespace slug',
  ),
}))

// Sub-prefix → panel header label mapping
const groupLabels = computed(() => ({
  'compose.ui.record-toolbar-buttons': t(
    'compose.settings.editor.group.record-toolbar',
    'Record Toolbar',
  ),
  'compose.ui': t('compose.settings.editor.group.ui', 'UI'),
  compose: t('compose.settings.editor.group.general', 'General'),
}))

const settingGroups = computed(() => {
  const groups = new Map()

  for (const setting of rawSettings.value) {
    const prefix = getGroupPrefix(setting.name)
    if (!groups.has(prefix)) {
      groups.set(prefix, {
        key: prefix,
        label: groupLabels.value[prefix] || capitalize(prefix.replace('compose.', '')),
        settings: [],
      })
    }
    groups.get(prefix).settings.push(setting)
  }

  return Array.from(groups.values())
})

function getGroupPrefix(name) {
  const sorted = Object.keys(groupLabels.value).sort((a, b) => b.length - a.length)
  for (const prefix of sorted) {
    if (name.startsWith(prefix + '.') || name === prefix) {
      return prefix
    }
  }
  return 'compose'
}

function capitalize(str) {
  return str.charAt(0).toUpperCase() + str.slice(1)
}

function markDirty(setting) {
  dirtySettings.value = new Set([...dirtySettings.value, setting.name])
}

async function loadSettings() {
  loading.value = true
  try {
    const result = await $SystemAPI.settingsList({ prefix: 'compose.' })
    rawSettings.value = (result || []).map(s => ({
      name: s.name,
      label: settingLabels.value[s.name] || s.name.replace('compose.', '').split('.').pop(),
      value: parseValue(s.value),
    }))
  } catch (e) {
    $toast.toastErrorHandler(t('notification.settings.fetch.error', 'Failed to load settings'))(e)
  } finally {
    loading.value = false
  }
}

function parseValue(raw) {
  if (raw === null || raw === undefined) return ''
  if (typeof raw === 'boolean') return raw
  if (typeof raw === 'number') return raw
  if (typeof raw === 'object') {
    if ('@value' in raw) return raw['@value']
    return JSON.stringify(raw)
  }
  return raw
}

async function handleSave() {
  saving.value = true
  try {
    const values = rawSettings.value
      .filter(s => dirtySettings.value.has(s.name))
      .map(s => ({ name: s.name, value: s.value }))

    await $SystemAPI.settingsUpdate({ values })
    dirtySettings.value = new Set()
    $toast.toastSuccess(t('notification.settings.update.success', 'Settings saved'))
  } catch (e) {
    $toast.toastErrorHandler(t('notification.settings.update.error', 'Failed to save settings'))(e)
  } finally {
    saving.value = false
  }
}

onMounted(() => loadSettings())
</script>
