<template>
  <div class="flex items-center gap-1">
    <div
      class="reference-chip flex items-center gap-2 bg-surface border border-surface flex-1 text-sm"
      :style="chipStyle"
      @click="focusInput"
    >
      <i class="pi pi-link text-primary" />
      <input
        ref="inputEl"
        :value="label"
        class="reference-input flex-1 min-w-0 bg-transparent border-none outline-none text-color text-sm"
        @input="emit('update:label', $event.target.value)"
      />
    </div>
    <Button
      icon="pi pi-times"
      text
      rounded
      size="small"
      severity="secondary"
      @click="emit('clear')"
    />
  </div>
</template>

<script setup>
import { computed, ref } from 'vue'

const props = defineProps({
  label: {
    type: String,
    default: 'Reference',
  },
  size: {
    type: String,
    default: 'normal',
  },
})

const emit = defineEmits(['click', 'clear', 'update:label'])

const inputEl = ref(null)

function focusInput() {
  inputEl.value?.focus()
}

const chipStyle = computed(() => {
  const prefix = props.size === 'small' ? '--p-form-field-sm' : '--p-form-field'
  return {
    padding: `var(${prefix}-padding-y) var(${prefix}-padding-x)`,
    borderRadius: 'var(--p-form-field-border-radius)',
  }
})
</script>

<style scoped>
.reference-chip:focus-within {
  outline: none;
  border-color: var(--p-inputtext-focus-border-color);
}

.reference-input {
  font: inherit;
}
</style>
