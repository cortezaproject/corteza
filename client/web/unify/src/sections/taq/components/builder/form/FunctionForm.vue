<template>
  <DynamicForm
    :processed-segments="processedSegments"
    :show-reference-toggle="true"
    @update:value="handleValueUpdate"
    @toggle-reference="onToggleReference"
    @clear-reference="handleClearReference"
    @update-reference-source="handleReferenceSourceUpdate"
  />
</template>

<script setup>
import { computed } from 'vue'
import { useSegmentForm } from '@/sections/taq/composables/useSegmentForm'
import DynamicForm from './DynamicForm.vue'

const props = defineProps({
  functionDef: {
    type: Object,
    required: true,
  },
  arguments: {
    type: Array,
    default: () => [],
  },
  upstreamResults: {
    type: Array,
    default: () => [],
  },
  nodes: {
    type: Array,
    default: () => [],
  },
})

const emit = defineEmits(['update:arguments', 'toggleReference'])

const parameters = computed(() => props.functionDef.parameters || [])

function getParam(argumentName) {
  return parameters.value.find(p => p.argumentName === argumentName)
}

// Get value for a non-aggregate argument
function getValue(argumentName) {
  const expr = props.arguments.find(a => a.argumentName === argumentName)
  // If this arg has a reference, don't return the value (it's in reference mode)
  if (expr?.scope && (expr?.source || expr?.expr)) return null
  return expr?.value ?? null
}

// Get aggregate values as { [target]: { value, scope?, source? } } object (for FieldValueMap)
function getAggregateValue(argumentName) {
  const exprs = props.arguments.filter(a => a.argumentName === argumentName)
  if (exprs.length === 0) return null
  const result = {}
  for (const expr of exprs) {
    if (expr.target) {
      if (expr.scope && expr.expr) {
        // Reference row
        result[expr.target] = { value: '', scope: expr.scope, source: expr.expr }
      } else {
        // Literal value row
        result[expr.target] = { value: expr.value ?? expr.expr ?? '' }
      }
    }
  }
  return Object.keys(result).length > 0 ? result : null
}

// Get reference info for an argument (scope + source)
function getReferenceInfo(argumentName) {
  const expr = props.arguments.find(a => a.argumentName === argumentName)
  if (expr?.scope && (expr?.source || expr?.expr)) {
    return { scope: expr.scope, source: expr.source || expr.expr }
  }
  return null
}

/**
 * Resolve a reference value at design time.
 * For trigger references: look up the trigger node's constraints for the matching property.
 */
function resolveReferenceValue(scope, source) {
  // Find the node whose handle (data.ref) matches the scope
  const node = props.nodes.find(n => n.data?.ref === scope)
  if (!node) return null

  // For trigger nodes, look up constraints
  if (node.type === 'trigger') {
    const constraints = node.data?.constraints || []
    const constraint = constraints.find(c => c.name === source)
    if (constraint?.values?.length) {
      return constraint.values[0]['@value'] ?? null
    }
  }

  return null
}

function onUpdate(argumentName, value) {
  const param = getParam(argumentName)
  const isTypeAggregate = param?.types && (param.types.includes('FieldValueMap') || param.types.includes('Array'))
  const isAgg = param?.aggregate || isTypeAggregate || false

  let newArgs = [...props.arguments]

  if (value === null || value === undefined) {
    // Cascade clear: null out or remove the argument
    if (isAgg) {
      newArgs = newArgs.filter(a => a.argumentName !== argumentName)
    } else {
      const idx = newArgs.findIndex(a => a.argumentName === argumentName)
      if (idx !== -1) {
        newArgs[idx] = {
          ...newArgs[idx],
          value: null,
          scope: undefined,
          source: undefined,
          expr: undefined,
        }
      }
    }
  } else if (isAgg && typeof value === 'object') {
    // Remove all existing entries for this aggregate argument
    newArgs = newArgs.filter(a => a.argumentName !== argumentName)
    // Add one Expr per target/value pair (supports both literal and reference rows)
    for (const [target, rowData] of Object.entries(value)) {
      if (rowData && typeof rowData === 'object' && 'scope' in rowData && rowData.scope) {
        // Reference row
        newArgs.push({
          argumentName,
          target,
          type: 'String',
          scope: rowData.scope,
          expr: rowData.source,
          value: undefined,
        })
      } else {
        // Literal value row
        const val = rowData && typeof rowData === 'object' ? rowData.value : rowData
        newArgs.push({
          argumentName,
          target,
          type: 'String',
          value: val,
        })
      }
    }
  } else {
    // Non-aggregate: find existing or create — clear any reference fields
    const idx = newArgs.findIndex(a => a.argumentName === argumentName)
    const expr = {
      argumentName,
      type: param?.types?.[0] || 'Any',
      value,
      scope: undefined,
      source: undefined,
      expr: undefined,
    }
    if (idx !== -1) {
      newArgs[idx] = { ...newArgs[idx], ...expr }
    } else {
      newArgs.push(expr)
    }
  }

  emit('update:arguments', newArgs)
}

// Handle setting a reference on an argument
// When `target` is provided, it targets a specific row in an aggregate argument (e.g. FieldValueMap)
function onReferenceSelect(argumentName, { scope, source }, target) {
  const param = getParam(argumentName)
  let newArgs = [...props.arguments]

  if (target) {
    // Per-row reference for aggregate arguments
    // Find existing entry for this target, or add new one
    const idx = newArgs.findIndex(a => a.argumentName === argumentName && a.target === target)
    const entry = {
      argumentName,
      target,
      type: 'String',
      scope,
      expr: source,
      value: undefined,
    }

    if (idx !== -1) {
      newArgs[idx] = { ...newArgs[idx], ...entry }
    } else {
      newArgs.push(entry)
    }
  } else {
    // Whole-argument reference (non-aggregate)
    const idx = newArgs.findIndex(a => a.argumentName === argumentName)
    const expr = {
      argumentName,
      type: param?.types?.[0] || 'Any',
      scope,
      expr: source,
      value: undefined,
    }

    if (idx !== -1) {
      newArgs[idx] = { ...newArgs[idx], ...expr }
    } else {
      newArgs.push(expr)
    }
  }

  emit('update:arguments', newArgs)
}

// Handle clearing a reference
function handleClearReference(argumentName) {
  let newArgs = [...props.arguments]
  const idx = newArgs.findIndex(a => a.argumentName === argumentName)
  if (idx !== -1) {
    newArgs[idx] = {
      ...newArgs[idx],
      scope: undefined,
      source: undefined,
      expr: undefined,
      value: null,
    }
  }
  emit('update:arguments', newArgs)
}

// Handle editing a reference source (user typed in the reference chip)
function handleReferenceSourceUpdate(argumentName, newSource) {
  let newArgs = [...props.arguments]
  const idx = newArgs.findIndex(a => a.argumentName === argumentName)
  if (idx !== -1 && newArgs[idx].scope) {
    newArgs[idx] = { ...newArgs[idx], expr: newSource, source: undefined }
    emit('update:arguments', newArgs)
  }
}

// Expose onReferenceSelect for parent components
defineExpose({ onReferenceSelect })

const { processedSegments, updateValue } = useSegmentForm({
  segments: () => props.functionDef.segments || [],
  parameters: () => parameters.value,
  getValue,
  getAggregateValue,
  onUpdate,
  getReferenceInfo,
  resolveReferenceValue,
  upstreamResults: () => props.upstreamResults,
})

// Automatically seed defaults into arguments on mount if they are missing
import { onMounted } from 'vue'

onMounted(() => {
  let changed = false
  const newArgs = [...props.arguments]

  processedSegments.value.forEach(segment => {
    segment.sections?.forEach(section => {
      section.inputs?.forEach(input => {
        // If argument is missing entirely and it has a default, add it
        const exists = newArgs.some(a => a.argumentName === input.argument)
        if (!exists && input.defaultValue !== undefined) {
          const param = getParam(input.argument)
          newArgs.push({
            argumentName: input.argument,
            type: param?.types?.[0] || 'Any',
            value: input.defaultValue,
            scope: undefined,
            source: undefined,
          })
          changed = true
        }
      })
    })
  })

  if (changed) {
    emit('update:arguments', newArgs)
  }
})

function handleValueUpdate(argumentName, value) {
  updateValue(argumentName, value)
}

// Enrich toggleReference with the parameter's accepted types
// Handles both plain string (argument name) and { argument, target } (per-row FieldValueMap/Array)
function onToggleReference(payload) {
  const argumentName = typeof payload === 'string' ? payload : payload.argument
  const target = typeof payload === 'object' ? payload.target : undefined
  const param = getParam(argumentName)

  // If selecting a reference for an entire array argument, use param.types.
  // If selecting a reference for a single ELEMENT in the array (target is defined), it can be Any (or String).
  const expectedTypes = target ? ['Any'] : (param?.types || ['Any'])

  emit('toggleReference', {
    name: argumentName,
    types: expectedTypes,
    target,
  })
}
</script>
