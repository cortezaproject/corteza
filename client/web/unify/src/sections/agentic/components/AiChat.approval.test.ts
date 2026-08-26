import { describe, it, expect, beforeAll, vi } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import AiChat from './AiChat.vue'

vi.mock('vue-i18n', () => ({ useI18n: () => ({ t: (k: string) => k }) }))

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

function newConversation() {
  return { messages: [] as any[], traceHistory: [] as any[], conversationID: null as any }
}

function mountChat(agentExec: any, conversation = newConversation()) {
  const w = mount(AiChat, {
    props: {
      agent: { agentID: '1' },
      conversation,
      showTrace: false,
    },
    global: {
      provide: {
        $SystemAPI: { agentExec },
        $toast: { toastErrorHandler: () => () => {} },
      },
      mocks: { $t: (k: string) => k },
      stubs: { CChatMessages: true, AiTrace: true, Button: true },
    },
  })

  return { w, conversation }
}

function contents(conversation: any): string {
  return conversation.messages.map((m: any) => String(m.content ?? '')).join('\n')
}

describe('AiChat approvals', () => {
  // The editor's chat is where an agent is tested, so it is the first place a
  // paused run shows up. It kept the output fallback the shared chat had
  // already dropped, and printed the whole response object as the agent's
  // answer — a wall of JSON where the reply belongs, with no way to approve.
  it('does not print the response object when a run pauses', async () => {
    const { w, conversation } = mountChat(
      vi.fn().mockResolvedValue({
        conversationID: '9',
        output: '',
        status: 'awaiting_approval',
        pendingApproval: { tool: 'compose_record_create', title: 'Create record', risk: 'write' },
      }),
    )

    await (w.vm as any).$.setupState.sendChatMessage('add a card')
    await flushPromises()

    expect(contents(conversation)).not.toContain('awaiting_approval')
    expect(contents(conversation)).not.toContain('conversationID')
    expect((w.vm as any).$.setupState.pending?.tool).toBe('compose_record_create')
  })

  // Without this the run is simply stuck: the agent asked and nothing can answer.
  it('resumes the paused run with the approved tool', async () => {
    const exec = vi
      .fn()
      .mockResolvedValueOnce({
        conversationID: '9',
        output: '',
        status: 'awaiting_approval',
        pendingApproval: { tool: 'compose_record_create', risk: 'write' },
      })
      .mockResolvedValueOnce({ conversationID: '9', output: 'Added it.', status: 'complete' })

    const { w, conversation } = mountChat(exec)
    const vm = (w.vm as any).$.setupState

    await vm.sendChatMessage('add a card')
    await flushPromises()
    await vm.onApprove(false)
    await flushPromises()

    expect(exec).toHaveBeenCalledTimes(2)
    expect(exec.mock.calls[1][0]).toMatchObject({
      input: '',
      conversationID: '9',
      approvedTools: ['compose_record_create'],
    })
    expect(contents(conversation)).toContain('Added it.')
    expect(vm.pending).toBeNull()
  })

  // Approving "for this chat" must be remembered, or every call asks again.
  it('remembers a chat-scoped approval for the next call', async () => {
    const exec = vi
      .fn()
      .mockResolvedValueOnce({
        conversationID: '9',
        output: '',
        status: 'awaiting_approval',
        pendingApproval: { tool: 'compose_record_create', risk: 'write' },
      })
      .mockResolvedValue({ conversationID: '9', output: 'ok', status: 'complete' })

    const { w } = mountChat(exec)
    const vm = (w.vm as any).$.setupState

    await vm.sendChatMessage('add a card')
    await flushPromises()
    await vm.onApprove(true)
    await flushPromises()
    await vm.sendChatMessage('add another')
    await flushPromises()

    expect(exec.mock.calls[1][0].approvedTools).toEqual(['compose_record_create'])
  })

  it('says so when the tool is refused, rather than going quiet', async () => {
    const { w, conversation } = mountChat(
      vi.fn().mockResolvedValue({
        conversationID: '9',
        output: '',
        status: 'awaiting_approval',
        pendingApproval: { tool: 'compose_record_delete', risk: 'destructive' },
      }),
    )
    const vm = (w.vm as any).$.setupState

    await vm.sendChatMessage('delete everything')
    await flushPromises()
    vm.onDeny()

    expect(vm.pending).toBeNull()
    expect(contents(conversation)).toContain('approval.denied')
  })

  // A finished run still needs the fallback: an agent that returns nothing
  // useful should not leave the turn blank.
  it('keeps the output fallback for a run that finished', async () => {
    const { w, conversation } = mountChat(
      vi.fn().mockResolvedValue({ conversationID: '9', response: { text: 'plain answer' } }),
    )

    await (w.vm as any).$.setupState.sendChatMessage('hi')
    await flushPromises()

    expect(contents(conversation)).toContain('plain answer')
  })
})
