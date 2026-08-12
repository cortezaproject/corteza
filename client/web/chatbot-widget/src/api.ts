import type { ChatbotConfig, Session } from './types'

function apiRoot(scriptSrc: string): string {
  try {
    const u = new URL(scriptSrc)
    return `${u.origin}/api/widget/v1`
  } catch {
    return '/api/widget/v1'
  }
}

// Extracts a human-readable error from a failed widget HTTP response. The
// backend writes plain text via http.Error with a "widget: " prefix — strip
// it so the message lands cleanly in the chat as a system message.
async function errorFromResponse(r: Response, fallback: string): Promise<Error> {
  const raw = (await r.text().catch(() => '')).trim()
  const msg = raw.replace(/^widget:\s*/i, '').trim() || fallback
  return new Error(msg)
}

export class WidgetAPI {
  private base: string
  private key: string
  private token = ''

  constructor(scriptSrc: string, key: string) {
    this.base = apiRoot(scriptSrc)
    this.key = key
  }

  setToken(t: string) {
    this.token = t
  }

  get origin(): string {
    return this.base.replace(/\/api\/widget\/v1$/, '')
  }

  async config(): Promise<ChatbotConfig> {
    const r = await fetch(`${this.base}/config?widgetKey=${encodeURIComponent(this.key)}`, {
      credentials: 'omit',
    })
    if (!r.ok) throw await errorFromResponse(r, 'config failed')
    return r.json()
  }

  async openSession(): Promise<Session> {
    const r = await fetch(`${this.base}/session?widgetKey=${encodeURIComponent(this.key)}`, {
      method: 'POST',
      credentials: 'omit',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ widgetKey: this.key }),
    })
    if (!r.ok) throw await errorFromResponse(r, 'session failed')
    return r.json()
  }

  async sendMessage(sessionID: string, input: string): Promise<void> {
    const r = await this.post(sessionID, 'submit', { type: 'message', data: { input } })
    if (r.status !== 202) throw await errorFromResponse(r, 'send failed')
  }

  async submitConsent(sessionID: string, accepted: boolean): Promise<void> {
    const r = await this.post(sessionID, 'submit', { type: 'consent', data: { accepted } })
    if (r.status !== 204) throw await errorFromResponse(r, 'consent failed')
  }

  // submitForm posts the form values. Returns a map of field→error when the
  // server rejects validation (HTTP 422). Resolves with null on success.
  async submitForm(
    sessionID: string,
    fields: Record<string, string>,
  ): Promise<Record<string, string> | null> {
    const r = await this.post(sessionID, 'submit', { type: 'form', data: { fields } })
    if (r.status === 204) return null
    if (r.status === 422) {
      const body = (await r.json().catch(() => ({}))) as { errors?: Record<string, string> }
      return body.errors || {}
    }
    throw await errorFromResponse(r, 'form submit failed')
  }

  async advanceStep(sessionID: string): Promise<void> {
    const r = await this.post(sessionID, 'advance-step', {})
    if (r.status !== 204) throw await errorFromResponse(r, 'advance failed')
  }

  async closeSession(sessionID: string): Promise<void> {
    const r = await this.post(sessionID, 'close', {})
    if (r.status !== 204) throw await errorFromResponse(r, 'close failed')
  }

  async requestHandoff(sessionID: string, reason?: string): Promise<{ handoffID: string }> {
    const r = await this.post(sessionID, 'handoff', { reason: reason || '' })
    if (!r.ok) throw await errorFromResponse(r, 'handoff failed')
    return r.json()
  }

  async closeHandoff(sessionID: string, handoffID: string): Promise<void> {
    const r = await this.post(sessionID, 'handoff-complete', { handoffID })
    if (r.status !== 204) throw await errorFromResponse(r, 'handoff-complete failed')
  }

  openStream(sessionID: string): EventSource {
    const url = `${this.base}/session/${encodeURIComponent(sessionID)}/stream?widgetKey=${encodeURIComponent(
      this.key,
    )}&token=${encodeURIComponent(this.token)}`
    return new EventSource(url, { withCredentials: false })
  }

  private post(sessionID: string, path: string, body: any): Promise<Response> {
    return fetch(
      `${this.base}/session/${encodeURIComponent(sessionID)}/${path}?widgetKey=${encodeURIComponent(this.key)}`,
      {
        method: 'POST',
        credentials: 'omit',
        headers: {
          'Content-Type': 'application/json',
          Authorization: `Bearer ${this.token}`,
        },
        body: JSON.stringify(body),
      },
    )
  }
}
