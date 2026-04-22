<template>
  <div class="dynamic-form">
    <template v-for="segment in processedSegments" :key="segment.key">
      <!-- Segment title -->
      <div v-if="segment.title" class="text-sm font-medium text-color mb-2 mt-4 first:mt-0">
        {{ segment.title }}
      </div>

      <div v-for="section in segment.sections" :key="section.key" class="flex flex-col gap-3">
        <!-- Section title -->
        <div v-if="section.title" class="text-xs text-muted-color mb-2 mt-3">
          {{ section.title }}
        </div>

        <div class="flex flex-col gap-4">
          <DynamicInput
            v-for="input in section.inputs"
            :key="input.key"
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
        </div>
      </div>
    </template>
  </div>
</template>

<script setup>
import DynamicInput from './DynamicInput.vue'

defineProps({
  processedSegments: {
    type: Array,
    default: () => [],
  },
  showReferenceToggle: {
    type: Boolean,
    default: false,
  },
})

defineEmits(['update:value', 'toggleReference', 'clearReference', 'updateReferenceSource'])
</script>
