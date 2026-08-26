import { describe, it, expect, vi } from 'vitest'
import { useAgentTurn } from './useAgentTurn'

function setup(exec: any, over: Partial<any> = {}) {
  const replies: string[] = []
  let conversationID: string | null = null

  const turn = useAgentTurn({
    agentID: () => '1',
    conversationID: () => conversationID,
    exec,
    onReply: c => replies.push(c),
    onResponse: (res: any) => {
      if (res?.conversationID) conversationID = res.conversationID
    },
    deniedMessage: () => 'refused',
    ...over,
  })

  return { turn, replies, convID: () => conversationID }
}

const paused = {
  conversationID: '9',
  output: '',
  status: 'awaiting_approval',
  pendingApproval: { tool: 'compose_record_create', title: 'Create record', risk: 'write' },
}

describe('useAgentTurn', () => {
  // The fallback exists so an agent that returns nothing useful does not leave
  // the turn blank; on a paused run it printed the response object instead.
  it('says nothing when a run stops to ask', async () => {
    const { turn, replies } = setup(vi.fn().mockResolvedValue(paused))

    await turn.send('add a card')

    expect(replies).toEqual([])
    expect(turn.pending.value?.tool).toBe('compose_record_create')
    expect(turn.pending.value?.label).toBe('Create record')
  })

  it('falls back to the response text on a run that finished', async () => {
    const { turn, replies } = setup(
      vi.fn().mockResolvedValue({ response: { text: 'plain answer' } }),
    )

    await turn.send('hi')

    expect(replies).toEqual(['plain answer'])
    expect(turn.pending.value).toBeNull()
  })

  // A pending approval left standing after the run resumes would ask twice.
  it('clears the question once the run carries on', async () => {
    const exec = vi
      .fn()
      .mockResolvedValueOnce(paused)
      .mockResolvedValueOnce({ conversationID: '9', output: 'Added it.', status: 'complete' })
    const { turn, replies } = setup(exec)

    await turn.send('add a card')
    await turn.approve(false)

    expect(replies).toEqual(['Added it.'])
    expect(turn.pending.value).toBeNull()
  })

  it('sends a one-off approval without remembering it', async () => {
    const exec = vi
      .fn()
      .mockResolvedValueOnce(paused)
      .mockResolvedValue({ conversationID: '9', output: 'ok', status: 'complete' })
    const { turn } = setup(exec)

    await turn.send('add a card')
    await turn.approve(false)
    await turn.send('add another')

    expect(exec.mock.calls[1][0]).toMatchObject({
      input: '',
      conversationID: '9',
      approvedTools: ['compose_record_create'],
    })
    expect(exec.mock.calls[2][0].approvedTools).toBeUndefined()
  })

  it('remembers a chat-scoped approval for every later call', async () => {
    const exec = vi
      .fn()
      .mockResolvedValueOnce(paused)
      .mockResolvedValue({ conversationID: '9', output: 'ok', status: 'complete' })
    const { turn } = setup(exec)

    await turn.send('add a card')
    await turn.approve(true)
    await turn.send('add another')

    expect(exec.mock.calls[2][0].approvedTools).toEqual(['compose_record_create'])
  })

  // The memory is keyed by conversation, so a new one asks again.
  it('does not carry an approval into another conversation', async () => {
    let conversationID: string | null = '9'
    const exec = vi.fn().mockResolvedValue({ output: 'ok', status: 'complete' })

    const turn = useAgentTurn({
      agentID: () => '1',
      conversationID: () => conversationID,
      exec,
      onReply: () => {},
      deniedMessage: () => 'refused',
    })

    turn.approvedTools.value = { '1:9': ['compose_record_create'] }
    await turn.send('one')
    expect(exec.mock.calls[0][0].approvedTools).toEqual(['compose_record_create'])

    conversationID = '10'
    await turn.send('two')
    expect(exec.mock.calls[1][0].approvedTools).toBeUndefined()
  })

  it('says so when the tool is refused, rather than going quiet', async () => {
    const { turn, replies } = setup(vi.fn().mockResolvedValue(paused))

    await turn.send('delete everything')
    turn.deny()

    expect(replies).toEqual(['refused'])
    expect(turn.pending.value).toBeNull()
  })

  it('reports an error without leaving the chat busy', async () => {
    const errors: any[] = []
    const { turn } = setup(vi.fn().mockRejectedValue(new Error('boom')), {
      onError: (e: any) => errors.push(e),
    })

    await turn.send('hi')

    expect(errors).toHaveLength(1)
    expect(turn.executing.value).toBe(false)
  })
})
