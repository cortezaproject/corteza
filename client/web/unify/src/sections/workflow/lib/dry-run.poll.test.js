import { describe, it, expect } from 'vitest'
import { pollOutcome } from './dry-run'

// The dry run used to poll until completedAt appeared. The server sets that on
// completed, failed and canceled only, so a workflow that stopped on a prompt
// or a delay step polled at 1Hz forever.
describe('pollOutcome', () => {
  it('keeps polling a session that is still running', () => {
    expect(pollOutcome({ status: 'started' })).toBe('running')
  })

  it('stops on a session waiting for someone to answer a prompt', () => {
    expect(pollOutcome({ status: 'prompted', completedAt: null })).toBe('suspended')
  })

  it('stops on a session waiting out a delay step', () => {
    expect(pollOutcome({ status: 'suspended', completedAt: null })).toBe('suspended')
  })

  it.each(['completed', 'failed', 'canceled'])('finishes on %s', status => {
    expect(pollOutcome({ status, completedAt: '2026-08-24T14:17:38+02:00' })).toBe('finished')
  })

  // completedAt is the server's own terminal marker, so it outranks a status
  // that has not caught up with it.
  it('finishes whenever completedAt is set, whatever the status says', () => {
    expect(pollOutcome({ status: 'prompted', completedAt: '2026-08-24T14:17:38+02:00' })).toBe(
      'finished',
    )
  })

  it('gives up on a running session once the ceiling passes', () => {
    expect(pollOutcome({ status: 'started' }, { expired: true })).toBe('expired')
  })

  // A suspended session is waiting on purpose and can still be resumed, so the
  // clock must not turn it into an abandoned one.
  it('never expires a suspended session', () => {
    expect(pollOutcome({ status: 'prompted' }, { expired: true })).toBe('suspended')
    expect(pollOutcome({ status: 'suspended' }, { expired: true })).toBe('suspended')
  })

  // An unrecognised or missing status is not a reason to stop; the ceiling
  // still bounds it.
  it('keeps polling an unrecognised status', () => {
    expect(pollOutcome({ status: 'something-new' })).toBe('running')
    expect(pollOutcome({})).toBe('running')
    expect(pollOutcome(undefined)).toBe('running')
  })
})
