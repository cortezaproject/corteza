<template>
  <div class="flex flex-col gap-2 w-full">
    <Select
      v-model="selectedPreset"
      :options="mergedOptions"
      :optionLabel="optionLabel"
      optionValue="value"
      :placeholder="placeholder"
      :disabled="disabled"
      class="w-full"
      @change="onPresetChange"
    >
      <template #option="slotProps">
        {{ getTranslatedLabel(slotProps.option.label) }}
      </template>
      <template #value="slotProps">
        <template v-if="slotProps.value">
           {{ getTranslatedLabel(mergedOptions.find(o => o.value === slotProps.value)?.label) }}
        </template>
        <template v-else>
          {{ placeholder }}
        </template>
      </template>
    </Select>
    
    <InputText
      v-if="selectedPreset === 'custom'"
      v-model="customValue"
      placeholder="* * * * *"
      :disabled="disabled"
      @input="onCustomInput"
    />
  </div>
</template>

<script setup>
import { ref, watch, computed } from 'vue'
import { useI18n } from 'vue-i18n'

const { t, te } = useI18n()

const props = defineProps({
  modelValue: {
    type: String,
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

const selectedPreset = ref(null)
const customValue = ref('')

const mergedOptions = computed(() => {
  return [
    ...props.options,
    { label: 'Custom', value: 'custom' }
  ]
})

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
    'Every month on the 1st': 'builder.triggers.interval.everyMonth',
    'Custom': 'builder.triggers.interval.custom'
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
    if (!newVal) {
      selectedPreset.value = null
      customValue.value = ''
      return
    }
    
    const matchedOption = props.options.find(o => o.value === newVal)
    if (matchedOption) {
      selectedPreset.value = matchedOption.value
      customValue.value = ''
    } else {
      selectedPreset.value = 'custom'
      customValue.value = newVal
    }
  },
  { immediate: true }
)

function onPresetChange() {
  if (selectedPreset.value === 'custom') {
    emit('update:modelValue', customValue.value)
  } else {
    customValue.value = ''
    emit('update:modelValue', selectedPreset.value)
  }
}

function onCustomInput() {
  emit('update:modelValue', customValue.value)
}
</script>
