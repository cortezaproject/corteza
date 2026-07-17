<template>
  <Drawer
    :visible="visible"
    @update:visible="$emit('update:visible', $event)"
    position="right"
    class="event-drawer"
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

    <div v-if="cfg" class="flex flex-col gap-5">
      <!-- Badge row — type/severity/status, derived from the category's own
           column config so it stays in sync with CategoryView's list (review
           has no severity, so that badge is simply absent). -->
      <div class="flex items-center gap-2 flex-wrap">
        <EventBadge v-if="typeCol" :value="record?.[typeCol.key]" variant="type" />
        <EventBadge v-if="severityCol" :value="record?.[severityCol.key]" variant="severity" />
        <EventBadge v-if="statusCol" :value="record?.[statusCol.key]" variant="status" />
      </div>

      <!-- Read-only summary — every form field except title (already the
           header), label + value pairs. Textareas (and any field flagged
           `full`) span both columns, same rule GovernanceForm itself uses. -->
      <div class="grid grid-cols-1 sm:grid-cols-2 gap-x-6 gap-y-4">
        <div
          v-for="field in summaryFields"
          :key="field.key"
          :class="{ 'sm:col-span-2': isFullField(field) }"
        >
          <div class="text-xs font-medium text-muted-color uppercase tracking-wide mb-1">
            {{ $t(field.labelKey) }}
          </div>
          <div
            class="text-sm text-color"
            :class="{ 'whitespace-pre-wrap': field.type === 'textarea' }"
          >
            {{ fieldValue(field) }}
          </div>
        </div>
      </div>

      <Divider />

      <!-- Backlog items linked to this record — see stores/backlogItems.js.
           Remove-only here + an "Add item" button that opens BacklogItemDialog
           pre-scoped to this event (see addItemVisible/lockedEvent below);
           full edit (priority/assignee/status/due) happens from the Backlog
           page itself. Moved here from EventDetailDialog, which is now the
           pure edit form. -->
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
            @click="addItemVisible = true"
          />
        </div>
        <CFormItemList
          :items="backlogItems"
          :remove-label="$t('general.label.remove')"
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
    </div>

    <template #footer>
      <Button
        :label="$t('general.label.edit')"
        icon="pi pi-pencil"
        size="small"
        class="w-full"
        @click="$emit('edit')"
      />
    </template>
  </Drawer>

  <!-- "Add item" flow — the full create/edit dialog, pre-scoped to this
       record via `lockedEvent` (hides the category/linked-event selects; see
       BacklogItemDialog). A Dialog opened from within a Drawer: both are
       Portal-teleported and manage their own stacking z-index, so this just
       works as a sibling. -->
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
  // [{ label, value }] used to populate BacklogItemDialog's assignee field.
  userOptions: { type: Array, default: () => [] },
})
defineEmits(['update:visible', 'edit'])

const { t } = useI18n()
const $toast = inject('$toast')
const backlogStore = useBacklogItemsStore()
const { confirmDelete } = useConfirmDelete()

const cfg = computed(() => CATEGORY_CONFIG[props.category] || null)

// --- Badge row -----------------------------------------------------------
// Reuse the category's own column config rather than re-deriving field keys —
// stays in sync with CategoryView's list by construction.
const typeCol = computed(() => cfg.value?.columns.find(c => c.kind === 'type') || null)
const severityCol = computed(() => cfg.value?.columns.find(c => c.kind === 'severity') || null)
const statusCol = computed(() => cfg.value?.columns.find(c => c.kind === 'status') || null)

// --- Read-only summary -----------------------------------------------------
// Every field in the category's form schema except title (rendered in the
// header instead).
const summaryFields = computed(() =>
  (cfg.value?.formSchema || []).flatMap(section => section.fields).filter(f => f.key !== 'title'),
)

const isFullField = field => !!(field.full || field.type === 'textarea')

const formatDate = v => {
  if (!v) return '—'
  const d = new Date(v)
  return isNaN(d.getTime()) ? String(v) : d.toLocaleDateString()
}

// User-select fields already carry the resolved display name under
// `record[field.key]` (see stores/events.js#mapRow) — no special-casing
// needed beyond dates.
function fieldValue(field) {
  const v = props.record ? props.record[field.key] : ''
  if (field.type === 'date') return formatDate(v)
  return v || '—'
}

// --- Backlog items -----------------------------------------------------------
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

// Reset the "add item" dialog whenever the drawer (re)opens or the record it's
// showing changes (e.g. a new row click while the drawer stays open) — never
// leave it locked to a stale record.
watch([() => props.visible, () => props.record], ([visible]) => {
  if (visible) addItemVisible.value = false
})
</script>

<style scoped>
/* Tailwind's scale has no width utility around 30rem (and arbitrary bracket
   values are off the table here), so the drawer's width lives in plain CSS —
   same idiom as CategoryView's .category-list. */
.event-drawer {
  width: 30rem;
  max-width: 90vw;
}
</style>
