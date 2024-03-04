<template>
  <div>
    <p v-if="!!message" v-html="message" />

    <div class="d-flex align-items-between justify-content-center gap-2">
      <b-button
        @click="$emit('submit', { value: pRaw('confirmButtonValue', true, 'Boolean') })"
        :variant="pVal('confirmButtonVariant', 'primary')"
        :disabled="loading"
      >
        {{ pVal('confirmButtonLabel', 'Yes') }}
      </b-button>

      <b-button
        @click="$emit('submit', { value: pRaw('rejectButtonValue', false, 'Boolean') })"
        :disabled="loading"
        :variant="pVal('rejectButtonVariant', 'light')"
      >
        {{ pVal('rejectButtonLabel', 'No') }}
      </b-button>
    </div>
  </div>
</template>
<script lang="js">
import base from './base.vue'

export default {
  extends: base,
  name: 'c-prompt-choice',
}
</script>
