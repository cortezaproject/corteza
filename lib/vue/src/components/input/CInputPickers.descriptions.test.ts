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

// A description that does not fit is wrapped, never clipped: the whole point of
// the second line is the text on it, and an ellipsis hides exactly the part that
// made two options tell apart.
//
// Three parts hold that together and none works alone:
//  - the option wrapper is a shrink-to-fit flex item of the option row, so it
//    needs `w-full` to reach the row's width and `min-w-0` to be allowed back
//    below its content's width,
//  - the option row itself is `white-space: nowrap`, so the wrapper carries
//    `whitespace-normal` to let the text break and `break-words` to keep an
//    unbroken token — a URL in a description — inside the row; both properties
//    inherit, so the name on the first line wraps by the same rule,
//  - and the overlay is capped, which is what makes the text wrap at all: the
//    panel is sized by its content, so with nothing bounding it a long
//    description opens one very wide line instead of several readable ones.
//    PrimeVue already sets the overlay's `min-width` to the trigger's width
//    inline, and `min-width` wins over `max-width`, so the cap changes nothing
//    for a picker already wider than it.
const WIDTH_RULE = [...PICKERS.map(([file]) => file), 'CInputKnowledgeBase.vue']

const OVERLAY_CAP = `:pt="{ overlay: { class: 'max-w-lg' } }"`

describe('a picker description wraps to the width of its row', () => {
  it.each(WIDTH_RULE)('%s does not cap the description at a fixed width', file => {
    expect(read(file)).not.toContain('max-w-64')
  })

  it.each(WIDTH_RULE)('%s clips no text anywhere', file => {
    expect(read(file)).not.toContain('truncate')
  })

  it.each(WIDTH_RULE)('%s lets the option wrapper fill the row and wrap in it', file => {
    expect(read(file)).toContain(
      '<div class="flex flex-col w-full min-w-0 whitespace-normal break-words">',
    )
  })

  it.each(WIDTH_RULE)(
    '%s caps the overlay so a long description wraps rather than widens it',
    file => {
      const text = read(file)
      // One cap per dropdown, not one per file: Module and Role each render a
      // MultiSelect and a Select.
      const caps = text.split(OVERLAY_CAP).length - 1
      const dropdowns = text.split('<template #option=').length - 1
      expect(caps).toBe(dropdowns)
      expect(caps).toBeGreaterThan(0)
    },
  )
})
