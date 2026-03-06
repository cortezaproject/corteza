<template>
  <div class="flex flex-col gap-3">
    <div class="flex flex-col gap-1">
      <label class="text-primary font-medium text-sm">{{ $t('block.content.label') }}</label>
      <Textarea
        v-model="body"
        rows="12"
        class="w-full font-mono text-sm"
        :placeholder="$t('block.content.placeholder')"
      />
    </div>

    <small class="text-muted-color">
      {{ $t('block.content.interpolationFootnote') }}
      <code>${record.values.fieldName}</code>,
      <code>${recordID}</code>,
      <code>${ownerID}</code>,
      <code>${userID}</code>,
      <code>${user.name}</code>
    </small>
  </div>
</template>

<script setup>
import { computed } from 'vue'

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
