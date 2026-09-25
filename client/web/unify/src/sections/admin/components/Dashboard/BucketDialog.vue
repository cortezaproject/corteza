<template>
  <Dialog
    :visible="visible"
    modal
    :header="title"
    class="w-full max-w-3xl"
    data-testid="dashboard-bucket-dialog"
    @update:visible="$emit('update:visible', $event)"
  >
    <div v-if="error" class="text-base text-red-500">{{ $t('dashboard.loadFailed') }}</div>

    <div v-else-if="!events" class="flex items-center justify-center py-10">
      <ProgressSpinner style="width: 32px; height: 32px" />
    </div>

    <div v-else class="flex flex-col gap-5">
      <div class="flex flex-wrap items-baseline gap-x-6 gap-y-1">
        <div>
          <span class="text-3xl font-semibold leading-none text-color">
            {{ events.total.toLocaleString() }}
          </span>
          <span class="text-base text-muted-color ml-2">{{ countLabel }}</span>
        </div>
        <span v-if="events.total > items.length" class="text-sm text-muted-color">
          {{ $t('dashboard.bucket.showingNewest', { n: items.length }) }}
        </span>
      </div>

      <section v-if="events.ranking.length" class="flex flex-col gap-1.5">
        <h3 class="text-base font-semibold text-color">{{ $t('dashboard.bucket.byResource') }}</h3>
        <ul class="flex flex-col gap-1 text-base">
          <li v-for="r in events.ranking" :key="r.key" class="flex items-center gap-3">
            <span class="w-44 shrink-0 break-words text-muted-color text-sm">
              {{ resourceTypeLabel(r.key) }}
            </span>
            <span class="flex-1 h-1.5 rounded-full bg-emphasis overflow-hidden">
              <span
                class="block h-full rounded-full"
                :style="{ width: barWidth(r.count) + '%', backgroundColor: barColor }"
              />
            </span>
            <span class="w-14 text-right tabular-nums text-sm">{{ r.count.toLocaleString() }}</span>
          </li>
        </ul>
      </section>

      <section class="flex flex-col gap-1.5">
        <h3 class="text-base font-semibold text-color">{{ listLabel }}</h3>
        <AttentionList
          :items="items"
          :empty-text="$t('dashboard.empty')"
          @select="$emit('select', $event)"
        />
      </section>
    </div>

    <template #footer>
      <Button :label="$t('general.label.close')" @click="$emit('update:visible', false)" />
    </template>
  </Dialog>
</template>

<script setup>
import { computed, inject, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { filters } from '@planetcrust/human-vue'
import AttentionList from './AttentionList.vue'
import { SERIES, seriesColor, useIsDark } from './chartTheme'
import { describeEvent, resourceTypeLabel } from './eventLabel'
import { useActors } from './useActors'

const { locFullDateTime } = filters

const props = defineProps({
  visible: { type: Boolean, default: false },
  // activity | signins | workflows | taqs
  kind: { type: String, default: '' },
  // ISO strings bounding the clicked bucket
  from: { type: String, default: '' },
  to: { type: String, default: '' },
  // the bucket's tooltip label, as the dialog title
  label: { type: String, default: '' },
})
const emit = defineEmits(['update:visible', 'select'])

const { t } = useI18n()
const $SystemAPI = inject('$SystemAPI')
const actors = useActors($SystemAPI)
const isDark = useIsDark()

const events = ref(null)
const error = ref(null)

let pending = null
async function load() {
  events.value = null
  error.value = null
  if (!props.kind) return

  const req = $SystemAPI.statsEventsCancellable({
    kind: props.kind,
    from: props.from,
    to: props.to,
  })
  pending?.cancel()
  pending = req

  try {
    const result = await req.response()
    if (pending !== req) return
    events.value = result
    actors.resolve([
      ...result.entries.map(e => e.actorID),
      ...result.sessions.map(s => s.createdBy),
    ])
  } catch (e) {
    if (pending === req && !e?.__CANCEL__) error.value = e
  }
}

watch(
  () => [props.visible, props.kind, props.from, props.to],
  ([visible]) => {
    if (visible) load()
  },
  { immediate: true },
)

const title = computed(() => `${t(`dashboard.bucket.${props.kind}`)} · ${props.label}`)
const countLabel = computed(() => t(`dashboard.bucket.count.${props.kind}`))
const listLabel = computed(() => t(`dashboard.bucket.list.${props.kind}`))
const barColor = computed(() => seriesColor('activity', isDark.value))

const maxCount = computed(() => Math.max(1, ...(events.value?.ranking || []).map(r => r.count)))
const barWidth = count => Math.max(2, Math.round((count / maxCount.value) * 100))

// Rows in the shape AttentionList and EventDialog share.
const items = computed(() => {
  if (!events.value) return []

  if (props.kind === 'workflows') {
    return events.value.sessions.map(s => ({
      id: s.sessionID,
      kind: 'session',
      at: s.createdAt,
      title:
        s.workflowName || s.workflowHandle || t('dashboard.runs.workflow', { id: s.workflowID }),
      detail:
        s.error ||
        [t(`dashboard.runs.${s.status}`, s.status), s.eventType].filter(Boolean).join(' · '),
      time: locFullDateTime(s.createdAt),
      to: { name: 'automation.sessions.view', params: { sessionID: s.sessionID } },
      color:
        s.status === 'failed'
          ? SERIES.failed
          : s.status === 'completed'
            ? SERIES.completed
            : SERIES.running,
      raw: s,
    }))
  }

  return events.value.entries.map(e => ({
    id: e.actionID,
    kind: 'log',
    at: e.timestamp,
    title:
      props.kind === 'signins'
        ? actors.label(e.actorID) || e.meta?.email || describeEvent(e.resource, e.action)
        : describeEvent(e.resource, e.action),
    detail: e.error || e.description || (props.kind === 'signins' ? e.meta?.email : ''),
    time: locFullDateTime(e.timestamp),
    to:
      props.kind === 'taqs' && e.meta?.automationID
        ? { name: 'automation.taq.edit', params: { automationID: e.meta.automationID } }
        : { name: 'system.actionLog' },
    color: e.error
      ? SERIES.error
      : props.kind === 'taqs'
        ? SERIES.completed
        : seriesColor('activity', isDark.value),
    raw: e,
  }))
})
</script>
