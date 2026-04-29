import { computed } from 'vue'
import { useI18n } from 'vue-i18n'

/**
 * Composable that provides shared segment form logic.
 * Handles dependency resolution, cascade clearing, disabled states,
 * and builds processedSegments for the DynamicForm renderer.
 *
 * Type-specific forms (FunctionForm, TriggerForm) provide value
 * read/write callbacks and render DynamicForm with the result.
 */
export function useSegmentForm(options: {
  segments: () => any[]
  parameters: () => any[]
  getValue: (_argumentName: string) => any
  getAggregateValue: (_argumentName: string) => any
  onUpdate: (_argumentName: string, _value: any) => void
  getReferenceInfo?: (_argumentName: string) => { scope: string; source: string } | null
  resolveReferenceValue?: (_scope: string, _source: string) => any
  upstreamResults?: () => any[]
}) {
  const { t } = useI18n()

  function getParam(argumentName: string) {
    return options.parameters().find((p: any) => p.argumentName === argumentName)
  }

  function isRequired(argumentName: string) {
    return getParam(argumentName)?.required || false
  }

  /**
   * Get the effective value for an argument, falling back to
   * resolving its reference if it's in reference mode.
   */
  function getEffectiveValue(argumentName: string) {
    const value = options.getValue(argumentName)
    if (value != null) return value

    // If the arg has a reference, try to resolve it at design time
    const ref = options.getReferenceInfo?.(argumentName)
    if (ref) {
      return options.resolveReferenceValue?.(ref.scope, ref.source) ?? null
    }
    return null
  }

  function resolveContextProps(context: any) {
    if (!context?.dependsOn) return {}

    const resolved: Record<string, any> = {}
    for (const [propName, sourceArgument] of Object.entries(context.dependsOn)) {
      resolved[propName] = getEffectiveValue(sourceArgument as string)
    }
    return resolved
  }

  function getLabelForArgument(argumentName: string) {
    for (const segment of options.segments()) {
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

  function resolveDisabledState(context: any) {
    if (!context?.dependsOn) return { disabled: false, disabledPlaceholder: '' }

    const missingDeps: string[] = []
    for (const [, sourceArgument] of Object.entries(context.dependsOn)) {
      if (!getEffectiveValue(sourceArgument as string)) {
        missingDeps.push(getLabelForArgument(sourceArgument as string).toLowerCase())
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

  function getDependentArguments(sourceArgument: string) {
    const dependents: string[] = []
    for (const segment of options.segments()) {
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

  /**
   * Build a human-readable label for a reference like "HTTP Request → response"
   */
  function buildReferenceLabel(_scope: string, source: string) {
    return source
  }

  const processedSegments = computed(() => {
    return (options.segments() || []).map((segment: any, sIdx: number) => ({
      key: `segment-${sIdx}`,
      title: segment.meta?.short || null,
      sections: (segment.sections || []).map((section: any, secIdx: number) => ({
        key: `section-${sIdx}-${secIdx}`,
        title: section.meta?.short || null,
        inputs: (section.elements || [])
          .filter((el: any) => el.input)
          .map((element: any, elIdx: number) => {
            const param = getParam(element.input.argument)
            const AGGREGATE_TYPES = ['FieldValueMap', 'Array']
            const isAgg = param?.aggregate || AGGREGATE_TYPES.includes(element.input.type) || false
            const { disabled, disabledPlaceholder } = resolveDisabledState(element.input.context)

            // Reference state
            const refInfo = options.getReferenceInfo?.(element.input.argument) ?? null
            const isReference = !!refInfo
            const referenceLabel = refInfo ? buildReferenceLabel(refInfo.scope, refInfo.source) : ''

            const inputOptions = element.input.options || []

            return {
              key: `input-${sIdx}-${secIdx}-${elIdx}`,
              type: element.input.type,
              label: element.input.label,
              placeholder: element.input.placeholder || t('builder.form.selectPlaceholder', { field: element.input.label || element.input.argument }),
              disabledPlaceholder,
              disabled,
              argument: element.input.argument,
              required: element.input.required || isRequired(element.input.argument),
              contextProps: resolveContextProps(element.input.context),
              value: isAgg
                ? options.getAggregateValue(element.input.argument)
                : options.getValue(element.input.argument) ?? element.input.default,
              defaultValue: element.input.default,
              isReference,
              referenceLabel,
              options: inputOptions,
              isAggregate: isAgg,
              description: element.input.description || '',
            }
          }),
      })),
    }))
  })

  function updateValue(argumentName: string, value: any) {
    // Cascade clear dependents when source value changes
    const oldValue = options.getValue(argumentName)
    if (oldValue !== value) {
      const dependents = getDependentArguments(argumentName)
      for (const depArg of dependents) {
        options.onUpdate(depArg, null)
      }
    }

    options.onUpdate(argumentName, value)
  }

  return {
    processedSegments,
    updateValue,
  }
}
