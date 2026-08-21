import { describe, it, expect } from 'vitest'
import { readdirSync, readFileSync } from 'node:fs'
import { fileURLToPath } from 'node:url'
import { dirname, join } from 'node:path'
import yaml from 'js-yaml'
import { vueSources, gateFor, indexesOf, lineAt, windowAround } from './control-gates.js'

// Restoring a soft-deleted resource is offered from a dozen screens. The word,
// the glyph and the severity are the same on every one of them, so the control
// reads as one thing rather than a dozen near-misses.

const SECTIONS = dirname(fileURLToPath(import.meta.url))
const LOCALE = join(SECTIONS, '../../../../../locale/en/human-webapp')

const LABEL = 'general.label.restore'
const ICON = 'pi pi-replay'
const SEVERITY = 'warn'

const sources = vueSources(SECTIONS)

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
