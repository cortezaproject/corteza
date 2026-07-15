<template>
  <!-- Field-level old→new diff for an action-log event (Action.delta).
       This is the view's signature: "what did it become?" is the question the
       log exists to answer, so it renders inline in the row rather than hiding
       behind a disclosure.

       Borrows the vernacular of a version-control diff — struck-through old,
       plain new — because that is the idiom this content already lives in, and
       it is the one place in the row where colour is allowed to speak. Kept
       borderless and monospaced: at three rows per event a bordered card reads
       as a second table inside the list. -->
  <dl v-if="rows.length" class="grid grid-cols-[auto_minmax(0,1fr)] gap-x-3 gap-y-1">
    <template v-for="row in rows" :key="row.key">
      <dt class="text-xs text-muted-color break-words">{{ row.key }}</dt>
      <dd class="flex flex-wrap items-baseline gap-x-1.5 gap-y-0.5 min-w-0 text-xs font-mono">
        <span class="text-red-600 line-through decoration-red-400/60 break-all dark:text-red-400">
          {{ row.old }}
        </span>
        <i class="pi pi-arrow-right text-xs text-muted-color shrink-0" />
        <span class="text-emerald-700 break-all dark:text-emerald-400">{{ row.new }}</span>
      </dd>
    </template>
  </dl>
</template>

<script setup>
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()

const props = defineProps({
  // actionlog.Action.delta — [{ key, old: [...], new: [...] }, …]
  delta: { type: Array, default: () => [] },
  // Cap the rows rendered; 0 shows all. Keeps a 15-field update from swallowing
  // the viewport while still answering the question for the common 1–3 field case.
  limit: { type: Number, default: 0 },
})

// The server sends each side as a slice (a field can hold multiple values), so
// join for display. An absent/empty side means the field was unset, which reads
// far better as "(empty)" than as a blank gap the eye skips over.
function side(v) {
  if (v === null || v === undefined) return t('project.dashboard.allEvents.diff.empty')
  const parts = (Array.isArray(v) ? v : [v]).filter(x => x !== null && x !== undefined && x !== '')
  return parts.length ? parts.join(', ') : t('project.dashboard.allEvents.diff.empty')
}

const rows = computed(() => {
  const all = (props.delta || [])
    .filter(c => c && c.key)
    .map(c => ({ key: c.key, old: side(c.old), new: side(c.new) }))
  return props.limit > 0 ? all.slice(0, props.limit) : all
})
</script>
