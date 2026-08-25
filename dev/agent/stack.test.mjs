// Tests for the stack resolver. No browser and no server:
//
//   node --test dev/agent/stack.test.mjs
//
// What these pin is one thing: a checkout answers for ITS OWN stack. Every
// script here used to carry the literal 1043/5173, so a call made from a
// worktree reached the primary — same JWT secret, same user IDs, different
// database — and came back looking like a pass.
import { test } from 'node:test'
import assert from 'node:assert/strict'
import { execFileSync } from 'node:child_process'
import { mkdtempSync, mkdirSync, rmSync, writeFileSync } from 'node:fs'
import { tmpdir } from 'node:os'
import { dirname, join } from 'node:path'
import { fileURLToPath } from 'node:url'

const HERE = dirname(fileURLToPath(import.meta.url))
const STACK = join(HERE, 'stack.sh')

// resolve({...files}) — build a throwaway checkout and read what stack.sh
// makes of it. `env` goes to the process, so it also covers the override rule.
function resolve({ serverEnv, e2e, playwright, registry, env = {} } = {}) {
  const root = mkdtempSync(join(tmpdir(), 'stack-'))
  const unify = join(root, 'client', 'web', 'unify')
  mkdirSync(join(root, 'server'), { recursive: true })
  mkdirSync(unify, { recursive: true })

  if (serverEnv !== undefined) writeFileSync(join(root, 'server', '.env'), serverEnv)
  if (e2e !== undefined) writeFileSync(join(unify, '.env.e2e'), e2e)
  if (playwright !== undefined) writeFileSync(join(unify, 'playwright.config.ts'), playwright)
  if (registry !== undefined) {
    const dir = join(root, 'dev', 'agent', '.state', 'worktrees')
    mkdirSync(dir, { recursive: true })
    writeFileSync(join(dir, 'some-lane.json'), registry(root))
  }

  try {
    const out = execFileSync('bash', [STACK], {
      encoding: 'utf8',
      env: {
        ...process.env,
        HUMAN_API: '',
        HUMAN_WEBAPP: '',
        HUMAN_GIN: '',
        STACK_ROOT: root,
        ...env,
      },
    })
    return Object.fromEntries(
      out
        .split('\n')
        .filter(Boolean)
        .map(l => [l.slice(0, l.indexOf('=')), l.slice(l.indexOf('=') + 1)]),
    )
  } finally {
    rmSync(root, { recursive: true, force: true })
  }
}

test('the API port comes from the checkout, not from a literal', async t => {
  const cases = [
    {
      name: "a worktree's own port, not the primary's",
      serverEnv: 'DOMAIN=localhost:1543\nHTTP_ADDR=:1543\nENVIRONMENT=dev\n',
      want: 'http://localhost:1543/api',
    },
    {
      name: 'a host-qualified bind address',
      serverEnv: 'HTTP_ADDR=127.0.0.1:1743\n',
      want: 'http://localhost:1743/api',
    },
    {
      name: 'a commented assignment is not the value',
      serverEnv: '#HTTP_ADDR=:9999\nHTTP_ADDR=:1343\n',
      want: 'http://localhost:1343/api',
    },
    {
      // godotenv parses the whole file into a map before applying it, so the
      // server listens on the LAST assignment. Reading the first one reports a
      // port nothing is bound to the moment somebody appends an override.
      name: 'the last assignment wins, as the server sees it',
      serverEnv: 'HTTP_ADDR=:1343\nHTTP_ADDR=:1443\n',
      want: 'http://localhost:1443/api',
    },
    {
      name: 'an empty assignment falls back rather than resolving to nothing',
      serverEnv: 'HTTP_ADDR=\n',
      want: 'http://localhost:1043/api',
    },
    {
      name: 'no env file at all — the checkout is not set up yet',
      want: 'http://localhost:1043/api',
    },
  ]

  for (const c of cases) {
    await t.test(c.name, () => {
      assert.equal(resolve({ serverEnv: c.serverEnv }).HUMAN_API, c.want)
    })
  }
})

test('HTTP_API_BASE_URL is part of the API URL and not of the base', () => {
  const s = resolve({ serverEnv: 'HTTP_ADDR=:1543\nHTTP_API_BASE_URL=/gateway\n' })
  assert.equal(s.HUMAN_API, 'http://localhost:1543/gateway')
  assert.equal(s.HUMAN_BASE, 'http://localhost:1543')
  assert.equal(s.HUMAN_AUTH, 'http://localhost:1543/auth')
})

test('the webapp URL comes from .env.e2e, whose fallback is the config literal', async t => {
  await t.test("a worktree's own vite port", () => {
    const s = resolve({
      e2e: 'E2E_BASE_URL=http://localhost:5174\nE2E_USER=agent@local.dev\n',
      playwright: "baseURL: process.env.E2E_BASE_URL || 'http://localhost:5173',\n",
    })
    assert.equal(s.HUMAN_WEBAPP, 'http://localhost:5174')
  })

  await t.test('no .env.e2e yet — the registry, not the config literal', () => {
    // dev/setup.sh generating .env.e2e for a lane reads this. The literal in
    // playwright.config.ts is 5173 in every checkout, so trusting it here
    // writes the primary's webapp into a worktree's own config.
    const s = resolve({
      serverEnv: 'HTTP_ADDR=:1143\n',
      playwright: "baseURL: process.env.E2E_BASE_URL || 'http://localhost:5173',\n",
      registry: root =>
        JSON.stringify({ name: 'lane', path: root, api: 1143, vite: 5174 }, null, 2),
    })
    assert.equal(s.HUMAN_WEBAPP, 'http://localhost:5174')
  })

  await t.test('no .env.e2e and no registry — the config literal', () => {
    const s = resolve({
      playwright: "baseURL: process.env.E2E_BASE_URL || 'http://localhost:5179',\n",
    })
    assert.equal(s.HUMAN_WEBAPP, 'http://localhost:5179')
  })

  await t.test('neither — the primary', () => {
    assert.equal(resolve().HUMAN_WEBAPP, 'http://localhost:5173')
  })
})

test("gin's proxy port is only recorded in the worktree registry", async t => {
  await t.test('a registered lane', () => {
    const s = resolve({
      serverEnv: 'HTTP_ADDR=:1143\n',
      registry: root =>
        JSON.stringify({ name: 'lane', path: root, slot: 1, api: 1143, gin: 3101 }, null, 2),
    })
    assert.equal(s.HUMAN_GIN, 'http://localhost:3101')
  })

  await t.test('an entry for a different checkout is not this one', () => {
    const s = resolve({
      registry: () =>
        JSON.stringify({ name: 'other', path: '/elsewhere', api: 1143, gin: 3201 }, null, 2),
    })
    assert.equal(s.HUMAN_GIN, 'http://localhost:3001')
  })

  await t.test('a stale entry, for the port this checkout no longer serves', () => {
    const s = resolve({
      serverEnv: 'HTTP_ADDR=:1943\n',
      registry: root => JSON.stringify({ name: 'lane', path: root, api: 1143, gin: 3101 }, null, 2),
    })
    assert.equal(s.HUMAN_GIN, 'http://localhost:3001')
  })

  await t.test('the primary, with no registry at all', () => {
    assert.equal(resolve().HUMAN_GIN, 'http://localhost:3001')
  })
})

test('an explicit environment value outranks the files', () => {
  const s = resolve({
    serverEnv: 'HTTP_ADDR=:1543\n',
    e2e: 'E2E_BASE_URL=http://localhost:5174\n',
    env: { HUMAN_API: 'http://127.0.0.1:1943/api', HUMAN_WEBAPP: 'http://127.0.0.1:5999' },
  })
  assert.equal(s.HUMAN_API, 'http://127.0.0.1:1943/api')
  assert.equal(s.HUMAN_BASE, 'http://127.0.0.1:1943')
  assert.equal(s.HUMAN_WEBAPP, 'http://127.0.0.1:5999')
})
