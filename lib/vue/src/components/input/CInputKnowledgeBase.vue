<template>
  <div class="flex flex-col gap-3">
    <!-- Select + Create button -->
    <div class="flex gap-2">
      <Select
        v-model="pickerSelection"
        :options="availableOptions"
        :option-label="getOptionLabel"
        :placeholder="placeholder"
        :disabled="disabled"
        :loading="loading"
        class="flex-1"
        filter
        fluid
        @update:model-value="onPickerSelect"
        @show="onShow"
      >
        <template #option="{ option }">
          <div class="flex flex-col">
            <span>{{ option.title || option.handle || option.knowledgeBaseID }}</span>
            <small v-if="option.description" class="text-muted-color truncate max-w-64">
              {{ option.description }}
            </small>
          </div>
        </template>
      </Select>
      <Button
        icon="pi pi-plus"
        :aria-label="createLabel"
        v-tooltip.top="createLabel"
        @click="openCreateDialog"
        :disabled="disabled"
      />
    </div>

    <!-- Selected KBs as rows -->
    <div v-if="selectedEntries.length" class="flex flex-col gap-2">
      <div
        v-for="entry in selectedEntries"
        :key="entry.kb.knowledgeBaseID"
        class="flex items-center gap-3 p-3 border border-surface rounded-lg"
      >
        <div class="flex flex-col gap-0.5 flex-1 min-w-0">
          <span class="font-medium text-color text-sm truncate">
            {{ entry.kb.title || entry.kb.handle || entry.kb.knowledgeBaseID }}
          </span>
          <small v-if="entry.kb.description" class="text-muted-color text-xs truncate">
            {{ entry.kb.description }}
          </small>
        </div>

        <div class="flex items-center gap-1 shrink-0">
          <Button
            icon="pi pi-pencil"
            severity="secondary"
            text
            size="small"
            @click="openEditDialog(entry.kb)"
          />
          <Button
            icon="pi pi-trash"
            severity="danger"
            text
            size="small"
            @click="removeEntry(entry)"
          />
        </div>
      </div>
    </div>

    <!-- Create/Edit KB Dialog -->
    <Dialog
      v-model:visible="dialogVisible"
      :header="editingKB ? editLabel : createDialogLabel"
      modal
      :style="{ width: '48rem' }"
      :closable="!saving"
    >
      <div class="flex flex-col gap-4">
        <!-- Title -->
        <div class="flex flex-col gap-1">
          <label for="kb-title" class="font-medium text-primary text-sm">
            {{ titleLabel }}
          </label>
          <InputText id="kb-title" v-model="dialogForm.title" />
        </div>

        <!-- Description -->
        <div class="flex flex-col gap-1">
          <label for="kb-description" class="font-medium text-primary text-sm">
            {{ descriptionLabel }}
          </label>
          <Textarea id="kb-description" v-model="dialogForm.description" rows="4" autoResize />
          <small class="text-muted-color">
            {{ descriptionHelp }}
          </small>
        </div>

        <!-- Compose Context -->
        <div class="flex flex-col gap-2">
          <label class="font-medium text-primary text-sm">
            {{ composeContextLabel }}
          </label>
          <div class="flex flex-col gap-3">
            <div
              v-for="(nsCtx, idx) in dialogForm.context.namespaces"
              :key="idx"
              class="flex items-start gap-2 border border-surface rounded-lg p-3"
            >
              <div class="flex flex-col gap-2 flex-1">
                <CInputNamespace
                  :model-value="nsCtx.namespaceID"
                  @update:model-value="onNamespaceChange(nsCtx, $event)"
                  :placeholder="namespacePlaceholder"
                />
                <MultiSelect
                  v-if="nsCtx.namespaceID"
                  :model-value="getModuleObjects(nsCtx)"
                  @update:model-value="onModuleSelect(nsCtx, $event)"
                  :options="nsCtx._moduleOptions || []"
                  option-label="name"
                  :placeholder="modulesPlaceholder"
                  display="chip"
                  filter
                  fluid
                  :loading="nsCtx._loadingModules"
                  @show="fetchModulesForNamespace(nsCtx)"
                />
              </div>
              <Button
                icon="pi pi-trash"
                severity="danger"
                text
                rounded
                size="small"
                @click="removeNamespaceContext(idx)"
              />
            </div>

            <Button
              :label="addNamespaceLabel"
              icon="pi pi-plus"
              severity="secondary"
              outlined
              size="small"
              @click="addNamespaceContext"
            />
          </div>
        </div>
      </div>

      <template #footer>
        <div class="flex items-center justify-between w-full">
          <CInputDelete
            v-if="editingKB && editingKB.canDeleteKnowledgeBase !== false"
            :label="deleteLabel"
            :message="deleteMessage"
            :header="editingKB?.title || editingKB?.handle || deleteLabel"
            :disabled="saving"
            size="small"
            @confirm="deleteKnowledgeBase"
          />
          <span v-else />
          <div class="flex gap-2">
            <Button
              :label="cancelLabel"
              severity="secondary"
              outlined
              size="small"
              @click="dialogVisible = false"
              :disabled="saving"
            />
            <Button
              :label="saveLabel"
              severity="primary"
              size="small"
              @click="saveKnowledgeBase"
              :loading="saving"
            />
          </div>
        </div>
      </template>
    </Dialog>
  </div>
</template>

<script setup>
import { computed, inject, nextTick, onBeforeUnmount, onMounted, ref, watch } from 'vue'
import CInputNamespace from './CInputNamespace.vue'
import CInputDelete from './CInputDelete.vue'

defineOptions({ inheritAttrs: false })

const props = defineProps({
  modelValue: {
    type: Array,
    default: () => [],
  },
  placeholder: {
    type: String,
    default: 'Select knowledge base to add...',
  },
  disabled: {
    type: Boolean,
    default: false,
  },
  // Labels (i18n injected from parent)
  createLabel: { type: String, default: 'Create new' },
  createDialogLabel: { type: String, default: 'Create Knowledge Base' },
  editLabel: { type: String, default: 'Edit Knowledge Base' },
  titleLabel: { type: String, default: 'Title' },
  descriptionLabel: { type: String, default: 'Description' },
  descriptionHelp: {
    type: String,
    default:
      'Free-text knowledge the agent should know about. Supports any text — policies, domain context, instructions, etc.',
  },
  composeContextLabel: { type: String, default: 'Compose Context' },
  namespacePlaceholder: { type: String, default: 'Select namespace...' },
  modulesPlaceholder: { type: String, default: 'Select modules...' },
  addNamespaceLabel: { type: String, default: 'Add namespace' },
  saveLabel: { type: String, default: 'Save' },
  cancelLabel: { type: String, default: 'Cancel' },
  deleteLabel: { type: String, default: 'Delete' },
  deleteMessage: { type: String, default: 'Are you sure you want to permanently delete this knowledge base?' },
})

const emit = defineEmits(['update:modelValue'])

const $SystemAPI = inject('$SystemAPI')
const $ComposeAPI = inject('$ComposeAPI')

const allKnowledgeBases = ref([])
const loading = ref(false)
const pickerSelection = ref(null)

// Each entry: { kb: <full KB object> }
const selectedEntries = ref([])

// Dialog state
const dialogVisible = ref(false)
const editingKB = ref(null)
const saving = ref(false)
const dialogForm = ref(emptyForm())

let cancelCurrentRequest = null

function emptyForm() {
  return {
    title: '',
    description: '',
    context: { namespaces: [] },
  }
}

function getOptionLabel(kb) {
  if (!kb) return ''
  return kb.title || kb.handle || kb.knowledgeBaseID
}

// Options not yet selected
const availableOptions = computed(() => {
  const selectedIds = new Set(selectedEntries.value.map(e => e.kb.knowledgeBaseID))
  return allKnowledgeBases.value.filter(kb => !selectedIds.has(kb.knowledgeBaseID))
})

// --- Fetching KBs ---

async function fetchKnowledgeBases() {
  if (cancelCurrentRequest) {
    cancelCurrentRequest()
    cancelCurrentRequest = null
  }

  loading.value = true
  try {
    const { response, cancel } = $SystemAPI.knowledgeBaseListCancellable({
      limit: 100,
      sort: 'title ASC',
    })
    cancelCurrentRequest = cancel

    const result = await response()
    allKnowledgeBases.value = Array.isArray(result) ? result : result.set || []
  } catch (e) {
    if (e?.message !== 'canceled') {
      allKnowledgeBases.value = []
    }
  } finally {
    loading.value = false
    cancelCurrentRequest = null
  }
}

function onShow() {
  if (allKnowledgeBases.value.length === 0) {
    fetchKnowledgeBases()
  }
}

// --- Selection handling ---

function onPickerSelect(kb) {
  if (!kb) return
  if (!selectedEntries.value.find(e => e.kb.knowledgeBaseID === kb.knowledgeBaseID)) {
    selectedEntries.value.push({ kb })
    emitValue()
  }
  // Clear the picker
  pickerSelection.value = null
}

function removeEntry(entry) {
  selectedEntries.value = selectedEntries.value.filter(
    e => e.kb.knowledgeBaseID !== entry.kb.knowledgeBaseID,
  )
  emitValue()
  nextTick(() => {
    pickerSelection.value = null
  })
}

function emitValue() {
  emit(
    'update:modelValue',
    selectedEntries.value.map(e => e.kb.knowledgeBaseID),
  )
}

// Sync from modelValue (e.g. on load)
function syncFromModelValue() {
  const ids = props.modelValue || []
  if (ids.length === 0) {
    selectedEntries.value = []
    return
  }

  // Build entries from known KBs
  const newEntries = []
  for (const id of ids) {
    // Already in entries?
    const existing = selectedEntries.value.find(
      e => e.kb.knowledgeBaseID === id || e.kb.knowledgeBaseID === String(id),
    )
    if (existing) {
      newEntries.push(existing)
      continue
    }

    // Find in loaded KBs
    const kb = allKnowledgeBases.value.find(
      k => k.knowledgeBaseID === id || k.knowledgeBaseID === String(id),
    )
    if (kb) {
      newEntries.push({ kb })
    } else {
      // Need to load individually
      loadKnowledgeBaseById(id)
    }
  }

  selectedEntries.value = newEntries
}

async function loadKnowledgeBaseById(kbID) {
  if (!kbID || kbID === '0' || !$SystemAPI) return
  try {
    const kb = await $SystemAPI.knowledgeBaseRead({ knowledgeBaseID: kbID })
    if (kb) {
      if (!allKnowledgeBases.value.find(o => o.knowledgeBaseID === kb.knowledgeBaseID)) {
        allKnowledgeBases.value = [...allKnowledgeBases.value, kb]
      }
      if (!selectedEntries.value.find(e => e.kb.knowledgeBaseID === kb.knowledgeBaseID)) {
        selectedEntries.value.push({ kb })
      }
    }
  } catch {
    // KB not found
  }
}

// --- Dialog ---

function openCreateDialog() {
  editingKB.value = null
  dialogForm.value = emptyForm()
  dialogVisible.value = true
}

function openEditDialog(kb) {
  editingKB.value = kb
  dialogForm.value = {
    title: kb.title || '',
    description: kb.description || '',
    context: kb.context ? JSON.parse(JSON.stringify(kb.context)) : { namespaces: [] },
  }
  // Pre-fetch modules for each namespace context
  dialogForm.value.context.namespaces.forEach(nsCtx => {
    nsCtx._moduleOptions = []
    nsCtx._loadingModules = false
    if (nsCtx.namespaceID) {
      fetchModulesForNamespace(nsCtx)
    }
  })
  dialogVisible.value = true
}

async function saveKnowledgeBase() {
  saving.value = true
  try {
    const context = {
      namespaces: dialogForm.value.context.namespaces
        .filter(ns => ns.namespaceID)
        .map(ns => ({
          namespaceID: ns.namespaceID,
          moduleIDs: ns.moduleIDs || [],
        })),
    }

    let kb
    if (editingKB.value) {
      kb = await $SystemAPI.knowledgeBaseUpdate({
        knowledgeBaseID: editingKB.value.knowledgeBaseID,
        title: dialogForm.value.title,
        description: dialogForm.value.description,
        context: context.namespaces.length > 0 ? context : undefined,
        updatedAt: editingKB.value.updatedAt,
      })

      // Update in allKnowledgeBases
      const idx = allKnowledgeBases.value.findIndex(o => o.knowledgeBaseID === kb.knowledgeBaseID)
      if (idx >= 0) allKnowledgeBases.value[idx] = kb

      // Update in selected entries
      const entry = selectedEntries.value.find(e => e.kb.knowledgeBaseID === kb.knowledgeBaseID)
      if (entry) entry.kb = kb
    } else {
      kb = await $SystemAPI.knowledgeBaseCreate({
        title: dialogForm.value.title,
        description: dialogForm.value.description,
        context: context.namespaces.length > 0 ? context : undefined,
      })

      // Add to list and auto-select
      allKnowledgeBases.value = [...allKnowledgeBases.value, kb]
      selectedEntries.value.push({ kb })
      emitValue()
    }

    dialogVisible.value = false
  } catch {
    // silent
  } finally {
    saving.value = false
  }
}

async function deleteKnowledgeBase() {
  if (!editingKB.value) return
  saving.value = true
  try {
    await $SystemAPI.knowledgeBaseDelete({
      knowledgeBaseID: editingKB.value.knowledgeBaseID,
    })

    // Remove from options and selection
    const id = editingKB.value.knowledgeBaseID
    allKnowledgeBases.value = allKnowledgeBases.value.filter(o => o.knowledgeBaseID !== id)
    selectedEntries.value = selectedEntries.value.filter(e => e.kb.knowledgeBaseID !== id)
    emitValue()

    dialogVisible.value = false
  } catch {
    // silent
  } finally {
    saving.value = false
  }
}

// --- Compose Context helpers ---

function addNamespaceContext() {
  dialogForm.value.context.namespaces.push({
    namespaceID: null,
    moduleIDs: [],
    _moduleOptions: [],
    _loadingModules: false,
  })
}

function removeNamespaceContext(idx) {
  dialogForm.value.context.namespaces.splice(idx, 1)
}

async function fetchModulesForNamespace(nsCtx) {
  if (!nsCtx.namespaceID || !$ComposeAPI) return

  nsCtx._loadingModules = true
  try {
    const result = await $ComposeAPI.moduleList({
      namespaceID: nsCtx.namespaceID,
      limit: 100,
      sort: 'name ASC',
    })
    nsCtx._moduleOptions = result.set || []
  } catch {
    nsCtx._moduleOptions = []
  } finally {
    nsCtx._loadingModules = false
  }
}

function getModuleObjects(nsCtx) {
  if (!nsCtx.moduleIDs?.length || !nsCtx._moduleOptions?.length) return []
  return nsCtx._moduleOptions.filter(
    m => nsCtx.moduleIDs.includes(m.moduleID) || nsCtx.moduleIDs.includes(String(m.moduleID)),
  )
}

function onNamespaceChange(nsCtx, namespaceID) {
  nsCtx.namespaceID = namespaceID
  nsCtx.moduleIDs = []
}

function onModuleSelect(nsCtx, modules) {
  nsCtx.moduleIDs = modules.map(m => m.moduleID)
}

// --- Watchers ---

watch(
  () => props.modelValue,
  () => syncFromModelValue(),
  { deep: true },
)

watch(
  () => allKnowledgeBases.value.length,
  () => {
    if (props.modelValue?.length && allKnowledgeBases.value.length) {
      syncFromModelValue()
    }
  },
)

onMounted(() => {
  fetchKnowledgeBases()
})

onBeforeUnmount(() => {
  if (cancelCurrentRequest) {
    cancelCurrentRequest()
  }
})
</script>
