<template>
  <div class="dynamic-form">
    <template v-for="(segment, sIdx) in functionDef.segments" :key="`segment-${sIdx}`">
      <!-- Segment title -->
      <div v-if="segment.meta?.short" class="text-sm font-medium text-color mb-2 mt-4 first:mt-0">
        {{ segment.meta.short }}
      </div>

      <template v-for="(section, secIdx) in segment.sections" :key="`section-${sIdx}-${secIdx}`">
        <!-- Section title -->
        <div v-if="section.meta?.short" class="text-xs text-muted-color mb-2 mt-3">
          {{ section.meta.short }}
        </div>

        <div class="flex flex-col gap-3">
          <template v-for="(element, elIdx) in section.elements" :key="`element-${sIdx}-${secIdx}-${elIdx}`">
            <DynamicInput
              v-if="element.input"
              :type="element.input.type"
              :label="element.input.label"
              :placeholder="element.input.placeholder"
              :required="isRequired(element.input.argument)"
              :visual="element.input.visual"
              :model-value="getConfigValue(element.input.argument)"
              @update:model-value="updateConfig(element.input.argument, $event)"
            />
          </template>
        </div>
      </template>
    </template>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import DynamicInput from './DynamicInput.vue'

const props = defineProps({
  functionDef: {
    type: Object,
    required: true,
  },
  config: {
    type: Object,
    default: () => ({}),
  },
})

const emit = defineEmits(['update:config'])

const parameters = computed(() => {
  return props.functionDef.parameters || []
})

function isRequired(argumentName) {
  const param = parameters.value.find(p => p.name === argumentName)
  return param?.required || false
}

function getConfigValue(argumentName) {
  return props.config[argumentName] ?? null
}

function updateConfig(argumentName, value) {
  emit('update:config', { key: argumentName, value })
}
</script>
