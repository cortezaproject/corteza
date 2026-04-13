<template>
  <div class="flex gap-3 items-stretch">
    <!-- Available fields -->
    <div class="flex-1 flex flex-col border border-surface rounded-border overflow-hidden min-w-0">
      <div class="flex items-center justify-between py-1 pl-2 pr-0 border-b border-surface">
        <div class="flex items-center gap-2">
          <span class="text-sm font-medium text-color">{{ availableLabel }}</span>
          <span class="text-xs text-muted-color">({{ filteredAvailable.length }})</span>
        </div>
        <Button
          :label="selectAllLabel"
          severity="secondary"
          text
          size="small"
          :disabled="filteredAvailable.length === 0"
          @click="selectAll"
        />
      </div>
      <div class="p-2 border-b border-surface">
        <CInputSearch v-model="availableSearch" :placeholder="searchPlaceholder" size="small" />
      </div>
      <div
        :class="['flex-1 overflow-y-auto p-1.5 flex flex-col gap-2', listClass]"
        @dragover.prevent
        @drop="onDropToAvailable"
      >
        <div
          v-for="field in filteredAvailable"
          :key="field.name"
          class="flex items-center gap-2 py-1.5 px-2.5 border border-surface rounded-border cursor-pointer select-none text-sm text-color transition-colors hover:bg-emphasis"
          draggable="true"
          @dragstart="onDragStart($event, field, 'available')"
          @click.prevent="selectField(field)"
        >
          <span class="overflow-hidden text-ellipsis whitespace-nowrap">
            {{ field.label || field.name }}
          </span>
        </div>
        <div
          v-if="filteredAvailable.length === 0 && availableSearch"
          class="flex items-center justify-center h-full p-4 text-muted-color text-xs"
        >
          {{ noItemsLabel }}
        </div>
      </div>
    </div>

    <!-- Selected fields -->
    <div class="flex-1 flex flex-col border border-surface rounded-border overflow-hidden min-w-0">
      <div class="flex items-center justify-between py-1 pl-2 pr-0 border-b border-surface">
        <div class="flex items-center gap-2">
          <span class="text-sm font-medium text-color">{{ selectedLabel }}</span>
          <span class="text-xs text-muted-color">({{ selected.length }})</span>
        </div>
        <Button
          :label="unselectAllLabel"
          severity="secondary"
          text
          size="small"
          :disabled="selected.length === 0"
          @click="deselectAll"
        />
      </div>
      <div class="p-2 border-b border-surface">
        <CInputSearch v-model="selectedSearch" :placeholder="searchPlaceholder" size="small" />
      </div>
      <div
        :class="['flex-1 overflow-y-auto p-1.5 flex flex-col gap-2', listClass]"
        @dragover.prevent
        @dragenter.prevent
        @drop="onDropToSelected($event)"
      >
        <div
          v-for="(field, index) in filteredSelected"
          :key="field.name"
          class="flex items-center gap-2 py-1.5 px-2.5 border border-surface rounded-border cursor-pointer select-none text-sm text-color transition-colors hover:bg-emphasis"
          :class="{ 'border-t-2 !border-t-primary': dropTargetIndex === index }"
          draggable="true"
          @dragstart="onDragStart($event, field, 'selected')"
          @dragover.prevent="onDragOverSelectedItem($event, index)"
          @dragleave="onDragLeaveSelectedItem"
          @drop.stop="onDropToSelectedAt($event, index)"
          @click.prevent="deselectField(field)"
        >
          <span class="overflow-hidden text-ellipsis whitespace-nowrap">
            {{ field.label || field.name }}
          </span>
        </div>
        <div
          v-if="filteredSelected.length === 0 && selectedSearch"
          class="flex items-center justify-center h-full p-4 text-muted-color text-xs"
        >
          {{ noItemsLabel }}
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed, ref } from 'vue'
import CInputSearch from './CInputSearch.vue'

const props = defineProps({
  /** All fields from the module (must include name, label) */
  allFields: {
    type: Array,
    default: () => [],
  },
  /** Currently selected field names (ordered) */
  modelValue: {
    type: Array,
    default: () => [],
  },
  /** Label for the available panel header */
  availableLabel: {
    type: String,
    default: 'Available fields',
  },
  /** Label for the selected panel header */
  selectedLabel: {
    type: String,
    default: 'Selected fields',
  },
  /** Label for the select all button */
  selectAllLabel: {
    type: String,
    default: 'Select all',
  },
  /** Label for the unselect all button */
  unselectAllLabel: {
    type: String,
    default: 'Unselect all',
  },
  /** Placeholder for the search inputs */
  searchPlaceholder: {
    type: String,
    default: 'Search...',
  },
  /** Text shown when no items match */
  noItemsLabel: {
    type: String,
    default: 'No items found',
  },
  /** Tailwind class controlling the max-height of each list panel */
  listClass: {
    type: String,
    default: 'max-h-80',
  },
})

const emit = defineEmits(['update:modelValue'])

// Search state
const availableSearch = ref('')
const selectedSearch = ref('')

// Drag state
const dragSource = ref(null) // 'available' | 'selected'
const dragField = ref(null)
const dropTargetIndex = ref(-1)

// Resolved selected fields (in order)
const selected = computed(() => {
  return props.modelValue.map(name => props.allFields.find(f => f.name === name)).filter(Boolean)
})

// Available fields (not selected)
const available = computed(() => {
  const selectedSet = new Set(props.modelValue)
  return props.allFields.filter(f => !selectedSet.has(f.name))
})

// Filtered lists
const filteredAvailable = computed(() => {
  if (!availableSearch.value) return available.value
  const q = availableSearch.value.toLowerCase()
  return available.value.filter(
    f => (f.label || f.name).toLowerCase().includes(q) || f.name.toLowerCase().includes(q),
  )
})

const filteredSelected = computed(() => {
  if (!selectedSearch.value) return selected.value
  const q = selectedSearch.value.toLowerCase()
  return selected.value.filter(
    f => (f.label || f.name).toLowerCase().includes(q) || f.name.toLowerCase().includes(q),
  )
})

// --- Actions ---

function selectField(field) {
  emit('update:modelValue', [...props.modelValue, field.name])
}

function deselectField(field) {
  emit(
    'update:modelValue',
    props.modelValue.filter(n => n !== field.name),
  )
}

function selectAll() {
  const allNames = props.allFields.map(f => f.name)
  emit('update:modelValue', allNames)
}

function deselectAll() {
  emit('update:modelValue', [])
}

// --- Drag and Drop ---

function onDragStart(event, field, source) {
  dragSource.value = source
  dragField.value = field
  event.dataTransfer.effectAllowed = 'move'
  event.dataTransfer.setData('text/plain', field.name)
}

function onDropToAvailable() {
  if (!dragField.value) return
  if (dragSource.value === 'selected') {
    deselectField(dragField.value)
  }
  resetDrag()
}

function onDragOverSelectedItem(event, index) {
  event.preventDefault()
  dropTargetIndex.value = index
}

function onDragLeaveSelectedItem() {
  dropTargetIndex.value = -1
}

function onDropToSelected(event) {
  if (!dragField.value) return

  if (dragSource.value === 'available') {
    emit('update:modelValue', [...props.modelValue, dragField.value.name])
  }
  resetDrag()
}

function onDropToSelectedAt(event, targetIndex) {
  if (!dragField.value) return

  const fieldName = dragField.value.name
  const currentNames = [...props.modelValue]

  if (dragSource.value === 'available') {
    const targetField = filteredSelected.value[targetIndex]
    const realIndex = currentNames.indexOf(targetField?.name)
    if (realIndex >= 0) {
      currentNames.splice(realIndex, 0, fieldName)
    } else {
      currentNames.push(fieldName)
    }
    emit('update:modelValue', currentNames)
  } else if (dragSource.value === 'selected') {
    const fromIndex = currentNames.indexOf(fieldName)
    if (fromIndex < 0) return

    currentNames.splice(fromIndex, 1)

    const targetField = filteredSelected.value[targetIndex]
    let toIndex = currentNames.indexOf(targetField?.name)
    if (toIndex < 0) toIndex = currentNames.length
    currentNames.splice(toIndex, 0, fieldName)

    emit('update:modelValue', currentNames)
  }

  resetDrag()
}

function resetDrag() {
  dragSource.value = null
  dragField.value = null
  dropTargetIndex.value = -1
}
</script>
