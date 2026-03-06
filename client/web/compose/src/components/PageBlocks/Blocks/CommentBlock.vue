<template>
  <PageBlock :block="block">
    <div class="p-3 flex flex-col gap-2">
      <div v-if="comments.length" class="flex flex-col gap-3">
        <div v-for="(c, i) in comments" :key="i" class="border border-surface rounded p-3">
          <div class="text-sm text-muted-color mb-1">
            {{ c.createdBy || $t('block.comment.noComments') }}
          </div>
          <div>{{ c.content }}</div>
        </div>
      </div>
      <div v-else class="text-muted-color italic text-center py-4">
        {{ $t('block.comment.noComments') }}
      </div>
    </div>
  </PageBlock>
</template>

<script setup>
import { computed } from 'vue'
import PageBlock from './PageBlock.vue'

defineProps({
  block: { type: Object, required: true },
  namespace: { type: Object, default: () => ({}) },
  page: { type: Object, default: () => ({}) },
})

// Comments are stored as records in a comment module
// Full implementation requires integration with record store
const comments = computed(() => [])
</script>
