// Global trim-on-blur.
//
// Leading/trailing whitespace in text inputs is never meaningful in our config
// and form fields (expressions, handles, names, labels, …) and only causes
// subtle bugs — e.g. a value that reads back wrong or fails a comparison. This
// plugin installs a single document-level `focusout` listener that trims the
// value of text inputs/textareas when they lose focus, then re-dispatches an
// `input` event so Vue's v-model (and PrimeVue's internal binding) pick up the
// trimmed value.
//
// Notes:
//  - Only leading/trailing whitespace is removed; internal whitespace (e.g.
//    multi-line expression bodies) is preserved.
//  - Fires on blur, not per keystroke, so typing spaces mid-edit is unaffected.
//  - `type=password` is excluded so credentials/secrets are never silently
//    altered. Opt any field out with the `data-no-trim` attribute.
//  - contenteditable/rich-text editors are not <input>/<textarea>, so they are
//    naturally excluded.

const TRIMMABLE_INPUT_TYPES = new Set(['text', 'search', 'email', 'url', 'tel', ''])

function shouldTrim(el) {
  if (!el || el.disabled || el.readOnly) return false
  if (el.hasAttribute('data-no-trim')) return false

  const tag = el.tagName
  if (tag === 'TEXTAREA') return true
  if (tag !== 'INPUT') return false

  // el.type normalizes unknown/absent types to 'text'
  return TRIMMABLE_INPUT_TYPES.has((el.getAttribute('type') || 'text').toLowerCase())
}

function onFocusOut(event) {
  const el = event.target
  if (!shouldTrim(el)) return

  const trimmed = el.value.trim()
  if (trimmed === el.value) return

  el.value = trimmed
  // Let v-model / PrimeVue observe the trimmed value. Only `input` is
  // dispatched: native `change` already fires on blur, so re-dispatching it
  // could double-trigger `@change` handlers (e.g. save-on-change).
  el.dispatchEvent(new Event('input', { bubbles: true }))
}

export default {
  install() {
    // `focusout` bubbles (unlike `blur`), so a single capture at the document
    // root covers every field in the app, present and future.
    document.addEventListener('focusout', onFocusOut, true)
  },
}
