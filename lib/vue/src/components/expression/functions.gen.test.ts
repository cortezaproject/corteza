import fs from 'fs'
import yaml from 'js-yaml'
import { join } from 'path'
import { describe, expect, it } from 'vitest'

import { functionNames, render } from '../../../tools/codegen/expr-functions.js'

// vitest runs from the package root (lib/vue)
const source = join(process.cwd(), '../../server/pkg/expr/expr_functions.yaml')
const generated = join(process.cwd(), 'src/components/expression/functions.gen.ts')

describe('functions.gen.ts', () => {
  it('matches the server definitions — run `pnpm codegen` in lib/vue when it does not', () => {
    const names = functionNames(yaml.load(fs.readFileSync(source, 'utf8')))

    expect(names.length).toBeGreaterThan(0)
    expect(fs.readFileSync(generated, 'utf8')).toBe(render(names))
  })
})
