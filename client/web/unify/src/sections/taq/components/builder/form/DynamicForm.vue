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

        <!-- Primary fields -->
        <div class="flex flex-col gap-4">
          <DynamicFormField
            v-for="input in primaryInputs(section)"
            :key="input.key"
            :input="input"
            :show-reference-toggle="showReferenceToggle"
            @update:value="(arg, val) => $emit('update:value', arg, val)"
            @toggle-reference="$emit('toggleReference', $event)"
            @clear-reference="$emit('clearReference', $event)"
            @update-reference-source="(arg, val) => $emit('updateReferenceSource', arg, val)"
          />
        </div>

        <!-- Advanced fields, collapsed by default -->
        <div v-if="advancedInputs(section).length" class="flex flex-col gap-3">
          <button
            type="button"
            class="flex items-center gap-2 text-sm text-muted-color hover:text-color w-fit"
            @click="toggleAdvanced(section.key)"
          >
            <span :class="isAdvancedOpen(section.key) ? 'pi pi-chevron-down' : 'pi pi-chevron-right'" />
            {{ $t('builder.form.advanced', 'Advanced options') }}
          </button>

          <div v-if="isAdvancedOpen(section.key)" class="flex flex-col gap-4">
            <DynamicFormField
              v-for="input in advancedInputs(section)"
              :key="input.key"
              :input="input"
              :show-reference-toggle="showReferenceToggle"
              @update:value="(arg, val) => $emit('update:value', arg, val)"
              @toggle-reference="$emit('toggleReference', $event)"
              @clear-reference="$emit('clearReference', $event)"
              @update-reference-source="(arg, val) => $emit('updateReferenceSource', arg, val)"
            />
          </div>
        </div>
      </div>
    </template>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import DynamicFormField from './DynamicFormField.vue'

defineProps({
  processedSegments: {
    type: Array,
    default: () => [],
  },
  showReferenceToggle: {
    type: Boolean,
    default: false,
  },
})

defineEmits(['update:value', 'toggleReference', 'clearReference', 'updateReferenceSource'])

function primaryInputs(section) {
  return (section.inputs || []).filter(i => !i.advanced)
}
function advancedInputs(section) {
  return (section.inputs || []).filter(i => i.advanced)
}

const openAdvanced = ref(new Set())
function toggleAdvanced(key) {
  const next = new Set(openAdvanced.value)
  next.has(key) ? next.delete(key) : next.add(key)
  openAdvanced.value = next
}
function isAdvancedOpen(key) {
  return openAdvanced.value.has(key)
}
</script>
