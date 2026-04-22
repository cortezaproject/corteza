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
      @update:label="$emit('updateReferenceSource', argument, $event)"
    />

    <!-- Normal input mode -->
    <div v-else class="w-full">
      <div :class="{ 'flex gap-1 items-center': showReferenceToggle && !isAggregate }">
        <component
          :is="inputComponent"
          class="w-full"
          :model-value="displayValue"
          @update:model-value="onValueUpdate"
          @toggle-row-reference="onRowReferenceToggle"
          :placeholder="effectivePlaceholder"
          :disabled="disabled"
          :options="options"
          :complete-on-focus="hasOptions"
          v-bind="$attrs"
          @focus="onInputFocusOrClick"
          @click="onInputFocusOrClick"
        />
        <Button
          v-if="showReferenceToggle && !isAggregate"
          icon="pi pi-link"
          text
          rounded
          size="small"
          :severity="isReferenceActive ? 'primary' : 'secondary'"
          :title="$t('builder.form.referenceToggle')"
          @click.stop="onInputClick()"
        />
      </div>
    </div>
  </div>
</template>

<script setup>
import { resolveInputComponent } from './inputs/registry'
import CInputArray from './inputs/CInputArray.vue'
import CReferenceChip from './CReferenceChip.vue'
import { computed, inject, ref } from 'vue'

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
  options: {
    type: Array,
    default: () => [],
  },
  aggregate: {
    type: Boolean,
    default: false,
  },
})

const emit = defineEmits(['update:modelValue', 'toggleReference', 'clearReference', 'updateReferenceSource'])

const inputComponent = computed(() => isAggregate.value ? CInputArray : resolveInputComponent(props.type))

// Injected from Builder.vue — tracks which argument has the reference panel open
const activeReferenceArgument = inject('activeReferenceArgument', ref(null))
const isReferenceActive = computed(() => {
  if (!activeReferenceArgument.value) return false
  return activeReferenceArgument.value.name === props.argument
})

// Aggregate types handle their own per-row references (no whole-argument reference)
const AGGREGATE_TYPES = ['FieldValueMap', 'Array']
const isAggregate = computed(() => props.aggregate || AGGREGATE_TYPES.includes(props.type))

// Options support: map stored ID ↔ option object for CInputSelect
const hasOptions = computed(() => props.options.length > 0)

const displayValue = computed(() => {
  if (hasOptions.value && props.modelValue) {
    return props.options.find(o => o.value === props.modelValue) || props.modelValue
  }
  return props.modelValue
})

function onValueUpdate(val) {
  if (hasOptions.value && val && typeof val === 'object' && val.value !== undefined) {
    emit('update:modelValue', val.value)
  } else {
    emit('update:modelValue', val)
  }
}

const effectivePlaceholder = computed(() => {
  if (props.disabled && props.disabledPlaceholder) {
    return props.disabledPlaceholder
  }

  return props.placeholder
})

function onInputClick() {
  emit('toggleReference', props.argument)
}

function onInputFocusOrClick() {
  if (props.showReferenceToggle && !isAggregate.value) {
    emit('toggleReference', props.argument)
  }
}

// Handle per-row reference toggle from aggregate inputs (e.g. FieldValueMap)
function onRowReferenceToggle(targetField) {
  emit('toggleReference', { argument: props.argument, target: targetField })
}
</script>
