<template>
  <div class="flex flex-col gap-4">
    <div class="flex items-center gap-2">
      <Checkbox v-model="field.options.outputPlain" inputId="outputPlain" :binary="true" />
      <label for="outputPlain" class="cursor-pointer">{{ $t('field.kind.email.preventToLink') }}</label>
    </div>

    <CConfiguratorMultiDelimiter />
  </div>
</template>

<script setup>
import { inject } from 'vue'
import CConfiguratorMultiDelimiter from '../CConfiguratorMultiDelimiter.vue'

const field = inject('fieldDraft')
</script>
