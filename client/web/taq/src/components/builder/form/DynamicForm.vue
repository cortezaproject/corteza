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
            v-bind="input.contextProps"
            :model-value="input.value"
            @update:model-value="updateArgument(input.argument, $event)"
          />
        </div>
      </div>
    </template>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import DynamicInput from './DynamicInput.vue'

const { t } = useI18n()

const props = defineProps({
  functionDef: {
    type: Object,
    required: true,
  },
  arguments: {
    type: Array,
    default: () => [],
  },
})

const emit = defineEmits(['update:arguments'])

const parameters = computed(() => props.functionDef.parameters || [])

function getParam(argumentName) {
  return parameters.value.find((p) => p.argumentName === argumentName)
}

function isRequired(argumentName) {
  return getParam(argumentName)?.required || false
}

// Get value for a non-aggregate argument
function getArgValue(argumentName) {
  const expr = props.arguments.find((a) => a.argumentName === argumentName)
  return expr?.value ?? null
}

// Get aggregate values as { [target]: value } object (for FieldValueMap)
function getAggregateValue(argumentName) {
  const exprs = props.arguments.filter((a) => a.argumentName === argumentName)
  if (exprs.length === 0) return null
  const result = {}
  for (const expr of exprs) {
    if (expr.target) {
      result[expr.target] = expr.value ?? expr.expr ?? ''
    }
  }
  return Object.keys(result).length > 0 ? result : null
}

function resolveContextProps(context) {
  if (!context?.dependsOn) return {}

  const resolved = {}
  for (const [propName, sourceArgument] of Object.entries(context.dependsOn)) {
    resolved[propName] = getArgValue(sourceArgument)
  }

  return resolved
}

// Look up label for an argument name from segment elements
function getLabelForArgument(argumentName) {
  for (const segment of props.functionDef.segments || []) {
    for (const section of segment.sections || []) {
      for (const element of section.elements || []) {
        if (element.input?.argument === argumentName) {
          return element.input.label || argumentName
        }
      }
    }
  }
  return argumentName
}

// Check if an input should be disabled based on unresolved dependencies
function resolveDisabledState(context) {
  if (!context?.dependsOn) return { disabled: false, disabledPlaceholder: '' }

  const missingDeps = []
  for (const [, sourceArgument] of Object.entries(context.dependsOn)) {
    if (!getArgValue(sourceArgument)) {
      missingDeps.push(getLabelForArgument(sourceArgument).toLowerCase())
    }
  }

  if (missingDeps.length > 0) {
    return {
      disabled: true,
      disabledPlaceholder: t('builder.form.selectFirst', { field: missingDeps.join(', ') }),
    }
  }

  return { disabled: false, disabledPlaceholder: '' }
}

const processedSegments = computed(() => {
  return (props.functionDef.segments || []).map((segment, sIdx) => ({
    key: `segment-${sIdx}`,
    title: segment.meta?.short || null,
    sections: (segment.sections || []).map((section, secIdx) => ({
      key: `section-${sIdx}-${secIdx}`,
      title: section.meta?.short || null,
      inputs: (section.elements || [])
        .filter((el) => el.input)
        .map((element, elIdx) => {
          const param = getParam(element.input.argument)
          const isAgg = param?.aggregate || false
          const { disabled, disabledPlaceholder } = resolveDisabledState(element.input.context)
          return {
            key: `input-${sIdx}-${secIdx}-${elIdx}`,
            type: element.input.type,
            label: element.input.label,
            placeholder: element.input.placeholder,
            disabledPlaceholder,
            disabled,
            argument: element.input.argument,
            required: isRequired(element.input.argument),
            contextProps: resolveContextProps(element.input.context),
            value: isAgg
              ? getAggregateValue(element.input.argument)
              : getArgValue(element.input.argument),
          }
        }),
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

function updateArgument(argumentName, value) {
  const param = getParam(argumentName)
  const isAgg = param?.aggregate || false

  let newArgs = [...props.arguments]

  if (isAgg && typeof value === 'object' && value !== null) {
    // Remove all existing entries for this aggregate argument
    newArgs = newArgs.filter((a) => a.argumentName !== argumentName)
    // Add one Expr per target/value pair
    for (const [target, val] of Object.entries(value)) {
      newArgs.push({
        argumentName,
        target,
        type: param?.types?.[0] || 'Any',
        value: val,
      })
    }
  } else {
    // Non-aggregate: find existing or create
    const idx = newArgs.findIndex((a) => a.argumentName === argumentName)
    const expr = {
      argumentName,
      type: param?.types?.[0] || 'Any',
      value,
    }
    if (idx !== -1) {
      newArgs[idx] = { ...newArgs[idx], ...expr }
    } else {
      newArgs.push(expr)
    }
  }

  // Cascade clear dependents when source value changes
  const oldValue = getArgValue(argumentName)
  if (oldValue !== value) {
    const dependents = getDependentArguments(argumentName)
    for (const depArg of dependents) {
      const depParam = getParam(depArg)
      if (depParam?.aggregate) {
        // Remove all aggregate entries
        newArgs = newArgs.filter((a) => a.argumentName !== depArg)
      } else {
        const depIdx = newArgs.findIndex((a) => a.argumentName === depArg)
        if (depIdx !== -1) {
          newArgs[depIdx] = { ...newArgs[depIdx], value: null }
        }
      }
    }
  }

  emit('update:arguments', newArgs)
}
</script>
