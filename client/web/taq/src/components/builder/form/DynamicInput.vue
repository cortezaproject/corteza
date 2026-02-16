<template>
  <div class="flex flex-col gap-1">
    <label v-if="label" class="text-sm font-medium text-primary">
      {{ label }}
      <span v-if="required" class="text-red-500">*</span>
    </label>

    <!-- Reference chip mode (only for non-aggregate inputs) -->
    <CReferenceChip
      v-if="isReference && !isAggregate"
      :label="referenceLabel"
      @click="$emit('toggleReference', argument)"
      @clear="$emit('clearReference', argument)"
    />

    <!-- Normal input mode — click opens reference panel (skip for aggregate types) -->
    <div v-else @click="!isAggregate && onInputClick()">
      <component
        :is="inputComponent"
        :model-value="modelValue"
        @update:model-value="$emit('update:modelValue', $event)"
        @toggle-row-reference="onRowReferenceToggle"
        :placeholder="effectivePlaceholder"
        :disabled="disabled"
        v-bind="$attrs"
      />
    </div>
  </div>
</template>

<script setup>
import { resolveInputComponent } from './inputs/registry'
import CReferenceChip from './CReferenceChip.vue'
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()

const props = defineProps({
  modelValue: {
    type: [String, Number, Object, Array],
    default: null,
  },
  type: {
    type: String,
    default: 'Text',
  },
  label: {
    type: String,
    default: '',
  },
  placeholder: {
    type: String,
    default: '',
  },
  disabledPlaceholder: {
    type: String,
    default: '',
  },
  required: {
    type: Boolean,
    default: false,
  },
  disabled: {
    type: Boolean,
    default: false,
  },
  argument: {
    type: String,
    default: '',
  },
  isReference: {
    type: Boolean,
    default: false,
  },
  referenceLabel: {
    type: String,
    default: '',
  },
  showReferenceToggle: {
    type: Boolean,
    default: false,
  },
})

const emit = defineEmits(['update:modelValue', 'toggleReference', 'clearReference'])

const inputComponent = computed(() => resolveInputComponent(props.type))

// Aggregate types handle their own per-row references (no whole-argument reference)
const AGGREGATE_TYPES = ['FieldValueMap']
const isAggregate = computed(() => AGGREGATE_TYPES.includes(props.type))

const effectivePlaceholder = computed(() => {
  if (props.disabled && props.disabledPlaceholder) {
    return props.disabledPlaceholder
  }
  return (
    props.placeholder ||
    (props.label ? t('builder.form.selectPlaceholder', { field: props.label.toLowerCase() }) : '')
  )
})

function onInputClick() {
  emit('toggleReference', props.argument)
}

// Handle per-row reference toggle from aggregate inputs (e.g. FieldValueMap)
function onRowReferenceToggle(targetField) {
  emit('toggleReference', { argument: props.argument, target: targetField })
}
</script>
