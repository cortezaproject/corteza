<template>
  <div class="flex flex-col gap-4">
    <p v-if="message" class="m-0 whitespace-pre-wrap" v-html="message" />
    <div class="flex flex-wrap gap-2">
      <Button
        v-bind="vS('confirmButtonVariant', 'primary')"
        class="flex-1"
        :disabled="loading"
        :label="pVal('confirmButtonLabel', tF('general.label.yes', 'Yes'))"
        @click="$emit('submit', { value: pRaw('confirmButtonValue', true, 'Boolean') })"
      />
      <Button
        v-bind="vS('rejectButtonVariant', 'light')"
        class="flex-1"
        :disabled="loading"
        :label="pVal('rejectButtonLabel', tF('general.label.no', 'No'))"
        @click="$emit('submit', { value: pRaw('rejectButtonValue', false, 'Boolean') })"
      />
    </div>
  </div>
</template>

<script>
import base from './base.vue'

export default {
  name: 'CPromptChoice',
  extends: base,
  emits: ['submit'],
}
</script>
