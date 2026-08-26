import { describe, it, expect, beforeAll, vi } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import CAgentChat from './CAgentChat.vue'
import { makeAgentChatTranslations } from './translations'

beforeAll(() => {
  // PrimeVue overlays bind a matchMedia listener on mount; jsdom has none.
  // @ts-expect-error assigning a stub over a property jsdom does not define
  window.matchMedia = (q: string) => ({
    matches: false,
    media: q,
    onchange: null,
    addEventListener() {},
    removeEventListener() {},
    addListener() {},
    removeListener() {},
    dispatchEvent: () => false,
  })
})

const translations = makeAgentChatTranslations((k: string) => k, 'x.')

function mountChat(agentExec: ReturnType<typeof vi.fn>) {
  setActivePinia(createPinia())

  const api = {
    agentExec,
    agentList: async () => ({ set: [{ agentID: '1', meta: { short: 'A' } }] }),
    aiConversationList: async () => ({ set: [] }),
  }

  return mount(CAgentChat, {
    props: { translations, allowedAgentIDs: ['1'], defaultAgentID: '1', autoResume: false },
    global: { provide: { $SystemAPI: api } },
  })
}

// Every message the store holds for one agent, whichever conversation it is in.
function messagesFor(vm: any, agentID: string): string[] {
  const convs = vm.$.setupState.agentStore.conversations[agentID] || []
  return convs.flatMap((c: any) => (c.messages || []).map((m: any) => String(m.content ?? '')))
}

describe('CAgentChat approvals', () => {
  // A run that stopped to ask has nothing to say yet. The output fallback used
  // to print the whole response object as the agent's answer — a wall of JSON
  // where the reply belongs.
  it('does not print the response object when a run pauses', async () => {
    const w = mountChat(
      vi.fn().mockResolvedValue({
        conversationID: '9',
        output: '',
        status: 'awaiting_approval',
        pendingApproval: { tool: 'compose_record_update', risk: 'write' },
      }),
    )
    await flushPromises()

    const vm = w.vm as any
    await vm.$.setupState.runAgent('1', 'change it')
    await flushPromises()

    expect(vm.$.setupState.pending?.tool).toBe('compose_record_update')

    // Assert on what the store was actually given. The rendered text lags and
    // the filtered view depends on which agent is active, so a check against
    // either passes whether or not the JSON was ever added.
    const contents = messagesFor(vm, '1')
    expect(contents.join('\n')).not.toContain('conversationID')
    expect(contents.join('\n')).not.toContain('awaiting_approval')
  })

  // Approving "for this chat" must be remembered, or every call asks again.
  it('remembers a chat-scoped approval and sends it back', async () => {
    const exec = vi
      .fn()
      .mockResolvedValueOnce({
        conversationID: '9',
        output: '',
        status: 'awaiting_approval',
        pendingApproval: { tool: 'compose_record_update', risk: 'write' },
      })
      .mockResolvedValueOnce({ conversationID: '9', output: 'done', status: 'complete' })

    const w = mountChat(exec)
    await flushPromises()

    const vm = w.vm as any
    await vm.$.setupState.runAgent('1', 'change it')
    await flushPromises()

    await vm.$.setupState.onApprove(true)
    await flushPromises()

    expect(exec).toHaveBeenCalledTimes(2)
    expect(exec.mock.calls[1][0].approvedTools).toEqual(['compose_record_update'])
    expect(vm.$.setupState.pending).toBeNull()
  })

  // "Allow once" sends the approval with the call and does not keep it, so the
  // next use of the same tool asks again.
  it('does not remember a one-off approval', async () => {
    const exec = vi
      .fn()
      .mockResolvedValueOnce({
        conversationID: '9',
        output: '',
        status: 'awaiting_approval',
        pendingApproval: { tool: 'compose_record_update' },
      })
      .mockResolvedValueOnce({ conversationID: '9', output: 'done', status: 'complete' })
      .mockResolvedValueOnce({
        conversationID: '9',
        output: '',
        status: 'awaiting_approval',
        pendingApproval: { tool: 'compose_record_update' },
      })

    const w = mountChat(exec)
    await flushPromises()

    const vm = w.vm as any
    await vm.$.setupState.runAgent('1', 'change it')
    await flushPromises()
    await vm.$.setupState.onApprove(false)
    await flushPromises()

    expect(exec.mock.calls[1][0].approvedTools).toEqual(['compose_record_update'])

    await vm.$.setupState.runAgent('1', 'change it again')
    await flushPromises()
    expect(exec.mock.calls[2][0].approvedTools).toBeUndefined()
  })
})
