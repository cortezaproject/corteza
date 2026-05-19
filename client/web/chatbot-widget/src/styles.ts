import type { Styling } from './types'

export const baseCSS = `
:host { all: initial; }
.hb-root {
  position: fixed; bottom: 16px; right: 16px; z-index: 2147483646;
  font-family: var(--hb-font, system-ui, sans-serif);
  font-size: var(--hb-font-base, 14px);
  color: var(--hb-text, #111);
}
.hb-root.hb-pos-bottom-left { right: auto; left: 16px; }
.hb-root.hb-pos-bottom-center { right: auto; left: 50%; transform: translateX(-50%); }
.hb-root.hb-pos-top-right { bottom: auto; top: 16px; }
.hb-root.hb-pos-top-left { bottom: auto; top: 16px; right: auto; left: 16px; }
.hb-root.hb-pos-top-center { bottom: auto; top: 16px; right: auto; left: 50%; transform: translateX(-50%); }
.hb-root.hb-pos-left-middle { bottom: auto; top: 50%; right: auto; left: 16px; transform: translateY(-50%); }
.hb-root.hb-pos-right-middle { bottom: auto; top: 50%; transform: translateY(-50%); }
.hb-root.hb-contained {
  position: absolute; bottom: 16px; right: 16px; z-index: 1;
}
.hb-root.hb-contained.hb-pos-bottom-left { right: auto; left: 16px; }
.hb-root.hb-contained.hb-pos-bottom-center { right: auto; left: 50%; transform: translateX(-50%); }
.hb-root.hb-contained.hb-pos-top-right { bottom: auto; top: 16px; }
.hb-root.hb-contained.hb-pos-top-left { bottom: auto; top: 16px; right: auto; left: 16px; }
.hb-root.hb-contained.hb-pos-top-center { bottom: auto; top: 16px; right: auto; left: 50%; transform: translateX(-50%); }
.hb-root.hb-contained.hb-pos-left-middle { bottom: auto; top: 50%; right: auto; left: 16px; transform: translateY(-50%); }
.hb-root.hb-contained.hb-pos-right-middle { bottom: auto; top: 50%; transform: translateY(-50%); }
.hb-root.hb-contained .hb-panel { position: absolute; }
.hb-root.hb-pos-bottom-left .hb-panel { right: auto; left: 0; }
.hb-root.hb-pos-bottom-center .hb-panel { right: auto; left: 50%; transform: translateX(-50%); }
.hb-root.hb-pos-top-right .hb-panel {
  top: calc(var(--hb-launch-size, 56px) + 12px); bottom: auto;
}
.hb-root.hb-pos-top-left .hb-panel {
  top: calc(var(--hb-launch-size, 56px) + 12px); bottom: auto; right: auto; left: 0;
}
.hb-root.hb-pos-top-center .hb-panel {
  top: calc(var(--hb-launch-size, 56px) + 12px); bottom: auto; right: auto; left: 50%; transform: translateX(-50%);
}
.hb-root.hb-pos-left-middle .hb-panel {
  top: calc(var(--hb-launch-size, 56px) + 12px); bottom: auto; right: auto; left: 0;
}
.hb-root.hb-pos-right-middle .hb-panel {
  top: calc(var(--hb-launch-size, 56px) + 12px); bottom: auto;
}
.hb-launcher {
  width: var(--hb-launch-size, 56px); height: var(--hb-launch-size, 56px);
  border-radius: var(--hb-launch-radius, 999px);
  background: var(--hb-primary, #09344E); color: var(--hb-primary-text, #fff);
  border: none; cursor: pointer; box-shadow: 0 4px 12px rgba(0,0,0,0.18);
  display: flex; align-items: center; justify-content: center;
  pointer-events: auto;
}
.hb-launcher img { width: 60%; height: 60%; object-fit: contain; display: block; }
.hb-launcher.hb-launcher-default .hb-launcher-icon { display: flex; width: 50%; height: 50%; align-items: center; justify-content: center; line-height: 1; }
.hb-launcher.hb-launcher-default .hb-launcher-icon svg { width: 100%; height: 100%; display: block; }
.hb-launcher.hb-launcher-has-label {
  width: auto; padding: 0 16px; gap: 8px; border-radius: var(--hb-launch-radius, 999px);
  display: flex; align-items: center; justify-content: center;
}
.hb-launcher.hb-launcher-has-label img,
.hb-launcher.hb-launcher-has-label .hb-launcher-icon { width: 22px; height: 22px; flex: none; }
.hb-launcher .hb-launcher-label { font-weight: 600; font-size: 14px; line-height: 1; }
.hb-panel {
  position: absolute; bottom: calc(var(--hb-launch-size, 56px) + 12px); right: 0;
  width: 360px; max-width: calc(100vw - 48px);
  height: 520px; max-height: calc(100vh - 120px);
  background: var(--hb-bg, #fff); border-radius: 12px;
  box-shadow: 0 20px 60px rgba(0,0,0,0.24); overflow: hidden;
  display: flex; flex-direction: column;
  pointer-events: auto;
}
.hb-header {
  position: relative;
  background: var(--hb-header, #09344E); color: var(--hb-header-text, #fff);
  padding: 12px 40px 12px 14px; font-weight: 600; font-size: var(--hb-font-heading, 16px);
  display: flex; align-items: center; gap: 8px;
}
.hb-header img { height: 24px; max-width: 140px; object-fit: contain; display: block; flex: none; }
.hb-close {
  position: absolute; top: 50%; right: 8px; transform: translateY(-50%);
  width: 28px; height: 28px; padding: 0;
  background: transparent; border: none; border-radius: 6px;
  color: inherit; cursor: pointer;
  display: flex; align-items: center; justify-content: center;
  opacity: 0.8;
}
.hb-close:hover { opacity: 1; background: rgba(255,255,255,0.12); }
.hb-close svg { width: 18px; height: 18px; }
.hb-panel { font-size: var(--hb-font-base, 14px); }
.hb-body { flex: 1; overflow-y: auto; padding: 12px; display: flex; flex-direction: column; gap: 10px; font-size: var(--hb-font-base, 14px); }
.hb-msg { max-width: 85%; padding: 8px 12px; border-radius: 14px; word-wrap: break-word; overflow-wrap: anywhere; font-size: var(--hb-font-base, 14px); }
.hb-msg.user { align-self: flex-end; background: var(--hb-user-bubble, #09344E); color: #fff; white-space: pre-wrap; }
.hb-msg.agent { align-self: flex-start; background: var(--hb-agent-bubble, #f4f4f5); color: var(--hb-text, #111); line-height: 1.45; }
.hb-msg.system { align-self: center; max-width: 90%; background: transparent; color: #9ca3af; font-size: var(--hb-font-small, 12px); font-style: italic; padding: 4px 8px; }
.hb-msg.agent a { color: var(--hb-primary, #09344E); text-decoration: underline; }
.hb-msg.agent code { background: rgba(0,0,0,0.06); padding: 1px 5px; border-radius: 4px; font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, monospace; font-size: 0.92em; }
.hb-msg.agent pre { background: rgba(0,0,0,0.06); padding: 8px 10px; border-radius: 8px; overflow-x: auto; margin: 6px 0; }
.hb-msg.agent pre code { background: none; padding: 0; border-radius: 0; font-size: 0.88em; }
.hb-msg.agent ul { margin: 6px 0; padding-left: 20px; }
.hb-msg.agent strong { font-weight: 600; }
.hb-footer { border-top: 1px solid #e5e7eb; padding: 8px; display: flex; flex-wrap: wrap; gap: 6px; align-items: flex-end; font-size: var(--hb-font-base, 14px); }
.hb-input {
  flex: 1; border: 1px solid #d1d5db; border-radius: 8px; padding: 8px;
  font-family: inherit; color: inherit; background: var(--hb-bg, #fff);
  font-size: var(--hb-font-base, 14px);
  resize: none; max-height: 120px; overflow-y: auto; line-height: 1.4;
}
.hb-input:focus { outline: none; border-color: var(--hb-primary, #09344E); }
.hb-send {
  background: var(--hb-primary, #09344E); color: var(--hb-primary-text, #fff);
  border: none; border-radius: 8px; padding: 0; cursor: pointer;
  width: 36px; height: 36px; flex-shrink: 0;
  display: flex; align-items: center; justify-content: center;
  font-size: var(--hb-font-base, 14px);
}
.hb-send svg { width: 18px; height: 18px; }
.hb-send:disabled { opacity: 0.45; cursor: not-allowed; }
.hb-sr-only { position: absolute; width: 1px; height: 1px; padding: 0; margin: -1px; overflow: hidden; clip: rect(0,0,0,0); border: 0; }
.hb-sugg { display: flex; flex-wrap: wrap; gap: 6px; }
.hb-sugg button { background: #f4f4f5; border: 1px solid #e5e7eb; border-radius: 999px; padding: 6px 10px; cursor: pointer; font-family: inherit; font-size: var(--hb-font-base, 14px); }
.hb-form { display: flex; flex-direction: column; gap: 8px; }
.hb-form label { font-size: var(--hb-font-small, 12px); color: #6b7280; }
.hb-form input, .hb-form textarea {
  border: 1px solid #d1d5db; border-radius: 8px; padding: 6px 8px; font-family: inherit; font-size: var(--hb-font-base, 14px); background: var(--hb-bg, #fff); color: inherit;
}
.hb-scenario-form {
  display: flex; flex-direction: column; gap: 10px;
  border: 1px solid #e5e7eb; border-radius: 12px; padding: 12px; background: var(--hb-bg, #fff);
}
.hb-scenario-form .hb-field { display: flex; flex-direction: column; gap: 4px; align-items: stretch; }
.hb-scenario-form .hb-field label {
  text-align: left; font-size: var(--hb-font-small, 12px); color: #6b7280; font-weight: 500;
}
.hb-scenario-form .hb-field input, .hb-scenario-form .hb-field textarea {
  border: 1px solid #d1d5db; border-radius: 8px; padding: 6px 8px; font-family: inherit; font-size: var(--hb-font-base, 14px); background: var(--hb-bg, #fff); color: inherit;
}
.hb-scenario-form .hb-field input:focus, .hb-scenario-form .hb-field textarea:focus {
  outline: none; border-color: var(--hb-primary, #09344E);
}
.hb-form-actions { display: flex; justify-content: flex-end; }
.hb-submit {
  background: var(--hb-primary, #09344E); color: var(--hb-primary-text, #fff);
  border: none; border-radius: 8px; padding: 6px 14px; cursor: pointer; font-family: inherit; font-size: var(--hb-font-base, 14px);
}
.hb-submit:hover { opacity: 0.92; }
.hb-typing { align-self: flex-start; padding: 10px 14px; display: inline-flex; gap: 4px; background: var(--hb-agent-bubble, #f4f4f5); border-radius: 14px; }
.hb-typing span {
  width: 6px; height: 6px; border-radius: 999px; background: #9ca3af;
  animation: hb-bounce 1.1s infinite ease-in-out;
}
.hb-typing span:nth-child(2) { animation-delay: 0.15s; }
.hb-typing span:nth-child(3) { animation-delay: 0.30s; }
@keyframes hb-bounce {
  0%, 80%, 100% { transform: translateY(0); opacity: 0.4; }
  40%           { transform: translateY(-4px); opacity: 1; }
}
.hb-info { font-size: var(--hb-font-small, 12px); color: #6b7280; text-align: center; padding: 10px; }
.hb-conv-actions { display: flex; gap: 6px; flex-wrap: wrap; flex-basis: 100%; width: 100%; padding-top: 6px; order: 2; }
.hb-footer .hb-input { order: 1; }
.hb-footer .hb-send { order: 1; }
.hb-action {
  background: transparent; color: var(--hb-primary, #09344E);
  border: 1px solid #d1d5db; border-radius: 999px;
  padding: 4px 10px; cursor: pointer; font-family: inherit;
  font-size: var(--hb-font-small, 12px);
}
.hb-action:hover { background: rgba(0,0,0,0.04); }
.hb-action:disabled { opacity: 0.45; cursor: not-allowed; }
.hb-handoff-badge {
  background: #fef3c7; color: #78350f;
  border-bottom: 1px solid #fde68a;
  padding: 8px 12px; font-size: var(--hb-font-small, 12px);
  display: flex; align-items: center; gap: 8px;
}
.hb-handoff-badge .hb-action { border-color: #fcd34d; color: #78350f; }
.hb-msg.agent .hb-operator-tag {
  font-size: var(--hb-font-small, 12px); font-weight: 600;
  color: #047857; margin-bottom: 2px;
}
.hb-msg.agent .hb-msg-body { white-space: pre-wrap; }
.hb-field-error {
  font-size: var(--hb-font-small, 12px); color: #b91c1c; padding-top: 2px;
}
.hb-hidden { display: none !important; }
`

export function applyStyling(root: HTMLElement, s: Styling) {
  const set = (k: string, v: string | undefined) => {
    if (v) root.style.setProperty(k, v)
  }
  set('--hb-primary', s.colors.primary)
  set('--hb-primary-text', s.colors.primaryText)
  set('--hb-header', s.colors.header)
  set('--hb-header-text', s.colors.headerText)
  set('--hb-bg', s.colors.background)
  set('--hb-text', s.colors.text)
  set('--hb-user-bubble', s.colors.userBubble)
  set('--hb-agent-bubble', s.colors.agentBubble)
  set('--hb-font', s.fontFamily)
  set('--hb-font-base', s.fontSizes.base)
  set('--hb-font-small', s.fontSizes.small)
  set('--hb-font-heading', s.fontSizes.heading)
  set('--hb-launch-size', s.launcher.size)
  set('--hb-launch-radius', s.launcher.shape === 'square' ? '12px' : '999px')
}
