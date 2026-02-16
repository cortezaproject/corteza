<template>
  <component
    :is="viewerComponent"
    :field="field"
    :record="record"
    :namespace="namespace"
    :value-only="valueOnly"
    :extra-options="extraOptions"
    :disable-click="disableClick"
  />
</template>

<script setup>
import { computed } from 'vue'
import { resolveFieldViewer } from './registry'

const props = defineProps({
  field: {
    type: Object,
    required: true,
  },
  record: {
    type: Object,
    required: true,
  },
  namespace: {
    type: Object,
    default: () => ({}),
  },
  valueOnly: {
    type: Boolean,
    default: false,
  },
  extraOptions: {
    type: Object,
    default: () => ({}),
  },
  disableClick: {
    type: Boolean,
    default: false,
  },
})

const viewerComponent = computed(() => resolveFieldViewer(props.field.kind))
</script>
