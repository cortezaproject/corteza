<template>
  <!-- Floating right sidebar — same shell as the app's other right panels
       (.right-sidebar from useTheme.ts: fixed, rounded, no overlay/dim), so
       it coexists with the page and is closed by the rightSidebarStore when
       another panel (TAQ config, notifications, agent…) opens. -->
  <Transition
    enter-active-class="transition-transform duration-300 ease-in-out"
    enter-from-class="translate-x-full"
    enter-to-class="translate-x-0"
    leave-active-class="transition-transform duration-300 ease-in-out"
    leave-from-class="translate-x-0"
    leave-to-class="translate-x-full"
  >
    <div v-if="visible && cfg" class="right-sidebar event-drawer flex flex-col">
      <div class="flex items-center gap-2.5 min-w-0 px-4 py-3 border-b border-surface shrink-0">
        <KindIcon :config="cfg.badge" size="lg" plain-icon />
        <div class="min-w-0 flex-1">
          <DialogEyebrow>{{ $t(cfg.singularKey) }}</DialogEyebrow>
          <div class="font-semibold truncate leading-tight">
            {{ record?.title || $t('project.dashboard.event.untitled') }}
          </div>
        </div>
        <Button
          icon="pi pi-times"
          severity="secondary"
          text
          rounded
          size="small"
          :aria-label="$t('general.label.close')"
          @click="$emit('update:visible', false)"
        />
      </div>

      <div class="flex-1 min-h-0 overflow-y-auto p-4 flex flex-col gap-5">
      <!-- Read-only summary — every form field except title (already the
           header), label + value pairs. Textareas (and any field flagged
           `full`) span both columns, same rule GovernanceForm itself uses.
           Badge-like fields (type/severity/risk/status) render as their pills
           right where they sit in the field order — no separate badge row. -->
      <div class="grid grid-cols-1 sm:grid-cols-2 gap-x-6 gap-y-4">
        <div
          v-for="field in summaryFields"
          :key="field.key"
          :class="{ 'sm:col-span-2': isFullField(field) }"
        >
          <div class="text-xs font-medium text-muted-color uppercase tracking-wide mb-1">
            {{ $t(field.labelKey) }}
          </div>
          <EventBadge
            v-if="badgeVariant(field) && record?.[field.key]"
            :value="record[field.key]"
            :variant="badgeVariant(field)"
            size="md"
          />
          <div
            v-else
            class="text-sm text-color"
            :class="{ 'whitespace-pre-wrap': field.type === 'textarea' }"
          >
            {{ fieldValue(field) }}
          </div>
        </div>
      </div>

      <Divider />

      <!-- Backlog items linked to this record — see stores/backlogItems.js.
           Row click inspects/edits the item, "Add item" creates one; both go
           through BacklogItemDialog pre-scoped to this event (see
           lockedEvent below). Remove stays inline with a confirm. -->
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
            @click="openAddItem"
          />
        </div>
        <CFormItemList
          :items="backlogItems"
          :remove-label="$t('general.label.remove')"
          @select="onBacklogSelect"
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

      <div class="px-4 py-3 border-t border-surface shrink-0">
        <Button
          :label="$t('general.label.edit')"
          icon="pi pi-pencil"
          size="small"
          class="w-full"
          @click="$emit('edit')"
        />
      </div>
    </div>
  </Transition>

  <!-- Add/inspect flow — the full create/edit dialog, pre-scoped to this
       record via `lockedEvent` (hides the category/linked-event selects; see
       BacklogItemDialog). "Add item" opens it blank; clicking a list row
       opens it on that item for inspect/edit/delete. -->
  <BacklogItemDialog
    v-if="cfg"
    v-model:visible="itemDialogVisible"
    :record="selectedBacklogItem"
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

// --- Badges in the summary -------------------------------------------------
// Fields that render as pills instead of plain text: severity/risk/status
// declare a `badge` hint on their schema config (shared with the dialog
// selects); the category's type field is identified via the column config
// (kind === 'type'), same source CategoryView's list uses.
const typeCol = computed(() => cfg.value?.columns.find(c => c.kind === 'type') || null)
const badgeVariant = field => field.badge || (field.key === typeCol.value?.key ? 'type' : '')

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

// Add/inspect dialog — BacklogItemDialog locked to this record's category +
// id so the payload always lands here regardless of what the (hidden)
// category/linked-event fields would otherwise default to. "Add item" opens
// it blank (`selectedBacklogItem` null → create); a list-row click opens it
// on that item (edit/delete).
const itemDialogVisible = ref(false)
const selectedBacklogItem = ref(null)
const lockedEvent = computed(() =>
  props.record ? { category: props.category, eventID: props.record.id } : null,
)

function openAddItem() {
  selectedBacklogItem.value = null
  itemDialogVisible.value = true
}

function onBacklogSelect(item) {
  selectedBacklogItem.value = item
  itemDialogVisible.value = true
}

async function onBacklogItemSave(id, payload) {
  const editing = !!id
  try {
    if (editing) await backlogStore.update(id, payload)
    else await backlogStore.add(payload)
    $toast.toastSuccess(
      t('project.dashboard.backlog.singular'),
      t(editing ? 'project.dashboard.backlog.toast.updated' : 'project.dashboard.backlog.toast.created'),
    )
    return true
  } catch (err) {
    $toast.toastErrorHandler(
      t(editing ? 'project.dashboard.backlog.toast.updateFailed' : 'project.dashboard.backlog.toast.createFailed'),
    )(err)
    return false
  }
}

async function onBacklogItemDelete(id) {
  try {
    await backlogStore.remove(id)
    $toast.toastSuccess(t('project.dashboard.backlog.singular'), t('project.dashboard.backlog.toast.deleted'))
    return true
  } catch (err) {
    $toast.toastErrorHandler(t('project.dashboard.backlog.toast.deleteFailed'))(err)
    return false
  }
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
  if (visible) {
    itemDialogVisible.value = false
    selectedBacklogItem.value = null
  }
})
</script>

<style scoped>
/* Widen the shared .right-sidebar shell (its default width is the
   --right-sidebar-width token) — the summary + backlog need more room than
   the notification-style panels. Plain CSS since Tailwind's scale has no
   30rem width step. */
.event-drawer {
  width: 30rem;
  max-width: calc(100vw - 1.5rem);
}
</style>
