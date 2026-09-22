import { expect, test, type Locator, type Page } from '@playwright/test'

// Contract under test:
// - views/automation/Script/Index.intent.md: the read-only Corredor inventory —
//   its banner (off / unreachable / connected with the last refresh + a Refresh
//   that re-fetches), rows grouped by extension root, a kind badge, one chip per
//   trigger with its constraints, the iterator chip, the security chips, and the
//   Server/Client toggle pair that narrows the list.
//
// The data under test is the fixture extension in `dev/fixtures/corredor`, whose
// scripts are shaped so that every branch of the view has a row: a client script
// in a bundle, a server script, a before-event script with constraints, an
// interval script with `runAs`, and an iterator.
//
// What each group would have caught:
// - banner: a Refresh that replaced the rows with the raw axios envelope (the
//   endpoint answers `{set, enabled, connected}`, not an array) — the list came
//   back empty while the banner still said Connected.
// - grouping: kind/bundle/extension are read off the script NAME, since the wire
//   `type`/`bundle` fields are omitted when empty; an off-by-one in that path put
//   every client script in a group named after its bundle instead of its
//   extension, and left the loose scripts in a group called `undefined`.
// - badges/chips: a constraint with no `op` must render as the equality the
//   server assumes (`module = agent-contact`), not as `module agent-contact`;
//   and an iterator script has no `triggers` at all, so a view that reads the
//   trigger rows only showed it as a script that fires on nothing.
// - toggles: the pair narrows only while exactly one is on — both on used to
//   mean "no kind matches" and emptied the screen.

const CORREDOR_OFF = 'Corredor is turned off on this server'
const CORREDOR_UNREACHABLE = 'the server cannot reach it'

// A row the fixture extension must have deployed for this spec to mean anything.
const FIXTURE_SCRIPT = 'Greet contact (client)'

/** One extension group of the inventory, addressed by its heading.
 *  `c-resource-table` is CResourceTable's own root class — the group boundary
 *  has no role, and asserting inside the right group is the point here. */
function scriptGroup(page: Page, heading: string | RegExp): Locator {
  return page.locator('.c-resource-table').filter({ has: page.getByText(heading, { exact: true }) })
}

/** The inventory row for a script, by the label the script declares. */
function scriptRow(scope: Locator | Page, label: string): Locator {
  return scope.locator('tbody tr').filter({ hasText: label })
}

/** The filter switch of a labelled toggle. Label and control are not
 *  associated, so the switch is found inside the innermost element that holds
 *  both. */
function filterToggle(page: Page, label: RegExp): Locator {
  return page
    .locator('div')
    .filter({ has: page.getByText(label) })
    .filter({ has: page.getByRole('switch') })
    .last()
    .getByRole('switch')
}

/** Loads the inventory and says why it cannot be asserted on, empty when it can. */
async function openScriptList(page: Page): Promise<string> {
  await page.goto('/admin/automation/scripts')

  // The view renders a spinner until the list resolves; the filter form is the
  // part that is there in every banner state.
  await expect(page.getByText('Search query')).toBeVisible({ timeout: 30000 })

  if (await page.getByText(CORREDOR_OFF).count()) {
    return 'Corredor is off on this stack (CORREDOR_ENABLED) — bring the stack up with Corredor: dev/agent/worktree.sh up'
  }

  if (await page.getByText(CORREDOR_UNREACHABLE).count()) {
    return 'the server cannot reach Corredor on this stack — check the Corredor gRPC address (CORREDOR_ADDR)'
  }

  if ((await scriptRow(page, FIXTURE_SCRIPT).count()) === 0) {
    return 'the dev/fixtures/corredor extension is not deployed on this stack — the inventory carries no fixture scripts to assert on'
  }

  return ''
}

/** The count a toggle's label carries, e.g. 7 from `Server scripts (7)`. */
async function toggleCount(page: Page, label: RegExp): Promise<number> {
  const text = await page.getByText(label).first().innerText()
  return Number(text.match(/\((\d+)\)/)?.[1] ?? '-1')
}

test.describe('corredor script inventory', () => {
  test('the connected banner names the last refresh and Refresh keeps the rows', async ({
    page,
  }) => {
    const unavailable = await openScriptList(page)
    test.skip(!!unavailable, unavailable)

    // Relative time and the exact instant beside it — the view shows both so the
    // "4 seconds ago" is never the only thing an admin has to trust.
    await expect(page.getByText(/Connected · scripts refreshed .+ \(.+\)/)).toBeVisible()

    const group = scriptGroup(page, 'agent-sandbox')
    const before = await group.locator('tbody tr').count()
    expect(before).toBeGreaterThan(0)

    const refresh = page.getByRole('button', { name: 'Refresh' })
    await expect(refresh).toBeVisible()

    // Refresh is the view's only round-trip after mount, and the rows have to
    // survive it.
    const refetch = page.waitForResponse(r => r.url().includes('/system/automation/'))
    await refresh.click()
    await refetch

    await expect(scriptRow(group, FIXTURE_SCRIPT)).toHaveCount(1)
    await expect(group.locator('tbody tr')).toHaveCount(before)
    await expect(page.getByText(/Connected · scripts refreshed .+ \(.+\)/)).toBeVisible()
  })

  test('rows are grouped by the extension they came from', async ({ page }) => {
    const unavailable = await openScriptList(page)
    test.skip(!!unavailable, unavailable)

    const extension = scriptGroup(page, 'agent-sandbox')
    const loose = scriptGroup(page, 'Scripts outside any extension')

    await expect(extension).toHaveCount(1)
    await expect(loose).toHaveCount(1)

    // A client script and a server script of the same extension group together,
    // however differently their names are laid out.
    await expect(scriptRow(extension, 'Greet contact (client)')).toHaveCount(1)
    await expect(scriptRow(extension, 'Activate contact (server)')).toHaveCount(1)

    // A script sitting directly under its kind belongs to no extension.
    await expect(scriptRow(loose, 'System ping (server)')).toHaveCount(1)
    await expect(scriptRow(extension, 'System ping (server)')).toHaveCount(0)
    await expect(scriptRow(loose, 'Greet contact (client)')).toHaveCount(0)
  })

  test('a row names its kind, its triggers, its iterator and its security', async ({ page }) => {
    const unavailable = await openScriptList(page)
    test.skip(!!unavailable, unavailable)

    const group = scriptGroup(page, 'agent-sandbox')

    // Kind badges: where the script runs, and for a client script which bundle
    // the browser loads it in.
    await expect(
      scriptRow(group, 'Activate contact (server)').getByText('Server', { exact: true }),
    ).toBeVisible()
    await expect(
      scriptRow(group, 'Greet contact (client)').getByText('Client · compose', { exact: true }),
    ).toBeVisible()

    // Trigger chip plus one chip per constraint, the missing operator filled in.
    const beforeCreate = scriptRow(group, 'Contact: default company (server)')
    await expect(
      beforeCreate.getByText('beforeCreate · compose:record', { exact: true }),
    ).toBeVisible()
    await expect(beforeCreate.getByText('module = agent-contact', { exact: true })).toBeVisible()
    await expect(beforeCreate.getByText('namespace = agent-sandbox', { exact: true })).toBeVisible()

    // Security: who a deferred script runs as, which is the only thing that
    // explains its permissions.
    await expect(
      scriptRow(group, 'Task: heartbeat (server, interval)').getByText('run as agent@local.dev', {
        exact: true,
      }),
    ).toBeVisible()

    // An iterator script declares no triggers: the iterator chip is the whole
    // account of when it runs and what it walks.
    await expect(
      scriptRow(group, 'Contact: sweep dormant (server, iterator)').getByText(
        /^iterator · onInterval · compose:record · update · \* \* \* \* \*$/,
      ),
    ).toBeVisible()
  })

  test('the Server and Client toggles narrow the list', async ({ page }) => {
    const unavailable = await openScriptList(page)
    test.skip(!!unavailable, unavailable)

    const rows = page.locator('tbody tr')
    const serverScripts = await toggleCount(page, /^Server scripts \(\d+\)$/)
    const clientScripts = await toggleCount(page, /^Client scripts \(\d+\)$/)
    expect(serverScripts).toBeGreaterThan(0)
    expect(clientScripts).toBeGreaterThan(0)
    await expect(rows).toHaveCount(serverScripts + clientScripts)

    const server = filterToggle(page, /^Server scripts \(\d+\)$/)
    const client = filterToggle(page, /^Client scripts \(\d+\)$/)

    await server.click()
    await expect(server).toBeChecked()
    await expect(rows).toHaveCount(serverScripts)
    await expect(scriptRow(page, 'Greet contact (client)')).toHaveCount(0)
    await expect(scriptRow(page, 'Activate contact (server)')).toHaveCount(1)

    // Both on states no preference, so the pair stops narrowing.
    await client.click()
    await expect(rows).toHaveCount(serverScripts + clientScripts)
    await expect(scriptRow(page, 'Greet contact (client)')).toHaveCount(1)

    await server.click()
    await expect(rows).toHaveCount(clientScripts)
    await expect(scriptRow(page, 'Activate contact (server)')).toHaveCount(0)
    await expect(scriptRow(page, 'Greet contact (client)')).toHaveCount(1)
  })
})
