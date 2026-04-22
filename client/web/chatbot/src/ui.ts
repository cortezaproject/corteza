import type { ChatbotConfig, Scenario } from './types'
import { Engine, type Message } from './engine'
import { applyStyling, baseCSS } from './styles'
import { renderMarkdown } from './md'

const SVG_CLOSE =
  '<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2" stroke-linecap="round"><path d="M6 6l12 12M6 18L18 6"/></svg>'
const SVG_SEND =
  '<svg viewBox="0 0 24 24" fill="currentColor"><path d="M3.4 20.6L21 12 3.4 3.4 3 10l12 2-12 2z"/></svg>'
// PrimeIcons "comments" glyph — used when no custom launcher icon is set.
const SVG_COMMENTS =
  '<svg viewBox="0 0 14 14" fill="currentColor" xmlns="http://www.w3.org/2000/svg"><path d="M10.5,0h-8C1.119,0,0,1.119,0,2.5v5C0,8.881,1.119,10,2.5,10H3v2.5c0,0.202,0.122,0.385,0.309,0.462C3.371,12.988,3.436,13,3.5,13c0.13,0,0.259-0.051,0.354-0.146L6.707,10H10.5C11.881,10,13,8.881,13,7.5v-5C13,1.119,11.881,0,10.5,0z"/></svg>'

export interface WidgetUIOptions {
  container?: HTMLElement
  contained?: boolean
  startOpen?: boolean
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

  onUserInput: (text: string) => void = () => {}
  onFormSubmit: (values: Record<string, string>) => void = () => {}

  constructor(cfg: ChatbotConfig, engine: Engine, opts: WidgetUIOptions = {}) {
    this.cfg = cfg
    this.engine = engine

    const host = document.createElement('div')
    host.setAttribute('data-human-chatbot', '')
    const parent = opts.container || document.body
    const attach = () => {
      if (!host.isConnected) parent.appendChild(host)
    }
    attach()
    this.host = host

    // Only observe body when mounting to body — SPA hosts that re-render
    // document.body (e.g. Vue mount('body')) wipe our widget. For contained
    // previews the parent is stable, so no observer needed.
    if (!opts.container) {
      this.mo = new MutationObserver(() => attach())
      this.mo.observe(document.body, { childList: true })
    }

    this.root = host.attachShadow({ mode: 'open' })
    const style = document.createElement('style')
    style.textContent = baseCSS
    this.root.appendChild(style)

    const rootEl = document.createElement('div')
    rootEl.className = 'hb-root' + (opts.contained ? ' hb-contained' : '')
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
      img.src = cfg.styling.logoURL
      header.appendChild(img)
    }
    const title = document.createElement('span')
    title.textContent = cfg.styling.launcher.label || 'Assistant'
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
      if (ev.type === 'scenario') this.renderScenario(ev.scenario)
      else if (ev.type === 'message') this.appendMessage(ev.message)
      else if (ev.type === 'typing') this.setTyping(ev.on)
      else if (ev.type === 'agent_start') this.beginAgentStream()
      else if (ev.type === 'agent_delta') this.appendAgentDelta(ev.text)
      else if (ev.type === 'agent_end') this.endAgentStream()
      else if (ev.type === 'error') this.appendMessage({ role: 'system', content: ev.error })
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
    const cls = m.role === 'user' ? 'user' : m.role === 'agent' ? 'agent' : 'system'
    row.className = `hb-msg ${cls}`
    if (m.role === 'agent') {
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
  }

  private renderScenario(s: Scenario) {
    this.footerEl.className = 'hb-footer hb-hidden'
    this.footerEl.innerHTML = ''

    switch (s.type) {
      case 'static_message':
        return this.renderStaticMessage(s)
      case 'form':
        return this.renderForm(s)
      case 'conversation':
        return this.renderConversation()
    }
  }

  private renderStaticMessage(s: Scenario) {
    const cfg = s.config || {}
    const row = document.createElement('div')
    row.className = 'hb-msg agent'
    if (cfg.isMarkdown) {
      row.innerHTML = renderMarkdown(cfg.message || '')
    } else {
      row.textContent = cfg.message || ''
    }
    this.bodyEl.appendChild(row)
    this.bodyEl.scrollTop = this.bodyEl.scrollHeight
  }

  private renderForm(s: Scenario) {
    const cfg = s.config || {}
    const fields: Array<{ name: string; label: string; type?: string; required?: boolean }> =
      cfg.fields || []
    const form = document.createElement('form')
    form.className = 'hb-scenario-form'
    fields.forEach(f => {
      const row = document.createElement('div')
      row.className = 'hb-field'
      const lbl = document.createElement('label')
      lbl.textContent = f.label || f.name
      row.appendChild(lbl)
      const input = document.createElement('input')
      input.name = f.name
      input.type = f.type || 'text'
      if (f.required) input.required = true
      row.appendChild(input)
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
      this.onFormSubmit(values)
    })
    this.bodyEl.appendChild(form)
  }

  private renderConversation() {
    this.footerEl.className = 'hb-footer'
    this.footerEl.innerHTML = ''
    const input = document.createElement('textarea')
    input.className = 'hb-input'
    input.rows = 1
    input.placeholder =
      (this.engine.current?.config as { placeholder?: string } | undefined)?.placeholder ||
      'Type a message…'

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
      send.disabled = input.value.trim().length === 0
    }
    const submit = () => {
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
  }

}
