<template>
  <!-- Quick time-range presets (stock-chart style). Purely the button row —
       callers own what a preset selection actually does (ActivityPanel wires
       it to its filter window and releases the selection on manual edits;
       Overview/CategoryView wire it straight to their trend reload), so this
       stays a dumb v-model over the preset key. -->
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
import { RANGES } from '@/sections/project/config/trend'
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'

defineProps({
  // The active preset key (one of RANGES' keys), or null when the selection
  // has been released (e.g. a manual filter edit) — SelectButton then renders
  // with nothing pressed.
  modelValue: { type: String, default: null },
})
defineEmits(['update:modelValue'])

const { t } = useI18n()

const options = computed(() =>
  RANGES.map(r => ({ key: r.key, label: t(`project.dashboard.range.${r.key}`) })),
)
</script>
