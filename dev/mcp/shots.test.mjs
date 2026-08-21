// Tests for the screenshot retention rule. No browser and no server:
//
//   node --test dev/mcp/shots.test.mjs
import { test } from 'node:test'
import assert from 'node:assert/strict'
import { execFile } from 'node:child_process'
import { promisify } from 'node:util'
import { mkdtempSync, readdirSync, rmSync, utimesSync, writeFileSync } from 'node:fs'
import { tmpdir } from 'node:os'
import { fileURLToPath } from 'node:url'
import { join } from 'node:path'
import { GRACE_MS, KEEP, pruneShots, shotName } from './shots.mjs'

const HERE = fileURLToPath(new URL('.', import.meta.url))
const HOUR = 60 * 60 * 1000

// Shots named the way the driver names them, aged an hour apart, newest first.
// Going through shotName() is deliberate: a name the sweep no longer recognises
// is the whole bug back again, silently.
function seed(dir, count, now) {
  const names = []
  for (let i = 0; i < count; i++) {
    const at = now - i * HOUR
    const name = shotName(at, 1000 + i)
    writeFileSync(join(dir, name), 'x')
    utimesSync(join(dir, name), at / 1000, at / 1000)
    names.push(name)
  }
  return names
}

function withDir(fn) {
  const dir = mkdtempSync(join(tmpdir(), 'shots-'))
  try {
    fn(dir)
  } finally {
    rmSync(dir, { recursive: true, force: true })
  }
}

test('keeps the newest and deletes the rest', () => {
  withDir(dir => {
    const now = Date.now()
    const names = seed(dir, 10, now)

    const deleted = pruneShots(dir, { keep: 3, now })

    assert.deepEqual(readdirSync(dir).sort(), names.slice(0, 3).sort())
    assert.deepEqual(deleted.sort(), names.slice(3).sort())
  })
})

test('never deletes a shot younger than the grace window', () => {
  withDir(dir => {
    const now = Date.now()
    // Every one of these is minutes old, and there are far more than keep
    // allows — a burst of concurrent sessions, none of them finished reading.
    const fresh = seed(dir, 8, now).map((name, i) => {
      const at = now - i * 60 * 1000
      utimesSync(join(dir, name), at / 1000, at / 1000)
      return name
    })

    const deleted = pruneShots(dir, { keep: 2, now })

    assert.deepEqual(deleted, [])
    assert.equal(readdirSync(dir).length, fresh.length)
  })
})

test('sweeps the old and spares the fresh in the same directory', () => {
  withDir(dir => {
    const now = Date.now()
    const old = seed(dir, 4, now - GRACE_MS - HOUR)
    const fresh = seed(dir, 2, now)

    const deleted = pruneShots(dir, { keep: 2, now })

    // The two fresh ones fill keep, so all four old ones are candidates.
    assert.deepEqual(deleted.sort(), old.sort())
    assert.deepEqual(readdirSync(dir).sort(), fresh.sort())
  })
})

test('leaves files it did not write alone', () => {
  withDir(dir => {
    const now = Date.now()
    const strangers = ['notes.txt', 'ui-session-http-localhost-5173.json', 'keep-me.png']
    for (const name of strangers) {
      writeFileSync(join(dir, name), 'x')
      const at = now - 10 * HOUR
      utimesSync(join(dir, name), at / 1000, at / 1000)
    }
    seed(dir, 3, now - GRACE_MS - HOUR)

    const deleted = pruneShots(dir, { keep: 0, now })

    assert.equal(deleted.length, 3)
    assert.deepEqual(readdirSync(dir).sort(), [...strangers].sort())
  })
})

test('a directory that does not exist is not an error', () => {
  assert.deepEqual(pruneShots(join(tmpdir(), 'shots-does-not-exist-' + Date.now())), [])
})

test('the default rule bounds the directory', () => {
  withDir(dir => {
    const now = Date.now()
    seed(dir, KEEP + 25, now - GRACE_MS - HOUR)

    pruneShots(dir, { now })

    assert.equal(readdirSync(dir).length, KEEP)
  })
})

// The property the unique naming exists for: whatever the retention rule does,
// two sessions taking a screenshot at the same moment must not write to one
// path. Real processes, because the pid is what carries it.
test('concurrent processes get distinct names', async () => {
  const run = promisify(execFile)
  const script = `import { shotName } from ${JSON.stringify(join(HERE, 'shots.mjs'))}
process.stdout.write(shotName())`

  const results = await Promise.all(
    Array.from({ length: 8 }, () => run(process.execPath, ['--input-type=module', '-e', script])),
  )

  const names = results.map(r => r.stdout.trim())
  assert.equal(new Set(names).size, names.length, `collided: ${names.join(' ')}`)
})
