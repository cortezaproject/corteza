// @vitest-environment node
import fs from 'fs'
import yaml from 'js-yaml'
import { fileURLToPath } from 'url'
import { describe, expect, it } from 'vitest'

import { functionNames, render } from '../../../tools/codegen/expr-functions.js'

const source = fileURLToPath(
  new URL('../../../../../server/pkg/expr/expr_functions.yaml', import.meta.url),
)
const generated = fileURLToPath(new URL('./functions.gen.ts', import.meta.url))

describe('functions.gen.ts', () => {
  it('matches the server definitions — run `pnpm codegen` in lib/vue when it does not', () => {
    const names = functionNames(yaml.load(fs.readFileSync(source, 'utf8')))

    expect(names.length).toBeGreaterThan(0)
    expect(fs.readFileSync(generated, 'utf8')).toBe(render(names))
  })
})
