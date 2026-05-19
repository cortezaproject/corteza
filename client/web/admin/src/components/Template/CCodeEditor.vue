<template>
  <div ref="editorContainer" class="code-editor" />
</template>

<script setup>
import { ref, watch, onMounted, onBeforeUnmount, nextTick } from 'vue'
import { EditorState } from '@codemirror/state'
import { EditorView, keymap, lineNumbers, highlightActiveLine, highlightActiveLineGutter } from '@codemirror/view'
import { defaultKeymap, indentWithTab } from '@codemirror/commands'
import { html } from '@codemirror/lang-html'
import { json } from '@codemirror/lang-json'
import { oneDark } from '@codemirror/theme-one-dark'
import { syntaxHighlighting, defaultHighlightStyle, bracketMatching, foldGutter } from '@codemirror/language'

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
    keymap.of([...defaultKeymap, indentWithTab]),
    EditorView.lineWrapping,
    EditorView.updateListener.of((update) => {
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

  if (props.dark) {
    extensions.push(oneDark)
  }

  if (props.readOnly) {
    extensions.push(EditorState.readOnly.of(true))
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

watch(() => props.modelValue, (newVal) => {
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
})

watch(() => props.language, () => {
  // Recreate editor when language changes
  const currentDoc = view ? view.state.doc.toString() : props.modelValue
  destroyEditor()
  nextTick(() => {
    createEditor()
    if (currentDoc !== props.modelValue) {
      emit('update:modelValue', currentDoc)
    }
  })
})

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
