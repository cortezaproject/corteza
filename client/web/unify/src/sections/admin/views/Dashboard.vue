<template>
  <Teleport to="#topbar-title" defer>
    <span>{{ $t('dashboard.title') }}</span>
  </Teleport>

  <div class="flex flex-col h-full p-3 sm:p-6 gap-4 overflow-y-auto min-w-0">
    <div class="flex flex-wrap items-center justify-end gap-2">
      <CManualScriptButtons
        resource-type="system"
        ui-page="dashboard"
        ui-slot="toolbar"
        container-class="flex flex-wrap gap-2"
        @click="handleScriptButton"
      />
      <RangeSelect v-model="range" />
      <Button
        icon="pi pi-refresh"
        severity="secondary"
        text
        size="small"
        :loading="loading"
        :aria-label="$t('dashboard.reload')"
        @click="reload"
      />
    </div>

    <Message v-if="error" severity="error" :closable="false">
      {{ $t('dashboard.loadFailed') }}
      <Button :label="$t('dashboard.reload')" link size="small" @click="reload" />
    </Message>

    <div v-if="!stats && loading" class="flex items-center justify-center flex-1">
      <ProgressSpinner style="width: 40px; height: 40px" />
    </div>

    <!-- Resource first: what exists, in what state, then how it moves. -->
    <div
      v-else-if="stats"
      class="flex flex-col gap-4 min-w-0"
      :class="{ 'opacity-60 transition-opacity': loading }"
    >
      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-3">
        <button
          v-for="tile in d.tiles.value"
          :key="tile.key"
          type="button"
          class="flex flex-col gap-2 rounded-border border border-surface bg-surface px-4 py-3 text-left hover:border-primary transition-colors min-w-0"
          :data-testid="`dashboard-tile-${tile.key}`"
          @click="inspect(tile.key)"
        >
          <div class="flex items-baseline justify-between gap-2">
            <span class="flex items-center gap-2 min-w-0">
              <i :class="tile.icon" class="text-muted-color" aria-hidden="true" />
              <span class="text-base font-medium text-color truncate">{{ tile.label }}</span>
            </span>
            <span v-if="tile.created" class="text-sm text-muted-color whitespace-nowrap">
              {{ $t('dashboard.tile.newInRange', { n: tile.created }) }}
            </span>
          </div>
          <div class="flex items-end justify-between gap-3">
            <div>
              <span class="text-3xl font-semibold leading-none text-color">
                {{ tile.live.toLocaleString() }}
              </span>
              <span v-if="tile.total !== tile.live" class="text-base text-muted-color ml-1">
                / {{ tile.total.toLocaleString() }}
              </span>
            </div>
            <Sparkline :values="tile.spark" :color="d.color('created')" class="shrink-0" />
          </div>
          <StatusBar :status="tile.status" :order="orderedStatuses(tile.resource, tile.status)" />
        </button>
      </div>

      <div class="grid grid-cols-1 xl:grid-cols-2 gap-4">
        <Card
          v-if="d.activity.value"
          :title="$t('dashboard.activity.title')"
          :subtitle="
            $t('dashboard.activity.subtitle', {
              entries: d.activity.value.total.toLocaleString(),
              errors: d.activity.value.errors.toLocaleString(),
            })
          "
        >
          <TrendChart
            :labels="bucketLabels"
            :range-labels="rangeLabels"
            :series="d.activity.value.series"
            :height="180"
            @select="openBucket('activity', $event)"
          />
        </Card>
        <Card
          v-if="d.signins.value"
          :title="$t('dashboard.signins.title')"
          :subtitle="
            $t('dashboard.signins.subtitle', {
              n: d.signins.value.users,
              live: d.signins.value.live,
            })
          "
        >
          <TrendChart
            :labels="bucketLabels"
            :range-labels="rangeLabels"
            :series="d.signins.value.series"
            :height="180"
            @select="openBucket('signins', $event)"
          />
        </Card>
      </div>

      <div class="grid grid-cols-1 xl:grid-cols-2 gap-4">
        <Card
          v-if="d.workflows.value || d.taqs.value"
          :title="$t('dashboard.runs.title')"
          :subtitle="runsSubtitle"
        >
          <div class="flex flex-col gap-5">
            <div v-if="d.workflows.value" class="min-w-0">
              <div class="text-sm text-muted-color mb-1">
                {{ $t('dashboard.runs.workflowRuns') }} ·
                {{
                  $t('dashboard.runs.subtitle', {
                    n: d.workflows.value.total,
                    failed: d.workflows.value.failed,
                  })
                }}
              </div>
              <TrendChart
                :labels="bucketLabels"
                :range-labels="rangeLabels"
                :series="d.workflows.value.series"
                stacked
                :height="180"
                @select="openBucket('workflows', $event)"
              />
            </div>
            <div v-if="d.taqs.value" class="min-w-0">
              <div class="text-sm text-muted-color mb-1">
                {{ $t('dashboard.runs.taqRuns') }} ·
                {{
                  $t('dashboard.runs.subtitle', {
                    n: d.taqs.value.total,
                    failed: d.taqs.value.failed,
                  })
                }}
              </div>
              <TrendChart
                :labels="bucketLabels"
                :range-labels="rangeLabels"
                :series="d.taqs.value.series"
                stacked
                :height="180"
                @select="openBucket('taqs', $event)"
              />
            </div>
          </div>
        </Card>

        <Card
          :title="$t('dashboard.attention.title')"
          :subtitle="$t('dashboard.attention.subtitle')"
        >
          <AttentionList
            :items="d.attention.value"
            :empty-text="$t('dashboard.attention.empty')"
            @select="openEvent"
          />
        </Card>
      </div>
    </div>

    <ResourceDialog v-model:visible="dialogOpen" :resource="dialogResource" :range="range" />
    <BucketDialog
      v-model:visible="bucketOpen"
      :kind="bucket.kind"
      :from="bucket.from"
      :to="bucket.to"
      :label="bucket.label"
      @select="openEvent"
    />
    <EventDialog v-model:visible="eventOpen" :item="event" />
  </div>
</template>

<script setup>
import { computed, inject, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { system } from '@planetcrust/human-js'
import { components } from '@planetcrust/human-vue'
import AttentionList from '../components/Dashboard/AttentionList.vue'
import BucketDialog from '../components/Dashboard/BucketDialog.vue'
import EventDialog from '../components/Dashboard/EventDialog.vue'
import Card from '../components/Dashboard/Card.vue'
import RangeSelect from '../components/Dashboard/RangeSelect.vue'
import ResourceDialog from '../components/Dashboard/ResourceDialog.vue'
import Sparkline from '../components/Dashboard/Sparkline.vue'
import StatusBar from '../components/Dashboard/StatusBar.vue'
import TrendChart from '../components/Dashboard/TrendChart.vue'
import { useIsDark } from '../components/Dashboard/chartTheme'
import { orderedStatuses } from '../components/Dashboard/resources'
import { useDashboardData } from '../components/Dashboard/useDashboardData'
import { useSystemStats } from '../components/Dashboard/useSystemStats'

const { CManualScriptButtons } = components
const { t } = useI18n()

const $ScriptBus = inject('$ScriptBus', null)
const $toast = inject('$toast', null)

const { range, stats, loading, error, reload, bucketLabels, rangeLabels } = useSystemStats()
const isDark = useIsDark()
const d = useDashboardData(stats, isDark, bucketLabels, rangeLabels)

const runsSubtitle = computed(() => {
  const total = (d.workflows.value?.total || 0) + (d.taqs.value?.total || 0)
  const failed = (d.workflows.value?.failed || 0) + (d.taqs.value?.failed || 0)
  return t('dashboard.runs.subtitle', { n: total, failed })
})

// A tile opens its resource in depth; the dialog fetches for the page's range.
const dialogOpen = ref(false)
const dialogResource = ref('')

function inspect(key) {
  dialogResource.value = key
  dialogOpen.value = true
}

// A clicked bar opens the events behind that bucket: from its start to the
// next bucket's start, or the range end for the last one.
const bucketOpen = ref(false)
const bucket = ref({ kind: '', from: '', to: '', label: '' })

function openBucket(kind, index) {
  const buckets = stats.value?.range?.buckets || []
  const day = buckets[index]
  if (!day) return

  const from = localDay(day)
  const to = buckets[index + 1] ? localDay(buckets[index + 1]) : new Date(stats.value.range.to)
  bucket.value = {
    kind,
    from: from.toISOString(),
    to: to.toISOString(),
    label: rangeLabels.value[index] || day,
  }
  bucketOpen.value = true
}

function localDay(s) {
  const [y, m, d] = s.split('-').map(Number)
  return new Date(y, m - 1, d)
}

// Any listed event opens in depth
const eventOpen = ref(false)
const event = ref(null)

function openEvent(item) {
  event.value = item
  eventOpen.value = true
}

async function handleScriptButton(button) {
  if (!$ScriptBus) return

  try {
    await $ScriptBus.Dispatch(
      { ...system.SystemEvent(), resourceType: button.resourceType, args: {} },
      button.script,
    )
  } catch (e) {
    console.error('Automation script failed:', e)
    $toast?.toastErrorHandler?.(t('notification.automation.scriptFailed'))(e)
  }
}
</script>
