import { describe, it, expect } from 'vitest'
import { canReadTrace } from './dry-run'

// A dry run executes the workflow and then polls its session for the trace.
// Execute and session-read are different permissions, so an execute-only user
// used to run the workflow successfully and then be shown a failure toast
// carrying the raw permission message from the first poll.
describe('canReadTrace', () => {
  const allow = () => true
  const deny = () => false

  it('is true only when both permissions are held', () => {
    expect(canReadTrace({ canManageWorkflowSessions: true }, allow)).toBe(true)
  })

  it('is false without sessions.manage on the workflow', () => {
    expect(canReadTrace({ canManageWorkflowSessions: false }, allow)).toBe(false)
  })

  it('is false without sessions.search on the automation component', () => {
    expect(canReadTrace({ canManageWorkflowSessions: true }, deny)).toBe(false)
  })

  it('asks for sessions.search on the automation component', () => {
    const asked = []
    canReadTrace({ canManageWorkflowSessions: true }, (r, o) => {
      asked.push([r, o])
      return true
    })
    expect(asked).toEqual([['automation/', 'sessions.search']])
  })

  // A Workflow model that never declared canManageWorkflowSessions reads
  // undefined here, which must not be mistaken for permission.
  it('is false when the flag is absent from the payload', () => {
    expect(canReadTrace({}, allow)).toBe(false)
    expect(canReadTrace(undefined, allow)).toBe(false)
  })
})
