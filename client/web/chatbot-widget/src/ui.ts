import type { ChatbotConfig, StepStartPayload } from './types'
import { Engine, type HandoffState, type Message } from './engine'
import { applyStyling, baseCSS } from './styles'
import { renderMarkdown } from './md'

const SVG_CLOSE =
  '<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2" stroke-linecap="round"><path d="M6 6l12 12M6 18L18 6"/></svg>'
const SVG_SEND =
  '<svg viewBox="0 0 24 24" fill="currentColor"><path d="M3.4 20.6L21 12 3.4 3.4 3 10l12 2-12 2z"/></svg>'
const SVG_COMMENTS =
  '<svg viewBox="0 0 14 14" fill="currentColor" xmlns="http://www.w3.org/2000/svg"><path d="M10.5,0h-8C1.119,0,0,1.119,0,2.5v5C0,8.881,1.119,10,2.5,10H3v2.5c0,0.202,0.122,0.385,0.309,0.462C3.371,12.988,3.436,13,3.5,13c0.13,0,0.259-0.051,0.354-0.146L6.707,10H10.5C11.881,10,13,8.881,13,7.5v-5C13,1.119,11.881,0,10.5,0z"/></svg>'

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
          return this.appendMessage({ role: 'system', content: ev.error })
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

  appendMessage(m: Message) {
    const row = document.createElement('div')
    const cls =
      m.role === 'user'
        ? 'user'
        : m.role === 'agent' || m.role === 'operator'
          ? 'agent'
          : 'system'
    row.className = `hb-msg ${cls}`
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
    }
  }

  private renderStaticMessage(p: StepStartPayload) {
    const cfg = p.config || {}
    const row = document.createElement('div')
    row.className = 'hb-msg agent'
    if (cfg.isMarkdown) {
      row.innerHTML = renderMarkdown(cfg.message || '')
    } else {
      row.textContent = cfg.message || ''
    }
    this.bodyEl.appendChild(row)
    this.bodyEl.scrollTop = this.bodyEl.scrollHeight

    // Determine whether more scenarios follow. Past-last → close button;
    // otherwise → next button. Auto-advance timing is driven by entry.ts.
    const lastIdx = this.cfg.scenarios.length - 1
    const isLast = p.scenarioIndex >= lastIdx
    if (cfg.showCloseButton === false) return

    this.footerEl.className = 'hb-footer'
    const btn = document.createElement('button')
    btn.className = 'hb-submit'
    btn.type = 'button'
    if (isLast) {
      btn.textContent = cfg.closeLabel || 'Close'
      btn.addEventListener('click', () => this.onCloseSession())
    } else {
      btn.textContent = cfg.nextLabel || 'Next'
      btn.addEventListener('click', () => this.onEndConversation())
    }
    this.footerEl.appendChild(btn)
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

    // Mirror the conversation scenario's footer: a small "End conversation"
    // action button so the visitor can bail out of a form step the same way
    // they would from a chat step. Reuses the same hb-conv-actions /
    // hb-action styling for visual parity.
    this.footerEl.className = 'hb-footer'
    this.footerEl.innerHTML = ''
    const formActions = document.createElement('div')
    formActions.className = 'hb-conv-actions'
    const end = document.createElement('button')
    end.className = 'hb-action'
    end.type = 'button'
    end.dataset.action = 'end'
    end.textContent = (p.config as { endLabel?: string } | undefined)?.endLabel || 'End conversation'
    end.addEventListener('click', () => this.onEndConversation())
    formActions.appendChild(end)
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

    if (this.cfg.handoff.enabled && !this.cfg.handoff.notImplemented) {
      const ho = document.createElement('button')
      ho.className = 'hb-action'
      ho.type = 'button'
      ho.dataset.action = 'handoff'
      ho.textContent = (p.config as { handoffLabel?: string } | undefined)?.handoffLabel || 'Talk to human'
      ho.addEventListener('click', () => this.onRequestHandoff())
      actions.appendChild(ho)
    }

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
      // Re-enable handoff button in the conv footer if present.
      this.footerEl.querySelectorAll<HTMLButtonElement>('button[data-action="handoff"]').forEach(b => {
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

    // Disable the request-handoff button while we are mid-handoff.
    this.footerEl.querySelectorAll<HTMLButtonElement>('button[data-action="handoff"]').forEach(b => {
      b.disabled = true
    })
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
