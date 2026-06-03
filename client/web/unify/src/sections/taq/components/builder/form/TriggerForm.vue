<template>
  <DynamicForm
    :processed-segments="processedSegments"
    @update:value="handleValueUpdate"
  />
</template>

<script setup>
import { computed } from 'vue'
import { useSegmentForm } from '@/sections/taq/composables/useSegmentForm'
import DynamicForm from './DynamicForm.vue'

const props = defineProps({
  triggerDef: {
    type: Object,
    required: true,
  },
  constraints: {
    type: Array,
    default: () => [],
  },
})

const emit = defineEmits(['update:constraints'])

// Derive parameters from catalog constraint definitions
const parameters = computed(() =>
  (props.triggerDef.constraints || []).map((c) => ({
    argumentName: c.name,
    types: c.types || [],
    required: false,
    aggregate: false,
  })),
)

// Find catalog constraint definition by name
function getCatalogConstraint(name) {
  return (props.triggerDef.constraints || []).find((c) => c.name === name)
}

// Read value from constraint: find by name, extract @value from first entry
function getValue(argumentName) {
  const constraint = props.constraints.find((c) => c.name === argumentName)
  if (!constraint?.values?.length) return null
  return constraint.values[0]['@value'] ?? null
}

// Not used for triggers currently
function getAggregateValue() {
  return null
}

function onUpdate(argumentName, value) {
  let newConstraints = [...props.constraints]

  if (value === null || value === undefined || value === '') {
    // Remove constraint when value is cleared
    newConstraints = newConstraints.filter((c) => c.name !== argumentName)
  } else {
    // Find the catalog constraint to get the @type
    const catalogDef = getCatalogConstraint(argumentName)
    const type = catalogDef?.types?.[0] || 'String'

    const constraintValue = {
      '@type': type,
      '@value': String(value),
    }

    const idx = newConstraints.findIndex((c) => c.name === argumentName)
    if (idx !== -1) {
      newConstraints[idx] = {
        ...newConstraints[idx],
        values: [constraintValue],
      }
    } else {
      newConstraints.push({
        name: argumentName,
        op: '=',
        values: [constraintValue],
      })
    }
  }

  emit('update:constraints', newConstraints)
}

const { processedSegments, updateValue } = useSegmentForm({
  segments: () => props.triggerDef.segments || [],
  parameters: () => parameters.value,
  getValue,
  getAggregateValue,
  onUpdate,
})

function handleValueUpdate(argumentName, value) {
  updateValue(argumentName, value)
}
</script>
