<template>
  <!-- One work item on the board — draggable by native HTML5 DnD (there is no
       drag/sort library in this repo; see the other native-DnD components in
       this section, e.g. wizard/steps/PagesStep.vue and lib/vue's
       CFormItemList.vue, for the same idiom). `dragging` dims the card while
       it's the drag source; the actual move happens in BoardPanel.vue once a
       BoardColumn reports a drop. Also click/keyboard-activatable — opens the
       read-only detail drawer in BoardPanel.vue (mirrors the dashboard's
       row-click idiom); this stays enabled even while `disabled` (viewing
       never needs write capability, only the drag move does). -->
  <div
    role="button"
    tabindex="0"
    class="rounded-md border border-surface bg-surface p-2.5 shadow-sm transition hover:bg-emphasis"
    :class="[dragging ? 'opacity-40' : '', disabled ? '' : 'cursor-grab active:cursor-grabbing']"
    :draggable="!disabled"
    @dragstart="onDragStart"
    @dragend="$emit('dragend')"
    @click="onClick"
    @keydown.enter.prevent="onClick"
    @keydown.space.prevent="onClick"
  >
    <div class="flex items-center gap-1.5 flex-wrap">
      <!-- Type chip — reuses CATEGORY_CONFIG's badge colors (same bg/ring/icon/
           text families as the category screens and the M&M nav rail), so a
           card reads as the same visual identity as the rest of the section.
           Backlog items aren't one of the five categories themselves, but each
           one is linked to one (see stores/backlogItems.js); its chip borrows
           that category's colors, same as BacklogItemDialog's title badge. -->
      <span
        class="inline-flex items-center gap-1 rounded-full ring-1 px-2 py-0.5 text-[11px] font-medium"
        :class="[typeBadge.bg, typeBadge.ring, typeBadge.text]"
      >
        <i :class="[typeBadge.icon, 'text-[10px]']" />
        {{ typeLabel }}
      </span>

      <!-- Revision chip — chain-wide only (BoardPanel's `!isRevisionScoped`,
           threaded down as `showRevision`). Same treatment CategoryPanel/
           BacklogView use for their revisionID column: an explicit
           "Unassigned" state rather than a blank cell. `item.revisionLabel`
           is resolved once in BoardPanel's boardItems (via useRevisionLabel)
           rather than here, so this stays a pure presentational leaf, same as
           typeBadge/typeLabel above. Omitted revision-scoped — every card on
           that board already shares the one open revision, so the chip would
           be redundant on every card (same reasoning CategoryPanel drops its
           revision column). -->
      <template v-if="showRevision">
        <span
          v-if="item.revisionLabel.unassigned"
          class="inline-flex items-center gap-1 text-[11px] text-muted-color italic"
        >
          <i class="pi pi-question-circle text-[10px]" />
          {{ item.revisionLabel.label }}
        </span>
        <Tag
          v-else
          :value="item.revisionLabel.label"
          :severity="item.revisionLabel.severity"
          class="!text-[11px]"
        />
      </template>
    </div>

    <p class="text-sm font-medium leading-snug mt-1.5 mb-2 line-clamp-2">
      {{ item.title || $t(untitledKey) }}
    </p>

    <div v-if="item.severity || item.priority" class="flex items-center flex-wrap gap-1.5 mb-1.5">
      <EventBadge v-if="item.severity" :value="item.severity" variant="severity" />
      <EventBadge v-if="item.priority" :value="item.priority" variant="priority" />
    </div>

    <div
      v-if="item.assignee || item.dueDate"
      class="flex items-center justify-between gap-2 pt-1.5 border-t border-surface text-xs text-muted-color"
    >
      <span v-if="item.assignee" class="inline-flex items-center gap-1 min-w-0">
        <i class="pi pi-user text-[10px] shrink-0" />
        <span class="truncate">{{ item.assignee }}</span>
      </span>
      <span v-else />
      <span v-if="item.dueDate" class="inline-flex items-center gap-1 shrink-0">
        <i class="pi pi-calendar text-[10px]" />
        {{ formattedDueDate }}
      </span>
    </div>
  </div>
</template>

<script setup>
import EventBadge from '@/sections/project/components/dashboard/EventBadge.vue'
import { CATEGORY_CONFIG } from '@/sections/project/config/categories'
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()

const props = defineProps({
  // Normalized board item — see BoardPanel.vue's `boardItems` computed:
  // { key, id, itemType, linkedCategory?, title, status, severity?,
  //   priority?, assignee, dueDate, revisionLabel }.
  item: { type: Object, required: true },
  disabled: { type: Boolean, default: false },
  // True while this exact card is the active drag source.
  dragging: { type: Boolean, default: false },
  // Chain-wide only — see BoardColumn's own prop comment for the full
  // reasoning; threaded straight through from there.
  showRevision: { type: Boolean, default: false },
})
const emit = defineEmits(['dragstart', 'dragend', 'click'])

// Neutral fallback badge — mirrors BacklogItemDialog.vue's NEUTRAL_BADGE, used
// only if a backlog item's linked category ever fails to resolve.
const NEUTRAL_BADGE = {
  icon: 'pi pi-th-large',
  bg: 'bg-emphasis',
  ring: 'ring-surface',
  text: 'text-color',
}

const isBacklog = computed(() => props.item.itemType === 'backlog')
const typeBadge = computed(() => {
  const key = isBacklog.value ? props.item.linkedCategory : props.item.itemType
  return CATEGORY_CONFIG[key]?.badge || NEUTRAL_BADGE
})
const typeLabel = computed(() =>
  isBacklog.value
    ? t('project.dashboard.backlog.singular')
    : t(CATEGORY_CONFIG[props.item.itemType]?.singularKey || ''),
)
const untitledKey = computed(() =>
  isBacklog.value ? 'project.dashboard.backlog.untitled' : 'project.dashboard.event.untitled',
)

const formattedDueDate = computed(() => {
  if (!props.item.dueDate) return ''
  const d = new Date(props.item.dueDate)
  return Number.isNaN(d.getTime()) ? String(props.item.dueDate) : d.toLocaleDateString()
})

function onDragStart(e) {
  if (props.disabled) return
  e.dataTransfer.effectAllowed = 'move'
  // BoardColumn's onDrop reads this back to identify which item moved — the
  // two components aren't otherwise related, so dataTransfer is the bridge
  // (mirrors PagesStep.vue's onDragStart, which sets the same MIME type for
  // Firefox's benefit).
  e.dataTransfer.setData('text/plain', props.item.key)
  emit('dragstart', props.item.key)
}

// Same key-based payload as dragstart — BoardPanel.vue looks the full item
// back up by key (findItem) rather than this component threading the whole
// object through two emit hops.
function onClick() {
  emit('click', props.item.key)
}
</script>
