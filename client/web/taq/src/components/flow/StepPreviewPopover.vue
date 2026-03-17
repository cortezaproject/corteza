<template>
  <Teleport to=".vue-flow__nodes">
  <Transition name="popover-fade">
  <div
    v-if="visible && hasContent"
    ref="popoverEl"
    class="absolute z-[9999] bg-surface border-[0.5px] border-surface rounded-xl shadow-lg p-3 w-[260px] max-h-[280px] overflow-auto transition-opacity duration-150"
    :class="[
      !alwaysShow ? 'pointer-events-none' : '',
      hiddenByHover ? 'opacity-0 pointer-events-none' : '',
    ]"
    :style="positionStyle"
    @mouseenter="onPreviewEnter"
  >

    <!-- Branch configuration preview -->
    <div v-if="node?.type === 'branch'" class="flex flex-col gap-1.5">
      <!-- Gateway type -->
      <div class="flex items-baseline gap-2 text-xs">
        <span class="text-primary shrink-0">{{ $t('builder.branch.gatewayType') }}:</span>
        <span class="text-color-emphasis">
          {{ gatewayLabel }}
        </span>
      </div>
      <!-- Branch paths -->
      <div v-if="branchPaths.length" class="w-full rounded-lg border border-surface overflow-hidden">
        <div class="flex text-xs font-medium text-muted-color bg-emphasis">
          <span class="px-2 py-1 flex-1">{{ $t('builder.configSidebar.branches') }}</span>
          <span class="px-2 py-1 flex-1">{{ $t('builder.branch.condition') }}</span>
        </div>
        <div
          v-for="path in branchPaths"
          :key="path.edgeId"
          class="flex text-xs border-t border-surface"
        >
          <span class="px-2 py-1 text-muted-color truncate flex-1" :title="path.label">
            {{ path.label }}
          </span>
          <span
            class="px-2 py-1 truncate flex-1"
            :class="path.expr ? 'font-mono text-color-emphasis' : 'italic text-muted-color'"
            :title="path.expr || ''"
          >
            {{ path.expr || path.defaultLabel }}
          </span>
        </div>
      </div>
    </div>

    <!-- Configuration preview list (steps / triggers) -->
    <div v-else-if="previewItems.length" class="flex flex-col gap-2">
      <div
        v-for="item in previewItems"
        :key="item.key"
        class="flex flex-col gap-0.5 text-xs"
      >
        <span class="text-primary">{{ item.label }}</span>

        <!-- Reference value -->
        <CViewReference
          v-if="item.isReference"
          :scope="item.refScope"
          :source="item.refSource"
          :nodes="nodes"
        />

        <!-- Type-resolved value -->
        <component
          v-else
          :is="item.viewerComponent"
          :model-value="item.value"
          v-bind="item.viewerProps"
        />
      </div>
    </div>

    <!-- No config: hide preview entirely instead of showing a message -->
  </div>
  </Transition>
  </Teleport>
</template>

<script setup>
import { computed, nextTick, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import CViewReference from '@/components/builder/form/viewers/CViewReference.vue'
import { resolveViewerComponent } from '@/components/builder/form/viewers/registry'
import { NODE_DIMENSIONS } from '@/utils/flow-constants'

const { t } = useI18n()

const props = defineProps({
  visible: { type: Boolean, default: false },
  node: { type: Object, default: null },
  functions: { type: Array, default: () => [] },
  triggers: { type: Array, default: () => [] },
  nodes: { type: Array, default: () => [] },
  edges: { type: Array, default: () => [] },
  alwaysShow: { type: Boolean, default: false },
  nodeEl: { type: Object, default: null },
})

const popoverEl = ref(null)
const isFlippedUp = ref(false)
const hiddenByHover = ref(false)

/**
 * When user hovers the preview in always-show mode,
 * hide it and track mouse to restore when mouse leaves the area.
 */
function onPreviewEnter() {
  if (!props.alwaysShow) return
  const el = popoverEl.value
  if (!el) return

  const rect = el.getBoundingClientRect()
  const pad = 8
  hiddenByHover.value = true

  function onMove(e) {
    if (
      e.clientX < rect.left - pad ||
      e.clientX > rect.right + pad ||
      e.clientY < rect.top - pad ||
      e.clientY > rect.bottom + pad
    ) {
      hiddenByHover.value = false
      document.removeEventListener('mousemove', onMove)
    }
  }

  document.addEventListener('mousemove', onMove)
}

const positionStyle = ref({})

// Calculate position using VueFlow coordinates (node.position)
watch(() => props.visible, async (isVisible) => {
  if (!isVisible) return

  isFlippedUp.value = false
  hiddenByHover.value = false

  if (!props.node?.position) return

  const nodeEl = props.nodeEl?.$el || props.nodeEl
  const nodeW = NODE_DIMENSIONS.WIDTH
  const nodeH = nodeEl ? nodeEl.offsetHeight : NODE_DIMENSIONS.HEIGHT
  const popoverW = 260
  const gap = 6

  // Center below node in VueFlow coordinate space
  const x = props.node.position.x + (nodeW / 2) - (popoverW / 2)
  const y = props.node.position.y + nodeH + gap

  positionStyle.value = {
    top: `${y}px`,
    left: `${x}px`,
  }

  // After render, check viewport clipping and flip if needed
  await nextTick()
  await nextTick()
  const el = popoverEl.value
  if (!el) return

  const rect = el.getBoundingClientRect()
  const vh = window.innerHeight
  const vw = window.innerWidth

  // If clipped on bottom, flip above
  if (rect.bottom > vh - 8) {
    const popoverH = el.offsetHeight
    const yAbove = props.node.position.y - popoverH - gap
    positionStyle.value = {
      top: `${yAbove}px`,
      left: `${x}px`,
    }
    isFlippedUp.value = true
  }

  // If clipped on left/right, re-check and adjust
  const rect2 = el.getBoundingClientRect()
  if (rect2.left < 8 || rect2.right > vw - 8) {
    // Re-center on the node using screen coords
    const nodeEl = props.nodeEl?.$el || props.nodeEl
    if (nodeEl) {
      const nodeRect = nodeEl.getBoundingClientRect()
      const shift = (nodeRect.left + nodeRect.width / 2 - rect2.width / 2) - rect2.left
      positionStyle.value = {
        ...positionStyle.value,
        left: `${parseFloat(positionStyle.value.left) + shift}px`,
      }
    }
  }
})


/**
 * Whether the preview has any content worth displaying
 */
const hasContent = computed(() => {
  if (!props.node) return false
  if (props.node.type === 'branch') return branchPaths.value.length > 0
  return previewItems.value.length > 0
})

/**
 * Branch gateway type label
 */
const gatewayLabel = computed(() => {
  const raw = props.node?.data?.nodeType
  if (raw === 'incl') return t('builder.branch.inclusive')
  return t('builder.branch.exclusive')
})

/**
 * Branch paths from outgoing edges
 */
const branchPaths = computed(() => {
  if (!props.node || props.node.type !== 'branch') return []

  const outgoing = props.edges
    .filter(e => e.source === props.node.id)
    .map((edge, index, arr) => {
      const isFirst = index === 0
      const isLast = index === arr.length - 1

      let label
      if (isFirst) label = t('builder.branch.if')
      else if (isLast) label = t('builder.branch.else')
      else label = t('builder.branch.elseIf')

      return {
        edgeId: edge.id,
        label,
        expr: edge.data?.expr || '',
        defaultLabel: isLast ? t('builder.branch.defaultPath') : t('builder.preview.notSet'),
      }
    })

  return outgoing
})

/**
 * Look up function/trigger definition from catalog
 */
const definition = computed(() => {
  if (!props.node) return null
  const nodeType = props.node.data?.nodeType

  if (props.node.type === 'trigger') {
    const resourceType = props.node.data?.resourceType
    return (
      props.triggers.find(
        tr => tr.eventType === nodeType && (!resourceType || tr.resourceType === resourceType),
      ) || null
    )
  }

  return props.functions.find(f => f.ref === nodeType) || null
})

/**
 * Get value for a non-aggregate argument (mirrors FunctionForm.getValue)
 */
function getArgValue(argumentName) {
  const args = props.node?.data?.arguments || []
  const expr = args.find(a => a.argumentName === argumentName)
  if (!expr) return null
  // If it has a reference, don't return the value
  if (expr.scope && expr.source) return null
  return expr.value ?? null
}

/**
 * Get reference info for an argument (mirrors FunctionForm.getReferenceInfo)
 */
function getReferenceInfo(argumentName) {
  const args = props.node?.data?.arguments || []
  const expr = args.find(a => a.argumentName === argumentName)
  if (expr?.scope && expr?.source) {
    return { scope: expr.scope, source: expr.source }
  }
  return null
}

/**
 * Get effective value: actual value or resolved reference value
 * Used for resolving context props (e.g. namespaceID for module)
 */
function getEffectiveValue(argumentName) {
  const value = getArgValue(argumentName)
  if (value != null) return value

  // If reference, resolve its value at design time
  const ref = getReferenceInfo(argumentName)
  if (ref) {
    const node = props.nodes.find(n => n.data?.ref === ref.scope)
    if (node?.type === 'trigger') {
      const constraint = (node.data?.constraints || []).find(c => c.name === ref.source)
      if (constraint?.values?.length) {
        return constraint.values[0]['@value'] ?? null
      }
    }
  }
  return null
}

/**
 * Resolve context props for dependent fields (mirrors useSegmentForm.resolveContextProps)
 * E.g. for Module input, resolves { namespaceID: <value of namespace argument> }
 */
function resolveContextProps(context) {
  if (!context?.dependsOn) return {}
  const resolved = {}
  for (const [propName, sourceArgument] of Object.entries(context.dependsOn)) {
    resolved[propName] = getEffectiveValue(sourceArgument)
  }
  return resolved
}

/**
 * Build the preview items from segments + arguments/constraints
 */
const previewItems = computed(() => {
  if (!props.node || !definition.value) return []

  const isTrigger = props.node.type === 'trigger'
  const segments = definition.value.segments || []
  const items = []

  for (const segment of segments) {
    for (const section of segment.sections || []) {
      for (const element of section.elements || []) {
        const input = element.input
        if (!input?.argument) continue

        const label = input.label || input.argument
        const inputType = input.type || 'Text'
        const inputOptions = input.options || []
        const hasOptions = inputOptions.length > 0
        const effectiveType = hasOptions ? 'Select' : inputType

        if (isTrigger) {
          // Triggers store constraints
          const constraints = props.node.data?.constraints || []
          const constraint = constraints.find(c => c.name === input.argument)
          const rawValue = constraint?.values?.length
            ? constraint.values.map(v => v['@value'] ?? v).join(', ')
            : null

          items.push({
            key: `${input.argument}`,
            label,
            value: rawValue,
            isReference: false,
            viewerComponent: resolveViewerComponent(effectiveType),
            viewerProps: hasOptions
              ? { options: inputOptions }
              : resolveContextProps(input.context),
          })
        } else {
          // Steps store arguments — mirrors FunctionForm logic
          const params = definition.value.parameters || []
          const param = params.find(p => p.argumentName === input.argument)
          const isAggregate = param?.aggregate || false

          if (isAggregate) {
            // Aggregate type (e.g. FieldValueMap) — collect all entries
            const allArgs = props.node.data?.arguments || []
            const entries = allArgs
              .filter(a => a.argumentName === input.argument && a.target)
              .map(a => ({
                target: a.target,
                value: a.value ?? a.expr ?? '',
                scope: a.scope || null,
                source: a.source || null,
                expr: a.expr || null,
              }))

            items.push({
              key: `${input.argument}`,
              label,
              value: entries,
              isReference: false,
              isAggregate: true,
              viewerComponent: resolveViewerComponent(inputType),
              viewerProps: {
                ...resolveContextProps(input.context),
                nodes: props.nodes,
              },
            })
          } else {
            const refInfo = getReferenceInfo(input.argument)

            if (refInfo) {
              // Reference mode
              items.push({
                key: `${input.argument}`,
                label,
                isReference: true,
                refScope: refInfo.scope,
                refSource: refInfo.source,
              })
            } else {
              // Literal value
              const value = getArgValue(input.argument)
              items.push({
                key: `${input.argument}`,
                label,
                value,
                isReference: false,
                viewerComponent: resolveViewerComponent(effectiveType),
                viewerProps: hasOptions
                  ? { options: inputOptions }
                  : resolveContextProps(input.context),
              })
            }
          }
        }
      }
    }
  }

  return items
})
</script>

<style scoped>
.popover-fade-enter-active {
  transition: opacity 150ms ease-out, transform 150ms ease-out;
}
.popover-fade-leave-active {
  transition: opacity 100ms ease-in, transform 100ms ease-in;
}
.popover-fade-enter-from {
  opacity: 0;
  transform: scale(0.95);
}
.popover-fade-leave-to {
  opacity: 0;
  transform: scale(0.95);
}
</style>
