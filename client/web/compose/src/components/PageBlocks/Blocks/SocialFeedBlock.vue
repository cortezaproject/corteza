<template>
  <PageBlock :block="block">
    <div class="flex items-center justify-center h-full p-3 text-muted-color">
      <div class="text-center">
        <i class="pi pi-twitter text-4xl mb-3" />
        <p class="font-medium text-color mb-2">
          {{ $t('block.socialFeed.label') }}
        </p>
        <p v-if="profileUrl" class="text-sm">
          {{ profileUrl }}
        </p>
        <p v-else class="text-sm italic">
          {{ $t('block.socialFeed.noInput') }}
        </p>
      </div>
    </div>
  </PageBlock>
</template>

<script setup>
import { computed } from 'vue'
import PageBlock from './PageBlock.vue'

const props = defineProps({
  block: { type: Object, required: true },
  namespace: { type: Object, default: () => ({}) },
  page: { type: Object, default: () => ({}) },
  record: { type: Object, default: undefined },
})

const options = computed(() => props.block.options || {})

const profileUrl = computed(() => {
  const { profileSourceField, profileUrl: staticUrl } = options.value

  if (profileSourceField && props.record) {
    const val = props.record.values?.[profileSourceField]
    if (Array.isArray(val) && val.length > 0) return val[0]
    if (val) return val
  }

  return staticUrl || ''
})
</script>
