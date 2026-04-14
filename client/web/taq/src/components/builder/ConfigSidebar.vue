<template>
  <div class="flex flex-col h-full flex-1 overflow-auto">
    <!-- Header -->
    <div class="flex items-start justify-between py-2 pr-2 group/header">
      <div class="flex items-center gap-3 flex-1 min-w-0 pr-2">
        <div class="w-10 h-10 rounded-border flex items-center justify-center shrink-0">
          <TaqIcon
            :icon="node.data?.icon"
            :fallback="node.type === 'trigger' ? DEFAULT_ICONS.TRIGGER : DEFAULT_ICONS.ACTION"
            class="text-lg text-primary"
          />
        </div>
        <div class="flex-1 min-w-0">
          <div v-if="isEditingLabel" class="flex items-center gap-2">
            <InputText
              ref="labelInputRef"
              v-model="editLabelValue"
              @keyup.enter="saveLabel"
              @blur="saveLabel"
              size="small"
              class="w-full"
            />
          </div>
          <h3 v-else class="text-lg font-semibold text-color truncate flex items-center gap-2">
            <span class="truncate">{{ node.data?.label || $t('builder.configSidebar.node') }}</span>
            <Button
              icon="pi pi-pencil"
              text
              rounded
              size="small"
              class="opacity-0 group-hover/header:opacity-100 transition-opacity !w-6 !h-6 !p-0 shrink-0"
              @click="startEditLabel"
            />
          </h3>
        </div>
      </div>
      <Button
        icon="pi pi-times"
        text
        rounded
        size="small"
        class="shrink-0 mt-1"
        @click="emit('close')"
      />
    </div>

    <!-- Description -->
    <div class="mb-5 px-2 group/desc relative">
      <div v-if="isEditingDescription">
        <Textarea
          ref="descriptionInputRef"
          v-model="editDescriptionValue"
          autoResize
          rows="2"
          class="w-full text-sm"
          @blur="saveDescription"
        />
      </div>
      <div v-else class="flex items-start gap-2 min-h-6">
        <div
          class="text-sm text-muted-color whitespace-pre-wrap"
          :class="{ 'italic opacity-75': !node.data?.description }"
        >
          {{ node.data?.description || $t('builder.configSidebar.noDescription') }}
        </div>
        <Button
          icon="pi pi-pencil"
          text
          rounded
          size="small"
          class="opacity-0 group-hover/desc:opacity-100 transition-opacity !w-6 !h-6 !p-0 shrink-0"
          @click="startEditDescription"
        />
      </div>
    </div>

    <div class="flex-1 overflow-auto p-2">
      <!-- Function form for step/iterator configuration -->
      <FunctionForm
        v-if="
          (node.type === 'iterator' || (node.type !== 'trigger' && node.type !== 'branch')) &&
          functionDefinition?.segments?.length
        "
        ref="functionFormRef"
        :function-def="functionDefinition"
        :arguments="node.data?.arguments || []"
        :upstream-results="upstreamResults"
        :nodes="nodes"
        @update:arguments="onArgumentsUpdate"
        @toggle-reference="onToggleReference"
        class="mb-4"
      />

      <!-- Standard Trigger form for generic trigger configuration -->
      <TriggerForm
        v-if="
          node.type === 'trigger' &&
          triggerDefinition?.segments?.length &&
          !triggerDefinition?.resourceType?.startsWith('automation:trigger-definition:')
        "
        :trigger-def="triggerDefinition"
        :constraints="node.data?.constraints || []"
        @update:constraints="onConstraintsUpdate"
        class="mb-4"
      />

      <!-- Trigger Schema Editor for dynamically defining inputs that the LLM resolves -->
      <div
        v-if="
          node.type === 'trigger' &&
          triggerDefinition?.resourceType?.startsWith('automation:trigger-definition:')
        "
        class="mb-4"
      >
        <TriggerSchemaEditor :trigger-def="triggerDefinition" />
      </div>

      <!-- Branch configuration (only for branch nodes) -->
      <div v-if="node.type === 'branch'" class="mb-4">
        <!-- Gateway type -->
        <div class="mb-4">
          <label class="text-sm font-medium text-color block mb-2">
            {{ $t('builder.branch.gatewayType') }}
          </label>
          <Select
            :model-value="gatewayType"
            :options="gatewayOptions"
            option-label="label"
            option-value="value"
            class="w-full"
            @update:model-value="onGatewayTypeChange"
          >
            <template #value="{ value }">
              <span>{{ gatewayOptions.find(o => o.value === value)?.label }}</span>
            </template>
            <template #option="{ option }">
              <div class="flex flex-col gap-1">
                <span class="font-medium">{{ option.label }}</span>
                <span class="text-xs text-muted-color">{{ option.description }}</span>
              </div>
            </template>
          </Select>
        </div>

        <!-- Branch outputs -->
        <div class="text-sm font-medium text-color mb-2">
          {{ $t('builder.configSidebar.branches') }}
        </div>
        <div class="flex flex-col gap-2">
          <Panel
            toggleable
            v-for="(data, index) in branchOutputs"
            :key="data.edgeId"
            :header="
              index === 0
                ? $t('builder.branch.if')
                : index === branchOutputs.length - 1
                  ? $t('builder.branch.else')
                  : $t('builder.branch.elseIf')
            "
            class="transition-opacity"
            :pt="{
              root: { style: 'overflow: hidden; min-width: 0' },
              toggleableContent: { style: 'overflow: hidden' },
              content: { style: 'padding: 0.5rem !important; overflow: hidden; min-width: 0' },
            }"
            :class="{
              'cursor-grab': index < branchOutputs.length - 1,
              'opacity-50': dragIndex === index,
            }"
            :draggable="index < branchOutputs.length - 1"
            @dragstart="onDragStart(index, $event)"
            @dragover.prevent="onDragOver(index)"
            @dragleave="onDragLeave"
            @drop.prevent="onDrop(index)"
            @dragend="onDragEnd"
          >
            <!-- Condition builder (not for the last/Else branch) -->
            <ConditionBuilder
              v-if="index < branchOutputs.length - 1"
              :model-value="data.condition"
              :edge-id="data.edgeId"
              @update:model-value="onBranchConditionChange(data.edgeId, $event)"
              @toggle-reference="onConditionToggleReference"
            />
            <span v-else class="text-sm text-muted-color italic">
              {{ $t('builder.branch.defaultPath') }}
            </span>
          </Panel>
        </div>
      </div>

      <!-- Spacer -->
      <div class="flex-1" />

      <!-- Branch-specific actions -->
      <div v-if="node.type === 'branch'" class="mb-4">
        <Button
          :label="$t('builder.configSidebar.addElseIf')"
          icon="pi pi-plus"
          severity="secondary"
          outlined
          class="w-full"
          @click="emit('addBranch')"
        />
      </div>
    </div>

    <!-- Delete action at bottom -->
    <div class="p-2 flex flex-col">
      <CInputDelete
        :label="$t('builder.configSidebar.deleteNode')"
        :message="$t('builder.confirmDelete.message')"
        :header="node.data?.label || $t('builder.configSidebar.node')"
        outlined
        size="small"
        @confirm="emit('delete')"
      />
    </div>
  </div>
</template>

<script setup>
import { components } from '@cortezaproject/corteza-vue-next'
import { DEFAULT_ICONS } from '@cortezaproject/corteza-js-next/src/automation/types/icon'
import { computed, nextTick, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { conditionToShort } from '@/utils/taq-parser'
import TaqIcon from '../common/TaqIcon.vue'
import ConditionBuilder from './condition/ConditionBuilder.vue'
import FunctionForm from './form/FunctionForm.vue'
import TriggerForm from './form/TriggerForm.vue'
import TriggerSchemaEditor from './form/TriggerSchemaEditor.vue'

const { CInputDelete } = components

const { t } = useI18n()

const props = defineProps({
  node: { type: Object, required: true },
  edges: { type: Array, default: () => [] },
  nodes: { type: Array, default: () => [] },
  functions: { type: Array, default: () => [] },
  triggers: { type: Array, default: () => [] },
  upstreamResults: { type: Array, default: () => [] },
})

import { provide, toRef } from 'vue'
provide('taq-nodes', toRef(props, 'nodes'))
provide('taq-upstream-results', toRef(props, 'upstreamResults'))

const functionFormRef = ref(null)

const emit = defineEmits([
  'close',
  'delete',
  'addBranch',
  'reorderBranches',
  'updateArguments',
  'updateConstraints',
  'toggleReference',
  'updateGatewayType',
  'updateBranchExpr',
  'updateMetadata',
])

const isEditingLabel = ref(false)
const editLabelValue = ref('')
const labelInputRef = ref(null)

const isEditingDescription = ref(false)
const editDescriptionValue = ref('')
const descriptionInputRef = ref(null)

function startEditLabel() {
  editLabelValue.value = props.node.data?.label || ''
  isEditingLabel.value = true
  nextTick(() => {
    labelInputRef.value?.$el?.focus()
  })
}

function saveLabel() {
  if (!isEditingLabel.value) return
  isEditingLabel.value = false
  if (editLabelValue.value !== props.node.data?.label) {
    emit('updateMetadata', { label: editLabelValue.value })
  }
}

function startEditDescription() {
  editDescriptionValue.value = props.node.data?.description || ''
  isEditingDescription.value = true
  nextTick(() => {
    descriptionInputRef.value?.$el?.focus()
  })
}

function saveDescription() {
  if (!isEditingDescription.value) return
  isEditingDescription.value = false
  if (editDescriptionValue.value !== props.node.data?.description) {
    emit('updateMetadata', { description: editDescriptionValue.value })
  }
}

// Look up function definition from store.functions by node.data.nodeType (which holds the function ref)
const functionDefinition = computed(() => {
  const nodeType = props.node.data?.nodeType
  if (!nodeType) return null
  return props.functions.find(f => f.ref === nodeType) || null
})

// Look up trigger definition from store.triggers by eventType
const triggerDefinition = computed(() => {
  if (props.node.type !== 'trigger') return null
  const eventType = props.node.data?.nodeType
  const resourceType = props.node.data?.resourceType
  if (!eventType) return null
  return (
    props.triggers.find(
      t => t.eventType === eventType && (!resourceType || t.resourceType === resourceType),
    ) || null
  )
})

// Handle argument updates from FunctionForm
function onArgumentsUpdate(args) {
  emit('updateArguments', args)
}

// Handle constraint updates from TriggerForm
function onConstraintsUpdate(constraints) {
  emit('updateConstraints', constraints)
}

// Handle reference toggle from FunctionForm
function onToggleReference(argumentInfo) {
  emit('toggleReference', argumentInfo)
}

// Apply a reference selection from the ReferencePanel
function applyReference(argumentName, ref, target) {
  functionFormRef.value?.onReferenceSelect(argumentName, ref, target)
}

// Expose applyReference for parent (Builder.vue)
defineExpose({ applyReference })

// Get branch outputs from edges
const branchOutputs = ref([])

watch(
  () => [props.edges, props.node.id],
  () => {
    if (props.node.type !== 'branch') {
      branchOutputs.value = []
      return
    }

    // Find edges originating from this branch
    const outputs = props.edges
      .filter(e => e.source === props.node.id)
      .map(edge => {
        const targetNode = props.nodes.find(n => n.id === edge.target)
        const condition = edge.data?.condition || null
        return {
          edgeId: edge.id,
          targetId: edge.target,
          label: targetNode?.data?.label || t('builder.configSidebar.end'),
          condition,
          conditionSummary: condition ? conditionToShort(condition) : '',
        }
      })

    branchOutputs.value = outputs
  },
  { immediate: true, deep: true },
)

// Drag-and-drop state for branch reordering
const dragIndex = ref(null)
const dropTarget = ref(null)

function onDragStart(index, event) {
  dragIndex.value = index
  event.dataTransfer.effectAllowed = 'move'
}

function onDragOver(index) {
  const lastIndex = branchOutputs.value.length - 1
  if (dragIndex.value !== null && dragIndex.value !== index && index < lastIndex) {
    dropTarget.value = index
  }
}

function onDragLeave() {
  dropTarget.value = null
}

function onDrop(index) {
  const lastIndex = branchOutputs.value.length - 1
  if (dragIndex.value === null || dragIndex.value === index || index >= lastIndex) return

  const items = [...branchOutputs.value]
  const [moved] = items.splice(dragIndex.value, 1)
  items.splice(index, 0, moved)
  branchOutputs.value = items

  emit(
    'reorderBranches',
    items.map(b => b.edgeId),
  )
  dragIndex.value = null
  dropTarget.value = null
}

function onDragEnd() {
  dragIndex.value = null
  dropTarget.value = null
}

// Gateway type options and state
const gatewayOptions = [
  {
    label: t('builder.branch.exclusive'),
    value: 'gatewayExclusive',
    description: t('builder.branch.exclusiveHint'),
  },
  {
    label: t('builder.branch.inclusive'),
    value: 'gatewayInclusive',
    description: t('builder.branch.inclusiveHint'),
  },
]

const gatewayType = computed(() => {
  const raw = props.node.data?.nodeType
  return raw === 'gatewayExclusive' || raw === 'gatewayInclusive' ? raw : 'gatewayExclusive'
})

function onGatewayTypeChange(value) {
  emit('updateGatewayType', value)
}

function onBranchConditionChange(edgeId, condition) {
  emit('updateBranchExpr', { edgeId, condition })
}

// Forward condition reference toggle to Builder.vue (opens ReferencePanel)
function onConditionToggleReference(ctx) {
  emit('toggleReference', `condition:${ctx.edgeId}:${ctx.side}:${ctx.rowIndex}`)
}
</script>

<style scoped>
:deep(.p-datatable-tbody > tr) {
  background: transparent !important;
}

:deep(.p-datatable-header) {
  display: none;
}

:deep(.p-datatable .p-datatable-thead) {
  display: none;
}
</style>
