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
import {
  mkdtempSync,
  mkdirSync,
  rmSync,
  writeFileSync,
  copyFileSync,
  readFileSync,
  existsSync,
} from 'node:fs'
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

// `up` starts each service in a subshell that forked from this script, and a
// fork inherits the caller's stdout. A caller reading `up` through a pipe then
// waits on that copy for as long as the service runs — the command looks hung
// and is killed by whatever timeout the caller has.
test('a started service does not hold the caller end of a pipe', async t => {
  const dir = mkdtempSync(join(tmpdir(), 'wt-svc-'))
  const pidfile = join(dir, 'svc.pid')
  const runner = join(dir, 'run.sh')
  t.after(() => {
    try {
      const leader = readFileSync(pidfile, 'utf8').trim()
      execFileSync('bash', ['-c', `kill -TERM -${leader} 2>/dev/null || true`])
    } catch {}
    rmSync(dir, { recursive: true, force: true })
  })

  // The shipped definition, lifted out of cmd_up rather than restated here.
  writeFileSync(
    runner,
    'set -euo pipefail\n' +
      `cd ${JSON.stringify(HERE)}\n` +
      `eval "$(awk '/^  start_svc\\(\\)/{f=1} f{print} f&&/^  \\}$/{exit}' worktree.sh | sed -E 's/^  //')"\n` +
      `start_svc ${JSON.stringify(pidfile)} ${JSON.stringify(join(dir, 'svc.log'))} ${JSON.stringify(dir)} sleep 120\n` +
      'echo returning\n',
  )

  // `cat` is the caller: it reads until every writer has let go. timeout's exit
  // code says whether one never did.
  const out = execFileSync(
    'bash',
    ['-c', `timeout 6 bash -c 'bash ${JSON.stringify(runner)} | cat'; echo "rc=$?"`],
    { encoding: 'utf8', timeout: 30000 },
  )

  assert.match(out, /returning/, 'start_svc never returned')
  assert.match(out, /rc=0/, 'the pipe stayed open after the service started — rc=124 is timeout')
  assert.ok(await waitFor(() => existsSync(pidfile)), 'the service never recorded its leader')
})

// `rm` and `land` delete the checkout they are often running inside, and
// $REPO_DIR is that checkout. Every database credential comes from the
// primary's DSN, which is found by asking git where the primary is — so asking
// after the removal answers nothing, and dropdb falls back to prompting for a
// password that no one is there to type.
test('rm from inside the worktree still resolves the database credentials', async t => {
  const root = mkdtempSync(join(tmpdir(), 'wt-rm-'))
  const primary = join(root, 'primary')
  const wt = join(root, 'wt')
  const bin = join(root, 'bin')
  const record = join(root, 'dropdb.args')
  const api = freePort()
  const vite = freePort(new Set([api]))
  t.after(() => rmSync(root, { recursive: true, force: true }))

  const git = (...args) => execFileSync('git', ['-C', primary, ...args], { stdio: 'ignore' })

  // A primary that is a real git repo: the removal has to be a real worktree
  // removal for the lookup to lose its answer.
  mkdirSync(join(primary, 'dev', 'agent'), { recursive: true })
  mkdirSync(join(primary, 'server'), { recursive: true })
  mkdirSync(join(primary, 'client', 'web', 'unify'), { recursive: true })
  for (const f of ['worktree.sh', 'common.sh', 'stack.sh', 'tty.sh'])
    copyFileSync(join(HERE, f), join(primary, 'dev', 'agent', f))
  writeFileSync(
    join(primary, 'server', '.env'),
    `HTTP_ADDR=:${api}\nDB_DSN=postgres://alice:s3cret@db.example:6543/base\n`,
  )
  writeFileSync(
    join(primary, 'client', 'web', 'unify', '.env.e2e'),
    `E2E_BASE_URL=http://localhost:${vite}\n`,
  )
  execFileSync('git', ['init', '-q', primary], { stdio: 'ignore' })
  git('add', '-A')
  git('-c', 'user.email=t@t', '-c', 'user.name=t', 'commit', '-qm', 'fixture')
  git('worktree', 'add', '-q', '-b', 'fake', wt)

  // .state is a symlink into the primary in a real worktree, and that is what
  // makes the ledger outlive the checkout.
  mkdirSync(join(primary, 'dev', 'agent', '.state', 'worktrees'), { recursive: true })
  execFileSync('ln', [
    '-s',
    join(primary, 'dev', 'agent', '.state'),
    join(wt, 'dev', 'agent', '.state'),
  ])
  writeFileSync(
    join(primary, 'dev', 'agent', '.state', 'worktrees', 'fake.json'),
    JSON.stringify(
      { name: 'fake', path: wt, slot: 7, branch: 'fake', api, vite, db: 'fixture_wt7' },
      null,
      2,
    ),
  )
  writeFileSync(join(primary, 'dev', 'agent', 'backlog.sh'), '#!/usr/bin/env bash\necho 0\n', {
    mode: 0o755,
  })

  // The stub stands in for the real thing and records what it was handed: an
  // empty -U with no password is exactly the state that makes libpq prompt.
  mkdirSync(bin, { recursive: true })
  writeFileSync(
    join(bin, 'dropdb'),
    '#!/usr/bin/env bash\nprintf "PGPASSWORD=%s\\n" "${PGPASSWORD-unset}" >"$RECORD"\n' +
      'printf "%s\\n" "$@" >>"$RECORD"\n',
    { mode: 0o755 },
  )

  const out = execFileSync(
    'bash',
    [
      '-c',
      `cd ${JSON.stringify(wt)} && PATH=${JSON.stringify(bin)}:$PATH dev/agent/worktree.sh rm fake`,
    ],
    { encoding: 'utf8', timeout: 60000, env: { ...process.env, RECORD: record } },
  )

  assert.match(out, /removed 'fake'/)
  const args = readFileSync(record, 'utf8')
  assert.match(args, /PGPASSWORD=s3cret/, `dropdb got no password: ${args}`)
  assert.match(args, /^-U$\nalice$/m, `dropdb got the wrong user: ${args}`)
  assert.match(args, /^-w$/m, 'dropdb may still prompt')
  assert.match(args, /^fixture_wt7$/m)
  assert.ok(!existsSync(wt), 'the checkout survived')
})
