import { describe, it, expect, vi } from 'vitest'
import { defineComponent } from 'vue'
import { mount } from '@vue/test-utils'

// A step-config field asks for what its editor can take: a selector asks you to
// select, a text box asks you to type. The corredorExec construct carries one of
// each — `script` is a CorredorScriptSelector, `args` an Expression.

vi.mock('vue-i18n', () => ({
  useI18n: () => ({ t: (k: string, p: Record<string, string>) => `${k}|${p?.field ?? ''}` }),
}))

import { useSegmentForm } from './useSegmentForm'

const CORREDOR_SEGMENTS = [
  {
    meta: {},
    sections: [
      {
        meta: {},
        elements: [
          {
            input: {
              type: 'CorredorScriptSelector',
              label: 'Script',
              argument: 'script',
              required: true,
              context: {},
            },
          },
          { input: { type: 'Expression', label: 'Arguments', argument: 'args', context: {} } },
        ],
      },
    ],
  },
]

const CORREDOR_PARAMETERS = [
  { argumentName: 'script', types: ['String'], required: true },
  { argumentName: 'args', types: ['Vars'] },
]

function placeholders(segments: any[], parameters: any[] = []) {
  let inputs: any[] = []

  const Probe = defineComponent({
    setup() {
      const { processedSegments } = useSegmentForm({
        segments: () => segments,
        parameters: () => parameters,
        getValue: () => null,
        getAggregateValue: () => null,
        onUpdate: () => {},
      })
      inputs = processedSegments.value.flatMap(s => s.sections.flatMap((sec: any) => sec.inputs))
      return {}
    },
    template: '<div/>',
  })

  mount(Probe)
  return Object.fromEntries(inputs.map(i => [i.argument, i.placeholder]))
}

describe('useSegmentForm placeholders', () => {
  it('asks to select a script and to enter its arguments', () => {
    const p = placeholders(CORREDOR_SEGMENTS, CORREDOR_PARAMETERS)

    expect(p.script).toBe('builder.form.selectPlaceholder|Script')
    expect(p.args).toBe('builder.form.enterPlaceholder|Arguments')
  })

  it('asks to enter a plain text field', () => {
    const p = placeholders([
      {
        sections: [
          { elements: [{ input: { type: 'String', label: 'Title', argument: 'title' } }] },
        ],
      },
    ])

    expect(p.title).toBe('builder.form.enterPlaceholder|Title')
  })

  it('keeps the placeholder the construct declares for itself', () => {
    const p = placeholders([
      {
        sections: [
          {
            elements: [
              {
                input: {
                  type: 'Expression',
                  label: 'Arguments',
                  argument: 'args',
                  placeholder: 'e.g. { "id": 1 }',
                },
              },
            ],
          },
        ],
      },
    ])

    expect(p.args).toBe('e.g. { "id": 1 }')
  })

  it('falls back to the argument name when the input carries no label', () => {
    const p = placeholders([
      { sections: [{ elements: [{ input: { type: 'Expression', argument: 'args' } }] }] },
    ])

    expect(p.args).toBe('builder.form.enterPlaceholder|args')
  })
})
