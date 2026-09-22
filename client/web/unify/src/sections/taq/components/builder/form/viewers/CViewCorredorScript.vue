<template>
  <span class="flex flex-col min-w-0" :title="scriptName">
    <span class="truncate" :class="scriptName ? 'text-color-emphasis' : 'italic text-muted-color'">
      {{ scriptName ? label : notSet }}
    </span>
    <!-- The stored value: the script's path inside the extension -->
    <small v-if="scriptName && label !== scriptName" class="truncate text-muted-color">
      {{ scriptName }}
    </small>
  </span>
</template>

<script setup>
import { inject, onMounted, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { resolveCorredorScriptLabel } from './corredor-scripts'

const { t } = useI18n()

const props = defineProps({
  modelValue: { type: [String, Number], default: null },
})

const $SystemAPI = inject('$SystemAPI', null)
const notSet = t('builder.preview.notSet')

const scriptName = ref('')
const label = ref('')

async function resolve(value) {
  const name = value ? String(value) : ''
  scriptName.value = name
  label.value = name
  if (!name) return

  const resolved = await resolveCorredorScriptLabel($SystemAPI, name)
  if (scriptName.value === name) label.value = resolved
}

watch(() => props.modelValue, resolve, { immediate: false })
onMounted(() => resolve(props.modelValue))
</script>
