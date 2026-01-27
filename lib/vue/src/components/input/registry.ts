import type { Component } from 'vue'
import CInputText from './CInputText.vue'
import CInputSelect from './CInputSelect.vue'
import CInputUser from './CInputUser.vue'

/**
 * Maps input.type (from API segments) to Vue component.
 * Extend this registry as new input types are needed.
 */
export const INPUT_REGISTRY: Record<string, Component> = {
  // User selectors
  UserSelector: CInputUser,
  User: CInputUser,

  // Text inputs
  Text: CInputText,
  String: CInputText,
  Number: CInputText,

  // Select/dropdown
  Select: CInputSelect,
  Dropdown: CInputSelect,
}

/**
 * Resolves an input type string to its corresponding Vue component.
 * Falls back to CInputText if type is not found.
 */
export function resolveInputComponent(type: string): Component {
  return INPUT_REGISTRY[type] || CInputText
}
