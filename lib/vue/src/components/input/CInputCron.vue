<template>
  <div class="flex flex-col gap-2 w-full">
    <Select
      v-model="selectedValue"
      :options="props.options"
      :optionLabel="optionLabel"
      optionValue="value"
      :placeholder="placeholder"
      :disabled="disabled"
      class="w-full"
      @change="onChange"
    >
      <template #option="slotProps">
        {{ getTranslatedLabel(slotProps.option.label) }}
      </template>
      <template #value="slotProps">
        <template v-if="slotProps.value">
           {{ getTranslatedLabel(props.options.find(o => o.value === slotProps.value)?.label) }}
        </template>
        <template v-else>
          {{ placeholder }}
        </template>
      </template>
    </Select>
  </div>
</template>

<script setup>
import { ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'

const { t, te } = useI18n()

const props = defineProps({
  modelValue: {
    type: [String, Array, Object],
    default: '',
  },
  options: {
    type: Array,
    default: () => [],
  },
  optionLabel: {
    type: String,
    default: 'label',
  },
  placeholder: {
    type: String,
    default: '',
  },
  disabled: {
    type: Boolean,
    default: false,
  },
})

const emit = defineEmits(['update:modelValue'])

const selectedValue = ref(null)

function getTranslatedLabel(label) {
  if (!label) return ''
  // Try to translate if an exact translation key exists or we can map it
  const keyMap = {
    'Every minute': 'builder.triggers.interval.everyMinute',
    'Every 5 minutes': 'builder.triggers.interval.every5Minutes',
    'Every 15 minutes': 'builder.triggers.interval.every15Minutes',
    'Every hour': 'builder.triggers.interval.everyHour',
    'Every day at midnight': 'builder.triggers.interval.everyDay',
    'Every Monday at midnight': 'builder.triggers.interval.everyMonday',
    'Every month on the 1st': 'builder.triggers.interval.everyMonth'
  }
  
  if (keyMap[label] && te(keyMap[label])) {
    return t(keyMap[label])
  } else if (te(label)) {
    return t(label)
  }
  return label
}

// Initialize from modelValue
watch(
  () => props.modelValue,
  (newVal) => {
    let actualValue = newVal
    while (actualValue !== null && actualValue !== undefined && typeof actualValue === 'object') {
      if (Array.isArray(actualValue)) {
        actualValue = actualValue.length > 0 ? actualValue[0] : ''
      } else if (actualValue.value !== undefined) {
        actualValue = actualValue.value
      } else {
        break
      }
    }

    if (!actualValue) {
      selectedValue.value = null
      return
    }
    
    selectedValue.value = actualValue
  },
  { immediate: true }
)

function onChange() {
  emit('update:modelValue', selectedValue.value)
}
</script>

