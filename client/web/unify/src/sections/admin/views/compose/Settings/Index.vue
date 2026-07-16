<template>
  <Teleport to="#topbar-title" defer>
    <span>{{ $t('compose.settings.editor.title') }}</span>
  </Teleport>

  <div v-if="loading" class="flex items-center justify-center h-full">
    <ProgressSpinner />
  </div>

  <div v-else class="flex flex-col h-full">
    <CViewContainer scroll>
      <!-- Basic / Attachments -->
      <Panel
        :header="$t('compose.settings.editor.group.general')"
        toggleable
        :collapsed="false"
        class="shadow"
      >
        <!-- Page attachments -->
        <h5 class="text-sm font-semibold mb-2">
          {{ $t('compose.settings.editor.basic.attachments.page') }}
        </h5>
        <div class="grid grid-cols-1 md:grid-cols-2 gap-4 mb-4">
          <div class="flex flex-col gap-1">
            <label class="font-medium text-sm">
              {{ $t('compose.settings.editor.basic.attachments.max-size') }}
            </label>
            <InputNumber v-model="settings['compose.page.attachments.max-size']" class="w-full" />
          </div>
          <div class="flex flex-col gap-1">
            <label class="font-medium text-sm">
              {{ $t('compose.settings.editor.basic.attachments.type.whitelist') }}
            </label>
            <span class="text-xs text-muted-color">
              {{ $t('compose.settings.editor.basic.attachments.type.description') }}
            </span>
            <InputText v-model="pageAttachmentWhitelist" class="w-full" />
          </div>
        </div>

        <Divider />

        <!-- Record attachments -->
        <h5 class="text-sm font-semibold mb-2">
          {{ $t('compose.settings.editor.basic.attachments.record') }}
        </h5>
        <div class="grid grid-cols-1 md:grid-cols-2 gap-4 mb-4">
          <div class="flex flex-col gap-1">
            <label class="font-medium text-sm">
              {{ $t('compose.settings.editor.basic.attachments.max-size') }}
            </label>
            <InputNumber v-model="settings['compose.record.attachments.max-size']" class="w-full" />
          </div>
          <div class="flex flex-col gap-1">
            <label class="font-medium text-sm">
              {{ $t('compose.settings.editor.basic.attachments.type.whitelist') }}
            </label>
            <span class="text-xs text-muted-color">
              {{ $t('compose.settings.editor.basic.attachments.type.description') }}
            </span>
            <InputText v-model="recordAttachmentWhitelist" class="w-full" />
          </div>
        </div>

        <Divider />

        <!-- Icon attachments -->
        <h5 class="text-sm font-semibold mb-2">
          {{ $t('compose.settings.editor.basic.attachments.icon') }}
        </h5>
        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div class="flex flex-col gap-1">
            <label class="font-medium text-sm">
              {{ $t('compose.settings.editor.basic.attachments.max-size') }}
            </label>
            <InputNumber v-model="settings['compose.icon.attachments.max-size']" class="w-full" />
          </div>
          <div class="flex flex-col gap-1">
            <label class="font-medium text-sm">
              {{ $t('compose.settings.editor.basic.attachments.type.whitelist') }}
            </label>
            <span class="text-xs text-muted-color">
              {{ $t('compose.settings.editor.basic.attachments.type.description') }}
            </span>
            <InputText v-model="iconAttachmentWhitelist" class="w-full" />
          </div>
        </div>
      </Panel>

      <!-- User Interface -->
      <Panel
        :header="$t('compose.settings.editor.ui.title')"
        toggleable
        :collapsed="false"
        class="shadow"
      >
        <!-- Sidebar -->
        <h5 class="text-sm font-semibold mb-2">
          {{ $t('compose.settings.editor.ui.sidebar.title') }}
        </h5>
        <div class="grid grid-cols-1 md:grid-cols-2 gap-4 mb-4">
          <CInputSwitch
            v-model="sidebar.hideNamespaceList"
            :label="$t('compose.settings.editor.ui.sidebar.hide-namespace-list')"
          />
          <CInputSwitch
            v-model="sidebar.hideNamespaceListLink"
            :label="$t('compose.settings.editor.ui.sidebar.hide-namespace-list-link')"
          />
        </div>

        <Divider />

        <!-- Record Toolbar -->
        <h5 class="text-sm font-semibold mb-2">
          {{ $t('compose.settings.editor.ui.record-toolbar.title') }}
        </h5>
        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
          <CInputSwitch
            v-model="recordToolbar.hideSubmit"
            :label="$t('compose.settings.editor.ui.record-toolbar.hide-submit')"
          />
          <CInputSwitch
            v-model="recordToolbar.hideDelete"
            :label="$t('compose.settings.editor.ui.record-toolbar.hide-delete')"
          />
          <CInputSwitch
            v-model="recordToolbar.hideEdit"
            :label="$t('compose.settings.editor.ui.record-toolbar.hide-edit')"
          />
          <CInputSwitch
            v-model="recordToolbar.hideNew"
            :label="$t('compose.settings.editor.ui.record-toolbar.hide-new')"
          />
          <CInputSwitch
            v-model="recordToolbar.hideClone"
            :label="$t('compose.settings.editor.ui.record-toolbar.hide-clone')"
          />
          <CInputSwitch
            v-model="recordToolbar.hideBack"
            :label="$t('compose.settings.editor.ui.record-toolbar.hide-back')"
          />
        </div>
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
import { computed, inject, onMounted, reactive, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { components } from '@planetcrust/human-vue'

const { CViewContainer } = components
const { t } = useI18n()

const $toast = inject('$toast')
const $SystemAPI = inject('$SystemAPI')

const loading = ref(false)
const saving = ref(false)
const settings = reactive({})

// UI sub-objects
const sidebar = reactive({
  hideNamespaceList: false,
  hideNamespaceListLink: false,
})

const recordToolbar = reactive({
  hideSubmit: false,
  hideDelete: false,
  hideEdit: false,
  hideNew: false,
  hideClone: false,
  hideBack: false,
})

// Computed whitelist getters/setters (comma-separated string ↔ array)
const pageAttachmentWhitelist = computed({
  get: () => (settings['compose.page.attachments.mimetypes'] || []).join(','),
  set: val => {
    settings['compose.page.attachments.mimetypes'] = parseMimeTypes(val)
  },
})

const recordAttachmentWhitelist = computed({
  get: () => (settings['compose.record.attachments.mimetypes'] || []).join(','),
  set: val => {
    settings['compose.record.attachments.mimetypes'] = parseMimeTypes(val)
  },
})

const iconAttachmentWhitelist = computed({
  get: () => (settings['compose.icon.attachments.mimetypes'] || []).join(','),
  set: val => {
    settings['compose.icon.attachments.mimetypes'] = parseMimeTypes(val)
  },
})

function parseMimeTypes(value) {
  return (value || '')
    .split(',')
    .map(v => v.replace(/ /g, ''))
    .filter(v => v.match(/^[-\w.]+\/[-\w/+.]+$/g) !== null)
}

function parseValue(raw) {
  if (raw === null || raw === undefined) return ''
  if (typeof raw === 'boolean') return raw
  if (typeof raw === 'number') return raw
  if (typeof raw === 'object') {
    if ('@value' in raw) return raw['@value']
    return raw
  }
  return raw
}

async function loadSettings() {
  loading.value = true
  try {
    const result = await $SystemAPI.settingsList({ prefix: 'compose.' })
    for (const s of result || []) {
      settings[s.name] = parseValue(s.value)
    }

    // Hydrate UI sub-objects
    const sidebarData = settings['compose.ui.sidebar'] || {}
    Object.assign(sidebar, sidebarData)

    const toolbarData = settings['compose.ui.record-toolbar'] || {}
    Object.assign(recordToolbar, toolbarData)
  } catch (e) {
    $toast.toastErrorHandler(t('notification.settings.fetch.error'))(e)
  } finally {
    loading.value = false
  }
}

async function handleSave() {
  saving.value = true
  try {
    const values = Object.entries(settings).map(([name, value]) => ({
      name,
      value,
    }))

    // Merge UI sub-objects back
    values.push(
      { name: 'compose.ui.sidebar', value: { ...sidebar } },
      { name: 'compose.ui.record-toolbar', value: { ...recordToolbar } },
    )

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
