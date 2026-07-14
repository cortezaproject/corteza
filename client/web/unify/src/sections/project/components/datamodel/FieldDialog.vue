<template>
  <Dialog
    v-model:visible="visible"
    modal
    :style="{ width: '54rem' }"
    :pt="{ content: { class: '!pt-2' }, footer: { class: 'flex justify-end gap-2' } }"
  >
    <template #header>
      <div class="flex items-center gap-2.5 min-w-0">
        <span
          class="inline-flex items-center justify-center w-8 h-8 rounded-md ring-1 ring-surface bg-emphasis shrink-0"
        >
          <i :class="[step === 'type' ? 'pi pi-tag' : typeCfg?.icon, 'text-muted-color']" />
        </span>
        <div class="min-w-0">
          <div class="text-[10px] uppercase tracking-wider text-muted-color leading-none mb-0.5">
            {{ $t('project.field.label') }}
            <template v-if="moduleName">· {{ moduleName }}</template>
          </div>
          <div class="font-semibold truncate leading-tight">
            {{
              step === 'type'
                ? $t('project.field.chooseType')
                : draft.name || $t('project.field.unnamed')
            }}
          </div>
        </div>
      </div>
    </template>

    <!-- Step 1: visual field-type picker (shown when adding, or via "Change
         type" when editing). -->
    <div v-if="step === 'type'" class="grid grid-cols-2 sm:grid-cols-3 gap-3">
      <button
        v-for="ft in fieldTypeList"
        :key="ft.id"
        type="button"
        class="group flex flex-col gap-2 p-3 rounded-border border text-left transition-colors hover:bg-emphasis"
        :class="
          ft.id === draft.type
            ? 'border-primary ring-1 ring-primary'
            : 'border-surface hover:border-primary/60'
        "
        @click="chooseType(ft.id)"
      >
        <div class="flex items-center gap-2 min-w-0">
          <span
            class="inline-flex items-center justify-center w-9 h-9 rounded-md ring-1 ring-surface bg-emphasis shrink-0"
          >
            <i :class="[ft.icon, 'text-muted-color group-hover:text-primary transition-colors']" />
          </span>
          <span class="font-medium text-sm truncate">{{ ft.label }}</span>
        </div>
        <span class="text-xs text-muted-color leading-snug">{{ ft.hint }}</span>
      </button>
    </div>

    <!-- Step 2: field configuration. -->
    <div v-else class="flex flex-col gap-5">
      <div class="grid grid-cols-1 sm:grid-cols-2 gap-5">
        <CFormGroup :label="$t('general.label.name')" required>
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
        <CFormGroup :label="$t('general.label.type')">
          <button
            type="button"
            class="flex items-center gap-2 w-full px-3 min-h-[2.25rem] rounded-md border border-surface text-left transition-colors enabled:hover:bg-emphasis disabled:cursor-default"
            :disabled="readonly"
            @click="step = 'type'"
          >
            <i :class="[typeCfg?.icon, 'text-muted-color shrink-0']" />
            <span class="text-sm flex-1 truncate">{{ typeLabel }}</span>
          </button>
        </CFormGroup>
      </div>

      <div class="flex flex-col gap-3">
        <CInputToggleCard
          v-model="draft.required"
          :label="$t('general.label.required')"
          :description="$t('project.field.requiredDescription')"
          :disabled="readonly"
        />
        <CInputToggleCard
          v-model="draft.multi"
          :label="$t('project.field.multipleValues')"
          :description="$t('project.field.multipleValuesDescription')"
          :disabled="readonly"
        />
      </div>

      <!-- Type-specific settings (record reference target, select options),
           separated from the general config above by a divider. -->
      <template v-if="isRecordRef(draft.type)">
        <div class="border-t border-surface" />
        <div class="flex flex-wrap gap-5">
          <CFormGroup
            :label="$t('project.field.targetModule')"
            :description="$t('project.field.targetModuleDescription')"
            required
          >
            <div>
              <Select
                :model-value="draft.targetModuleId"
                :options="moduleOptions"
                option-label="name"
                option-value="id"
                size="small"
                :placeholder="$t('project.field.targetModulePlaceholder')"
                class="w-72"
                :disabled="readonly"
                :invalid="submitted && !!targetError"
                @update:model-value="onTarget"
              />
              <ValidationMessage :message="submitted ? targetError : ''" />
            </div>
          </CFormGroup>
          <CFormGroup
            :label="$t('project.field.labelField')"
            :description="$t('project.field.labelFieldDescription')"
          >
            <Select
              v-model="draft.labelField"
              :options="labelFieldOptions"
              option-label="label"
              option-value="value"
              size="small"
              :placeholder="$t('project.field.labelFieldPlaceholder')"
              show-clear
              class="w-72"
              :disabled="readonly || !draft.targetModuleId"
            />
          </CFormGroup>
        </div>
      </template>

      <template v-else-if="draft.type === 'Select'">
        <div class="border-t border-surface" />
        <CFormGroup :label="$t('project.field.options')">
          <div class="flex flex-col gap-2">
            <div v-for="(opt, i) in draft.selectOptions" :key="i" class="flex items-center gap-2">
              <InputText
                v-model="opt.text"
                size="small"
                fluid
                :placeholder="$t('project.field.optionLabelPlaceholder')"
                :disabled="readonly"
                :invalid="submitted && !!optionsError && !(opt.text || '').trim()"
              />
              <Button
                v-if="!readonly"
                icon="pi pi-trash"
                severity="danger"
                text
                size="small"
                :aria-label="$t('general.label.remove')"
                @click="removeOption(i)"
              />
            </div>
            <Button
              v-if="!readonly"
              :label="$t('general.label.add')"
              icon="pi pi-plus"
              severity="secondary"
              size="small"
              outlined
              class="self-start"
              @click="addOption"
            />
            <ValidationMessage :message="submitted ? optionsError : ''" />
          </div>
        </CFormGroup>
      </template>

    </div>

    <template #footer>
      <Button
        v-if="readonly"
        :label="$t('general.label.close')"
        severity="secondary"
        text
        size="small"
        @click="visible = false"
      />
      <template v-else>
        <Button
          :label="$t('general.label.cancel')"
          severity="secondary"
          text
          size="small"
          @click="visible = false"
        />
        <!-- The picker commits on card click, so it shows no Save button. -->
        <Button
          v-if="step !== 'type'"
          :label="$t('general.label.save')"
          size="small"
          :loading="saving"
          @click="onSave"
        />
      </template>
    </template>
  </Dialog>
</template>

<script setup>
import ValidationMessage from '@/sections/project/components/ValidationMessage.vue'
import { FIELD_TYPES, fieldType, isRecordRef } from '@/sections/project/config/fieldTypes'
import { useProjectsStore } from '@/sections/project/stores/projects'
import { fieldName } from '@/sections/project/utils/fields'
import { computed, inject, reactive, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()

const props = defineProps({
  modelValue: { type: Boolean, default: false },
  project: { type: Object, required: true },
  // Owning module of the field being edited.
  moduleId: { type: String, default: null },
  // Field to edit.
  fieldId: { type: String, default: null },
  // Open as read-only (e.g. a submitted/approved step).
  readonly: { type: Boolean, default: false },
})

const emit = defineEmits(['update:modelValue'])

const store = useProjectsStore()
const $toast = inject('$toast')

const visible = computed({
  get: () => props.modelValue,
  set: v => emit('update:modelValue', v),
})

// Field types for the visual picker, with localized label + hint.
const fieldTypeList = computed(() =>
  FIELD_TYPES.map(ft => ({ id: ft.id, icon: ft.icon, label: t(ft.labelKey), hint: t(ft.hintKey) })),
)
// Display config for the currently selected type (header + form chip).
const typeCfg = computed(() => fieldType(draft.type))
const typeLabel = computed(() => (typeCfg.value ? t(typeCfg.value.labelKey) : draft.type))

const module = computed(() =>
  props.moduleId ? store.resourcesFor(props.project?.id).find(r => r.id === props.moduleId) : null,
)
const moduleName = computed(() => module.value?.name || '')
const field = computed(() =>
  props.fieldId ? (module.value?.fields || []).find(f => f.id === props.fieldId) : null,
)

// Modules a Record field can point at (anything but its own module).
const moduleOptions = computed(() =>
  store.resourcesFor(props.project?.id).filter(r => r.kind === 'module' && r.id !== props.moduleId),
)

// --- Draft (staged) edits; nothing persists until Save -----------------------
const draft = reactive({
  name: '',
  type: 'String',
  required: false,
  multi: false,
  sensitivity: null,
  helpText: '',
  multiLine: false,
  precision: 0,
  dateMode: 'datetime',
  targetModuleId: null,
  labelField: null,
  selectOptions: [],
})

// Two-step flow: pick a type, then configure. New fields start on the picker;
// editing an existing field jumps straight to the form.
const step = ref('form')

function initDraft() {
  submitted.value = false
  const f = field.value
  draft.name = f?.name || ''
  draft.type = f?.type || 'String'
  draft.required = !!f?.required
  draft.multi = !!f?.multi
  draft.sensitivity = f?.sensitivity || null
  draft.helpText = f?.helpText || ''
  draft.multiLine = !!f?.multiLine
  draft.precision = f?.precision ?? 0
  draft.dateMode = f?.dateMode || 'datetime'
  draft.targetModuleId = f?.targetModuleId || null
  draft.labelField = f?.labelField || null
  // Clone each option so edits stay staged until Save.
  draft.selectOptions = (f?.selectOptions || []).map(o => ({ value: o.value, text: o.text }))
  step.value = props.fieldId ? 'form' : 'type'
}

// Pick a type from the picker, then move on to configuration. Switching to a
// different kind clears any settings that no longer apply.
function chooseType(id) {
  const changed = draft.type !== id
  draft.type = id
  if (changed) onType()
  step.value = 'form'
}

watch(
  () => [props.modelValue, props.moduleId, props.fieldId],
  () => {
    if (props.modelValue) initDraft()
  },
  { immediate: true },
)

// --- Validation (checked on Save) --------------------------------------------
const submitted = ref(false)

const nameError = computed(() => {
  const name = (draft.name || '').trim()
  if (!name) return t('project.field.nameRequired')
  const key = fieldName(name).toLowerCase()
  const clash = (module.value?.fields || []).some(f => {
    // Skip the field being edited and any blank-named siblings (an empty name
    // normalizes to the `field` handle, which would clash spuriously).
    if (f.id === props.fieldId) return false
    const other = (f.name || '').trim()
    if (!other) return false
    return fieldName(other).toLowerCase() === key
  })
  return clash ? t('project.field.nameClash') : ''
})

const targetError = computed(() =>
  isRecordRef(draft.type) && !draft.targetModuleId ? t('project.field.targetRequired') : '',
)

const optionsError = computed(() => {
  if (draft.type !== 'Select') return ''
  const labels = draft.selectOptions.map(o => (o.text || '').trim()).filter(Boolean)
  if (!labels.length) return t('project.field.optionsRequired')
  // Values are machine handles derived from the labels — two labels that
  // collapse to the same handle would silently merge, so flag it.
  const values = labels.map(l => fieldName(l))
  if (values.some((v, i) => values.indexOf(v) !== i)) return t('project.field.optionsClash')
  return ''
})

const isValid = computed(() => !nameError.value && !targetError.value && !optionsError.value)

// Fields of the referenced module, by their (predicted) machine name.
const labelFieldOptions = computed(() => {
  const target = moduleOptions.value.find(m => m.id === draft.targetModuleId)
  return (target?.fields || []).map(tf => ({ label: tf.name, value: fieldName(tf.name) }))
})

// Clear type-specific settings when the kind changes away from them.
function onType() {
  if (!isRecordRef(draft.type)) {
    draft.targetModuleId = null
    draft.labelField = null
  }
  if (draft.type !== 'Select') draft.selectOptions = []
  if (draft.type !== 'String') draft.multiLine = false
  if (draft.type === 'Number') draft.precision = draft.precision ?? 0
  else draft.precision = null
  if (draft.type !== 'DateTime') draft.dateMode = 'datetime'
}

function onTarget(moduleId) {
  draft.targetModuleId = moduleId
  draft.labelField = null
}

function addOption() {
  draft.selectOptions.push({ value: '', text: '' })
}

function removeOption(i) {
  draft.selectOptions.splice(i, 1)
}

// --- Commit ------------------------------------------------------------------
const saving = ref(false)

async function onSave() {
  if (saving.value) return
  submitted.value = true
  if (!isValid.value) return
  saving.value = true
  try {
    const patch = {
      name: draft.name.trim(),
      type: draft.type,
      required: draft.required,
      multi: draft.multi,
      sensitivity: draft.sensitivity || null,
      helpText: draft.helpText.trim(),
      multiLine: draft.type === 'String' ? draft.multiLine : false,
      precision: draft.type === 'Number' ? (draft.precision ?? 0) : null,
      dateMode: draft.type === 'DateTime' ? draft.dateMode : 'datetime',
      targetModuleId: draft.type === 'Record' ? draft.targetModuleId || null : null,
      labelField: draft.type === 'Record' ? draft.labelField || null : null,
      // The stored value is a machine handle derived from the label, the same
      // way field names are generated.
      selectOptions:
        draft.type === 'Select'
          ? draft.selectOptions
              .filter(o => (o.text || '').trim())
              .map(o => ({ value: fieldName(o.text.trim()), text: o.text.trim() }))
          : [],
    }
    // No fieldId means we're adding a new field to the module.
    if (props.fieldId) {
      await store.updateField(props.project.projectID, props.moduleId, props.fieldId, patch)
    } else {
      await store.addField(props.project.projectID, props.moduleId, patch)
    }
    visible.value = false
  } catch (err) {
    // Stay open so the staged edits aren't lost.
    $toast.toastErrorHandler(t('project.field.toast.saveFailed'))(err)
  } finally {
    saving.value = false
  }
}
</script>
