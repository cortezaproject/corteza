<template>
  <div class="condition-builder">
    <!-- Empty state -->
    <div v-if="!modelValue" class="flex">
      <Button
        :label="$t('builder.condition.addCondition')"
        icon="pi pi-plus"
        outlined
        size="small"
        severity="secondary"
        @click="onAddFirst"
      />
    </div>

    <!-- Single condition (bare node, no group) -->
    <template v-else-if="!isGroup(modelValue)">
      <ConditionRow
        :node="modelValue"
        :edge-id="edgeId"
        :row-index="0"
        @update="onUpdate"
        @delete="onClear"
        @toggle-reference="$emit('toggleReference', $event)"
      />
      <div class="flex items-center gap-2 mt-2">
        <Button
          :label="$t('builder.condition.addCondition')"
          icon="pi pi-plus"
          outlined
          size="small"
          severity="secondary"
          @click="onAddSecond"
        />
      </div>
    </template>

    <!-- Multiple conditions (flat list with inline AND/OR connectors) -->
    <div v-else class="flex flex-col">
      <template v-for="(child, index) in modelValue.args" :key="childKey(child, index)">
        <!-- AND/OR connector between rows -->
        <div v-if="index > 0" class="flex items-center gap-2 py-1">
          <div class="flex-1 border-t border-surface" />
          <button
            class="text-xs font-semibold px-2 py-0.5 rounded-border cursor-pointer transition-colors"
            :class="
              modelValue.ref === 'and'
                ? 'bg-primary text-primary-contrast'
                : 'bg-highlight text-primary'
            "
            @click="onCombinatorToggle"
          >
            {{
              modelValue.ref === 'and' ? $t('builder.condition.and') : $t('builder.condition.or')
            }}
          </button>
          <div class="flex-1 border-t border-surface" />
        </div>

        <ConditionRow
          :node="child"
          :edge-id="edgeId"
          :row-index="index"
          @update="onChildUpdate(index, $event)"
          @delete="onChildDelete(index)"
          @toggle-reference="$emit('toggleReference', $event)"
        />
      </template>

      <!-- Add more -->
      <div class="flex items-center gap-2 mt-2">
        <Button
          :label="$t('builder.condition.addCondition')"
          icon="pi pi-plus"
          outlined
          size="small"
          severity="secondary"
          @click="onAddToGroup"
        />
      </div>
    </div>
  </div>
</template>

<script setup>
import ConditionRow from './ConditionRow.vue'

const props = defineProps({
  modelValue: { type: Object, default: null },
  edgeId: { type: String, default: '' },
})

const emit = defineEmits(['update:modelValue', 'toggleReference'])

function isGroup(node) {
  return node && (node.ref === 'and' || node.ref === 'or')
}

function cloneNode(node) {
  return JSON.parse(JSON.stringify(node))
}

function childKey(child, index) {
  if (!child) return `empty_${index}`
  const symbol = child.args?.[0]?.symbol || ''
  return `row_${index}_${child.ref}_${symbol}`
}

function makeEmptyComparison() {
  return {
    ref: 'eq',
    args: [{ symbol: '', meta: {} }, { value: { '@type': 'String', '@value': '' } }],
  }
}

function onUpdate(node) {
  emit('update:modelValue', node)
}

function onClear() {
  emit('update:modelValue', null)
}

function onAddFirst() {
  emit('update:modelValue', makeEmptyComparison())
}

// Wraps single condition into an AND group with a new empty row
function onAddSecond() {
  emit('update:modelValue', {
    ref: 'and',
    args: [cloneNode(props.modelValue), makeEmptyComparison()],
  })
}

// Toggle AND ↔ OR on the group
function onCombinatorToggle() {
  const updated = cloneNode(props.modelValue)
  updated.ref = updated.ref === 'and' ? 'or' : 'and'
  emit('update:modelValue', updated)
}

function onChildUpdate(index, updatedChild) {
  const updated = cloneNode(props.modelValue)
  updated.args[index] = updatedChild
  emit('update:modelValue', updated)
}

function onChildDelete(index) {
  const updated = cloneNode(props.modelValue)
  updated.args.splice(index, 1)

  if (updated.args.length === 1) {
    // Unwrap: promote single remaining child to bare node
    emit('update:modelValue', updated.args[0])
  } else if (updated.args.length === 0) {
    emit('update:modelValue', null)
  } else {
    emit('update:modelValue', updated)
  }
}

function onAddToGroup() {
  const updated = cloneNode(props.modelValue)
  updated.args.push(makeEmptyComparison())
  emit('update:modelValue', updated)
}
</script>
