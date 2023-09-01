<template>
  <div>
    <p v-html="message"></p>
    <b-form-group
      :label="label"
      label-class="text-primary"
    >
      <b-input
        :type="type"
        :disabled="loading"
        v-model="value"
      />
    </b-form-group>
    <b-button
      :disabled="loading"
      @click="$emit('submit', { value: { '@value': value, '@type': 'String' }})"
    >
      {{ pVal('buttonLabel', 'Submit') }}
    </b-button>
  </div>
</template>
<script lang="js">
import base from './base.vue'

const validTypes = [
  'text',
  'number',
  'email',
  'password',
  'search',
  'date',
  'time',
]

export default {
  extends: base,
  name: 'c-prompt-input',

  data () {
    return {
      value: undefined
    }
  },

  beforeMount() {
    this.value = this.pVal('value')
  },

  computed: {
    type () {
      const t = this.pVal('type', 'text')
      if (validTypes.indexOf(t) === -1) {
        return 'text'
      }

      return t
    },

    label () {
      return this.pVal('label', '')
    }
  },

}
</script>
