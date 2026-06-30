<template>
  <div class="p-3 border-t border-surface flex gap-2 items-end shrink-0 bg-surface">
    <Textarea
      ref="taRef"
      v-model="text"
      :placeholder="placeholder"
      :disabled="disabled"
      class="flex-1 resize-none"
      :class="scrollable ? '!overflow-y-auto' : '!overflow-y-hidden'"
      :style="{ maxHeight: `calc(${maxRows} * 1.5rem + 1rem)` }"
      rows="1"
      autoResize
      @input="updateScrollable"
      @keydown.enter.exact.prevent="submit"
    />
    <Button
      icon="pi pi-send"
      :disabled="disabled || !text.trim()"
      :loading="loading"
      @click="submit"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, nextTick, watch } from 'vue'

// Shared chat composer. Auto-growing textarea: grows with typed lines up to
// `maxRows`, then scrolls internally. Enter sends, Shift+Enter inserts a newline.
const props = defineProps({
  placeholder: {
    type: String,
    default: '',
  },
  // Disables typing and sending (e.g. while executing or when send not allowed).
  disabled: {
    type: Boolean,
    default: false,
  },
  // Shows a spinner on the send button without disabling typing.
  loading: {
    type: Boolean,
    default: false,
  },
  // Max visible rows before the textarea scrolls internally.
  maxRows: {
    type: Number,
    default: 6,
  },
})

const emit = defineEmits<{
  (_e: 'send', _input: string): void
}>()

const text = ref('')

// Show the scrollbar only when the textarea is actually capped at maxHeight and
// has overflowing content — autoResize otherwise leaves a permanent 1px scroll.
const taRef = ref()
const scrollable = ref(false)
function updateScrollable() {
  nextTick(() => {
    const el = taRef.value?.$el
    if (!el) return
    scrollable.value = el.scrollHeight > el.clientHeight + 1
  })
}
// Recompute after the text is cleared on send (textarea shrinks back).
watch(text, updateScrollable)

function submit() {
  const txt = text.value.trim()
  if (!txt || props.disabled) return
  emit('send', text.value)
  text.value = ''
}
</script>
