// Run: node --test dev/agent/claude-launch.test.mjs
//
// dev/claude.sh attaches the human-local MCP server only when asked. A session
// launched plainly must not carry the ~150 configurator tools, and one launched
// with HUMAN_MCP=1 or --human-mcp must carry exactly one --mcp-config pointing
// at dev/human-local.mcp.json. CLAUDE_BIN=echo turns the final exec into a
// printout of the argv that would have reached claude.
import { test } from 'node:test'
import assert from 'node:assert/strict'
import { spawnSync } from 'node:child_process'
import { fileURLToPath } from 'node:url'
import path from 'node:path'

const here = path.dirname(fileURLToPath(import.meta.url))
const launcher = path.join(here, '..', 'claude.sh')
const config = path.join(here, '..', 'human-local.mcp.json')

function launch(args, env = {}) {
  const r = spawnSync('bash', [launcher, ...args], {
    env: { ...process.env, CLAUDE_BIN: 'echo', HUMAN_MCP: '', ...env },
    encoding: 'utf8',
  })
  assert.equal(r.status, 0, r.stderr)
  return r.stdout.trim().split(/\s+/).filter(Boolean)
}

test('a plain launch passes the flags through and attaches no MCP config', () => {
  const argv = launch(['--model', 'opus'])
  assert.deepEqual(argv, ['--model', 'opus'])
})

test('HUMAN_MCP=1 attaches dev/human-local.mcp.json once, ahead of the flags', () => {
  const argv = launch(['--verbose'], { HUMAN_MCP: '1' })
  assert.deepEqual(argv, ['--mcp-config', config, '--verbose'])
})

test('--human-mcp is consumed by the launcher and does the same', () => {
  const argv = launch(['--human-mcp', '--verbose'])
  assert.deepEqual(argv, ['--mcp-config', config, '--verbose'])
})
