<template>
  <div>
    <p
      v-if="formatted"
      :style="{ 'white-space': field.options.useRichTextEditor && 'pre-line' }"
      :class="viewerClasses"
      v-html="formatted"
    />
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

const formatted = computed(() => {
  const v = value.value
  if (v === undefined || v === null) return ''

  const delimiter = props.field.options?.multiDelimiter || ', '

  if (props.field.isMulti && Array.isArray(v)) {
    return v.join(delimiter)
  }

  return String(v)
})

const viewerClasses = computed(() => {
  const classes = []
  const { fieldID } = props.field
  const { textStyles = {} } = props.extraOptions

  if (props.field.isMulti || props.field.options?.multiLine) {
    classes.push('multiline')
  } else if (textStyles.wrappedFields && !textStyles.wrappedFields.includes(fieldID)) {
    classes.push('text-nowrap')
  }

  return classes
})
</script>

<style scoped>
.multiline {
  white-space: pre-line;
}
</style>
