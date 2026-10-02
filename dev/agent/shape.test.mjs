// Run: node --test dev/agent/shape.test.mjs
//
// shape.sh answers "what does this call take" from server/*/rest.yaml, the
// codegen input, so a wrong verb or a forgotten path segment is caught before
// the request goes out. These pin the three views and the traps overlay.
import { test } from 'node:test'
import assert from 'node:assert/strict'
import { spawnSync } from 'node:child_process'
import { fileURLToPath } from 'node:url'
import path from 'node:path'

const here = path.dirname(fileURLToPath(import.meta.url))
const script = path.join(here, 'shape.sh')

function run(...args) {
  const r = spawnSync('bash', [script, ...args], { encoding: 'utf8' })
  return { status: r.status, out: r.stdout, err: r.stderr }
}

test('no argument lists every service with its groups', () => {
  const r = run()
  assert.equal(r.status, 0)
  assert.match(r.out, /^compose: .*\brecord\b/m)
  assert.match(r.out, /^system: .*\buser\b/m)
})

test('a group lists each call with method and full path', () => {
  const r = run('compose/record')
  assert.match(
    r.out,
    /^update\s+POST\s+\/compose\/namespace\/\{namespaceID\}\/module\/\{moduleID\}\/record\/\{recordID\}$/m,
  )
  assert.match(r.out, /^list\s+GET\s+/m)
})

test('a call shows the verb, scope path params, body fields and traps', () => {
  const r = run('compose/record', 'update')
  assert.match(
    r.out,
    /^POST \/compose\/namespace\/\{namespaceID\}\/module\/\{moduleID\}\/record\/\{recordID\}/,
  )
  assert.match(r.out, /^path:\n\s+\* namespaceID/m)
  assert.match(r.out, /^body \(JSON\):\n(.*\n)*\s+values\s+types\.RecordValueSet/m)
  assert.match(r.out, /^traps:\n(.*\n)*.*PUT is a bare HTTP 405/m)
  assert.match(r.out, /values.*is an array of \{name, value\}/)
})

test('the trailing-slash trap is shown for collection calls only', () => {
  assert.doesNotMatch(run('compose/record', 'update').out, /trailing slash/)
  assert.match(run('compose/record', 'list').out, /trailing slash/)
})

test('the three-segment form and a case-insensitive group both resolve', () => {
  assert.equal(run('compose/record/update').out, run('compose/record', 'update').out)
  assert.match(run('compose/pagelayout').out, /^list\s+GET/m)
})

test('--grep finds calls by name or path across services', () => {
  const r = run('--grep', 'undelete')
  const lines = r.out.trim().split('\n')
  assert.ok(lines.length >= 3, r.out)
  assert.ok(
    lines.every(l => /undelete/i.test(l)),
    r.out,
  )
})

test('an unknown group names the ones nearby and exits non-zero', () => {
  const r = run('compose/recrd')
  assert.notEqual(r.status, 0)
  assert.match(r.err, /no group 'recrd' in compose/)
})
