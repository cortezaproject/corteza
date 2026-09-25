<template>
  <Dialog
    :visible="visible"
    modal
    :header="title"
    class="w-full max-w-3xl"
    data-testid="dashboard-resource-dialog"
    @update:visible="$emit('update:visible', $event)"
  >
    <div v-if="error" class="text-base text-red-500">{{ $t('dashboard.loadFailed') }}</div>

    <div v-else-if="!detail" class="flex items-center justify-center py-10">
      <ProgressSpinner style="width: 32px; height: 32px" />
    </div>

    <div v-else class="flex flex-col gap-5">
      <div class="grid grid-cols-2 sm:grid-cols-4 gap-3">
        <div v-for="f in figures" :key="f.key" class="flex flex-col gap-0.5">
          <span class="text-sm text-muted-color">{{ f.label }}</span>
          <span class="text-2xl font-semibold leading-none text-color">
            {{ f.value.toLocaleString() }}
          </span>
        </div>
      </div>

      <StatusBar :status="detail.status" :order="order" />

      <section class="flex flex-col gap-2">
        <h3 class="text-base font-semibold text-color">{{ $t('dashboard.dialog.movement') }}</h3>
        <TrendChart
          :labels="bucketLabels"
          :range-labels="rangeLabels"
          :series="series"
          :height="180"
        />
      </section>

      <section class="flex flex-col gap-2">
        <h3 class="text-base font-semibold text-color">{{ $t('dashboard.dialog.recent') }}</h3>
        <div v-if="!detail.recent.length" class="text-base text-muted-color">
          {{ $t('dashboard.empty') }}
        </div>
        <ul v-else class="divide-y divide-surface">
          <li
            v-for="it in detail.recent"
            :key="it.id"
            class="flex items-center gap-3 py-1.5 min-w-0"
          >
            <span
              class="inline-block h-2 w-2 shrink-0 rounded-full"
              :style="{ backgroundColor: statusColor(it.status) }"
            />
            <div class="flex-1 min-w-0">
              <component
                :is="linkTo(it) ? 'router-link' : 'span'"
                :to="linkTo(it)"
                class="text-base text-color break-words"
                :class="linkTo(it) ? 'hover:underline' : ''"
              >
                {{ it.label || it.handle || it.id }}
              </component>
              <div
                v-if="it.handle && it.handle !== it.label"
                class="text-sm text-muted-color break-all"
              >
                {{ it.handle }}
              </div>
            </div>
            <span class="text-sm text-muted-color whitespace-nowrap">
              {{ statusLabel(it.status) }}
            </span>
            <span class="text-sm text-muted-color whitespace-nowrap">
              {{ locDate(it.createdAt) }}
            </span>
          </li>
        </ul>
      </section>
    </div>

    <template #footer>
      <Button :label="$t('dashboard.dialog.openList')" severity="secondary" @click="openList" />
      <Button :label="$t('general.label.close')" @click="$emit('update:visible', false)" />
    </template>
  </Dialog>
</template>

<script setup>
import { computed, inject, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRouter } from 'vue-router'
import { filters } from '@planetcrust/human-vue'
import StatusBar from './StatusBar.vue'
import TrendChart from './TrendChart.vue'
import { SERIES, seriesColor, statusColor, useIsDark } from './chartTheme'
import { RESOURCES, orderedStatuses } from './resources'
import { bucketLabelsFor, rangeParams } from './useSystemStats'

const { locDate } = filters

const props = defineProps({
  visible: { type: Boolean, default: false },
  // inventory resource key (users, roles, ...)
  resource: { type: String, default: '' },
  // active range preset key, so the dialog shows the same window as the page
  range: { type: String, default: '30d' },
})
const emit = defineEmits(['update:visible'])

const { t, te } = useI18n()
const router = useRouter()
const $SystemAPI = inject('$SystemAPI')
const isDark = useIsDark()

const detail = ref(null)
const error = ref(null)

const config = computed(() => RESOURCES.find(r => r.key === props.resource))
const title = computed(() => (props.resource ? t(`dashboard.resources.${props.resource}`) : ''))

// Each opening fetches afresh; a stale answer for a previous resource is dropped.
let pending = null
async function load() {
  detail.value = null
  error.value = null
  if (!props.resource) return

  const req = $SystemAPI.statsDetailCancellable({
    resource: props.resource,
    ...rangeParams(props.range),
  })
  pending?.cancel()
  pending = req

  try {
    const result = await req.response()
    if (pending === req) detail.value = result
  } catch (e) {
    if (pending === req && !e?.__CANCEL__) error.value = e
  }
}

watch(
  () => [props.visible, props.resource],
  ([visible]) => {
    if (visible) load()
  },
  { immediate: true },
)

const order = computed(() =>
  config.value && detail.value ? orderedStatuses(config.value, detail.value.status) : [],
)

const figures = computed(() => {
  if (!detail.value) return []
  const d = detail.value
  return [
    { key: 'total', label: t('dashboard.inventory.total'), value: d.total },
    { key: 'created', label: t('dashboard.dialog.created'), value: d.inRange.created },
    { key: 'updated', label: t('dashboard.dialog.updated'), value: d.inRange.updated },
    { key: 'deleted', label: t('dashboard.dialog.deleted'), value: d.inRange.deleted },
  ]
})

const bucketLabels = computed(() => (detail.value ? bucketLabelsFor(detail.value.range).short : []))
const rangeLabels = computed(() => (detail.value ? bucketLabelsFor(detail.value.range).long : []))

const series = computed(() =>
  detail.value
    ? [
        {
          key: 'created',
          name: t('dashboard.dialog.created'),
          color: seriesColor('created', isDark.value),
          data: detail.value.series.created,
        },
        {
          key: 'updated',
          name: t('dashboard.dialog.updated'),
          color: SERIES.running,
          data: detail.value.series.updated,
        },
        {
          key: 'deleted',
          name: t('dashboard.dialog.deleted'),
          color: SERIES.canceled,
          data: detail.value.series.deleted,
        },
      ]
    : [],
)

function statusLabel(key) {
  return te(`dashboard.status.${key}`) ? t(`dashboard.status.${key}`) : key
}

// Editor routes per resource; resources without a direct editor route link nowhere.
const EDIT_ROUTES = {
  users: id => ({ name: 'system.users.edit', params: { userID: id } }),
  roles: id => ({ name: 'system.roles.edit', params: { roleID: id } }),
  applications: id => ({ name: 'system.applications.edit', params: { applicationID: id } }),
  authClients: id => ({ name: 'system.authClients.edit', params: { authClientID: id } }),
  workflows: id => ({ name: 'automation.workflows.edit', params: { workflowID: id } }),
  taqs: id => ({ name: 'automation.taq.edit', params: { automationID: id } }),
  agents: id => ({ name: 'agentic.edit', params: { agentID: id } }),
  chatbots: id => ({ name: 'chatbot.edit', params: { chatbotID: id } }),
  projects: id => ({ name: 'project.overview', params: { projectId: id } }),
}

function linkTo(it) {
  const f = EDIT_ROUTES[props.resource]
  return f && !it.deletedAt ? f(it.id) : null
}

function openList() {
  emit('update:visible', false)
  if (config.value) router.push(config.value.route)
}
</script>
