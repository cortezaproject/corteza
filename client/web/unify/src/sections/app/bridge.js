// The in-page bridge, prefixed to a custom app's source.
//
// The same text the MCP skill hands an app author, so a page written to the
// skill already carries a copy of it: the snippet claims `window.human` only
// when nothing has, and the two copies together leave one object and send one
// handshake. No top-level `const`/`let` for the same reason — a second
// evaluation of a lexical declaration is a SyntaxError that takes the page
// with it.
//
// With no port inside 600 ms the page is running outside Human (an artifact
// preview), and calls fall back to whatever the page put in `window.SAMPLE`.
export const BRIDGE_SCRIPT = `
window.human = window.human || (function () {
  var port = null
  var seq = 0
  var waiting = {}
  var listeners = {}

  var ready = new Promise(function (resolve) {
    var timer = setTimeout(function () { resolve(false) }, 600)

    addEventListener('message', function (e) {
      if (e.source !== parent || !e.data || e.data.type !== 'human:port' || !e.ports[0]) return
      clearTimeout(timer)
      port = e.ports[0]
      port.onmessage = function (m) {
        if (m.data && m.data.event) {
          (listeners[m.data.event] || []).forEach(function (fn) {
            try { fn(m.data.payload) } catch (err) { setTimeout(function () { throw err }) }
          })
          return
        }
        var w = waiting[m.data.id]
        if (!w) return
        delete waiting[m.data.id]
        if (m.data.error) w.reject(new Error(m.data.error))
        else w.resolve(m.data.result)
      }
      resolve(true)
    })

    try {
      parent.postMessage({ type: 'human:hello', v: 2 }, '*')
    } catch (err) {
      resolve(false)
    }
  })

  function call (op, args) {
    return ready.then(function (live) {
      if (!live) {
        var sample = window.SAMPLE
        if (sample && sample[op]) return sample[op](args)
        return Promise.reject(new Error('not connected'))
      }
      return new Promise(function (resolve, reject) {
        var id = ++seq
        waiting[id] = { resolve: resolve, reject: reject }
        port.postMessage({ id: id, op: op, args: args })
      })
    })
  }

  return {
    ready: ready,
    call: call,
    on: function (event, fn) {
      (listeners[event] = listeners[event] || []).push(fn)
      return function () {
        listeners[event] = (listeners[event] || []).filter(function (f) { return f !== fn })
      }
    },
    records: {
      list: function (a) { return call('records.list', a) },
      read: function (a) { return call('records.read', a) },
      report: function (a) { return call('records.report', a) },
      create: function (a) { return call('records.create', a) },
      update: function (a) { return call('records.update', a) },
      delete: function (a) { return call('records.delete', a) },
      open: function (a) { return call('records.open', a) }
    },
    users: {
      search: function (a) { return call('users.search', a) }
    },
    toast: function (a) { return call('toast', a) },
    confirm: function (a) { return call('confirm', a) },
    prompt: function (a) { return call('prompt', a) },
    setTitle: function (text) { return call('title', { text: text }) },
    download: function (name, text) { return call('download', { name: name, text: text }) },
    modules: function () { return call('modules', {}) },
    user: function () { return call('user', {}) },
    theme: function () { return call('theme', {}) },
    context: function () { return call('context', {}) },
    navigate: function (a) { return call('navigate', a) },
    refresh: function () { return call('refresh', {}) },
    chatbot: {
      open: function (a) { return call('chatbot.open', a) },
      close: function (a) { return call('chatbot.close', a) }
    },
    automation: {
      run: function (a) { return call('automation.run', a) }
    },
    files: {
      list: function (a) { return call('files.list', a) },
      read: function (a) { return call('files.read', a) },
      upload: function (a) { return call('files.upload', a) }
    },
    resize: function (height) { return call('resize', { height: height }) }
  }
})()
`
