import { expect, test, type Locator, type Page } from '@playwright/test'

// Lifecycle + dashboard smoke, covering:
// - project.intent.md: locked lifecycle draft -> active via the wizard-header
//   request -> approve -> publish cycle; the live status is 'active', the
//   backend never sets 'published'.
// - sidebar/sidebar.intent.md: routing rule — a live (active) project opens
//   its dashboard (project.overview), anything else opens the wizard.
// - views/dashboard/*.intent.md: locked dashboard view set — Overview,
//   Category, Backlog, All Events, + a Reports stub — behind DashboardLayout's
//   left-rail nav.
//
// project.intent.md's WIP note: approval PERSISTENCE is session-local FE
// scaffolding (no backend governance store), so create -> request approval ->
// approve -> publish must happen in one continuous page session, never across
// a reload. Splitting that across separate Playwright test()s would give each
// a fresh browser context (and a fresh in-memory Pinia store), silently
// losing the request-approval state — so the create + publish cycle below is
// ONE test, not two, even though it covers two of the task's five concerns.
//
// A second wrinkle from config/roles.js: no role preset carries both
// requestApproval and grantApproval, and the project creator is added as a
// 'developer' (requestApproval only). A single e2e user can only walk the
// full cycle solo by self-promoting to a grantApproval role (e.g.
// 'governance-owner') via the Members dialog between the request and approve
// steps — which itself requires the members.manage RBAC permission
// (project.canManageMembers). If the logged-in e2e user doesn't hold it, the
// role <select> won't render and the test skips (tolerate, don't fail — see
// e2e.intent.md) rather than assuming a capability we can't verify from here.
test.describe.serial('project lifecycle & dashboard', () => {
  const projectName = `e2e-lifecycle-${Date.now()}`
  let projectId = ''
  let isLive = false

  async function findRow(page: Page, name: string): Promise<Locator> {
    // Search narrows server-side (handle OR meta name since the 2026-07-24
    // filter fix); the text filter still pins the exact row. The search is
    // debounced and rewrites the URL (?query=) when it lands — wait for that
    // rewrite BEFORE returning, or a row click races it and its navigation
    // is lost.
    await page.getByPlaceholder('Search projects').fill(name)
    await page.waitForURL(u => u.searchParams.get('query') === name)
    await page.waitForLoadState('networkidle')
    return page.locator('tbody tr').filter({ hasText: name })
  }

  async function openRowMenu(row: Locator): Promise<Locator> {
    await row.getByRole('button').click()
    // PrimeVue's TieredMenu popup exposes menuitems but no role=menu wrapper
    // — scope to the page and wait for the items instead.
    const menu = row.page().locator('body')
    await menu.getByRole('menuitem').first().waitFor()
    return menu
  }

  async function openMembersDialog(page: Page): Promise<Locator> {
    await page.getByRole('button', { name: 'Members' }).click()
    const dialog = page.getByRole('dialog', { name: 'Project members' })
    await expect(dialog).toBeVisible()
    return dialog
  }

  // Either the header's default close icon or the footer's explicit Close
  // button works — both just flip `visible` to false — so `.first()` avoids
  // having to disambiguate the two same-labelled buttons PrimeVue renders.
  async function closeDialog(dialog: Locator) {
    await dialog.getByRole('button', { name: 'Close' }).first().click()
    await expect(dialog).toBeHidden()
  }

  // Best-effort and idempotent: archive first (an active project can't be
  // deleted directly — see ProjectList.vue's row actions), then delete. Safe
  // to call more than once (e.g. from the afterAll safety net below) — a
  // missing row just means there's nothing left to clean up.
  async function cleanupProject(page: Page, name: string) {
    await page.goto('/project/projects')
    await page.waitForLoadState('networkidle')
    const row = await findRow(page, name)
    if ((await row.count()) === 0) return

    let menu = await openRowMenu(row)
    if ((await menu.getByRole('menuitem', { name: 'Unarchive' }).count()) === 0) {
      await menu.getByRole('menuitem', { name: 'Archive', exact: true }).click()
      await expect(row.getByText('Archived', { exact: true })).toBeVisible()
    } else {
      await page.keyboard.press('Escape')
    }

    menu = await openRowMenu(row)
    await menu.getByRole('menuitem', { name: 'Delete' }).click()
    await page.getByRole('alertdialog').getByRole('button', { name: 'Delete', exact: true }).click()
    await expect(row).toHaveCount(0)
  }

  test('creates a project, confirms draft routing, then drives request -> approve -> publish to Active in one session', async ({
    page,
  }) => {
    await page.goto('/project/projects')
    await page.waitForLoadState('networkidle')

    // --- Create: name/description only (deployer questions moved to Govern) ---
    await page.getByRole('button', { name: 'New Project' }).click()
    const createDialog = page.getByRole('dialog')
    await expect(createDialog).toBeVisible()
    await createDialog.getByPlaceholder('e.g. CRM & Sales').fill(projectName)
    await createDialog.getByRole('button', { name: 'Create project' }).click()

    await page.waitForURL(/\/project\/projects\/[^/]+\/wizard/)
    const match = new URL(page.url()).pathname.match(/\/project\/projects\/([^/]+)\/wizard/)
    projectId = match ? match[1] : ''
    expect(projectId).not.toBe('')

    // Auto-opened Members dialog (route.query.new === '1') — close it first.
    const autoMembers = page.getByRole('dialog', { name: 'Project members' })
    await expect(autoMembers).toBeVisible()
    await closeDialog(autoMembers)

    // --- Draft routing: /project/projects row -> wizard (sidebar.intent.md) --
    await page.goto('/project/projects')
    await page.waitForLoadState('networkidle')
    const row = await findRow(page, projectName)
    await expect(row).toBeVisible()
    await row.click()
    await page.waitForURL(/\/project\/projects\/[^/]+\/wizard/)

    // --- Publish cycle: request -> approve -> publish (project.intent.md) --
    const requestBtn = page.getByRole('button', { name: 'Request approval' })
    await expect(requestBtn).toBeVisible()
    await requestBtn.click()
    // The creator is only a 'developer' (requestApproval, not grantApproval —
    // config/roles.js has no single role with both), so the header's action
    // button disappears once submitted, until a grantApproval role is held.
    await expect(requestBtn).toBeHidden()

    const membersDialog = await openMembersDialog(page)
    const roleSelect = membersDialog.getByRole('combobox').first()
    if ((await roleSelect.count()) === 0) {
      // No members.manage permission on this project for the logged-in e2e
      // user — can't self-elevate to a grantApproval role, so this single
      // session can't finish the cycle. Tolerate, don't fail.
      await closeDialog(membersDialog)
      test.skip(
        true,
        'logged-in e2e user cannot manage project members (project.canManageMembers false) — cannot self-elevate to a grant-approval role to complete the publish cycle solo',
      )
    }
    await roleSelect.click()
    await page.getByRole('option', { name: 'Governance Owner' }).click()
    await closeDialog(membersDialog)

    const approveBtn = page.getByRole('button', { name: 'Approve project' })
    await expect(approveBtn).toBeVisible()
    await approveBtn.click()

    const publishBtn = page.getByRole('button', { name: 'Publish project' })
    await expect(publishBtn).toBeVisible()
    await publishBtn.click()

    const confirmDialog = page.getByRole('alertdialog')
    await expect(confirmDialog).toBeVisible()
    await confirmDialog.getByRole('button', { name: 'Publish project' }).click()

    // doPublish() hands off to the dashboard on success — the definitive
    // signal the project is now live (never assert a 'Published' status).
    await page.waitForURL(new RegExp(`/project/projects/${projectId}$`))
    isLive = true
  })

  test('routing flip: a live project row opens the dashboard, not the wizard', async ({ page }) => {
    test.skip(!isLive, 'publish cycle did not complete in the previous test')

    await page.goto('/project/projects')
    await page.waitForLoadState('networkidle')
    const row = await findRow(page, projectName)
    await expect(row).toBeVisible()
    await expect(row.getByText('Active', { exact: true })).toBeVisible()

    await row.click()
    await page.waitForURL(new RegExp(`/project/projects/${projectId}$`))
    expect(page.url()).not.toMatch(/\/wizard/)
  })

  test('dashboard view set renders: Overview, All Events, Backlog (+ nav entries for Category/Reports)', async ({
    page,
  }) => {
    test.skip(!isLive, 'publish cycle did not complete in the first test')

    await page.goto(`/project/projects/${projectId}`)
    await page.waitForLoadState('networkidle')

    // DashboardNav's left rail — scope past the app's own section sidebar
    // (also a nav landmark) via its unique 'Monitor' section heading.
    const nav = page.getByRole('navigation').filter({ hasText: 'Monitor' })
    await expect(nav).toBeVisible()
    for (const name of ['Dashboard', 'All Events', 'Backlog', 'Reports', 'Incidents']) {
      await expect(nav.getByRole('button', { name })).toBeVisible()
    }

    // Overview — already landed here on the initial navigation.
    await expect(page.getByRole('heading', { name: 'Dashboard' })).toBeVisible()

    await nav.getByRole('button', { name: 'All Events' }).click()
    await page.waitForURL(/\/events$/)
    await expect(page.getByRole('heading', { name: 'All Events' })).toBeVisible()

    await nav.getByRole('button', { name: 'Backlog' }).click()
    await page.waitForURL(/\/backlog$/)
    await expect(page.getByRole('heading', { name: 'Backlog' })).toBeVisible()

    // Category — directly clickable via the 'Incidents' rail entry, so this
    // isn't skipped (per the task: skip only if no category entry is
    // directly clickable).
    await nav.getByRole('button', { name: 'Incidents' }).click()
    await page.waitForURL(/\/category\/incident$/)
    await expect(page.getByRole('heading', { name: 'Incidents' })).toBeVisible()
  })

  test('teardown: archive and delete the created project', async ({ page }) => {
    await cleanupProject(page, projectName)
  })

  // Safety net: describe.serial skips remaining tests (including the
  // teardown test above) after any hard failure, but afterAll hooks always
  // run regardless — so this is what actually guarantees cleanup when an
  // earlier assertion throws. cleanupProject() is idempotent (a missing row
  // is a no-op), so re-running it here after a successful teardown test costs
  // one harmless list fetch.
  test.afterAll(async ({ browser }) => {
    const page = await browser.newPage({ storageState: 'e2e/.auth/state.json' })
    try {
      await cleanupProject(page, projectName)
    } catch (err) {
      console.warn('lifecycle-dashboard.spec.ts: teardown safety net failed:', err)
    } finally {
      await page.close()
    }
  })
})
