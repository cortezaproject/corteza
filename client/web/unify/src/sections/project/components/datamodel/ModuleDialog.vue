<template>
  <Dialog
    v-model:visible="visible"
    modal
    :style="{ width: '72rem' }"
    :pt="{ content: { class: '!pt-2' }, footer: { class: 'flex justify-end gap-2' } }"
  >
    <template #header>
      <div class="flex items-center gap-2.5 min-w-0">
        <span
          class="inline-flex items-center justify-center w-8 h-8 rounded-md ring-1 shrink-0"
          :class="[cfg.bg, cfg.ring]"
        >
          <i :class="[cfg.icon, cfg.text]" />
        </span>
        <div class="min-w-0">
          <div class="text-[10px] uppercase tracking-wider text-muted-color leading-none mb-0.5">
            Module
          </div>
          <div class="font-semibold truncate leading-tight">{{ draft.name || 'Unnamed' }}</div>
        </div>
      </div>
    </template>

    <div class="flex flex-col gap-5">
      <div class="grid grid-cols-1 sm:grid-cols-2 gap-x-5 gap-y-4">
        <CFormGroup label="Name">
          <InputText v-model="draft.name" size="small" fluid :disabled="readonly" autofocus />
        </CFormGroup>
        <CFormGroup label="Data sensitivity" description="Default for new fields">
          <Select
            v-model="draft.sensitivity"
            :options="SENSITIVITY_OPTIONS"
            option-label="label"
            option-value="id"
            size="small"
            fluid
            :disabled="readonly"
          />
        </CFormGroup>
      </div>

      <CFormGroup label="Fields">
        <template #actions>
          <Button
            v-if="!readonly"
            icon="pi pi-plus"
            label="Add field"
            severity="secondary"
            size="small"
            @click="addDraftField"
          />
        </template>
        <FieldsEditor v-model="draft.fields" :module-options="moduleOptions" :disabled="readonly" />
      </CFormGroup>
    </div>

    <template #footer>
      <Button v-if="readonly" label="Close" size="small" @click="visible = false" />
      <template v-else>
        <Button label="Cancel" severity="secondary" outlined size="small" @click="visible = false" />
        <Button label="Save" size="small" :disabled="!draft.name.trim()" :loading="saving" @click="onSave" />
      </template>
    </template>
  </Dialog>
</template>

<script setup>
import FieldsEditor from '@/sections/project/components/datamodel/FieldsEditor.vue'
import { kindConfig } from '@/sections/project/config/kinds'
import { SENSITIVITY_OPTIONS } from '@/sections/project/config/sensitivity'
import { useProjectsStore } from '@/sections/project/stores/projects'
import { useToast } from 'primevue/usetoast'
import { computed, reactive, ref, watch } from 'vue'

const props = defineProps({
  modelValue: { type: Boolean, default: false },
  project: { type: Object, required: true },
  // Existing module to edit; null opens the dialog in create mode.
  moduleId: { type: String, default: null },
  // Open as read-only (e.g. approver view, or a submitted/approved step).
  readonly: { type: Boolean, default: false },
})

const emit = defineEmits(['update:modelValue'])

const store = useProjectsStore()
const toast = useToast()

const visible = computed({
  get: () => props.modelValue,
  set: v => emit('update:modelValue', v),
})

const cfg = kindConfig('module')

const module = computed(() =>
  props.moduleId ? (props.project?.resources || []).find(r => r.id === props.moduleId) : null,
)
const moduleOptions = computed(() =>
  (props.project?.resources || []).filter(r => r.kind === 'module' && r.id !== props.moduleId),
)

// --- Draft (staged) edits; nothing persists until Save -----------------------
const draft = reactive({ name: '', fields: [], sensitivity: null })
// Once a new module is created on Save, remember its id so a re-save (before
// the dialog closes) updates it instead of creating another copy.
const createdId = ref(null)

function initDraft() {
  createdId.value = null
  const m = module.value
  draft.name = m?.name || ''
  draft.fields = (m?.fields || []).map(f => ({ ...f }))
  draft.sensitivity = m?.sensitivity || null
}

watch(
  () => [props.modelValue, props.moduleId],
  () => {
    if (props.modelValue) initDraft()
  },
  { immediate: true },
)

const nid = () => `f-${Math.random().toString(36).slice(2, 9)}`
const addDraftField = () => {
  // New fields inherit the module's sensitivity as their default.
  draft.fields.push({
    id: nid(),
    name: '',
    type: 'String',
    required: false,
    targetModuleId: null,
    sensitivity: draft.sensitivity || null,
  })
}

// --- Commit -------------------------------------------------------------------
const saving = ref(false)

async function onSave() {
  if (!draft.name.trim() || saving.value) return
  const pid = props.project.id
  saving.value = true
  try {
    let id = props.moduleId || createdId.value
    if (!id) {
      // createdId guards against duplicates when a later write fails and the
      // user re-saves: the module exists, so the retry only updates it.
      id = await store.addResource(pid, { kind: 'module', name: draft.name })
      createdId.value = id
    }
    if (!id) return
    await store.setFields(pid, id, draft.fields)
    await store.updateResource(pid, id, { name: draft.name, sensitivity: draft.sensitivity || '' })
    visible.value = false
  } catch (err) {
    // Stay open so the staged edits aren't lost.
    toast.add({ severity: 'error', summary: 'Could not save module', detail: err.message, life: 4000 })
  } finally {
    saving.value = false
  }
}
</script>
