import { expect } from 'chai'
import fixture from '../../fixtures/record-scenarios.json'

describe('Compose baseline fixtures', () => {
  it('contains a record, attachment and revision without production identifiers', () => {
    expect(fixture.contractVersion).to.equal('1.0.0')
    expect(fixture.fixtureID).to.equal('contract-v1')
    expect(fixture.record.recordID).to.equal('311-fixture-record')
    expect(fixture.attachments).to.have.length(1)
    expect(fixture.revisions).to.have.length(1)
    expect(JSON.stringify(fixture)).to.not.contain('corteza-server-1')
  })

  it('declares all read-only state scenarios', () => {
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
