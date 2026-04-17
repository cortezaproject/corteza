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
import { useComposeResourceStore } from '@planetcrust/human-vue/src/stores/useComposeResourceStore'

const { t } = useI18n()

const props = defineProps({
  modelValue: { type: [String, Number], default: null },
  namespaceID: { type: [String, Number], default: null },
})

const store = useComposeResourceStore()
const notSet = t('builder.preview.notSet')
const displayLabel = ref(notSet)

async function resolve(moduleID) {
  if (!moduleID || !props.namespaceID) {
    displayLabel.value = t('builder.preview.notSet')
    return
  }

  // Check cache first
  const cached = store.getModule(String(props.namespaceID), String(moduleID))
  if (cached) {
    displayLabel.value = cached.name || cached.handle || String(moduleID)
    return
  }

  // Async resolve
  try {
    const mod = await store.resolveModule(String(props.namespaceID), String(moduleID))
    displayLabel.value = mod?.name || mod?.handle || String(moduleID)
  } catch {
    displayLabel.value = String(moduleID)
  }
}

watch(() => props.modelValue, resolve, { immediate: false })
watch(() => props.namespaceID, () => resolve(props.modelValue))
onMounted(() => resolve(props.modelValue))
</script>
