<template>
  <i v-if="resolved?.type === 'name'" :class="'pi pi-' + resolved.value" />
  <img
    v-else-if="resolved?.type === 'url'"
    :src="resolved.value"
    class="inline-block w-[1em] h-[1em]"
  />
  <img
    v-else-if="resolved?.type === 'brand' && !failed"
    :src="brandIconUrl(resolved.value)"
    class="inline-block w-[1em] h-[1em]"
    @error="failed = true"
  />
  <i v-else-if="resolved?.type === 'brand' && failed" :class="'pi pi-' + brandFallback" />
  <img
    v-else-if="resolved?.type === 'attachment'"
    :src="resolved.value"
    class="inline-block w-[1em] h-[1em]"
  />
</template>

<script setup>
import { computed, ref, watch } from 'vue'
import { brandIconUrl } from '@planetcrust/human-js/src/automation/types/icon'

const props = defineProps({
  icon: { type: Object, default: undefined },
  fallback: { type: Object, default: undefined },
})

const resolved = computed(() => props.icon ?? props.fallback)

// A brand logo that fails to load falls back to a named icon, never a broken image.
const failed = ref(false)
const brandFallback = computed(() =>
  props.fallback?.type === 'name' ? props.fallback.value : 'folder',
)
watch(
  () => resolved.value?.value,
  () => {
    failed.value = false
  },
)
</script>
