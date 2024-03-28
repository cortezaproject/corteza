<template>
  <div :class="classes">
    <span
      v-for="(v, index) of formattedValue"
      :key="index"
      :class="{ 'd-block': field.options.multiDelimiter === '\n' }"
    >
      <span v-if="field.options.outputPlain">
        {{ fixUrl(v) }}{{ index !== formattedValue.length - 1 ? field.options.multiDelimiter : '' }}
      </span>

      <a
        v-else
        :href="fixUrl(v)"
        target="_blank"
        rel="noopener noreferrer"
        @click.stop
      >
        {{ fixUrl(v) }}{{ index !== formattedValue.length - 1 ? field.options.multiDelimiter : '' }}
      </a>
    </span>

    <errors :errors="errors" />
  </div>
</template>
<script>
import base from './base'
import { trimUrlFragment, trimUrlQuery, trimUrlPath, onlySecureUrl } from '../url'

export default {
  extends: base,

  computed: {
    formattedValue () {
      return this.field.isMulti ? this.value : [this.value].filter(v => v)
    },
  },

  methods: {
    fixUrl (value) {
      // run through all the attributes
      if (this.field.options.trimFragment) {
        value = trimUrlFragment(value)
      }
      if (this.field.options.trimQuery) {
        value = trimUrlQuery(value)
      }
      if (this.field.options.trimPath) {
        value = trimUrlPath(value)
      }
      if (this.field.options.onlySecure) {
        value = onlySecureUrl(value)
      }

      return value
    },
  },
}
</script>
