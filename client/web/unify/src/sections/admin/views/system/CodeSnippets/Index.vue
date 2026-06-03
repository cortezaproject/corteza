<template>
  <Teleport to="#topbar-title" defer>
    <span>{{ $t('system.code-snippets.editor.title') }}</span>
  </Teleport>

  <div v-if="loading" class="flex items-center justify-center h-full">
    <ProgressSpinner />
  </div>

  <div v-else class="flex flex-col h-full">
    <div class="container mx-auto p-4 flex-1 flex flex-col min-h-0 gap-4 overflow-y-auto">
      <Panel
        :header="$t('system.code-snippets.editor.title')"
        toggleable
        :collapsed="false"
        class="shadow"
      >
        <CResourceTable
          :items="codeSnippets"
          :fields="snippetFields"
          :action-items="getSnippetActions"
          primary-key="name"
          :row-class="() => 'cursor-pointer'"
          @row-click="({ index }) => openEditor(index)"
        >
          <template #header>
            <Button
              :label="$t('system.code-snippets.editor.code-snippets.new')"
              icon="pi pi-plus"
              size="small"
              @click="openEditor()"
            />
          </template>
          <template #body-enabled="{ data }">
            <i
              :class="
                data.enabled ? 'pi pi-check text-green-500' : 'pi pi-times text-muted-color'
              "
            />
          </template>
        </CResourceTable>
      </Panel>
    </div>

    <!-- Edit Dialog -->
    <Dialog v-model:visible="modal.open" :header="modal.title" modal class="w-full max-w-3xl">
      <div class="flex flex-col gap-4">
        <CInputSwitch
          v-model="modal.data.enabled"
          :label="$t('system.code-snippets.editor.code-snippets.enabled')"
        />

        <div class="flex flex-col gap-1">
          <label class="font-medium text-sm">
            {{ $t('system.code-snippets.editor.code-snippets.form.name.label') }}
          </label>
          <InputText v-model="modal.data.name" class="w-full" />
        </div>

        <div class="flex flex-col gap-1">
          <label class="font-medium text-sm">
            {{ $t('system.code-snippets.editor.code-snippets.form.code.label') }}
          </label>
          <span class="text-xs text-muted-color">
            {{ $t('system.code-snippets.editor.code-snippets.form.code.description') }}
          </span>
          <Textarea v-model="modal.data.script" rows="12" class="w-full font-mono text-sm" />
        </div>
      </div>

      <template #footer>
        <div class="flex items-center w-full gap-2">
          <Button
            v-if="modal.index >= 0"
            :label="$t('general.label.delete')"
            severity="danger"
            size="small"
            @click="confirmDeleteSnippet"
          />
          <div class="flex-1" />
          <Button
            :label="$t('general.label.cancel')"
            severity="secondary"
            outlined
            size="small"
            @click="modal.open = false"
          />
          <Button
            :label="$t('general.label.save')"
            size="small"
            :disabled="!modal.data.name || !modal.data.script"
            @click="saveSnippet"
          />
        </div>
      </template>
    </Dialog>
  </div>
</template>

<script setup>
import { inject, onMounted, reactive, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { components } from '@planetcrust/human-vue'

const { CResourceTable } = components
const { t } = useI18n()

const $toast = inject('$toast')
const $SystemAPI = inject('$SystemAPI')

const loading = ref(false)
const codeSnippets = ref([])

const modal = reactive({
  open: false,
  index: -1,
  title: '',
  data: { name: '', script: '', enabled: true },
})

function openEditor(index) {
  if (index !== undefined && index >= 0) {
    modal.index = index
    modal.title =
      codeSnippets.value[index].name || t('system.code-snippets.editor.code-snippets.add')
    modal.data = { ...codeSnippets.value[index] }
  } else {
    modal.index = -1
    modal.title = t('system.code-snippets.editor.code-snippets.add')
    modal.data = { name: '', script: '<' + 'script> ' + '</' + 'script>', enabled: true }
  }
  modal.open = true
}

function saveSnippet() {
  if (modal.index >= 0) {
    codeSnippets.value.splice(modal.index, 1, { ...modal.data })
  } else {
    codeSnippets.value.push({ ...modal.data })
  }
  modal.open = false
  persistSnippets('update')
}

function deleteSnippet(index) {
  codeSnippets.value.splice(index, 1)
  persistSnippets('delete')
}

function confirmDeleteSnippet() {
  deleteSnippet(modal.index)
  modal.open = false
}

const snippetFields = [
  { key: 'name', header: t('system.code-snippets.editor.code-snippets.table-headers.name') },
  {
    key: 'enabled',
    header: t('system.code-snippets.editor.code-snippets.table-headers.enabled'),
    headerStyle: 'width: 6rem',
    headerClass: 'text-center',
    bodyClass: 'text-center',
  },
]

function getSnippetActions(data, index) {
  return [
    {
      label: t('general.label.delete'),
      icon: 'pi pi-trash',
      class: 'text-red-500',
      command: () => deleteSnippet(index),
    },
  ]
}

async function loadSettings() {
  loading.value = true
  try {
    const result = await $SystemAPI.settingsList({ prefix: 'code-snippets' })
    if (result && result[0] && result[0].value) {
      codeSnippets.value = result[0].value
    } else {
      codeSnippets.value = []
    }
  } catch (e) {
    $toast.toastErrorHandler(t('notification.settings.fetch.error'))(e)
  } finally {
    loading.value = false
  }
}

async function persistSnippets(action) {
  try {
    await $SystemAPI.settingsUpdate({
      values: [{ name: 'code-snippets', value: codeSnippets.value }],
    })

    if (action === 'delete') {
      $toast.toastSuccess(t('notification.settings.code-snippet.delete.success'))
    } else {
      $toast.toastSuccess(t('notification.settings.code-snippet.update.success'))
    }
  } catch (e) {
    $toast.toastErrorHandler(t('notification.settings.code-snippet.update.error'))(e)
  }
}

onMounted(() => loadSettings())
</script>
