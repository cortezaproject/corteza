#!/usr/bin/env node
//
// Emit the colour scheme names the server embeds for compose_chart_create's
// colorScheme check, at server/compose/agentic/chart_color_schemes.json.
//
// These tables are the source of truth; the server only holds a copy so it can
// refuse a scheme the webapp cannot resolve. contract.test.ts asserts the copy
// matches what actually imports, so a family added here without a rerun of this
// script fails the suite rather than shipping a stale list.
//
//   node lib/js/src/shared/types/chart/colorschemes/emit-contract.mjs

import { readFileSync, writeFileSync } from 'node:fs'
import { dirname, resolve } from 'node:path'
import { fileURLToPath } from 'node:url'

const HERE = dirname(fileURLToPath(import.meta.url))
const FAMILIES = ['brewer', 'office', 'tableau']
const OUT = resolve(HERE, '../../../../../../../server/compose/agentic/chart_color_schemes.json')

// Each family file is `export default { Name: [...colours], ... }`; the keys are
// at one indent level, which is what separates them from the colour literals.
const names = FAMILIES.flatMap(family => {
  const src = readFileSync(resolve(HERE, `${family}.ts`), 'utf8')
  return [...src.matchAll(/^ {2}([A-Za-z0-9_]+):/gm)].map(m => `${family}.${m[1]}`)
}).sort()

if (names.length < 400) {
  console.error(`refusing to write ${names.length} names — the parse looks wrong`)
  process.exit(1)
}

writeFileSync(OUT, `${JSON.stringify(names, null, 2)}\n`)
console.log(`${names.length} scheme names -> ${OUT}`)
