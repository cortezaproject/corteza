import type { ChatbotConfig, Session } from './types'

// Base URL for the widget API. Always on the same origin as widget.js.
function apiRoot(scriptSrc: string): string {
  try {
    const u = new URL(scriptSrc)
    return `${u.origin}/api/widget/v1`
  } catch {
    return '/api/widget/v1'
  }
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
    if (!r.ok) throw new Error('widget: config failed')
    return r.json()
  }

  async openSession(): Promise<Session> {
    const r = await fetch(`${this.base}/session`, {
      method: 'POST',
      credentials: 'omit',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ widgetKey: this.key }),
    })
    if (!r.ok) throw new Error('widget: session failed')
    return r.json()
  }

  async sendMessage(sessionID: string, scenarioID: string, input: string): Promise<void> {
    const r = await fetch(
      `${this.base}/session/${encodeURIComponent(sessionID)}/message?widgetKey=${encodeURIComponent(this.key)}`,
      {
        method: 'POST',
        credentials: 'omit',
        headers: {
          'Content-Type': 'application/json',
          Authorization: `Bearer ${this.token}`,
        },
        body: JSON.stringify({ scenarioID, input }),
      },
    )
    if (r.status !== 202) throw new Error('widget: send failed')
  }

  async advanceStep(sessionID: string): Promise<{ nextScenarioID?: string; sessionComplete?: boolean }> {
    const r = await fetch(
      `${this.base}/session/${encodeURIComponent(sessionID)}/advance-step?widgetKey=${encodeURIComponent(this.key)}`,
      {
        method: 'POST',
        credentials: 'omit',
        headers: { Authorization: `Bearer ${this.token}` },
      },
    )
    if (!r.ok) throw new Error('widget: advance-step failed')
    return r.json()
  }

  // EventSource can't set headers → pass token as query.
  openStream(sessionID: string): EventSource {
    const url = `${this.base}/session/${encodeURIComponent(sessionID)}/stream?widgetKey=${encodeURIComponent(
      this.key,
    )}&token=${encodeURIComponent(this.token)}`
    return new EventSource(url, { withCredentials: false })
  }
}
