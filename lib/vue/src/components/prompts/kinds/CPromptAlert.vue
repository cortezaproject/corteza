<template>
  <div>
    <p
      v-if="!!message"
      v-html="message"
    />

    <div class="d-flex justify-content-end gap-1">
      <b-button
        :variant="pVal('buttonVariant', 'primary')"
        :disabled="loading"
        @click="$emit('submit', { confirmed: pRaw('buttonValue', true, 'Boolean') })"
      >
        {{ pVal('buttonLabel', 'OK') }}
      </b-button>
    </div>
  </div>
</template>
<script lang="js">
import base from './base.vue'

export default {
  name: 'CPromptAlert',
  extends: base,
}
</script>
