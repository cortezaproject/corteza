// Tests for the teardown half of worktree.sh. No repo, no database, no
// server:
//
//   node --test dev/agent/worktree.test.mjs
//
// What these pin is that `down` frees the slot's ports. The dev server runs
// under a watcher that gives it a process group of its own, so the group that
// `up` recorded does not contain the thing holding the port — and the watcher
// is piped into `tee`, so the same group signal takes the watcher down before
// it can stop its child. `down` reported success against that and left the
// port bound; the next `up` on the slot then found it busy.
import { test } from 'node:test'
import assert from 'node:assert/strict'
import { execFileSync } from 'node:child_process'
import { mkdtempSync, mkdirSync, rmSync, writeFileSync, copyFileSync, readFileSync } from 'node:fs'
import { tmpdir } from 'node:os'
import { dirname, join } from 'node:path'
import { fileURLToPath } from 'node:url'

const HERE = dirname(fileURLToPath(import.meta.url))

const sleep = ms => new Promise(r => setTimeout(r, ms))

// Asking node for a free port leaves a listening handle behind, and the test
// runner then waits on it forever. Picking one and checking it with the same
// tool worktree.sh uses costs nothing and leaves nothing open.
function freePort(taken = new Set()) {
  for (let i = 0; i < 200; i++) {
    const port = 20000 + Math.floor(Math.random() * 9000)
    if (!taken.has(port) && !portBusy(port)) {
      taken.add(port)
      return port
    }
  }
  throw new Error('no free port found')
}

function portBusy(port) {
  return (
    execFileSync('bash', ['-c', `ss -ltn "sport = :${port}" | grep -c LISTEN || true`])
      .toString()
      .trim() !== '0'
  )
}

// A checkout with just enough of dev/agent for worktree.sh to run, plus one
// registry entry naming ports nothing else uses.
function checkout(api, vite) {
  const root = mkdtempSync(join(tmpdir(), 'wt-down-'))
  const agent = join(root, 'dev', 'agent')
  mkdirSync(join(agent, '.state', 'worktrees'), { recursive: true })
  for (const f of ['worktree.sh', 'common.sh', 'stack.sh', 'tty.sh'])
    copyFileSync(join(HERE, f), join(agent, f))
  const path = join(root, 'wt')
  mkdirSync(join(path, '.run'), { recursive: true })
  // stack.sh resolves this checkout's ports from these two, and refuses to
  // source without them.
  mkdirSync(join(root, 'server'), { recursive: true })
  mkdirSync(join(root, 'client', 'web', 'unify'), { recursive: true })
  writeFileSync(join(root, 'server', '.env'), `HTTP_ADDR=:${api}\n`)
  writeFileSync(
    join(root, 'client', 'web', 'unify', '.env.e2e'),
    `E2E_BASE_URL=http://localhost:${vite}\n`,
  )
  writeFileSync(
    join(agent, '.state', 'worktrees', 'fake.json'),
    JSON.stringify({ name: 'fake', path, slot: 7, branch: 'fake', api, vite, db: 'none' }, null, 2),
  )
  return { root, agent, path }
}

// The shape that broke it: a group leader that dies the instant it is
// signalled (make), and the process actually holding the port sitting in a
// process group of its own (the watcher's child).
function startService(path, what, port) {
  const pidfile = join(path, '.run', `${what}.pid`)
  const holder =
    'import socket,time; ' +
    's=socket.socket(); s.setsockopt(socket.SOL_SOCKET,socket.SO_REUSEADDR,1); ' +
    `s.bind(("0.0.0.0",${port})); s.listen(1); time.sleep(300)`
  // setsid twice: once for the wrapper group that `up` records, once for the
  // listener, which is what the watcher does to the server.
  const leader = `setsid python3 -c ${JSON.stringify(holder)} >/dev/null 2>&1 & exec sleep 300`
  const start =
    `cd ${JSON.stringify(path)} && ` +
    `setsid bash -c 'echo $$ >"$1"; shift; exec "$@"' _ ${JSON.stringify(pidfile)} ` +
    `bash -c ${JSON.stringify(leader)} >/dev/null 2>&1 &`
  execFileSync('bash', ['-c', start], { stdio: 'ignore' })
  return pidfile
}

async function waitFor(fn, ms = 8000) {
  const until = Date.now() + ms
  while (Date.now() < until) {
    if (fn()) return true
    await sleep(100)
  }
  return false
}

test('down frees a port held outside the recorded process group', async t => {
  const taken = new Set()
  const api = freePort(taken)
  const vite = freePort(taken)
  const { root, agent, path } = checkout(api, vite)
  t.after(() => {
    try {
      execFileSync('bash', [join(agent, 'worktree.sh'), 'down', 'fake', '--force'], {
        stdio: 'ignore',
      })
    } catch {}
    rmSync(root, { recursive: true, force: true })
  })

  const pidfile = startService(path, 'server', api)
  assert.ok(await waitFor(() => portBusy(api)), 'fixture never bound the port')

  // The leader on file is not the process holding the port — that is the whole
  // point of the fixture, and asserting it keeps the test honest if the dev
  // stack ever stops working that way.
  const leader = readFileSync(pidfile, 'utf8').trim()
  const leaderPgid = execFileSync('bash', ['-c', `ps -o pgid= -p ${leader} | tr -d ' '`])
    .toString()
    .trim()
  const holderPgid = execFileSync('bash', [
    '-c',
    `p=$(ss -ltnp "sport = :${api}" | sed -nE 's/.*pid=([0-9]+).*/\\1/p' | head -1); ps -o pgid= -p $p | tr -d ' '`,
  ])
    .toString()
    .trim()
  assert.notEqual(holderPgid, leaderPgid, 'fixture did not put the listener in its own group')

  const out = execFileSync('bash', [join(agent, 'worktree.sh'), 'down', 'fake'], {
    encoding: 'utf8',
    timeout: 60000,
  })

  assert.match(out, /stopped server/)
  assert.ok(await waitFor(() => !portBusy(api), 3000), `port ${api} still bound after down`)
})
