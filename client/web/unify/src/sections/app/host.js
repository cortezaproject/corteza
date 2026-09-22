// The sandbox around a custom app: the two documents it runs in, and the rules
// every call it makes goes through. See app.intent.md.
//
// The outer document is a `srcdoc` string, so it can import nothing — its
// script is generated here as self-contained text. It holds no token and makes
// no API call: it owns the port to the app and relays each call to the shell,
// which answers it by running `dispatch` with the real clients.

// The outer document's own policy. `frame-src 'none'` is what stops the inner
// document navigating itself somewhere else; a policy on the inner document
// cannot block its own navigation.
export const CSP_OUTER = "frame-src 'none'"

// The app's policy: no network of any kind, inline script and style only,
// images from data: URLs.
export const CSP_INNER =
  "default-src 'none'; " +
  "script-src 'unsafe-inline' https://cdnjs.cloudflare.com; " +
  "style-src 'unsafe-inline'; " +
  'img-src data:; ' +
  "connect-src 'none'; " +
  "form-action 'none'"

// The most records one call may ask for.
export const MAX_LIMIT = 500

// A string as a JS literal safe to sit inside a `<script>` body: `</script`
// anywhere in it would end the element it is written into.
function embed(value) {
  return JSON.stringify(value).replace(/<\//g, '<\\/')
}

// The outer document's script, with the shell's origin baked in — every
// message it sends the shell, and every one it accepts, names that origin.
export function hostScriptSource({ origin }) {
  return `(function () {
  var SHELL = ${embed(origin)}
  var seq = 0
  var pending = {}
  var loads = 0

  var frame = document.createElement('iframe')
  frame.setAttribute('sandbox', 'allow-scripts')
  frame.srcdoc = window.__humanInner

  // A second load is the app navigating itself: the document that answered the
  // handshake is gone, so the frame goes with it.
  frame.addEventListener('load', function () {
    if (++loads < 2) return
    frame.remove()
    parent.postMessage({ type: 'human:navigated' }, SHELL)
  })

  addEventListener('message', function (e) {
    var data = e.data || {}

    if (e.source === parent) {
      if (e.origin !== SHELL || data.type !== 'human:result') return
      var waiting = pending[data.id]
      if (!waiting) return
      delete pending[data.id]
      waiting.port.postMessage({ id: waiting.id, result: data.result, error: data.error })
      return
    }

    if (e.source !== frame.contentWindow || data.type !== 'human:hello') return

    var channel = new MessageChannel()
    channel.port1.onmessage = function (m) {
      var call = m.data || {}
      var id = ++seq
      pending[id] = { port: channel.port1, id: call.id }
      parent.postMessage({ type: 'human:call', id: id, op: call.op, args: call.args }, SHELL)
    }
    // The app's origin is opaque, so '*' is the only target there is.
    e.source.postMessage({ type: 'human:port' }, '*', [channel.port2])
  })

  document.body.appendChild(frame)
})()`
}

// The whole outer document: the inner document as a string, then the host that
// puts it in a sandboxed frame.
export function buildOuterDocument({ source, hostScript, bridgeScript, cspInner = CSP_INNER }) {
  const inner = `<!doctype html>
<meta charset="utf-8">
<meta http-equiv="Content-Security-Policy" content="${cspInner}">
<script>${bridgeScript}</script>
${source}`

  // An explicit body: the host script appends to it as it runs, and a document
  // written without one has no body element until the parser reaches content.
  return `<!doctype html>
<html>
<head>
<meta charset="utf-8">
<meta http-equiv="Content-Security-Policy" content="${CSP_OUTER}">
<style>html,body{margin:0;padding:0;height:100%;overflow:hidden}iframe{display:block;width:100%;height:100%;border:0}</style>
</head>
<body>
<script>window.__humanInner = ${embed(inner)}</script>
<script>${hostScript}</script>
</body>
</html>`
}

// The app-facing shape of a record: `values` keyed by field name, a repeated
// field as an array, a field with nothing in it absent.
export function reshapeRecord(record) {
  const values = {}

  for (const { name, value } of record?.values || []) {
    if (value === null || value === undefined || value === '') continue
    if (name in values) {
      values[name] = Array.isArray(values[name]) ? [...values[name], value] : [values[name], value]
    } else {
      values[name] = value
    }
  }

  return {
    recordID: record?.recordID,
    values,
    ownedBy: record?.ownedBy,
    createdAt: record?.createdAt,
    updatedAt: record?.updatedAt,
  }
}

// Whether the app may touch this module: the refusal text, or null.
export function allowModule(meta, module) {
  if (!module) return 'no module was named'
  if (!(meta?.modules || []).includes(module)) {
    return `module "${module}" is not declared for this app`
  }
  return null
}

export function capLimit(limit) {
  const asked = Number(limit)
  if (!Number.isFinite(asked) || asked <= 0) return undefined
  return Math.min(Math.floor(asked), MAX_LIMIT)
}

function moduleIDFor(ctx, module) {
  const refusal = allowModule(ctx.meta, module)
  if (refusal) throw new Error(refusal)

  const moduleID = ctx.moduleIDs?.[module]
  if (!moduleID) {
    throw new Error(`module "${module}" was not found in namespace "${ctx.meta?.namespace}"`)
  }
  return moduleID
}

// Every operation an app can reach. Runs in the shell, as the viewer, under
// the viewer's permissions.
export async function dispatch(op, args = {}, ctx) {
  switch (op) {
    case 'records.list': {
      const moduleID = moduleIDFor(ctx, args.module)
      const { set = [], filter = {} } = await ctx.compose.recordList({
        namespaceID: ctx.namespaceID,
        moduleID,
        query: args.filter,
        sort: args.sort,
        limit: capLimit(args.limit),
        pageCursor: args.pageCursor,
      })
      return {
        records: set.map(reshapeRecord),
        refs: ctx.refs ? await ctx.refs(set) : {},
        nextPageCursor: filter.nextPage || null,
      }
    }

    case 'records.read': {
      const moduleID = moduleIDFor(ctx, args.module)
      const record = await ctx.compose.recordRead({
        namespaceID: ctx.namespaceID,
        moduleID,
        recordID: args.recordID,
      })
      return {
        record: reshapeRecord(record),
        refs: ctx.refs ? await ctx.refs([record]) : {},
      }
    }

    case 'records.report': {
      const moduleID = moduleIDFor(ctx, args.module)
      return ctx.compose.recordReport({
        namespaceID: ctx.namespaceID,
        moduleID,
        metrics: args.metrics,
        dimensions: args.dimensions,
        filter: args.filter,
      })
    }

    case 'user':
      return ctx.user()

    case 'theme':
      return ctx.theme()

    case 'resize':
      ctx.resize(args.height)
      return true

    default:
      throw new Error(`operation "${op}" is not available to an app`)
  }
}
