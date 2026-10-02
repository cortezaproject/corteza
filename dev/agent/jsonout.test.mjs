// Run: node --test dev/agent/jsonout.test.mjs
//
// jsonout.py is what api.sh prints through, so its shaping is what every
// reader of an API response sees: a --pick projection, the 20-item cut with
// its visible marker, and compact JSON when nothing is reading a terminal.
import { test } from 'node:test'
import assert from 'node:assert/strict'
import { spawnSync } from 'node:child_process'
import { fileURLToPath } from 'node:url'
import path from 'node:path'

const here = path.dirname(fileURLToPath(import.meta.url))
const script = path.join(here, 'jsonout.py')

function run(input, ...args) {
  const r = spawnSync('python3', [script, ...args], { input, encoding: 'utf8' })
  return { status: r.status, out: r.stdout, err: r.stderr }
}

const users = JSON.stringify({
  response: {
    set: Array.from({ length: 25 }, (_, i) => ({
      userID: String(i),
      email: `u${i}@x`,
      values: { name: `n${i}` },
    })),
  },
})

test('a scalar list projects one per line', () => {
  const r = run(users, '--pick', 'response.set[].email', '--limit', '0')
  assert.equal(r.status, 0)
  assert.equal(r.out.split('\n').filter(Boolean).length, 25)
  assert.equal(r.out.split('\n')[0], 'u0@x')
})

test('braces keep named keys, nested paths included', () => {
  const r = run(users, '--pick', 'response.set[0].{userID,values.name}', '--compact')
  assert.equal(r.out.trim(), '{"userID":"0","values.name":"n0"}')
})

test('an empty string prints as "" so a row never reads as missing', () => {
  const r = run('{"set":[{"slug":""},{"slug":"crm"}]}', '--pick', 'set[].slug')
  assert.deepEqual(r.out.split('\n').filter(Boolean), ['""', 'crm'])
})

test('a missing key is null, not an error', () => {
  const r = run(users, '--pick', 'response.nope.deeper')
  assert.equal(r.status, 0)
  assert.equal(r.out.trim(), 'null')
})

test('lists past the limit are cut with a visible marker', () => {
  const r = run(users, '--compact')
  const set = JSON.parse(r.out).response.set
  assert.equal(set.length, 21)
  assert.equal(set[20], '… 5 more (--all)')
})

test('--limit 0 keeps everything', () => {
  const r = run(users, '--limit', '0', '--compact')
  assert.equal(JSON.parse(r.out).response.set.length, 25)
})

test('the marker also reaches a projected list', () => {
  const r = run(users, '--pick', 'response.set[].email')
  const lines = r.out.split('\n').filter(Boolean)
  assert.equal(lines.length, 21)
  assert.equal(lines[20], '… 5 more (--all)')
})

test('compact is the default off a terminal; --pretty indents', () => {
  const compact = run('{"a":[1,2]}')
  assert.equal(compact.out.trim(), '{"a":[1,2]}')
  const pretty = run('{"a":[1,2]}', '--pretty')
  assert.match(pretty.out, /^\{\n  "a": \[/)
})

test('a non-JSON body exits 1 and prints nothing', () => {
  const r = run('<!doctype html><html>')
  assert.equal(r.status, 1)
  assert.equal(r.out, '')
})
