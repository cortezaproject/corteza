<template>
  <PageBlock :block="block" @refreshBlock="refresh">
    <div v-if="loading" class="flex items-center justify-center h-full">
      <ProgressSpinner style="width: 28px; height: 28px" />
    </div>

    <div v-else class="relative h-full">
      <ProgressBar
        :value="percentage"
        :show-value="false"
        :style="progressStyle"
        class="w-full h-full"
      />
      <div v-if="showValue" class="absolute inset-0 flex items-center justify-center pointer-events-none">
        <span class="text-sm font-medium" :style="textStyle">{{ displayLabel }}</span>
      </div>
    </div>
  </PageBlock>
</template>

<script setup>
import { ref, computed, watch, onMounted, onBeforeUnmount, inject } from 'vue'
import PageBlock from './PageBlock.vue'
import { evaluatePrefilter } from '../../../lib/record-filter'

const props = defineProps({
  block: { type: Object, required: true },
  namespace: { type: Object, default: () => ({}) },
  page: { type: Object, default: () => ({}) },
  record: { type: Object, default: undefined },
})

const $ComposeAPI = inject('$ComposeAPI', null)
const $auth = inject('$auth', {})
const $eventBus = inject('$eventBus', null)

const loading = ref(false)
const value = ref(0)
const min = ref(0)
const max = ref(100)

const options = computed(() => props.block.options || {})
const displayOpts = computed(() => options.value.display || {})
const showValue = computed(() => displayOpts.value.showValue !== false)

const percentage = computed(() => {
  const range = max.value - min.value
  if (range <= 0) return 0
  return Math.min(100, Math.max(0, ((value.value - min.value) / range) * 100))
})

const displayLabel = computed(() => {
  const showProgress = displayOpts.value.showProgress
  if (displayOpts.value.showRelative) {
    const pct = `${Math.round(percentage.value)}%`
    return showProgress ? `${pct} / 100%` : pct
  }
  return showProgress ? `${value.value} / ${max.value}` : `${value.value}`
})

// Map progress variant → PrimeVue Button severity token suffix
// Reading the same --p-button-{severity}-* tokens the Button component uses guarantees
// the progress fill stays in sync with the configurator's Button tags across themes.
const severityTokenMap = {
  primary: 'primary',
  secondary: 'secondary',
  success: 'success',
  warning: 'warn',
  danger: 'danger',
  info: 'info',
  dark: 'contrast',
}

const progressStyle = computed(() => {
  const sev = severityTokenMap[getActiveVariant()] || 'primary'
  return {
    '--p-progressbar-value-background': `var(--p-button-${sev}-background)`,
  }
})

const textStyle = computed(() => {
  const sev = severityTokenMap[getActiveVariant()] || 'primary'
  return { color: `var(--p-button-${sev}-color)` }
})

function getActiveVariant() {
  const thresholds = displayOpts.value.thresholds || []
  const defaultVariant = displayOpts.value.variant || 'primary'

  // Sort thresholds ascending
  const sorted = [...thresholds].sort((a, b) => (a.value || 0) - (b.value || 0))

  let variant = defaultVariant
  for (const th of sorted) {
    if (value.value >= (th.value || 0)) {
      variant = th.variant || defaultVariant
    }
  }
  return variant
}

async function refresh() {
  if (!$ComposeAPI || !props.block.fetch) {
    // Use default value if no API available
    value.value = options.value.value?.default || 0
    min.value = options.value.minValue?.default || 0
    max.value = options.value.maxValue?.default || 100
    return
  }

  loading.value = true

  try {
    const record = props.record
    const user = $auth?.user || {}

    const additionalOptions = {
      value: {
        filter: evaluatePrefilter(options.value.value?.filter || '', {
          record, user,
          recordID: record?.recordID || '0',
          ownerID: record?.ownedBy || '0',
          userID: user?.userID || '0',
        }),
      },
      minValue: {
        filter: evaluatePrefilter(options.value.minValue?.filter || '', {
          record, user,
          recordID: record?.recordID || '0',
          ownerID: record?.ownedBy || '0',
          userID: user?.userID || '0',
        }),
      },
      maxValue: {
        filter: evaluatePrefilter(options.value.maxValue?.filter || '', {
          record, user,
          recordID: record?.recordID || '0',
          ownerID: record?.ownedBy || '0',
          userID: user?.userID || '0',
        }),
      },
    }

    const result = await props.block.fetch(additionalOptions, $ComposeAPI, props.namespace.namespaceID)
    value.value = result.value ?? 0
    min.value = result.min ?? 0
    max.value = result.max ?? 100
  } catch (e) {
    console.error('Failed to fetch progress data:', e)
    value.value = options.value.value?.default || 0
    min.value = options.value.minValue?.default || 0
    max.value = options.value.maxValue?.default || 100
  } finally {
    loading.value = false
  }
}

onMounted(() => refresh())

watch(() => props.record?.recordID, () => refresh())
watch(() => props.block.options, () => refresh(), { deep: true })

const offRefetch = $eventBus?.on('refetch-records', () => refresh())

onBeforeUnmount(() => {
  offRefetch?.()
})
</script>
