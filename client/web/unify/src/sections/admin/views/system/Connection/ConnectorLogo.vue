<template>
  <span
    class="flex items-center justify-center rounded-md shrink-0 overflow-hidden bg-surface-50 dark:bg-surface-800 border border-surface-200 dark:border-surface-700"
    :class="sizeClass"
    :style="showImg ? undefined : tileStyle"
  >
    <img
      v-if="showImg"
      :src="src"
      :alt="name || slug"
      class="w-2/3 h-2/3 object-contain"
      @error="failed = true"
    />
    <span v-else class="font-semibold text-white" :class="textClass">{{ initial }}</span>
  </span>
</template>

<script setup>
import { computed, ref } from 'vue'
import { brandIconUrl } from '@planetcrust/human-js/src/automation/types/icon'

const props = defineProps({
  // Brand slug from the catalog meta.icon (e.g. "google-sheets").
  icon: { type: String, default: '' },
  name: { type: String, default: '' },
  size: { type: String, default: 'md' }, // md | lg
})

const failed = ref(false)
const slug = computed(() => (props.icon || '').trim().toLowerCase())
const showImg = computed(() => !!slug.value && !failed.value)
const src = computed(() => brandIconUrl(slug.value))
const initial = computed(() => (props.name || slug.value || '?').charAt(0).toUpperCase())

const sizeClass = computed(() => (props.size === 'lg' ? 'w-12 h-12' : 'w-10 h-10'))
const textClass = computed(() => (props.size === 'lg' ? 'text-xl' : 'text-base'))

// Deterministic tile colour so the fallback is stable per connector.
const tileStyle = computed(() => {
  const s = props.name || slug.value || '?'
  let h = 0
  for (let i = 0; i < s.length; i++) {
    h = (h * 31 + s.charCodeAt(i)) % 360
  }
  return `background-color: hsl(${h}, 52%, 45%)`
})
</script>
