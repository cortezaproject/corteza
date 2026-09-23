import { describe, it } from 'mocha'
import { expect } from 'chai'
import { execFileSync } from 'child_process'
import path from 'path'

// The scripting reference is generated from the files that define the scripting
// surface, so a change to any of them that the reference does not carry is a
// failure here rather than a doc nobody noticed had gone stale.
const tool = path.join(__dirname, '..', 'tools', 'scripting-reference.mjs')

describe('scripting reference', () => {
  it('matches what its sources say', () => {
    try {
      execFileSync(process.execPath, [tool, '--check'], { encoding: 'utf8', stdio: 'pipe' })
    } catch (e) {
      const { stdout, stderr } = e as { stdout?: string; stderr?: string }
      expect.fail(`${stdout ?? ''}${stderr ?? ''}`.trim() || 'the reference check failed')
    }
  })
})
