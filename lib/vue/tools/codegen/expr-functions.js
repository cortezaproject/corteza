// Writes src/components/expression/functions.gen.ts from the server's
// expression function definitions (server/pkg/expr/expr_functions.yaml).
import fs from 'fs'
import yaml from 'js-yaml'
import { fileURLToPath } from 'url'
import { dirname, join } from 'path'

const __dirname = dirname(fileURLToPath(import.meta.url))

const server = process.argv[2] || join(__dirname, '../../../../server')
const src = join(server, 'pkg/expr/expr_functions.yaml')
const dst = join(__dirname, '../../src/components/expression/functions.gen.ts')

export function functionNames(def) {
  return def.groups.flatMap(g => g.functions.map(f => f.name)).sort()
}

export function render(names) {
  return [
    '// This file is auto-generated from server/pkg/expr/expr_functions.yaml.',
    '// Run `pnpm codegen` in lib/vue to regenerate it.',
    '',
    '// Every function the expression language registers, sorted.',
    'export const EXPR_FUNCTIONS: readonly string[] = [',
    ...names.map(n => `  '${n}',`),
    ']',
    '',
  ].join('\n')
}

if (process.argv[1] === fileURLToPath(import.meta.url)) {
  const def = yaml.load(fs.readFileSync(src, 'utf8'))
  fs.writeFileSync(dst, render(functionNames(def)))
  console.log(`wrote ${dst}`)
}
