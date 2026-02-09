<template>
  <div class="flex flex-col gap-1">
    <label v-if="label" class="text-sm font-medium text-primary">
      {{ label }}
      <span v-if="required" class="text-red-500">*</span>
    </label>
    <component
      :is="inputComponent"
      :model-value="modelValue"
      @update:model-value="$emit('update:modelValue', $event)"
      :placeholder="effectivePlaceholder"
      :disabled="disabled"
      v-bind="$attrs"
    />
  </div>
</template>

<script setup>
import { resolveInputComponent } from '@cortezaproject/corteza-vue-next/src/components/input/registry'
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
})

defineEmits(['update:modelValue'])

const inputComponent = computed(() => resolveInputComponent(props.type))

const effectivePlaceholder = computed(() => {
  if (props.disabled && props.disabledPlaceholder) {
    return props.disabledPlaceholder
  }
  return props.placeholder || (props.label ? t('builder.form.selectPlaceholder', { field: props.label.toLowerCase() }) : '')
})
</script>
