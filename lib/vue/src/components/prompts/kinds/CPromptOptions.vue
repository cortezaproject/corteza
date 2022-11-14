<template>
  <div>
    <p v-html="message"></p>
    <b-form-group
      :label="pVal('label', 'Input')"
    >
      <b-form-select
        v-if="type === 'select'"
        :options="options"
        :disabled="loading"
        :multiple="multiple"
        v-model="value"
      >
        <template v-if="!multiple" #first>
          <b-form-select-option
            :value="undefined"
            disabled
          >
            -- Please select an option --
          </b-form-select-option>
        </template>
      </b-form-select>
      <b-form-radio-group
        v-if="type === 'radio'"
        v-model="value"
        :disabled="loading"
        :options="options"
      />
    </b-form-group>
    <b-button
      :disabled="loading"
      @click="$emit('submit', { value: encodeValue() })"
    >
      {{ pVal('buttonLabel', 'Submit') }}
    </b-button>
  </div>
</template>
<script lang="js">
import base from './base.vue'

const validTypes = [
  'select',
  'radio',
]

export default {
  extends: base,
  name: 'c-prompt-options',

  data () {
    return {
      value: undefined
    }
  },

  methods: {
    encodeValue () {
      if (Array.isArray(this.value)) {
        return {
          '@type': 'Array',
          '@value': this.value || []
        }
      } else {
        return { '@type': 'String', '@value': this.value }
      }
    },
  },

  beforeMount() {
    let value = this.pVal('value')

    if (this.pVal('multiselect', false)) {
      if (Array.isArray(value)) {
        value = value.map(v => v['@value'])
      } else {
        value = value ? [value] : []
      }
    }

    this.value = value
  },

  computed: {
    options() {
      const out = []
      const options = this.pVal('options', {})
      for (const value in options) {
        out.push({value, text: options[value]})
      }

      return out
    },

    type() {
      const t = this.pVal('type', 'text')
      if (validTypes.indexOf(t) === -1) {
        return 'select'
      }

      return t
    },

    multiple () {
      return this.pVal('multiselect', false)
    }
  },
}
</script>
