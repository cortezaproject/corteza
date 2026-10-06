import { createContext, runInContext } from 'node:vm'
import { describe, expect, it, vi } from 'vitest'
import { BRIDGE_SCRIPT } from './bridge'
import {
  cspInner,
  automationInput,
  navigationTarget,
  allowModule,
  describeField,
  dimensionRefs,
  downloadName,
  MAX_DOWNLOAD,
  allowWrite,
  toStoreValues,
  BRIDGE_VERSION,
  bridgeVersion,
  buildOuterDocument,
  capLimit,
  dispatch,
  hostScriptSource,
  INERT_NAVIGATION,
  NO_PEER_CONNECTIONS,
  REPORT_FAILURE,
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
    await expect(dispatch('records.purge', {}, ctx)).rejects.toThrow(
      'operation "records.purge" is not available to an app',
    )
  })

  it('answers context with where the host shows the app, and nothing where it says nothing', async () => {
    const where = { namespaceID: '1', pageID: '2', recordID: '3' }
    expect(await dispatch('context', {}, { ...context(), context: () => where })).toEqual(where)
    expect(await dispatch('context', {}, context())).toEqual({})
  })

  // The API calls the field a report groups by `dimension`, which is what an
  // author reading the tools is shown; a bridge taking only the plural dropped
  // it silently and the server answered "field dimensions is empty".
  it('groups a report by the field the API names, singular or plural', async () => {
    for (const args of [
      { module: 'agent-contact', dimension: 'status' },
      { module: 'agent-contact', dimensions: 'status' },
    ]) {
      const ctx = context()
      await dispatch('records.report', args, ctx)
      expect(ctx.compose.recordReport).toHaveBeenCalledWith(
        expect.objectContaining({ dimensions: 'status' }),
      )
    }
  })

  it('refuses a report that names no field to group by', async () => {
    const ctx = context()
    await expect(dispatch('records.report', { module: 'agent-contact' }, ctx)).rejects.toThrow(
      'a report needs a `dimension`: the field to group by',
    )
    expect(ctx.compose.recordReport).not.toHaveBeenCalled()
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

  it('gives the app its own address as base, and cancels navigation before any app script runs', () => {
    const doc = buildOuterDocument({
      source: '<script>app()</script>',
      hostScript: '',
      bridgeScript: 'BRIDGE',
    })
    const inner = JSON.parse(
      doc.match(/window.__humanInner = (".*?")<\/script>/s)[1].replace(/<\\\//g, '</'),
    )
    expect(inner).toContain('<base href="about:srcdoc">')
    expect(inner.indexOf(INERT_NAVIGATION)).toBeGreaterThan(-1)
    expect(inner.indexOf(INERT_NAVIGATION)).toBeLessThan(inner.indexOf('BRIDGE'))
    expect(inner.indexOf('BRIDGE')).toBeLessThan(inner.indexOf('app()'))
  })

  // A page that throws draws its headings and stops, which reads as a working
  // page with no data in it. The reporter has to be first, or it misses the
  // very errors that stop the page before anything else runs.
  it('carries a failure reporter ahead of everything the page runs', () => {
    const doc = buildOuterDocument({
      source: '<script>boom()</script>',
      hostScript: '',
      bridgeScript: 'BRIDGE',
    })
    const inner = JSON.parse(
      doc.match(/window.__humanInner = (".*?")<\/script>/s)[1].replace(/<\\\//g, '</'),
    )
    expect(inner).toContain(REPORT_FAILURE)
    expect(inner.indexOf(REPORT_FAILURE)).toBeLessThan(inner.indexOf(INERT_NAVIGATION))
    expect(inner.indexOf(REPORT_FAILURE)).toBeLessThan(inner.indexOf('boom()'))
  })

  it('takes peer connections away before the bridge or any app script runs', () => {
    const doc = buildOuterDocument({
      source: '<script>app()</script>',
      hostScript: '',
      bridgeScript: 'BRIDGE',
    })
    const inner = JSON.parse(
      doc.match(/window.__humanInner = (".*?")<\/script>/s)[1].replace(/<\\\//g, '</'),
    )
    expect(inner.indexOf(REPORT_FAILURE)).toBeLessThan(inner.indexOf(NO_PEER_CONNECTIONS))
    expect(inner.indexOf(NO_PEER_CONNECTIONS)).toBeLessThan(inner.indexOf('BRIDGE'))
    expect(inner.indexOf(NO_PEER_CONNECTIONS)).toBeLessThan(inner.indexOf('app()'))
  })

  it('leaves the window no peer connection constructor', () => {
    const sandbox = {
      RTCPeerConnection() {},
      webkitRTCPeerConnection() {},
      RTCSessionDescription() {},
    }
    sandbox.window = sandbox
    runInContext(NO_PEER_CONNECTIONS, createContext(sandbox))
    expect('RTCPeerConnection' in sandbox).toBe(false)
    expect('webkitRTCPeerConnection' in sandbox).toBe(false)
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

describe('changing records', () => {
  const meta = { namespace: 'agent-sandbox', modules: ['agent-contact'], writes: ['agent-contact'] }
  const fields = [
    { name: 'name', kind: 'String' },
    { name: 'active', kind: 'Bool' },
    { name: 'amount', kind: 'Number' },
    { name: 'tags', kind: 'String', multi: true },
  ]

  const writeContext = (over = {}) => ({
    ...context(meta),
    version: 2,
    fields: () => fields,
    consent: vi.fn().mockResolvedValue(true),
    ...over,
  })

  it('refuses a module the app reads but was not given to change', async () => {
    expect(allowWrite({ modules: ['a'], writes: [] }, 'a')).toBe(
      'module "a" is not declared as one this app may change',
    )
    expect(allowWrite({ modules: ['a'], writes: ['a'] }, 'a')).toBe(null)
    expect(allowWrite({ modules: [], writes: ['a'] }, 'a')).toBe(
      'module "a" is not declared for this app',
    )
  })

  it('writes values the way the webapp does, and refuses a field the module has not', () => {
    expect(
      toStoreValues({ name: 'Ana', active: true, amount: 12.5, tags: ['a', 'b'] }, fields),
    ).toEqual([
      { name: 'name', value: 'Ana' },
      { name: 'active', value: '1' },
      { name: 'amount', value: '12.5' },
      { name: 'tags', value: 'a' },
      { name: 'tags', value: 'b' },
    ])
    expect(toStoreValues({ active: false }, fields)).toEqual([{ name: 'active', value: '' }])
    expect(() => toStoreValues({ nope: 'x' }, fields)).toThrow(
      '"nope" is not a field of this module',
    )
  })

  it('creates a record once the viewer has agreed', async () => {
    const ctx = writeContext()
    ctx.compose.recordCreate = vi
      .fn()
      .mockResolvedValue({ recordID: '9', values: [{ name: 'name', value: 'Ana' }] })

    const out = await dispatch(
      'records.create',
      { module: 'agent-contact', values: { name: 'Ana' } },
      ctx,
    )

    expect(ctx.consent).toHaveBeenCalledWith('agent-contact')
    expect(ctx.compose.recordCreate.mock.calls[0][0]).toMatchObject({
      moduleID: '2',
      values: [{ name: 'name', value: 'Ana' }],
    })
    expect(out.record.values).toEqual({ name: 'Ana', active: false })
  })

  it('changes nothing when the viewer says no', async () => {
    const ctx = writeContext({ consent: vi.fn().mockResolvedValue(false) })
    ctx.compose.recordCreate = vi.fn()

    await expect(
      dispatch('records.create', { module: 'agent-contact', values: { name: 'Ana' } }, ctx),
    ).rejects.toThrow('did not agree')
    expect(ctx.compose.recordCreate).not.toHaveBeenCalled()
  })

  // The endpoint behind an update changes every record a filter matches, so an
  // update with no record named must never reach it.
  it('refuses an update that names no record, without calling the API', async () => {
    const ctx = writeContext()
    ctx.compose.recordPatch = vi.fn()

    await expect(
      dispatch('records.update', { module: 'agent-contact', values: { name: 'Ana' } }, ctx),
    ).rejects.toThrow('needs the recordID')
    expect(ctx.compose.recordPatch).not.toHaveBeenCalled()
  })

  it('changes only the named record and reads it back', async () => {
    const ctx = writeContext()
    ctx.compose.recordPatch = vi.fn().mockResolvedValue({})
    ctx.compose.recordRead = vi
      .fn()
      .mockResolvedValue({ recordID: '9', values: [{ name: 'name', value: 'Bo' }] })

    const out = await dispatch(
      'records.update',
      { module: 'agent-contact', recordID: '9', values: { name: 'Bo' } },
      ctx,
    )

    expect(ctx.compose.recordPatch.mock.calls[0][0]).toMatchObject({
      moduleID: '2',
      recordID: ['9'],
    })
    expect(out.record.values.name).toBe('Bo')
  })
})

describe('handing the viewer a file', () => {
  it('keeps the app to one file name', () => {
    expect(downloadName('leads.csv')).toBe('leads.csv')
    expect(downloadName('../../etc/passwd')).toBe('-..-etc-passwd')
    expect(downloadName('  ')).toBe('download.txt')
    expect(downloadName('.bashrc')).toBe('bashrc')
  })

  it('passes the file to the shell, which is the only one that can save it', async () => {
    const ctx = { ...context(), download: vi.fn() }
    await expect(dispatch('download', { name: 'a/b.csv', text: 'x,y' }, ctx)).resolves.toBe(true)
    expect(ctx.download).toHaveBeenCalledWith('a-b.csv', 'x,y')
  })

  it('refuses one too large to be worth handing over', async () => {
    const ctx = { ...context(), download: vi.fn() }
    await expect(
      dispatch('download', { name: 'big.csv', text: 'x'.repeat(MAX_DOWNLOAD + 1) }, ctx),
    ).rejects.toThrow('at most')
    expect(ctx.download).not.toHaveBeenCalled()
  })
})

describe('a breakdown grouped by a reference', () => {
  const company = {
    name: 'company',
    kind: 'Record',
    options: { moduleID: '9', labelField: 'name' },
  }
  const owner = { name: 'owner', kind: 'User' }
  const rows = [
    { count: 2, dimension_0: '100' },
    { count: 1, dimension_0: null },
  ]

  const reportContext = fields => ({
    ...context(),
    version: 2,
    fields: () => fields,
    userLabels: vi.fn().mockResolvedValue({ 100: 'Dev Agent' }),
    compose: {
      recordReport: vi.fn().mockResolvedValue(rows),
      moduleRead: vi.fn().mockResolvedValue({ fields: [{ name: 'name', kind: 'String' }] }),
      recordList: vi.fn().mockResolvedValue({
        set: [{ recordID: '100', values: [{ name: 'name', value: 'Acme' }] }],
      }),
    },
  })

  it('names the records a report groups by', async () => {
    expect(await dimensionRefs(reportContext([company]), '2', 'company', rows)).toEqual({
      100: 'Acme',
    })
  })

  it('names the people a report groups by', async () => {
    expect(await dimensionRefs(reportContext([owner]), '2', 'owner', rows)).toEqual({
      100: 'Dev Agent',
    })
  })

  it('leaves a breakdown by something else alone', async () => {
    const ctx = reportContext([{ name: 'status', kind: 'Select' }])
    expect(await dimensionRefs(ctx, '2', 'status', rows)).toEqual({})
    expect(ctx.compose.recordList).not.toHaveBeenCalled()
  })

  it('hands contract 2 the rows with their labels, and contract 1 the rows alone', async () => {
    const ctx = reportContext([company])
    const out = await dispatch(
      'records.report',
      { module: 'agent-contact', dimensions: 'company' },
      ctx,
    )
    expect(out).toEqual({ rows, refs: { 100: 'Acme' } })

    const old = await dispatch(
      'records.report',
      { module: 'agent-contact', dimensions: 'company' },
      { ...ctx, version: 1 },
    )
    expect(old).toEqual(rows)
  })
})

describe('what the app may read about its modules', () => {
  const fields = [
    { name: 'name', label: 'Full name', kind: 'String' },
    {
      name: 'status',
      kind: 'Select',
      options: { options: [{ value: 'lead', text: 'Lead' }, { value: 'active' }] },
    },
    { name: 'tags', kind: 'String', multi: true },
  ]

  it('describes a field by what it is called and what it may hold', () => {
    expect(describeField(fields[0])).toEqual({
      name: 'name',
      label: 'Full name',
      kind: 'String',
      multi: false,
    })
    // An option with no text of its own is known by its value, and carries its
    // wording under both names — `text` is what the module calls it, `label` is
    // what every other field here is called.
    expect(describeField(fields[1]).options).toEqual([
      { value: 'lead', text: 'Lead', label: 'Lead' },
      { value: 'active', text: 'active', label: 'active' },
    ])
    expect(describeField(fields[2]).multi).toBe(true)
  })

  it('answers with the declared modules, saying which may be changed', async () => {
    const ctx = {
      ...context({
        namespace: 'agent-sandbox',
        modules: ['agent-contact'],
        writes: ['agent-contact'],
      }),
      fields: () => fields,
    }

    expect(await dispatch('modules', {}, ctx)).toEqual([
      {
        handle: 'agent-contact',
        moduleID: '2',
        writable: true,
        fields: fields.map(describeField),
      },
    ])
  })

  it('says nothing about a module the app never declared', async () => {
    const ctx = { ...context({ namespace: 'x', modules: [] }), fields: () => fields }
    expect(await dispatch('modules', {}, ctx)).toEqual([])
  })
})

describe('sandbox policy', () => {
  it('admits a listed origin for script, style, font and image, never for a network call', () => {
    const csp = cspInner(['https://cdn.jsdelivr.net', 'https://fonts.gstatic.com'])
    expect(csp).toContain(
      "script-src 'unsafe-inline' https://cdnjs.cloudflare.com https://cdn.jsdelivr.net https://fonts.gstatic.com;",
    )
    expect(csp).toContain("style-src 'unsafe-inline' https://cdn.jsdelivr.net")
    expect(csp).toContain('font-src data: https://cdn.jsdelivr.net')
    expect(csp).toContain('img-src data: https://cdn.jsdelivr.net')
    expect(csp).toContain("connect-src 'none'")
  })

  it('drops anything that is not a bare https origin', () => {
    const csp = cspInner([
      'http://plain.test',
      'https://x.test; connect-src *',
      'https://*.test',
      'https://ok.test"',
    ])
    expect(csp).toBe(cspInner())
  })
})

describe('navigation', () => {
  it('takes a page, or a module with the record to open', () => {
    expect(navigationTarget({ page: 'board' })).toEqual({ page: 'board', module: '', recordID: '' })
    expect(navigationTarget({ module: 'task', recordID: 42 })).toEqual({
      page: '',
      module: 'task',
      recordID: '42',
    })
  })

  it('refuses what names no single place', () => {
    expect(() => navigationTarget({})).toThrow('needs a page')
    expect(() => navigationTarget({ page: 'a', module: 'b' })).toThrow('not both')
    expect(() => navigationTarget({ module: 'task' })).toThrow('name it with recordID')
    expect(() => navigationTarget({ page: 'a', recordID: '1 or 1' })).toThrow('not a record ID')
  })

  it('hands the target to the shell', async () => {
    const ctx = { ...context(), navigate: vi.fn() }
    await dispatch('navigate', { page: 'board' }, ctx)
    expect(ctx.navigate).toHaveBeenCalledWith({ page: 'board', module: '', recordID: '' })
  })
})

describe('files', () => {
  const fileCtx = () => ({
    ...context(),
    fields: () => [
      { name: 'photo', kind: 'File' },
      { name: 'name', kind: 'String' },
    ],
    compose: {
      recordRead: vi.fn().mockResolvedValue({
        values: [
          { name: 'photo', value: '11' },
          { name: 'photo', value: '12' },
        ],
      }),
    },
    attachment: vi.fn(id =>
      Promise.resolve({
        attachmentID: id,
        name: `f${id}.png`,
        meta: { original: { size: 10, mimetype: 'image/png' } },
      }),
    ),
    fileData: vi.fn().mockResolvedValue('data:image/png;base64,AA=='),
  })

  it('lists and reads the files a record field holds', async () => {
    const ctx = fileCtx()
    const args = { module: 'agent-contact', recordID: '5', field: 'photo' }
    expect(await dispatch('files.list', args, ctx)).toEqual([
      { attachmentID: '11', name: 'f11.png', mimetype: 'image/png', size: 10 },
      { attachmentID: '12', name: 'f12.png', mimetype: 'image/png', size: 10 },
    ])
    expect(await dispatch('files.read', { ...args, attachmentID: '12' }, ctx)).toMatchObject({
      attachmentID: '12',
      dataURL: 'data:image/png;base64,AA==',
    })
  })

  it('refuses a file the record does not hold, a field that holds no files, an undeclared module', async () => {
    const ctx = fileCtx()
    const args = { module: 'agent-contact', recordID: '5', field: 'photo' }
    await expect(dispatch('files.read', { ...args, attachmentID: '99' }, ctx)).rejects.toThrow(
      'is not one field',
    )
    await expect(dispatch('files.read', { ...args, field: 'name' }, ctx)).rejects.toThrow(
      'not a File field',
    )
    await expect(dispatch('files.list', { ...args, module: 'Lead' }, ctx)).rejects.toThrow(
      'not declared',
    )
    expect(ctx.fileData).not.toHaveBeenCalled()
  })

  it('refuses a file over the cap before reading it', async () => {
    const ctx = fileCtx()
    ctx.attachment = vi.fn().mockResolvedValue({
      attachmentID: '11',
      name: 'big',
      meta: { original: { size: 6 * 1024 * 1024 } },
    })
    await expect(
      dispatch('files.read', { module: 'agent-contact', recordID: '5', field: 'photo' }, ctx),
    ).rejects.toThrow('at most')
    expect(ctx.fileData).not.toHaveBeenCalled()
  })
})

describe('automations', () => {
  const runCtx = (agree = true) => ({
    ...context({ ...META, automations: ['close_deal'] }),
    consentToRun: vi.fn().mockResolvedValue(agree),
    runAutomation: vi.fn().mockResolvedValue({ ok: true }),
  })

  it('types what the app hands an automation', () => {
    expect(automationInput({ a: 'x', n: 2, f: 1.5, b: true })).toEqual({
      a: { '@type': 'String', '@value': 'x' },
      n: { '@type': 'Integer', '@value': 2 },
      f: { '@type': 'Float', '@value': 1.5 },
      b: { '@type': 'Boolean', '@value': true },
    })
    expect(() => automationInput({ o: { deep: 1 } })).toThrow('not a string, number or boolean')
  })

  it('runs a declared automation once the viewer agrees', async () => {
    const ctx = runCtx()
    expect(
      await dispatch('automation.run', { automation: 'close_deal', input: { n: 1 } }, ctx),
    ).toEqual({ ok: true })
    expect(ctx.runAutomation).toHaveBeenCalledWith('close_deal', {
      n: { '@type': 'Integer', '@value': 1 },
    })
  })

  it('refuses an undeclared automation before asking, and runs nothing the viewer refused', async () => {
    let ctx = runCtx()
    await expect(dispatch('automation.run', { automation: 'drop_all' }, ctx)).rejects.toThrow(
      'not declared',
    )
    expect(ctx.consentToRun).not.toHaveBeenCalled()

    ctx = runCtx(false)
    await expect(dispatch('automation.run', { automation: 'close_deal' }, ctx)).rejects.toThrow(
      'did not allow',
    )
    expect(ctx.runAutomation).not.toHaveBeenCalled()
  })
})

describe('page refresh', () => {
  it('tells the host after a change succeeds, and when the app asks', async () => {
    const ctx = {
      ...context({ ...META, automations: ['close_deal'] }),
      consentToRun: vi.fn().mockResolvedValue(true),
      runAutomation: vi.fn().mockResolvedValue({ ok: true }),
      changed: vi.fn(),
    }
    await dispatch('automation.run', { automation: 'close_deal' }, ctx)
    expect(ctx.changed).toHaveBeenCalledTimes(1)
    await dispatch('refresh', {}, ctx)
    expect(ctx.changed).toHaveBeenCalledTimes(2)
  })

  it('stays quiet for a read and for a change that failed', async () => {
    const ctx = { ...context(), changed: vi.fn(), consentToRun: vi.fn().mockResolvedValue(false) }
    await dispatch('records.list', { module: 'agent-contact' }, ctx)
    await expect(dispatch('automation.run', { automation: 'nope' }, ctx)).rejects.toThrow()
    expect(ctx.changed).not.toHaveBeenCalled()
  })
})

describe('chatbots', () => {
  const chatCtx = () => ({
    ...context({ ...META, chatbots: ['support'] }),
    openChatbot: vi.fn().mockResolvedValue(undefined),
    closeChatbot: vi.fn(),
  })

  it('opens and closes a declared chatbot through the shell', async () => {
    const ctx = chatCtx()
    expect(await dispatch('chatbot.open', { chatbot: 'support' }, ctx)).toBe(true)
    expect(ctx.openChatbot).toHaveBeenCalledWith('support')
    await dispatch('chatbot.close', { chatbot: 'support' }, ctx)
    expect(ctx.closeChatbot).toHaveBeenCalledWith('support')
  })

  it('refuses a chatbot the app did not declare', async () => {
    const ctx = chatCtx()
    await expect(dispatch('chatbot.open', { chatbot: 'sales' }, ctx)).rejects.toThrow(
      'not declared',
    )
    await expect(dispatch('chatbot.open', {}, ctx)).rejects.toThrow('needs the handle')
    expect(ctx.openChatbot).not.toHaveBeenCalled()
  })
})

describe('shell calls', () => {
  it('shows a toast in the severity asked, and info for anything else', async () => {
    const ctx = { ...context(), toast: vi.fn() }
    await dispatch('toast', { message: 'saved', severity: 'success' }, ctx)
    await dispatch('toast', { message: 'hm', severity: 'shout' }, ctx)
    expect(ctx.toast.mock.calls).toEqual([
      ['success', 'saved', ''],
      ['info', 'hm', ''],
    ])
    await expect(dispatch('toast', {}, ctx)).rejects.toThrow('needs a message')
  })

  it('asks the viewer through Human and hands back the answer', async () => {
    const ctx = {
      ...context(),
      ask: vi.fn().mockResolvedValueOnce(true).mockResolvedValueOnce('Ana'),
    }
    expect(await dispatch('confirm', { message: 'Sure?' }, ctx)).toBe(true)
    expect(await dispatch('prompt', { message: 'Name?', value: 'A' }, ctx)).toBe('Ana')
    expect(ctx.ask.mock.calls[1]).toEqual(['prompt', { message: 'Name?', title: '', value: 'A' }])
  })

  it('bounds what an app puts in Human chrome', async () => {
    const ctx = { ...context(), setTitle: vi.fn() }
    await dispatch('title', { text: 'x'.repeat(900) }, ctx)
    expect(ctx.setTitle.mock.calls[0][0]).toHaveLength(500)
  })

  it('opens a record of a declared module only', async () => {
    const ctx = { ...context(), openRecord: vi.fn() }
    await dispatch('records.open', { module: 'agent-contact', recordID: '5' }, ctx)
    expect(ctx.openRecord).toHaveBeenCalledWith({
      module: 'agent-contact',
      recordID: '5',
      edit: false,
    })
    await expect(dispatch('records.open', { module: 'Lead', recordID: '5' }, ctx)).rejects.toThrow(
      'not declared',
    )
    await expect(dispatch('records.open', { module: 'agent-contact' }, ctx)).rejects.toThrow(
      'needs the recordID',
    )
  })
})

describe('deleting', () => {
  const delCtx = (agree = true) => ({
    ...context({ ...META, deletes: ['agent-contact'] }),
    compose: { recordDelete: vi.fn().mockResolvedValue(undefined) },
    consentToDelete: vi.fn().mockResolvedValue(agree),
    changed: vi.fn(),
  })

  it('deletes one record of a module declared for it, once the viewer agrees', async () => {
    const ctx = delCtx()
    expect(await dispatch('records.delete', { module: 'agent-contact', recordID: '9' }, ctx)).toBe(
      true,
    )
    expect(ctx.compose.recordDelete).toHaveBeenCalledWith({
      namespaceID: '1',
      moduleID: '2',
      recordID: '9',
    })
    expect(ctx.changed).toHaveBeenCalled()
  })

  it('refuses a module not declared for deleting, and does nothing the viewer refused', async () => {
    let ctx = { ...delCtx(), meta: META }
    await expect(
      dispatch('records.delete', { module: 'agent-contact', recordID: '9' }, ctx),
    ).rejects.toThrow('not one this app may delete records from')
    ctx = delCtx(false)
    await expect(
      dispatch('records.delete', { module: 'agent-contact', recordID: '9' }, ctx),
    ).rejects.toThrow('did not agree to it deleting')
    expect(ctx.compose.recordDelete).not.toHaveBeenCalled()
  })
})

describe('uploading', () => {
  const upCtx = multi => ({
    ...context({ ...META, writes: ['agent-contact'] }),
    fields: () => [{ name: 'photo', kind: 'File', isMulti: multi }],
    compose: {
      recordRead: vi.fn().mockResolvedValue({ values: [{ name: 'photo', value: '11' }] }),
      recordPatch: vi.fn().mockResolvedValue(undefined),
    },
    consent: vi.fn().mockResolvedValue(true),
    uploadFile: vi.fn().mockResolvedValue({
      attachmentID: '12',
      name: 'a.txt',
      meta: { original: { size: 2, mimetype: 'text/plain' } },
    }),
  })
  const args = {
    module: 'agent-contact',
    recordID: '5',
    field: 'photo',
    name: 'a.txt',
    dataURL: 'data:text/plain;base64,aGk=',
  }

  it('stores the file and sets the field to it', async () => {
    const ctx = upCtx(false)
    expect(await dispatch('files.upload', args, ctx)).toEqual({
      attachmentID: '12',
      name: 'a.txt',
      mimetype: 'text/plain',
      size: 2,
    })
    const sent = ctx.uploadFile.mock.calls[0][0]
    expect(sent.file.name).toBe('a.txt')
    expect(new TextDecoder().decode(sent.file.data)).toBe('hi')
    expect(ctx.compose.recordPatch.mock.calls[0][0].values).toEqual([
      { name: 'photo', value: '12' },
    ])
  })

  it('adds to what a multi-value field holds', async () => {
    const ctx = upCtx(true)
    await dispatch('files.upload', args, ctx)
    expect(ctx.compose.recordPatch.mock.calls[0][0].values).toEqual([
      { name: 'photo', value: '11' },
      { name: 'photo', value: '12' },
    ])
  })

  it('refuses a module it may not change, and what is not a data: URL', async () => {
    await expect(dispatch('files.upload', args, { ...upCtx(false), meta: META })).rejects.toThrow()
    await expect(
      dispatch('files.upload', { ...args, dataURL: 'https://x' }, upCtx(false)),
    ).rejects.toThrow('data: URL')
  })
})

describe('users', () => {
  it('searches as the viewer, bounded', async () => {
    const ctx = { ...context(), searchUsers: vi.fn().mockResolvedValue([]) }
    await dispatch('users.search', { query: 'an', limit: 500 }, ctx)
    expect(ctx.searchUsers).toHaveBeenCalledWith('an', 50)
    await expect(dispatch('users.search', { query: 'a' }, ctx)).rejects.toThrow('two characters')
  })
})

describe('fields', () => {
  it('reads multi-value from the name the API uses', () => {
    expect(describeField({ name: 'tags', kind: 'String', isMulti: true }).multi).toBe(true)
  })
})

// A record ID is written into the request path, so `../` in one would reach
// another module, the namespace, or any endpoint the viewer may call.
describe('record IDs', () => {
  const anyCtx = () => ({
    ...context({
      ...META,
      writes: ['agent-contact'],
      deletes: ['agent-contact'],
    }),
    fields: () => [
      { name: 'photo', kind: 'File' },
      { name: 'name', kind: 'String' },
    ],
    compose: {
      recordRead: vi.fn().mockResolvedValue({ values: [] }),
      recordPatch: vi.fn().mockResolvedValue(undefined),
      recordDelete: vi.fn().mockResolvedValue(undefined),
    },
    consent: vi.fn().mockResolvedValue(true),
    consentToDelete: vi.fn().mockResolvedValue(true),
    uploadFile: vi.fn(),
    attachment: vi.fn(),
    changed: vi.fn(),
  })
  const calls = {
    'records.read': {},
    'records.update': { values: { name: 'x' } },
    'records.delete': {},
    'files.list': { field: 'photo' },
    'files.read': { field: 'photo' },
    'files.upload': { field: 'photo', name: 'a.txt', dataURL: 'data:text/plain;base64,aGk=' },
  }

  for (const [op, args] of Object.entries(calls)) {
    it(`refuses ${op} on a record ID that is not one, before the API is touched`, async () => {
      for (const recordID of ['../../3/record/4', '%2e%2e/3', '9/../10', '9?x=1', '9#']) {
        const ctx = anyCtx()
        await expect(
          dispatch(op, { module: 'agent-contact', recordID, ...args }, ctx),
        ).rejects.toThrow('is not a record ID')
        expect(ctx.compose.recordRead).not.toHaveBeenCalled()
        expect(ctx.compose.recordPatch).not.toHaveBeenCalled()
        expect(ctx.compose.recordDelete).not.toHaveBeenCalled()
        expect(ctx.uploadFile).not.toHaveBeenCalled()
      }
    })
  }

  it('reads the one record a numeric ID names', async () => {
    const ctx = anyCtx()
    await dispatch('records.read', { module: 'agent-contact', recordID: 9 }, ctx)
    expect(ctx.compose.recordRead).toHaveBeenCalledWith({
      namespaceID: '1',
      moduleID: '2',
      recordID: '9',
    })
  })
})
