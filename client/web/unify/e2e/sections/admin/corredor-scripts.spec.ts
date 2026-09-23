import { expect, test, type Locator, type Page } from '@playwright/test'

// Contract under test:
// - views/automation/Script/Index.intent.md: the read-only Corredor inventory —
//   its banner (off / unreachable / connected with the last refresh + a Refresh
//   that re-fetches), one resource list over the whole fetched set with the
//   extension named per row, a kind badge, one chip per trigger with its
//   constraints, the iterator chip, the security chips, the search box over
//   label/name/description, and the filter popover whose Server/Client pair
//   narrows the list.
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
// - extension column: kind/bundle/extension are read off the script NAME, since
//   the wire `type`/`bundle` fields are omitted when empty; an off-by-one in
//   that path named a client script's row after its bundle instead of its
//   extension, and left the loose scripts saying `undefined`.
// - badges/chips: a constraint with no `op` must render as the equality the
//   server assumes (`module = agent-contact`), not as `module agent-contact`;
//   and an iterator script has no `triggers` at all, so a view that reads the
//   trigger rows only showed it as a script that fires on nothing.
// - search: the box filters the fetched set over label, name AND description —
//   a search that read the label alone found nothing for a word the row shows.
// - filters: the Server/Client pair narrows only while exactly one is on — both
//   on used to mean "no kind matches" and emptied the screen.

const CORREDOR_OFF = 'Corredor is turned off on this server'
const CORREDOR_UNREACHABLE = 'the server cannot reach it'

const SEARCH_PLACEHOLDER = 'Search scripts by label, name or description'
const FILTER_BUTTON = 'Filter the scripts'

// A row the fixture extension must have deployed for this spec to mean anything.
const FIXTURE_SCRIPT = 'Greet contact (client)'

/** The inventory row for a script, by the label the script declares. */
function scriptRow(page: Page, label: string): Locator {
  return page.locator('tbody tr').filter({ hasText: label })
}

/** A filter of the popover, by the label that carries its live count. */
function filterCheckbox(page: Page, label: RegExp): Locator {
  return page.getByRole('checkbox', { name: label })
}

/** Opens the filter popover; its controls are not in the DOM until it is. */
async function openFilters(page: Page): Promise<void> {
  await page.getByRole('button', { name: FILTER_BUTTON }).click()
  await expect(page.getByText(/^Server scripts \(\d+\)$/)).toBeVisible()
}

async function closeFilters(page: Page): Promise<void> {
  await page.keyboard.press('Escape')
  await expect(page.getByText(/^Server scripts \(\d+\)$/)).toBeHidden()
}

/** Loads the inventory and says why it cannot be asserted on, empty when it can. */
async function openScriptList(page: Page): Promise<string> {
  await page.goto('/admin/automation/scripts')

  // The list renders behind a loading mask until it resolves; the search box is
  // the part that is there in every banner state.
  await expect(page.getByPlaceholder(SEARCH_PLACEHOLDER)).toBeVisible({ timeout: 30000 })

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

/** The count a filter's label carries, e.g. 7 from `Server scripts (7)`. */
async function filterCount(page: Page, label: RegExp): Promise<number> {
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

    const rows = page.locator('tbody tr')
    const before = await rows.count()
    expect(before).toBeGreaterThan(0)

    const refresh = page.getByRole('button', { name: 'Refresh' })
    await expect(refresh).toBeVisible()

    // Refresh is the view's only round-trip after mount, and the rows have to
    // survive it.
    const refetch = page.waitForResponse(r => r.url().includes('/system/automation/'))
    await refresh.click()
    await refetch

    await expect(scriptRow(page, FIXTURE_SCRIPT)).toHaveCount(1)
    await expect(rows).toHaveCount(before)
    await expect(page.getByText(/Connected · scripts refreshed .+ \(.+\)/)).toBeVisible()
  })

  test('one list holds every extension, each row naming the one it came from', async ({ page }) => {
    const unavailable = await openScriptList(page)
    test.skip(!!unavailable, unavailable)

    // One table: the inventory is a single resource list, not a table per
    // extension.
    await expect(page.locator('table')).toHaveCount(1)

    // A client script and a server script of the same extension sit in the same
    // list, however differently their names are laid out, and each row names the
    // extension rather than being grouped under it.
    await expect(
      scriptRow(page, 'Greet contact (client)').getByText('agent-sandbox', { exact: true }),
    ).toBeVisible()
    await expect(
      scriptRow(page, 'Activate contact (server)').getByText('agent-sandbox', { exact: true }),
    ).toBeVisible()

    // A script sitting directly under its kind belongs to no extension.
    await expect(
      scriptRow(page, 'System ping (server)').getByText('Outside any extension', { exact: true }),
    ).toBeVisible()

    // And the filter can narrow the one list to either of them.
    await openFilters(page)
    const looseScripts = await filterCount(page, /^Outside any extension \(\d+\)$/)
    expect(looseScripts).toBeGreaterThan(0)
    await page.getByRole('radio', { name: /^Outside any extension \(\d+\)$/ }).check()
    await closeFilters(page)

    await expect(page.locator('tbody tr')).toHaveCount(looseScripts)
    await expect(scriptRow(page, 'System ping (server)')).toHaveCount(1)
    await expect(scriptRow(page, 'Activate contact (server)')).toHaveCount(0)
  })

  test('a row names its kind, its triggers, its iterator and its security', async ({ page }) => {
    const unavailable = await openScriptList(page)
    test.skip(!!unavailable, unavailable)

    // Kind badges: where the script runs, and for a client script which bundle
    // the browser loads it in.
    await expect(
      scriptRow(page, 'Activate contact (server)').getByText('Server', { exact: true }),
    ).toBeVisible()
    await expect(
      scriptRow(page, 'Greet contact (client)').getByText('Client · compose', { exact: true }),
    ).toBeVisible()

    // Trigger chip plus one chip per constraint, the missing operator filled in.
    const beforeCreate = scriptRow(page, 'Contact: default company (server)')
    await expect(
      beforeCreate.getByText('beforeCreate · compose:record', { exact: true }),
    ).toBeVisible()
    await expect(beforeCreate.getByText('module = agent-contact', { exact: true })).toBeVisible()
    await expect(beforeCreate.getByText('namespace = agent-sandbox', { exact: true })).toBeVisible()

    // Security: who a deferred script runs as, which is the only thing that
    // explains its permissions.
    await expect(
      scriptRow(page, 'Task: heartbeat (server, interval)').getByText('run as agent@local.dev', {
        exact: true,
      }),
    ).toBeVisible()

    // An iterator script declares no triggers: the iterator chip is the whole
    // account of when it runs and what it walks.
    await expect(
      scriptRow(page, 'Contact: sweep dormant (server, iterator)').getByText(
        /^iterator · onInterval · compose:record · update · \* \* \* \* \*$/,
      ),
    ).toBeVisible()
  })

  test('the search box narrows over label, name and description', async ({ page }) => {
    const unavailable = await openScriptList(page)
    test.skip(!!unavailable, unavailable)

    const rows = page.locator('tbody tr')
    const all = await rows.count()
    const search = page.getByPlaceholder(SEARCH_PLACEHOLDER)

    // The label of one row, which no other row carries.
    await search.fill('sweep dormant')
    await expect(rows).toHaveCount(1)
    await expect(scriptRow(page, 'Contact: sweep dormant (server, iterator)')).toHaveCount(1)

    // The machine name, which is shown beneath the label.
    await search.fill('SinkHello')
    await expect(rows).toHaveCount(1)
    await expect(scriptRow(page, 'Sink: hello (server)')).toHaveCount(1)

    // Words only the description carries.
    await search.fill('signed sink request')
    await expect(rows).toHaveCount(1)
    await expect(scriptRow(page, 'Sink: hello (server)')).toHaveCount(1)

    await search.fill('')
    await expect(rows).toHaveCount(all)
  })

  test('the Server and Client filters narrow the list', async ({ page }) => {
    const unavailable = await openScriptList(page)
    test.skip(!!unavailable, unavailable)

    const rows = page.locator('tbody tr')

    await openFilters(page)
    const serverScripts = await filterCount(page, /^Server scripts \(\d+\)$/)
    const clientScripts = await filterCount(page, /^Client scripts \(\d+\)$/)
    expect(serverScripts).toBeGreaterThan(0)
    expect(clientScripts).toBeGreaterThan(0)
    await expect(rows).toHaveCount(serverScripts + clientScripts)

    const server = filterCheckbox(page, /^Server scripts \(\d+\)$/)
    const client = filterCheckbox(page, /^Client scripts \(\d+\)$/)

    await server.check()
    await expect(server).toBeChecked()
    await expect(rows).toHaveCount(serverScripts)
    await expect(scriptRow(page, 'Greet contact (client)')).toHaveCount(0)
    await expect(scriptRow(page, 'Activate contact (server)')).toHaveCount(1)

    // Both on states no preference, so the pair stops narrowing.
    await client.check()
    await expect(rows).toHaveCount(serverScripts + clientScripts)
    await expect(scriptRow(page, 'Greet contact (client)')).toHaveCount(1)

    await server.uncheck()
    await expect(rows).toHaveCount(clientScripts)
    await expect(scriptRow(page, 'Activate contact (server)')).toHaveCount(0)
    await expect(scriptRow(page, 'Greet contact (client)')).toHaveCount(1)
  })
})
