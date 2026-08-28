import { describe, it, expect, vi, afterEach } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'

// A multi-value field is stored as one Expr per value, all sharing the target:
// that repetition is what the server folds into a KVV (see the `auxVals` loop in
// server/compose/automation/ng_record_handlers.go). Collapsing them into a
// single Expr loses every value but the last.

vi.mock('vue-i18n', () => ({ useI18n: () => ({ t: k => k }) }))

vi.mock('./DynamicForm.vue', () => ({
  default: {
    name: 'DynamicForm',
    props: { processedSegments: { type: Array, default: () => [] } },
    emits: ['update:value'],
    template: '<div />',
  },
}))

import FunctionForm from './FunctionForm.vue'

const functionDef = {
  parameters: [
    { argumentName: 'namespace', types: ['ID'] },
    { argumentName: 'module', types: ['ID'] },
    { argumentName: 'values', types: ['FieldValueMap'], aggregate: true },
  ],
  segments: [
    {
      sections: [
        {
          elements: [
            { input: { argument: 'namespace', type: 'Namespace', label: 'Namespace' } },
            { input: { argument: 'module', type: 'Module', label: 'Module' } },
            { input: { argument: 'values', type: 'FieldValueMap', label: 'Values' } },
          ],
        },
      ],
    },
  ],
}

const baseArgs = [
  { argumentName: 'namespace', target: '', value: '100', type: 'ID' },
  { argumentName: 'module', target: '', value: '200', type: 'ID' },
]

let wrapper

async function mountForm(args) {
  wrapper = mount(FunctionForm, {
    props: { functionDef, arguments: args, upstreamResults: [], nodes: [] },
  })
  await flushPromises()
  return wrapper
}

function valuesInput() {
  const segments = wrapper.findComponent({ name: 'DynamicForm' }).props('processedSegments')
  return segments[0].sections[0].inputs.find(i => i.argument === 'values')
}

function emitValues(value) {
  return wrapper.findComponent({ name: 'DynamicForm' }).vm.$emit('update:value', 'values', value)
}

function lastArgs() {
  return wrapper.emitted('update:arguments').at(-1)[0]
}

afterEach(() => {
  wrapper?.unmount()
  wrapper = null
})

describe('function form — multi-value aggregate rows', () => {
  it('reads repeated targets back as an array', async () => {
    await mountForm([
      ...baseArgs,
      { argumentName: 'values', target: 'tags', value: 'a', type: 'String' },
      { argumentName: 'values', target: 'tags', value: 'b', type: 'String' },
      { argumentName: 'values', target: 'subject', value: 'hi', type: 'String' },
    ])

    expect(valuesInput().value).toEqual({
      tags: { value: ['a', 'b'] },
      subject: { value: 'hi' },
    })
  })

  it('leaves a target seen once as a plain value', async () => {
    await mountForm([
      ...baseArgs,
      { argumentName: 'values', target: 'subject', value: 'hi', type: 'String' },
    ])

    expect(valuesInput().value).toEqual({ subject: { value: 'hi' } })
  })

  it('writes one Expr per value, all sharing the target', async () => {
    await mountForm([...baseArgs])

    await emitValues({ tags: { value: ['a', 'b'] } })

    expect(lastArgs().filter(a => a.argumentName === 'values')).toEqual([
      { argumentName: 'values', target: 'tags', type: 'String', value: 'a' },
      { argumentName: 'values', target: 'tags', type: 'String', value: 'b' },
    ])
  })

  it('round-trips a multi-value row through the shape the server reads', async () => {
    await mountForm([...baseArgs])

    await emitValues({ tags: { value: ['a', 'b', 'c'] } })
    const saved = lastArgs()

    // What lands in storage has to be three Exprs, not one holding an array:
    // the server casts each Expr to a string, so an array in one slot is lost.
    expect(saved.filter(a => a.argumentName === 'values')).toHaveLength(3)
    expect(saved.every(a => !Array.isArray(a.value))).toBe(true)

    await mountForm(saved)
    expect(valuesInput().value).toEqual({ tags: { value: ['a', 'b', 'c'] } })
  })

  it('keeps an emptied list on screen as one blank Expr', async () => {
    await mountForm([...baseArgs])

    await emitValues({ tags: { value: [] } })

    expect(lastArgs().filter(a => a.argumentName === 'values')).toEqual([
      { argumentName: 'values', target: 'tags', type: 'String', value: '' },
    ])
  })

  it('still writes a single-value row as one Expr', async () => {
    await mountForm([...baseArgs])

    await emitValues({ subject: { value: 'hi' } })

    expect(lastArgs().filter(a => a.argumentName === 'values')).toEqual([
      { argumentName: 'values', target: 'subject', type: 'String', value: 'hi' },
    ])
  })

  it('leaves namespace and module in the first two slots the server reads', async () => {
    await mountForm([...baseArgs])

    await emitValues({ tags: { value: ['a', 'b'] } })

    const args = lastArgs()
    expect(args[0].argumentName).toBe('namespace')
    expect(args[1].argumentName).toBe('module')
  })

  it('does not fold a reference row into the array', async () => {
    await mountForm([
      ...baseArgs,
      { argumentName: 'values', target: 'owner', scope: 'trigger', expr: 'record.ownedBy' },
      { argumentName: 'values', target: 'tags', value: 'a', type: 'String' },
      { argumentName: 'values', target: 'tags', value: 'b', type: 'String' },
    ])

    expect(valuesInput().value).toEqual({
      owner: { value: '', scope: 'trigger', source: 'record.ownedBy' },
      tags: { value: ['a', 'b'] },
    })
  })
})
