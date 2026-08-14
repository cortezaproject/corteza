<template>
  <div
    ref="host"
    class="c-expression"
    :class="{ 'c-expression--disabled': disabled }"
    :data-dialect="dialect"
  />
</template>

<script setup>
import { onBeforeUnmount, onMounted, ref, watch } from 'vue'
import { Compartment, EditorState, Prec, RangeSetBuilder } from '@codemirror/state'
import {
  Decoration,
  EditorView,
  ViewPlugin,
  keymap,
  placeholder as cmPlaceholder,
} from '@codemirror/view'
import { defaultKeymap, history, historyKeymap } from '@codemirror/commands'
import {
  autocompletion,
  closeCompletion,
  completionKeymap,
  completionStatus,
  startCompletion,
} from '@codemirror/autocomplete'
import { linter } from '@codemirror/lint'
import { completionAt, lintExpression, tokenRanges } from './syntax'

const props = defineProps({
  modelValue: { type: String, default: '' },
  // 'ql' for a record filter, 'interpolation' for free text carrying ${} holes,
  // 'expr' for a server-evaluated expression (visibility, field conditions).
  dialect: { type: String, default: 'ql' },
  // ScopeEntry[] from buildScope() — what ${...} may reference here.
  scope: { type: Array, default: () => [] },
  // Fields of the module being queried; offered as bare QL identifiers.
  queryFields: { type: Array, default: () => [] },
  placeholder: { type: String, default: '' },
  disabled: { type: Boolean, default: false },
  // Opens this tall, so a filter has room to read as one.
  minLines: { type: Number, default: 2 },
  // Grows with the content up to this many lines, then scrolls.
  maxLines: { type: Number, default: 8 },
})

const emit = defineEmits(['update:modelValue'])

const host = ref(null)
const context = new Compartment()
let view = null
let applying = false

const holeMark = Decoration.mark({ class: 'c-expression__hole' })
const marks = {
  hole: holeMark,
  keyword: Decoration.mark({ class: 'c-expression__keyword' }),
  function: Decoration.mark({ class: 'c-expression__function' }),
  string: Decoration.mark({ class: 'c-expression__string' }),
  number: Decoration.mark({ class: 'c-expression__number' }),
}

function buildDecorations(state) {
  const builder = new RangeSetBuilder()
  for (const r of tokenRanges(state.doc.toString(), props.dialect)) {
    builder.add(r.from, r.to, marks[r.kind] || holeMark)
  }
  return builder.finish()
}

const highlighter = ViewPlugin.fromClass(
  class {
    constructor(v) {
      this.decorations = buildDecorations(v.state)
    }
    update(u) {
      if (u.docChanged || u.viewportChanged) this.decorations = buildDecorations(u.state)
    }
  },
  { decorations: v => v.decorations },
)

// Completion and lint both close over the current scope, so they are grouped
// in a compartment and swapped whole when the scope changes — a module still
// loading must not leave the editor checking against an empty field list.
function contextExtensions() {
  return [
    autocompletion({
      override: [
        cx => {
          const res = completionAt(
            cx.state.doc.toString(),
            cx.pos,
            props.scope,
            props.dialect,
            props.queryFields,
            cx.explicit,
          )
          if (!res || !res.options.length) return null
          return {
            from: res.from,
            to: res.to,
            // completionAt has already filtered by prefix and ordered by group.
            // CodeMirror's own filter would re-match labels against the typed
            // range, which for a bare `$` contains no letter any label shares —
            // so it would discard every option.
            filter: false,
            options: res.options.map(o => ({
              label: o.label,
              detail: o.detail,
              boost: o.boost,
              apply: (view, _c, from, to) => {
                const insert = o.insert ?? o.label
                view.dispatch({
                  changes: { from, to, insert },
                  selection: { anchor: from + (o.cursor ?? insert.length) },
                })
                if (o.retrigger) startCompletion(view)
              },
            })),
          }
        },
      ],
      activateOnTyping: true,
      icons: false,
    }),
    linter(v =>
      lintExpression(v.state.doc.toString(), props.scope, props.dialect).map(d => ({
        from: d.from,
        to: d.to,
        severity: d.severity,
        message: d.message,
      })),
    ),
  ]
}

function extensions() {
  return [
    history(),
    highlighter,
    EditorView.lineWrapping,
    cmPlaceholder(props.placeholder || ''),
    // Escape dismisses the suggestion list and stops there. CodeMirror's keymap
    // calls preventDefault but lets the event bubble, and these inputs sit
    // inside a dialog that closes on Escape — so without this, dismissing the
    // suggestions also discards the block being configured.
    //
    // Must outrank completionKeymap: that binding closes the completion first,
    // leaving completionStatus null by the time a lower-precedence handler asks
    // whether there was anything to dismiss.
    Prec.highest(
      EditorView.domEventHandlers({
        keydown(event, v) {
          if (event.key !== 'Escape' || !completionStatus(v.state)) return false
          event.stopPropagation()
          closeCompletion(v)
          return true
        },
      }),
    ),
    keymap.of([...completionKeymap, ...historyKeymap, ...defaultKeymap]),
    context.of(contextExtensions()),
    EditorState.readOnly.of(props.disabled),
    EditorView.editable.of(!props.disabled),
    EditorView.updateListener.of(u => {
      if (!u.docChanged || applying) return
      emit('update:modelValue', u.state.doc.toString())
    }),
    EditorView.theme({
      '&': {
        fontSize: '13px',
        fontFamily: 'ui-monospace, SFMono-Regular, "SF Mono", Menlo, Consolas, monospace',
      },
      '&.cm-focused': { outline: 'none' },
      '.cm-line': { padding: '0' },
      '.cm-scroller': { lineHeight: '1.5', maxHeight: `calc(${props.maxLines} * 1.5 * 13px)` },
      '.cm-content': { padding: '0', minHeight: `calc(${props.minLines} * 1.5 * 13px)` },
    }),
  ]
}

onMounted(() => {
  const doc = props.modelValue || ''
  view = new EditorView({
    // Cursor parked at the end, so a hint chip clicked before the editor has
    // been focused appends rather than landing in front of the expression.
    state: EditorState.create({
      doc,
      selection: { anchor: doc.length },
      extensions: extensions(),
    }),
    parent: host.value,
  })
})

onBeforeUnmount(() => {
  view?.destroy()
  view = null
})

watch(
  () => props.modelValue,
  v => {
    const next = v || ''
    if (!view || next === view.state.doc.toString()) return
    applying = true
    view.dispatch({ changes: { from: 0, to: view.state.doc.length, insert: next } })
    applying = false
  },
)

watch(
  () => [props.scope, props.queryFields, props.dialect],
  () => view?.dispatch({ effects: context.reconfigure(contextExtensions()) }),
  { deep: true },
)

// Used by CExpressionHint to drop a variable in at the cursor.
function insert(text) {
  if (!view) return
  const { from, to } = view.state.selection.main
  view.dispatch({ changes: { from, to, insert: text }, selection: { anchor: from + text.length } })
  view.focus()
}

defineExpose({ insert, focus: () => view?.focus() })
</script>

<style scoped>
.c-expression {
  width: 100%;
  padding: 0.5rem 0.75rem;
  background: var(--p-content-background);
  border: 1px solid var(--p-content-border-color);
  border-radius: var(--p-content-border-radius, 6px);
  color: var(--p-text-color);
  transition: border-color 0.2s;
}

.c-expression:focus-within {
  border-color: var(--p-primary-color);
}

.c-expression--disabled {
  background: var(--p-content-hover-background);
  opacity: 0.6;
}

.c-expression :deep(.c-expression__hole) {
  background: color-mix(in srgb, var(--p-primary-color) 14%, transparent);
  border-radius: 3px;
  color: var(--p-primary-color);
}

.c-expression :deep(.c-expression__keyword) {
  color: var(--p-primary-color);
  font-weight: 600;
}

.c-expression :deep(.c-expression__string) {
  color: var(--p-green-500, #22c55e);
}

.c-expression :deep(.c-expression__function) {
  color: var(--p-purple-500, #a855f7);
}

.c-expression :deep(.c-expression__number) {
  color: var(--p-orange-500, #f97316);
}
</style>
