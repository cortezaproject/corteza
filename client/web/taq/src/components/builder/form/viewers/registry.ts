import type { Component } from 'vue'
import CViewText from './CViewText.vue'
import CViewNamespace from './CViewNamespace.vue'
import CViewModule from './CViewModule.vue'
import CViewUser from './CViewUser.vue'
import CViewAgent from './CViewAgent.vue'
import CViewSelect from './CViewSelect.vue'
import CViewFieldValueMap from './CViewFieldValueMap.vue'

/**
 * Maps input.type (from API function definition segments) to a read-only viewer component.
 * Mirrors inputs/registry.ts but for preview display instead of editing.
 */
export const VIEWER_REGISTRY: Record<string, Component> = {
  // Agent viewers
  AgentSelector: CViewAgent,
  Agent: CViewAgent,

  // User viewers
  UserSelector: CViewUser,
  User: CViewUser,

  // Namespace viewers
  NamespaceSelector: CViewNamespace,
  Namespace: CViewNamespace,

  // Module viewers
  ModuleSelector: CViewModule,
  Module: CViewModule,

  // Text viewers (default)
  Text: CViewText,
  String: CViewText,
  Number: CViewText,

  // Select/dropdown — handled specially (needs options prop)
  Select: CViewSelect,
  Dropdown: CViewSelect,

  // DateTime — show as text for now
  DateTime: CViewText,

  // Field-value map (aggregate record values)
  FieldValueMap: CViewFieldValueMap,
}

/**
 * Resolves an input type string to its corresponding viewer component.
 * Falls back to CViewText if type is not found.
 */
export function resolveViewerComponent(type: string): Component {
  return VIEWER_REGISTRY[type] || CViewText
}
