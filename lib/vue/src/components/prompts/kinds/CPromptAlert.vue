<template>
  <div>
    <p v-html="message"></p>
    <div
      class="text-center m-2"
    >
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
