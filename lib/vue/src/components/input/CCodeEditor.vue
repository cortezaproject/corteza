<template>
  <div ref="editorContainer" class="code-editor" />
</template>

<script setup>
import { ref, watch, onMounted, onBeforeUnmount, nextTick } from 'vue'
import { EditorState, Prec } from '@codemirror/state'
import {
  EditorView,
  tooltips,
  keymap,
  lineNumbers,
  highlightActiveLine,
  highlightActiveLineGutter,
} from '@codemirror/view'
import { defaultKeymap, indentWithTab } from '@codemirror/commands'
import {
  autocompletion,
  clearSnippet,
  closeCompletion,
  completionKeymap,
  completionStatus,
} from '@codemirror/autocomplete'
import { html } from '@codemirror/lang-html'
import { json } from '@codemirror/lang-json'
import { oneDark } from '@codemirror/theme-one-dark'
import {
  syntaxHighlighting,
  defaultHighlightStyle,
  bracketMatching,
  foldGutter,
} from '@codemirror/language'

const props = defineProps({
  modelValue: {
    type: String,
    default: '',
  },
  language: {
    type: String,
    default: 'html', // 'html', 'json', 'text'
  },
  readOnly: {
    type: Boolean,
    default: false,
  },
  minHeight: {
    type: String,
    default: '300px',
  },
  dark: {
    type: Boolean,
    default: false,
  },
  // Suggestions while typing: the language's own (HTML tags and attributes)
  // and whatever these extensions add. Null leaves the editor without any.
  // Fixed when the editor is built.
  assist: {
    type: Array,
    default: null,
  },
})

const emit = defineEmits(['update:modelValue'])

const editorContainer = ref(null)
let view = null
let isUpdating = false

function getLanguageExtension() {
  switch (props.language) {
    case 'html':
      return html()
    case 'json':
      return json()
    default:
      return []
  }
}

function getExtensions() {
  const extensions = [
    lineNumbers(),
    highlightActiveLine(),
    highlightActiveLineGutter(),
    bracketMatching(),
    foldGutter(),
    syntaxHighlighting(defaultHighlightStyle, { fallback: true }),
    keymap.of([...(props.assist ? completionKeymap : []), ...defaultKeymap, indentWithTab]),
    EditorView.lineWrapping,
    EditorView.updateListener.of(update => {
      if (update.docChanged && !isUpdating) {
        emit('update:modelValue', update.state.doc.toString())
      }
    }),
    // Theme — match surface background
    EditorView.theme({
      '&': {
        minHeight: props.minHeight,
        fontSize: '13px',
        fontFamily: 'ui-monospace, SFMono-Regular, "SF Mono", Menlo, Consolas, monospace',
      },
      '.cm-content': {
        padding: '8px 0',
      },
      '.cm-gutters': {
        borderRight: '1px solid var(--p-content-border-color)',
        backgroundColor: 'var(--p-content-background)',
        color: 'var(--p-text-muted-color)',
      },
      '.cm-activeLineGutter': {
        backgroundColor: 'var(--p-content-hover-background)',
        color: 'var(--p-text-color)',
      },
      '&.cm-focused': {
        outline: 'none',
      },
      '.cm-activeLine': {
        backgroundColor: 'var(--p-content-hover-background)',
      },
    }),
  ]

  const lang = getLanguageExtension()
  if (lang) extensions.push(lang)

  if (props.assist) {
    extensions.push(
      autocompletion({ icons: false }),
      // In the body, so the editor's rounded clip does not cut a suggestion's
      // note or a hover note short; above a dialog the editor may sit in.
      tooltips({ parent: document.body }),
      // Escape dismisses the suggestion list, or leaves the fields of a call
      // just written, and stops there, so an editor inside a dialog does not
      // also close the dialog. It outranks completionKeymap, which would
      // otherwise close the list first.
      Prec.highest(
        EditorView.domEventHandlers({
          keydown(event, v) {
            if (event.key !== 'Escape') return false
            if (completionStatus(v.state)) closeCompletion(v)
            else if (!clearSnippet(v)) return false
            event.stopPropagation()
            return true
          },
        }),
      ),
      EditorView.theme({
        '.cm-tooltip': {
          zIndex: 3000,
          backgroundColor: 'var(--p-content-background)',
          color: 'var(--p-text-color)',
          border: '1px solid var(--p-content-border-color)',
          borderRadius: 'var(--p-content-border-radius)',
        },
        '.cm-tooltip-autocomplete > ul > li[aria-selected]': {
          backgroundColor: 'var(--p-highlight-background)',
          color: 'var(--p-highlight-color)',
        },
        '.cm-completionDetail': {
          color: 'var(--p-text-muted-color)',
        },
      }),
      ...props.assist,
    )
  }

  if (props.dark) {
    extensions.push(oneDark)
  }

  if (props.readOnly) {
    // Both: the first refuses the change, the second stops the element taking
    // the keystroke at all — with only the first, typing shows text the
    // document never took.
    extensions.push(EditorState.readOnly.of(true), EditorView.editable.of(false))
  }

  return extensions
}

function createEditor() {
  if (!editorContainer.value) return

  const state = EditorState.create({
    doc: props.modelValue || '',
    extensions: getExtensions(),
  })

  view = new EditorView({
    state,
    parent: editorContainer.value,
  })
}

function destroyEditor() {
  if (view) {
    view.destroy()
    view = null
  }
}

// Puts text where the cursor is, replacing any selection, with every line after
// the first indented as far as the line it lands on.
function insert(text) {
  if (!view) {
    emit('update:modelValue', (props.modelValue || '') + text)
    return
  }
  const { from, to } = view.state.selection.main
  const indent = /^\s*/.exec(view.state.doc.lineAt(from).text)[0]
  const insertText = text.replace(/\n/g, `\n${indent}`)
  view.dispatch({
    changes: { from, to, insert: insertText },
    selection: { anchor: from + insertText.length },
    scrollIntoView: true,
  })
  view.focus()
}

defineExpose({ insert })

watch(
  () => props.modelValue,
  newVal => {
    if (!view) return
    const currentVal = view.state.doc.toString()
    if (newVal !== currentVal) {
      isUpdating = true
      view.dispatch({
        changes: {
          from: 0,
          to: view.state.doc.length,
          insert: newVal || '',
        },
      })
      isUpdating = false
    }
  },
)

watch(
  () => [props.language, props.readOnly],
  () => {
    // Recreate the editor: both are fixed when it is built
    const currentDoc = view ? view.state.doc.toString() : props.modelValue
    destroyEditor()
    nextTick(() => {
      createEditor()
      if (currentDoc !== props.modelValue) {
        emit('update:modelValue', currentDoc)
      }
    })
  },
)

onMounted(() => {
  createEditor()
})

onBeforeUnmount(() => {
  destroyEditor()
})
</script>

<style scoped>
.code-editor {
  width: 100%;
  overflow: hidden;
  border-radius: var(--p-content-border-radius);
  border: 1px solid var(--p-content-border-color);
}

.code-editor :deep(.cm-editor) {
  height: 100%;
}

.code-editor :deep(.cm-scroller) {
  overflow: auto;
}
</style>
