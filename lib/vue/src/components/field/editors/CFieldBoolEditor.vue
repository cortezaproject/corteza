<template>
  <!-- Switch: No ──[switch]── Yes -->
  <div v-if="field.options?.switch" class="flex items-center gap-2">
    <span class="text-sm text-muted-color">{{ falseLabel }}</span>
    <ToggleSwitch
      :model-value="boolValue"
      :disabled="disabled"
      @update:model-value="onUpdate"
    />
    <span class="text-sm text-muted-color">{{ trueLabel }}</span>
  </div>

  <!-- Checkbox: [✓] Field label -->
  <div v-else class="flex items-center gap-2">
    <Checkbox
      :input-id="inputId"
      :model-value="boolValue"
      :binary="true"
      :disabled="disabled"
      @update:model-value="onUpdate"
    />
    <label :for="inputId" class="cursor-pointer text-sm">
      {{ field.label || field.name }}
    </label>
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
  modelValue: {
    type: String,
    default: '',
  },
  disabled: {
    type: Boolean,
    default: false,
  },
})

const emit = defineEmits(['update:modelValue'])

const boolValue = computed(() => props.modelValue === '1')
const inputId = computed(() => `bool-${props.field.fieldID || props.field.name}`)
const trueLabel = computed(() => props.field.options?.trueLabel || t('general.label.yes'))
const falseLabel = computed(() => props.field.options?.falseLabel || t('general.label.no'))

function onUpdate(value) {
  emit('update:modelValue', value ? '1' : '')
}
</script>
