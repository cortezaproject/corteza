import type { Component } from 'vue'
import CInputFieldValueMap from './CInputFieldValueMap.vue'
import CInputModule from './CInputModule.vue'
import CInputNamespace from './CInputNamespace.vue'
import CInputRecord from './CInputRecord.vue'
import CInputSelect from './CInputSelect.vue'
import CInputText from './CInputText.vue'
import CInputUser from './CInputUser.vue'

/**
 * Maps input.type (from API segments) to Vue component.
 * Extend this registry as new input types are needed.
 */
export const INPUT_REGISTRY: Record<string, Component> = {
  // User selectors
  UserSelector: CInputUser,
  User: CInputUser,

  // Namespace selectors
  NamespaceSelector: CInputNamespace,
  Namespace: CInputNamespace,

  // Module selectors
  ModuleSelector: CInputModule,
  Module: CInputModule,

  // Record selectors
  RecordSelector: CInputRecord,
  Record: CInputRecord,

  // Text inputs
  Text: CInputText,
  String: CInputText,
  Number: CInputText,

  // Select/dropdown
  Select: CInputSelect,
  Dropdown: CInputSelect,

  // Field-value map (for aggregate record values)
  FieldValueMap: CInputFieldValueMap,
}

/**
 * Resolves an input type string to its corresponding Vue component.
 * Falls back to CInputText if type is not found.
 */
export function resolveInputComponent(type: string): Component {
  return INPUT_REGISTRY[type] || CInputText
}
