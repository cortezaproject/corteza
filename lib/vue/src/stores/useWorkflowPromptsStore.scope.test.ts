import { describe, it, expect, vi } from 'vitest'
import { createPinia, setActivePinia } from 'pinia'
import { createApp } from 'vue'
import { useWorkflowPromptsStore } from './useWorkflowPromptsStore'

// The prompt list is the whole user's, not this tab's. A kind this webapp
// cannot render used to be dropped outright, which left its session waiting on
// a dialog no screen here would ever show. Resuming past it is right only for a
// session this screen started — another tab may be the one meant to handle it.

function withAPI() {
  const sessionResumeState = vi.fn().mockResolvedValue({})
  const app = createApp({})
  app.provide('$AutomationAPI', { sessionResumeState })
  const pinia = createPinia()
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  ;(pinia as any)._a = app
  app.use(pinia)
  setActivePinia(pinia)
  return { store: useWorkflowPromptsStore(), sessionResumeState }
}

// refetchRecords is webapps:['compose']; alert has no restriction.
// Ids go through the HumanID cast, so they must look like ids.
const S1 = '111111111111111111'
const X1 = '222222222222222222'
const X2 = '333333333333333333'
const prompt = (ref: string, sessionID: string, stateID: string) =>
  ({ ref, sessionID, stateID }) as never

describe('prompts this webapp cannot render', () => {
  it('keeps a kind the webapp does render', () => {
    const { store } = withAPI()

    store.newPrompt(prompt('alert', S1, X1), 'workflow')
    expect(store.prompts).toHaveLength(1)
  })

  it('does not resume a foreign kind for a session it did not start', async () => {
    const { store, sessionResumeState } = withAPI()

    store.newPrompt(prompt('refetchRecords', S1, X1), 'workflow')

    expect(store.prompts).toHaveLength(0)
    expect(store.skipped).toHaveLength(0)
    expect(sessionResumeState).not.toHaveBeenCalled()
  })

  it('resumes past a foreign kind for a session it started', async () => {
    const { store, sessionResumeState } = withAPI()
    store.ownSession(S1)

    store.newPrompt(prompt('refetchRecords', S1, X1), 'workflow')
    await Promise.resolve()

    expect(store.prompts).toHaveLength(0)
    expect(store.skipped).toHaveLength(1)
    expect(sessionResumeState).toHaveBeenCalledWith({
      sessionID: S1,
      stateID: X1,
      input: {},
    })
  })

  // The realtime prompt push beats the exec response that names the session, so
  // this ordering is the normal one, not the edge case.
  it('steps over a prompt that arrived before the session was claimed', async () => {
    const { store, sessionResumeState } = withAPI()

    store.newPrompt(prompt('refetchRecords', S1, X1), 'workflow')
    expect(sessionResumeState).not.toHaveBeenCalled()

    store.ownSession(S1)
    await Promise.resolve()

    expect(store.skipped).toHaveLength(1)
    expect(sessionResumeState).toHaveBeenCalledWith({
      sessionID: S1,
      stateID: X1,
      input: {},
    })
  })

  it('forgets a session when it is disowned', async () => {
    const { store, sessionResumeState } = withAPI()
    store.ownSession(S1)
    store.newPrompt(prompt('refetchRecords', S1, X1), 'workflow')
    await Promise.resolve()
    expect(store.skipped).toHaveLength(1)

    store.disownSession(S1)
    expect(store.skipped).toHaveLength(0)

    store.newPrompt(prompt('refetchRecords', S1, X2), 'workflow')
    expect(sessionResumeState).toHaveBeenCalledTimes(1)
  })

  it('still renders a foreign kind in the webapp that owns it', () => {
    const { store } = withAPI()
    store.newPrompt(prompt('refetchRecords', S1, X1), 'compose')
    expect(store.prompts).toHaveLength(1)
  })
})
