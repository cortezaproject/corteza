import { ref } from 'vue'

// One agent turn, however the chat around it stores its messages.
//
// Two chats run agents — the sidebar and the page block through CAgentChat, and
// the agent editor's test chat through AiChat — and each assembled the request
// and read the response for itself. The approval gate was added to one and not
// the other, so a paused run printed the whole response object where the reply
// belongs and left the run with no way to continue. The parts that must not
// diverge live here; where a chat keeps its messages does not.

export type PendingApproval = {
  tool: string
  label: string
  risk?: string
  agentID: string
}

export type AgentTurnOptions = {
  // The agent this turn runs, read at call time: one chat switches agents.
  agentID: () => string | null | undefined
  conversationID: () => string | null | undefined

  exec: (req: Record<string, unknown>) => Promise<any>

  // Where a reply goes, and what the caller wants to do with the whole
  // response — record a trace, keep the returned context.
  onReply: (content: string, res: any, agentID: string) => void
  onResponse?: (res: any, agentID: string) => void
  onError?: (err: any, agentID: string) => void

  // Sent with every run when it returns anything.
  context?: () => Record<string, unknown> | undefined

  // What the agent says when the tool is refused.
  deniedMessage: () => string
}

export function useAgentTurn(opts: AgentTurnOptions) {
  const executing = ref(false)

  // What the agent stopped to ask about, if anything.
  const pending = ref<PendingApproval | null>(null)

  // Tools the user has approved, per agent and conversation. An approval is a
  // UX memory, not a control: the server asks again on a conversation it has
  // not been told about, and nothing here can grant what the invoking user
  // could not do anyway.
  const approvedTools = ref<Record<string, string[]>>({})

  function approvalKey(agentID: string) {
    return `${agentID}:${opts.conversationID() || 'new'}`
  }

  function isPaused(res: any) {
    return res?.status === 'awaiting_approval'
  }

  // A run that stopped to ask usually has nothing to say yet, and the fallback
  // would print the whole response object as the agent's answer. Only fall back
  // when the run actually finished.
  function replyOf(res: any): string {
    if (isPaused(res)) return res?.output || ''

    return (
      res?.output || (typeof res === 'string' ? res : res?.response?.text || JSON.stringify(res))
    )
  }

  // run covers both a fresh question and the resume that follows an approval; a
  // resume carries no input of its own, only the tools it may use.
  async function run(input: string, once: string[] = []) {
    const agentID = opts.agentID()
    if (!agentID) return

    executing.value = true
    pending.value = null

    try {
      const conversationID = opts.conversationID()
      const context = opts.context?.()
      const approved = [...(approvedTools.value[approvalKey(agentID)] || []), ...once]

      const res = await opts.exec({
        agentID,
        input,
        ...(conversationID ? { conversationID } : {}),
        ...(context && Object.keys(context).length > 0 ? { context } : {}),
        ...(approved.length ? { approvedTools: approved } : {}),
      })

      opts.onResponse?.(res, agentID)

      const reply = replyOf(res)
      if (reply) opts.onReply(reply, res, agentID)

      pending.value =
        isPaused(res) && res?.pendingApproval?.tool
          ? {
              tool: res.pendingApproval.tool,
              label: res.pendingApproval.title || res.pendingApproval.tool,
              risk: res.pendingApproval.risk,
              agentID,
            }
          : null

      return res
    } catch (err: any) {
      opts.onError?.(err, agentID)
    } finally {
      executing.value = false
    }
  }

  async function send(input: string) {
    return run(input)
  }

  async function approve(forChat: boolean) {
    const p = pending.value
    if (!p) return

    if (forChat) {
      const key = approvalKey(p.agentID)
      approvedTools.value = {
        ...approvedTools.value,
        [key]: [...(approvedTools.value[key] || []), p.tool],
      }
    }

    pending.value = null

    // A one-off approval is sent with the call and not remembered, so the next
    // use of the same tool asks again.
    return run('', forChat ? [] : [p.tool])
  }

  function deny() {
    const p = pending.value
    if (!p) return

    pending.value = null
    opts.onReply(opts.deniedMessage(), null, p.agentID)
  }

  return { executing, pending, approvedTools, send, run, approve, deny }
}
