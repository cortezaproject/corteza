<template>
  <div>
    <span v-if="formatted">{{ formatted }}</span>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()

const props = defineProps({
  field: {
    type: Object,
    required: true,
  },
  record: {
    type: Object,
    required: true,
  },
  namespace: {
    type: Object,
    default: () => ({}),
  },
  valueOnly: {
    type: Boolean,
    default: false,
  },
  extraOptions: {
    type: Object,
    default: () => ({}),
  },
  disableClick: {
    type: Boolean,
    default: false,
  },
})

const value = computed(() => {
  if (props.field.isSystem) {
    return props.record[props.field.name]
  }
  return props.record?.values?.[props.field.name]
})

const formatted = computed(() => {
  const v = value.value
  if (!v && v !== 0) return null

  const delimiter = props.field.options?.multiDelimiter || ', '

  // Delegate to field.formatValue if available (handles moment-based formatting)
  if (props.field.formatValue) {
    if (props.field.isMulti && Array.isArray(v)) {
      return v.map(val => props.field.formatValue(val)).join(delimiter)
    }
    return props.field.formatValue(v)
  }

  // Fallback: use browser's locale formatting
  if (props.field.isMulti && Array.isArray(v)) {
    return v
      .map(val => formatDate(val))
      .filter(Boolean)
      .join(delimiter)
  }

  return formatDate(v)
})

function formatDate(val) {
  if (!val) return null

  try {
    const date = new Date(val)
    if (isNaN(date.getTime())) return String(val)

    const opts = props.field.options || {}

    if (opts.outputRelative) {
      return formatRelative(date)
    }

    if (opts.onlyTime) {
      return date.toLocaleTimeString()
    }

    if (opts.onlyDate) {
      return date.toLocaleDateString()
    }

    return date.toLocaleString()
  } catch {
    return String(val)
  }
}

function formatRelative(date) {
  const now = new Date()
  const diffMs = now.getTime() - date.getTime()
  const diffSecs = Math.floor(Math.abs(diffMs) / 1000)
  const diffMins = Math.floor(diffSecs / 60)
  const diffHours = Math.floor(diffMins / 60)
  const diffDays = Math.floor(diffHours / 24)

  const isFuture = diffMs < 0

  if (diffSecs < 60)
    return isFuture
      ? t('field.kind.dateTime.relative.inFewSeconds')
      : t('field.kind.dateTime.relative.fewSecondsAgo')
  if (diffMins < 60)
    return isFuture
      ? t('field.kind.dateTime.relative.inMinutes', { n: diffMins })
      : t('field.kind.dateTime.relative.minutesAgo', { n: diffMins })
  if (diffHours < 24)
    return isFuture
      ? t('field.kind.dateTime.relative.inHours', { n: diffHours })
      : t('field.kind.dateTime.relative.hoursAgo', { n: diffHours })
  if (diffDays < 30)
    return isFuture
      ? t('field.kind.dateTime.relative.inDays', { n: diffDays })
      : t('field.kind.dateTime.relative.daysAgo', { n: diffDays })

  return date.toLocaleDateString()
}
</script>
