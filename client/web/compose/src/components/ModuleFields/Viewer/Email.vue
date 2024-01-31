<template>
  <div :class="classes">
    <span
      v-for="(v, index) of formattedValue"
      :key="index"
      :class="{ 'd-block': field.options.multiDelimiter === '\n' }"
    >
      <span v-if="field.options.outputPlain">
        {{ v }}{{ index !== formattedValue.length - 1 ? field.options.multiDelimiter : '' }}
      </span>

      <a
        v-else
        :href="'mailto:' + formattedValue"
        target="_blank"
        rel="noopener noreferrer"
        @click.stop
      >
        {{ v }}{{ index !== formattedValue.length - 1 ? field.options.multiDelimiter : '' }}
      </a>
    </span>

    <errors :errors="errors" />
  </div>
</template>
<script>
import base from './base'

export default {
  extends: base,

  computed: {
    formattedValue () {
      return this.field.isMulti ? this.value : [this.value].filter(v => v)
    },
  },
}
</script>
