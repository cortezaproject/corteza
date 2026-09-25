<template>
  <SelectButton
    :model-value="modelValue"
    :options="options"
    option-label="label"
    option-value="key"
    size="small"
    :allow-empty="false"
    @update:model-value="$emit('update:modelValue', $event)"
  />
</template>

<script setup>
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { RANGES } from './useSystemStats'

defineProps({ modelValue: { type: String, default: '30d' } })
defineEmits(['update:modelValue'])

const { t } = useI18n()
const options = computed(() =>
  RANGES.map(r => ({ key: r.key, label: t(`dashboard.range.${r.key}`) })),
)
</script>
