// Every guarded editor, audited the same way: open an existing record from its
// list, leave without touching anything, and assert it stays quiet — then type
// into it and assert it warns.
//
// Run: node dev/agent/drive.mjs dev/agent/checks/editor-guards.mjs
//      node dev/agent/drive.mjs dev/agent/checks/editor-guards.mjs --only "role"
//
// Both halves are the point. An editor that never warns loses work; one that
// warns on load — because it dirties itself resolving defaults, or captures its
// baseline before the data lands — trains people to click through the warning,
// which loses work just as effectively.
import { drive, check } from '../drive.mjs'

// Editors reached from a list. Only the list path is named: the check clicks
// the first row, so no resource ID is pinned here to go stale.
const EDITORS = [
  { name: 'user', list: '/admin/system/users' },
  { name: 'user group', list: '/admin/system/user-groups' },
  { name: 'role', list: '/admin/system/roles' },
  { name: 'application', list: '/admin/system/applications' },
  { name: 'auth client', list: '/admin/system/auth-clients' },
  { name: 'llm provider', list: '/admin/system/llm-providers' },
  { name: 'template', list: '/admin/system/templates' },
  { name: 'queue', list: '/admin/system/queues' },
  { name: 'api gateway route', list: '/admin/system/api-gateway' },
  { name: 'connection', list: '/admin/system/connections' },
  { name: 'data source', list: '/admin/system/data-sources' },
  { name: 'workflow', list: '/admin/automation/workflows' },
  { name: 'taq', list: '/admin/automation/taq' },
  { name: 'federation node', list: '/admin/federation/nodes' },
]

// Opens the first row of a list. Returns null when the list is empty — an
// absent fixture is a gap in the run, not a failing editor, and the check says
// which so an empty report is never read as a clean one.
async function openFirstRow(page, list) {
  await page.open(list)
  // An empty table still renders a row — its "no records" message — and clicking
  // that navigates nowhere, which reads as a broken editor rather than a gap.
  const rows = page.locator('.p-datatable-tbody > tr:not(.p-datatable-empty-message)')
  await rows
    .first()
    .waitFor({ state: 'visible', timeout: 20000 })
    .catch(() => {})
  if ((await rows.count()) === 0) return null

  await rows.first().click()
  await page.raw.waitForTimeout(2500)
  return page.path()
}

async function leaveViaBack(page) {
  page.clearDialog()
  const back = page.locator('[data-testid="editor-back"]').first()
  if (!(await back.count())) return { left: false, dialog: null }
  await back.click()
  await page.raw.waitForTimeout(2000)
  return { left: true, dialog: page.dialog() }
}

for (const editor of EDITORS) {
  drive(`${editor.name}: untouched does not warn`, async page => {
    const at = await openFirstRow(page, editor.list)
    if (!check(`${editor.name} list has a row to open`, !!at, 'list is empty')) return
    if (!check('landed in an editor', at !== editor.list, `still at ${at}`)) return

    const { left, dialog } = await leaveViaBack(page)
    if (!check('editor has a Back button', left, 'no [data-testid="editor-back"]')) return
    check('no confirm on a clean leave', !dialog, dialog || '')
  })

  drive(`${editor.name}: warns once edited`, async page => {
    const at = await openFirstRow(page, editor.list)
    if (!check(`${editor.name} list has a row to open`, !!at, 'list is empty')) return

    const input = page.locator('input[type="text"]:not([readonly])').first()
    if (!(await input.count())) {
      check('editor has a text input to type in', false, 'none found')
      return
    }
    // Typed, not filled: a fill() sets the value with no pointer or key event,
    // so it exercises a path no user can take.
    await input.click()
    await page.raw.keyboard.type('drive check edit')
    await page.raw.keyboard.press('Tab')
    await page.raw.waitForTimeout(1200)

    const { left, dialog } = await leaveViaBack(page)
    if (!check('editor has a Back button', left, 'no [data-testid="editor-back"]')) return
    check('leaving prompts a confirm', !!dialog, dialog || 'no dialog')
  })
}
