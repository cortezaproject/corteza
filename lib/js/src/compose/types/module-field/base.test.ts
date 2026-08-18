import { expect } from 'chai'
import { ModuleField } from './base'
import { ModuleFieldBool } from './bool'
import { NoID } from '../../../cast'

describe('check module field casting', () => {
  it('simple assignment', () => {
    const f = new ModuleField({
      name: 'fname',
      kind: 'number',
    })

    expect(f.name).to.equal('fname')
  })
})

describe('required capability', () => {
  it('a bool field cannot be marked required', () => {
    const f = new ModuleFieldBool({ name: 'agree' })
    expect(f.cap.required).to.equal(false)
  })

  it('a new bool field drops an incoming required flag', () => {
    const f = new ModuleFieldBool({ name: 'agree', isRequired: true })
    expect(f.isRequired).to.equal(false)
  })

  it('an existing bool field keeps the flag it was saved with', () => {
    const f = new ModuleFieldBool({ fieldID: '12345', name: 'agree', isRequired: true })
    expect(f.fieldID).to.not.equal(NoID)
    expect(f.isRequired).to.equal(true)
  })

  it('other kinds are still requirable', () => {
    const f = new ModuleField({ name: 'title', kind: 'String', isRequired: true })
    expect(f.isRequired).to.equal(true)
  })
})
