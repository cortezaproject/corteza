<template>
  <div
    class="position-relative"
  >
    <c-ace-editor
      v-model="expressionValue"
      :lang="lang"
      :min-height="minHeight"
      :show-line-numbers="showLineNumbers"
      :font-size="fontSize"
      :show-popout="showPopout"
      :auto-complete="autoComplete"
      :border="border"
      :auto-complete-suggestions="expressionAutoCompleteValues"
      resizable
    />
  </div>
</template>

<script>
import { components } from '@cortezaproject/corteza-vue-next'
import { EXPRESSION_EDITOR_AUTO_COMPLETE_VALUES } from '../lib/editor-auto-complete.js'

const { CAceEditor } = components

export default {
  components: {
    CAceEditor,
  },

  props: {
    value: {
      type: String,
      default: '',
    },

    lang: {
      type: String,
      default: 'text',
    },

    minHeight: {
      type: String,
      default: '6rem',
    },

    showLineNumbers: {
      type: Boolean,
      default: false,
    },

    fontSize: {
      type: String,
      default: '14px',
    },

    border: {
      type: Boolean,
      default: true,
    },

    showPopout: {
      type: Boolean,
      default: true,
    },

    autoComplete: {
      type: Boolean,
      default: true,
    },
  },

  data () {
    return {
      expressionAutoCompleteValues: EXPRESSION_EDITOR_AUTO_COMPLETE_VALUES,
    }
  },

  computed: {
    expressionValue: {
      get () {
        return this.value
      },

      set (value = '') {
        this.$emit('update:value', value)
      },
    },
  },
}
</script>

<style scoped>
.popout {
  z-index: 7;
  bottom: 0;
  right: 0;
}
</style>
