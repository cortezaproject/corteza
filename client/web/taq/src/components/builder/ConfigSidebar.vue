<template>
  <div class="flex flex-col h-full">
    <!-- Header -->
    <div class="flex items-center justify-between mb-4">
      <h3 class="text-lg font-semibold text-color">
        {{ node.data?.label || $t('builder.configSidebar.node') }}
      </h3>
      <Button icon="pi pi-times" text rounded size="small" @click="emit('close')" />
    </div>

    <!-- Description -->
    <div v-if="node.data?.description" class="text-sm text-muted-color mb-4">
      {{ node.data.description }}
    </div>

    <!-- Function form for step configuration -->
    <FunctionForm
      v-if="node.type !== 'trigger' && functionDefinition?.segments?.length"
      ref="functionFormRef"
      :function-def="functionDefinition"
      :arguments="node.data?.arguments || []"
      :upstream-results="upstreamResults"
      :nodes="nodes"
      @update:arguments="onArgumentsUpdate"
      @toggle-reference="onToggleReference"
      class="mb-4"
    />

    <!-- Trigger form for trigger configuration -->
    <TriggerForm
      v-if="node.type === 'trigger' && triggerDefinition?.segments?.length"
      :trigger-def="triggerDefinition"
      :constraints="node.data?.constraints || []"
      @update:constraints="onConstraintsUpdate"
      class="mb-4"
    />

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
      <DataTable
        :value="branchOutputs"
        :reorderableRows="true"
        @row-reorder="onRowReorder"
        :pt="{
          table: { class: 'w-full' },
          bodyrow: { class: 'border-b border-surface last:border-0' },
        }"
      >
        <Column
          rowReorder
          headerStyle="width: 1.5rem; padding: 0"
          bodyStyle="width: 1.5rem; padding: 0.25rem"
          :pt="{ rowreordericon: { class: 'pi pi-bars text-muted-color cursor-grab' } }"
        />
        <Column field="label" header="">
          <template #body="{ data, index }">
            <div class="flex flex-col gap-1 py-2">
              <div class="flex items-center gap-2">
                <span class="text-sm text-color truncate flex-1">{{ data.label }}</span>
                <Tag
                  v-if="index === 0"
                  :value="$t('builder.branch.if')"
                  severity="info"
                  class="shrink-0"
                />
                <Tag
                  v-else-if="index === branchOutputs.length - 1"
                  :value="$t('builder.branch.else')"
                  severity="secondary"
                  class="shrink-0"
                />
                <Tag
                  v-else
                  :value="$t('builder.branch.elseIf')"
                  severity="warn"
                  class="shrink-0"
                />
              </div>
              <!-- Expression input (not for the last/Else branch) -->
              <InputText
                v-if="index < branchOutputs.length - 1"
                :model-value="data.expr || ''"
                :placeholder="$t('builder.branch.exprPlaceholder')"
                size="small"
                class="w-full"
                style="font-family: monospace; font-size: 0.8rem"
                @update:model-value="onBranchExprChange(data.edgeId, $event)"
              />
              <span v-else class="text-sm text-muted-color italic">
                {{ $t('builder.branch.defaultPath') }}
              </span>
            </div>
          </template>
        </Column>
      </DataTable>
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

    <!-- Delete action at bottom -->
    <div class="pt-4 border-t border-surface">
      <CInputDelete
        :label="$t('builder.configSidebar.deleteNode')"
        :message="$t('builder.confirmDelete.message')"
        :header="node.data?.label || $t('builder.configSidebar.node')"
        outlined
        class="w-full"
        @confirm="emit('delete')"
      />
    </div>
  </div>
</template>

<script setup>
import { components } from '@cortezaproject/corteza-vue-next'
import { computed, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import FunctionForm from './form/FunctionForm.vue'
import TriggerForm from './form/TriggerForm.vue'

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
])

const functionFormRef = ref(null)

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
        return {
          edgeId: edge.id,
          targetId: edge.target,
          label: targetNode?.data?.label || t('builder.configSidebar.end'),
          expr: edge.data?.expr || '',
        }
      })

    branchOutputs.value = outputs
  },
  { immediate: true, deep: true },
)

// Handle row reorder
function onRowReorder(event) {
  branchOutputs.value = event.value
  emit(
    'reorderBranches',
    branchOutputs.value.map(b => b.edgeId),
  )
}

// Gateway type options and state
const gatewayOptions = [
  { label: t('builder.branch.exclusive'), value: 'excl', description: t('builder.branch.exclusiveHint') },
  { label: t('builder.branch.inclusive'), value: 'incl', description: t('builder.branch.inclusiveHint') },
]

const gatewayType = computed(() => {
  const raw = props.node.data?.nodeType
  // Normalize: existing branches may have empty ref or non-gateway ref
  return raw === 'excl' || raw === 'incl' ? raw : 'excl'
})

function onGatewayTypeChange(value) {
  emit('updateGatewayType', value)
}

function onBranchExprChange(edgeId, expr) {
  emit('updateBranchExpr', { edgeId, expr })
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
