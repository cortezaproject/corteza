import type { Component } from 'vue'
import { defineAsyncComponent } from 'vue'

type FieldEntry = {
  editor: Component
  viewer: Component
}

const CFieldStringEditor = defineAsyncComponent(() => import('./editors/CFieldStringEditor.vue'))
const CFieldStringViewer = defineAsyncComponent(() => import('./viewers/CFieldStringViewer.vue'))

export const FIELD_REGISTRY: Record<string, FieldEntry> = {
  String: {
    editor: CFieldStringEditor,
    viewer: CFieldStringViewer,
  },
  Number: {
    editor: defineAsyncComponent(() => import('./editors/CFieldNumberEditor.vue')),
    viewer: defineAsyncComponent(() => import('./viewers/CFieldNumberViewer.vue')),
  },
  Bool: {
    editor: defineAsyncComponent(() => import('./editors/CFieldBoolEditor.vue')),
    viewer: defineAsyncComponent(() => import('./viewers/CFieldBoolViewer.vue')),
  },
  DateTime: {
    editor: defineAsyncComponent(() => import('./editors/CFieldDateTimeEditor.vue')),
    viewer: defineAsyncComponent(() => import('./viewers/CFieldDateTimeViewer.vue')),
  },
  Select: {
    editor: defineAsyncComponent(() => import('./editors/CFieldSelectEditor.vue')),
    viewer: defineAsyncComponent(() => import('./viewers/CFieldSelectViewer.vue')),
  },
  Email: {
    editor: CFieldStringEditor,
    viewer: defineAsyncComponent(() => import('./viewers/CFieldEmailViewer.vue')),
  },
  Url: {
    editor: CFieldStringEditor,
    viewer: defineAsyncComponent(() => import('./viewers/CFieldUrlViewer.vue')),
  },
}

export function resolveFieldEditor(kind: string): Component {
  return FIELD_REGISTRY[kind]?.editor || CFieldStringEditor
}

export function resolveFieldViewer(kind: string): Component {
  return FIELD_REGISTRY[kind]?.viewer || CFieldStringViewer
}
