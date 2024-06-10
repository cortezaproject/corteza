<template>
  <div
    class="position-relative"
  >
    <ace-editor
      ref="aceeditor"
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
import { faExpandAlt } from '@fortawesome/free-solid-svg-icons'

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

    autoComplete: {
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

    autoCompleteSuggestions: {
      type: [Array, Object],
      default: () => ([]),
    },

    initExpressions: {
      type: Boolean,
      required: false,
    },

    fontFamily: {
      type: String,
      default: '',
    },

    placeholder: {
      type: String,
      default: '',
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
      require('brace/mode/json')

      require('brace/snippets/text')
      require('brace/snippets/html')
      require('brace/snippets/css')
      require('brace/snippets/scss')
      require('brace/snippets/json')
      require('brace/snippets/javascript')
      require('brace/snippets/json')

      require('brace/theme/chrome')
      require('brace/ext/language_tools')
      require('brace/ext/emmet')

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
        cursorStyle: 'smooth',
        // // minLines: this.height,
        // maxPixelHeight: 100,

        ...(this.autoComplete && {
          enableBasicAutocompletion: true,
          enableLiveAutocompletion: true,
          enableSnippets: true,
          enableEmmet: true,
        }),

        ...(this.fontFamily && { fontFamily: this.fontFamily }),
        ...(this.fontSize && { fontSize: this.fontSize }),
      })

      editor.on('input', this.updatePlaceholder)
      this.updatePlaceholder(undefined, editor)

      if (this.initExpressions) {
        this.processExpressionAutoComplete(editor)
      } else {
        // eslint-disable-next-line @typescript-eslint/no-this-alias
        const self = this
        const staticWordCompleter = {
          getCompletions: function (editor, session, pos, prefix, callback) {
            const autoCompleteSuggestions = self.autoCompleteSuggestions
            callback(
              null,
              autoCompleteSuggestions.map(function ({ caption, value, meta }) {
                return {
                  caption,
                  value,
                  meta,
                }
              }),
            )
          },
        }

        editor.completers.push(staticWordCompleter)
      }
    },

    updatePlaceholder (_, editor) {
      if (!this.placeholder) return

      const shouldShow = !editor.session.getValue().length
      let node = editor.renderer.emptyMessageNode

      if (!shouldShow && node) {
        editor.renderer.scroller.removeChild(editor.renderer.emptyMessageNode)
        editor.renderer.emptyMessageNode = null
      } else if (shouldShow && !node) {
        node = editor.renderer.emptyMessageNode = document.createElement('div')
        node.textContent = this.placeholder
        node.className = 'ace_placeholder'
        node.style.padding = '7px 10px'
        node.style.position = 'absolute'
        node.style.zIndex = 9
        node.style.opacity = 0.5
        editor.renderer.scroller.appendChild(node)
      }
    },

    processExpressionAutoComplete (editor) {
      const staticWordCompleter = {
        getCompletions: (editor, session, pos, prefix, callback) => {
          const context = this.getContext(editor, session, pos)
          const suggestions = this.getSuggestionsForContext(context)

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
              meta: 'variable',
              completer: {
                insertMatch: function (insertEditor, data) {
                  const insertValue = data.value

                  insertEditor.jumpToMatching()
                  const line = session.getLine(pos.row)
                  let lastSpaceIndex = line.lastIndexOf(' ') >= 0 ? line.lastIndexOf(' ') : 0

                  if (lastSpaceIndex > 0) {
                    lastSpaceIndex += 1
                  }

                  insertEditor.session.replace({
                    start: { row: pos.row, column: lastSpaceIndex },
                    end: { row: pos.row, column: pos.column },
                  }, insertValue)
                },
              },
            }
          }))
        },
      }

      editor.completers = [staticWordCompleter]

      editor.commands.on('afterExec', function (e) {
        if (['insertstring', 'Return'].includes(e.command.name) || /^[\w.($]$/.test(e.args)) {
          editor.execCommand('startAutocomplete')
        }
      })

      editor.renderer.setScrollMargin(7, 7)
      editor.renderer.setPadding(10)
    },

    getContext (editor, session, pos) {
      const line = session.getLine(pos.row)
      const lastSpaceIndex = line.lastIndexOf(' ') >= 0 ? line.lastIndexOf(' ') : 0
      const textBeforeCursor = line.slice(lastSpaceIndex, pos.column)
      const context = textBeforeCursor.split('.').slice(0, -1).join('.').trim()

      return context
    },

    getSuggestionsForContext (context) {
      const suggestions = this.autoCompleteSuggestions

      return suggestions[context] || []
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
