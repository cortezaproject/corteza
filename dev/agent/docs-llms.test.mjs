// Run: node --test dev/agent/docs-llms.test.mjs   (builds the docs site, ~10 s)
//
// docs/llms.txt and docs/llms-full.txt are the committed copies of what the
// VitePress llms plugin emits, so an LLM can read the product docs from the
// repo without a build. They go stale with every docs edit; this rebuilds into
// a scratch directory and compares. `make docs-llms` refreshes them.
import { test } from 'node:test'
import assert from 'node:assert/strict'
import { spawnSync } from 'node:child_process'
import { mkdtempSync, readFileSync, rmSync } from 'node:fs'
import { tmpdir } from 'node:os'
import { fileURLToPath } from 'node:url'
import path from 'node:path'

const here = path.dirname(fileURLToPath(import.meta.url))
const docs = path.join(here, '..', '..', 'docs')

test('docs/llms*.txt match a fresh docs build', { timeout: 120_000 }, () => {
  const out = mkdtempSync(path.join(tmpdir(), 'llms-'))
  try {
    const r = spawnSync('npx', ['vitepress', 'build', '--outDir', out], { cwd: docs, encoding: 'utf8' })
    assert.equal(r.status, 0, r.stderr.slice(-2000))
    for (const name of ['llms.txt', 'llms-full.txt']) {
      const fresh = readFileSync(path.join(out, name), 'utf8')
      const committed = readFileSync(path.join(docs, name), 'utf8')
      assert.equal(committed, fresh, `docs/${name} is stale; run make docs-llms and commit it`)
    }
  } finally {
    rmSync(out, { recursive: true, force: true })
  }
})
