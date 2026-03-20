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
import { computed } from 'vue'
import { components } from '@cortezaproject/corteza-vue-next'

const { CRichTextInput } = components

const props = defineProps({
  block: { type: Object, required: true },
  namespace: { type: Object, default: () => ({}) },
  page: { type: Object, default: () => ({}) },
})

const emit = defineEmits(['update:block'])

const body = computed({
  get: () => props.block.options?.body || '',
  set: v => {
    emit('update:block', {
      ...props.block,
      options: { ...props.block.options, body: v },
    })
  },
})
</script>
