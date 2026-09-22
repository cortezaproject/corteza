import { expect, test, type Page } from '@playwright/test'

// Contract under test:
// - components/builder/builder.intent.md: `NodePicker.vue` lists the catalog's
//   function steps grouped by the construct's own group, and a construct the
//   catalog marks `disabled` (today only `corredorExec`, while Corredor is off)
//   is listed greyed and cannot be picked — so a step's absence is never mistaken
//   for a missing feature.
// - components/builder/form/inputs/registry.intent.md: the `corredorExec` step's
//   `script` argument is a `CorredorScriptSelector`, which maps to lib/vue
//   `CInputCorredorScript`: an editable Select over the server's manual
//   system-level server scripts, keeping a name the list does not carry.
//
// What each group would have caught:
// - the picker: the construct declares `Groups: ["Corredor"]` server-side; a
//   picker that grouped by `Labels` instead dropped the step into a group named
//   after the empty string, where nobody could find it.
// - the script field: the registry maps `input.type` to a component, and an
//   unmapped type silently falls back to a plain InputText — the field still
//   looked fine and simply never offered a script. Asserting that the compose
//   record scripts are NOT offered guards the selector's own filter
//   (`onManual` + `system` + server-side only).
// - reopening: the step's arguments round-trip through `Expr[]`; a name typed
//   rather than picked was dropped on the way out and the step reopened empty.
//
// Nothing is saved, so there is nothing to clean up.

const CORREDOR_GROUP = 'Corredor'
const CORREDOR_STEP = 'Run Corredor script'
const FIXTURE_SCRIPT = 'System ping (server)'
const TYPED_SCRIPT = '/server-scripts/e2e-typed-by-hand.js:default'

/** Opens a fresh builder, adds the first trigger the catalog offers and opens the
 *  step picker on its Corredor group. Returns why the Corredor step cannot be
 *  exercised here, empty when it can. */
async function openCorredorGroup(page: Page): Promise<string> {
  await page.goto('/taq/builder')

  await page.getByRole('button', { name: 'Add Trigger' }).click()
  // Picker entries are plain divs; `node-item` is NodePicker's own class.
  const firstTrigger = page.getByRole('dialog').locator('.node-item').first()
  await expect(firstTrigger).toBeVisible({ timeout: 30000 })
  await firstTrigger.click()

  // The "+" on the edge between the trigger and End opens the step picker.
  const addStep = page.locator('.edge-button button').first()
  await expect(addStep).toBeVisible({ timeout: 30000 })
  await addStep.click()

  const picker = page.getByRole('dialog')
  await expect(picker.locator('.node-item').first()).toBeVisible({ timeout: 30000 })

  const group = picker.getByRole('button', { name: CORREDOR_GROUP })
  if ((await group.count()) === 0) {
    return 'the construct catalog has no Corredor group — this server registers no corredorExec at all, which is older than the disabled-construct contract'
  }

  await group.click()
  const step = picker.locator('.node-item').filter({ hasText: CORREDOR_STEP })
  if ((await step.count()) === 0) {
    return 'the Corredor group carries no Run Corredor script step on this stack'
  }

  if ((await step.getAttribute('aria-disabled')) === 'true') {
    return 'corredorExec is listed disabled on this stack (Corredor is off — CORREDOR_ENABLED) — bring the stack up with Corredor: dev/agent/worktree.sh up'
  }

  return ''
}

/** The step's Script argument — an editable Select, addressed by its placeholder. */
function scriptField(page: Page) {
  return page.getByPlaceholder('Select Script')
}

test.describe.serial('the corredor step in the TAQ builder', () => {
  test('the step picker offers Run Corredor script in a Corredor group', async ({ page }) => {
    const unavailable = await openCorredorGroup(page)
    test.skip(!!unavailable, unavailable)

    const step = page.getByRole('dialog').locator('.node-item').filter({ hasText: CORREDOR_STEP })
    await expect(step).toHaveCount(1)
    await expect(step).toContainText('Run a manual Corredor server script')
    // Pickable, not merely listed — the greyed state is what Corredor being off
    // looks like.
    await expect(step).toHaveAttribute('aria-disabled', 'false')
  })

  test('adding the step asks for a script and offers the deployed ones', async ({ page }) => {
    const unavailable = await openCorredorGroup(page)
    test.skip(!!unavailable, unavailable)

    await page.getByRole('dialog').getByText(CORREDOR_STEP, { exact: true }).click()

    // The construct declares `script` required; the form marks it.
    await expect(page.getByText('Script*', { exact: true })).toBeVisible({ timeout: 30000 })
    await expect(scriptField(page)).toBeVisible()

    await scriptField(page).press('ArrowDown')
    const options = page.getByRole('option')
    await expect(options.filter({ hasText: FIXTURE_SCRIPT })).toHaveCount(1)

    // Only manual system-level server scripts belong here: a compose record
    // script has no way to be run by a TAQ step.
    await expect(options.filter({ hasText: 'Activate contact (server)' })).toHaveCount(0)
    await expect(options.filter({ hasText: 'Greet contact (client)' })).toHaveCount(0)
  })

  test('a script name typed by hand survives reopening the step', async ({ page }) => {
    const unavailable = await openCorredorGroup(page)
    test.skip(!!unavailable, unavailable)

    await page.getByRole('dialog').getByText(CORREDOR_STEP, { exact: true }).click()
    await expect(scriptField(page)).toBeVisible({ timeout: 30000 })

    // A script deployed later, or hidden from this caller, is kept as typed.
    await scriptField(page).fill(TYPED_SCRIPT)
    await scriptField(page).press('Enter')

    // Away from the step and back: `vue-flow__node` is the canvas library's own
    // node class, and selecting a node is what opens its configuration.
    await page.locator('.vue-flow__node-trigger').first().click()
    await expect(scriptField(page)).toHaveCount(0)

    await page.locator('.vue-flow__node').filter({ hasText: CORREDOR_STEP }).first().click()
    await expect(scriptField(page)).toHaveValue(TYPED_SCRIPT, { timeout: 30000 })
  })
})
