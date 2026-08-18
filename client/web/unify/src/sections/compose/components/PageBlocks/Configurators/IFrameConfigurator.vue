<template>
  <div class="flex flex-col gap-3">
    <CFormGroup
      v-if="isRecordPage"
      :label="$t('block.iframe.srcFieldLabel')"
      :description="$t('block.iframe.srcFieldDesc')"
    >
      <CInputModuleField
        v-model="srcField"
        :module-i-d="page?.moduleID"
        :kinds="['Url']"
        :placeholder="$t('block.iframe.pickURLField')"
      />
    </CFormGroup>

    <CFormGroup
      :label="$t('block.iframe.srcLabel')"
      :description="isRecordPage ? $t('block.iframe.srcDesc') : ''"
    >
      <CInputExpression
        ref="srcUrlInput"
        v-model="srcUrl"
        dialect="interpolation"
        :scope="scope"
        :placeholder="$t('block.content.urlPlaceholder')"
      />
    </CFormGroup>

    <CExpressionHint :scope="scope" @insert="srcUrlInput?.insert($event)" />
  </div>
</template>

<script setup>
import { computed, inject, ref } from 'vue'
import { useExpressionScope } from '@/sections/compose/composables/useExpressionScope'

const props = defineProps({
  namespace: { type: Object, default: () => ({}) },
  page: { type: Object, default: () => ({}) },
})

const block = inject('blockDraft')

const isRecordPage = computed(() => !!props.page?.moduleID && props.page.moduleID !== '0')

const srcUrlInput = ref(null)
const { scope } = useExpressionScope({ page: computed(() => props.page) })

function updateOptions(key, value) {
  if (!block.value.options) block.value.options = {}
  block.value.options[key] = value
}

const srcUrl = computed({
  get: () => block.value.options?.src || block.value.options?.url || '',
  set: v => updateOptions('src', v),
})

const srcField = computed({
  get: () => block.value.options?.srcField || '',
  set: v => updateOptions('srcField', v),
})
</script>
