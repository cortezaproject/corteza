<template>
  <PageBlock :block="block">
    <div v-if="src !== 'about:blank'" class="h-full">
      <iframe
        ref="iframeRef"
        :src="src"
        class="w-full h-full border-0"
        :title="block.title || $t('block.iframe.label')"
        sandbox="allow-scripts allow-same-origin allow-popups allow-forms"
      />
    </div>
    <div v-else class="flex items-center justify-center h-full p-3 text-muted-color italic">
      {{ $t('block.iframe.noInput') }}
    </div>
  </PageBlock>
</template>

<script setup>
import { computed, inject, ref } from 'vue'
import PageBlock from './PageBlock.vue'
import { evaluatePrefilter } from '../../../lib/record-filter'

const props = defineProps({
  block: { type: Object, required: true },
  namespace: { type: Object, default: () => ({}) },
  page: { type: Object, default: () => ({}) },
  record: { type: Object, default: undefined },
})

const $auth = inject('$auth', {})
const iframeRef = ref(null)

const src = computed(() => {
  const { srcField, src: srcUrl } = props.block.options || {}
  const blank = 'about:blank'
  let url = srcUrl

  // If srcField is set and we have a record, use the field value
  if (srcField && props.record) {
    url = props.record.values?.[srcField]
  }

  if (!url) return blank

  // Interpolate variables like ${record.values.X}, ${userID}, etc.
  const record = props.record
  const user = $auth?.user || {}

  const interpolatedURL = evaluatePrefilter(url, {
    record,
    user,
    recordID: record?.recordID || '0',
    ownerID: record?.ownedBy || '0',
    userID: user?.userID || '0',
  })

  return interpolatedURL || blank
})
</script>
