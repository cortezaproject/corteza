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
import { useModuleStore } from '@planetcrust/human-vue'

const { t } = useI18n()

const props = defineProps({
  modelValue: { type: [String, Number], default: null },
  namespaceID: { type: [String, Number], default: null },
})

const $ComposeAPI = inject('$ComposeAPI')
const moduleStore = useModuleStore()
const notSet = t('builder.preview.notSet')
const displayLabel = ref(notSet)

async function resolve(moduleID) {
  if (!moduleID || !props.namespaceID) {
    displayLabel.value = t('builder.preview.notSet')
    return
  }

  try {
    const mod = await moduleStore.findByID({
      namespaceID: String(props.namespaceID),
      moduleID: String(moduleID),
    })
    displayLabel.value = mod?.name || mod?.handle || String(moduleID)
  } catch {
    displayLabel.value = String(moduleID)
  }
}

watch(() => props.modelValue, resolve, { immediate: false })
watch(
  () => props.namespaceID,
  () => resolve(props.modelValue),
)
onMounted(() => resolve(props.modelValue))
</script>
