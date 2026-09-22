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
  CInputCorredorScript,
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

  // Corredor script selectors (value is the script name)
  CorredorScriptSelector: CInputCorredorScript,
  CorredorScript: CInputCorredorScript,

  // Typed values — a plain text box, which is also where an unmapped type lands
  Text: InputText,
  String: InputText,
  Number: InputText,
  Expression: InputText,

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

/**
 * Whether a type's value is typed in rather than picked from a list — true for
 * every type whose editor is a plain text box, unmapped ones included.
 */
export function isTypedValueInput(type: string): boolean {
  return resolveInputComponent(type) === InputText
}
