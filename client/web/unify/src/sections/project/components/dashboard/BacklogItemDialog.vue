<template>
  <Dialog
    :visible="visible"
    @update:visible="$emit('update:visible', $event)"
    modal
    :style="{ width: '46rem' }"
    :pt="{ content: { class: '!pt-2' }, footer: { class: 'flex justify-between gap-2 p-3' } }"
  >
    <template #header>
      <div class="flex items-center gap-2.5 min-w-0">
        <KindIcon :config="badge" size="lg" plain-icon />
        <div class="min-w-0">
          <DialogEyebrow>{{ $t('project.dashboard.backlog.singular') }}</DialogEyebrow>
          <div class="font-semibold truncate leading-tight">
            {{
              record
                ? record.title || $t('project.dashboard.backlog.untitled')
                : $t('project.dashboard.backlog.newButton')
            }}
          </div>
        </div>
      </div>
    </template>

    <!-- Create + edit share one dialog (unlike events, which split
         NewEventDialog/EventDetailDialog) — the field set is small enough
         that the split isn't worth it. Category/linked-event are plain
         selects rather than GovernanceForm's `source: 'users'` mechanism,
         since their options come from local config/store lookups, not a
         prop; assignee still uses `source: 'users'` for the injected
         directory, same as the event dialogs. -->
    <GovernanceForm
      :schema="resolvedSchema"
      :model-value="model"
      :columns="2"
      :disabled="saving || deleting"
      :errors="fieldErrors"
      :submitted="submitted"
      @update:model-value="onModelUpdate"
    />

    <template #footer>
      <Button
        v-if="record"
        :label="$t('general.label.delete')"
        icon="pi pi-trash"
        severity="danger"
        text
        size="small"
        :loading="deleting"
        :disabled="saving"
        @click="onDeleteClick"
      />
      <span v-else />
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
import { CATEGORY_CONFIG, CATEGORY_ORDER } from '@/sections/project/config/categories'
import { EVENT_STATUS } from '@/sections/project/config/eventForm'
import { useEventsStore } from '@/sections/project/stores/events'
import { useConfirmDelete } from '@planetcrust/human-vue'
import { computed, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'

const props = defineProps({
  visible: { type: Boolean, default: false },
  // The clicked row (a store-mapped backlog item — see stores/
  // backlogItems.js#mapRow) — null for create.
  record: { type: Object, default: null },
  // [{ label, value }] used to populate the assignee field.
  userOptions: { type: Array, default: () => [] },
  // Async save/delete handlers owned by the parent (it holds the stores +
  // toast + create/update branching — see views/dashboard/BacklogView.vue).
  // `onSave(id, payload)` — `id` is null in create mode. Resolve truthy on
  // success, falsy (or throw) on failure — either keeps the dialog open with
  // the user's input intact so they can retry.
  onSave: { type: Function, required: true },
  onDelete: { type: Function, required: true },
})
const emit = defineEmits(['update:visible', 'saved', 'deleted'])

const { t } = useI18n()
const { confirmDelete } = useConfirmDelete()
const eventsStore = useEventsStore()

const PRIORITY_OPTIONS = ['High', 'Medium', 'Low']

// Neutral title-bar badge for create mode / an unrecognized category (mirrors
// the old BacklogView page badge); once a category is picked the icon swaps
// to that category's own badge, live.
const NEUTRAL_BADGE = { icon: 'pi pi-th-large', bg: 'bg-emphasis', ring: 'ring-surface', text: 'text-color' }
const badge = computed(() => CATEGORY_CONFIG[model.value.category]?.badge || NEUTRAL_BADGE)

const categoryOptions = computed(() =>
  CATEGORY_ORDER.map(key => ({ label: t(CATEGORY_CONFIG[key].singularKey), value: key })),
)

// Linked-event choices are the selected category's events (label = title);
// empty until a category is picked.
const eventOptions = computed(() => {
  const cat = model.value.category
  if (!cat) return []
  return eventsStore
    .byCategory(cat)
    .map(e => ({ label: e.title || t('project.dashboard.event.untitled'), value: e.id }))
})

const schema = computed(() => [
  {
    fields: [
      {
        key: 'category',
        labelKey: 'project.dashboard.backlog.f.category',
        type: 'select',
        options: categoryOptions.value,
        optionLabel: 'label',
        optionValue: 'value',
        required: true,
      },
      {
        key: 'eventID',
        labelKey: 'project.dashboard.backlog.f.linkedEvent',
        type: 'select',
        options: eventOptions.value,
        optionLabel: 'label',
        optionValue: 'value',
        filter: true,
        required: true,
      },
      {
        key: 'title',
        labelKey: 'project.dashboard.columns.title',
        type: 'text',
        full: true,
        placeholderKey: 'project.dashboard.event.placeholder.title',
        required: true,
      },
      {
        key: 'description',
        labelKey: 'project.dashboard.event.f.description',
        type: 'textarea',
        full: true,
      },
      {
        key: 'assignee',
        labelKey: 'project.dashboard.backlog.f.assignee',
        type: 'select',
        source: 'users',
        filter: true,
      },
      {
        key: 'priority',
        labelKey: 'project.dashboard.backlog.f.priority',
        type: 'select',
        options: PRIORITY_OPTIONS,
      },
      {
        key: 'status',
        labelKey: 'project.dashboard.event.f.status',
        type: 'select',
        options: EVENT_STATUS,
      },
      {
        key: 'dateDue',
        labelKey: 'project.dashboard.event.f.dateDue',
        type: 'date',
      },
    ],
  },
])

// Inject the dynamic user directory into the assignee field, same mechanism
// as NewEventDialog/EventDetailDialog.
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

// Collected form values, keyed by field.key. Reset every time the dialog
// opens — from the record in edit mode, or sane create-mode defaults
// (priority Medium / status Open, matching the events store's own default).
const model = ref({})

function buildModel() {
  const r = props.record
  if (!r) {
    return {
      category: '',
      eventID: null,
      title: '',
      description: '',
      assignee: null,
      priority: 'Medium',
      status: 'Open',
      dateDue: null,
    }
  }
  const d = r.dateDue ? new Date(r.dateDue) : null
  return {
    category: r.category || '',
    eventID: r.eventID != null ? String(r.eventID) : null,
    title: r.title || '',
    description: r.description || '',
    assignee: r.assigneeId || null,
    priority: r.priority || 'Medium',
    status: r.status || 'Open',
    dateDue: d && !Number.isNaN(d.getTime()) ? d : null,
  }
}

// GovernanceForm emits the whole merged object on every field change; when
// the category changes, the linked-event options change under it too, so
// drop a stale eventID that no longer belongs to the new category rather
// than silently submitting a mismatched pair.
function onModelUpdate(next) {
  if (next.category !== model.value.category) {
    const validIds = next.category ? eventsStore.byCategory(next.category).map(e => e.id) : []
    if (!validIds.includes(next.eventID)) next.eventID = null
  }
  model.value = next
}

// --- Validation ------------------------------------------------------------------
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
    const ok = await props.onSave(props.record?.id || null, { ...model.value })
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
    header: t('project.dashboard.backlog.confirmDelete.header'),
    message: t('project.dashboard.backlog.confirmDelete.message', {
      name: props.record?.title || t('project.dashboard.backlog.untitled'),
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
