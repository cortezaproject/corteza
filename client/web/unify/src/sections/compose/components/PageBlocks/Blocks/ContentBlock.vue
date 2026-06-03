<template>
  <PageBlock :block="block">
    <div class="p-3" style="white-space: pre-wrap" v-html="contentBody" />
  </PageBlock>
</template>

<script setup>
import { computed, inject } from 'vue'
import PageBlock from './PageBlock.vue'
import { evaluatePrefilter } from '../../../lib/record-filter'

const props = defineProps({
  block: {
    type: Object,
    required: true,
  },
  namespace: {
    type: Object,
    default: () => ({}),
  },
  record: {
    type: Object,
    default: undefined,
  },
})

const $auth = inject('$auth', {})

const contentBody = computed(() => {
  const { body = '' } = props.block.options || {}
  if (!body) return ''

  try {
    const record = props.record
    const user = $auth?.user || {}

    return evaluatePrefilter(body, {
      record,
      user,
      recordID: record?.recordID || '0',
      ownerID: record?.ownedBy || '0',
      userID: user?.userID || '0',
    })
  } catch {
    return body
  }
})
</script>
