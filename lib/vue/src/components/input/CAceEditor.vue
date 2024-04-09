<template>
  <div
    class="position-relative"
  >
    <ace-editor
      v-model="editorValue"
      :lang="lang"
      :mode="lang"
      theme="chrome"
      width="100%"
      :height="height"
      :class="{ 'border-0 rounded-0': !border }"
      v-on="$listeners"
      @init="editorInit"
    />

    <b-button
      v-if="showPopout"
      variant="link"
      class="popout position-absolute px-2 py-1 mr-3"
      @click="$emit('open')"
    >
      <font-awesome-icon
        :icon="['fas', 'expand-alt']"
      />
    </b-button>
  </div>
</template>

<script>
import AceEditor from 'vue2-ace-editor'
import { library } from '@fortawesome/fontawesome-svg-core'
import { faExpandAlt} from '@fortawesome/free-solid-svg-icons'

library.add(faExpandAlt)

export default {
  components: {
    AceEditor,
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

    height: {
      type: String,
      default: '80',
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
      default: false,
    },

    readOnly: {
      type: Boolean,
      default: false,
    },

    highlightActiveLine: {
      type: Boolean,
      default: false,
    },

    showPrintMargin: {
      type: Boolean,
      default: false,
    },
  },

  computed: {
    editorValue: {
      get () {
        return this.value
      },

      set (value = '') {
        this.$emit('update:value', value)
      },
    },
  },

  methods: {
    editorInit (editor) {
      require('brace/mode/text')
      require('brace/mode/html')
      require('brace/mode/css')
      require('brace/mode/scss')
      require('brace/mode/json')
      require('brace/mode/javascript')
      require('brace/theme/chrome')

      editor.setOptions({
        tabSize: 2,
        fontSize: this.fontSize,
        wrap: true,
        indentedSoftWrap: false,
        showPrintMargin: this.showPrintMargin,
        showLineNumbers: this.showLineNumbers,
        showGutter: this.showLineNumbers,
        displayIndentGuides: this.lang !== 'text',
        useWorker: false,
        readOnly: this.readOnly,
        highlightActiveLine: this.highlightActiveLine
      })
    },
  },
}
</script>

<style lang="scss" scoped>
.popout {
  z-index: 7;
  bottom: 0;
  right: 0;
}
</style>

<style lang="scss">
// Remove from server/assets/src/scss/main/18201141_custom_webapp.scss when all ace-editors use c-ace-editor
.ace_editor {
  color: var(--black) !important;
  background-color: var(--white) !important;
  border-radius: 0.25rem;
  border: 2px solid var(--extra-light);

  .ace_gutter {
    background-color: var(--light) !important;
    color: var(--black) !important;

    .ace_gutter-active-line {
      background-color: var(--extra-light) !important;
    }
  }
}
</style>
