<template>
  <div>
    <p
      v-if="!!message"
      v-html="message"
    />

    <div class="d-flex align-items-between justify-content-center gap-2">
      <b-button
        :variant="pVal('confirmButtonVariant', 'primary')"
        :disabled="loading"
        @click="$emit('submit', { value: pRaw('confirmButtonValue', true, 'Boolean') })"
      >
        {{ pVal('confirmButtonLabel', 'Yes') }}
      </b-button>

      <b-button
        :disabled="loading"
        :variant="pVal('rejectButtonVariant', 'light')"
        @click="$emit('submit', { value: pRaw('rejectButtonValue', false, 'Boolean') })"
      >
        {{ pVal('rejectButtonLabel', 'No') }}
      </b-button>
    </div>
  </div>
</template>
<script lang="js">
import base from './base.vue'

export default {
  name: 'CPromptChoice',
  extends: base,
}
</script>
