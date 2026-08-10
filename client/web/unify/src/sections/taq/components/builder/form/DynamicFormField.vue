<template>
  <div class="flex flex-col gap-1">
    <DynamicInput
      :type="input.type"
      :label="input.label"
      :placeholder="input.placeholder"
      :disabled-placeholder="input.disabledPlaceholder"
      :disabled="input.disabled"
      :required="input.required"
      :argument="input.argument"
      :options="input.options || []"
      :aggregate="input.isAggregate || false"
      :is-reference="input.isReference"
      :reference-label="input.referenceLabel"
      :show-reference-toggle="showReferenceToggle"
      v-bind="input.contextProps"
      :model-value="input.value"
      @update:model-value="$emit('update:value', input.argument, $event)"
      @toggle-reference="$emit('toggleReference', $event)"
      @clear-reference="$emit('clearReference', $event)"
      @update-reference-source="(arg, val) => $emit('updateReferenceSource', arg, val)"
    />
    <p v-if="input.description" class="text-xs text-muted-color m-0">
      {{ input.description }}
    </p>
  </div>
</template>

<script setup>
import DynamicInput from './DynamicInput.vue'

defineProps({
  input: { type: Object, required: true },
  showReferenceToggle: { type: Boolean, default: false },
})

defineEmits(['update:value', 'toggleReference', 'clearReference', 'updateReferenceSource'])
</script>
