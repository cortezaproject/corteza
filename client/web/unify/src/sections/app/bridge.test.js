import { createContext, runInContext } from 'node:vm'
import { describe, expect, it, vi } from 'vitest'
import { BRIDGE_SCRIPT } from './bridge'
import {
  allowModule,
  BRIDGE_VERSION,
  bridgeVersion,
  buildOuterDocument,
  capLimit,
  dispatch,
  hostScriptSource,
  labelFieldOf,
  recordLabels,
  referenceTargets,
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

describe('bridge contract version', () => {
  it('reads the number the handshake in the page sends', () => {
    expect(bridgeVersion(`parent.postMessage({ type: 'human:hello', v: 1 }, '*')`)).toBe(1)
    expect(bridgeVersion(`parent.postMessage({ type: "human:hello", v:2 }, '*')`)).toBe(2)
  })

  it('treats a handshake that names no number as the first contract', () => {
    expect(bridgeVersion(`parent.postMessage({ type: 'human:hello' }, '*')`)).toBe(1)
  })

  it('gives a page with no bridge of its own the current contract', () => {
    expect(bridgeVersion('<h1>hi</h1>')).toBe(BRIDGE_VERSION)
  })

  it('ships the copy the shell prefixes at the current contract', () => {
    expect(bridgeVersion(BRIDGE_SCRIPT)).toBe(BRIDGE_VERSION)
  })
})

describe('typed values (contract 2)', () => {
  const fields = [
    { name: 'active', kind: 'Bool' },
    { name: 'archived', kind: 'Bool' },
    { name: 'amount', kind: 'Number' },
    { name: 'scores', kind: 'Number', multi: true },
    { name: 'name', kind: 'String' },
  ]
  const record = {
    recordID: '1',
    values: [
      { name: 'active', value: '1' },
      { name: 'archived', value: null },
      { name: 'amount', value: '12.50' },
      { name: 'scores', value: '3' },
      { name: 'scores', value: '4' },
      { name: 'name', value: 'Alice' },
    ],
  }

  it('makes a Bool true or false and always present, and a Number a number', () => {
    expect(reshapeRecord(record, fields, 2).values).toEqual({
      active: true,
      archived: false,
      amount: 12.5,
      scores: [3, 4],
      name: 'Alice',
    })
  })

  it('gives a Bool the store holds nothing for false, not absent', () => {
    expect(reshapeRecord({ values: [] }, fields, 2).values).toEqual({
      active: false,
      archived: false,
    })
  })

  it('leaves the first contract exactly as it was', () => {
    expect(reshapeRecord(record, fields, 1).values).toEqual({
      active: '1',
      amount: '12.50',
      scores: ['3', '4'],
      name: 'Alice',
    })
  })
})

describe('record references', () => {
  const company = { name: 'company', kind: 'Record', options: { moduleID: '9' } }

  it('picks the named label field, else the first field', () => {
    const module = { fields: [{ name: 'name' }, { name: 'code' }] }
    expect(labelFieldOf(module, 'code').name).toBe('code')
    expect(labelFieldOf(module, '').name).toBe('name')
    expect(labelFieldOf(module, 'gone').name).toBe('name')
    expect(labelFieldOf(null, 'x')).toBe(null)
  })

  it('collects what each record points at, by target module', () => {
    const targets = referenceTargets(
      [company],
      [
        { values: [{ name: 'company', value: '100' }] },
        { values: [{ name: 'company', value: '0' }] },
        {
          values: [
            { name: 'company', value: '100' },
            { name: 'other', value: '5' },
          ],
        },
      ],
    )
    expect(Object.keys(targets)).toEqual(['9'])
    expect([...targets['9'].ids]).toEqual(['100'])
  })

  const compose = (modules, records) => ({
    moduleRead: vi.fn(({ moduleID }) =>
      modules[moduleID] ? Promise.resolve(modules[moduleID]) : Promise.reject(new Error('no')),
    ),
    recordList: vi.fn(({ moduleID, recordID }) =>
      Promise.resolve({
        set: (records[moduleID] || []).filter(r => recordID.includes(r.recordID)),
      }),
    ),
  })

  it('labels a referenced record by its label field, in one call per module', async () => {
    const api = compose(
      { 9: { fields: [{ name: 'name', kind: 'String' }] } },
      { 9: [{ recordID: '100', values: [{ name: 'name', value: 'Acme' }] }] },
    )
    const labels = await recordLabels(
      api,
      '1',
      referenceTargets([company], [{ values: [{ name: 'company', value: '100' }] }]),
    )
    expect(labels).toEqual({ 100: 'Acme' })
    expect(api.recordList).toHaveBeenCalledTimes(1)
    expect(api.recordList.mock.calls[0][0]).toMatchObject({
      moduleID: '9',
      recordID: ['100'],
      limit: 1,
    })
  })

  it('follows a label that is itself a reference one level further', async () => {
    const api = compose(
      {
        9: { fields: [{ name: 'owner', kind: 'Record', options: { moduleID: '8' } }] },
        8: { fields: [{ name: 'title', kind: 'String' }] },
      },
      {
        9: [{ recordID: '100', values: [{ name: 'owner', value: '200' }] }],
        8: [{ recordID: '200', values: [{ name: 'title', value: 'Globex' }] }],
      },
    )
    const labels = await recordLabels(
      api,
      '1',
      referenceTargets([company], [{ values: [{ name: 'company', value: '100' }] }]),
    )
    expect(labels).toEqual({ 100: 'Globex' })
  })

  it('leaves a record the viewer cannot read unlabelled, and does not fail', async () => {
    const api = compose({}, {})
    const labels = await recordLabels(
      api,
      '1',
      referenceTargets([company], [{ values: [{ name: 'company', value: '100' }] }]),
    )
    expect(labels).toEqual({})
    expect(api.recordList).not.toHaveBeenCalled()
  })

  it('puts record labels beside user labels in what records.list returns', async () => {
    const ctx = context()
    ctx.version = 2
    ctx.fields = () => [company, { name: 'active', kind: 'Bool' }]
    ctx.refs = vi.fn().mockResolvedValue({ 7: 'Dev Agent' })
    ctx.compose.moduleRead = vi
      .fn()
      .mockResolvedValue({ fields: [{ name: 'name', kind: 'String' }] })
    ctx.compose.recordList = vi
      .fn()
      .mockResolvedValueOnce({
        set: [{ recordID: '1', ownedBy: '7', values: [{ name: 'company', value: '100' }] }],
        filter: {},
      })
      .mockResolvedValueOnce({
        set: [{ recordID: '100', values: [{ name: 'name', value: 'Acme' }] }],
      })

    const out = await dispatch('records.list', { module: 'agent-contact' }, ctx)
    expect(out.records[0].values).toEqual({ company: '100', active: false })
    expect(out.refs).toEqual({ 7: 'Dev Agent', 100: 'Acme' })
  })
})
