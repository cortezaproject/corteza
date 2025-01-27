<template>
  <div>
    <p
      v-if="!!message"
      v-html="message"
    />

    <b-form-group
      :label="pVal('label', 'Input')"
      label-class="text-primary"
    >
      <b-form-select
        v-if="type === 'select'"
        v-model="value"
        :options="options"
        :disabled="loading"
        :multiple="multiple"
      >
        <template
          v-if="!multiple"
          #first
        >
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
      variant="primary"
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
  name: 'CPromptOptions',
  extends: base,

  data () {
    return {
      value: undefined,
    }
  },

  computed: {
    options () {
      const out = []
      const options = this.pVal('options', {})
      for (const value in options) {
        out.push({ value, text: options[value] })
      }

      return out
    },

    type () {
      const t = this.pVal('type', 'text')
      if (validTypes.indexOf(t) === -1) {
        return 'select'
      }

      return t
    },

    multiple () {
      return this.pVal('multiselect', false)
    },
  },

  beforeMount () {
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

  methods: {
    encodeValue () {
      if (Array.isArray(this.value)) {
        return {
          '@type': 'Array',
          '@value': this.value || [],
        }
      } else {
        return { '@type': 'String', '@value': this.value }
      }
    },
  },
}
</script>
