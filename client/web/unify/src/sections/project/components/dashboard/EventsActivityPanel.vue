<template>
  <!-- Audit-event summary for the overview: pulse + stat trio on the left, a
       short recent-activity feed on the right. Renders nothing when the log is
       unavailable (the viewer lacks action-log.read), so non-admins simply do
       not see it. -->
  <section v-if="!failed">
    <div class="flex items-center justify-between mb-3">
      <h2 class="text-sm font-semibold uppercase tracking-wide text-muted-color">
        {{ $t('project.dashboard.eventsBand.title') }}
      </h2>
      <RouterLink
        :to="{ name: 'project.overview.events', params: { projectId } }"
        class="text-xs text-primary-500 hover:underline"
      >
        {{ $t('project.dashboard.eventsBand.viewAll') }}
      </RouterLink>
    </div>

    <div class="grid grid-cols-1 lg:grid-cols-2 gap-3">
      <!-- Pulse + stats -->
      <div class="rounded-lg border border-surface bg-surface-0 dark:bg-surface-900 p-4 flex flex-col gap-3">
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
        <CategoryTrendChart bare :height="88" :labels="metrics.labels" :series="metrics.series" />
      </div>

      <!-- Recent activity feed -->
      <div class="rounded-lg border border-surface bg-surface-0 dark:bg-surface-900 p-4">
        <div class="text-sm font-medium text-color mb-2">
          {{ $t('project.dashboard.eventsBand.recent') }}
        </div>
        <div v-if="recent.length" class="flex flex-col max-h-72 overflow-y-auto">
          <EventTimelineItem
            v-for="e in recent"
            :key="e.actionID"
            :data="e"
            :actor-name="events.actorName(e)"
          />
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
import { reactive, ref, watch } from 'vue'

const props = defineProps({
  projectId: { type: [String, Number], default: '' },
})

const events = useEventActivity()

const metrics = reactive({ labels: [], series: [], total: 0, actors: 0, errors: 0 })
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
