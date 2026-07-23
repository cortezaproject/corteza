#!/usr/bin/env node
// Claude Code Stop hook: block ending the turn while intent drift exists.
import { execSync } from 'node:child_process'
import { resolve, dirname } from 'node:path'
import { fileURLToPath } from 'node:url'

const ROOT = resolve(dirname(fileURLToPath(import.meta.url)), '../..')

let input = ''
for await (const chunk of process.stdin) input += chunk
let stopHookActive = false
try {
  stopHookActive = JSON.parse(input)?.stop_hook_active === true
} catch {}

try {
  execSync('node .intent/intent.mjs check --changed', { cwd: ROOT, stdio: ['ignore', 'pipe', 'pipe'] })
  process.exit(0)
} catch (e) {
  const report = (e.stderr ?? '').toString().trim()
  if (stopHookActive) {
    // Second pass still failing — warn but do not loop forever.
    console.log(`intent drift remains (not blocking again):\n${report}`)
    process.exit(0)
  }
  console.error(`Intent drift — reconcile the listed intent docs and run "node .intent/intent.mjs sync" before finishing:\n${report}`)
  process.exit(2)
}
