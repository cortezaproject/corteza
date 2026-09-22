import { createContext, runInContext } from 'node:vm'
import { describe, expect, it, vi } from 'vitest'
import { BRIDGE_SCRIPT } from './bridge'
import {
  allowModule,
  buildOuterDocument,
  capLimit,
  dispatch,
  hostScriptSource,
  reshapeRecord,
} from './host'

const META = { namespace: 'agent-sandbox', modules: ['agent-contact'] }

const context = (meta = META) => ({
  compose: {
    recordList: vi.fn().mockResolvedValue({ set: [], filter: {} }),
    recordRead: vi.fn(),
    recordReport: vi.fn(),
  },
  meta,
  namespaceID: '1',
  moduleIDs: { 'agent-contact': '2' },
})

describe('record reshaping', () => {
  it('keys values by field name', () => {
    expect(
      reshapeRecord({
        recordID: '1',
        ownedBy: '7',
        createdAt: 'then',
        values: [
          { name: 'name', value: 'Alice Novak' },
          { name: 'company', value: 'Acme' },
        ],
      }),
    ).toEqual({
      recordID: '1',
      ownedBy: '7',
      createdAt: 'then',
      values: { name: 'Alice Novak', company: 'Acme' },
    })
  })

  it('collects a repeated field into an array', () => {
    const { values } = reshapeRecord({
      values: [
        { name: 'tag', value: 'one' },
        { name: 'tag', value: 'two' },
        { name: 'tag', value: 'three' },
      ],
    })
    expect(values).toEqual({ tag: ['one', 'two', 'three'] })
  })

  it('leaves a field with nothing in it out entirely', () => {
    // The app-facing shape has no empty strings in it: a field the record does
    // not carry and one carrying '' read the same to a page.
    const { values } = reshapeRecord({
      values: [
        { name: 'name', value: 'Alice Novak' },
        { name: 'company', value: '' },
        { name: 'email', value: null },
      ],
    })
    expect(values).toEqual({ name: 'Alice Novak' })
  })
})

describe('module allowlist', () => {
  it('admits a module the source declares', () => {
    expect(allowModule(META, 'agent-contact')).toBeNull()
  })

  it('refuses one it does not, naming it', () => {
    expect(allowModule(META, 'Lead')).toBe('module "Lead" is not declared for this app')
  })

  it('refuses everything when the source declares nothing', () => {
    expect(allowModule(undefined, 'agent-contact')).toBe(
      'module "agent-contact" is not declared for this app',
    )
  })
})

describe('limit cap', () => {
  it('caps what an app asks for', () => {
    expect(capLimit(5000)).toBe(500)
    expect(capLimit(50)).toBe(50)
  })

  it('leaves the server default in place when nothing sensible is asked', () => {
    expect(capLimit(undefined)).toBeUndefined()
    expect(capLimit(0)).toBeUndefined()
    expect(capLimit('lots')).toBeUndefined()
  })
})

describe('operation dispatch', () => {
  it('refuses an undeclared module before the API is touched', async () => {
    const ctx = context()
    await expect(dispatch('records.list', { module: 'Lead' }, ctx)).rejects.toThrow(
      'module "Lead" is not declared for this app',
    )
    expect(ctx.compose.recordList).not.toHaveBeenCalled()
  })

  it('refuses an operation that is not in the table', async () => {
    const ctx = context()
    await expect(dispatch('records.create', {}, ctx)).rejects.toThrow(
      'operation "records.create" is not available to an app',
    )
  })

  it('passes a declared module through as its id, with the limit capped', async () => {
    const ctx = context()
    await dispatch('records.list', { module: 'agent-contact', limit: 5000 }, ctx)
    expect(ctx.compose.recordList).toHaveBeenCalledWith(
      expect.objectContaining({ namespaceID: '1', moduleID: '2', limit: 500 }),
    )
  })
})

describe('outer document', () => {
  const build = source =>
    buildOuterDocument({
      source,
      bridgeScript: BRIDGE_SCRIPT,
      hostScript: hostScriptSource({ origin: 'https://human.test' }),
    })

  it('denies itself child frames, which is what pins the app in place', () => {
    expect(build('<p>hello</p>')).toContain(
      `<meta http-equiv="Content-Security-Policy" content="frame-src 'none'">`,
    )
  })

  it('escapes a closing script tag in the source it carries', () => {
    const doc = build('<script>alert(1)</script>')
    expect(doc).not.toContain('alert(1)</script>')
    expect(doc).toContain(String.raw`alert(1)<\/script>`)
  })
})

describe('in-page bridge', () => {
  it('leaves one object however many copies of it the page runs', () => {
    // An app written to the MCP skill pastes the snippet itself and we prefix
    // it too, so a second evaluation has to be a no-op rather than a
    // SyntaxError that takes the page with it.
    const sandbox = {
      setTimeout: () => 0,
      clearTimeout: () => {},
      addEventListener: () => {},
      parent: { postMessage: () => {} },
    }
    sandbox.window = sandbox
    createContext(sandbox)

    runInContext(BRIDGE_SCRIPT, sandbox)
    const first = sandbox.human
    runInContext(BRIDGE_SCRIPT, sandbox)

    expect(first).toBeTypeOf('object')
    expect(sandbox.human).toBe(first)
  })
})
