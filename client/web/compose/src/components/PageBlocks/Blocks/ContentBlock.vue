<template>
  <PageBlock :block="block">
    <div class="p-3" v-html="body" />
  </PageBlock>
</template>

<script setup>
import { computed } from 'vue'
import PageBlock from './PageBlock.vue'

const props = defineProps({
  block: {
    type: Object,
    required: true,
  },
  namespace: {
    type: Object,
    default: () => ({}),
  },
})

const body = computed(() => props.block.options?.body || '')
</script>
