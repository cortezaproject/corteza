import { expect, test } from '@playwright/test'

// Locked members contract (sections/project/project.intent.md +
// components/project/project.intent.md): users hold per-project roles, and
// MembersDialog.vue — auto-opened once via the `?new=1` deep link right
// after create (Wizard.vue) — is the only management surface. Groups are
// intended but not built yet; not tested here.
//
// Serial: every test after the first depends on the project the first test
// creates (module-scoped `projectId`); `test.skip` guards each dependent
// test in case an earlier one didn't run/failed. Cleanup lives in
// `afterAll`, which Playwright always runs — including after an earlier
// assertion failure — so the created project doesn't linger in dev data.
test.describe.serial('project members dialog', () => {
  const projectName = `e2e-members-${Date.now()}`
  let projectId: string | undefined

  test('creating a project auto-opens the members dialog via ?new=1, then strips the query', async ({
    page,
  }) => {
    await page.goto('/project/projects')
    await page.waitForLoadState('networkidle')

    await page.getByRole('button', { name: 'New Project' }).click()

    // Step 1 — name/description (ProjectList/NewProjectDialog two-step create).
    const createDialog = page.getByRole('dialog')
    await createDialog.getByPlaceholder('e.g. CRM & Sales').fill(projectName)
    await createDialog.getByRole('button', { name: 'Next' }).click()

    // Step 2 — deployer categories. Every question already defaults to "No"
    // (NewProjectDialog.vue's defaultDeployer()), so creating needs no
    // further input; this only proves the step is traversable, not its
    // content (per the deployer/FRIA WIP note in project.intent.md).
    await createDialog.getByRole('button', { name: 'Create project' }).click()

    await page.waitForURL(/\/project\/projects\/[^/]+\/wizard/)
    projectId = page.url().match(/\/project\/projects\/([^/?]+)\/wizard/)?.[1]
    expect(projectId, 'wizard URL should carry the new project id').toBeTruthy()

    // Locked deep link: Wizard.vue's `route.query.new === '1'` watcher opens
    // the members dialog once, then strips the flag via router.replace.
    const membersDialog = page.getByRole('dialog', { name: 'Project members' })
    await expect(membersDialog).toBeVisible()
    await expect(page).not.toHaveURL(/[?&]new=1/)

    await page.keyboard.press('Escape')
    await expect(membersDialog).toBeHidden()
  })

  test('the auto-open is once-only — reloading the wizard does not reopen it', async ({ page }) => {
    test.skip(!projectId, 'no project created by the previous test')

    // A fresh navigation to the (now query-less) wizard URL is exactly what
    // Wizard.vue's own comment says a refresh/back-navigation lands on after
    // the strip — the flag is gone from the URL, so the watcher never fires.
    await page.goto(`/project/projects/${projectId}/wizard`)
    await page.waitForLoadState('networkidle')

    await expect(page).not.toHaveURL(/[?&]new=1/)
    await expect(page.getByRole('dialog', { name: 'Project members' })).toBeHidden()
  })

  test('reopening the members dialog manually shows the creator as a member with role/capability columns', async ({
    page,
  }) => {
    test.skip(!projectId, 'no project created by the first test')

    await page.goto(`/project/projects/${projectId}/wizard`)
    await page.waitForLoadState('networkidle')

    // Manual entry point: the wizard-header "Members" button (Wizard.intent.md).
    await page.getByRole('button', { name: 'Members' }).click()
    const dialog = page.getByRole('dialog', { name: 'Project members' })
    await expect(dialog).toBeVisible()

    // Capability-column shape (MembersDialog.vue): member, role preset, then
    // the four derived capability flags, then the descriptive resources cell.
    for (const header of [
      'Member',
      'Role',
      'Read',
      'Write',
      'Request approval',
      'Grant approval',
      'Resources',
    ]) {
      await expect(dialog.getByRole('columnheader', { name: header, exact: true })).toBeVisible()
    }

    // A freshly created project's only member is its creator — this test
    // session's user — per stores/projects.js create(): "adds the creator as
    // a developer".
    const rows = dialog.locator('tbody tr')
    await expect(rows).toHaveCount(1)

    // Assert the role resolves to one of the defined presets (config/roles.js)
    // rather than pinning the exact default beyond what the backend documents.
    const presetLabels = [
      'Developer',
      'Governance Owner',
      'Security Owner',
      'Junior Developer',
      'Executive Authority',
      'Infrastructure Administrator',
    ]
    await expect(rows.first().getByText(new RegExp(`^(${presetLabels.join('|')})$`))).toBeVisible()
  })

  test.afterAll(async ({ browser }) => {
    if (!projectId) return

    const context = await browser.newContext({ storageState: 'e2e/.auth/state.json' })
    const page = await context.newPage()
    try {
      await page.goto('/project/projects')
      await page.waitForLoadState('networkidle')

      // Search narrows server-side (handle OR meta name since the 2026-07-24
      // filter fix); the text filter still pins the exact row.
      await page.getByPlaceholder('Search projects').fill(projectName)
      const row = page.locator('tbody tr').filter({ hasText: projectName })
      await expect(row).toHaveCount(1)

      // Archive first (defensive, in case delete requires a non-draft
      // lifecycle status), then delete — mirrors ProjectList.vue's per-row
      // actions menu (rename / archive-unarchive / delete).
      await row.getByRole('button').click()
      await page.getByRole('menuitem', { name: 'Archive' }).click()
      await expect(row.getByText('Archived', { exact: true })).toBeVisible()

      await row.getByRole('button').click()
      await page.getByRole('menuitem', { name: 'Delete' }).click()
      await page
        .getByRole('alertdialog')
        .getByRole('button', { name: 'Delete', exact: true })
        .click()
      await expect(row).toHaveCount(0)
    } finally {
      await context.close()
    }
  })
})
