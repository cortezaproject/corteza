<template>
  <div>
    <span
      v-for="(v, index) in resolvedValues"
      :key="index"
      :class="{ block: isNewlineDelimiter, 'mt-1': isNewlineDelimiter && index !== 0 }"
    >
      <Tag
        v-if="isBadgeDisplay"
        :value="v.text"
        :pt="{ root: { style: v.style } }"
        class="mr-1"
      />
      <span v-else>{{ v.text }}{{ index !== resolvedValues.length - 1 ? delimiter : '' }}</span>
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

const isBadgeDisplay = computed(() => props.field.options?.displayType === 'badge')
const delimiter = computed(() => props.field.options?.multiDelimiter || ', ')
const isNewlineDelimiter = computed(() => delimiter.value === '\n')

const resolvedValues = computed(() => {
  const v = value.value
  let vals

  if (props.field.isMulti) {
    vals = Array.isArray(v) ? v : []
  } else {
    vals = v !== undefined && v !== null ? [v] : []
  }

  const options = props.field.options?.options || []

  return vals.map(val => resolveValue(val, options)).filter(item => item && item.text)
})

function resolveValue(val, options) {
  const opt = options.find(o => o.value === val) || { text: val }

  return {
    text: opt.text || val,
    style: getOptionStyle(opt),
  }
}

// Theme-independent defaults so badges are always visible
const DEFAULT_BADGE_TEXT_COLOR = '#FFFFFFFF'
const DEFAULT_BADGE_BG_COLOR = '#09344EFF'

function toColor(hex) {
  if (!hex || hex === 'transparent') return null
  return hex.startsWith('#') ? hex : '#' + hex
}

function getOptionStyle(opt) {
  const style = {}

  if (isBadgeDisplay.value) {
    const optStyle = opt.style || {}
    style.color = toColor(optStyle.textColor) || DEFAULT_BADGE_TEXT_COLOR
    style.backgroundColor = toColor(optStyle.backgroundColor) || DEFAULT_BADGE_BG_COLOR
  }

  return style
}
</script>
