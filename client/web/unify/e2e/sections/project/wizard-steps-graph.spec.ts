import { expect, test } from '@playwright/test'

// Contract under test:
// - sections/project/project.intent.md (locked): the Build step set is
//   exactly data model, connections, automations, agents, chatbots, pages,
//   roles, permissions, users, in that order (source of truth:
//   config/pipeline.js STEPS) — and every resource kind has a CreateDialog +
//   DetailDialog (the "resource dialog standard").
// - components/graph/graph.intent.md (locked): the resource graph is the
//   Build tab's primary surface, and clicking any graph node opens that
//   resource's EDITABLE detail dialog, never a read-only one.
// wizard.spec.ts already covers the three-tab shape (Build / Govern / Manage
// & Monitor) — not repeated here.
//
// Unlike the other e2e specs, this one needs a project it can create
// resources on and inspect the graph for, so it creates and tears down its
// own project instead of relying on whatever dev data happens to exist.
// Setup/teardown run as beforeAll/afterAll (not as ordinary tests): inside a
// `describe.serial` block a failing test skips the rest of the block, which
// would skip a "teardown" test too — hooks always run, so they're what makes
// cleanup robust to an earlier assertion failing.

const BASE_URL = process.env.E2E_BASE_URL || 'http://localhost:5173'
const STORAGE_STATE = 'e2e/.auth/state.json'

// Pipeline order — config/pipeline.js STEPS filtered to tab: 'build'; labels
// are locale/en/human-webapp/project.yaml's `steps.<key>.label`.
const BUILD_STEP_LABELS = [
  'Data Model',
  'Connections',
  'Automations',
  'Agents',
  'Chatbots',
  'Pages',
  'Roles',
  'Permissions',
  'Users',
]

test.describe.serial('wizard Build pipeline + resource graph contract', () => {
  const projectName = `e2e-wizard-graph-${Date.now()}`
  let wizardPath = ''

  test.beforeAll(async ({ browser }) => {
    const context = await browser.newContext({ baseURL: BASE_URL, storageState: STORAGE_STATE })
    try {
      const page = await context.newPage()

      // Create flow (project.intent.md, locked 2026-07-28: name + description
      // only — the AI Act Deployer/FRIA questions moved to the Govern tab).
      await page.goto('/project/projects')
      await page.getByRole('button', { name: 'New Project' }).click()

      const detailsDialog = page.getByRole('dialog', { name: 'Project details' })
      await expect(detailsDialog).toBeVisible()
      await detailsDialog.getByRole('textbox').first().fill(projectName)
      await detailsDialog.getByRole('button', { name: 'Create project' }).click()

      await page.waitForURL(/\/wizard/)
      wizardPath = new URL(page.url()).pathname

      // A just-created project auto-opens the Members dialog (?new=1) —
      // close it so it doesn't shadow the Build tab in the tests below.
      try {
        await page
          .getByRole('dialog', { name: 'Project members' })
          .waitFor({ state: 'visible', timeout: 5000 })
        await page.keyboard.press('Escape')
      } catch {
        // Didn't auto-open this run — nothing to close.
      }
    } finally {
      await context.close()
    }
  })

  test('Build tab lists the nine locked steps in pipeline order', async ({ page }) => {
    await page.goto(wizardPath)
    // The wizard normalizes its query via router.replace on entry — wait it
    // out, then assert with the retrying matcher (a one-shot allTextContents
    // read raced that navigation and died with "execution context destroyed").
    await page.waitForLoadState('networkidle')
    const stepNav = page.getByRole('navigation').filter({ hasText: 'Data Model' })
    await expect(stepNav.getByRole('button')).toHaveText(BUILD_STEP_LABELS)
  })

  test("resource graph renders as the Build tab's primary surface", async ({ page }) => {
    await page.goto(wizardPath)
    await expect(page.getByRole('heading', { name: 'Resources', level: 3 })).toBeVisible()
    await expect(page.getByTitle('Reload from the backend')).toBeVisible()
    // A fresh project has nothing built yet, so the pane shows its own empty
    // state rather than a chart — still the contract surface, just no nodes.
    await expect(page.getByText('Resources appear here as you create them')).toBeVisible()
  })

  test('creating a module surfaces it in the graph and opens an editable detail dialog', async ({
    page,
  }) => {
    await page.goto(wizardPath)

    const modulesChip = page.getByRole('button').filter({ hasText: 'Modules' })
    const countBefore = Number((await modulesChip.locator('span').nth(1).innerText()).trim())

    // Resource dialog standard: every kind's CreateDialog opens with a name
    // input (module.intent.md-level; ModuleCreateDialog.vue).
    await page.getByRole('button', { name: 'New module' }).click()
    const createDialog = page.getByRole('dialog')
    await expect(createDialog).toBeVisible()
    const nameInput = createDialog.getByRole('textbox').first()
    await expect(nameInput).toBeVisible()

    const moduleName = `e2e-module-${Date.now()}`
    await nameInput.fill(moduleName)
    await createDialog.getByRole('button', { name: 'Create module' }).click()
    await expect(createDialog).toBeHidden()

    // The module joins the resource graph: the Modules metric chip (counted
    // straight off the backend graph payload) increments, and the graph's
    // echarts canvas — inert while the graph was empty — now renders a node.
    await expect(modulesChip.locator('span').nth(1)).toHaveText(String(countBefore + 1))
    await expect(page.locator('canvas').first()).toBeVisible()

    // graph.intent.md (locked): clicking any graph node opens the EDITABLE
    // detail dialog. The graph renders to a <canvas> via echarts'
    // CanvasRenderer — there are no per-node DOM handles, and node position
    // comes from a force-directed layout with no stable coordinates to target
    // — so a graph node can't be clicked reliably from this suite. Opening
    // the same module through the data-model step's list row instead
    // exercises the identical inspectResource(kind, id) path a graph node
    // click uses (see Wizard.vue's provide/inject wiring), just from a
    // DOM-addressable surface.
    await page.getByText(moduleName, { exact: true }).click()

    const detailDialog = page.getByRole('dialog')
    await expect(detailDialog).toBeVisible()
    const detailNameInput = detailDialog.getByRole('textbox').first()
    await expect(detailNameInput).toBeVisible()
    await expect(detailNameInput).toBeEnabled()
    await expect(detailNameInput).toHaveValue(moduleName)
  })

  test.afterAll(async ({ browser }) => {
    const context = await browser.newContext({ baseURL: BASE_URL, storageState: STORAGE_STATE })
    try {
      const page = await context.newPage()

      await page.goto('/project/projects')

      // Search narrows server-side (handle OR meta name since the 2026-07-24
      // filter fix); the text filter still pins the exact row.
      await page.getByPlaceholder('Search projects…').fill(projectName)
      const row = page.locator('tbody tr').filter({ hasText: projectName })
      try {
        await expect(row).toHaveCount(1, { timeout: 15000 })
      } catch {
        return // Nothing to clean up (e.g. beforeAll never got this far).
      }

      // The row's action-menu trigger is an icon-only PrimeVue Button with no
      // accessible name (the same gap wizard.spec.ts routes around for tab
      // icons) — scope a CSS locator to the matched row rather than the page.
      await row.locator('.row-action-btn').click()
      await page.getByRole('menuitem', { name: 'Archive' }).click()

      await row.locator('.row-action-btn').click()
      await page.getByRole('menuitem', { name: 'Delete' }).click()

      const confirmDialog = page.getByRole('alertdialog', { name: 'Delete project' })
      await expect(confirmDialog).toBeVisible()
      await confirmDialog.getByRole('button', { name: 'Delete', exact: true }).click()

      await expect(row).toHaveCount(0, { timeout: 15000 })
    } finally {
      await context.close()
    }
  })
})
