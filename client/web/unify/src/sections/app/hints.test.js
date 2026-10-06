import { createContext, runInContext } from 'node:vm'
import { readFileSync } from 'node:fs'
import { fileURLToPath } from 'node:url'
import { dirname, join } from 'node:path'
import { describe, expect, it } from 'vitest'
import yaml from 'js-yaml'
import { EditorState } from '@codemirror/state'
import { EditorView } from '@codemirror/view'
import { html } from '@codemirror/lang-html'
import { CompletionContext } from '@codemirror/autocomplete'
import { BRIDGE_SCRIPT } from './bridge'
import {
  BRIDGE_CALLS,
  CONTEXT_KEYS,
  bridgeAssist,
  inScript,
  noteKey,
  snippets,
  suggestAt,
} from './hints'

const LOCALE = join(
  dirname(fileURLToPath(import.meta.url)),
  '../../../../../../locale/en/human-webapp/app.yaml',
)
const app = yaml.load(readFileSync(LOCALE, 'utf8'))
const t = key =>
  key
    .split('.')
    .slice(1)
    .reduce((node, k) => node?.[k], app) ?? key

// Every member of the object BRIDGE_SCRIPT leaves on window.human, nested
// ones as `group.member`.
function bridgeMembers() {
  const sandbox = {
    setTimeout: () => 0,
    clearTimeout: () => {},
    addEventListener: () => {},
    parent: { postMessage: () => {} },
  }
  sandbox.window = sandbox
  createContext(sandbox)
  runInContext(BRIDGE_SCRIPT, sandbox)

  const names = []
  for (const [key, value] of Object.entries(sandbox.human)) {
    names.push(key)
    if (value && typeof value === 'object' && typeof value.then !== 'function') {
      for (const inner of Object.keys(value)) names.push(`${key}.${inner}`)
    }
  }
  return names.sort()
}

const labels = found => (found?.options || []).map(o => o.label || o.call.name.split('.').pop())

const at = (text, declared) => suggestAt(text.replace('|', ''), text.indexOf('|'), declared)

describe('bridge hints catalog', () => {
  it('has an entry for every member of the bridge, and none it lacks', () => {
    expect(BRIDGE_CALLS.map(c => c.name).sort()).toEqual(bridgeMembers())
  })

  it('lists a group with exactly the calls under it', () => {
    for (const group of BRIDGE_CALLS.filter(c => c.group)) {
      const under = BRIDGE_CALLS.filter(c => c.name.startsWith(`${group.name}.`))
      expect([...group.group].sort()).toEqual(under.map(c => c.name.split('.')[1]).sort())
    }
  })

  it('has a note in the locale for every call, context key and kind', () => {
    for (const call of BRIDGE_CALLS)
      expect(t(noteKey(call.name)), call.name).not.toBe(noteKey(call.name))
    for (const key of CONTEXT_KEYS) expect(app.bridge.contextKey[key], key).toBeTruthy()
    for (const kind of ['module', 'automation', 'chatbot', 'context', 'param']) {
      expect(app.bridge.kind[kind], kind).toBeTruthy()
    }
    for (const { key } of snippets()) expect(app.snippets[key], key).toBeTruthy()
  })

  it('writes the declared module into the examples', () => {
    expect(snippets('task').find(s => s.key === 'readRecord').text).toContain("module: 'task'")
  })
})

describe('suggestAt', () => {
  it('offers the bridge after human.', () => {
    const got = labels(at('human.|'))
    expect(got).toContain('records')
    expect(got).toContain('context')
    expect(got).not.toContain('list')
  })

  it('offers a group’s calls after human.<group>.', () => {
    expect(labels(at('await human.records.re|')).sort()).toEqual([
      'create',
      'delete',
      'list',
      'open',
      'read',
      'report',
      'update',
    ])
    expect(at('await human.records.re|').from).toBe('await human.records.'.length)
  })

  it('offers the declared names inside module, automation and chatbot strings', () => {
    const declared = { modules: ['task', 'project'], automations: ['notify'], chatbots: ['help'] }
    expect(labels(at("human.records.list({ module: 'ta|", declared))).toEqual(['task', 'project'])
    expect(at("human.records.list({ module: 'ta|", declared).from).toBe(
      "human.records.list({ module: '".length,
    )
    expect(labels(at('human.automation.run({ automation: "|', declared))).toEqual(['notify'])
    expect(labels(at('human.chatbot.open({ chatbot: `|', declared))).toEqual(['help'])
  })

  it('offers what human.context() answers on a name that holds it', () => {
    const page = 'const ctx = await human.context()\nctx.|'
    expect(labels(at(page))).toEqual(CONTEXT_KEYS)
    expect(labels(at('const other = 1\nother.|'))).toEqual([])
  })

  it('offers the block’s parameters under params', () => {
    const declared = { params: ['status', 'limit'] }
    expect(labels(at('const ctx = await human.context()\nctx.params.|', declared))).toEqual([
      'status',
      'limit',
    ])
    expect(labels(at('const { params } = await human.context()\nparams.|', declared))).toEqual([
      'status',
      'limit',
    ])
    expect(labels(at('const { params: p } = await human.context()\np.|', declared))).toEqual([
      'status',
      'limit',
    ])
    expect(labels(at('(await human.context()).params.|', declared))).toEqual(['status', 'limit'])
  })

  it('offers the keys not yet taken inside a destructuring of human.context()', () => {
    expect(labels(at('const { recordID, | } = await human.context()'))).toEqual(
      CONTEXT_KEYS.filter(k => k !== 'recordID'),
    )
  })
})

describe('bridgeAssist in an editor', () => {
  const stateFor = (text, declared = {}) =>
    EditorState.create({
      doc: text.replace('|', ''),
      selection: { anchor: text.indexOf('|') },
      extensions: [html(), bridgeAssist(t, () => declared)],
    })

  const complete = async (text, declared) => {
    const state = stateFor(text, declared)
    const pos = state.selection.main.head
    const results = await Promise.all(
      state
        .languageDataAt('autocomplete', pos)
        .map(s => s(new CompletionContext(state, pos, true))),
    )
    return results.flatMap(r => (r ? r.options : []))
  }

  it('suggests the bridge inside a script', async () => {
    const options = await complete('<script>\nhuman.rec|\n</script>')
    const records = options.find(o => o.label === 'records')
    expect(records).toBeTruthy()
    expect(records.info).toContain('human.records.{ list, read')
    expect(records.info).toContain(app.bridge.records)
  })

  it('suggests the bridge in an on… attribute', async () => {
    const options = await complete('<button onclick="human.to|">x</button>')
    expect(options.map(o => o.label)).toContain('toast')
  })

  it('suggests nothing of the bridge in the page’s text', async () => {
    const options = await complete('<p>human.|</p>')
    expect(options.map(o => o.label)).not.toContain('records')
  })

  it('reads what was declared afresh on every suggestion', async () => {
    const declared = { modules: ['task'] }
    const state = stateFor("<script>human.records.list({ module: '|", declared)
    const pos = state.selection.main.head
    const run = () =>
      Promise.all(
        state
          .languageDataAt('autocomplete', pos)
          .map(s => s(new CompletionContext(state, pos, true))),
      ).then(rs => rs.flatMap(r => (r ? r.options.map(o => o.label) : [])))

    expect(await run()).toContain('task')
    declared.modules = ['project']
    expect(await run()).toContain('project')
  })

  it('writes the call with its arguments when a suggestion is taken', async () => {
    const state = stateFor('<script>\nhuman.records.li|\n</script>')
    const view = new EditorView({ state, parent: document.body })
    const pos = state.selection.main.head
    const [result] = (
      await Promise.all(
        state
          .languageDataAt('autocomplete', pos)
          .map(s => s(new CompletionContext(state, pos, true))),
      )
    ).filter(Boolean)
    const list = result.options.find(o => o.label === 'list')

    list.apply(view, list, result.from, pos)
    expect(view.state.doc.toString()).toContain("human.records.list({ module: 'module' })")
    view.destroy()
  })

  it('knows where JavaScript is', () => {
    const state = stateFor('<p>a</p><script>b|</script>')
    expect(inScript(state, state.selection.main.head)).toBe(true)
    expect(inScript(state, 4)).toBe(false)
  })
})
