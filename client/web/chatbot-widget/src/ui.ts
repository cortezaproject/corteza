import type { ChatbotConfig, StepStartPayload } from './types'
import { Engine, type HandoffState, type Message } from './engine'
import { applyStyling, baseCSS } from './styles'
import { renderMarkdown } from '@planetcrust/human-js'

const SVG_CLOSE =
  '<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2" stroke-linecap="round"><path d="M6 6l12 12M6 18L18 6"/></svg>'
const SVG_SEND =
  '<svg viewBox="0 0 24 24" fill="currentColor"><path d="M3.4 20.6L21 12 3.4 3.4 3 10l12 2-12 2z"/></svg>'
const SVG_COMMENTS =
  '<svg viewBox="0 0 14 14" fill="currentColor" xmlns="http://www.w3.org/2000/svg"><path d="M10.5,0h-8C1.119,0,0,1.119,0,2.5v5C0,8.881,1.119,10,2.5,10H3v2.5c0,0.202,0.122,0.385,0.309,0.462C3.371,12.988,3.436,13,3.5,13c0.13,0,0.259-0.051,0.354-0.146L6.707,10H10.5C11.881,10,13,8.881,13,7.5v-5C13,1.119,11.881,0,10.5,0z"/></svg>'

function capitalizeFirst(s: string): string {
  if (!s) return s
  return s.charAt(0).toUpperCase() + s.slice(1)
}

// sanitizeRichHTML strips the input to an allow-list of tags/attributes
// suitable for consent-step body content (TOS, privacy notices). Anything
// outside the list is replaced with its text content. Used because the
// editor emits raw HTML from a rich-text editor and we don't want to push
// that straight into innerHTML.
const HB_CONSENT_ALLOWED_TAGS = new Set([
  'p', 'br', 'span', 'strong', 'em', 'u', 's', 'a', 'ul', 'ol', 'li',
  'h1', 'h2', 'h3', 'h4', 'blockquote', 'code', 'pre',
])
// Inline-style declarations the RTE produces that we want to preserve so the
// chatbot renders text the way the author saw it in the editor (TipTap's
// TextStyle / Color extensions emit these on <span>). Anything outside the
// list is dropped. We additionally reject values containing url(), expression()
// or javascript: to block CSS-vector XSS.
const HB_CONSENT_ALLOWED_STYLE_PROPS = new Set([
  'color', 'background-color', 'font-weight', 'font-style',
  'text-decoration', 'text-decoration-line', 'text-align',
])
function sanitizeStyle(raw: string): string {
  return raw
    .split(';')
    .map(decl => decl.trim())
    .filter(Boolean)
    .map(decl => {
      const idx = decl.indexOf(':')
      if (idx < 0) return ''
      const prop = decl.slice(0, idx).trim().toLowerCase()
      const val = decl.slice(idx + 1).trim()
      if (!HB_CONSENT_ALLOWED_STYLE_PROPS.has(prop)) return ''
      if (/url\(|expression\(|javascript:|@import/i.test(val)) return ''
      return `${prop}: ${val}`
    })
    .filter(Boolean)
    .join('; ')
}
function sanitizeRichHTML(html: string): string {
  if (!html) return ''
  const tpl = document.createElement('template')
  tpl.innerHTML = html
  hbScrubNode(tpl.content)
  return tpl.innerHTML
}
function hbScrubNode(node: Node): void {
  const children = Array.from(node.childNodes)
  for (const child of children) {
    if (child.nodeType === Node.ELEMENT_NODE) {
      const el = child as Element
      const tag = el.tagName.toLowerCase()
      if (!HB_CONSENT_ALLOWED_TAGS.has(tag)) {
        const replacement = document.createTextNode(el.textContent || '')
        el.replaceWith(replacement)
        continue
      }
      // Capture safe attributes before stripping everything else.
      let safeHref = ''
      if (tag === 'a') {
        const href = el.getAttribute('href') || ''
        if (/^(https?:|mailto:|tel:|\/|#)/i.test(href)) safeHref = href
      }
      const safeStyle = sanitizeStyle(el.getAttribute('style') || '')
      const attrs = Array.from(el.attributes).map(a => a.name)
      for (const a of attrs) el.removeAttribute(a)
      if (tag === 'a' && safeHref) {
        el.setAttribute('href', safeHref)
        el.setAttribute('target', '_blank')
        el.setAttribute('rel', 'noopener noreferrer')
      }
      if (safeStyle) el.setAttribute('style', safeStyle)
      hbScrubNode(el)
    } else if (child.nodeType !== Node.TEXT_NODE) {
      child.parentNode?.removeChild(child)
    }
  }
}

export interface WidgetUIOptions {
  container?: HTMLElement
  contained?: boolean
  startOpen?: boolean
  onToggle?: (open: boolean) => void
}

export class WidgetUI {
  private root: ShadowRoot
  private host: HTMLElement
  private panelEl!: HTMLDivElement
  private bodyEl!: HTMLDivElement
  private footerEl!: HTMLDivElement
  private typingEl: HTMLDivElement | null = null
  private currentAgentEl: HTMLDivElement | null = null
  private agentBuffer = ''
  private engine: Engine
  private cfg: ChatbotConfig
  private mo: MutationObserver | null = null
  private onToggle?: (open: boolean) => void
  private currentForm: HTMLFormElement | null = null
  // Mirrors the typing-indicator state so rapid double-sends can't trigger
  // overlapping runtime calls (the BE conversation save races on a stale
  // UpdatedAt otherwise). Updated from the typing on/off events; gates the
  // composer disabled state alongside the empty-input check.
  private agentBusy = false

  onUserInput: (text: string) => void = () => {}
  onFormSubmit: (values: Record<string, string>) => void = () => {}
  onConsentDecision: (accepted: boolean) => void = () => {}
  onRequestHandoff: () => void = () => {}
  onCancelHandoff: () => void = () => {}
  onEndConversation: () => void = () => {}
  onCloseSession: () => void = () => {}

  constructor(cfg: ChatbotConfig, engine: Engine, opts: WidgetUIOptions = {}) {
    this.cfg = cfg
    this.engine = engine
    this.onToggle = opts.onToggle

    const host = document.createElement('div')
    host.setAttribute('data-human-chatbot', '')
    const parent = opts.container || document.body
    const attach = () => {
      if (!host.isConnected) parent.appendChild(host)
    }
    attach()
    this.host = host

    if (!opts.container) {
      this.mo = new MutationObserver(() => attach())
      this.mo.observe(document.body, { childList: true })
    }

    this.root = host.attachShadow({ mode: 'open' })
    const style = document.createElement('style')
    style.textContent = baseCSS
    this.root.appendChild(style)

    const rootEl = document.createElement('div')
    const allowedPositions = [
      'bottom-right', 'bottom-left', 'bottom-center',
      'top-right', 'top-left', 'top-center',
      'left-middle', 'right-middle',
    ]
    const position = allowedPositions.includes(cfg.styling.launcher.position)
      ? cfg.styling.launcher.position
      : 'bottom-right'
    rootEl.className = 'hb-root hb-pos-' + position + (opts.contained ? ' hb-contained' : '')
    this.root.appendChild(rootEl)
    applyStyling(rootEl, cfg.styling)

    const launcher = document.createElement('button')
    launcher.className = 'hb-launcher'
    const iconVisible = cfg.styling.launcher.iconVisible !== false
    const buttonLabel = cfg.styling.launcher.buttonLabel || ''
    if (iconVisible) {
      if (cfg.styling.launcher.iconURL) {
        const img = document.createElement('img')
        img.src = cfg.styling.launcher.iconURL
        img.alt = cfg.styling.launcher.label || 'Chat'
        launcher.appendChild(img)
      } else {
        const iconWrap = document.createElement('span')
        iconWrap.className = 'hb-launcher-icon'
        iconWrap.innerHTML = SVG_COMMENTS
        launcher.appendChild(iconWrap)
        launcher.classList.add('hb-launcher-default')
      }
    }
    if (buttonLabel) {
      const lbl = document.createElement('span')
      lbl.className = 'hb-launcher-label'
      lbl.textContent = buttonLabel
      launcher.appendChild(lbl)
      launcher.classList.add('hb-launcher-has-label')
    }
    launcher.addEventListener('click', () => this.togglePanel())
    rootEl.appendChild(launcher)

    this.panelEl = document.createElement('div')
    this.panelEl.className = 'hb-panel' + (opts.startOpen ? '' : ' hb-hidden')
    rootEl.appendChild(this.panelEl)

    const header = document.createElement('div')
    header.className = 'hb-header'
    if (cfg.styling.logoURL) {
      const img = document.createElement('img')
      img.className = 'hb-logo'
      img.src = cfg.styling.logoURL
      img.alt = cfg.styling.launcher.label || 'Logo'
      header.appendChild(img)
    }
    const title = document.createElement('span')
    title.className = 'hb-title'
    title.textContent = cfg.styling.launcher.label || ''
    header.appendChild(title)

    const closeBtn = document.createElement('button')
    closeBtn.className = 'hb-close'
    closeBtn.setAttribute('aria-label', 'Close chat')
    closeBtn.innerHTML = SVG_CLOSE
    closeBtn.addEventListener('click', () => this.togglePanel())
    header.appendChild(closeBtn)

    this.panelEl.appendChild(header)

    this.bodyEl = document.createElement('div')
    this.bodyEl.className = 'hb-body'
    this.panelEl.appendChild(this.bodyEl)

    this.footerEl = document.createElement('div')
    this.footerEl.className = 'hb-footer hb-hidden'
    this.panelEl.appendChild(this.footerEl)

    this.engine.on(ev => {
      switch (ev.type) {
        case 'step_start':
          return this.renderStep(ev.payload)
        case 'step_complete':
          return
        case 'message':
          return this.appendMessage(ev.message)
        case 'typing':
          return this.setTyping(ev.on)
        case 'agent_start':
          return this.beginAgentStream()
        case 'agent_delta':
          return this.appendAgentDelta(ev.text)
        case 'agent_end':
          return this.endAgentStream()
        case 'form_error':
          return this.showFormErrors(ev.errors)
        case 'handoff_change':
          return this.renderHandoffBadge(ev.state)
        case 'session_closed':
          return this.renderClosed()
        case 'error':
          return this.appendMessage(
            { role: 'system', content: capitalizeFirst(ev.error) },
            { error: true },
          )
      }
    })
  }

  private beginAgentStream() {
    this.setTyping(false)
    const row = document.createElement('div')
    row.className = 'hb-msg agent'
    this.bodyEl.appendChild(row)
    this.currentAgentEl = row
    this.agentBuffer = ''
    this.bodyEl.scrollTop = this.bodyEl.scrollHeight
  }

  private appendAgentDelta(text: string) {
    if (!this.currentAgentEl) this.beginAgentStream()
    this.agentBuffer += text
    if (this.currentAgentEl) this.currentAgentEl.innerHTML = renderMarkdown(this.agentBuffer)
    this.bodyEl.scrollTop = this.bodyEl.scrollHeight
  }

  private endAgentStream() {
    this.currentAgentEl = null
    this.agentBuffer = ''
  }

  togglePanel() {
    this.panelEl.classList.toggle('hb-hidden')
    if (this.onToggle) this.onToggle(!this.panelEl.classList.contains('hb-hidden'))
  }

  destroy() {
    if (this.mo) {
      this.mo.disconnect()
      this.mo = null
    }
    this.host.remove()
  }

  appendMessage(m: Message, opts: { error?: boolean } = {}) {
    const row = document.createElement('div')
    const cls =
      m.role === 'user'
        ? 'user'
        : m.role === 'agent' || m.role === 'operator'
          ? 'agent'
          : 'system'
    row.className = `hb-msg ${cls}${opts.error ? ' error' : ''}`
    if (m.role === 'operator') {
      const tag = document.createElement('div')
      tag.className = 'hb-operator-tag'
      tag.textContent = m.operator ? m.operator : 'Operator'
      row.appendChild(tag)
      const body = document.createElement('div')
      body.className = 'hb-msg-body'
      body.textContent = m.content
      row.appendChild(body)
    } else if (m.role === 'agent') {
      row.innerHTML = renderMarkdown(m.content)
    } else {
      row.textContent = m.content
    }
    this.bodyEl.appendChild(row)
    this.bodyEl.scrollTop = this.bodyEl.scrollHeight
  }

  setTyping(on: boolean) {
    if (on && !this.typingEl) {
      const el = document.createElement('div')
      el.className = 'hb-typing'
      el.innerHTML = '<span></span><span></span><span></span>'
      this.bodyEl.appendChild(el)
      this.typingEl = el
      this.bodyEl.scrollTop = this.bodyEl.scrollHeight
    } else if (!on && this.typingEl) {
      this.typingEl.remove()
      this.typingEl = null
    }
    this.agentBusy = on
    this.refreshComposerBusyState()
  }

  // Look up the currently-mounted conversation composer and toggle its
  // disabled state to match `agentBusy`. The send button still respects the
  // empty-input rule when re-enabled. No-op when the composer isn't mounted
  // (e.g. on a form / static step).
  private refreshComposerBusyState() {
    const input = this.footerEl.querySelector<HTMLTextAreaElement>('.hb-input')
    const send = this.footerEl.querySelector<HTMLButtonElement>('.hb-send')
    if (!input || !send) return
    input.disabled = this.agentBusy
    if (this.agentBusy) {
      send.disabled = true
    } else {
      send.disabled = input.value.trim().length === 0
    }
  }

  // Appends the persistent "Talk to human" button to a per-step actions
  // container so it sits inline with whatever step-specific buttons exist
  // (e.g. End conversation, Next, Close). No-op if handoff is disabled.
  // Pre-disables the button when a handoff is already requested/active so a
  // step transition mid-handoff doesn't re-enable it.
  private appendHandoffAction(actions: HTMLElement) {
    if (!this.cfg.handoff.enabled || this.cfg.handoff.notImplemented) return
    const btn = document.createElement('button')
    btn.className = 'hb-action'
    btn.type = 'button'
    btn.dataset.action = 'handoff'
    btn.textContent = 'Talk to human'
    if (this.engine.handoffState.phase !== 'idle') btn.disabled = true
    btn.addEventListener('click', () => this.onRequestHandoff())
    actions.appendChild(btn)
  }

  private renderStep(p: StepStartPayload) {
    this.footerEl.className = 'hb-footer hb-hidden'
    this.footerEl.innerHTML = ''
    this.currentForm = null

    switch (p.type) {
      case 'static_message':
        return this.renderStaticMessage(p)
      case 'form':
        return this.renderForm(p)
      case 'conversation':
        return this.renderConversation(p)
      case 'consent':
        return this.renderConsent(p)
    }
  }

  private renderStaticMessage(p: StepStartPayload) {
    const cfg = p.config || {}
    // Static messages reuse the consent-step body styling so author HTML/RTE
    // content flows full-width with neutral typography rather than living in
    // an agent-bubble. No buttons under the body — static steps advance via
    // auto-advance / scenario flow.
    const block = document.createElement('div')
    block.className = 'hb-consent-block'

    const body = document.createElement('div')
    body.className = 'hb-consent-body'
    body.innerHTML = sanitizeRichHTML(cfg.message || '')
    block.appendChild(body)

    this.bodyEl.appendChild(block)
    this.bodyEl.scrollTop = this.bodyEl.scrollHeight

    // Static messages advance via auto-advance / scenario flow — no Next/Close
    // button. The only footer affordance is the optional "Talk to human"
    // button when handoff is enabled.
    const actions = document.createElement('div')
    actions.className = 'hb-conv-actions'
    this.appendHandoffAction(actions)

    if (actions.children.length === 0) return
    this.footerEl.className = 'hb-footer'
    this.footerEl.appendChild(actions)
  }

  private renderForm(p: StepStartPayload) {
    const cfg = p.config || {}
    const fields: Array<{ name: string; label: string; type?: string; required?: boolean }> =
      cfg.fields || []
    const form = document.createElement('form')
    form.className = 'hb-scenario-form'
    fields.forEach(f => {
      const row = document.createElement('div')
      row.className = 'hb-field'
      row.dataset.field = f.name
      const lbl = document.createElement('label')
      lbl.textContent = f.label || f.name
      row.appendChild(lbl)
      const input = document.createElement('input')
      input.name = f.name
      input.type = f.type || 'text'
      if (f.required) input.required = true
      row.appendChild(input)
      const err = document.createElement('div')
      err.className = 'hb-field-error hb-hidden'
      row.appendChild(err)
      form.appendChild(row)
    })
    const actions = document.createElement('div')
    actions.className = 'hb-form-actions'
    const submit = document.createElement('button')
    submit.className = 'hb-submit'
    submit.type = 'submit'
    submit.textContent = cfg.submitLabel || 'Submit'
    actions.appendChild(submit)
    form.appendChild(actions)

    form.addEventListener('submit', ev => {
      ev.preventDefault()
      const values: Record<string, string> = {}
      fields.forEach(f => {
        const el = form.elements.namedItem(f.name) as HTMLInputElement | null
        if (el) values[f.name] = el.value
      })
      // Immediately disable form and show submitted state
      submit.disabled = true
      submit.textContent = '✓'
      form.querySelectorAll('input, textarea, select').forEach(el => {
        (el as HTMLInputElement).disabled = true
      })
      this.onFormSubmit(values)
    })
    this.bodyEl.appendChild(form)
    this.currentForm = form

    // Form Submit lives inside the form itself. The only footer affordance is
    // the optional "Talk to human" button when handoff is enabled — there is
    // no End conversation here (the visitor either submits or abandons).
    const formActions = document.createElement('div')
    formActions.className = 'hb-conv-actions'
    this.appendHandoffAction(formActions)
    if (formActions.children.length === 0) return
    this.footerEl.className = 'hb-footer'
    this.footerEl.innerHTML = ''
    this.footerEl.appendChild(formActions)
  }

  // showFormErrors displays inline validation errors on the current form.
  // No-op if no form is mounted.
  private showFormErrors(errors: Record<string, string>) {
    const form = this.currentForm
    if (!form) return
    form.querySelectorAll<HTMLDivElement>('.hb-field-error').forEach(el => {
      el.textContent = ''
      el.classList.add('hb-hidden')
    })
    Object.entries(errors).forEach(([name, msg]) => {
      const row = form.querySelector<HTMLDivElement>(`.hb-field[data-field="${CSS.escape(name)}"]`)
      if (!row) return
      const errEl = row.querySelector<HTMLDivElement>('.hb-field-error')
      if (!errEl) return
      errEl.textContent = msg
      errEl.classList.remove('hb-hidden')
    })
  }

  private renderConversation(p: StepStartPayload) {
    this.footerEl.className = 'hb-footer'
    this.footerEl.innerHTML = ''
    const input = document.createElement('textarea')
    input.className = 'hb-input'
    input.rows = 1
    input.placeholder =
      (p.config as { placeholder?: string } | undefined)?.placeholder || 'Type a message…'

    const send = document.createElement('button')
    send.className = 'hb-send'
    send.type = 'button'
    send.innerHTML = SVG_SEND + '<span class="hb-sr-only">Send</span>'
    send.disabled = true
    send.setAttribute('aria-label', 'Send')

    const autogrow = () => {
      input.style.height = 'auto'
      input.style.height = Math.min(input.scrollHeight, 120) + 'px'
    }
    const refreshDisabled = () => {
      send.disabled = this.agentBusy || input.value.trim().length === 0
    }
    const submit = () => {
      // Front-stop the rapid double-send: if the agent is still working on
      // the previous turn, the BE save would race with this one and trip
      // the conversation isStale check. The typing event re-enables the
      // composer as soon as the runtime is idle again.
      if (this.agentBusy) return
      const v = input.value.trim()
      if (!v) return
      this.onUserInput(v)
      input.value = ''
      autogrow()
      refreshDisabled()
    }

    send.addEventListener('click', submit)
    input.addEventListener('input', () => {
      autogrow()
      refreshDisabled()
    })
    input.addEventListener('keydown', ev => {
      if (ev.key === 'Enter' && !ev.shiftKey) {
        ev.preventDefault()
        submit()
      }
    })

    this.footerEl.appendChild(input)
    this.footerEl.appendChild(send)
    // Apply current busy state if the agent is mid-stream when this composer
    // mounts (e.g. after a step transition while a response is still
    // generating).
    this.refreshComposerBusyState()

    const actions = document.createElement('div')
    actions.className = 'hb-conv-actions'

    this.appendHandoffAction(actions)

    const end = document.createElement('button')
    end.className = 'hb-action'
    end.type = 'button'
    end.dataset.action = 'end'
    end.textContent = (p.config as { endLabel?: string } | undefined)?.endLabel || 'End conversation'
    end.addEventListener('click', () => this.onEndConversation())
    actions.appendChild(end)

    this.footerEl.appendChild(actions)
  }

  private renderHandoffBadge(state: HandoffState) {
    let badge = this.panelEl.querySelector<HTMLDivElement>('.hb-handoff-badge')
    if (state.phase === 'idle') {
      if (badge) badge.remove()
      // Re-enable any handoff buttons (persistent bar + any per-step buttons).
      this.panelEl.querySelectorAll<HTMLButtonElement>('button[data-action="handoff"]').forEach(b => {
        b.disabled = false
      })
      return
    }
    if (!badge) {
      badge = document.createElement('div')
      badge.className = 'hb-handoff-badge'
      this.panelEl.insertBefore(badge, this.bodyEl)
    }
    // Reset before re-rendering — handoff state may transition from
    // 'requested' (with Cancel button) to 'active' on the same badge node.
    badge.textContent = ''
    if (state.phase === 'requested') {
      badge.appendChild(document.createTextNode('Waiting for an operator…'))
    } else {
      badge.appendChild(document.createTextNode('Connected to '))
      const name = document.createElement('span')
      name.className = 'hb-operator-name'
      name.textContent = state.operator || 'operator'
      badge.appendChild(name)
    }

    if (state.phase === 'requested') {
      const cancel = document.createElement('button')
      cancel.className = 'hb-action'
      cancel.type = 'button'
      cancel.textContent = 'Cancel'
      cancel.addEventListener('click', () => this.onCancelHandoff())
      badge.appendChild(document.createTextNode(' '))
      badge.appendChild(cancel)
    }

    // Disable the request-handoff buttons while we are mid-handoff (persistent
    // bar + any per-step buttons).
    this.panelEl.querySelectorAll<HTMLButtonElement>('button[data-action="handoff"]').forEach(b => {
      b.disabled = true
    })
  }

  private renderConsent(p: StepStartPayload) {
    const cfg = (p.config || {}) as {
      body?: string
      acceptLabel?: string
      rejectLabel?: string
    }

    // Consent body sits in the message stream styled like a system info
    // message (centered, muted) rather than an agent bubble. Accept/Reject
    // live directly under the text so the decision stays anchored to the
    // content the visitor is consenting to.
    const block = document.createElement('div')
    block.className = 'hb-consent-block'

    const body = document.createElement('div')
    body.className = 'hb-consent-body'
    body.innerHTML = sanitizeRichHTML(cfg.body || '')
    block.appendChild(body)

    const decision = document.createElement('div')
    decision.className = 'hb-consent-decision'

    const reject = document.createElement('button')
    reject.className = 'hb-action hb-consent-reject'
    reject.type = 'button'
    reject.dataset.action = 'consent-reject'
    reject.textContent = cfg.rejectLabel || 'Reject'
    reject.addEventListener('click', () => {
      reject.disabled = true
      accept.disabled = true
      this.onConsentDecision(false)
    })

    const accept = document.createElement('button')
    accept.className = 'hb-submit hb-consent-accept'
    accept.type = 'button'
    accept.dataset.action = 'consent-accept'
    accept.textContent = cfg.acceptLabel || 'Accept'
    accept.addEventListener('click', () => {
      accept.disabled = true
      reject.disabled = true
      this.onConsentDecision(true)
    })

    decision.appendChild(reject)
    decision.appendChild(accept)
    block.appendChild(decision)
    this.bodyEl.appendChild(block)
    this.bodyEl.scrollTop = this.bodyEl.scrollHeight

    // Footer only carries the persistent Talk to human button (when handoff
    // is enabled); when handoff is disabled the footer stays hidden.
    const actions = document.createElement('div')
    actions.className = 'hb-conv-actions'
    this.appendHandoffAction(actions)
    if (actions.children.length === 0) return
    this.footerEl.className = 'hb-footer'
    this.footerEl.appendChild(actions)
  }

  private renderClosed() {
    this.footerEl.className = 'hb-footer hb-hidden'
    this.footerEl.innerHTML = ''
    const row = document.createElement('div')
    row.className = 'hb-msg system'
    row.textContent = 'Conversation ended.'
    this.bodyEl.appendChild(row)
    this.bodyEl.scrollTop = this.bodyEl.scrollHeight
  }
}
