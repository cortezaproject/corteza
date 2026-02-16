<template>
  <div>
    <span>{{ formatted }}</span>
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
  const { trueLabel, falseLabel } = props.field.options || {}
  const v = value.value

  if (v === '1' || v === true) {
    return trueLabel || t('general.label.yes', 'Yes')
  } else if (v === '0' || v === false) {
    return falseLabel || t('general.label.no', 'No')
  }

  return ''
})
</script>
