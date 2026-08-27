// Tests for the reporting helpers:
//
//   node --test dev/agent/tty.test.mjs
//
// Two things, and both have already broken once. Sourcing this file must not
// fail — every script here runs under `set -euo pipefail`, where a sourced
// file ending in a false test kills the caller with no message at all. And
// nothing may reach a pipe decorated: these outputs are read by tests, greps
// and log files, none of which want escape sequences.
import { test } from 'node:test'
import assert from 'node:assert/strict'
import { execFileSync } from 'node:child_process'
import { dirname, join } from 'node:path'
import { fileURLToPath } from 'node:url'

const HERE = dirname(fileURLToPath(import.meta.url))
const TTY = join(HERE, 'tty.sh')
const ESC = '['

// run(script, env) — a `set -euo pipefail` shell that sources tty.sh. stdout is
// a pipe here, which is exactly the undecorated case.
function run(script, env = {}) {
  return execFileSync('bash', ['-c', `set -euo pipefail; source ${TTY}; ${script}`], {
    encoding: 'utf8',
    env: { ...process.env, NO_COLOR: '', FORCE_COLOR: '', ...env },
  })
}

test('sourcing it succeeds whether or not colour is on', async t => {
  for (const [name, env] of [
    ['colour off — a pipe', {}],
    ['colour forced on', { FORCE_COLOR: '1' }],
    ['NO_COLOR set', { NO_COLOR: '1' }],
  ]) {
    await t.test(name, () => {
      assert.equal(run('echo sourced', env).trim(), 'sourced')
    })
  }
})

test('a pipe gets no escape sequences', () => {
  const out = run('section S; ok o; step l d; warn w; note n; bad b 2>&1')
  assert.ok(!out.includes(ESC), 'decorated a pipe')
  // …and the glyphs are still there: they are ordinary characters, and the
  // tests reading this output match on them.
  assert.match(out, /✓ o/)
  assert.match(out, /✗ b/)
})

test('a terminal gets them, and the precedence is stated', () => {
  assert.ok(run('ok o', { FORCE_COLOR: '1' }).includes(`${ESC}32m`))
  // FORCE_COLOR outranks NO_COLOR: the flag on the command beats the standing
  // preference, which is the rule every other tool follows.
  assert.ok(run('ok o', { FORCE_COLOR: '1', NO_COLOR: '1' }).includes(`${ESC}32m`))
})

test('GREP_COLOUR_WHEN follows the same switch', () => {
  assert.equal(run('echo "$GREP_COLOUR_WHEN"').trim(), 'never')
  assert.equal(run('echo "$GREP_COLOUR_WHEN"', { FORCE_COLOR: '1' }).trim(), 'always')
})
