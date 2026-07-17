<template>
  <!-- Audit-event summary for the overview: pulse + stat trio on the left, a
       short recent-activity feed on the right. Renders nothing when the log is
       unavailable (the viewer lacks action-log.read), so non-admins simply do
       not see it. -->
  <section v-if="!failed">
    <h2 class="text-sm font-semibold uppercase tracking-wide text-muted-color mb-3">
      {{ $t('project.dashboard.eventsBand.title') }}
    </h2>

    <div class="grid grid-cols-1 lg:grid-cols-2 gap-3">
      <!-- Pulse + stats -->
      <div class="rounded-lg border border-surface bg-surface p-4 flex flex-col gap-3">
        <div class="flex gap-6">
          <div class="flex flex-col">
            <span class="text-2xl font-semibold text-color leading-none">{{ metrics.total }}</span>
            <span class="text-xs text-muted-color">
              {{ $t('project.dashboard.allEvents.metrics.events') }}
            </span>
          </div>
          <div class="flex flex-col">
            <span class="text-2xl font-semibold text-color leading-none">{{ metrics.actors }}</span>
            <span class="text-xs text-muted-color">
              {{ $t('project.dashboard.allEvents.metrics.people') }}
            </span>
          </div>
          <div class="flex flex-col">
            <span
              class="text-2xl font-semibold leading-none"
              :class="metrics.errors ? 'text-red-500' : 'text-color'"
            >
              {{ metrics.errors }}
            </span>
            <span class="text-xs text-muted-color">
              {{ $t('project.dashboard.allEvents.metrics.errors') }}
            </span>
          </div>
        </div>
        <!-- The pulse fills whatever height the grid row gives the card (the
             feed next door usually decides it); min-h keeps it readable when
             the cards stack on one column. -->
        <div class="flex-1 min-h-24">
          <CategoryTrendChart
            bare
            height="fill"
            :labels="metrics.labels"
            :range-labels="metrics.rangeLabels"
            :series="metrics.series"
          />
        </div>
      </div>

      <!-- Recent activity feed -->
      <div class="rounded-lg border border-surface bg-surface p-4">
        <div class="flex items-center justify-between gap-2 mb-2">
          <div class="text-sm font-medium text-color">
            {{ $t('project.dashboard.eventsBand.recent') }}
          </div>
          <RouterLink
            :to="{ name: 'project.overview.events', params: { projectId } }"
            class="text-xs text-primary-500 hover:underline shrink-0"
          >
            {{ $t('project.dashboard.eventsBand.viewAll') }}
          </RouterLink>
        </div>
        <div v-if="recent.length" class="flex flex-col max-h-72 overflow-y-auto">
          <template v-for="group in recentGroups" :key="group.key">
            <!-- Same day separators as All Events — the identical bg-emphasis
                 band, sticky within the card's own scroll (solid token, so
                 rows never show through underneath). -->
            <div
              class="sticky top-0 z-10 px-2 py-1.5 bg-emphasis rounded text-xs font-semibold uppercase tracking-wide text-muted-color"
            >
              {{ group.label }}
            </div>
            <EventTimelineItem
              v-for="e in group.items"
              :key="e.actionID"
              :data="e"
              :actor-name="events.actorName(e)"
              :collapsible="false"
            />
          </template>
        </div>
        <p v-else class="text-sm text-muted-color py-6 text-center">
          {{ $t('project.dashboard.eventsBand.noRecent') }}
        </p>
      </div>
    </div>
  </section>
</template>

<script setup>
import CategoryTrendChart from '@/sections/project/components/dashboard/CategoryTrendChart.vue'
import EventTimelineItem from '@/sections/project/components/dashboard/EventTimelineItem.vue'
import { useEventActivity } from '@/sections/project/composables/useEventActivity'
import { filters } from '@planetcrust/human-vue'
import { computed, reactive, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'

const props = defineProps({
  projectId: { type: [String, Number], default: '' },
})

const { t } = useI18n()
const { locDate } = filters
const events = useEventActivity()

// Group the feed by calendar day, newest first — same labels as All Events
// (Today/Yesterday/date), so the two views read as one system.
const recentGroups = computed(() => {
  const out = []
  let current = null
  for (const e of recent.value) {
    const d = e.timestamp ? new Date(e.timestamp) : null
    const key = d && !Number.isNaN(d.getTime()) ? d.toDateString() : 'unknown'
    if (!current || current.key !== key) {
      current = { key, label: dayLabel(d), items: [] }
      out.push(current)
    }
    current.items.push(e)
  }
  return out
})

function dayLabel(d) {
  if (!d || Number.isNaN(d.getTime())) return t('project.dashboard.allEvents.unknownDate')
  const today = new Date()
  const yesterday = new Date()
  yesterday.setDate(today.getDate() - 1)
  if (d.toDateString() === today.toDateString()) return t('project.dashboard.allEvents.today')
  if (d.toDateString() === yesterday.toDateString()) return t('project.dashboard.allEvents.yesterday')
  return locDate(d)
}

const metrics = reactive({ labels: [], rangeLabels: [], series: [], total: 0, actors: 0, errors: 0 })
const recent = ref([])
const failed = ref(false)

async function load(pid) {
  if (!pid) return
  try {
    // Metrics gate visibility: if the report is forbidden, hide the panel.
    const m = await events.loadMetrics(pid)
    Object.assign(metrics, m)
    failed.value = false
    recent.value = await events.loadRecent(pid).catch(() => [])
  } catch (e) {
    console.error('Failed to load event activity', e)
    failed.value = true
  }
}

watch(() => props.projectId, load, { immediate: true })
</script>
