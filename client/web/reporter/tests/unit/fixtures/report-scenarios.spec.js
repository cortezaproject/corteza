import { expect } from 'chai'
import fixture from '../../fixtures/report-scenarios.json'

describe('Reporter baseline fixtures', () => {
  it('contains a non-empty builder response', () => {
    expect(fixture.contractVersion).to.equal('1.0.0')
    expect(fixture.fixtureID).to.equal('contract-v1')
    expect(fixture.reports).to.have.length(1)
    expect(fixture.builder.rows).to.have.length(1)
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
    expect(fixture.scenarios.terminal.operation_status).to.equal('FAILED')
    expect(fixture.scenarios.terminal.error.error).to.equal('OPERATION_FAILED')
    expect(fixture.scenarios['version-conflict'].current_version).to.equal(2)
  })
})
