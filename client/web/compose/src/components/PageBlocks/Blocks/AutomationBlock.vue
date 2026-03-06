<template>
  <PageBlock :block="block">
    <div v-if="buttons.length" class="flex flex-wrap gap-2 p-3">
      <Button
        v-for="(btn, i) in buttons"
        :key="i"
        :label="btn.label || $t('block.automation.noLabel')"
        :severity="mapVariant(btn.variant)"
        @click="runButton(btn)"
      />
    </div>
    <div v-else class="flex items-center justify-center h-full p-3 text-muted-color italic">
      {{ $t('block.automation.noScripts') }}
    </div>
  </PageBlock>
</template>

<script setup>
import { computed, inject } from 'vue'
import { useI18n } from 'vue-i18n'
import PageBlock from './PageBlock.vue'

const { t } = useI18n()

const props = defineProps({
  block: { type: Object, required: true },
  namespace: { type: Object, default: () => ({}) },
  page: { type: Object, default: () => ({}) },
})

const $toast = inject('$toast')

const buttons = computed(() => props.block.options?.buttons || [])

function mapVariant(variant) {
  const map = {
    primary: undefined,
    secondary: 'secondary',
    light: 'secondary',
    dark: 'contrast',
    success: 'success',
    danger: 'danger',
    warning: 'warn',
  }
  return map[variant] || undefined
}

async function runButton(btn) {
  // Automation execution — placeholder for workflow/script execution
  if (btn.workflowID) {
    $toast?.toastInfo(t('block.automation.noScript'))
  } else if (btn.script) {
    $toast?.toastInfo(t('block.automation.noScript'))
  }
}
</script>
