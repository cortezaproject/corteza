import { expect } from 'chai'
import { UIHooks } from './ui-hooks'

describe(__filename, () => {
  it('should convert script to button', () => {
    const hooks = new UIHooks('test')

    hooks.Register(
      {
        name: 'scriptName',
        label: 'btnLabel',
        triggers: [{
          resourceTypes: ['system'],
          eventTypes: ['onManual'],
          weight: 42,
          constraints: [],
          uiProps: [
            { name: 'variant', value: 'danger' },
            { name: 'app', value: 'test' },
            { name: 'page', value: 'index' },
            { name: 'slot', value: 'header' },
          ],
        }],
      },
    )

    const bb = hooks.Find('system', 'index', 'header')

    expect(bb).to.have.lengthOf(1)
    const b = bb[0]
    expect(b).to.have.property('label').equal('btnLabel')
    expect(b).to.have.property('script').equal('scriptName')
    expect(b).to.have.property('resourceType').equal('system')
    expect(b).to.have.property('slot').equal('header')
    expect(b).to.have.property('variant').equal('danger')
  })

  it('should remove existing scripts', () => {
    const hooks = new UIHooks('test')
    const s = {
      name: 'scriptName',
      label: 'btnLabel',
      triggers: [{
        resourceTypes: ['system'],
        eventTypes: ['onManual'],
        constraints: [],
        uiProps: [
          { name: 'app', value: 'test' },
          { name: 'page', value: 'index' },
          { name: 'slot', value: 'header' },
        ],
      }],
    }

    let bb

    hooks.Register(s)
    bb = hooks.Find('system', 'index', 'header')
    expect(bb).to.have.lengthOf(1)
    hooks.Register(s)
    bb = hooks.Find('system', 'index', 'header')
    expect(bb).to.have.lengthOf(1)
    hooks.Register(s, { ...s, name: 'anotherScript' })
    bb = hooks.Find('system', 'index', 'header')
    expect(bb).to.have.lengthOf(2)
    hooks.Register(s)
    bb = hooks.Find('system', 'index', 'header')
    expect(bb).to.have.lengthOf(2)
  })
})
