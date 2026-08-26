import { describe, it, expect } from 'vitest'
import { readFileSync } from 'node:fs'
import { globSync } from 'node:fs'
import { join, dirname, relative } from 'node:path'
import { fileURLToPath } from 'node:url'

// A list's name cell carries the resource name over a muted subtitle. The
// subtitle sets `truncate`, which is `nowrap` — so without an upper bound the
// cell grows to the full run of text and pushes every later column off the
// right-hand edge. `max-w-full` is not a bound: it resolves against a cell that
// is itself sized by this text.

const here = dirname(fileURLToPath(import.meta.url))
const files = globSync('**/*.vue', { cwd: here })
  .map(f => join(here, f))
  .filter(f => readFileSync(f, 'utf8').includes('#body-name'))

const CLASS_ATTR = /class="([^"]*)"/g

function uncappedSubtitles(source) {
  const bad = []
  for (const [, classes] of source.matchAll(CLASS_ATTR)) {
    if (!classes.includes('text-muted-color') || !classes.includes('truncate')) continue
    const capped = classes.split(/\s+/).some(c => c.startsWith('max-w-') && c !== 'max-w-full')
    if (!capped) bad.push(classes)
  }
  return bad
}

describe('list name-cell subtitles', () => {
  it('covers every list view that renders a name cell', () => {
    expect(files.length).toBeGreaterThan(15)
  })

  it.each(files.map(f => [relative(here, f), f]))('%s bounds its subtitle', (_name, file) => {
    expect(uncappedSubtitles(readFileSync(file, 'utf8'))).toEqual([])
  })
})
