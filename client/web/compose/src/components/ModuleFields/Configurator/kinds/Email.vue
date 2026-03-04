<template>
  <div class="flex flex-col gap-4">
    <div class="flex items-center gap-2">
      <Checkbox v-model="field.options.outputPlain" inputId="outputPlain" :binary="true" />
      <label for="outputPlain" class="cursor-pointer">{{ $t('field.kind.email.preventToLink') }}</label>
    </div>

    <CConfiguratorMultiDelimiter :field="field" />
  </div>
</template>

<script setup>
import CConfiguratorMultiDelimiter from '../CConfiguratorMultiDelimiter.vue'

defineProps({
  field: {
    type: Object,
    required: true,
  },
})
</script>
