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
import { useUserResolver } from '@planetcrust/human-vue/src/composables/useUserResolver'

const { t } = useI18n()

const props = defineProps({
  modelValue: { type: [String, Number], default: null },
})

const { resolveUser, formatUser } = useUserResolver()
const notSet = t('builder.preview.notSet')
const displayLabel = ref(notSet)

async function resolve(userID) {
  if (!userID) {
    displayLabel.value = t('builder.preview.notSet')
    return
  }

  try {
    const user = await resolveUser(String(userID))
    displayLabel.value = formatUser(user) || String(userID)
  } catch {
    displayLabel.value = String(userID)
  }
}

watch(() => props.modelValue, resolve, { immediate: false })
onMounted(() => resolve(props.modelValue))
</script>
