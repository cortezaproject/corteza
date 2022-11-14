import { expect } from 'chai'
import { DummyHandler, Handler } from './handlers'

describe('handler', () => {
  it('should match compatible event', () => {
    const h = new Handler(DummyHandler, { eventTypes: ['e1'], resourceTypes: ['r1'] })
    expect(h.Match({ resourceType: 'r1', eventType: 'e1' })).to.equal(true)
  })

  it('should not match incompatible events', () => {
    const h = new Handler(DummyHandler, { eventTypes: ['e1'], resourceTypes: ['r1'] })
    expect(h.Match({ resourceType: 'r2', eventType: 'e1' })).to.equal(false)
    expect(h.Match({ resourceType: 'r1', eventType: 'e2' })).to.equal(false)
    expect(h.Match({ resourceType: 'r2', eventType: 'e2' })).to.equal(false)
  })
})
