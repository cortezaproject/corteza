// Where this checkout's dev stack listens, for the node side of the toolkit.
//
// One resolver, shared with the bash scripts: stack.sh reads server/.env and
// .env.e2e, so a browser check run inside a worktree drives that worktree's
// webapp rather than the primary's. An exported HUMAN_API / HUMAN_WEBAPP still
// wins — stack.sh inherits this process's environment.
import { execFileSync } from 'node:child_process'
import { dirname, join } from 'node:path'
import { fileURLToPath } from 'node:url'

const HERE = dirname(fileURLToPath(import.meta.url))

export function stack() {
  const out = execFileSync('bash', [join(HERE, 'stack.sh')], { encoding: 'utf8' })
  return Object.fromEntries(
    out
      .split('\n')
      .filter(Boolean)
      .map(line => {
        const at = line.indexOf('=')
        return [line.slice(0, at), line.slice(at + 1)]
      }),
  )
}
