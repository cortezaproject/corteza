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
            :required="input.required"
            v-bind="input.contextProps"
            :model-value="input.value"
            @update:model-value="updateConfig(input.argument, $event)"
          />
        </div>
      </div>
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

const parameters = computed(() => props.functionDef.parameters || [])

function isRequired(argumentName) {
  return parameters.value.find(p => p.name === argumentName)?.required || false
}

function resolveContextProps(context) {
  if (!context?.dependsOn) return {}

  const resolved = {}
  // Resolve dependencies: { "namespaceID": "namespace" } -> { namespaceID: config.namespace }
  for (const [propName, sourceArgument] of Object.entries(context.dependsOn)) {
    resolved[propName] = props.config[sourceArgument] ?? null
  }

  return resolved
}

// Pre-process all segments/sections/inputs with resolved values
// This computed tracks props.config changes properly
const processedSegments = computed(() => {
  return (props.functionDef.segments || []).map((segment, sIdx) => ({
    key: `segment-${sIdx}`,
    title: segment.meta?.short || null,
    sections: (segment.sections || []).map((section, secIdx) => ({
      key: `section-${sIdx}-${secIdx}`,
      title: section.meta?.short || null,
      inputs: (section.elements || [])
        .filter(el => el.input)
        .map((element, elIdx) => ({
          key: `input-${sIdx}-${secIdx}-${elIdx}`,
          type: element.input.type,
          label: element.input.label,
          placeholder: element.input.placeholder,
          argument: element.input.argument,
          required: isRequired(element.input.argument),
          contextProps: resolveContextProps(element.input.context),
          value: props.config[element.input.argument] ?? null,
        })),
    })),
  }))
})

function getDependentArguments(sourceArgument) {
  const dependents = []
  for (const segment of props.functionDef.segments || []) {
    for (const section of segment.sections || []) {
      for (const element of section.elements || []) {
        if (!element.input?.context?.dependsOn) continue

        const dependsOnValues = Object.values(element.input.context.dependsOn)
        if (dependsOnValues.includes(sourceArgument)) {
          dependents.push(element.input.argument)
        }
      }
    }
  }
  return dependents
}

function updateConfig(argumentName, value) {
  emit('update:config', { key: argumentName, value })

  // Cascade clear dependents when source value changes
  if (props.config[argumentName] !== value) {
    const dependents = getDependentArguments(argumentName)
    for (const depArg of dependents) {
      emit('update:config', { key: depArg, value: null })
    }
  }
}
</script>
