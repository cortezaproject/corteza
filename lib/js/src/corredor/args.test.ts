/* eslint-disable @typescript-eslint/ban-ts-ignore,no-unused-expressions */

import { describe, it } from 'mocha'
import { expect } from 'chai'
import { Args } from './args'
import { CortezaTypes } from './args-corteza'
import { Module } from '../compose/types/module'

describe(__filename, () => {
  it('should provide getter for a given arg', () => {
    const args = new Args({ foo: 'bar' })
    expect(args).to.haveOwnProperty('foo')
    expect(args).property('foo').eq('bar')
  })

  it('should not have undefined property', () => {
    const args = new Args({})
    expect(args).not.to.haveOwnProperty('foo')
  })

  it('should use caster', () => {
    const module = { moduleID: '42' }
    const args = new Args({ module }, CortezaTypes)
    expect(args).to.haveOwnProperty('$module')
    expect(args).to.haveOwnProperty('rawModule')
    expect(args).property('$module').instanceOf(Module)
    expect(args).property('rawModule').to.deep.eq(module)
  })

  it('should properly handle pre-casted variables', () => {
    const module = new Module({ moduleID: '42' })
    const args = new Args({ module }, CortezaTypes)
    expect(args).to.haveOwnProperty('$module')
    expect(args).to.haveOwnProperty('rawModule')
    expect(args).property('$module').instanceOf(Module)
    expect(args).property('rawModule').to.deep.eq(module)
  })
})
