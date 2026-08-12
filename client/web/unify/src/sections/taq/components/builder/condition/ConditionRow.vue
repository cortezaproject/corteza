<template>
  <div class="condition-row flex items-start gap-1">
    <div class="flex flex-col gap-1 flex-1 min-w-0">
      <!-- Row 1: variable -->
      <div class="flex items-center gap-1">
        <!-- Variable: reference chip when linked, else input + link button -->
        <CReferenceChip
          v-if="hasVariableRef"
          :label="variableRefLabel"
          size="small"
          class="flex-1 min-w-0"
          @click="emit('toggleReference', { side: 'variable', edgeId, rowIndex })"
          @clear="onClearVariable"
          @update:label="onVariableRefEdit"
        />
        <div v-else class="flex items-center gap-1 flex-1 min-w-0">
          <InputText
            :model-value="symbolStr"
            :placeholder="$t('builder.condition.selectValue')"
            class="flex-1 min-w-0"
            size="small"
            @focus="emit('toggleReference', { side: 'variable', edgeId, rowIndex })"
            @update:model-value="onSymbolChange"
          />
          <Button
            icon="pi pi-link"
            text
            rounded
            size="small"
            :severity="isVariableRefActive ? 'primary' : 'secondary'"
            @click="emit('toggleReference', { side: 'variable', edgeId, rowIndex })"
          />
        </div>
      </div>

      <!-- Row 2: operator -->
      <Select
        :model-value="node.ref || 'eq'"
        :options="operatorOptions"
        option-label="label"
        option-value="value"
        size="small"
        class="w-full"
        @update:model-value="onOperatorChange"
      />

      <!-- Row 3: value (hidden for null checks) -->
      <template v-if="!isNullCheck">
        <CReferenceChip
          v-if="hasValueRef"
          :label="valueRefLabel"
          size="small"
          class="min-w-0"
          @click="emit('toggleReference', { side: 'value', edgeId, rowIndex })"
          @clear="onClearValue"
          @update:label="onValueRefEdit"
        />
        <div
          v-else
          class="flex items-center gap-1"
          @focusin="emit('toggleReference', { side: 'value', edgeId, rowIndex })"
        >
          <CFieldEditor
            :field="resolvedFieldDef"
            :model-value="valueVal"
            class="flex-1 min-w-0"
            @update:model-value="onValueChange"
          />
          <Button
            icon="pi pi-link"
            text
            rounded
            size="small"
            :severity="isValueRefActive ? 'primary' : 'secondary'"
            @click="emit('toggleReference', { side: 'value', edgeId, rowIndex })"
          />
        </div>
      </template>
    </div>

    <!-- Delete (outside the card) -->
    <Button
      icon="pi pi-trash"
      text
      rounded
      size="small"
      severity="danger"
      class="shrink-0"
      @click="emit('delete')"
    />
  </div>
</template>

<script setup>
import { computed, inject, ref as vueRef, watch, watchEffect } from 'vue'
import { useI18n } from 'vue-i18n'
import CReferenceChip from '../form/CReferenceChip.vue'
import { CFieldEditor } from '@planetcrust/human-vue/src/components/field'
import { useModuleStore } from '@planetcrust/human-vue'

const { t } = useI18n()
const $ComposeAPI = inject('$ComposeAPI')
const moduleStore = useModuleStore()

const props = defineProps({
  node: { type: Object, required: true },
  edgeId: { type: String, default: '' },
  rowIndex: { type: Number, default: 0 },
})

const emit = defineEmits(['update', 'delete', 'toggleReference'])

// Injected upstream results for resolving scope handles to step names
const upstreamResults = inject('taq-upstream-results', vueRef([]))

// Injected active reference argument to show active state on link buttons
const activeReferenceArgument = inject('activeReferenceArgument', vueRef(null))

const isVariableRefActive = computed(() => {
  const arg = activeReferenceArgument.value
  return (
    arg?.isCondition &&
    arg.edgeId === props.edgeId &&
    arg.side === 'variable' &&
    arg.rowIndex === props.rowIndex
  )
})

const isValueRefActive = computed(() => {
  const arg = activeReferenceArgument.value
  return (
    arg?.isCondition &&
    arg.edgeId === props.edgeId &&
    arg.side === 'value' &&
    arg.rowIndex === props.rowIndex
  )
})

function resolveLabel(scope, symbol) {
  return symbol
}

// Operator definitions — must match server/pkg/ast/eval.go refs
const OPERATORS = ['eq', 'ne', 'lt', 'lte', 'gt', 'gte', 'isNull', 'isNotNull']

const operatorOptions = computed(() =>
  OPERATORS.map(op => ({
    value: op,
    label: t(`builder.condition.operators.${op}`),
  })),
)

const isNullCheck = computed(() => {
  return props.node.ref === 'isNull' || props.node.ref === 'isNotNull'
})

// --- Variable (left side) ---

const leftArg = computed(() => props.node.args?.[0])

const hasVariableRef = computed(() => {
  const left = leftArg.value
  return left?.symbol && left?.meta?.scope
})

const variableRefLabel = computed(() => {
  const left = leftArg.value
  if (!left) return ''
  const scope = left.meta?.scope || ''
  const symbol = left.symbol || ''
  return scope ? resolveLabel(scope, symbol) : symbol
})

const symbolStr = computed(() => leftArg.value?.symbol || '')

function makeFieldDef(kind = 'String', options = {}) {
  return { kind, name: 'value', label: '', options, isMulti: false, isRequired: false }
}

const resolvedFieldDef = vueRef(makeFieldDef())
const userChangedSymbol = vueRef(false)

// Derived for use in onValueChange typing logic
const resolvedVariableType = computed(() => resolvedFieldDef.value.kind)

watch(resolvedVariableType, (newKind, oldKind) => {
  if (oldKind && newKind !== oldKind && !hasValueRef.value && userChangedSymbol.value) {
    onClearValue()
  }
  userChangedSymbol.value = false
})

watchEffect(async () => {
  const left = leftArg.value
  if (!left?.symbol || !left?.meta?.scope) {
    resolvedFieldDef.value = makeFieldDef()
    return
  }
  const scope = left.meta.scope
  const symbol = left.symbol

  const step = upstreamResults.value.find(s => s.handle === scope)
  if (!step) {
    resolvedFieldDef.value = makeFieldDef()
    return
  }

  const topLevelName = symbol.split('.')[0]
  const topLevel = (step.properties || step.results || []).find(r => r.sourceName === topLevelName)
  if (!topLevel) {
    resolvedFieldDef.value = makeFieldDef()
    return
  }

  if (symbol === topLevelName) {
    resolvedFieldDef.value = makeFieldDef(topLevel.types?.[0] || 'String')
    return
  }

  if (topLevel.types?.includes('ComposeRecord') && symbol.startsWith(`${topLevelName}.`)) {
    const fieldPath = symbol.substring(topLevelName.length + 1)

    if (fieldPath === 'recordID' || fieldPath === 'moduleID' || fieldPath === 'namespaceID') {
      resolvedFieldDef.value = makeFieldDef('ID')
      return
    }
    if (
      fieldPath === 'ownedBy' ||
      fieldPath === 'createdBy' ||
      fieldPath === 'updatedBy' ||
      fieldPath === 'deletedBy'
    ) {
      resolvedFieldDef.value = makeFieldDef('UserSelector')
      return
    }
    if (fieldPath === 'createdAt' || fieldPath === 'updatedAt' || fieldPath === 'deletedAt') {
      resolvedFieldDef.value = makeFieldDef('DateTime')
      return
    }
    if (fieldPath.startsWith('values.')) {
      const customFieldName = fieldPath.substring(7)
      if (topLevel.namespaceID && topLevel.moduleID) {
        try {
          const mod = await moduleStore.findByID({
            namespaceID: topLevel.namespaceID,
            moduleID: topLevel.moduleID,
          })
          const field = mod.fields?.find(f => f.name === customFieldName)
          if (field) {
            resolvedFieldDef.value = {
              kind: field.kind,
              name: field.name,
              label: field.label || field.name,
              options: field.options || {},
              isMulti: false,
              isRequired: false,
            }
            return
          }
        } catch {
          // fall through to default field def
        }
      }
    }
  }

  resolvedFieldDef.value = makeFieldDef()
})

// --- Value (right side) ---

const rightArg = computed(() => props.node.args?.[1])

const hasValueRef = computed(() => {
  const right = rightArg.value
  return right?.symbol && right?.meta?.scope
})

const valueRefLabel = computed(() => {
  const right = rightArg.value
  if (!right) return ''
  const scope = right.meta?.scope || ''
  const symbol = right.symbol || ''
  return scope ? resolveLabel(scope, symbol) : symbol
})

const valueVal = computed(() => {
  const right = rightArg.value
  if (!right) return ''
  if (right.symbol) return right.symbol
  if (right.value && typeof right.value === 'object') {
    return right.value['@value'] ?? ''
  }
  return ''
})

// --- Deep clone helper ---

function cloneNode(node) {
  return JSON.parse(JSON.stringify(node))
}

// --- Event handlers ---

function onSymbolChange(val) {
  userChangedSymbol.value = true
  const newNode = cloneNode(props.node)
  if (!newNode.args) newNode.args = [null, null]
  newNode.args[0] = { symbol: val, meta: {} }
  if (!newNode.args[1] && !isNullCheck.value) {
    newNode.args[1] = { value: { '@type': 'String', '@value': '' } }
  }
  emit('update', newNode)
}

function onClearVariable() {
  const newNode = cloneNode(props.node)
  if (!newNode.args) newNode.args = [null, null]
  newNode.args[0] = { symbol: '', meta: {} }
  emit('update', newNode)
}

function onVariableRefEdit(val) {
  userChangedSymbol.value = true
  const newNode = cloneNode(props.node)
  if (!newNode.args) newNode.args = [null, null]
  const scope = newNode.args[0]?.meta?.scope || ''
  newNode.args[0] = { symbol: val, meta: { scope } }
  emit('update', newNode)
}

function onClearValue() {
  const newNode = cloneNode(props.node)
  if (!newNode.args) newNode.args = [{ symbol: '', meta: {} }]
  newNode.args[1] = { value: { '@type': 'String', '@value': '' } }
  emit('update', newNode)
}

function onValueRefEdit(val) {
  const newNode = cloneNode(props.node)
  if (!newNode.args) newNode.args = [null, null]
  const scope = newNode.args[1]?.meta?.scope || ''
  newNode.args[1] = { symbol: val, meta: { scope } }
  emit('update', newNode)
}

function onOperatorChange(op) {
  const newNode = cloneNode(props.node)
  newNode.ref = op

  if (op === 'isNull' || op === 'isNotNull') {
    newNode.args = (newNode.args || []).slice(0, 1)
  } else if (!newNode.args || newNode.args.length < 2) {
    if (!newNode.args) newNode.args = [{ symbol: '', meta: {} }]
    newNode.args[1] = { value: { '@type': 'String', '@value': '' } }
  }

  emit('update', newNode)
}

function onValueChange(val) {
  const newNode = cloneNode(props.node)
  if (!newNode.args) newNode.args = [{ symbol: '', meta: {} }]

  let typedVal

  if (resolvedVariableType.value === 'Boolean' || resolvedVariableType.value === 'Bool') {
    typedVal = { '@type': 'Boolean', '@value': Boolean(val) }
  } else if (
    resolvedVariableType.value === 'DateTime' ||
    resolvedVariableType.value === 'UserSelector' ||
    resolvedVariableType.value === 'ID' ||
    resolvedVariableType.value === 'String'
  ) {
    typedVal = { '@type': 'String', '@value': String(val || '') }
  } else {
    // Detect type: boolean > number > string
    const lower = String(val).toLowerCase().trim()
    if (lower === 'true' || lower === 'false') {
      typedVal = { '@type': 'Boolean', '@value': lower === 'true' }
    } else if (val !== '' && !isNaN(Number(val)) && String(val).trim() !== '') {
      if (Number.isInteger(Number(val))) {
        typedVal = { '@type': 'Integer', '@value': Number(val) }
      } else {
        typedVal = { '@type': 'Float', '@value': Number(val) }
      }
    } else {
      typedVal = { '@type': 'String', '@value': val }
    }
  }

  newNode.args[1] = { value: typedVal }
  emit('update', newNode)
}

/**
 * Called from parent when a reference is selected from ReferencePanel.
 * Sets the symbol + scope on the appropriate side.
 */
function applyReference(side, scope, source) {
  const newNode = cloneNode(props.node)
  if (!newNode.args) newNode.args = [null, null]

  const refNode = { symbol: source, meta: { scope } }
  if (side === 'variable') {
    newNode.args[0] = refNode
    if (!newNode.args[1] && !isNullCheck.value) {
      newNode.args[1] = { value: { '@type': 'String', '@value': '' } }
    }
  } else {
    newNode.args[1] = refNode
  }

  emit('update', newNode)
}

defineExpose({ applyReference })
</script>
