<template>
  <div class="flex flex-col gap-3">
    <CRichTextInput v-model="body" class="w-full" />

    <InterpolationFootnote :is-record-page="isRecordPage" />
  </div>
</template>

<script setup>
import { computed, inject } from 'vue'
import { components } from '@planetcrust/human-vue'
import InterpolationFootnote from '@/sections/compose/components/Common/InterpolationFootnote.vue'

const { CRichTextInput } = components

const props = defineProps({
  namespace: { type: Object, default: () => ({}) },
  page: { type: Object, default: () => ({}) },
})

const isRecordPage = computed(() => !!props.page?.moduleID && props.page.moduleID !== '0')

const block = inject('blockDraft')

const body = computed({
  get: () => block.value.options?.body || '',
  set: v => {
    if (!block.value.options) block.value.options = {}
    block.value.options.body = v
  },
})
</script>
