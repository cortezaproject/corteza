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

const { t } = useI18n()

const props = defineProps({
  modelValue: { type: [String, Number], default: null },
})

const $SystemAPI = inject('$SystemAPI')
const notSet = t('builder.preview.notSet')
const displayLabel = ref(notSet)

async function resolve(agentID) {
  if (!agentID || agentID === '0') {
    displayLabel.value = t('builder.preview.notSet')
    return
  }

  try {
    const agent = await $SystemAPI.agentRead({ agentID: String(agentID) })
    displayLabel.value = agent.meta?.short || agent.handle || String(agentID)
  } catch {
    displayLabel.value = String(agentID)
  }
}

watch(() => props.modelValue, resolve, { immediate: false })
onMounted(() => resolve(props.modelValue))
</script>
