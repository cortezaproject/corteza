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

    <!-- Inline "Step 3 — Add Backlog Items" widget from the reference demo:
         queue follow-up work titles locally (no backend id yet — items are
         only created once the event itself is created, see onSubmit); items
         start unassigned. Purely a title list — no priority/assignee picker
         here, matching the widget's own scope in the demo. -->
    <div class="flex flex-col gap-2 mt-6">
      <span class="text-sm font-medium text-color">
        {{ $t('project.dashboard.event.backlogWidget.title') }}
      </span>
      <CFormItemList
        :items="backlogDraft"
        :remove-label="$t('general.label.remove')"
        :hide-remove="saving"
        @remove="(item, index) => backlogDraft.splice(index, 1)"
      >
        <template #default="{ item }">
          <CFormItemContent :title="item.title" />
        </template>
      </CFormItemList>
      <div class="flex items-center gap-2">
        <InputText
          v-model="newBacklogTitle"
          :placeholder="$t('project.dashboard.event.backlogWidget.placeholder')"
          :disabled="saving"
          size="small"
          fluid
          @keyup.enter="addBacklogDraft"
        />
        <Button
          icon="pi pi-plus"
          :label="$t('general.label.add')"
          severity="secondary"
          outlined
          size="small"
          :disabled="saving || !newBacklogTitle.trim()"
          @click="addBacklogDraft"
        />
      </div>
    </div>

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
import { useRevisionOptions } from '@/sections/project/composables/useRevisionOptions'
import { computed, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'

const props = defineProps({
  visible: { type: Boolean, default: false },
  category: { type: String, required: true },
  schema: { type: Array, default: () => [] },
  // [{ label, value }] used to populate fields flagged `source: 'users'`.
  userOptions: { type: Array, default: () => [] },
  // Offers the schema's revisionID field (config/eventForm.js REVISION_FIELD)
  // at create time, letting the item be pre-assigned to a revision from the
  // dashboard's chain-wide "+ New" flow. Left false (the default) by
  // board-scoped callers — the wizard board's quick-add (ManageBoard.vue,
  // unmodified) and CategoryPanel.vue's revision-scoped mode — whose create
  // flow already assigns the new item to the open revision implicitly (see
  // stores/events.js#add's trailing `revisionId` arg); showing a field there
  // that gets silently overridden would only confuse. Editing an item always
  // shows this field regardless (see EventDetailDialog) — this prop only
  // gates CREATE.
  allowRevisionSelect: { type: Boolean, default: false },
  // Async create handler owned by the parent (it holds the store + toast).
  // Called as `onCreate(payload, backlogTitles)` — `backlogTitles` is the
  // queued list of titles from the widget below (may be empty). Resolves
  // truthy on success, falsy (or throws) on failure — either keeps the dialog
  // open with the user's input intact so they can fix and retry.
  onCreate: { type: Function, required: true },
})
const emit = defineEmits(['update:visible', 'created'])

const { t } = useI18n()
const { revisionOptions, ensureLoaded: ensureRevisionsLoaded } = useRevisionOptions()

// Inject the dynamic user directory into any `source: 'users'` field so owner /
// approver selects render names and yield user IDs; inject the revision chain
// into the `source: 'revisions'` field the same way. The revision field itself
// is dropped entirely (not just disabled) when `allowRevisionSelect` is false,
// so its `default: null` never seeds `model.revisionID` (see the prop's own
// comment).
const resolvedSchema = computed(() =>
  props.schema.map(section => ({
    ...section,
    fields: section.fields
      .filter(f => f.source !== 'revisions' || props.allowRevisionSelect)
      .map(f => {
        if (f.source === 'users') {
          return { ...f, options: props.userOptions, optionLabel: 'label', optionValue: 'value' }
        }
        if (f.source === 'revisions') {
          return {
            ...f,
            options: revisionOptions.value,
            optionLabel: 'label',
            optionValue: 'value',
          }
        }
        return f
      }),
  })),
)

// Collected form values, keyed by field.key. Reset on each open, seeded with
// whatever defaults the (already user-resolved) schema declares — e.g.
// severity/risk/status — so required fields never open on an empty
// selection. Fields without a `default` (type selects, user/owner selects,
// review's frequency/scope, …) are left absent, same as before.
const model = ref({})

function buildDefaults() {
  const m = {}
  for (const section of resolvedSchema.value) {
    for (const f of section.fields) {
      if (f.default !== undefined) m[f.key] = f.default
    }
  }
  return m
}

// Queued backlog-item titles (Step 3's inline widget) — [{ title }]. Reset on
// each open alongside `model`; flattened to plain strings and handed to
// `onCreate` as its second argument on submit.
const backlogDraft = ref([])
const newBacklogTitle = ref('')

function addBacklogDraft() {
  const title = newBacklogTitle.value.trim()
  if (!title) return
  backlogDraft.value.push({ title })
  newBacklogTitle.value = ''
}

// "New {type}" where {type} is the singular category label.
const createLabel = computed(() =>
  t('project.dashboard.newButton', {
    type: t(`project.dashboard.categorySingular.${props.category}`),
  }),
)
const headerText = createLabel

// --- Validation ------------------------------------------------------------------
// Required fields come straight off the (already user-resolved) schema; empty
// string/nullish values fail regardless of field type (text/select/date all
// collapse to that check — none of the schemas use 0/false as a real value).
const submitted = ref(false)

const requiredFields = computed(() =>
  resolvedSchema.value.flatMap(s => s.fields.filter(f => f.required)),
)

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
      model.value = buildDefaults()
      backlogDraft.value = []
      newBacklogTitle.value = ''
      submitted.value = false
      if (props.allowRevisionSelect) ensureRevisionsLoaded()
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
    const ok = await props.onCreate(
      { ...model.value },
      backlogDraft.value.map(b => b.title),
    )
    if (ok !== false) {
      emit('created')
      close()
    }
  } finally {
    saving.value = false
  }
}
</script>
