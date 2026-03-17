<template>
  <span
    class="truncate"
    :class="displayLabel !== notSet ? 'text-color-emphasis' : 'italic text-muted-color'"
    :title="displayLabel"
  >
    {{ displayLabel }}
  </span>
</template>

<script setup>
import { onMounted, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { useComposeResourceStore } from '@cortezaproject/corteza-vue-next/src/stores/useComposeResourceStore'

const { t } = useI18n()

const props = defineProps({
  modelValue: { type: [String, Number], default: null },
})

const store = useComposeResourceStore()
const notSet = t('builder.preview.notSet')
const displayLabel = ref(notSet)

async function resolve(id) {
  if (!id) {
    displayLabel.value = t('builder.preview.notSet')
    return
  }

  // Check cache first
  const cached = store.getNamespace(String(id))
  if (cached) {
    displayLabel.value = cached.name || cached.slug || String(id)
    return
  }

  // Async resolve
  try {
    const ns = await store.resolveNamespace(String(id))
    displayLabel.value = ns?.name || ns?.slug || String(id)
  } catch {
    displayLabel.value = String(id)
  }
}

watch(() => props.modelValue, resolve, { immediate: false })
onMounted(() => resolve(props.modelValue))
</script>
