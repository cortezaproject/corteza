import { describe, it, expect } from 'vitest'
import { readFileSync, globSync } from 'node:fs'
import { join, dirname, relative } from 'node:path'
import { fileURLToPath } from 'node:url'
import yaml from 'js-yaml'

// A list footer reads "One {single}" for a single row and "{n} {plural}"
// otherwise. Both nouns come from one place — general.label.<noun> — so a list
// cannot end up naming itself with a create-button label ("One New Role") or
// with its plural twice ("One Agents").

const here = dirname(fileURLToPath(import.meta.url))
const repo = join(here, '..', '..', '..', '..', '..')
const labels = yaml.load(
  readFileSync(join(repo, 'locale/en/human-webapp/general.yaml'), 'utf8'),
).label

const PAIR = /resourceSingle:\s*\$?t\('([^']+)'\),\s*\n\s*resourcePlural:\s*\$?t\('([^']+)'\)/g

const sites = globSync('**/*.vue', { cwd: here }).flatMap(f => {
  const source = readFileSync(join(here, f), 'utf8')
  return [...source.matchAll(PAIR)].map(([, single, plural], i) => [`${f}#${i}`, single, plural])
})

describe('list footer nouns', () => {
  it('covers every list that names its resource', () => {
    expect(sites.length).toBeGreaterThan(25)
  })

  it.each(sites)('%s reads both nouns from one general.label entry', (_at, single, plural) => {
    const noun = single.match(/^general\.label\.(.+)\.single$/)?.[1]
    expect(noun, `singular key ${single}`).toBeTruthy()
    expect(plural).toBe(`general.label.${noun}.plural`)
  })

  it.each(sites)('%s names a resource whose two forms differ', (_at, single) => {
    const entry = labels[single.match(/^general\.label\.(.+)\.single$/)?.[1]]
    expect(entry?.single).toBeTruthy()
    expect(entry?.plural).toBeTruthy()
    expect(entry.single).not.toBe(entry.plural)
  })
})
