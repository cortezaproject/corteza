import { describe, it, expect } from 'vitest'
import { readdirSync, readFileSync, statSync } from 'node:fs'
import { fileURLToPath } from 'node:url'
import { dirname, join } from 'node:path'
import yaml from 'js-yaml'

// Restoring a soft-deleted resource is offered from a dozen screens. The word,
// the glyph and the severity are the same on every one of them, so the control
// reads as one thing rather than a dozen near-misses.

const SECTIONS = dirname(fileURLToPath(import.meta.url))
const LOCALE = join(SECTIONS, '../../../../../locale/en/human-webapp')

const LABEL = 'general.label.restore'
const ICON = 'pi pi-replay'
const SEVERITY = 'warn'

function walk(dir) {
  return readdirSync(dir).flatMap(entry => {
    const path = join(dir, entry)
    if (statSync(path).isDirectory()) return walk(path)
    return path.endsWith('.vue') ? [path] : []
  })
}

const sources = walk(SECTIONS).map(path => ({
  path,
  name: path.slice(SECTIONS.length + 1),
  text: readFileSync(path, 'utf8'),
}))

const bundle = readdirSync(LOCALE)
  .filter(f => f.endsWith('.yaml'))
  .reduce((acc, f) => {
    const flatten = (node, path) => {
      if (node && typeof node === 'object') {
        for (const [k, v] of Object.entries(node)) flatten(v, [...path, k])
      } else {
        acc[path.join('.')] = node
      }
    }
    flatten(yaml.load(readFileSync(join(LOCALE, f), 'utf8')) || {}, [f.slice(0, -5)])
    return acc
  }, {})

describe('restore controls', () => {
  it('is offered from more than one screen, so the sweep below means something', () => {
    const screens = sources.filter(s => s.text.includes(LABEL))
    expect(screens.length).toBeGreaterThan(15)
  })

  // A menu entry is an object literal; a toolbar entry is a <Button>. Both
  // carry the label and the icon as adjacent lines, so one window catches each.
  const withLabel = sources.filter(s => s.text.includes(LABEL))

  it.each(withLabel.map(s => [s.name, s.text]))('%s uses the restore glyph', (_name, text) => {
    for (const idx of indexesOf(text, LABEL)) {
      expect(windowAround(text, idx)).toContain(ICON)
    }
  })

  it.each(
    withLabel
      .filter(s => /<Button[^>]*\n[^>]*general\.label\.restore/.test(s.text))
      .map(s => [s.name, s.text]),
  )('%s gives the restore button warn severity', (_name, text) => {
    for (const idx of indexesOf(text, LABEL)) {
      const window = windowAround(text, idx)
      // Only the standalone buttons take a severity; menu entries carry none.
      if (window.includes('<Button')) {
        expect(window).toContain(`severity="${SEVERITY}"`)
      }
    }
  })

  it('never labels the action "Undelete" anywhere a user can read it', () => {
    const undeleteLabels = Object.entries(bundle)
      .filter(([, v]) => typeof v === 'string' && /^undelete$/i.test(v.trim()))
      .map(([k]) => k)
    expect(undeleteLabels).toEqual([])
  })

  it('spells the permission operation "Restore" too, in every service', () => {
    const permissionText = Object.entries(bundle)
      .filter(([k]) => k.startsWith('permissions.') && k.includes('.undelete.'))
      .map(([, v]) => v)
      .filter(v => typeof v === 'string')
    expect(permissionText.length).toBeGreaterThan(0)
    for (const text of permissionText) {
      expect(text).not.toMatch(/undelete/i)
    }
  })

  // The pre-existing bug this guards: agentic and chatbot offered restore to
  // everyone, and the server refuses it without the delete permission. The
  // deleted state is `deletedAt` on a stored resource and `deleted` on the
  // settings screen, whose providers are marked for deletion before saving.
  it('gates every restore control on a permission flag and on the deleted state', () => {
    for (const { name, text } of withLabel) {
      for (const idx of indexesOf(text, LABEL)) {
        const gate = gateFor(text, idx)
        expect(gate, `${name}: restore is not gated on the deleted state`).toMatch(
          /\bdeleted(At)?\b/,
        )
        expect(gate, `${name}: restore is not gated on a permission`).toMatch(
          /can(Delete|Undelete)\w*/,
        )
      }
    }
  })

  it('leaves no restore control on the glyphs that mean undo or reload', () => {
    for (const { name, text } of withLabel) {
      for (const idx of indexesOf(text, LABEL)) {
        const window = windowAround(text, idx)
        expect(window, name).not.toContain('pi-undo')
        expect(window, name).not.toContain('pi-refresh')
      }
    }
  })
})

// The condition that actually decides whether this control is offered — the
// element's own v-if for a toolbar button, the enclosing `if` for a menu entry.
// A neighbouring branch's condition must not count, or a control with no gate
// at all passes on its sibling's permission check.
function gateFor(text, idx) {
  // The label's own line too: the settings screen picks the word with a
  // ternary rather than branching around two pushes.
  let gate = ownGate(text, idx) + ' && ' + lineAt(text, idx)
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

function ownGate(text, idx) {
  const tagStart = text.lastIndexOf('<Button', idx)
  const tagEnd = text.indexOf('/>', idx)
  // A toolbar button: the label sits inside the element that carries the v-if.
  if (tagStart !== -1 && tagEnd !== -1 && tagStart < idx && text.indexOf('>', tagStart) > idx) {
    const tag = text.slice(tagStart, tagEnd)
    return [...tag.matchAll(/v-(?:else-)?if="([^"]*)"/g)].map(m => m[1]).join(' && ')
  }
  // A menu entry: every `if (…)` block still open at this point.
  return enclosingConditions(text, idx).join(' && ')
}

function enclosingConditions(text, idx) {
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
function conditionBefore(text, braceIdx) {
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

function lineAt(text, idx) {
  const start = text.lastIndexOf('\n', idx) + 1
  const end = text.indexOf('\n', idx)
  return text.slice(start, end === -1 ? text.length : end)
}

function indexesOf(text, needle) {
  const out = []
  let i = text.indexOf(needle)
  while (i !== -1) {
    out.push(i)
    i = text.indexOf(needle, i + 1)
  }
  return out
}

// The label and its icon sit within a few lines of each other in both shapes.
function windowAround(text, idx) {
  return text.slice(Math.max(0, idx - 300), idx + 300)
}
