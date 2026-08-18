import { describe, expect, it } from 'vitest'

import {
  fromWireArgument,
  fromWireArguments,
  toWireArgument,
  toWireArguments,
  toWireCondition,
  fromWireCondition,
} from '@/sections/taq/utils/reference-binding'

describe('toWireArgument', () => {
  it('turns an identity scope into a dotted expression over the global scope', () => {
    // What the reference panel produces for Invoker → Email. Sent as-is the
    // automation is rejected with scope.unknown and never registers.
    expect(
      toWireArgument({
        argumentName: 'description',
        type: 'String',
        scope: 'invoker',
        expr: 'email',
      }),
    ).toEqual({
      argumentName: 'description',
      type: 'String',
      scope: '',
      expr: 'invoker.email',
      source: undefined,
    })
  })

  it('does the same for runner', () => {
    expect(toWireArgument({ scope: 'runner', expr: 'userID' })).toMatchObject({
      scope: '',
      expr: 'runner.userID',
    })
  })

  it('reads the field from source when that is where it sits', () => {
    expect(toWireArgument({ scope: 'invoker', source: 'name' })).toMatchObject({
      scope: '',
      expr: 'invoker.name',
    })
  })

  it('leaves an ordinary step reference alone', () => {
    const arg = { scope: 'step_2', expr: 'counter' }
    expect(toWireArgument(arg)).toEqual(arg)
  })

  it('leaves a literal alone', () => {
    const arg = { argumentName: 'title', value: 'hello', type: 'String' }
    expect(toWireArgument(arg)).toEqual(arg)
  })

  it('leaves an identity scope with no field alone', () => {
    const arg = { scope: 'invoker' }
    expect(toWireArgument(arg)).toEqual(arg)
  })
})

describe('fromWireArgument', () => {
  it('restores the scope the panel groups by', () => {
    expect(fromWireArgument({ scope: '', expr: 'invoker.email' })).toMatchObject({
      scope: 'invoker',
      expr: 'email',
    })
  })

  it('keeps a dotted path beyond the first field intact', () => {
    expect(fromWireArgument({ scope: '', expr: 'runner.meta.theme' })).toMatchObject({
      scope: 'runner',
      expr: 'meta.theme',
    })
  })

  it('leaves an unrelated expression alone', () => {
    const arg = { scope: '', expr: 'record.values.name' }
    expect(fromWireArgument(arg)).toEqual(arg)
  })

  it('leaves a scoped reference alone', () => {
    const arg = { scope: 'step_2', expr: 'counter' }
    expect(fromWireArgument(arg)).toEqual(arg)
  })
})

describe('round trip', () => {
  it('survives builder → wire → builder unchanged', () => {
    const args = [
      { argumentName: 'recipient', value: '123', type: 'ID' },
      { argumentName: 'description', scope: 'invoker', expr: 'email', type: 'String' },
      { argumentName: 'title', scope: 'step_2', expr: 'counter', type: 'String' },
    ]
    const back = fromWireArguments(toWireArguments(args))
    expect(back).toMatchObject([
      { argumentName: 'recipient', value: '123' },
      { argumentName: 'description', scope: 'invoker', expr: 'email' },
      { argumentName: 'title', scope: 'step_2', expr: 'counter' },
    ])
  })

  it('keeps the scope truthy so the builder still reads it as a reference', () => {
    // Every "is this a reference" test in the forms is a truthiness check on
    // scope; an identity binding that came back with scope '' would render as
    // a literal and lose its chip.
    const [arg] = fromWireArguments(toWireArguments([{ scope: 'invoker', expr: 'email' }]))
    expect(arg.scope).toBeTruthy()
  })

  it('tolerates undefined', () => {
    expect(toWireArguments(undefined)).toEqual([])
    expect(fromWireArguments(undefined)).toEqual([])
  })
})

describe('conditions', () => {
  const lit = { value: { '@type': 'String', '@value': 'x' } }

  it('turns an identity scope into a dotted symbol on the global scope', () => {
    // The evaluator defaults to the global scope and splits the symbol on '.',
    // so this is the shape a condition can actually resolve.
    const built = { ref: 'eq', args: [{ symbol: 'email', meta: { scope: 'invoker' } }, lit] }
    expect(toWireCondition(built)).toEqual({
      ref: 'eq',
      args: [{ symbol: 'invoker.email', meta: { scope: 'global' } }, lit],
    })
  })

  it('restores it on the way back', () => {
    const wire = { ref: 'eq', args: [{ symbol: 'runner.name', meta: { scope: 'global' } }, lit] }
    expect(fromWireCondition(wire)).toEqual({
      ref: 'eq',
      args: [{ symbol: 'name', meta: { scope: 'runner' } }, lit],
    })
  })

  it('reaches into and/or groups', () => {
    const built = {
      ref: 'and',
      args: [
        { ref: 'eq', args: [{ symbol: 'email', meta: { scope: 'invoker' } }, lit] },
        { ref: 'eq', args: [{ symbol: 'counter', meta: { scope: 'step_2' } }, lit] },
      ],
    }
    const wire = toWireCondition(built) as any
    expect(wire.args[0].args[0]).toEqual({ symbol: 'invoker.email', meta: { scope: 'global' } })
    expect(wire.args[1].args[0]).toEqual({ symbol: 'counter', meta: { scope: 'step_2' } })
    expect(fromWireCondition(wire)).toEqual(built)
  })

  it('leaves an ordinary step reference and a literal alone', () => {
    const built = { ref: 'eq', args: [{ symbol: 'counter', meta: { scope: 'step_2' } }, lit] }
    expect(toWireCondition(built)).toEqual(built)
    expect(fromWireCondition(built)).toEqual(built)
  })

  it('tolerates null', () => {
    expect(toWireCondition(null)).toBeNull()
    expect(fromWireCondition(null)).toBeNull()
  })
})
