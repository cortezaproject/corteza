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
    <div v-if="visible" class="right-sidebar backlog-drawer flex flex-col">
      <div class="flex items-center gap-2.5 min-w-0 px-4 py-3 border-b border-surface shrink-0">
        <KindIcon :config="badge" size="lg" plain-icon />
        <div class="min-w-0 flex-1">
          <DialogEyebrow>{{ $t('project.dashboard.backlog.singular') }}</DialogEyebrow>
          <div class="font-semibold truncate leading-tight">
            {{ record?.title || $t('project.dashboard.backlog.untitled') }}
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
        <!-- Read-only summary — hand-rolled rows (the backlog schema is small
             and category/linked-event need custom rendering, unlike
             EventDetailDrawer's schema-driven grid). -->
        <div class="grid grid-cols-1 sm:grid-cols-2 gap-x-6 gap-y-4">
          <div>
            <div class="text-xs font-medium text-muted-color uppercase tracking-wide mb-1">
              {{ $t('project.dashboard.backlog.f.category') }}
            </div>
            <div class="text-sm text-color">{{ categoryLabel }}</div>
          </div>

          <div>
            <div class="text-xs font-medium text-muted-color uppercase tracking-wide mb-1">
              {{ $t('project.dashboard.backlog.f.linkedEvent') }}
            </div>
            <button
              v-if="linkedEvent"
              type="button"
              class="flex items-center gap-2 min-w-0 text-left group"
              @click="$emit('open-event')"
            >
              <KindIcon :config="CATEGORY_CONFIG[record.category]?.badge" size="sm" />
              <span class="text-sm text-color truncate group-hover:underline">
                {{ linkedEvent.title || $t('project.dashboard.event.untitled') }}
              </span>
            </button>
            <span v-else class="font-mono text-xs text-muted-color">#{{ record?.eventID }}</span>
          </div>

          <div class="sm:col-span-2">
            <div class="text-xs font-medium text-muted-color uppercase tracking-wide mb-1">
              {{ $t('project.dashboard.event.f.description') }}
            </div>
            <div class="whitespace-pre-wrap text-sm text-color">
              {{ record?.description || '—' }}
            </div>
          </div>

          <div>
            <div class="text-xs font-medium text-muted-color uppercase tracking-wide mb-1">
              {{ $t('project.dashboard.backlog.f.assignee') }}
            </div>
            <div class="text-sm text-color">{{ record?.assignee || '—' }}</div>
          </div>

          <div>
            <div class="text-xs font-medium text-muted-color uppercase tracking-wide mb-1">
              {{ $t('project.dashboard.backlog.f.priority') }}
            </div>
            <EventBadge
              v-if="record?.priority"
              :value="record.priority"
              variant="priority"
              size="md"
            />
            <div v-else class="text-sm text-color">—</div>
          </div>

          <div>
            <div class="text-xs font-medium text-muted-color uppercase tracking-wide mb-1">
              {{ $t('project.dashboard.event.f.status') }}
            </div>
            <EventBadge v-if="record?.status" :value="record.status" variant="status" size="md" />
            <div v-else class="text-sm text-color">—</div>
          </div>

          <div>
            <div class="text-xs font-medium text-muted-color uppercase tracking-wide mb-1">
              {{ $t('project.dashboard.event.f.dateDue') }}
            </div>
            <div class="text-sm text-color">{{ formatDate(record?.dateDue) }}</div>
          </div>
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
</template>

<script setup>
import DialogEyebrow from '@/sections/project/components/DialogEyebrow.vue'
import EventBadge from '@/sections/project/components/dashboard/EventBadge.vue'
import KindIcon from '@/sections/project/components/KindIcon.vue'
import { CATEGORY_CONFIG } from '@/sections/project/config/categories'
import { useEventsStore } from '@/sections/project/stores/events'
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'

const props = defineProps({
  visible: { type: Boolean, default: false },
  // The clicked row — a store-mapped backlog item: id, category, eventID,
  // title, description, assignee (resolved display name), priority, status,
  // dateDue. Null until a row is clicked.
  record: { type: Object, default: null },
})
defineEmits(['update:visible', 'edit', 'open-event'])

const { t: $t } = useI18n()
const eventsStore = useEventsStore()

// Neutral badge — same one BacklogItemDialog falls back to — used only if
// the record's category is somehow unresolved; the header otherwise borrows
// the category's own badge, same idiom as EventDetailDrawer.
const NEUTRAL_BADGE = {
  icon: 'pi pi-th-large',
  bg: 'bg-emphasis',
  ring: 'ring-surface',
  text: 'text-color',
}
const badge = computed(() => CATEGORY_CONFIG[props.record?.category]?.badge || NEUTRAL_BADGE)

const categoryLabel = computed(() => {
  const key = CATEGORY_CONFIG[props.record?.category]?.singularKey
  return key ? $t(key) : '—'
})

// Resolve the linked event via the events store (already loaded alongside
// the backlog store by DashboardLayout) — null when the event no longer
// exists (deleted) or hasn't loaded, in which case the row falls back to the
// raw `#<eventID>` reference (see template).
const linkedEvent = computed(() => {
  if (!props.record) return null
  return (
    eventsStore
      .byCategory(props.record.category)
      .find(e => e.id === String(props.record.eventID)) || null
  )
})

const formatDate = v => {
  if (!v) return '—'
  const d = new Date(v)
  return isNaN(d.getTime()) ? String(v) : d.toLocaleDateString()
}
</script>

<style scoped>
/* Widen the shared .right-sidebar shell (its default width is the
   --right-sidebar-width token) — same width as EventDetailDrawer's
   .event-drawer. Plain CSS since Tailwind's scale has no 30rem width step. */
.backlog-drawer {
  width: 30rem;
  max-width: calc(100vw - 1.5rem);
}
</style>
