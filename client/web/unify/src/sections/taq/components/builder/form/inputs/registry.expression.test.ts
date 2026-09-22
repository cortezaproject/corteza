import { describe, it, expect } from 'vitest'
import InputText from 'primevue/inputtext'
import { resolveInputComponent, isTypedValueInput } from './registry'
import { resolveViewerComponent } from '../viewers/registry'
import CViewText from '../viewers/CViewText.vue'

// `Expression` is the construct library's free-typed value: the corredorExec
// step's `args`, and every usersCreate/usersUpdate scalar. It is edited and
// shown as text, and it is a typed value rather than something picked.
const SELECTORS = [
  'CorredorScriptSelector',
  'UserSelector',
  'NamespaceSelector',
  'ModuleSelector',
  'RecordSelector',
  'WorkflowSelector',
  'AgentSelector',
  'Select',
  'Dropdown',
  'DateTime',
  'Cron',
]

const TYPED_VALUES = ['Expression', 'Text', 'String', 'Number']

describe('TAQ registries — Expression', () => {
  it('edits an Expression input with a text box', () => {
    expect(resolveInputComponent('Expression')).toBe(InputText)
  })

  it('shows an Expression value as text', () => {
    expect(resolveViewerComponent('Expression')).toBe(CViewText)
  })

  it.each(TYPED_VALUES)('counts %s as a typed value', type => {
    expect(isTypedValueInput(type)).toBe(true)
  })

  it('counts a type it does not know as a typed value, the way it renders one', () => {
    expect(isTypedValueInput('SomethingTheServerAddedLater')).toBe(true)
  })

  it.each(SELECTORS)('does not count %s as a typed value', type => {
    expect(isTypedValueInput(type)).toBe(false)
  })
})
