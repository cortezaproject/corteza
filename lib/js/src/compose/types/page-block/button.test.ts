import { describe, it } from 'mocha'
import { expect } from 'chai'
import { PageBlockAutomation } from './automation'
import { PageBlockRecordList } from './record-list'

// A button the editor configured must come back the same after a round trip:
// the editor tells a TAQ button apart by the trigger it names, and offers a
// trigger it cannot see on a button as if it were still free.
describe('automation button round trip', () => {
  const configured = {
    label: 'Run it',
    variant: 'primary',
    resourceType: 'compose:record',
    automationID: '42',
    triggerHandle: 'onDemand',
    scriptType: 'taq',
    enabled: true,
  }

  it('keeps the trigger handle and the kind through the block', () => {
    const block = new PageBlockAutomation({ options: { buttons: [configured] } })
    const [button] = block.options.buttons

    expect(button.triggerHandle).to.equal('onDemand')
    expect(button.scriptType).to.equal('taq')
    expect(button.automationID).to.equal('42')
  })

  it('keeps them through a record list block and its JSON', () => {
    const block = new PageBlockRecordList({ options: { selectionButtons: [configured] } })
    const [button] = JSON.parse(JSON.stringify(block)).options.selectionButtons

    expect(button.triggerHandle).to.equal('onDemand')
    expect(button.scriptType).to.equal('taq')
  })

  it('leaves both unset on a button that names neither', () => {
    const block = new PageBlockAutomation({
      options: { buttons: [{ label: 'Script', script: '/server-scripts/x.js:default' }] },
    })
    const [button] = block.options.buttons

    expect(button.triggerHandle).to.equal(undefined)
    expect(button.scriptType).to.equal(undefined)
  })
})

describe('automation button visibility', () => {
  it('keeps the condition through the block and its JSON', () => {
    const block = new PageBlockAutomation({
      options: { buttons: [{ label: 'Send', visibility: { expression: 'record.values.ok' } }] },
    })
    const [button] = JSON.parse(JSON.stringify(block)).options.buttons

    expect(button.visibility).to.deep.equal({ expression: 'record.values.ok' })
  })

  it('carries no condition when none or an empty one is given', () => {
    const block = new PageBlockAutomation({
      options: {
        buttons: [{ label: 'Plain' }, { label: 'Blank', visibility: { expression: '' } }],
      },
    })
    const [plain, blank] = JSON.parse(JSON.stringify(block)).options.buttons

    expect(plain).to.not.have.property('visibility')
    expect(blank).to.not.have.property('visibility')
  })
})
