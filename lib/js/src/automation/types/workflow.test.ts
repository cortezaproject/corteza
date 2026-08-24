import { expect } from 'chai'
import { Workflow } from './workflow'

// The REST payload's can* block is the only thing telling the webapp which
// controls a user may operate. A flag the model does not declare is dropped in
// Apply() and reads false for everyone, so each one the server sends has to
// round-trip.
describe('Workflow permission flags', () => {
  const flags = [
    'canGrant',
    'canUpdateWorkflow',
    'canDeleteWorkflow',
    'canUndeleteWorkflow',
    'canExecuteWorkflow',
    'canManageWorkflowTriggers',
    'canManageWorkflowSessions',
  ] as const

  it('carries every can* flag the REST payload sends', () => {
    const wf = new Workflow(
      Object.fromEntries(flags.map(f => [f, true])) as Record<string, boolean>,
    )

    for (const f of flags) {
      expect(wf[f], f).to.equal(true)
    }
  })

  it('defaults every can* flag to false when the payload omits it', () => {
    const wf = new Workflow({})

    for (const f of flags) {
      expect(wf[f], f).to.equal(false)
    }
  })
})
