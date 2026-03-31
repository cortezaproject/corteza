<template>
  <div class="flex flex-col gap-4">
    <p v-if="message" class="m-0 whitespace-pre-wrap" v-html="message" />
    <div class="flex justify-end">
      <Button
        :disabled="loading"
        :label="pVal('buttonLabel', 'Ok')"
        @click="$emit('submit', { confirmed: pRaw(undefined, true, 'Boolean') })"
      />
    </div>
  </div>
</template>

<script>
import base from './base.vue'

export default {
  name: 'CPromptAlert',
  extends: base,
  emits: ['submit'],
}
</script>
