// Run: node --test dev/agent/logs.test.mjs
//
// logs.sh reads errors first: the default view keeps error-level entries and
// the [devwatch] markers that place them, hides the warn banner a restart
// repeats, and says what it hid. -w adds warnings, -a is the raw tail.
import { test } from 'node:test'
import assert from 'node:assert/strict'
import { spawnSync } from 'node:child_process'
import { writeFileSync, mkdtempSync } from 'node:fs'
import { tmpdir } from 'node:os'
import { fileURLToPath } from 'node:url'
import path from 'node:path'

const here = path.dirname(fileURLToPath(import.meta.url))
const script = path.join(here, 'logs.sh')

const sample =
  [
    '[devwatch] 10:23:54 building',
    '[devwatch] 10:24:01 serving /x/dev-bin',
    JSON.stringify({
      level: 'warn',
      ts: 1790928805.5,
      logger: 'dal.connection',
      msg: 'could not connect',
    }),
    JSON.stringify({ level: 'warn', ts: 1790928805.7, msg: 'no SMTP servers found' }),
    JSON.stringify({ level: 'info', ts: 1790928806.0, msg: 'listening' }),
    JSON.stringify({
      level: 'error',
      ts: 1790928807.0,
      logger: 'compose.record',
      msg: 'create failed',
      error: 'not unique',
    }),
    'panic: runtime error: index out of range',
  ].join('\n') + '\n'

const dir = mkdtempSync(path.join(tmpdir(), 'logs-test-'))
const log = path.join(dir, 'dev.log')
writeFileSync(log, sample)

function run(...args) {
  const r = spawnSync('bash', [script, ...args], {
    env: { ...process.env, HUMAN_DEV_LOG: log, NO_COLOR: '1' },
    encoding: 'utf8',
  })
  return { status: r.status, lines: r.stdout.split('\n').filter(Boolean) }
}

test('the default keeps errors and markers, hides warn and info, and says so', () => {
  const r = run()
  assert.equal(r.status, 0)
  assert.deepEqual(r.lines.slice(0, 2), [
    '[devwatch] 10:23:54 building',
    '[devwatch] 10:24:01 serving /x/dev-bin',
  ])
  assert.match(r.lines[2], /ERROR compose\.record: create failed — not unique$/)
  assert.equal(r.lines[3], 'panic: runtime error: index out of range')
  assert.match(r.lines[4], /^— hidden: 2 warn, 1 info; -w shows warnings/)
  assert.equal(r.lines.length, 5)
})

test('-w adds the warnings in order', () => {
  const r = run('-w')
  assert.match(r.lines[2], /WARN dal\.connection: could not connect$/)
  assert.match(r.lines[3], /WARN no SMTP servers found$/)
  assert.match(r.lines[4], /ERROR/)
  assert.match(r.lines[6], /^— hidden: 1 info;/)
})

test('-a is the raw tail', () => {
  const r = run('-a')
  assert.equal(r.lines.length, 7)
  assert.equal(r.lines[2], sample.split('\n')[2])
})

test('a pattern greps the raw tail as before', () => {
  const r = run('smtp')
  assert.equal(r.lines.length, 1)
  assert.match(r.lines[0], /SMTP/)
})

test('a quiet log says there was nothing to show', () => {
  writeFileSync(
    path.join(dir, 'quiet.log'),
    JSON.stringify({ level: 'info', ts: 1, msg: 'ok' }) + '\n',
  )
  const r = spawnSync('bash', [script], {
    env: { ...process.env, HUMAN_DEV_LOG: path.join(dir, 'quiet.log') },
    encoding: 'utf8',
  })
  assert.equal(r.stdout.trim(), 'no errors in the last 400 lines (1 info hidden)')
})
