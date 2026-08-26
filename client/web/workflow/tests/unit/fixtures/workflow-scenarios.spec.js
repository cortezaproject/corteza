import { expect } from 'chai'
import fixture from '../../fixtures/workflow-scenarios.json'

describe('Workflow baseline fixtures', () => {
  it('contains a workflow and a dirty editor state', () => {
    expect(fixture.contractVersion).to.equal('1.0.0')
    expect(fixture.fixtureID).to.equal('contract-v1')
    expect(fixture.workflows).to.have.length(1)
    expect(fixture.unsaved.dirty).to.equal(true)
  })

  it('declares all error and empty states', () => {
    expect(Object.keys(fixture.scenarios)).to.have.members([
      'success', 'empty', 'forbidden', 'not-found', 'validation', 'retryable', 'terminal', 'version-conflict',
    ])
    expect(fixture.scenarios.empty.total).to.equal(0)
    expect(fixture.scenarios.forbidden.status).to.equal(403)
    expect(fixture.scenarios['not-found'].status).to.equal(404)
    expect(fixture.scenarios.validation.error).to.equal('VALIDATION_ERROR')
    expect(fixture.scenarios.retryable.retryable).to.equal(true)
    expect(fixture.scenarios.terminal.error).to.equal('OPERATION_FAILED')
    expect(fixture.scenarios['version-conflict'].current_version).to.equal(2)
  })
})
