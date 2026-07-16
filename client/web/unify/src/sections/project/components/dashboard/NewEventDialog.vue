<template>
  <Dialog
    :visible="visible"
    @update:visible="$emit('update:visible', $event)"
    modal
    :header="headerText"
    :style="{ width: '46rem' }"
    :pt="{ content: { class: '!pt-2' } }"
  >
    <!-- Reuse the per-category governance form; two-column layout with
         textareas spanning both columns. `model` is a plain object keyed by
         field.key and is reset every time the dialog opens. -->
    <GovernanceForm
      :schema="resolvedSchema"
      v-model="model"
      :columns="2"
      :disabled="saving"
      :errors="fieldErrors"
      :submitted="submitted"
    />

    <template #footer>
      <div class="flex items-center justify-end gap-2 w-full">
        <Button
          :label="$t('general.label.cancel')"
          severity="secondary"
          text
          size="small"
          :disabled="saving"
          @click="close"
        />
        <Button :label="createLabel" size="small" :loading="saving" @click="onSubmit" />
      </div>
    </template>
  </Dialog>
</template>

<script setup>
import GovernanceForm from '@/sections/project/components/wizard/GovernanceForm.vue'
import { computed, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'

const props = defineProps({
  visible: { type: Boolean, default: false },
  category: { type: String, required: true },
  schema: { type: Array, default: () => [] },
  // [{ label, value }] used to populate fields flagged `source: 'users'`.
  userOptions: { type: Array, default: () => [] },
  // Async create handler owned by the parent (it holds the store + toast).
  // Resolves truthy on success, falsy (or throws) on failure — either keeps
  // the dialog open with the user's input intact so they can fix and retry.
  onCreate: { type: Function, required: true },
})
const emit = defineEmits(['update:visible', 'created'])

const { t } = useI18n()

// Inject the dynamic user directory into any `source: 'users'` field so owner /
// approver selects render names and yield user IDs.
const resolvedSchema = computed(() =>
  props.schema.map(section => ({
    ...section,
    fields: section.fields.map(f =>
      f.source === 'users'
        ? { ...f, options: props.userOptions, optionLabel: 'label', optionValue: 'value' }
        : f,
    ),
  })),
)

// Collected form values, keyed by field.key. Reset on each open.
const model = ref({})

// "New {type}" where {type} is the singular category label.
const createLabel = computed(() =>
  t('project.dashboard.newButton', { type: t(`project.dashboard.categorySingular.${props.category}`) }),
)
const headerText = createLabel

// --- Validation ------------------------------------------------------------------
// Required fields come straight off the (already user-resolved) schema; empty
// string/nullish values fail regardless of field type (text/select/date all
// collapse to that check — none of the schemas use 0/false as a real value).
const submitted = ref(false)

const requiredFields = computed(() => resolvedSchema.value.flatMap(s => s.fields.filter(f => f.required)))

const isEmpty = v => v === null || v === undefined || (typeof v === 'string' && v.trim() === '')

const fieldErrors = computed(() => {
  const errs = {}
  for (const f of requiredFields.value) {
    if (isEmpty(model.value[f.key])) errs[f.key] = t('project.dashboard.event.validation.required')
  }
  return errs
})

const isValid = computed(() => Object.keys(fieldErrors.value).length === 0)

// Opening the dialog starts from a blank form so no stale values leak between
// categories or repeated creates.
watch(
  () => props.visible,
  v => {
    if (v) {
      model.value = {}
      submitted.value = false
    }
  },
)

function close() {
  emit('update:visible', false)
}

// Busy state while the parent's create call is in flight; the dialog only
// closes (and drops the draft) once that resolves successfully, so a failed
// create leaves the form exactly as the user left it for a retry.
const saving = ref(false)

async function onSubmit() {
  if (saving.value) return
  submitted.value = true
  if (!isValid.value) return
  saving.value = true
  try {
    const ok = await props.onCreate({ ...model.value })
    if (ok !== false) {
      emit('created')
      close()
    }
  } finally {
    saving.value = false
  }
}
</script>
