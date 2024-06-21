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
  </div>
</template>

<script>
import AceEditor from 'vue2-ace-editor'

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
      default: '80px',
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

    suggestionTree: {
      type: Object,
      default: () => ({})
    }
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
      require('brace/mode/javascript')

      require('brace/theme/chrome')
      require('brace/ext/language_tools')

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
        highlightActiveLine: this.highlightActiveLine,
        enableBasicAutocompletion: true,
        enableLiveAutocompletion: true,
      })
      
      const staticWordCompleter = {
        getCompletions: (editor, session, pos, prefix, callback) => {
          const context = this.getContext(editor, session, pos);
          const suggestions = this.getSuggestionsForContext(context);

          callback(null, suggestions.map(suggestion => ({
            caption: suggestion,
            value: suggestion,
            meta: "variable"
          })));
        }
      };

      editor.completers = [staticWordCompleter]

      editor.commands.on("afterExec", function (e) {
        if (e.command.name == "insertstring" && /^[\w.(${]$/.test(e.args)) {
          editor.execCommand("startAutocomplete");
        }
      });

      editor.renderer.setScrollMargin(7, 7)
      editor.renderer.setPadding(10)
    },
    getContext (editor, session, pos) {
      const line = session.getLine(pos.row).replace(/[\$\{\'\(\)]/g, '');
      const lastSpaceIndex = line.lastIndexOf(' ') >= 0 ? line.lastIndexOf(' ') : 0;
      const textBeforeCursor = line.slice(lastSpaceIndex, pos.column);
      const context = textBeforeCursor.split('.').slice(0, -1).join('.').trim();

      return context;
    },
    getSuggestionsForContext (context) {
      const suggestions = this.suggestionTree;
      return suggestions[context] || [];
    },
  },
}
</script>

<style lang="scss">
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
