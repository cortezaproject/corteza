<template>
  <div class="flex flex-col gap-3">
    <CRichTextInput v-model="body" class="w-full" />

    <CExpressionHint :scope="scope" :insertable="false" />
  </div>
</template>

<script setup>
import { computed, inject } from 'vue'
import { components } from '@planetcrust/human-vue'
import { useExpressionScope } from '@/sections/compose/composables/useExpressionScope'

const { CRichTextInput } = components

const props = defineProps({
  namespace: { type: Object, default: () => ({}) },
  page: { type: Object, default: () => ({}) },
})

const { scope } = useExpressionScope({ page: computed(() => props.page) })

const block = inject('blockDraft')

const body = computed({
  get: () => block.value.options?.body || '',
  set: v => {
    if (!block.value.options) block.value.options = {}
    block.value.options.body = v
  },
})
</script>
