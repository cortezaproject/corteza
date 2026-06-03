// Admin-side client for the chatbot preview API. Mirrors the public widget
// API surface but talks to /api/system/chatbot/preview/* and authenticates
// using the Corteza access token exposed by $SystemAPI.accessTokenFn().

export interface PreviewSession {
  sessionID: string
  token: string
  conversationID: string
  dbSessionID: string
}

export interface SystemAPILike {
  baseURL: string
  accessTokenFn?: () => string
}

function apiOrigin(baseURL: string): string {
  try {
    return new URL(baseURL).origin
  } catch {
    return window.location.origin
  }
}

// Extracts a human-readable error from a failed preview HTTP response. The
// backend writes plain text via http.Error with a "preview: " prefix — strip
// it so the message lands cleanly in the chat as a system message.
async function errorFromResponse(r: Response, fallback: string): Promise<Error> {
  const raw = (await r.text().catch(() => '')).trim()
  const msg = raw.replace(/^preview:\s*/i, '').trim() || fallback
  return new Error(msg)
}

export class PreviewClient {
  private base: string
  private origin: string
  private tokenFn: () => string

  constructor(systemAPI: SystemAPILike) {
    const sysBase = systemAPI.baseURL.replace(/\/+$/, '')
    // SystemAPI baseURL ends in /system. Replace to reach /chatbot/preview siblings.
    this.base = sysBase.replace(/\/system$/, '') + '/system/chatbot/preview'
    this.origin = apiOrigin(systemAPI.baseURL)
    this.tokenFn = systemAPI.accessTokenFn || (() => '')
  }

  get apiOrigin(): string {
    return this.origin
  }

  async openSession(chatbot: any): Promise<PreviewSession> {
    const r = await this.post('/session', { chatbot })
    if (!r.ok) throw await errorFromResponse(r, 'openSession failed')
    return r.json()
  }

  async startSession(sessionID: string): Promise<void> {
    const r = await this.post(`/session/${encodeURIComponent(sessionID)}/start`, {})
    if (r.status !== 204) throw await errorFromResponse(r, 'startSession failed')
  }

  async sendMessage(sessionID: string, input: string): Promise<void> {
    const r = await this.post(`/session/${encodeURIComponent(sessionID)}/submit`, {
      type: 'message',
      data: { input },
    })
    if (r.status !== 202) throw await errorFromResponse(r, 'sendMessage failed')
  }

  async submitConsent(sessionID: string, accepted: boolean): Promise<void> {
    const r = await this.post(`/session/${encodeURIComponent(sessionID)}/submit`, {
      type: 'consent',
      data: { accepted },
    })
    if (r.status !== 204) throw await errorFromResponse(r, 'submitConsent failed')
  }

  async submitForm(
    sessionID: string,
    fields: Record<string, string>,
  ): Promise<Record<string, string> | null> {
    const r = await this.post(`/session/${encodeURIComponent(sessionID)}/submit`, {
      type: 'form',
      data: { fields },
    })
    if (r.status === 204) return null
    if (r.status === 422) {
      const body = (await r.json().catch(() => ({}))) as { errors?: Record<string, string> }
      return body.errors || {}
    }
    throw await errorFromResponse(r, 'submitForm failed')
  }

  async advanceStep(sessionID: string): Promise<void> {
    const r = await this.post(`/session/${encodeURIComponent(sessionID)}/advance-step`, {})
    if (r.status !== 204) throw await errorFromResponse(r, 'advanceStep failed')
  }

  async closeSession(sessionID: string): Promise<void> {
    const r = await this.post(`/session/${encodeURIComponent(sessionID)}/close`, {})
    if (r.status !== 204) throw await errorFromResponse(r, 'closeSession failed')
  }

  async requestHandoff(sessionID: string): Promise<{ handoffID: string }> {
    const r = await this.post(`/session/${encodeURIComponent(sessionID)}/handoff`, {})
    if (!r.ok) throw await errorFromResponse(r, 'requestHandoff failed')
    return r.json()
  }

  async closeHandoff(sessionID: string): Promise<void> {
    const r = await this.post(`/session/${encodeURIComponent(sessionID)}/handoff-complete`, {})
    if (r.status !== 204) throw await errorFromResponse(r, 'closeHandoff failed')
  }

  async sendOperatorMessage(sessionID: string, message: string, operator?: string): Promise<void> {
    const r = await this.post(`/session/${encodeURIComponent(sessionID)}/operator-message`, {
      message,
      operator,
    })
    if (r.status !== 204) throw await errorFromResponse(r, 'operator-message failed')
  }

  openStream(sessionID: string): EventSource {
    // EventSource can't set headers — pass the bearer token as a query so the
    // admin auth middleware can pick it up. Corteza accepts ?jwt= for SSE GETs.
    const tok = this.tokenFn()
    const qs = tok ? `?jwt=${encodeURIComponent(tok)}` : ''
    return new EventSource(`${this.base}/session/${encodeURIComponent(sessionID)}/stream${qs}`, {
      withCredentials: true,
    })
  }

  private post(path: string, body: any): Promise<Response> {
    const tok = this.tokenFn()
    const headers: Record<string, string> = {
      'Content-Type': 'application/json',
    }
    if (tok) headers.Authorization = `Bearer ${tok}`
    return fetch(this.base + path, {
      method: 'POST',
      credentials: 'include',
      headers,
      body: JSON.stringify(body),
    })
  }
}
