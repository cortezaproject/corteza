<template>
  <div class="flex flex-col gap-3">
    <div class="flex flex-col gap-1 border border-surface rounded-border">
      <CRichTextInput v-model="body" class="w-full" />
    </div>

    <small class="text-muted-color">
      {{ $t('block.content.interpolationFootnote') }}
      <code>${record.values.fieldName}</code>
      ,
      <code>${recordID}</code>
      ,
      <code>${ownerID}</code>
      ,
      <code>${userID}</code>
      ,
      <code>${user.name}</code>
    </small>
  </div>
</template>

<script setup>
import { computed, inject } from 'vue'
import { components } from '@planetcrust/human-vue'

const { CRichTextInput } = components

defineProps({
  namespace: { type: Object, default: () => ({}) },
  page: { type: Object, default: () => ({}) },
})

const block = inject('blockDraft')

const body = computed({
  get: () => block.value.options?.body || '',
  set: v => {
    if (!block.value.options) block.value.options = {}
    block.value.options.body = v
  },
})
</script>
