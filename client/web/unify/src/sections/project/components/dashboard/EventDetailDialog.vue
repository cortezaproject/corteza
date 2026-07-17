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

    <!-- Backlog items linked to this record — see stores/backlogItems.js.
         Remove-only here + an "Add item" button that opens BacklogItemDialog
         pre-scoped to this event (see addItemVisible/lockedEvent below); full
         edit (priority/assignee/status/due) happens from the Backlog page
         itself. -->
    <template v-if="cfg">
      <Divider />
      <div class="flex flex-col gap-2">
        <div class="flex items-center justify-between gap-2">
          <span class="text-sm font-medium text-color">
            {{ $t('project.dashboard.backlog.section.title', { count: backlogItems.length }) }}
          </span>
          <Button
            icon="pi pi-plus"
            :label="$t('project.dashboard.backlog.section.addItem')"
            severity="secondary"
            outlined
            size="small"
            :disabled="saving || deleting"
            @click="addItemVisible = true"
          />
        </div>
        <CFormItemList
          :items="backlogItems"
          :remove-label="$t('general.label.remove')"
          :hide-remove="saving || deleting"
          @remove="onBacklogRemove"
        >
          <template #default="{ item }">
            <CFormItemContent :title="item.title" :subtitle="item.assignee" />
          </template>
          <template #actions="{ item }">
            <EventBadge :value="item.priority" variant="priority" />
          </template>
        </CFormItemList>
      </div>
    </template>

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

  <!-- "Add item" flow — the full create/edit dialog, pre-scoped to this
       record via `lockedEvent` (hides the category/linked-event selects; see
       BacklogItemDialog). A Dialog opened from within a Dialog: PrimeVue's
       Dialog teleports to the body and manages its own stacking z-index (same
       as confirmDelete already stacking over this dialog), so this just works
       as a sibling — closing it doesn't touch `visible` above. -->
  <BacklogItemDialog
    v-if="cfg"
    v-model:visible="addItemVisible"
    :record="null"
    :locked-event="lockedEvent"
    :user-options="userOptions"
    :on-save="onBacklogItemSave"
    :on-delete="onBacklogItemDelete"
  />
</template>

<script setup>
import BacklogItemDialog from '@/sections/project/components/dashboard/BacklogItemDialog.vue'
import DialogEyebrow from '@/sections/project/components/DialogEyebrow.vue'
import EventBadge from '@/sections/project/components/dashboard/EventBadge.vue'
import KindIcon from '@/sections/project/components/KindIcon.vue'
import GovernanceForm from '@/sections/project/components/wizard/GovernanceForm.vue'
import { CATEGORY_CONFIG } from '@/sections/project/config/categories'
import { useBacklogItemsStore } from '@/sections/project/stores/backlogItems'
import { useConfirmDelete } from '@planetcrust/human-vue'
import { computed, inject, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'

const props = defineProps({
  visible: { type: Boolean, default: false },
  category: { type: String, required: true },
  // The clicked row (a store-mapped event: resolved owner names plus raw
  // `<key>Id` companions, and every other field as the raw record value) —
  // null until a row is clicked.
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

const $toast = inject('$toast')
const backlogStore = useBacklogItemsStore()

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

// --- Backlog items -----------------------------------------------------------
// Items linked to this record — see stores/backlogItems.js#byEvent.
const backlogItems = computed(() =>
  props.record ? backlogStore.byEvent(props.category, props.record.id) : [],
)

// "Add item" dialog — BacklogItemDialog in create mode, locked to this
// record's category + id so the payload always lands here regardless of what
// the (hidden) category/linked-event fields would otherwise default to.
const addItemVisible = ref(false)
const lockedEvent = computed(() =>
  props.record ? { category: props.category, eventID: props.record.id } : null,
)

async function onBacklogItemSave(id, payload) {
  try {
    await backlogStore.add(payload)
    $toast.toastSuccess(t('project.dashboard.backlog.singular'), t('project.dashboard.backlog.toast.created'))
    return true
  } catch (err) {
    $toast.toastErrorHandler(t('project.dashboard.backlog.toast.createFailed'))(err)
    return false
  }
}

// Unreachable in practice — the dialog is always opened with `record: null`
// (create-only) here, so it never renders a Delete button, but `onDelete` is
// a required prop.
async function onBacklogItemDelete() {
  return false
}

function onBacklogRemove(item) {
  confirmDelete({
    header: t('project.dashboard.backlog.confirmDelete.header'),
    message: t('project.dashboard.backlog.confirmDelete.message', {
      name: item.title || t('project.dashboard.backlog.untitled'),
    }),
    onConfirm: async () => {
      try {
        await backlogStore.remove(item.id)
      } catch (err) {
        $toast.toastErrorHandler(t('project.dashboard.backlog.toast.deleteFailed'))(err)
      }
    },
  })
}

watch(
  () => props.visible,
  v => {
    if (v) {
      model.value = buildModel()
      addItemVisible.value = false
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
