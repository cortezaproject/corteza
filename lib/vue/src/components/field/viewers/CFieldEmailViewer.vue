<template>
  <div :class="viewerClasses">
    <span v-for="(v, index) in formattedValues" :key="index" :class="{ block: isNewlineDelimiter }">
      <span v-if="field.options.outputPlain || disableClick">
        {{ v }}{{ index !== formattedValues.length - 1 ? delimiter : '' }}
      </span>

      <a
        v-else
        :href="'mailto:' + v"
        target="_blank"
        rel="noopener noreferrer"
        class="text-primary hover:underline"
        @click.stop
      >
        {{ v }}{{ index !== formattedValues.length - 1 ? delimiter : '' }}
      </a>
    </span>
  </div>
</template>

<script setup>
import { computed } from 'vue'

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

const value = computed(() => {
  if (props.field.isSystem) {
    return props.record[props.field.name]
  }
  return props.record?.values?.[props.field.name]
})

const formattedValues = computed(() => {
  const v = value.value
  if (props.field.isMulti && Array.isArray(v)) {
    return v.filter(val => val)
  }
  return [v].filter(val => val)
})

const delimiter = computed(() => props.field.options?.multiDelimiter || ', ')

const isNewlineDelimiter = computed(() => delimiter.value === '\n')
</script>
