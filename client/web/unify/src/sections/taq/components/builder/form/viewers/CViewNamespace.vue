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
import { inject, onMounted, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { useNamespaceStore } from '@planetcrust/human-vue'

const { t } = useI18n()

const props = defineProps({
  modelValue: { type: [String, Number], default: null },
})

const $ComposeAPI = inject('$ComposeAPI')
const namespaceStore = useNamespaceStore()
const notSet = t('builder.preview.notSet')
const displayLabel = ref(notSet)

async function resolve(id) {
  if (!id) {
    displayLabel.value = t('builder.preview.notSet')
    return
  }

  try {
    const ns = await namespaceStore.findByID({ namespaceID: String(id) })
    displayLabel.value = ns?.name || ns?.slug || String(id)
  } catch {
    displayLabel.value = String(id)
  }
}

watch(() => props.modelValue, resolve, { immediate: false })
onMounted(() => resolve(props.modelValue))
</script>
