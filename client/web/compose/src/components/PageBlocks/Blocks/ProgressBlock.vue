<template>
  <PageBlock :block="block">
    <div v-if="loading" class="flex items-center justify-center h-full">
      <ProgressSpinner style="width: 28px; height: 28px" />
    </div>

    <div v-else class="flex items-center h-full p-3">
      <div class="flex-1">
        <ProgressBar
          :value="percentage"
          :show-value="showValue"
          :style="progressStyle"
        >
          <template v-if="showValue" #default>
            <span class="text-sm font-medium">
              {{ displayLabel }}
            </span>
          </template>
        </ProgressBar>
      </div>
    </div>
  </PageBlock>
</template>

<script setup>
import { ref, computed, watch, onMounted, onBeforeUnmount, inject } from 'vue'
import { useI18n } from 'vue-i18n'
import PageBlock from './PageBlock.vue'
import { evaluatePrefilter } from '../../../lib/record-filter'

const { t } = useI18n()

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
  if (displayOpts.value.showRelative) {
    return `${Math.round(percentage.value)}%`
  }
  return `${value.value}`
})

const progressStyle = computed(() => {
  const variant = getActiveVariant()
  const colorMap = {
    primary: undefined,
    success: 'var(--p-green-500)',
    warning: 'var(--p-yellow-500)',
    danger: 'var(--p-red-500)',
    info: 'var(--p-blue-500)',
    dark: 'var(--p-surface-700)',
    secondary: 'var(--p-surface-400)',
  }
  const color = colorMap[variant]
  return color ? { '--p-progressbar-value-background': color } : {}
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
