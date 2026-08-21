import { describe, it, expect } from 'vitest'
import { fileURLToPath } from 'node:url'
import { dirname } from 'node:path'
import { vueSources, gateFor, indexesOf } from './control-gates.js'

// A control the server will refuse is not offered. Two mechanisms answer that
// question and they are not interchangeable: a component operation through
// useRBACStore.can() for anything not yet created, and a can* flag on the
// resource's own payload for anything that exists.

const SECTIONS = dirname(fileURLToPath(import.meta.url))
const sources = vueSources(SECTIONS)

// A create control names no resource, so only the component operation can gate
// it. `<x>.create` and `<x>.manage` are the operations that exist.
const CREATE_GATE =
  /\bcan(Create|Add)\w*|\.can\(\s*['"][a-z]+\/['"]\s*,\s*['"][\w-]+\.(create|manage)['"]/

// An update control has the resource in hand, so it gates on that resource's
// own flag — canUpdate*, canManage*, or the block's own canSave/canEdit computed.
const UPDATE_GATE =
  /\bcan(Update|Manage|Save|Edit|Grant|Suspend|Unsuspend|Publish|Revise|Execute|Write)\w*/

// Views that legitimately carry none of this.
function skip(name) {
  return (
    name.includes('.test.') ||
    // the project section gates wholesale on the project role preset, not RBAC
    name.startsWith('project/')
  )
}

const views = sources.filter(s => !skip(s.name))

describe('permission-gated controls', () => {
  it('sweeps enough views to mean something', () => {
    expect(views.length).toBeGreaterThan(100)
  })

  // The pre-existing bug this guards: ten admin lists drew "New X" for a user
  // with no create right, and seven editors drew Save for one who could not
  // update. Both silently failed at the server.
  const creators = views.filter(
    s => /<(Button|CRouterLinkButton)\b/.test(s.text) && /list\.new|list\.add-button/.test(s.text),
  )

  it.each(creators.map(s => [s.name, s.text]))('%s gates its create control', (name, text) => {
    for (const idx of indexesOf(text, 'list.new').concat(indexesOf(text, 'list.add-button'))) {
      // only the control itself, not the resourceSingle translation entry
      const line = text.slice(text.lastIndexOf('\n', idx) + 1, text.indexOf('\n', idx))
      if (!line.includes(':label=') && !line.includes('label:')) continue
      expect(gateFor(text, idx), `${name}: create control is not gated on a permission`).toMatch(
        CREATE_GATE,
      )
    }
  })

  const savers = views.filter(s => s.text.includes('type="submit"'))

  it.each(savers.map(s => [s.name, s.text]))('%s gates its submit control', (name, text) => {
    for (const idx of indexesOf(text, 'type="submit"')) {
      const gate = gateFor(text, idx)
      // A create-only form has nothing to check a flag against.
      if (!/\bisEdit\b|\bisCreate\b/.test(text)) continue
      expect(gate, `${name}: Save is offered regardless of the update permission`).toMatch(
        UPDATE_GATE,
      )
    }
  })

  const granters = views.filter(s => s.text.includes('<CPermissionsButton'))

  it.each(granters.map(s => [s.name, s.text]))('%s gates its permissions control', (name, text) => {
    for (const idx of indexesOf(text, '<CPermissionsButton')) {
      expect(
        gateFor(text, idx),
        `${name}: the permissions button is not gated on canGrant`,
      ).toMatch(/canGrant/)
    }
  })
})
