<template>
  <div class="c-color-picker inline-flex items-center">
    <!-- Trigger: slot or default swatch -->
    <slot name="trigger" :toggle="toggle" :color="modelValue">
      <button
        ref="swatchBtn"
        type="button"
        class="c-color-swatch"
        :style="swatchStyle"
        @click="toggle"
      >
        <svg v-if="shape === 'circle'" viewBox="0 0 32 32" :style="{ width, height }">
          <defs>
            <pattern :id="patternId" width="12" height="12" patternUnits="userSpaceOnUse">
              <rect fill="#ccc" x="0" y="0" width="6" height="6" />
              <rect fill="#ccc" x="6" y="6" width="6" height="6" />
              <rect fill="#fff" x="6" y="0" width="6" height="6" />
              <rect fill="#fff" x="0" y="6" width="6" height="6" />
            </pattern>
          </defs>
          <circle cx="16" cy="16" r="15" :fill="`url(#${patternId})`" />
          <circle cx="16" cy="16" r="15" :fill="displayColor" />
          <circle cx="16" cy="16" r="15" fill="none" stroke="var(--p-content-border-color)" stroke-width="1" />
        </svg>
        <svg v-else viewBox="0 0 32 32" :style="{ width, height }">
          <defs>
            <pattern :id="patternId" width="12" height="12" patternUnits="userSpaceOnUse">
              <rect fill="#ccc" x="0" y="0" width="6" height="6" />
              <rect fill="#ccc" x="6" y="6" width="6" height="6" />
              <rect fill="#fff" x="6" y="0" width="6" height="6" />
              <rect fill="#fff" x="0" y="6" width="6" height="6" />
            </pattern>
          </defs>
          <rect x="0.5" y="0.5" width="31" height="31" rx="4" :fill="`url(#${patternId})`" />
          <rect x="0.5" y="0.5" width="31" height="31" rx="4" :fill="displayColor" />
          <rect x="0.5" y="0.5" width="31" height="31" rx="4" fill="none" stroke="var(--p-content-border-color)" stroke-width="1" />
        </svg>
      </button>
    </slot>

    <span v-if="showText" class="ml-2 text-sm text-muted-color">
      {{ modelValue || emptyLabel }}
    </span>

    <!-- Popover picker -->
    <Popover ref="popover" class="c-color-popover">
      <div class="flex flex-col gap-3 p-2">
        <!-- Color picker -->
        <ColorPicker
          v-model="pickerHex"
          inline
          class="c-color-picker-inline"
        />

        <!-- Alpha slider -->
        <div class="flex items-center gap-2">
          <label class="text-xs text-muted-color whitespace-nowrap">Alpha</label>
          <input
            v-model.number="alpha"
            type="range"
            min="0"
            max="255"
            class="c-alpha-slider flex-1"
          >
          <span class="text-xs w-8 text-right">{{ Math.round(alpha / 255 * 100) }}%</span>
        </div>

        <!-- Hex display -->
        <div class="flex items-center gap-2">
          <InputText
            v-model="hexInput"
            class="flex-1 text-xs"
            size="small"
            @keydown.enter="applyHexInput"
          />
        </div>

        <!-- Actions -->
        <div v-if="defaultValue" class="flex justify-end">
          <Button
            :label="'Default'"
            size="small"
            severity="secondary"
            text
            @click="resetDefault"
          />
        </div>
      </div>
    </Popover>
  </div>
</template>

<script setup>
import { ref, computed, watch } from 'vue'

const props = defineProps({
  modelValue: {
    type: String,
    default: '#000000FF',
  },
  defaultValue: {
    type: String,
    default: '',
  },
  showText: {
    type: Boolean,
    default: false,
  },
  emptyLabel: {
    type: String,
    default: 'transparent',
  },
  width: {
    type: String,
    default: '28px',
  },
  height: {
    type: String,
    default: '28px',
  },
  shape: {
    type: String,
    default: 'square',
    validator: (v) => ['square', 'circle'].includes(v),
  },
})

const emit = defineEmits(['update:modelValue'])

// Unique pattern ID to avoid SVG conflicts when multiple pickers exist
const patternId = `checker-${Math.random().toString(36).slice(2, 8)}`

const popover = ref(null)
const swatchBtn = ref(null)

// Internal state
const pickerHex = ref('000000')
const alpha = ref(255)

// Parse incoming modelValue into pickerHex + alpha
function parseColor(val) {
  if (!val || val === 'transparent') {
    pickerHex.value = '000000'
    alpha.value = 0
    return
  }

  // Strip #
  let hex = val.replace(/^#/, '')

  if (hex.length === 8) {
    // #RRGGBBAA
    pickerHex.value = hex.substring(0, 6)
    alpha.value = parseInt(hex.substring(6, 8), 16)
  } else if (hex.length === 6) {
    pickerHex.value = hex
    alpha.value = 255
  } else if (hex.length === 3) {
    pickerHex.value = hex[0] + hex[0] + hex[1] + hex[1] + hex[2] + hex[2]
    alpha.value = 255
  } else {
    pickerHex.value = '000000'
    alpha.value = 255
  }
}

// Initialize
parseColor(props.modelValue)

watch(() => props.modelValue, (val) => {
  if (!internalUpdate) {
    parseColor(val)
  }
  internalUpdate = false
})

// Emit live on every change
let internalUpdate = false

function emitColor() {
  internalUpdate = true
  if (alpha.value === 0) {
    emit('update:modelValue', 'transparent')
  } else {
    emit('update:modelValue', hex8.value)
  }
}

watch(pickerHex, () => {
  // If fully transparent and user picks a new color, restore full opacity
  if (alpha.value === 0) alpha.value = 255
  emitColor()
})
watch(alpha, () => emitColor())

// Computed hex8 from pickerHex + alpha
const hex8 = computed(() => {
  const a = Math.round(alpha.value).toString(16).padStart(2, '0')
  return `#${pickerHex.value}${a}`.toUpperCase()
})

// Display color for swatch (CSS rgba)
const displayColor = computed(() => {
  // When modelValue is empty and a defaultValue exists, show the default color in the swatch
  if ((!props.modelValue || props.modelValue === 'transparent') && props.defaultValue) {
    const dv = props.defaultValue.trim()

    // Handle rgb/rgba format directly
    if (dv.startsWith('rgb')) {
      return dv
    }

    // Handle hex format
    const hex = dv.replace(/^#/, '')
    if (hex.length >= 6) {
      const r = parseInt(hex.substring(0, 2), 16)
      const g = parseInt(hex.substring(2, 4), 16)
      const b = parseInt(hex.substring(4, 6), 16)
      const a = hex.length === 8 ? parseInt(hex.substring(6, 8), 16) / 255 : 1
      return `rgba(${r}, ${g}, ${b}, ${a})`
    }
  }

  const hex = pickerHex.value
  const r = parseInt(hex.substring(0, 2), 16)
  const g = parseInt(hex.substring(2, 4), 16)
  const b = parseInt(hex.substring(4, 6), 16)
  const a = alpha.value / 255
  return `rgba(${r}, ${g}, ${b}, ${a})`
})

// Swatch style
const swatchStyle = computed(() => ({
  width: props.width,
  height: props.height,
  borderRadius: props.shape === 'circle' ? '50%' : '4px',
}))

// Editable hex input
const hexInput = computed({
  get: () => hex8.value,
  set: () => { /* handled by applyHexInput */ },
})

function applyHexInput() {
  const el = document.querySelector('.c-color-popover input[type="text"]')
  if (el) {
    parseColor(el.value)
  }
}

function toggle(event) {
  popover.value?.toggle(event)
}

function resetDefault() {
  parseColor(props.defaultValue)
  emitColor()
}

defineExpose({ toggle })
</script>

<style>
.c-color-swatch {
  border: none;
  background: none;
  padding: 0;
  cursor: pointer;
  display: inline-flex;
  flex-shrink: 0;
}

.c-color-swatch:hover {
  opacity: 0.85;
}

.c-color-popover .p-popover-content {
  padding: 0 !important;
}

.c-color-picker-inline .p-colorpicker-preview {
  display: none;
}

/* Alpha slider */
.c-alpha-slider {
  -webkit-appearance: none;
  appearance: none;
  height: 8px;
  border-radius: 4px;
  background: linear-gradient(to right,
    transparent 0%,
    var(--p-text-color, #000) 100%
  ),
  repeating-conic-gradient(#ccc 0% 25%, #fff 0% 50%) 50% / 8px 8px;
  outline: none;
}

.c-alpha-slider::-webkit-slider-thumb {
  -webkit-appearance: none;
  appearance: none;
  width: 14px;
  height: 14px;
  border-radius: 50%;
  background: white;
  border: 2px solid var(--p-content-border-color);
  cursor: pointer;
}

.c-alpha-slider::-moz-range-thumb {
  width: 14px;
  height: 14px;
  border-radius: 50%;
  background: white;
  border: 2px solid var(--p-content-border-color);
  cursor: pointer;
}
</style>
