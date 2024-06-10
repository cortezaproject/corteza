<template>
  <div
    class="position-relative"
  >
    <c-ace-editor
      v-model="editorValue"
      :initializeEditor="editorInit"
      v-bind="{ ...$attrs, ...$props }"
      v-on="$listeners"
    />
  </div>
</template>

<script>
import CAceEditor from './CAceEditor.vue'

export default {
  components: {
    CAceEditor
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

          callback(null, suggestions.map(suggestion => {
            let caption = ''
            let value = ''

            if (typeof suggestion === 'string') {
              caption = suggestion
              value = suggestion
            } else {
              caption = suggestion.caption
              value = suggestion.value
            }

            return {
              caption,
              value,
              meta: "variable",
              completer: {
                insertMatch: function (insertEditor, data) {
                  let insertValue = data.value;

                  insertEditor.jumpToMatching();
                  const line = session.getLine(pos.row)
                  let lastSpaceIndex = line.lastIndexOf(' ') >= 0 ? line.lastIndexOf(' ') : 0;

                  if (lastSpaceIndex > 0) {
                    lastSpaceIndex += 1
                  }

                  insertEditor.session.replace({
                    start: { row: pos.row, column: lastSpaceIndex },
                    end: { row: pos.row, column: pos.column }
                  }, insertValue);
                }
              }
            }
          }))
        }
      }

      editor.completers = [staticWordCompleter]

      editor.commands.on("afterExec", function (e) {
        if (["insertstring", "Return"].includes(e.command.name) || /^[\w.($]$/.test(e.args)) {
          editor.execCommand("startAutocomplete");
        }
      });

      editor.renderer.setScrollMargin(7, 7)
      editor.renderer.setPadding(10)
    },
    getContext (editor, session, pos) {
      const line = session.getLine(pos.row)
      const lastSpaceIndex = line.lastIndexOf(' ') >= 0 ? line.lastIndexOf(' ') : 0;
      const textBeforeCursor = line.slice(lastSpaceIndex, pos.column);
      const context = textBeforeCursor.split('.').slice(0, -1).join('.').trim();

      return context
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
