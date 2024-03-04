<template>
  <div>
    <p v-if="!!message" v-html="message" />

    <div class="d-flex justify-content-end gap-1">
      <b-button
        @click="$emit('submit', { confirmed: pRaw('buttonValue', true, 'Boolean') })"
        :variant="pVal('buttonVariant', 'primary')"
        :disabled="loading"
      >
        {{ pVal('buttonLabel', 'OK') }}
      </b-button>
    </div>
  </div>
</template>
<script lang="js">
import base from './base.vue'

export default {
  extends: base,
  name: 'c-prompt-alert',
}
</script>
