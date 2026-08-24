import { describe, it, expect } from 'vitest'
import { readFileSync } from 'node:fs'
import { fileURLToPath } from 'node:url'
import { dirname, join } from 'node:path'

// A picker names a resource and nothing else unless it is told to. Every
// resource below carries a description the server already stores, so the
// dropdown shows it as a second line — picking a role or a module by name
// alone is a guess when two of them are named alike.
//
// The list is explicit: it is what stops the second line being deleted as
// decoration, not a discovery of pickers that ought to have one.

const HERE = dirname(fileURLToPath(import.meta.url))

// picker file → [what its option template renders, what guards it]
const PICKERS: Array<[string, string, string]> = [
  ['CInputAgent.vue', 'option.meta.description', 'option.meta?.description'],
  ['CInputChart.vue', 'option.meta.description', 'option.meta?.description'],
  ['CInputLLM.vue', 'option.meta.description', 'option.meta?.description'],
  ['CInputModule.vue', 'option.meta.description', 'option.meta?.description'],
  ['CInputNamespace.vue', 'option.meta.description', 'option.meta?.description'],
  ['CInputRole.vue', 'option.meta.description', 'option.meta?.description'],
  ['CInputUserGroup.vue', 'option.meta.description', 'option.meta?.description'],
  ['CInputWorkflow.vue', 'option.meta.description', 'option.meta?.description'],
  // NgAutomation keeps its description as a top-level field, not under meta.
  ['CInputTAQ.vue', 'option.description', 'option.description'],
]

const read = (file: string) => readFileSync(join(HERE, file), 'utf8')

describe('resource pickers show the resource description', () => {
  it.each(PICKERS)('%s renders %s in its option template', (file, expr) => {
    const text = read(file)
    expect(text).toContain('#option')
    expect(text).toContain(`{{ ${expr} }}`)
  })

  // Rendering it unguarded prints an empty muted line under every option that
  // has no description, which is worse than not showing one at all.
  it.each(PICKERS)('%s hides the line when the description is empty', (file, _expr, guard) => {
    expect(read(file)).toContain(`v-if="${guard}"`)
  })

  it('covers the pickers that pick a described resource', () => {
    expect(PICKERS.length).toBeGreaterThanOrEqual(9)
  })
})
