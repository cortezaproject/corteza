// Suggestions and hover notes for the `window.human` bridge, for the editors
// a custom app's or a Custom block's page is written in. Every call of
// BRIDGE_SCRIPT has an entry here; each one's note is `app.bridge.<op>` in the
// locale, with dots as underscores.
import { snippetCompletion, startCompletion } from '@codemirror/autocomplete'
import { syntaxTree } from '@codemirror/language'
import { EditorState } from '@codemirror/state'
import { EditorView, hoverTooltip } from '@codemirror/view'

// The bridge as it reads after `human.`: a group holds calls, a call has the
// arguments it takes and what it answers, and `snippet` is what accepting it
// writes, with `${}` fields Tab moves through.
export const BRIDGE_CALLS = [
  {
    name: 'records',
    group: ['list', 'read', 'report', 'create', 'update', 'delete', 'open'],
  },
  {
    name: 'records.list',
    args: '{ module, filter?, sort?, limit?, pageCursor? }',
    result: '{ records, refs, nextPageCursor }',
    snippet: "list({ module: '${module}' })",
  },
  {
    name: 'records.read',
    args: '{ module, recordID }',
    result: '{ record, refs }',
    snippet: "read({ module: '${module}', recordID: ${recordID} })",
  },
  {
    name: 'records.report',
    args: '{ module, dimension, metrics?, filter? }',
    result: '{ rows, refs }',
    snippet: "report({ module: '${module}', dimension: '${field}' })",
  },
  {
    name: 'records.create',
    args: '{ module, values }',
    result: '{ record }',
    snippet: "create({ module: '${module}', values: { ${} } })",
  },
  {
    name: 'records.update',
    args: '{ module, recordID, values }',
    result: '{ record }',
    snippet: "update({ module: '${module}', recordID: ${recordID}, values: { ${} } })",
  },
  {
    name: 'records.delete',
    args: '{ module, recordID }',
    result: 'true',
    snippet: "delete({ module: '${module}', recordID: ${recordID} })",
  },
  {
    name: 'records.open',
    args: '{ module, recordID, edit? }',
    result: 'true',
    snippet: "open({ module: '${module}', recordID: ${recordID} })",
  },
  { name: 'users', group: ['search'] },
  {
    name: 'users.search',
    args: '{ query, limit? }',
    result: '[{ userID, name, email }]',
    snippet: "search({ query: '${}' })",
  },
  { name: 'files', group: ['list', 'read', 'upload'] },
  {
    name: 'files.list',
    args: '{ module, recordID, field }',
    result: '[{ attachmentID, name, mimetype, size }]',
    snippet: "list({ module: '${module}', recordID: ${recordID}, field: '${field}' })",
  },
  {
    name: 'files.read',
    args: '{ module, recordID, field, attachmentID? }',
    result: '{ dataURL, … }',
    snippet: "read({ module: '${module}', recordID: ${recordID}, field: '${field}' })",
  },
  {
    name: 'files.upload',
    args: '{ module, recordID, field, name, dataURL }',
    result: 'the stored file',
    snippet:
      "upload({ module: '${module}', recordID: ${recordID}, field: '${field}', name: ${name}, dataURL: ${dataURL} })",
  },
  { name: 'chatbot', group: ['open', 'close'] },
  {
    name: 'chatbot.open',
    args: '{ chatbot }',
    result: 'true',
    snippet: "open({ chatbot: '${chatbot}' })",
  },
  {
    name: 'chatbot.close',
    args: '{ chatbot }',
    result: 'true',
    snippet: "close({ chatbot: '${chatbot}' })",
  },
  { name: 'automation', group: ['run'] },
  {
    name: 'automation.run',
    args: '{ automation, input? }',
    result: '{ ok: true }',
    snippet: "run({ automation: '${automation}', input: { ${} } })",
  },
  {
    name: 'context',
    args: '',
    result: '{ namespaceID, namespace, pageID, moduleID, recordID, params }',
    snippet: 'context()',
  },
  {
    name: 'modules',
    args: '',
    result: '[{ handle, moduleID, writable, fields }]',
    snippet: 'modules()',
  },
  { name: 'user', args: '', result: '{ userID, name, email }', snippet: 'user()' },
  { name: 'theme', args: '', result: '{ dark, colors }', snippet: 'theme()' },
  {
    name: 'toast',
    args: '{ message, severity?, title? }',
    result: 'true',
    snippet: "toast({ message: '${}', severity: '${success}' })",
  },
  {
    name: 'confirm',
    args: '{ message, title?, accept?, reject? }',
    result: 'true or false',
    snippet: "confirm({ message: '${}' })",
  },
  {
    name: 'prompt',
    args: '{ message, title?, value? }',
    result: 'the text, or null',
    snippet: "prompt({ message: '${}' })",
  },
  { name: 'setTitle', args: 'text', result: 'true', snippet: "setTitle('${}')" },
  { name: 'download', args: 'name, text', result: 'true', snippet: "download('${name}', ${text})" },
  {
    name: 'navigate',
    args: '{ page } | { page, recordID } | { module, recordID }',
    result: 'true',
    snippet: "navigate({ page: '${}' })",
  },
  { name: 'refresh', args: '', result: 'true', snippet: 'refresh()' },
  { name: 'resize', args: 'height', result: 'true', snippet: 'resize(${height})' },
  {
    name: 'on',
    args: 'event, fn',
    result: 'a function that stops listening',
    snippet: "on('${refresh}', () => {\n  ${}\n})",
  },
  { name: 'ready', result: 'Promise<boolean>', snippet: 'ready' },
  { name: 'call', args: 'op, args', result: 'the call’s answer', snippet: "call('${}', { ${} })" },
]

// What `human.context()` answers.
export const CONTEXT_KEYS = ['namespaceID', 'namespace', 'pageID', 'moduleID', 'recordID', 'params']

const byName = Object.fromEntries(BRIDGE_CALLS.map(c => [c.name, c]))

export function noteKey(name) {
  return `app.bridge.${name.replace(/\./g, '_')}`
}

export function signature(call) {
  if (call.group) return `human.${call.name}.{ ${call.group.join(', ')} }`
  if (call.args === undefined) return `human.${call.name} → ${call.result}`
  return `human.${call.name}(${call.args}) → ${call.result}`
}

// Inside JavaScript: a <script> element or an on… attribute.
export function inScript(state, pos) {
  for (let n = syntaxTree(state).resolveInner(pos, -1); n; n = n.parent) {
    if (n.name === 'Script') return true
  }
  return false
}

// The names in the page that hold what `human.context()` answered, and those
// that hold its `params`.
function contextNames(text) {
  const whole = new Set()
  const params = new Set()
  const assigned = /(?:const|let|var)\s+(\w+)\s*=\s*await\s+human\.context\(\)/g
  for (const m of text.matchAll(assigned)) whole.add(m[1])
  const destructured = /(?:const|let|var)\s*\{([^}]*)\}\s*=\s*await\s+human\.context\(\)/g
  for (const m of text.matchAll(destructured)) {
    for (const part of m[1].split(',')) {
      const [key, alias] = part.split(':').map(s => s.trim())
      if (key === 'params') params.add(alias || key)
    }
  }
  for (const name of whole) {
    const fromWhole = new RegExp(`(?:const|let|var)\\s+(\\w+)\\s*=\\s*${name}\\.params\\b`, 'g')
    for (const m of text.matchAll(fromWhole)) params.add(m[1])
  }
  return { whole, params }
}

// What to suggest at `pos` in `text`, given what the page declared:
// `{modules, automations, chatbots, params}`, each a list of names. Answers
// `{from, options: [{label, kind, ...}]}` or null.
export function suggestAt(text, pos, declared = {}) {
  const lineStart = text.lastIndexOf('\n', pos - 1) + 1
  const before = text.slice(lineStart, pos)
  const after = text.slice(
    pos,
    text.indexOf('\n', pos) === -1 ? undefined : text.indexOf('\n', pos),
  )
  let m

  if ((m = /\bhuman\.(\w+)\.(\w*)$/.exec(before)) && byName[m[1]]?.group) {
    return {
      from: pos - m[2].length,
      options: byName[m[1]].group.map(n => ({ kind: 'call', call: byName[`${m[1]}.${n}`] })),
    }
  }

  if ((m = /\bhuman\.(\w*)$/.exec(before))) {
    return {
      from: pos - m[1].length,
      options: BRIDGE_CALLS.filter(c => !c.name.includes('.')).map(c => ({
        kind: 'call',
        call: c,
      })),
    }
  }

  if ((m = /\b(module|automation|chatbot)\s*:\s*['"`]([\w-]*)$/.exec(before))) {
    const names = declared[`${m[1]}s`] || []
    return {
      from: pos - m[2].length,
      options: names.map(label => ({ kind: m[1], label })),
    }
  }

  const { whole, params } = contextNames(text)
  const paramOptions = () => (declared.params || []).map(label => ({ kind: 'param', label }))

  if ((m = /\b(\w+)\.params\.(\w*)$/.exec(before)) && whole.has(m[1])) {
    return { from: pos - m[2].length, options: paramOptions() }
  }

  if ((m = /\(await\s+human\.context\(\)\)\.params\.(\w*)$/.exec(before))) {
    return { from: pos - m[1].length, options: paramOptions() }
  }

  if ((m = /\b(\w+)\.(\w*)$/.exec(before)) && params.has(m[1])) {
    return { from: pos - m[2].length, options: paramOptions() }
  }

  if (
    ((m = /\b(\w+)\.(\w*)$/.exec(before)) && whole.has(m[1])) ||
    (m = /\(await\s+human\.context\(\)\)\.()(\w*)$/.exec(before))
  ) {
    return {
      from: pos - m[2].length,
      options: CONTEXT_KEYS.map(label => ({ kind: 'context', label })),
    }
  }

  // Inside `{ … } = await human.context()`.
  if (
    (m = /(?:const|let|var)\s*\{([\w\s,:]*?)(\w*)$/.exec(before)) &&
    /^\w*[\w\s,:]*\}\s*=\s*await\s+human\.context\(\)/.test(after)
  ) {
    const taken = new Set(m[1].split(',').map(s => s.split(':')[0].trim()))
    return {
      from: pos - m[2].length,
      options: CONTEXT_KEYS.filter(k => !taken.has(k)).map(label => ({ kind: 'context', label })),
    }
  }

  return null
}

function callOption(call, t) {
  const label = call.name.split('.').pop()
  const info = `${signature(call)}\n\n${t(noteKey(call.name))}`
  if (call.group) {
    return {
      label,
      type: 'namespace',
      detail: `{ ${call.group.join(', ')} }`,
      info,
      apply: (view, _c, from, to) => {
        view.dispatch({
          changes: { from, to, insert: `${label}.` },
          selection: { anchor: from + label.length + 1 },
        })
        startCompletion(view)
      },
    }
  }
  return snippetCompletion(call.snippet, {
    label,
    type: call.args === undefined ? 'property' : 'function',
    detail: call.args === undefined ? '' : `(${call.args})`,
    info,
  })
}

// The CodeMirror extensions an editor takes as its `assist`: suggestions and
// hover notes for the bridge. `declared()` answers what the page declared, read
// fresh on every suggestion.
export function bridgeAssist(t, declared = () => ({})) {
  const source = cx => {
    if (!inScript(cx.state, cx.pos)) return null
    const found = suggestAt(cx.state.doc.toString(), cx.pos, declared())
    if (!found || !found.options.length) return null
    return {
      from: found.from,
      validFor: /^[\w-]*$/,
      options: found.options.map(o => {
        if (o.kind === 'call') return callOption(o.call, t)
        return {
          label: o.label,
          type: o.kind === 'context' || o.kind === 'param' ? 'property' : 'constant',
          detail: t(`app.bridge.kind.${o.kind}`),
          info: o.kind === 'context' ? t(`app.bridge.contextKey.${o.label}`) : undefined,
        }
      }),
    }
  }

  const hover = hoverTooltip((view, pos) => {
    if (!inScript(view.state, pos)) return null
    const line = view.state.doc.lineAt(pos)
    for (const m of line.text.matchAll(/\bhuman\.(\w+)(?:\.(\w+))?/g)) {
      const start = line.from + m.index
      const end = start + m[0].length
      if (pos < start || pos > end) continue
      const call = (m[2] && byName[`${m[1]}.${m[2]}`]) || byName[m[1]]
      if (!call) return null
      return {
        pos: start,
        end,
        above: true,
        create() {
          const dom = document.createElement('div')
          dom.className = 'cm-bridge-hint'
          const sig = document.createElement('code')
          sig.textContent = signature(call)
          const note = document.createElement('div')
          note.textContent = t(noteKey(call.name))
          dom.append(sig, note)
          return { dom }
        },
      }
    }
    return null
  })

  return [
    // Global language data, so it reaches the JavaScript nested in the HTML.
    EditorState.languageData.of(() => [{ autocomplete: source }]),
    hover,
    EditorView.theme({
      '.cm-bridge-hint': {
        display: 'flex',
        flexDirection: 'column',
        gap: '4px',
        maxWidth: '32rem',
        padding: '6px 8px',
        whiteSpace: 'normal',
        overflowWrap: 'break-word',
      },
      '.cm-completionInfo': {
        whiteSpace: 'pre-line',
        maxWidth: '32rem',
      },
    }),
  ]
}

// The examples the editors' Insert menu offers, each labelled
// `app.snippets.<key>`. `module` is a module the page declared, or a stand-in
// when it declared none.
export function snippets(module = 'module_handle') {
  return [
    {
      key: 'starter',
      text: `<!doctype html>
<html>
  <head>
    <meta charset="utf-8" />
    <style>
      body { font-family: system-ui, sans-serif; margin: 0; padding: 16px; }
    </style>
  </head>
  <body>
    <ul id="list"></ul>
    <script>
      async function load() {
        const { records } = await human.records.list({ module: '${module}', limit: 20 })
        const list = document.getElementById('list')
        list.replaceChildren(
          ...records.map(r => {
            const li = document.createElement('li')
            li.textContent = r.recordID
            return li
          }),
        )
      }
      human.on('refresh', load)
      load()
    </script>
  </body>
</html>
`,
    },
    {
      key: 'context',
      text: `const { recordID, params } = await human.context()
`,
    },
    {
      key: 'readRecord',
      text: `const { recordID } = await human.context()
const { record, refs } = await human.records.read({ module: '${module}', recordID })
`,
    },
    {
      key: 'listAll',
      text: `const records = []
let pageCursor
do {
  const page = await human.records.list({ module: '${module}', limit: 200, pageCursor })
  records.push(...page.records)
  pageCursor = page.nextPageCursor
} while (pageCursor)
`,
    },
    {
      key: 'report',
      text: `const { rows, refs } = await human.records.report({ module: '${module}', dimension: 'field_name' })
// rows: [{ dimension_0: value, count }]; dimension_0 is null for records with no value
`,
    },
    {
      key: 'modules',
      text: `const modules = await human.modules()
const { fields } = modules.find(m => m.handle === '${module}')
`,
    },
    {
      key: 'create',
      text: `const { record } = await human.records.create({ module: '${module}', values: {} })
`,
    },
    {
      key: 'update',
      text: `const { record } = await human.records.update({ module: '${module}', recordID, values: {} })
`,
    },
    {
      key: 'confirmToast',
      text: `if (await human.confirm({ message: 'Are you sure?' })) {
  human.toast({ message: 'Done', severity: 'success' })
}
`,
    },
    {
      key: 'refresh',
      text: `human.on('refresh', () => load())
`,
    },
  ]
}
