import { readdirSync, readFileSync, statSync } from 'node:fs'
import { join } from 'node:path'

// Which condition actually decides whether a control is offered, and every .vue
// file that could carry one. Shared by the sweeps that assert a control is
// gated — restore-controls.test.js and permission-controls.test.js.

export function walk(dir) {
  return readdirSync(dir).flatMap(entry => {
    const path = join(dir, entry)
    if (statSync(path).isDirectory()) return walk(path)
    return path.endsWith('.vue') ? [path] : []
  })
}

export function vueSources(root) {
  return walk(root).map(path => ({
    path,
    name: path.slice(root.length + 1),
    text: readFileSync(path, 'utf8'),
  }))
}

// The condition that actually decides whether this control is offered — the
// element's own v-if for a toolbar button, the enclosing `if` for a menu entry.
// A neighbouring branch's condition must not count, or a control with no gate
// at all passes on its sibling's permission check.
export function gateFor(text, idx) {
  // The label's own line too: the settings screen picks the word with a
  // ternary rather than branching around two pushes.
  let gate = [ownGate(text, idx), ...ancestorGates(text, idx), lineAt(text, idx)].join(' && ')
  const seen = new Set()
  // A gate can name a computed that names another one, so follow the chain
  // rather than only its first link.
  for (let hop = 0; hop < 3; hop++) {
    let grew = ''
    for (const [id] of gate.matchAll(/[A-Za-z_$][\w$]*/g)) {
      if (seen.has(id)) continue
      seen.add(id)
      // Only a computed counts. Folding in every identifier would drag the
      // whole menu-building block along with `const items = []`, and a
      // sibling branch's permission check would pass for a control with none.
      const def = new RegExp(
        `(?:const\\s+${id}\\s*=\\s*computed\\(|^\\s*${id}\\s*\\(\\)\\s*\\{)[\\s\\S]{0,400}`,
        'm',
      ).exec(text)
      if (def) grew += '\n' + def[0]
    }
    if (!grew) break
    gate += grew
  }
  return gate
}

export function ownGate(text, idx) {
  const tag = enclosingTag(text, idx)
  // A toolbar control: the attributes sit on the element the label is inside.
  // `:disabled` counts as a gate — a control the user can see but not press is
  // still gated, and compose spells it that way.
  if (tag !== null) {
    return [...tag.matchAll(/(?:v-(?:else-)?if|:disabled|:visible)="([^"]*)"/g)]
      .map(m => m[1])
      .join(' && ')
  }
  // A menu entry: every `if (…)` block still open at this point.
  return enclosingConditions(text, idx).join(' && ')
}

// The v-if of every element still open at this offset. A control is often
// gated by the toolbar <div> that wraps it rather than on the control itself,
// and a sweep that only reads the element's own attributes calls that ungated.
export function ancestorGates(text, idx) {
  const out = []
  const stack = []
  const tag = /<(\/?)([A-Za-z][\w.-]*)((?:"[^"]*"|'[^']*'|[^>"'])*?)(\/?)>/g
  let m
  while ((m = tag.exec(text)) !== null) {
    if (m.index >= idx) break
    const [, closing, name, attrs, selfClose] = m
    if (closing) {
      const at = stack.map(f => f.name).lastIndexOf(name)
      if (at !== -1) stack.length = at
    } else if (!selfClose && !VOID_TAGS.has(name)) {
      stack.push({ name, attrs })
    }
  }
  for (const f of stack) {
    for (const a of f.attrs.matchAll(/v-(?:else-)?if="([^"]*)"/g)) out.push(a[1])
  }
  return out
}

const VOID_TAGS = new Set(['br', 'hr', 'img', 'input', 'link', 'meta', 'source'])

// The start tag the given offset sits inside, attributes only. Null when the
// offset is in <script> or between elements — a menu entry has no tag, and
// mistaking an enclosing <div> for one lets a control pass on a layout condition.
export function enclosingTag(text, idx) {
  const open = text.lastIndexOf('<', idx)
  if (open === -1) return null
  const close = text.indexOf('>', open)
  // The offset must fall inside this tag's own attribute list.
  if (close === -1 || close < idx) return null
  if (!/^<[A-Za-z]/.test(text.slice(open, open + 2))) return null
  return text.slice(open, close)
}

export function enclosingConditions(text, idx) {
  const out = []
  let depth = 0
  for (let i = idx; i >= 0; i--) {
    const c = text[i]
    if (c === '}') depth++
    else if (c === '{') {
      if (depth > 0) {
        depth--
        continue
      }
      const cond = conditionBefore(text, i)
      if (cond !== null) out.push(cond)
    }
  }
  return out
}

// The `if (…)` immediately before an opening brace, found by matching that
// condition's own parentheses. A regex would run back to an earlier `if` in
// the same neighbourhood and swallow a sibling branch's condition with it.
export function conditionBefore(text, braceIdx) {
  let i = braceIdx - 1
  while (i >= 0 && /\s/.test(text[i])) i--
  if (text[i] !== ')') return null
  let depth = 0
  const close = i
  for (; i >= 0; i--) {
    if (text[i] === ')') depth++
    else if (text[i] === '(') {
      depth--
      if (depth === 0) break
    }
  }
  if (i < 0) return null
  let j = i - 1
  while (j >= 0 && /\s/.test(text[j])) j--
  if (text.slice(j - 1, j + 1) !== 'if') return null
  return text.slice(i + 1, close)
}

export function lineAt(text, idx) {
  const start = text.lastIndexOf('\n', idx) + 1
  const end = text.indexOf('\n', idx)
  return text.slice(start, end === -1 ? text.length : end)
}

export function indexesOf(text, needle) {
  const out = []
  let i = text.indexOf(needle)
  while (i !== -1) {
    out.push(i)
    i = text.indexOf(needle, i + 1)
  }
  return out
}

// The label and its icon sit within a few lines of each other in both shapes.
export function windowAround(text, idx) {
  return text.slice(Math.max(0, idx - 300), idx + 300)
}
