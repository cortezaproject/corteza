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
      <div class="grid grid-cols-1 sm:grid-cols-2 gap-x-5">
        <CFormGroup label="Name" required>
          <div>
            <InputText
              v-model="draft.name"
              size="small"
              fluid
              :disabled="readonly"
              :invalid="submitted && !!nameError"
              autofocus
            />
            <ValidationMessage :message="submitted ? nameError : ''" />
          </div>
        </CFormGroup>
      </div>
      <CFormGroup label="Description">
        <Textarea v-model="draft.description" rows="2" auto-resize fluid :disabled="readonly" />
      </CFormGroup>

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
        <FieldsEditor
          ref="fieldsEditorRef"
          v-model="draft.fields"
          :module-options="moduleOptions"
          :disabled="readonly"
          :show-sensitivity="project.mode === 'gated'"
          :errors="submitted ? fieldErrors : {}"
        />
      </CFormGroup>
    </div>

    <template #footer>
      <Button v-if="readonly" label="Close" size="small" @click="visible = false" />
      <template v-else>
        <Button
          label="Cancel"
          severity="secondary"
          outlined
          size="small"
          @click="visible = false"
        />
        <Button label="Save" size="small" :loading="saving" @click="onSave" />
      </template>
    </template>
  </Dialog>
</template>

<script setup>
import ValidationMessage from '@/sections/project/components/ValidationMessage.vue'
import FieldsEditor from '@/sections/project/components/datamodel/FieldsEditor.vue'
import { kindConfig } from '@/sections/project/config/kinds'
import { useProjectsStore } from '@/sections/project/stores/projects'
import { fieldName } from '@/sections/project/utils/fields'
import { useToast } from 'primevue/usetoast'
import { computed, nextTick, reactive, ref, watch } from 'vue'

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
const draft = reactive({ name: '', description: '', fields: [] })
// Once a new module is created on Save, remember its id so a re-save (before
// the dialog closes) updates it instead of creating another copy.
const createdId = ref(null)

function initDraft() {
  createdId.value = null
  submitted.value = false
  const m = module.value
  draft.name = m?.name || ''
  draft.description = m?.description || ''
  draft.fields = (m?.fields || []).map(f => ({ ...f, selectOptions: [...(f.selectOptions || [])] }))
}

// --- Validation ----------------------------------------------------------------
// Checked on Save; invalid states only show after the first attempt so a
// fresh dialog isn't covered in red.
const submitted = ref(false)

const nameError = computed(() => (draft.name.trim() ? '' : 'Name is required'))

// Per-field structural errors, keyed by field id:
// { name?, target?, options? } per field.
const fieldErrors = computed(() => {
  const out = {}
  const seen = new Map() // machine name → first label using it
  for (const f of draft.fields) {
    const errs = {}
    if (!(f.name || '').trim()) {
      errs.name = 'Field name is required'
    } else {
      const key = fieldName(f.name).toLowerCase()
      if (seen.has(key)) errs.name = `Clashes with “${seen.get(key)}”`
      else seen.set(key, f.name)
    }
    if (f.type === 'Record' && !f.targetModuleId) errs.target = 'Target module is required'
    if (f.type === 'Select' && !(f.selectOptions || []).length) {
      errs.options = 'Add at least one option'
    }
    if (Object.keys(errs).length) out[f.id] = errs
  }
  return out
})

const isValid = computed(() => !nameError.value && !Object.keys(fieldErrors.value).length)

watch(
  () => [props.modelValue, props.moduleId],
  () => {
    if (props.modelValue) initDraft()
  },
  { immediate: true },
)

const fieldsEditorRef = ref(null)
const nid = () => `f-${Math.random().toString(36).slice(2, 9)}`
const addDraftField = async () => {
  const id = nid()
  draft.fields.push({
    id,
    name: '',
    type: 'String',
    required: false,
    multi: false,
    targetModuleId: null,
    labelField: null,
    selectOptions: [],
    sensitivity: null,
  })
  // Land the cursor in the new row's name input.
  await nextTick()
  fieldsEditorRef.value?.focusName(id)
}

// --- Commit -------------------------------------------------------------------
const saving = ref(false)

async function onSave() {
  if (saving.value) return
  submitted.value = true
  if (!isValid.value) return
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
    // Modules carry no sensitivity level — only their fields are classified
    // (in the gated Data Sensitivity step).
    await store.updateResource(pid, id, { name: draft.name, description: draft.description })
    visible.value = false
  } catch (err) {
    // Stay open so the staged edits aren't lost.
    toast.add({
      severity: 'error',
      summary: 'Could not save module',
      detail: err.message,
      life: 4000,
    })
  } finally {
    saving.value = false
  }
}
</script>
