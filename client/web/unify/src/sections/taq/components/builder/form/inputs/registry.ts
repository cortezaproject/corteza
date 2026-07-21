import type { Component } from 'vue'
import InputText from 'primevue/inputtext'
import {
  CInputSelect,
  CInputUser,
  CInputNamespace,
  CInputModule,
  CInputRecord,
  CInputDateTime,
  CInputAgent,
  CInputSwitch,
  CInputCron,
  CInputWorkflow,
} from '@planetcrust/human-vue/src/components/input'
import CInputFieldValueMap from './CInputFieldValueMap.vue'
import CInputWorkflowInputMap from './CInputWorkflowInputMap.vue'
import CInputArray from './CInputArray.vue'

/**
 * Maps input.type (from API function definition segments) to Vue component.
 * Used by DynamicInput to resolve automation step config inputs.
 */
export const INPUT_REGISTRY: Record<string, Component> = {
  // Agent selectors
  AgentSelector: CInputAgent,
  Agent: CInputAgent,

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

  // Workflow selectors
  WorkflowSelector: CInputWorkflow,
  Workflow: CInputWorkflow,

  // Text inputs
  Text: InputText,
  String: InputText,
  Number: InputText,

  // Array/List inputs
  Array: CInputArray,

  // Booleans
  Boolean: CInputSwitch,

  // Select/dropdown
  Select: CInputSelect,
  Dropdown: CInputSelect,

  // Date/Time/Cron
  DateTime: CInputDateTime,
  Cron: CInputCron,
  Interval: CInputCron,

  // Field-value map (for aggregate record values)
  FieldValueMap: CInputFieldValueMap,

  // Workflow input map (one referenceable field per declared workflow input)
  WorkflowInputMap: CInputWorkflowInputMap,
}

/**
 * Resolves an input type string to its corresponding Vue component.
 * Falls back to InputText if type is not found.
 */
export function resolveInputComponent(type: string): Component {
  return INPUT_REGISTRY[type] || InputText
}
