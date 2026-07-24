#!/usr/bin/env node
// Claude Code PostToolUse hook (Edit|Write): when a covered+enforced file is edited,
// inject the governing intent doc into context so it gets reconciled.
import { execSync } from 'node:child_process'
import { relative, resolve, dirname } from 'node:path'
import { fileURLToPath } from 'node:url'

const ROOT = resolve(dirname(fileURLToPath(import.meta.url)), '../..')

let input = ''
for await (const chunk of process.stdin) input += chunk
let filePath
try {
  filePath = JSON.parse(input)?.tool_input?.file_path
} catch {
  process.exit(0)
}
if (!filePath) process.exit(0)
const rel = relative(ROOT, resolve(filePath)).replaceAll('\\', '/')
if (rel.startsWith('..') || rel.endsWith('.intent.md') || rel.endsWith('INTENT.md')) process.exit(0)

// Only speak up for files that are enforced and drifting; stay silent otherwise.
try {
  execSync(`node .intent/intent.mjs check --changed`, { cwd: ROOT, encoding: 'utf8', stdio: ['ignore', 'pipe', 'pipe'] })
  process.exit(0)
} catch (e) {
  const report = (e.stderr ?? '').toString()
  // Match only actual drift lines ("  <file>  (why)  → …"), not advice text
  const line = report.split('\n').find((l) => l.trim().startsWith(rel + '  ') && l.includes('→'))
  if (!line) process.exit(0)
  const doc = line.split('→').pop()?.replace(/update & sync:/, '').trim()
  console.log(
    JSON.stringify({
      hookSpecificOutput: {
        hookEventName: 'PostToolUse',
        additionalContext: `INTENT: you edited ${rel}, governed by ${doc ?? 'an intent doc'}. Before ending the turn: reconcile that doc with the change (update it, or confirm intent is unchanged), then run: node .intent/intent.mjs sync ${rel}`,
      },
    }),
  )
}
