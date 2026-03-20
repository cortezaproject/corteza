<template>
  <div
    class="condition-group border border-surface rounded-border p-2"
    :class="{ 'border-l-2 border-l-primary': !isRoot }"
  >
    <!-- Children with inline AND/OR connectors between them -->
    <div class="flex flex-col">
      <template v-for="(child, index) in node.args" :key="childKey(child, index)">
        <!-- AND/OR connector between rows (not before the first) -->
        <div v-if="index > 0" class="flex items-center gap-2 py-2">
          <div class="flex-1 border-t border-surface" />
          <button
            class="text-xs font-semibold px-2 py-0.5 rounded-border cursor-pointer transition-colors"
            :class="
              node.ref === 'and' ? 'bg-primary text-primary-contrast' : 'bg-highlight text-primary'
            "
            @click="onCombinatorToggle"
          >
            {{ node.ref === 'and' ? $t('builder.condition.and') : $t('builder.condition.or') }}
          </button>
          <div class="flex-1 border-t border-surface" />
        </div>

        <!-- Nested group -->
        <ConditionGroup
          v-if="isGroup(child)"
          :node="child"
          :depth="depth + 1"
          :edge-id="edgeId"
          :is-root="false"
          @update="onChildUpdate(index, $event)"
          @delete="onChildDelete(index)"
          @toggle-reference="$emit('toggleReference', $event)"
        />

        <!-- Comparison row -->
        <ConditionRow
          v-else
          :node="child"
          :edge-id="edgeId"
          :row-index="index"
          @update="onChildUpdate(index, $event)"
          @delete="onChildDelete(index)"
          @toggle-reference="$emit('toggleReference', $event)"
        />
      </template>
    </div>

    <!-- Actions -->
    <div class="flex items-center gap-2 mt-2">
      <Button
        :label="$t('builder.condition.addCondition')"
        icon="pi pi-plus"
        outlined
        size="small"
        severity="secondary"
        @click="onAddCondition"
      />
      <Button
        v-if="depth < 3"
        :label="$t('builder.condition.addGroup')"
        icon="pi pi-plus"
        outlined
        size="small"
        severity="secondary"
        @click="onAddGroup"
      />

      <!-- Delete group (not for root) -->
      <Button
        v-if="!isRoot"
        icon="pi pi-trash"
        text
        rounded
        size="small"
        severity="danger"
        class="ml-auto"
        @click="emit('delete')"
      />
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import ConditionRow from './ConditionRow.vue'

const { t } = useI18n()

const props = defineProps({
  node: { type: Object, required: true },
  depth: { type: Number, default: 0 },
  edgeId: { type: String, default: '' },
  isRoot: { type: Boolean, default: false },
})

const emit = defineEmits(['update', 'delete', 'toggleReference'])

function isGroup(child) {
  return child && (child.ref === 'and' || child.ref === 'or')
}

function childKey(child, index) {
  if (!child) return `empty_${index}`
  if (child.ref === 'and' || child.ref === 'or') {
    return `group_${index}_${child.ref}_${(child.args || []).length}`
  }
  const symbol = child.args?.[0]?.symbol || ''
  return `row_${index}_${child.ref}_${symbol}`
}

function cloneNode(node) {
  return JSON.parse(JSON.stringify(node))
}

// Toggle between AND and OR by clicking the connector pill
function onCombinatorToggle() {
  const updated = cloneNode(props.node)
  updated.ref = updated.ref === 'and' ? 'or' : 'and'
  emit('update', updated)
}

function onChildUpdate(index, updatedChild) {
  const updated = cloneNode(props.node)
  updated.args[index] = updatedChild
  emit('update', updated)
}

function onChildDelete(index) {
  const updated = cloneNode(props.node)
  updated.args.splice(index, 1)

  if (updated.args.length === 1) {
    emit('update', updated.args[0])
  } else if (updated.args.length === 0) {
    if (props.isRoot) {
      emit('update', null)
    } else {
      emit('delete')
    }
  } else {
    emit('update', updated)
  }
}

function makeEmptyComparison() {
  return {
    ref: 'eq',
    args: [{ symbol: '', meta: {} }, { value: { '@type': 'String', '@value': '' } }],
  }
}

function onAddCondition() {
  const updated = cloneNode(props.node)
  updated.args.push(makeEmptyComparison())
  emit('update', updated)
}

function onAddGroup() {
  const updated = cloneNode(props.node)
  updated.args.push({
    ref: 'and',
    args: [makeEmptyComparison()],
  })
  emit('update', updated)
}
</script>
