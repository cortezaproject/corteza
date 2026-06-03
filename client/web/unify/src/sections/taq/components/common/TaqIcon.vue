<template>
  <i v-if="resolved?.type === 'name'" :class="'pi pi-' + resolved.value" />
  <img v-else-if="resolved?.type === 'url'" :src="resolved.value" class="inline-block w-[1em] h-[1em]" />
  <img v-else-if="resolved?.type === 'attachment'" :src="resolved.value" class="inline-block w-[1em] h-[1em]" />
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  icon: { type: Object, default: undefined },
  fallback: { type: Object, default: undefined },
})

const resolved = computed(() => props.icon ?? props.fallback)
</script>
