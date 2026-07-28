import { expect, test } from '@playwright/test'

// Contract: sections/project/views/ProjectList.intent.md (list at
// /project/projects: list/create/open/archive/delete) and
// sections/project/components/project/project.intent.md (NewProjectDialog's
// locked create flow — name/description only; the AI Act Deployer category
// questions moved out per the 2026-07-28 ruling and now live in the Govern
// tab's FRIA flow instead — see FriaDeterminationStep.vue).

// Tracks the project created by the "create flow" test below so the shared
// afterEach can clean it up from the dev stack — set as soon as the name is
// generated (before any dialog interaction) so cleanup is attempted even if
// an assertion fails partway through the flow.
let createdName: string | null = null

test.afterEach(async ({ page }) => {
  if (!createdName) return
  const name = createdName
  createdName = null

  await page.goto('/project/projects')
  await page.waitForLoadState('networkidle')

  const row = page.locator('tbody tr').filter({ hasText: name })
  if ((await row.count()) === 0) return // never created, or already cleaned up

  await row.first().hover()
  await row.first().locator('button.row-action-btn').click()

  const deleteItem = page.getByRole('menuitem', { name: 'Delete' })
  if ((await deleteItem.count()) === 0) return
  await deleteItem.click()

  const confirmAccept = page.getByRole('alertdialog').getByRole('button', { name: 'Delete' })
  if (await confirmAccept.count()) {
    await confirmAccept.click()
    await page.waitForLoadState('networkidle')
  }
})

test('project list renders at /project/projects', async ({ page }) => {
  await page.goto('/project/projects')
  await page.waitForLoadState('networkidle')

  // Tolerate existing dev data: either populated rows or the empty state.
  const rows = page.locator('tbody tr')
  const empty = page.getByText('No projects yet.')
  await expect(rows.first().or(empty)).toBeVisible()
})

test('create flow: name + description lands in the wizard', async ({ page }) => {
  await page.goto('/project/projects')
  await page.waitForLoadState('networkidle')

  createdName = `e2e-list-${Date.now()}`

  await page.getByRole('button', { name: 'New Project' }).click()

  // Single step — name + description only (locked 2026-07-28: the AI Act
  // Deployer category questions moved out to the Govern tab's FRIA flow).
  const dialog = page.getByRole('dialog')
  await expect(dialog).toBeVisible()
  await expect(dialog.getByText('Project details')).toBeVisible()

  const nameInput = dialog.getByPlaceholder('e.g. CRM & Sales')
  const descriptionInput = dialog.getByPlaceholder('Optional short description')
  await expect(nameInput).toBeVisible()
  await expect(descriptionInput).toBeVisible()
  await nameInput.fill(createdName)

  await dialog.getByRole('button', { name: 'Create project' }).click()

  // Just-created handoff: lands in the wizard with ?new=1 (ProjectList.vue /
  // Wizard.vue contract).
  await page.waitForURL(/\/wizard/)
  await page.waitForLoadState('networkidle')

  // The members dialog auto-opens once for a just-created project (Wizard.vue
  // route.query.new watcher) — covered by another spec, not asserted here.
  // Just close it if it appeared so the page is left in a clean state.
  const membersDialog = page.getByRole('dialog', { name: 'Project members' })
  if (await membersDialog.isVisible().catch(() => false)) {
    await page.keyboard.press('Escape')
  }
})
