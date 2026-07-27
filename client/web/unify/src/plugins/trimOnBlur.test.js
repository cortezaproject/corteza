import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'
import TrimOnBlurPlugin from './trimOnBlur'

describe('trimOnBlur plugin', () => {
  beforeEach(() => {
    TrimOnBlurPlugin.install()
  })

  afterEach(() => {
    document.body.innerHTML = ''
  })

  function mount(html) {
    const wrap = document.createElement('div')
    wrap.innerHTML = html
    document.body.appendChild(wrap)
    return wrap.firstElementChild
  }

  function blur(el) {
    el.dispatchEvent(new Event('focusout', { bubbles: true }))
  }

  it('trims leading/trailing whitespace on a text input at blur', () => {
    const el = mount('<input type="text" value="  hello  " />')
    blur(el)
    expect(el.value).toBe('hello')
  })

  it('defaults to trimming inputs with no type', () => {
    const el = mount('<input value="  x  " />')
    blur(el)
    expect(el.value).toBe('x')
  })

  it('preserves internal whitespace in a textarea', () => {
    const el = mount('<textarea>\n  a +\n  b  \n</textarea>')
    blur(el)
    expect(el.value).toBe('a +\n  b')
  })

  it('dispatches an input event so v-model updates', () => {
    const el = mount('<input type="text" value=" v " />')
    const spy = vi.fn()
    el.addEventListener('input', spy)
    blur(el)
    expect(spy).toHaveBeenCalledOnce()
    expect(el.value).toBe('v')
  })

  it('does not dispatch input when nothing changes', () => {
    const el = mount('<input type="text" value="clean" />')
    const spy = vi.fn()
    el.addEventListener('input', spy)
    blur(el)
    expect(spy).not.toHaveBeenCalled()
  })

  it('leaves password fields untouched', () => {
    const el = mount('<input type="password" value="  pw  " />')
    blur(el)
    expect(el.value).toBe('  pw  ')
  })

  it('honors data-no-trim opt-out', () => {
    const el = mount('<input type="text" data-no-trim value="  keep  " />')
    blur(el)
    expect(el.value).toBe('  keep  ')
  })

  it('ignores disabled and readonly fields', () => {
    const dis = mount('<input type="text" disabled value="  a  " />')
    const ro = mount('<input type="text" readonly value="  b  " />')
    blur(dis)
    blur(ro)
    expect(dis.value).toBe('  a  ')
    expect(ro.value).toBe('  b  ')
  })

  it('ignores non-text input types', () => {
    const el = mount('<input type="checkbox" />')
    // value getter on checkbox is "on"; just ensure no throw and no change
    expect(() => blur(el)).not.toThrow()
  })
})
