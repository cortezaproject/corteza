<template>
  <Dialog
    :visible="visible"
    @update:visible="$emit('update:visible', $event)"
    modal
    :style="{ width: '46rem' }"
    :pt="{ content: { class: '!pt-2' }, footer: { class: 'flex justify-between gap-2 p-3' } }"
  >
    <template #header>
      <div v-if="cfg" class="flex items-center gap-2.5 min-w-0">
        <KindIcon :config="cfg.badge" size="lg" plain-icon />
        <div class="min-w-0">
          <DialogEyebrow>{{ $t(cfg.singularKey) }}</DialogEyebrow>
          <div class="font-semibold truncate leading-tight">
            {{ record?.title || $t('project.dashboard.event.untitled') }}
          </div>
        </div>
      </div>
    </template>

    <!-- Reuse the per-category governance form; two-column layout with
         textareas spanning both columns. `model` is seeded from the record's
         raw field values (owner selects from their `<key>Id` companion) every
         time the dialog opens. -->
    <GovernanceForm
      v-if="cfg"
      :schema="resolvedSchema"
      v-model="model"
      :columns="2"
      :disabled="saving || deleting"
      :errors="fieldErrors"
      :submitted="submitted"
    />

    <template #footer>
      <Button
        :label="$t('general.label.delete')"
        icon="pi pi-trash"
        severity="danger"
        text
        size="small"
        :loading="deleting"
        :disabled="saving"
        @click="onDeleteClick"
      />
      <div class="flex items-center gap-2">
        <Button
          :label="$t('general.label.cancel')"
          severity="secondary"
          text
          size="small"
          :disabled="saving || deleting"
          @click="close"
        />
        <Button
          :label="$t('general.label.save')"
          size="small"
          :loading="saving"
          :disabled="deleting"
          @click="onSubmit"
        />
      </div>
    </template>
  </Dialog>
</template>

<script setup>
import DialogEyebrow from '@/sections/project/components/DialogEyebrow.vue'
import KindIcon from '@/sections/project/components/KindIcon.vue'
import GovernanceForm from '@/sections/project/components/wizard/GovernanceForm.vue'
import { CATEGORY_CONFIG } from '@/sections/project/config/categories'
import { useConfirmDelete } from '@planetcrust/human-vue'
import { computed, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'

const props = defineProps({
  visible: { type: Boolean, default: false },
  category: { type: String, required: true },
  // The clicked row (a store-mapped event: resolved owner names plus raw
  // `<key>Id` companions, backlog as an array, and every other field as the
  // raw record value) — null until a row is clicked.
  record: { type: Object, default: null },
  // [{ label, value }] used to populate fields flagged `source: 'users'`.
  userOptions: { type: Array, default: () => [] },
  // Async save/delete handlers owned by the parent (it holds the store +
  // toast). Resolve truthy on success, falsy (or throw) on failure — either
  // keeps the dialog open with the user's input intact so they can retry.
  onSave: { type: Function, required: true },
  onDelete: { type: Function, required: true },
})
const emit = defineEmits(['update:visible', 'saved', 'deleted'])

const { t } = useI18n()
const { confirmDelete } = useConfirmDelete()

const cfg = computed(() => CATEGORY_CONFIG[props.category] || null)
const schema = computed(() => cfg.value?.formSchema || [])

// Inject the dynamic user directory into any `source: 'users'` field so owner /
// approver selects render names and yield user IDs — same as NewEventDialog.
const resolvedSchema = computed(() =>
  schema.value.map(section => ({
    ...section,
    fields: section.fields.map(f =>
      f.source === 'users'
        ? { ...f, options: props.userOptions, optionLabel: 'label', optionValue: 'value' }
        : f,
    ),
  })),
)

// Collected form values, keyed by field.key. Reset from the record every time
// the dialog opens.
const model = ref({})

// User-select fields carry the record's raw id (mapRow keeps it under
// `<key>Id`) rather than the resolved display name; date fields need a Date
// object for the picker. Everything else is the raw record value as-is.
function buildModel() {
  const r = props.record
  const m = {}
  for (const section of resolvedSchema.value) {
    for (const f of section.fields) {
      if (f.source === 'users') {
        m[f.key] = r?.[`${f.key}Id`] || null
      } else if (f.type === 'date') {
        const d = r?.[f.key] ? new Date(r[f.key]) : null
        m[f.key] = d && !Number.isNaN(d.getTime()) ? d : null
      } else {
        m[f.key] = r?.[f.key] ?? ''
      }
    }
  }
  // Backlog isn't editable via this form, but the update call still expects it
  // — carry the record's current value through untouched so saving never
  // silently clears it.
  m.backlog = r?.backlog ?? []
  return m
}

// --- Validation ------------------------------------------------------------------
// Same idiom as NewEventDialog: required fields off the (already
// user-resolved) schema, blank on submit only.
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

watch(
  () => props.visible,
  v => {
    if (v) {
      model.value = buildModel()
      submitted.value = false
    }
  },
)

function close() {
  emit('update:visible', false)
}

// --- Save --------------------------------------------------------------------
const saving = ref(false)

async function onSubmit() {
  if (saving.value || deleting.value) return
  submitted.value = true
  if (!isValid.value) return
  saving.value = true
  try {
    const ok = await props.onSave(props.record.id, { ...model.value })
    if (ok !== false) {
      emit('saved')
      close()
    }
  } finally {
    saving.value = false
  }
}

// --- Delete (confirm first) ----------------------------------------------------
const deleting = ref(false)

function onDeleteClick() {
  confirmDelete({
    header: t('project.dashboard.event.confirmDelete.header'),
    message: t('project.dashboard.event.confirmDelete.message', {
      name: props.record?.title || t('project.dashboard.event.untitled'),
    }),
    onConfirm: handleDelete,
  })
}

async function handleDelete() {
  if (deleting.value || saving.value) return
  deleting.value = true
  try {
    const ok = await props.onDelete(props.record.id)
    if (ok !== false) {
      emit('deleted')
      close()
    }
  } finally {
    deleting.value = false
  }
}
</script>
